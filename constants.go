package lycron


const (

	// 常规
	Mode = "debug"
	// web 相关
	CodeSuccess = 0
	CodeSystemError = 500 	//系统错误
	CodeParamError = 40000  //参数错误

	WebAdder = "127.0.0.1:8083"

	// etcd 相关
	DefaultLock  ="/%s/lock/"
    Node = "/ly/node/"
    Proc = "/ly/proc/"
    Cmd = "/ly/cmd/"
    Once =  "/ly/once/"
    CsCtl = "/ly/csctl/"
    Lock = "/ly/lock/"
    Group = "/ly/group/"
	Notifier =  "/ly/notifier/"
)

var codeText = map[int]string{
	CodeSuccess:     "ok",
	CodeSystemError: "the system had a bad cold",
	CodeParamError:  "invalid param",
}

var (
	defaultNotice Notice
	defaultStore Store
)

func CodeText(code int) string {
	if msg, ok := codeText[code]; ok {
		return msg
	}
	return ""
}