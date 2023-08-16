package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
)

var (
	cli *req.Client
)

func init() {
	cli = req.NewClient()
	cli.SetTimeout(time.Minute)
}

type ReqParams map[string]string

func GET(url string, pathParams, queryParams ReqParams) ([]byte, error) {
	r := cli.R()
	if len(url) > 0 {
		r.SetPathParams(pathParams)
	}
	if len(queryParams) > 0 {
		r.SetQueryParams(queryParams)
	}
	resp, err := r.
		SetHeader("Content-Type", "application/json").
		EnableDump().
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New("请求错误")
	}
	bodyBytes, err := resp.ToBytes()
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

type Body map[string]interface{}

func POST(uri string, body Body) (*gjson.Result, error) {
	r := cli.R()
	if len(body) > 0 {
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r.SetBody(bodyJson)
	}
	resp, err := r.
		SetHeader("Content-Type", "application/json").
		EnableDump().
		Post(fmt.Sprintf("%s%s", cli.BaseURL, uri))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New("请求错误")
	}
	bodyBytes, err := resp.ToBytes()
	if err != nil {
		return nil, err
	}
	result := gjson.ParseBytes(bodyBytes)
	return &result, nil
}
