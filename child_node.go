package lycron

import (
	"fmt"
	"strconv"
)

type ChildNode struct {
	ID string
	PID string
	PIDFile string
	IP string
	client *Client
	Endpoints   []string
	Username    string
	Password    string
	DialTimeout int64
	ReqTimeout  int
	OnWorkerStart func()
}

func (cn *ChildNode)Run(){
	cn.WorkerStart()
}

func (cn *ChildNode)WorkerStart(){
	if cn.OnWorkerStart == nil{
		cn.DefaultWorkerStart()
	}else{
		cn.OnWorkerStart()
	}
}

func (cn *ChildNode)DefaultWorkerStart(){

}

func (cn *ChildNode)SetClient(){
	ec := EtcdConfig{
		Endpoints: cn.Endpoints,
		Username  : cn.Username,
		Password   : cn.Password,
		DialTimeout : cn.DialTimeout,
		ReqTimeout : cn.ReqTimeout,
	}
	var err error
	cn.client,err = NewClient(ec)
	if err != nil{
		fmt.Println("etcd client init err")
	}
}

func (cn *ChildNode)register(){

}

func (cn *ChildNode)exist() (pid int, err error) {
	resp, err := cn.client.get(Node + cn.ID)
	if err != nil {
		return
	}
	//不存在，则返回
	if len(resp.Kvs) == 0 {
		return -1, nil
	}
	//如果存在却读取不出来，则删除并返回
	if pid, err = strconv.Atoi(string(resp.Kvs[0].Value)); err != nil {
		if _, err = cn.client.delete(Node + cn.ID); err != nil {
			return
		}
		return -1, nil
	}
}