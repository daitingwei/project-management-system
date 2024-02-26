package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"test.com/project-api/api/rpc"
	"test.com/project-api/pkg/model/user"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/user/login"
)

type HandlerUser struct {
}

func New() *HandlerUser {
	return &HandlerUser{}
}

func (*HandlerUser) getCaptcha(ctx *gin.Context) {
	result := &common.Result{}
	mobile := ctx.PostForm("mobile")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rsp, err := rpc.LoginServiceClient.GetCaptcha(c, &login.CaptchaMessage{Mobile: mobile})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(rsp.Code))
}

func (u *HandlerUser) register(c *gin.Context) {
	// 1.接收参数 参数模型
	result := &common.Result{}
	var req user.RegisterReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	// 2.校验参数 判断参数是否合法
	if err := req.Verify(); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, err.Error()))
		return
	}
	// 3.调用user grpc服务 获取响应
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &login.RegisterMessage{}
	err = copier.Copy(msg, req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy有误"))
		return
	}
	_, err = rpc.LoginServiceClient.Register(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	// 4.返回结果
	c.JSON(http.StatusOK, result.Success("注册成功"))
}

// GetIp 获取ip函数
func GetIp(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}

func (u *HandlerUser) login(c *gin.Context) {
	// 1.接收参数 参数模型
	result := &common.Result{}
	var req user.LoginReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	// 2.调用user grpc 完成登录
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &login.LoginMessage{}
	err = copier.Copy(msg, req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy有误"))
		return
	}
	msg.Ip = GetIp(c)
	loginRsp, err := rpc.LoginServiceClient.Login(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	rsp := &user.LoginRsp{}
	err = copier.Copy(rsp, loginRsp)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy有误"))
		return
	}
	// 4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (u *HandlerUser) myOrgList(c *gin.Context) {
	result := &common.Result{}
	memberIdStr, _ := c.Get("memberId")
	memberId := memberIdStr.(int64)
	resp, err := rpc.LoginServiceClient.MyOrgList(context.Background(), &login.UserMessage{MemId: memberId})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var orgs []*user.OrganizationList
	if resp.OrganizationList != nil {
		copier.Copy(&orgs, resp.OrganizationList)
	} else {
		orgs = []*user.OrganizationList{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  orgs,
		"total": len(orgs),
	}))
}

func (u *HandlerUser) saveOrganization(c *gin.Context) {
	result := &common.Result{}
	memberIdStr, _ := c.Get("memberId")
	memberId := memberIdStr.(int64)

	name := c.PostForm("name")
	address := c.PostForm("address")
	if name == "" {
		c.JSON(http.StatusOK, result.Fail(400, "组织名称不能为空"))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	org, err := rpc.LoginServiceClient.SaveOrganization(ctx, &login.OrganizationMessage{
		Name:     name,
		Address:  address,
		MemberId: memberId,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(org))
}

func (u *HandlerUser) editOrganization(c *gin.Context) {
	result := &common.Result{}
	orgCode := c.PostForm("organizationCode")
	name := c.PostForm("name")
	address := c.PostForm("address")
	if orgCode == "" {
		c.JSON(http.StatusOK, result.Fail(400, "组织code不能为空"))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	org, err := rpc.LoginServiceClient.UpdateOrganization(ctx, &login.OrganizationMessage{
		Code:    orgCode,
		Name:    name,
		Address: address,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(org))
}

func (u *HandlerUser) logout(c *gin.Context) {
	result := &common.Result{}
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusOK, result.Fail(401, "未登录"))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rsp, err := rpc.LoginServiceClient.Logout(ctx, &login.LogoutMessage{Token: token})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	if rsp.Success {
		c.JSON(http.StatusOK, result.Success("退出登录成功"))
	} else {
		c.JSON(http.StatusOK, result.Fail(500, "退出登录失败"))
	}
}
