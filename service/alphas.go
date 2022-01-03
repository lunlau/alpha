package service

import (
	"context"

	"git.code.oa.com/v/live/store/live_ec_common/rlog"

	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/ten_video_live/live_management_log"

	"git.code.oa.com/v/live/console/live_management_log/common/code"
	dao "git.code.oa.com/v/live/console/live_management_log/dao"
	"git.code.oa.com/v/live/console/live_management_log/models"
	"git.code.oa.com/v/live/store/live_ec_common/livecom"
)

//go:generate mockgen -destination ./../mock/service/
//mock_live_management_log_export_service.go -package service git.code.oa.com/video_app_live/live_management_log/
//service LiveManagementLogExportService
// LiveManagementLogService defines service
type LiveManagementLogExportService interface {
	// ManagementLogAdd 增加操作日志
	ManagementLogAdd(ctx context.Context, req *pb.ManagementLogAddRequest,
		rsp *pb.ManagementLogAddResponse) (err error)
	// ManagementLogQueryList 查询多条操作日志
	ManagementLogQueryList(ctx context.Context, req *pb.ManagementLogQueryListRequest,
		rsp *pb.ManagementLogQueryListResponse) (err error)
	// LogQueryByPid 日志查询
	LogQueryByPid(ctx context.Context,
		req *pb.LogQueryByPidRequest, rsp *pb.LogQueryByPidResponse) (err error)
}

//liveManagementLogExportServiceImpl 日志service实现
type liveManagementLogExportServiceImpl struct {
	lmDao dao.LiveManagementLogDao
}

// LogQueryByPid 日志查询
func (s *liveManagementLogExportServiceImpl) LogQueryByPid(ctx context.Context,
	req *pb.LogQueryByPidRequest, rsp *pb.LogQueryByPidResponse) error {
	var pageIndex int32 = 0
	var pageSize int32 = 0
	if req.PageSize == 0 {
		pageSize = code.DefaultPageSize
	} else {
		pageSize = req.PageIndex
	}
	if req.PageIndex == 0 {
		pageIndex = code.DefaultPageIndex
	} else {
		pageIndex = req.PageIndex
	}
	id := req.PageCtx
	if id == "" {
		var err error
		id, err = s.lmDao.ManagementLogQueryMaxId(ctx)
		if err != nil {
			rlog.RErrorContextf(ctx, "ManagementLogQueryMaxId get faild ")
			return err
		}
	}
	reqQuery := &pb.ManagementLogQueryListRequest{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Condition: []*pb.QueryListCondition{{
			ConditionKey:   "pid",
			ConditionValue: req.Pid,
			Operation:      "EQUAL",
		},
		},
	}
	reqQuery.PageCtx = id
	rspQuery := &pb.ManagementLogQueryListResponse{}
	errQuery := s.ManagementLogQueryList(ctx, reqQuery, rspQuery)
	if errQuery != nil {
		rlog.RErrorContextf(ctx, "ManagementLogQueryList get faild  ")
		return errQuery
	}
	rsp.LiveManagementLogInfo = make([]*pb.LiveManagementLogInfo, 0)
	for _, item := range rspQuery.LiveManagementLogInfo {
		rsp.LiveManagementLogInfo = append(rsp.LiveManagementLogInfo, item)
	}
	rsp.Total = rspQuery.Total
	rsp.PageSize = rspQuery.PageSize
	return nil
}

// GetLiveManagementLogExportService GetLiveManagementLogExportService
func GetLiveManagementLogExportService() LiveManagementLogExportService {
	lm := liveManagementLogExportServiceImpl{}
	lm.lmDao = dao.GetLiveManagementLogDao()
	return &lm
}

// checkManagementLogAddRequest checkManagementLogAddRequest
func checkManagementLogAddRequest(req *pb.ManagementLogAddRequest) error {
	// check
	if req == nil || req.OperationBizSign == pb.OperationBizSignType_INVALID_SIGN_TYPE ||
		req.OperationType == pb.OperationType_INVALID_TYPE {
		return errs.New(code.BIZ_SIGN_IS_INVALID, "invalid OperationBizSign or OperationType")
	}
	return nil
}

// ReportInfo ReportInfo
type ReportInfo struct {
	OperatorName    string
	RpcType         string //函数类型
	ReportType      string // 操作类型
	OperatType      string // 操作，增加修改，删除
	ManagementScene string //页面key
}

// management_scene 页面key
func recordLog(ctx context.Context, reportInfo ReportInfo) {
	log.WithContextFields(ctx,
		"operate_name", reportInfo.OperatorName, // AttaID扩展字段
		"rpc_type", reportInfo.RpcType, //
		"report_type", reportInfo.ReportType, //
		"operat_type", reportInfo.OperatType, //
		"management_scene", reportInfo.ManagementScene, //
	)
}

// ManagementLogAdd 添加日志功能
func (s *liveManagementLogExportServiceImpl) ManagementLogAdd(ctx context.Context, req *pb.ManagementLogAddRequest,
	rsp *pb.ManagementLogAddResponse) error {
	//check参数
	_ = checkManagementLogAddRequest(req)
	// sql encode
	lm := models.LiveManagementLogDTO{}
	err := livecom.SimpleCopyProperties(&lm, req)
	//数据类型不一致, copy不到重新赋值, 重新赋值
	lm.OperationBizSign = pb.OperationBizSignType_name[int32(req.OperationBizSign)]
	lm.OperationType = pb.OperationType_name[int32(req.OperationType)]
	// 补充参数
	recordLog(ctx, ReportInfo{
		OperatorName:    req.OperationAccount,
		RpcType:         req.OperationBizName,
		ReportType:      pb.OperationType_name[int32(req.OperationType)],
		ManagementScene: req.OperationBizName,
	})
	//执行业务
	rows, err := s.lmDao.ManagementLogAdd(ctx, lm)
	if err != nil {
		log.ErrorContextf(ctx, "ManagementLogAdd  result:%+v", err)
		return nil
	}
	log.InfoContextf(ctx, "GetProgramList success, rows : %v", rows)
	return nil
}

// ManagementLogQueryList 查询日志列表service
func (s *liveManagementLogExportServiceImpl) ManagementLogQueryList(ctx context.Context,
	req *pb.ManagementLogQueryListRequest, rsp *pb.ManagementLogQueryListResponse) error {
	if len(req.Condition) == 0 {
		rsp.ErrMsg = "query condition not set"
		return nil
	}
	var queryRequest = models.ManagementLogQueryListRequest{}
	if req.PageSize > 0 && req.PageSize <= 200 {
		queryRequest.PageSize = req.PageSize
	} else {
		queryRequest.PageSize = 20
	}
	for _, conItem := range req.Condition {
		var queryCondition = livecom.QueryListCondition{}
		queryCondition.Operation = conItem.Operation
		queryCondition.ConditionValue = conItem.ConditionValue
		queryCondition.ConditionKey = conItem.ConditionKey
		queryRequest.Condition = append(queryRequest.Condition, &queryCondition)
	}
	// 插入最大id
	id := req.PageCtx
	if id == "" {
		var err error
		id, err = s.lmDao.ManagementLogQueryMaxId(ctx)
		if err != nil {
			rlog.RErrorContextf(ctx, "ManagementLogQueryMaxId get faild ")
			return err
		}
	}
	var queryCondition = livecom.QueryListCondition{}
	queryCondition.Operation = "EQUAL"
	queryCondition.ConditionValue = id
	queryCondition.ConditionKey = "live_management_log_id"
	queryRequest.Condition = append(queryRequest.Condition, &queryCondition)

	var queryRespons = models.ManagementLogQueryListResponse{}
	//处理请求
	err := s.lmDao.ManagementLogQueryList(ctx, &queryRequest, &queryRespons)
	if err != nil {
		log.Errorf("ManagementLogQueryList failed, %v", err)
		return nil
	}
	//组装回包
	rsp.PageSize = queryRespons.PageSize
	rsp.Total = queryRespons.Total
	for _, conItem := range queryRespons.LiveManagementLogDTOS {
		var info = pb.LiveManagementLogInfo{}
		err := livecom.SimpleCopyProperties(&info, conItem)
		//数据类型不一致,转换一下
		info.OperationBizSign = pb.OperationBizSignType(pb.OperationBizSignType_value[conItem.OperationBizSign])
		info.OperationType = pb.OperationType(pb.OperationType_value[conItem.OperationType])
		if err != nil {
			log.Errorf("copy request failed")
			return nil
		}
		rsp.LiveManagementLogInfo = append(rsp.LiveManagementLogInfo, &info)
	}
	return nil
}
