package project

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"

	"test.com/project-api/api/rpc"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-common/min"
	"test.com/project-grpc/user/login"
)

type HandlerIndex struct{}

func NewIndex() *HandlerIndex {
	return &HandlerIndex{}
}

// uploadAvatar 上传头像到 MinIO，更新 ms_member.avatar
func (h *HandlerIndex) uploadAvatar(c *gin.Context) {
	result := &common.Result{}
	memberId := c.GetInt64("memberId")
	if memberId == 0 {
		c.JSON(http.StatusOK, result.Fail(-1, "未登录"))
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(-1, "获取文件失败: "+err.Error()))
		return
	}

	open, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(-1, "打开文件失败"))
		return
	}
	defer open.Close()

	buf := make([]byte, file.Size)
	open.Read(buf)

	ext := path.Ext(file.Filename)
	objectName := fmt.Sprintf("avatar/%d%s", memberId, ext)
	bucketName := min.GetBucket()

	minioClient, err := min.GetMinioClient()
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(-1, "MinIO 客户端获取失败"))
		return
	}
	_, err = minioClient.Put(context.Background(), bucketName, objectName, buf, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(-1, "上传失败: "+err.Error()))
		return
	}

	avatarUrl := "http://" + min.GetEndpoint() + "/" + bucketName + "/" + objectName

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, grpcErr := rpc.LoginServiceClient.EditPersonal(ctx, &login.MemberMessage{
		Id:     memberId,
		Avatar: avatarUrl,
	})
	if grpcErr != nil {
		code, msg := errs.ParseGrpcError(grpcErr)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{"url": avatarUrl}))
}

// editPersonal 更新个人基本信息（name、description、avatar）
func (h *HandlerIndex) editPersonal(c *gin.Context) {
	result := &common.Result{}
	memberId := c.GetInt64("memberId")
	name := c.PostForm("name")
	avatar := c.PostForm("avatar")
	description := c.PostForm("description")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.LoginServiceClient.EditPersonal(ctx, &login.MemberMessage{
		Id:          memberId,
		Name:        name,
		Avatar:      avatar,
		Description: description,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

// uploadImg 上传备注图片到 MinIO
func (h *HandlerIndex) uploadImg(c *gin.Context) {
	result := &common.Result{}
	memberId := c.GetInt64("memberId")
	if memberId == 0 {
		c.JSON(http.StatusOK, result.Fail(-1, "未登录"))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(-1, "获取表单失败: "+err.Error()))
		return
	}

	files := form.File["image[]"]
	if len(files) == 0 {
		c.JSON(http.StatusOK, result.Fail(-1, "未找到图片文件"))
		return
	}

	minioClient, err := min.GetMinioClient()
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(-1, "MinIO 客户端获取失败"))
		return
	}

	bucketName := min.GetBucket()
	var urls []string

	for _, file := range files {
		open, err := file.Open()
		if err != nil {
			continue
		}
		defer open.Close()

		buf := make([]byte, file.Size)
		open.Read(buf)

		ext := path.Ext(file.Filename)
		objectName := fmt.Sprintf("task/%d/%d%s", memberId, time.Now().UnixNano(), ext)

		_, err = minioClient.Put(context.Background(), bucketName, objectName, buf, file.Size, file.Header.Get("Content-Type"))
		if err != nil {
			continue
		}

		avatarUrl := "http://" + min.GetEndpoint() + "/" + bucketName + "/" + objectName
		urls = append(urls, avatarUrl)
	}

	if len(urls) == 0 {
		c.JSON(http.StatusOK, result.Fail(-1, "上传失败"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno": 0,
		"data":  urls,
	})
}
