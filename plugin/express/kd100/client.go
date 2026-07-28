package kd100

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/lsjhtang/kid/logger"
)

type Client struct {
	BaseUrl string
}

// DoPost 发起 POST 请求
func (cli *Client) DoPost(url string, params *url.Values, receiver interface{}) error {
	resp, err := http.PostForm(cli.BaseUrl+url, *params)
	if err != nil {
		return fmt.Errorf("请求失败：%s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("关闭响应失败：%s", err.Error())
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码：%d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败：%s", err.Error())
	}
	if receiver != nil {
		if err := json.Unmarshal(body, receiver); err != nil {
			return fmt.Errorf("解析响应失败：%s", err.Error())
		}
	}
	return nil
}
