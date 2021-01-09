package lycron

import "fmt"

type NameNode struct {
	WebAdder string
	Mode string
	Notice Notice
	Store Store
	client *Client
 	Endpoints   []string
	Username    string
	Password    string
	DialTimeout int64
	ReqTimeout  int
	OnMasterStart func()
}

func (nn *NameNode)Run(){

    nn.SetNotice()
    nn.SetStore()
    nn.SetClient()
	nn.MasterStart()
    //监听报警
    go nn.client.monitorNotice()
    //监控节点 下线相关
	go nn.client.monitorNodes()

	_ = startServer(nn.WebAdder,nn.Mode)
}

func(nn *NameNode)MasterStart(){
	if nn.OnMasterStart == nil{
		nn.DefaultMasterStart()
	}else{
		nn.OnMasterStart()
	}
}

func(nn *NameNode)DefaultMasterStart(){

}

func (nn *NameNode)SetNotice(){
	if nn.Notice == nil{
		defaultNotice = &Ding{}
	}else{
		defaultNotice = nn.Notice
	}
}

func (nn *NameNode)SetStore(){
	if nn.Store == nil{
		defaultStore = &FileStore{}
	}else{
		defaultStore = nn.Store
	}
}

func (nn *NameNode)SetClient(){
	ec := EtcdConfig{
		Endpoints: nn.Endpoints,
		Username  : nn.Username,
		Password   : nn.Password,
		DialTimeout : nn.DialTimeout,
		ReqTimeout : nn.ReqTimeout,
	}
	var err error
	nn.client,err = NewClient(ec)
	if err != nil{
		fmt.Println("etcd client init err")
	}
}









