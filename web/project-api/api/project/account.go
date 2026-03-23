package project

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"test.com/project-api/api/rpc"
	"test.com/project-api/pkg/model"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/account"
)

type HandlerAccount struct {
}

func (a *HandlerAccount) account(c *gin.Context) {
	result := &common.Result{}
	var req *model.AccountReq
	_ = c.ShouldBind(&req)
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		MemberId:         memberId,
		OrganizationCode: c.GetString("organizationCode"),
		Page:             int64(req.Page),
		PageSize:         int64(req.PageSize),
		SearchType:       int32(req.SearchType),
		DepartmentCode:   req.DepartmentCode,
	}
	response, err := rpc.AccountServiceClient.Account(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var list []*model.MemberAccount
	copier.Copy(&list, response.AccountList)
	if list == nil {
		list = []*model.MemberAccount{}
	}
	var authList []*model.ProjectAuth
	copier.Copy(&authList, response.AuthList)
	if authList == nil {
		authList = []*model.ProjectAuth{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total":    response.Total,
		"page":     req.Page,
		"list":     list,
		"authList": authList,
	}))
}

func (a *HandlerAccount) allList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		OrganizationCode: c.GetString("organizationCode"),
		SearchType:       -1,
	}
	response, err := rpc.AccountServiceClient.SaveAccount(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var list []*model.MemberAccount
	copier.Copy(&list, response.AccountList)
	if list == nil {
		list = []*model.MemberAccount{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}

func (a *HandlerAccount) forbid(c *gin.Context) {
	result := &common.Result{}
	accountCode := c.PostForm("accountCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		AccountCode: accountCode,
		Authorize:   "forbid",
	}
	_, err := rpc.AccountServiceClient.SaveAccount(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(nil))
}

func (a *HandlerAccount) resume(c *gin.Context) {
	result := &common.Result{}
	accountCode := c.PostForm("accountCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		AccountCode: accountCode,
		Authorize:   "resume",
	}
	_, err := rpc.AccountServiceClient.SaveAccount(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(nil))
}

func (a *HandlerAccount) add(c *gin.Context) {
	result := &common.Result{}
	var req struct {
		Name           string `json:"name" form:"name"`
		Mobile         string `json:"mobile" form:"mobile"`
		Email          string `json:"email" form:"email"`
		DepartmentCode string `json:"departmentCode" form:"departmentCode"`
		Authorize      string `json:"authorize" form:"authorize"`
		Avatar         string `json:"avatar" form:"avatar"`
	}
	_ = c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		OrganizationCode: c.GetString("organizationCode"),
		Name:             req.Name,
		Mobile:           req.Mobile,
		Email:            req.Email,
		DepartmentCode:   req.DepartmentCode,
		Authorize:        req.Authorize,
		Avatar:           req.Avatar,
	}
	_, err := rpc.AccountServiceClient.SaveAccount(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(nil))
}

func (a *HandlerAccount) edit(c *gin.Context) {
	result := &common.Result{}
	var req struct {
		Code           string `json:"code" form:"code"`
		Name           string `json:"name" form:"name"`
		Mobile         string `json:"mobile" form:"mobile"`
		Email          string `json:"email" form:"email"`
		DepartmentCode string `json:"departmentCode" form:"departmentCode"`
		Authorize      string `json:"authorize" form:"authorize"`
		Avatar         string `json:"avatar" form:"avatar"`
		Position       string `json:"position" form:"position"`
	}
	_ = c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		AccountCode:      req.Code,
		OrganizationCode: c.GetString("organizationCode"),
		Name:             req.Name,
		Mobile:           req.Mobile,
		Email:            req.Email,
		DepartmentCode:   req.DepartmentCode,
		Authorize:        req.Authorize,
		Avatar:           req.Avatar,
		Position:         req.Position,
	}
	_, err := rpc.AccountServiceClient.SaveAccount(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(nil))
}

func (a *HandlerAccount) auth(c *gin.Context) {
	result := &common.Result{}
	accountCode := c.PostForm("id")
	authorize := c.PostForm("auth")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		AccountCode: accountCode,
		Authorize:   authorize,
	}
	_, err := rpc.AccountServiceClient.SaveAccount(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(nil))
}

func (a *HandlerAccount) del(c *gin.Context) {
	result := &common.Result{}
	accountCode := c.PostForm("accountCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		AccountCode: accountCode,
		Authorize:   "del",
	}
	_, err := rpc.AccountServiceClient.SaveAccount(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(nil))
}

func (a *HandlerAccount) read(c *gin.Context) {
	result := &common.Result{}
	code := c.PostForm("code")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		AccountCode: code,
	}
	response, err := rpc.AccountServiceClient.SaveAccount(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var ma *model.MemberAccount
	if len(response.AccountList) > 0 {
		ma = &model.MemberAccount{}
		copier.Copy(ma, response.AccountList[0])
	}
	c.JSON(http.StatusOK, result.Success(ma))
}

func (a *HandlerAccount) syncDetail(c *gin.Context) {
	result := &common.Result{}
	// 同步详情：直接返回当前登录用户的 account 信息
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		MemberId:         memberId,
		OrganizationCode: c.GetString("organizationCode"),
		Page:             1,
		PageSize:         1,
		SearchType:       1,
	}
	response, err := rpc.AccountServiceClient.Account(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var ma *model.MemberAccount
	if len(response.AccountList) > 0 {
		ma = &model.MemberAccount{}
		copier.Copy(ma, response.AccountList[0])
	}
	c.JSON(http.StatusOK, result.Success(ma))
}

func (a *HandlerAccount) joinByInviteLink(c *gin.Context) {
	result := &common.Result{}
	// 邀请链接加入：前端传 inviteCode，这里简单返回成功（具体逻辑在 project_member 模块）
	c.JSON(http.StatusOK, result.Success(nil))
}

func NewAccount() *HandlerAccount {
	return &HandlerAccount{}
}
