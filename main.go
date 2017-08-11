package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/levigross/grequests"
	"github.com/tidwall/buntdb"
)

const Failed = `{"status": "fail"}`

type Account struct {
	AppID  string `json:"appid"`
	Secret string `json:"secret"`
}

// 定义Basic Auth的用户名和密码用来防止接口被恶意访问
var Auth = map[string]string{
	"user": "api",
	"pass": "wechat",
}

// 封装grequests的Get方法获取字符串内容
func getPage(urlPath string, params map[string]string) string {
	ro := &grequests.RequestOptions{
		Params: params,
	}

	res, _ := grequests.Get(urlPath, ro)
	return res.String()
}

// 检查request的Basic Auth用户名和密码
func isValidAuth(ctx echo.Context) bool {
	ctx.Response().Header().Set("WWW-Authenticate", `Basic realm="unixs.org"`)
	if uname, upass, ok := ctx.Request().BasicAuth(); ok {
		if Auth["user"] == uname && Auth["pass"] == upass {
			return ok
		}
	}
	return false
}

func main() {
	var err error

	env := &Env{}                          //初始化运行时环境
	env.At = &Token{}                      // 获取到的access_token映射到Token结构中
	env.Acc = make(map[string]string)      // 从配置文件读取到的AppID和Secret存放到env.Acc中
	env.GetAccounts("account.json")        // 读取配置文件中的微信公众号AppID和AppSecret
	env.DB, err = buntdb.Open("wechat.db") // 创建一个K/V数据库用来保存access_token
	if err != nil {
		log.Fatal(err)
	}
	defer env.DB.Close()

	e := echo.New()

	e.GET("/token", func(ctx echo.Context) error {
		if !isValidAuth(ctx) {
			return ctx.String(http.StatusUnauthorized, "401 Authorization Required")
		}

		appid := ctx.QueryParam("appid")
		if appid == "" {
			log.Println("ERROR: 没有提供AppID参数")
			return ctx.String(http.StatusNotFound, Failed)
		}

		if secret, isExist := env.Acc[appid]; isExist {
			var access_token string
			var record_time string

			var content struct {
				Status      string `json:"status"`
				AccessToken string `json:"access_token"`
			}

			// 查询数据库中是否已经存在这个AppID的access_token
			record_time = env.GetValue(appid, "timestamp")
			access_token = env.GetValue(appid, "access_token")

		GetToken:
			// 如果在Access Token数据库中不存在这个appid的token就重新获取
			if access_token == "" {
				tjs := env.At.Get(appid, secret)

				// 没获得access_token就返回Failed消息
				if tjs == "" {
					log.Println("ERROR: 没有获得access_token")
					return ctx.String(http.StatusNotFound, Failed)
				}

				//获取Token之后更新运行时环境，然后返回access_token
				env.UpdateTokens(appid)

				content.Status = "success"
				content.AccessToken = env.At.AccessToken
				return ctx.JSON(http.StatusOK, content)
			}
			goto CheckTime

		CheckTime:
			// 如果数据库中已经存在了Token，就检查过期时间，如果过期了就去GetToken获取
			curTime := time.Now().Unix()

			expire_time, _ := strconv.ParseInt(record_time, 10, 64)
			timeout, _ := strconv.ParseInt(env.GetValue(appid, "expires_in"), 10, 64)

			if curTime >= expire_time+timeout {
				goto GetToken
			}

			content.Status = "success"
			content.AccessToken = access_token
			return ctx.JSON(http.StatusOK, content)
		}
		log.Println("ERROR: AppID不存在")
		// 如果提交的appid不在配置文件中，就返回Failed消息
		return ctx.String(http.StatusNotFound, Failed)
	})

	e.Start(":8000")
}
