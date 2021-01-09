package lycron

import (
	"bytes"
	"encoding/json"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
)

type Notice interface {
	Server()
	Send(msg *Message)
}

type Message struct {
	Subject string  //主题
	Body    string  //内容
	To      []string //给谁
}

type Ding struct {
	HttPApi string
}

type DingAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingMessage struct {
	MsgType string            `json:"msgtype"`
	Text    map[string]string `json:"text"`
	At      DingAt           `json:"at"`
}

func (d *Ding)Send(msg *Message){
	//=================================api形式的改用钉钉通知
	dingMessage:= &DingMessage{}
	dingMessage.MsgType = "text"
	dingMessage.Text = map[string]string{"content":msg.Body}
	dingMessage.At.AtMobiles =msg.To
	dingMessage.At.IsAtAll = false
	//==========================================
	//body, err := json.Marshal(msg)
	body, err := json.Marshal(dingMessage)
	if err != nil {
		log.Warnf("http api send msg[%+v] err: %s", msg, err.Error())
		return
	}

	req, err := http.NewRequest("POST", d.HttPApi, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Warnf("http api send msg[%+v] err: %s", msg, err.Error())
		return
	}
	defer resp.Body.Close()

	//fmt.Println("返回码",resp.StatusCode)

	if resp.StatusCode == 200 {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("http api send msg[%+v] err: %s", msg, err.Error())
		return
	}
	log.Warnf("http api send msg[%+v] err: %s", msg, string(data))
	return
}

func(d *Ding)Server(){

}
