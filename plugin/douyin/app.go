package douyin

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lsjhtang/kid/logger"
)

type App struct {
	opts *AppOptions
}

// GetAssessToken 获取授权码
// 文档 https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-permission/get-access-token
func (app *App) GetAssessToken(code string) (*AccessTokenResponse, error) {
	var resp AccessTokenResponse
	if err := app.doPost("access_token", map[string]interface{}{
		"grant_type":    "authorization_code",
		"client_key":    app.opts.ClientKey,
		"client_secret": app.opts.ClientSecret,
		"code":          code,
	}, &resp); err != nil {
		return nil, err
	}
	if resp.Data.ErrorCode != 0 {
		return nil, fmt.Errorf("app get access token error: code: %d msg: %s",
			resp.Data.ErrorCode, resp.Data.Description)
	}
	return &resp, nil
}

// GetUserinfo 获取用户信息
// 文档 https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/get-account-open-info
func (app *App) GetUserinfo(accessToken, openid string) (resp *UserinfoResponse, err error) {
	if err = app.doPost("userinfo", map[string]interface{}{
		"access_token": accessToken,
		"open_id":      openid,
	}, resp); err != nil {
		return nil, err
	}
	if resp.Data.ErrorCode != 0 {
		return nil, fmt.Errorf("app get user info error: code: %d msg: %s",
			resp.Data.ErrorCode, resp.Data.Description)
	}
	return
}

func (app *App) doPost(uri string, body map[string]interface{}, respData interface{}) error {
	var payload string
	for k, v := range body {
		payload += fmt.Sprintf("%s=%s&", k, v)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", douyinUrl, uri), strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("close app response body error: %s", err.Error())
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(respBody, respData); err != nil {
		return err
	}
	return nil
}

// DecryptUserPhone 解密用户手机号
// 文档 https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/phone-number-decode-demo
func (app *App) DecryptUserPhone(encryptedPhone string) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(encryptedPhone)
	if err != nil {
		return "", err
	}
	key := []byte(app.opts.ClientSecret)
	iv := []byte(app.opts.ClientSecret)[:16]
	phone, err := app.aesDecrypt(decodeString, key, iv)
	if err != nil {
		return "", err
	}
	return string(phone), nil
}

func (app *App) aesDecrypt(cryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("iv length must equal block size")
	}
	encrypter := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cryptData))
	encrypter.CryptBlocks(origData, cryptData)
	return app.pkcs5UnPadding(origData)
}

func (app *App) pkcs5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	end := length - unpadding
	if end > length || end < 0 {
		return []byte{}, fmt.Errorf("pkcs5 unpadding error")
	}
	return origData[:end], nil
}

func NewApp(opts *AppOptions) *App {
	return &App{opts: opts}
}
