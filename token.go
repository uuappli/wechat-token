package main

import (
	"encoding/json"
)

const AccessTokenAPI = "https://api.weixin.qq.com/cgi-bin/token"

type Token struct {
	AccessToken string `json:"access_token"`
	Expire      int    `json:"expires_in"`
}

func (t *Token) JSON(at string) string {
	if err := json.Unmarshal([]byte(at), t); err != nil {
		return ""
	}

	return at
}

// 获取AppID的access_token
func (t *Token) Get(appid string, secret string) string {
	var args = map[string]string{
		"appid":      appid,
		"secret":     secret,
		"grant_type": "client_credential",
	}

	at := getPage(AccessTokenAPI, args)
	tjs := t.JSON(at)

	return tjs
}
