package qqapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type httpHelper struct {
	logger *zap.Logger
	client *http.Client
}

func newHttpHelper(logger *zap.Logger, client *http.Client) *httpHelper {
	return &httpHelper{logger: logger, client: client}
}

//http-GetPic
func (c *httpHelper) GetPic(url string) ([]byte, error) {
	resp, err := c.client.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

//http-Get
func (c *httpHelper) Get(url string, result qqResult) error {
	resp, err := c.client.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		c.logger.Error("http-get请求结果：", zap.String("rsp", string(body)))
		return err
	}

	if result.GetErr() != 0 {
		err = errors.New(result.GetErrMessage())
		//err = xerrors.NewCustomError("QQHttpGet", result.GetErrMessage())
		return err
	}

	return nil
}

//http-Json-Post
func (c *httpHelper) Post(url string, data interface{}, contentType string, result qqResult) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		c.logger.Error("http-post请求结果：", zap.Any("rsp", resp))
		return err
	}

	return nil
}

//http-FormatData-Post
func (c *httpHelper) FormatDataPost(url string, data *bytes.Buffer, contentType string, result qqResult) error {
	resp, err := c.client.Post(url, contentType, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		c.logger.Error("http-post请求结果：", zap.Any("rsp", resp))
		return err
	}

	return nil
}
