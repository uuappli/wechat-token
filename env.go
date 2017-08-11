package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/tidwall/buntdb"
)

type Env struct {
	Acc  map[string]string
	Conf string
	DB   *buntdb.DB
	At   *Token
}

// 读取配置文件中的appid和secret值到一个map中
func (e *Env) GetAccounts(file string) {
	accounts := make([]Account, 1)

	e.Conf = file
	if file == "" {
		e.Conf = "account.json"
	}

	if _, err := os.Stat(e.Conf); err != nil {
		os.Exit(1)
	}

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		os.Exit(1)
	}

	json.Unmarshal(raw, &accounts)

	for _, a := range accounts {
		e.Acc[a.AppID] = a.Secret
	}
	return
}

func (e *Env) GetValue(appid string, key string) string {
	var value string

	err := e.DB.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(appid + "_" + key)
		if err != nil {
			return err
		}
		value = v
		return nil
	})
	if err != nil {
		value = ""
	}

	return value
}

// 更新AppID上下文环境中的Access Token和到期时间
func (e *Env) UpdateTokens(appid string) {
	timestamp := time.Now().Unix()

	e.DB.Update(func(tx *buntdb.Tx) error {
		tx.Set(appid+"_timestamp", strconv.FormatInt(timestamp, 10), nil)
		tx.Set(appid+"_access_token", e.At.AccessToken, nil)
		tx.Set(appid+"_expires_in", strconv.Itoa(e.At.Expire), nil)
		return nil
	})
}
