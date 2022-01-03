package models

import (
	"git.code.oa.com/v/live/console/live_management_log/common/entity"
	"git.code.oa.com/v/live/store/live_ec_common/livecom"
)

//LiveManagementLogDTO 适配转换协议
type LiveManagementLogDTO struct {
	LiveManagementLogId  string
	FirstModule          string
	SecondModule         string
	ExternModule         string
	OperationBizSign     string
	OperationBizSignName string
	OperationType        string
	OperationBizName     string
	OperationAccount     string
	OperationAccountName string
	OperationAccountIp   string
	OperationRequest     string
	OperationResponse    string
	OperationOriginData  string
	ExteranData          string
	CreateTime           string
	UpdateTime           string
	RequestUrl           string
	OperationAdgroup     string
	OperationRecord      string
	PrimaryKeyId         string
	Pid                  string
}

//ManagementLogQueryListRequest 适配请求
type ManagementLogQueryListRequest struct {
	livecom.PaginationRequest
	Condition []*livecom.QueryListCondition
}

//ManagementLogQueryListResponse 适配响应
type ManagementLogQueryListResponse struct {
	LiveManagementLogDTOS []*LiveManagementLogDTO
	entity.PaginationResponse
}

// Enum value maps for OperationBizSignType.
var (
	OperationBizSignTypeName = map[string]string{
		"INVALID_SIGN_TYPE":                          "非法标识",
		"CHANNEL_PAGE_LIVE_CONTENT_SIGIN_TYPE":       "直播管理",
		"SIGN_TOTAL_LIVE_CONTENT_OPERATOR_SIGN_TYPE": "直播内容管理",
		"FIRST_EXTERN_SIGN_TYPE":                     "扩展标识1",
		"SECOND_EXTERN_SIGN_TYPE":                    "扩展标识2",
	}
)
