package project

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"test.com/project-api/api/rpc"
	"test.com/project-api/pkg/model"
	"test.com/project-api/pkg/model/pro"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/project"
	task_service_v1 "test.com/project-grpc/task"
)

type HandlerProject struct {
}

func (p *HandlerProject) index(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.IndexMessage{}
	indexResponse, err := rpc.ProjectServiceClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	menus := indexResponse.Menus
	var ms []*model.Menu
	copier.Copy(&ms, menus)
	c.JSON(http.StatusOK, result.Success(ms))
}

func (p *HandlerProject) myProjectList(c *gin.Context) {
	result := &common.Result{}
	// 1. 获取参数
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	page := &model.Page{}
	page.Bind(c)
	selectBy := c.PostForm("selectBy")
	// 兼容前端传参：recycle=1 → deleted，archive=1 → archive
	if selectBy == "" && c.PostForm("recycle") == "1" {
		selectBy = "deleted"
	}
	if selectBy == "" && c.PostForm("archive") == "1" {
		selectBy = "archive"
	}
	msg := &project.ProjectRpcMessage{
		MemberId:   memberId,
		MemberName: memberName,
		SelectBy:   selectBy,
		Page:       page.Page,
		PageSize:   page.PageSize}
	myProjectResponse, err := rpc.ProjectServiceClient.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 调试日志
	zap.L().Info("myProjectList debug", zap.Any("Pm", myProjectResponse.Pm), zap.Any("Total", myProjectResponse.Total))

	// 用 JSON 中间转换解决 copier 字段名不匹配问题
	jsonBytes, _ := json.Marshal(myProjectResponse.Pm)
	var pms []*pro.ProjectAndMember
	json.Unmarshal(jsonBytes, &pms)
	if pms == nil {
		pms = []*pro.ProjectAndMember{}
	}
	// 补齐 deleted_time（JSON 中间转换时 key 大小写不匹配导致丢失）
	for i, v := range myProjectResponse.Pm {
		if i < len(pms) {
			pms[i].DeletedTime = v.DeletedTime
		}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": myProjectResponse.Total,
	}))
}

func (p *HandlerProject) projectTemplate(c *gin.Context) {
	result := &common.Result{}
	// 1. 获取参数
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	page := &model.Page{}
	page.Bind(c)
	viewTypeStr := c.PostForm("viewType")
	viewType, _ := strconv.ParseInt(viewTypeStr, 10, 64)
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		MemberName:       memberName,
		ViewType:         int32(viewType),
		Page:             page.Page,
		PageSize:         page.PageSize,
		OrganizationCode: c.GetString("organizationCode")}
	templateResponse, err := rpc.ProjectServiceClient.FindProjectTemplate(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	var pms []*pro.ProjectTemplate
	copier.Copy(&pms, templateResponse.Ptm)
	if pms == nil {
		pms = []*pro.ProjectTemplate{}
	}
	for _, v := range pms {
		if v.TaskStages == nil {
			v.TaskStages = []*pro.TaskStagesOnlyName{}
		}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms, // null nil -> []
		"total": templateResponse.Total,
	}))
}

func (p *HandlerProject) projectTemplateSave(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	organizationCode := c.GetString("organizationCode")
	name := c.PostForm("name")
	description := c.PostForm("description")
	cover := c.PostForm("cover")
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		OrganizationCode: organizationCode,
		Name:              name,
		Description:       description,
		Cover:             cover,
	}
	template, err := rpc.ProjectServiceClient.SaveProjectTemplate(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(template))
}

func (p *HandlerProject) projectTemplateEdit(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	code := c.PostForm("code")
	name := c.PostForm("name")
	description := c.PostForm("description")
	cover := c.PostForm("cover")
	msg := &project.ProjectRpcMessage{
		Code:        code,
		Name:        name,
		Description: description,
		Cover:       cover,
	}
	rsp, err := rpc.ProjectServiceClient.UpdateProjectTemplate(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Ok))
}

func (p *HandlerProject) projectTemplateDelete(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	code := c.PostForm("code")
	msg := &project.ProjectRpcMessage{
		Code: code,
	}
	rsp, err := rpc.ProjectServiceClient.DeleteProjectTemplate(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Ok))
}

func (p *HandlerProject) projectSave(c *gin.Context) {
	result := &common.Result{}
	// 1. 获取参数
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	organizationCode := c.GetString("organizationCode")
	req := &pro.SaveProjectRequest{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, result.Fail(-1, "参数错误"))
		return
	}
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		OrganizationCode: organizationCode,
		TemplateCode: req.TemplateCode,
		Name:         req.Name,
		Id:           int64(req.Id),
		Description:  req.Description,
	}
	saveProject, err := rpc.ProjectServiceClient.SaveProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	rsp := &pro.SaveProject{}
	copier.Copy(rsp, saveProject)
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (p *HandlerProject) readProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	detail, err := rpc.ProjectServiceClient.FindProjectDetail(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, MemberId: memberId})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	pd := &pro.ProjectDetail{}
	copier.Copy(pd, detail)
	c.JSON(http.StatusOK, result.Success(pd))
}

func (p *HandlerProject) recycleProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.UpdateDeletedProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, Deleted: true})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) recoveryProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.UpdateDeletedProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, Deleted: false})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) collectProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	collectType := c.PostForm("type")
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.UpdateCollectProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, CollectType: collectType, MemberId: memberId})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) editProject(c *gin.Context) {
	result := &common.Result{}
	var req *pro.ProjectReq
	_ = c.ShouldBind(&req)
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.UpdateProjectMessage{}
	copier.Copy(msg, req)
	msg.MemberId = memberId
	_, err := rpc.ProjectServiceClient.UpdateProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) getLogBySelfProject(c *gin.Context) {
	result := &common.Result{}
	var page = &model.Page{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.ProjectRpcMessage{
		MemberId: c.GetInt64("memberId"),
		Page:     page.Page,
		PageSize: page.PageSize,
	}
	projectLogResponse, err := rpc.ProjectServiceClient.GetLogBySelfProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var list []*model.ProjectLog
	copier.Copy(&list, projectLogResponse.List)
	if list == nil {
		list = []*model.ProjectLog{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}

func (p *HandlerProject) nodeList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := rpc.ProjectServiceClient.NodeList(ctx, &project.ProjectRpcMessage{})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var list []*model.ProjectNodeTree
	copier.Copy(&list, response.Nodes)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"nodes": list,
	}))
}

func (p *HandlerProject) FindProjectByMemberId(memberId int64, projectCode string, taskCode string) (*pro.Project, bool, bool, *errs.BError) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.ProjectRpcMessage{
		MemberId:    memberId,
		ProjectCode: projectCode,
		TaskCode:    taskCode,
	}
	projectResponse, err := rpc.ProjectServiceClient.FindProjectByMemberId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		return nil, false, false, errs.NewError(errs.ErrorCode(code), msg)
	}
	if projectResponse.Project == nil {
		return nil, false, false, nil
	}
	pr := &pro.Project{}
	copier.Copy(pr, projectResponse.Project)
	return pr, true, projectResponse.IsOwner, nil
}

func (p *HandlerProject) quit(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.QuitProject(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		MemberId:    c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) archive(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.ArchiveProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) recoveryArchive(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.RecoveryArchive(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		MemberId:    c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) deleteProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.DeleteProject(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		MemberId:    c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) analysis(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.ProjectServiceClient.ProjectAnalysis(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		MemberId:    c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp))
}

func (p *HandlerProject) searchInviteMember(c *gin.Context) {
	result := &common.Result{}
	keyword := c.PostForm("keyword")
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.ProjectServiceClient.SearchInviteMember(ctx, &project.ProjectRpcMessage{
		Name:             keyword,
		ProjectCode:      projectCode,
		MemberId:         c.GetInt64("memberId"),
		OrganizationCode: c.GetString("organizationCode"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	list := resp.List
	if list == nil {
		list = []*project.ProjectMemberMessage{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}

func (p *HandlerProject) inviteMember(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	memberCode := c.PostForm("memberCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.InviteMember(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		MemberCode:  memberCode,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (p *HandlerProject) joinByInviteLink(c *gin.Context) {
	result := &common.Result{}
	inviteCode := c.PostForm("inviteCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.JoinByInviteLink(ctx, &project.ProjectRpcMessage{
		InviteCode: inviteCode,
		MemberId:   c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (p *HandlerProject) removeMember(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	memberCode := c.PostForm("memberCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.RemoveMember(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		MemberCode:  memberCode,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (p *HandlerProject) listForInvite(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.ProjectServiceClient.ListForInvite(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		MemberId:    c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	list := resp.List
	if list == nil {
		list = []*project.ProjectMemberMessage{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}

func (p *HandlerProject) sourceLinkDelete(c *gin.Context) {
	result := &common.Result{}
	sourceCode := c.PostForm("sourceCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.DeleteSourceLink(ctx, &task_service_v1.TaskReqMessage{
		FileCode: sourceCode,
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func New() *HandlerProject {
	return &HandlerProject{}
}
