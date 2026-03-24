package project

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"test.com/project-api/api/rpc"
	"test.com/project-api/pkg/model"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/auth"
)

type HandlerAuth struct {
}

func (a *HandlerAuth) authList(c *gin.Context) {
	result := &common.Result{}
	organizationCode := c.GetString("organizationCode")
	var page = &model.Page{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthReqMessage{
		OrganizationCode: organizationCode,
		Page:             page.Page,
		PageSize:         page.PageSize,
	}
	response, err := rpc.AuthServiceClient.AuthList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var authList []*model.ProjectAuth
	copier.Copy(&authList, response.List)
	if authList == nil {
		authList = []*model.ProjectAuth{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total": response.Total,
		"list":  authList,
		"page":  page.Page,
	}))
}

func (a *HandlerAuth) apply(c *gin.Context) {
	result := &common.Result{}
	var req *model.ProjectAuthReq
	c.ShouldBind(&req)
	var nodes []string
	if req.Nodes != "" {
		json.Unmarshal([]byte(req.Nodes), &nodes)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthReqMessage{
		Action: req.Action,
		AuthId: req.Id,
		Nodes:  nodes,
	}
	applyResponse, err := rpc.AuthServiceClient.Apply(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var list []*model.ProjectNodeAuthTree
	copier.Copy(&list, applyResponse.List)
	var checkedList []string
	copier.Copy(&checkedList, applyResponse.CheckedList)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":        list,
		"checkedList": checkedList,
	}))
}

func (a *HandlerAuth) add(c *gin.Context) {
	result := &common.Result{}
	var req *model.ProjectAuthReq
	c.ShouldBind(&req)
	organizationCode := c.GetString("organizationCode")
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthReqMessage{
		OrganizationCode: organizationCode,
		MemberId:         memberId,
		Title:            req.Title,
		Desc:             req.Desc,
		Status:            int32(req.Status),
		Sort:             int32(req.Sort),
		Type:             req.Type,
	}
	_, err := rpc.AuthServiceClient.SaveAuth(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (a *HandlerAuth) edit(c *gin.Context) {
	result := &common.Result{}
	var req *model.ProjectAuthReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthReqMessage{
		AuthId: req.Id,
		Title:  req.Title,
		Desc:   req.Desc,
		Status: int32(req.Status),
		Sort:   int32(req.Sort),
		Type:   req.Type,
	}
	_, err := rpc.AuthServiceClient.UpdateAuth(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (a *HandlerAuth) forbid(c *gin.Context) {
	result := &common.Result{}
	var req struct {
		Id     int64 `form:"id" json:"id"`
		Status int32 `form:"status" json:"status"`
	}
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthReqMessage{
		AuthId: req.Id,
		Status: req.Status,
	}
	_, err := rpc.AuthServiceClient.UpdateAuth(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (a *HandlerAuth) resume(c *gin.Context) {
	result := &common.Result{}
	var req struct {
		Id     int64 `form:"id" json:"id"`
		Status int32 `form:"status" json:"status"`
	}
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthReqMessage{
		AuthId: req.Id,
		Status: req.Status,
	}
	_, err := rpc.AuthServiceClient.UpdateAuth(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (a *HandlerAuth) setDefault(c *gin.Context) {
	result := &common.Result{}
	var req struct {
		Id int64 `form:"id" json:"id"`
	}
	c.ShouldBind(&req)
	organizationCode := c.GetString("organizationCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.AuthServiceClient.SetDefaultAuth(ctx, &auth.AuthReqMessage{
		AuthId:           req.Id,
		OrganizationCode: organizationCode,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (a *HandlerAuth) del(c *gin.Context) {
	result := &common.Result{}
	var req struct {
		Id int64 `form:"id" json:"id"`
	}
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.AuthServiceClient.DeleteAuth(ctx, &auth.AuthReqMessage{
		AuthId: req.Id,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (a *HandlerAuth) GetAuthNodes(c *gin.Context) ([]string, error) {
	memberId := c.GetInt64("memberId")
	msg := &auth.AuthReqMessage{
		MemberId: memberId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := rpc.AuthServiceClient.AuthNodesByMemberId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		return nil, errs.NewError(errs.ErrorCode(code), msg)
	}
	return response.List, err
}

func NewAuth() *HandlerAuth {
	return &HandlerAuth{}
}
