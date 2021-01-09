package lycron

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/prometheus/common/log"
	"time"
	"context"
)

type Client struct {
	*clientv3.Client
	reqTimeout time.Duration
}

type EtcdConfig struct {
	Endpoints   []string //["http://127.0.0.1:2379"]
	Username    string // ""
	Password    string // ""
	DialTimeout int64 // 2 单位秒
	ReqTimeout  int   //// etcd客户端的,请求超时时间，单位秒  之所以放在外层，是因为是通过控制上下文时间来达成效果的
	conf clientv3.Config
}

func (e EtcdConfig) Copy() clientv3.Config {
	e.conf.Endpoints=e.Endpoints
	e.conf.Username=e.Username
	e.conf.Password=e.Password
	if e.DialTimeout > 0{
		e.conf.DialTimeout= time.Duration(e.DialTimeout) * time.Second
	}
	return e.conf
}


func NewClient(etcdConf  EtcdConfig) (c *Client, err error) {
	cli, err := clientv3.New(etcdConf.Copy())
	if err != nil {
		return
	}
	c = &Client{
		Client: cli,
		reqTimeout: time.Duration(etcdConf.ReqTimeout) * time.Second,
	}
	return
}

func (c *Client) Watch(key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return c.Client.Watch(context.Background(), key, opts...)
}


func (c *Client)monitorNodes(){
	//watch node kv 变化
	rch :=  c.Watch(Node, clientv3.WithPrefix())
	for re := range rch {
		for _,ev := range re.Events{
			switch ev.Type{
			case clientv3.EventTypeDelete:
				fmt.Println("掉线了一个。。。。")
			}
		}
	}
}

func (c *Client)monitorNotice(){
	rch := c.Watch(Notifier, clientv3.WithPrefix())
	var err error
	for resp := range rch{
		for _, ev := range resp.Events {
			switch {
			case ev.IsCreate(), ev.IsModify():
				msg := new(Message)
				if err = json.Unmarshal(ev.Kv.Value, msg); err != nil {
					log.Warnf("msg[%s] umarshal err: %s", string(ev.Kv.Value), err.Error())
					continue
				}
				defaultNotice.Server()
				defaultNotice.Send(msg)
			}
		}
	}
}

func (c *Client) put(key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	ctx, cancel := NewETimeoutContext(c)
	defer cancel()
	return c.Client.Put(ctx, key, val, opts...)
}

func (c *Client)get(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := NewETimeoutContext(c)
	defer cancel()
	return c.Client.Get(ctx, key, opts...)
}

func (c *Client) delete(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	ctx, cancel := NewETimeoutContext(c)
	defer cancel()
	return c.Client.Delete(ctx, key, opts...)
}

type eTimeoutContext struct {
	context.Context
	eEndpoints []string
}

func NewETimeoutContext(c *Client) (context.Context, context.CancelFunc) {
	//设置一个超时的context
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	eCtx := &eTimeoutContext{}
	eCtx.Context = ctx
	eCtx.eEndpoints = c.Endpoints()
	return eCtx, cancel
}
