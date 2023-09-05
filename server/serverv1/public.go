package serverv1

import (
	"context"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"

	"github.com/busyfree/leaf-go/rpc/common"
	"github.com/busyfree/leaf-go/util/conf"
)

type Public struct{}

func (s *Public) Segment(ctx context.Context, req *common.SegmentKeyReq) (*common.Result, error) {
	var (
		resp = &common.Result{Id: 0, Status: common.Status_Status_Exception, Msg: "error"}
	)
	qpsMaps := conf.GetStrMapStr("SENTINEL_RES_QPS")
	if _, ok := qpsMaps["api_segment"]; ok {
		// Entry 方法用于埋点
		option := sentinel.WithTrafficType(base.Inbound)
		e, b := sentinel.Entry("api_segment", option)
		if b != nil {
			resp.Msg = "服务超载"
			return resp, nil
		}
		defer e.Exit()
	}

	key := req.GetKey()
	if len(key) == 0 {
		resp.Msg = "missing key"
		return resp, nil
	}
	r := segmentService.Get(ctx, key)
	resp.Id = r.Id
	resp.Status = common.Status_Status_Success
	resp.Msg = "ok"
	return resp, nil
}

func (s *Public) Snowflake(ctx context.Context, req *common.SegmentKeyReq) (*common.Result, error) {
	var (
		resp = &common.Result{Id: 0, Status: common.Status_Status_Exception}
	)
	qpsMaps := conf.GetStrMapStr("SENTINEL_RES_QPS")
	if _, ok := qpsMaps["api_snowflake"]; ok {
		// Entry 方法用于埋点
		option := sentinel.WithTrafficType(base.Inbound)
		e, b := sentinel.Entry("honda_accord_api_wxlogin", option)
		if b != nil {
			resp.Msg = "服务超载"
			return resp, nil
		}
		defer e.Exit()
	}
	r := snowflakeService.Get(ctx, req.GetKey())
	resp.Id = r.Id
	resp.Status = common.Status_Status_Success
	resp.Msg = "ok"
	return resp, nil
}
