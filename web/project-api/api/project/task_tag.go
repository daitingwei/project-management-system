package project

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test.com/project-api/api/rpc"
	common "test.com/project-common"
	"test.com/project-common/errs"
	task_service_v1 "test.com/project-grpc/task"
)

type HandlerTaskTag struct {
}

func (t *HandlerTaskTag) list(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := rpc.TaskServiceClient.TaskTagList(ctx, &task_service_v1.TaskReqMessage{
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
		list = []*task_service_v1.TaskTagMessage{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}

func (t *HandlerTaskTag) save(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	name := c.PostForm("name")
	color := c.PostForm("color")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := rpc.TaskServiceClient.SaveTaskTag(ctx, &task_service_v1.TaskReqMessage{
		ProjectCode: projectCode,
		Name:        name,
		Color:       color,
		MemberId:    c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp.Tag))
}

func (t *HandlerTaskTag) edit(c *gin.Context) {
	result := &common.Result{}
	tagCode := c.PostForm("tagCode")
	name := c.PostForm("name")
	color := c.PostForm("color")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := rpc.TaskServiceClient.EditTaskTag(ctx, &task_service_v1.TaskReqMessage{
		TagCode:  tagCode,
		Name:     name,
		Color:    color,
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp.Tag))
}

func (t *HandlerTaskTag) del(c *gin.Context) {
	result := &common.Result{}
	tagCode := c.PostForm("tagCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := rpc.TaskServiceClient.DeleteTaskTag(ctx, &task_service_v1.TaskReqMessage{
		TagCode:  tagCode,
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func NewTaskTag() *HandlerTaskTag {
	return &HandlerTaskTag{}
}
