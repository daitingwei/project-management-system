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
	"test.com/project-grpc/department"
)

type HandlerDepartmentMember struct {
}

func (d *HandlerDepartmentMember) list(c *gin.Context) {
	result := &common.Result{}
	var req model.DepartmentMemberReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.DepartmentServiceClient.DepartmentMemberList(ctx, &department.DepartmentReqMessage{
		DepartmentCode:   req.DepartmentCode,
		OrganizationCode: c.GetString("organizationCode"),
		Page:             req.Page,
		PageSize:         req.PageSize,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var list []*model.DepartmentMember
	copier.Copy(&list, resp.List)
	if list == nil {
		list = []*model.DepartmentMember{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  list,
		"total": resp.Total,
	}))
}

func (d *HandlerDepartmentMember) searchInviteMember(c *gin.Context) {
	result := &common.Result{}
	keyword := c.PostForm("keyword")
	departmentCode := c.PostForm("departmentCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.DepartmentServiceClient.SearchDepartmentMember(ctx, &department.DepartmentReqMessage{
		DepartmentCode:   departmentCode,
		OrganizationCode: c.GetString("organizationCode"),
		Keyword:          keyword,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var list []*model.DepartmentMember
	copier.Copy(&list, resp.List)
	if list == nil {
		list = []*model.DepartmentMember{}
	}
	// 前端直接用 res.data 作为数组
	c.JSON(http.StatusOK, result.Success(list))
}

func (d *HandlerDepartmentMember) inviteMember(c *gin.Context) {
	result := &common.Result{}
	departmentCode := c.PostForm("departmentCode")
	accountCode := c.PostForm("accountCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.DepartmentServiceClient.InviteDepartmentMember(ctx, &department.DepartmentReqMessage{
		DepartmentCode:   departmentCode,
		OrganizationCode: c.GetString("organizationCode"),
		AccountCode:      accountCode,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (d *HandlerDepartmentMember) removeMember(c *gin.Context) {
	result := &common.Result{}
	departmentCode := c.PostForm("departmentCode")
	accountCode := c.PostForm("accountCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.DepartmentServiceClient.RemoveDepartmentMember(ctx, &department.DepartmentReqMessage{
		DepartmentCode:   departmentCode,
		OrganizationCode: c.GetString("organizationCode"),
		AccountCode:      accountCode,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (d *HandlerDepartmentMember) detail(c *gin.Context) {
	result := &common.Result{}
	// 前端传 {code: memberCode, organization: orgCode}
	accountCode := c.PostForm("code")
	departmentCode := c.PostForm("departmentCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.DepartmentServiceClient.DepartmentMemberDetail(ctx, &department.DepartmentReqMessage{
		AccountCode:      accountCode,
		DepartmentCode:   departmentCode,
		OrganizationCode: c.GetString("organizationCode"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var res model.DepartmentMember
	if resp.Member != nil {
		copier.Copy(&res, resp.Member)
	}
	c.JSON(http.StatusOK, result.Success(res))
}

func NewDepartmentMember() *HandlerDepartmentMember {
	return &HandlerDepartmentMember{}
}
