package code

const (
	// 每次请求的列表的最大数量限制
	MaxNum = 20
	// 解析上下文错误
	MYSQL_CLIENT_NAME_ID     = "trpc.mysql.xxx.xxx"
	LiveOperatorLogTableName = "t_live_management_log"
)
const (
	LiveLogMysqlErr = "LiveLogMysqlErr"

	DefaultPageIndex = 1
	DefaultPageSize  = 1000
)

const (
	UN_DELETED int = 1
	IS_DELETED int = 2
)

// codeText 定义错误码的映射关系
var codeText = map[int]string{
	0: "success",
	// 系统功能错误错误
	542001: " generate id failed",
}

// 定义错误码
var (
	SUCCESS int = 0
	//非法处理错误码
	PAR_ERRPARAM int = 540001

	//系统错误码
	SYS_GENERATE_ID_FAILED    int = 542001
	INVALID_KEY_AND_OPERATION int = 542001

	//日志业务错误码544101 ～ 544199
	BIZ_NAME_IS_EMPTY   int = 544101
	BIZ_SIGN_IS_INVALID int = 544102
)
