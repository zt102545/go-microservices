package consts

const (
	Info_OK = 200
)

const (
	Error_Normal           = 5000 //常规错误
	Error_Unknown          = 5001 // 未知异常
	Error_LackParam        = 5002 // 缺少参数
	Error_DBError          = 5003 // 数据库操作异常
	Error_PURCHASE_FAILURE = 5004
)

var CodeMap = map[int64]string{
	200:  "success",
	5000: "fail",
	5001: "fail",
	5002: "fail",
	5003: "fail",
	5004: "fail",
}
