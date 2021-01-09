package lycron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
)

func startServer(Adder string,RunMode string) error{
	if Adder == ""{
		Adder = WebAdder
	}
	if RunMode == ""{
		RunMode = Mode
	}
	gin.SetMode(RunMode)
	engine := gin.New()
	registerRoute(engine)
	err:= runEngine(engine, Adder)
	if err !=nil{
		fmt.Println(err)
		return err
	}
	return nil
}

// TODO 添加自定义 路由 + 中间件
func registerRoute(router *gin.Engine){
	router.GET("/ping", Ping)
}


func runEngine(engine *gin.Engine, addr string) error {
	server := endless.NewServer(addr, engine)
	err := server.ListenAndServe()
	return err
}

func suc(c *gin.Context, data interface{}) {
	if IsNil(data){
		data=make([]string,0)
	}
	c.JSON(http.StatusOK, gin.H{
		"errcode":     CodeSuccess,
		"errmsg":      CodeText(CodeSuccess),
		"data":        data,
	})
	c.Abort()
}

func err(c *gin.Context, code int, msg ...string) {
	//获取错误信息
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = CodeText(code)
		if len(message)<=0{
			message=http.StatusText(code)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"errcode":        code,
		"errmsg":         message,
		"request_uri": c.Request.URL.Path,
	})
	c.Abort()
}

func genRequest(c *gin.Context, request interface{}) (err error) {
	//读取客户端传递的body
	body, err := readBody(c)
	if err != nil {
		return
	}
	//json解码
	err = json.Unmarshal(body, request)
	if err == nil {
		validate := validator.New()
		errValidate := validate.Struct(request)
		if errValidate != nil {
			//记录验证错误的参数
			return errValidate
		}
	}else{
		return  err
	}
	return nil
}

func readBody(c *gin.Context) (body []byte, err error) {
	// ReadAll 读取 r 中的所有数据，返回读取的数据和遇到的错误。
	//如果读取成功，则 err 返回 nil，而不是 EOF，因为 ReadAll 定义为读取所有数据，所以不会把 EOF 当做错误处理。
	body, err = ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	//包装一个io.Reader，返回一个io.ReadCloser(读出来之后需要再归还额)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return
}

func Ping(ctx *gin.Context) {
	suc(ctx, "pong")
	return
}