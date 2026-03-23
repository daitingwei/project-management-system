package project

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"test.com/project-api/api/rpc"
	"test.com/project-api/pkg/model"
	"test.com/project-api/pkg/model/pro"
	"test.com/project-api/pkg/model/tasks"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-common/min"
	"test.com/project-common/tms"
	task_service_v1 "test.com/project-grpc/task"
)

type HandlerTask struct {
}

func (t *HandlerTask) taskStages(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 1.获取参数 校验参数的合法性
	projectCode := c.PostForm("projectCode")
	page := &model.Page{}
	page.Bind(c)
	// 2.调用grpc服务
	msg := &task_service_v1.TaskReqMessage{
		MemberId:    c.GetInt64("memberId"),
		ProjectCode: projectCode,
		Page:        page.Page,
		PageSize:    page.PageSize,
	}
	stages, err := rpc.TaskServiceClient.TaskStages(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	// 3.处理响应
	var list []*tasks.TaskStagesResp
	copier.Copy(&list, stages.List)
	if list == nil {
		list = []*tasks.TaskStagesResp{}
	}
	for _, v := range list {
		v.TasksLoading = true  // 任务加载状态
		v.FixedCreator = false // 添加任务按钮定位
		v.ShowTaskCard = false // 是否显示创建卡片
		v.Tasks = []int{}
		v.DoneTasks = []int{}
		v.UnDoneTasks = []int{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  list,
		"total": stages.Total,
		"page":  page.Page,
	}))
}

func (t *HandlerTask) memberProjectList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 1.获取参数 校验参数的合法性
	projectCode := c.PostForm("projectCode")
	page := &model.Page{}
	page.Bind(c)
	// 2.调用grpc服务
	msg := &task_service_v1.TaskReqMessage{
		MemberId:    c.GetInt64("memberId"),
		ProjectCode: projectCode,
		Page:        page.Page,
		PageSize:    page.PageSize,
	}
	resp, err := rpc.TaskServiceClient.MemberProjectList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	var list []*pro.MemberProjectResp
	copier.Copy(&list, resp.List)
	if list == nil {
		list = []*pro.MemberProjectResp{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  list,
		"total": resp.Total,
		"page":  page.Page,
	}))

}

func (t *HandlerTask) taskList(c *gin.Context) {
	result := &common.Result{}
	stageCode := c.PostForm("stageCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list, err := rpc.TaskServiceClient.TaskList(ctx, &task_service_v1.TaskReqMessage{StageCode: stageCode, MemberId: c.GetInt64("memberId")})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var taskDisplayList []*tasks.TaskDisplay
	copier.Copy(&taskDisplayList, list.List)
	if taskDisplayList == nil {
		taskDisplayList = []*tasks.TaskDisplay{}
	}
	// 返回给前端的数据 一定不要是null
	for _, v := range taskDisplayList {
		if v.Tags == nil {
			v.Tags = []int{}
		}
		if v.ChildCount == nil {
			v.ChildCount = []int{}
		}
	}
	c.JSON(http.StatusOK, result.Success(taskDisplayList))
}

func (t *HandlerTask) saveTask(c *gin.Context) {
	result := &common.Result{}
	var req *tasks.TaskSaveReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		ProjectCode: req.ProjectCode,
		Name:        req.Name,
		StageCode:   req.StageCode,
		AssignTo:    req.AssignTo,
		MemberId:    c.GetInt64("memberId"),
	}
	taskMessage, err := rpc.TaskServiceClient.SaveTask(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	td := &tasks.TaskDisplay{}
	copier.Copy(td, taskMessage)
	if td != nil {
		if td.Tags == nil {
			td.Tags = []int{}
		}
		if td.ChildCount == nil {
			td.ChildCount = []int{}
		}
	}
	c.JSON(http.StatusOK, result.Success(td))
}

func (t *HandlerTask) taskSort(c *gin.Context) {
	result := &common.Result{}
	var req *tasks.TaskSortReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		PreTaskCode:  req.PreTaskCode,
		NextTaskCode: req.NextTaskCode,
		ToStageCode:  req.ToStageCode,
	}
	_, err := rpc.TaskServiceClient.TaskSort(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (t *HandlerTask) taskStagesSort(c *gin.Context) {
	result := &common.Result{}
	var req *tasks.TaskStagesSortReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskStagesSortReqMessage{
		PreStageCode:  req.PreStageCode,
		NextStageCode: req.NextStageCode,
		ProjectCode:   req.ProjectCode,
	}
	_, err := rpc.TaskServiceClient.TaskStagesSort(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (t *HandlerTask) taskDone(c *gin.Context) {
	result := &common.Result{}
	var req struct {
		TaskCode string `form:"taskCode" json:"taskCode"`
		Done     string `form:"done" json:"done"`
	}
	c.ShouldBind(&req)
	zap.L().Info("taskDone req", zap.String("taskCode", req.TaskCode), zap.String("done", req.Done))
	done := 0
	if req.Done != "" {
		var err error
		done, err = strconv.Atoi(req.Done)
		if err != nil {
			done = 0
		}
	}
	zap.L().Info("taskDone req", zap.String("taskCode", req.TaskCode), zap.String("done", req.Done), zap.Int("doneInt", done))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		TaskCode: req.TaskCode,
		Done:     int32(done),
		MemberId: c.GetInt64("memberId"),
	}
	zap.L().Info("taskDone req", zap.String("taskCode", req.TaskCode), zap.Int32("done", msg.Done))
	_, err := rpc.TaskServiceClient.TaskDone(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (t *HandlerTask) myTaskList(c *gin.Context) {
	result := &common.Result{}
	var req *tasks.MyTaskReq
	c.ShouldBind(&req)
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	taskType := int32(req.TaskType)
	if taskType == 0 {
		taskType = 1
	}
	msg := &task_service_v1.TaskReqMessage{
		MemberId: memberId,
		TaskType: taskType,
		Type:     int32(req.Type),
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	myTaskListResponse, err := rpc.TaskServiceClient.MyTaskList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var myTaskList []*tasks.MyTaskDisplay
	copier.Copy(&myTaskList, myTaskListResponse.List)
	if myTaskList == nil {
		myTaskList = []*tasks.MyTaskDisplay{}
	}
	for _, v := range myTaskList {
		v.ProjectInfo = tasks.ProjectInfo{
			Name: v.ProjectName,
			Code: v.ProjectCode,
		}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  myTaskList,
		"total": myTaskListResponse.Total,
	}))
}

func (t *HandlerTask) readTask(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		MemberId: c.GetInt64("memberId"),
	}
	taskMessage, err := rpc.TaskServiceClient.ReadTask(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	td := &tasks.TaskDisplay{}
	copier.Copy(td, taskMessage)
	if td != nil {
		if td.Tags == nil {
			td.Tags = []int{}
		}
		if td.ChildCount == nil {
			td.ChildCount = []int{}
		}
	}
	c.JSON(200, result.Success(td))
}

func (t *HandlerTask) listTaskMember(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	page := &model.Page{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		MemberId: c.GetInt64("memberId"),
		Page:     page.Page,
		PageSize: page.PageSize,
	}
	taskMemberResponse, err := rpc.TaskServiceClient.ListTaskMember(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var tms []*tasks.TaskMember
	copier.Copy(&tms, taskMemberResponse.List)
	if tms == nil {
		tms = []*tasks.TaskMember{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  tms,
		"total": taskMemberResponse.Total,
		"page":  page.Page,
	}))
}

func (t *HandlerTask) taskLog(c *gin.Context) {
	result := &common.Result{}
	var req *model.TaskLogReq
	c.ShouldBind(&req)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		TaskCode: req.TaskCode,
		MemberId: c.GetInt64("memberId"),
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		All:      int32(req.All),
		Comment:  int32(req.Comment),
	}
	taskLogResponse, err := rpc.TaskServiceClient.TaskLog(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var tms []*model.ProjectLogDisplay
	copier.Copy(&tms, taskLogResponse.List)
	if tms == nil {
		tms = []*model.ProjectLogDisplay{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  tms,
		"total": taskLogResponse.Total,
		"page":  req.Page,
	}))
}

func (t *HandlerTask) taskWorkTimeList(c *gin.Context) {
	taskCode := c.PostForm("taskCode")
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		MemberId: c.GetInt64("memberId"),
	}
	taskWorkTimeResponse, err := rpc.TaskServiceClient.TaskWorkTimeList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var tms []*model.TaskWorkTime
	copier.Copy(&tms, taskWorkTimeResponse.List)
	if tms == nil {
		tms = []*model.TaskWorkTime{}
	}
	c.JSON(http.StatusOK, result.Success(tms))
}

func (t *HandlerTask) saveTaskWorkTime(c *gin.Context) {
	result := &common.Result{}
	var req *model.SaveTaskWorkTimeReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		TaskCode:  req.TaskCode,
		MemberId:  c.GetInt64("memberId"),
		Content:   req.Content,
		Num:       int32(req.Num),
		BeginTime: tms.ParseTime(req.BeginTime),
	}
	_, err := rpc.TaskServiceClient.SaveTaskWorkTime(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (t *HandlerTask) uploadFiles(c *gin.Context) {
	result := &common.Result{}
	req := model.UploadFileReq{}
	c.ShouldBind(&req)
	multipartForm, _ := c.MultipartForm()
	file := multipartForm.File
	uploadFile := file["file"][0]
	bucketName := "msproject"
	key := bucketName + "/" + req.Filename

	minioClient, err := min.New(
		"localhost:9009",
		"admin",
		"admin123456",
		false,
	)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(-999, err.Error()))
		return
	}

	if req.TotalChunks == 1 {
		open, err := uploadFile.Open()
		if err != nil {
			c.JSON(http.StatusOK, result.Fail(-999, err.Error()))
			return
		}
		defer open.Close()
		buf := make([]byte, req.CurrentChunkSize)
		open.Read(buf)
		_, err = minioClient.Put(
			context.Background(),
			bucketName,
			req.Filename,
			buf,
			int64(req.TotalSize),
			uploadFile.Header.Get("Content-Type"),
		)
		if err != nil {
			c.JSON(http.StatusOK, result.Fail(-999, err.Error()))
			return
		}
	}
	if req.TotalChunks > 1 {
		open, err := uploadFile.Open()
		if err != nil {
			c.JSON(http.StatusOK, result.Fail(-999, err.Error()))
			return
		}
		defer open.Close()
		buf := make([]byte, req.CurrentChunkSize)
		open.Read(buf)
		formatInt := strconv.FormatInt(int64(req.ChunkNumber), 10)
		_, err = minioClient.Put(
			context.Background(),
			bucketName,
			req.Filename+"_"+formatInt,
			buf,
			int64(req.CurrentChunkSize),
			uploadFile.Header.Get("Content-Type"),
		)
		if err != nil {
			c.JSON(http.StatusOK, result.Fail(-999, err.Error()))
			return
		}
		if req.TotalChunks == req.ChunkNumber {
			_, err := minioClient.Compose(context.Background(), bucketName, req.Filename, req.TotalChunks)
			if err != nil {
				c.JSON(http.StatusOK, result.Fail(-999, err.Error()))
				return
			}
			var chunkKeys []string
			for i := 1; i <= req.TotalChunks; i++ {
				chunkKeys = append(chunkKeys, req.Filename+"_"+strconv.FormatInt(int64(i), 10))
			}
			minioClient.DeleteObjects(context.Background(), bucketName, chunkKeys)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	minioEndpoint := min.GetEndpoint()
	fileUrl := "http://" + minioEndpoint + "/" + key
	msg := &task_service_v1.TaskFileReqMessage{
		TaskCode:         req.TaskCode,
		ProjectCode:      req.ProjectCode,
		OrganizationCode: c.GetString("organizationCode"),
		PathName:         key,
		FileName:         req.Filename,
		Size:             int64(req.TotalSize),
		Extension:        path.Ext(key),
		FileUrl:          fileUrl,
		FileType:         file["file"][0].Header.Get("Content-Type"),
		MemberId:         c.GetInt64("memberId"),
		LinkType:         "task",
	}
	if req.TotalChunks == req.ChunkNumber || req.TotalChunks == 1 {
		_, err := rpc.TaskServiceClient.SaveTaskFile(ctx, msg)
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
		}
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"file":        key,
		"hash":        "",
		"key":         key,
		"url":         "http://" + minioEndpoint + "/" + key,
		"projectName": req.ProjectName,
	}))
}

func (t *HandlerTask) taskSources(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	sources, err := rpc.TaskServiceClient.TaskSources(ctx, &task_service_v1.TaskReqMessage{TaskCode: taskCode})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var slList []*model.SourceLink
	copier.Copy(&slList, sources.List)
	if slList == nil {
		slList = []*model.SourceLink{}
	}
	c.JSON(http.StatusOK, result.Success(slList))
}

func (t *HandlerTask) createComment(c *gin.Context) {
	result := &common.Result{}
	req := model.CommentReq{}
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		TaskCode:       req.TaskCode,
		CommentContent: req.Comment,
		Mentions:       req.Mentions,
		MemberId:       c.GetInt64("memberId"),
	}
	_, err := rpc.TaskServiceClient.CreateComment(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) saveTaskStages(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	name := c.PostForm("name")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		ProjectCode: projectCode,
		Name:        name,
		MemberId:    c.GetInt64("memberId"),
	}
	resp, err := rpc.TaskServiceClient.SaveTaskStages(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	if len(resp.List) > 0 {
		c.JSON(http.StatusOK, result.Success(resp.List[0]))
	} else {
		c.JSON(http.StatusOK, result.Success(nil))
	}
}

func (t *HandlerTask) editTaskStages(c *gin.Context) {
	result := &common.Result{}
	stageCode := c.PostForm("stageCode")
	name := c.PostForm("name")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		StageCode: stageCode,
		Name:      name,
		MemberId:  c.GetInt64("memberId"),
	}
	resp, err := rpc.TaskServiceClient.EditTaskStages(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp))
}

func (t *HandlerTask) deleteTaskStages(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stageCode := c.PostForm("code")
	msg := &task_service_v1.TaskReqMessage{
		StageCode: stageCode,
		MemberId:  c.GetInt64("memberId"),
	}

	_, err := rpc.TaskServiceClient.DeleteTaskStages(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) fileIndex(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	projectCode := c.PostForm("projectCode")
	deleted := c.DefaultPostForm("deleted", "0")
	page := &model.Page{}
	page.Bind(c)

	msg := &task_service_v1.TaskReqMessage{
		ProjectCode: projectCode,
		Deleted:     int32(parseInt(deleted)),
		Page:        page.Page,
		PageSize:    page.PageSize,
	}

	resp, err := rpc.TaskServiceClient.FileIndex(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  resp.List,
		"total": resp.Total,
		"page":  page.Page,
	}))
}

func (t *HandlerTask) fileRead(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fileCode := c.PostForm("fileCode")
	msg := &task_service_v1.TaskReqMessage{
		FileCode: fileCode,
	}

	resp, err := rpc.TaskServiceClient.FileRead(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp))
}

func (t *HandlerTask) fileEdit(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fileCode := c.PostForm("fileCode")
	title := c.PostForm("title")
	msg := &task_service_v1.TaskReqMessage{
		FileCode: fileCode,
		Title:    title,
	}

	resp, err := rpc.TaskServiceClient.FileEdit(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp))
}

func (t *HandlerTask) fileRecycle(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fileCode := c.PostForm("fileCode")
	msg := &task_service_v1.TaskReqMessage{
		FileCode: fileCode,
	}

	_, err := rpc.TaskServiceClient.FileRecycle(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) fileRecovery(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fileCode := c.PostForm("fileCode")
	msg := &task_service_v1.TaskReqMessage{
		FileCode: fileCode,
	}

	_, err := rpc.TaskServiceClient.FileRecovery(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) fileDelete(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fileCode := c.PostForm("fileCode")
	msg := &task_service_v1.TaskReqMessage{
		FileCode: fileCode,
	}

	_, err := rpc.TaskServiceClient.FileDelete(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) fileDownload(c *gin.Context) {
	result := &common.Result{}
	fileCode := c.Query("fileCode")
	if fileCode == "" {
		fileCode = c.PostForm("fileCode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fileResp, err := rpc.TaskServiceClient.FileRead(ctx, &task_service_v1.TaskReqMessage{FileCode: fileCode})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	if fileResp.File == nil || fileResp.File.PathName == "" {
		c.JSON(http.StatusOK, result.Fail(-1, "文件不存在"))
		return
	}

	minioClient, err := min.GetMinioClient()
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(-999, err.Error()))
		return
	}

	bucket := min.GetBucket()
	object, err := minioClient.Get(ctx, bucket, fileResp.File.PathName)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(-999, err.Error()))
		return
	}
	defer object.Close()

	fileName := fileResp.File.FileName
	if fileResp.File.Extension != "" {
		fileName = fileName + fileResp.File.Extension
	}

	contentType := fileResp.File.FileType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.QueryEscape(fileName))
	c.Header("Content-Length", strconv.FormatInt(fileResp.File.Size, 10))
	c.Header("Access-Control-Allow-Origin", "*")

	http.ServeContent(c.Writer, c.Request, fileName, time.Now(), object)
}

func parseInt(s string) int {
	i := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			i = i*10 + int(c-'0')
		}
	}
	return i
}

func (t *HandlerTask) editTask(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 前端用 querystring 序列化，发的是 application/x-www-form-urlencoded
	// 用 c.Request.PostForm 读取，区分"字段未传"和"字段传了空值"
	c.Request.ParseForm()
	form := c.Request.PostForm

	taskCode := form.Get("taskCode")

	// 时间字段：未传 → math.MinInt64（service 层跳过）；传了空串 → -1（service 层清零）；有值 → 正常解析
	const timeNotSet = int64(-9223372036854775808) // math.MinInt64
	beginTime := timeNotSet
	if _, ok := form["begin_time"]; ok {
		s := form.Get("begin_time")
		if s == "" {
			beginTime = -1
		} else {
			beginTime = tms.ParseTime(s)
		}
	}
	endTime := timeNotSet
	if _, ok := form["end_time"]; ok {
		s := form.Get("end_time")
		if s == "" {
			endTime = -1
		} else {
			endTime = tms.ParseTime(s)
		}
	}

	priStr := form.Get("pri")
	workTimeStr := form.Get("work_time")

	// status: 未传 → -1（service 层跳过）；传了值（包括0）→ 正常更新
	statusVal := int32(-1)
	if _, ok := form["status"]; ok {
		statusVal = int32(parseInt(form.Get("status")))
	}
	// pri: 未传 → -1（service 层跳过）；传了值（包括0）→ 正常更新
	priVal := int32(-1)
	if _, ok := form["pri"]; ok {
		priVal = int32(parseInt(priStr))
	}

	msg := &task_service_v1.TaskReqMessage{
		TaskCode:  taskCode,
		MemberId:  c.GetInt64("memberId"),
		Name:      form.Get("name"),
		Content:   form.Get("description"),
		Pri:       priVal,
		Status:    statusVal,
		Num:       int32(parseInt(workTimeStr)),
		BeginTime: beginTime,
		EndTime:   endTime,
	}
	taskMessage, err := rpc.TaskServiceClient.EditTask(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	td := &tasks.TaskDisplay{}
	copier.Copy(td, taskMessage)
	if td != nil {
		if td.Tags == nil {
			td.Tags = []int{}
		}
		if td.ChildCount == nil {
			td.ChildCount = []int{}
		}
	}
	c.JSON(http.StatusOK, result.Success(td))
}

func (t *HandlerTask) batchAssignTask(c *gin.Context) {
	result := &common.Result{}
	type req struct {
		TaskCodes     []string `form:"taskCodes" json:"taskCodes"`
		BatchAssignTo string   `form:"assignTo" json:"assignTo"`
	}
	var r req
	c.ShouldBind(&r)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.BatchAssignTask(ctx, &task_service_v1.TaskReqMessage{
		TaskCodes:     r.TaskCodes,
		BatchAssignTo: r.BatchAssignTo,
		MemberId:      c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) editTaskWorkTime(c *gin.Context) {
	result := &common.Result{}
	var req *model.SaveTaskWorkTimeReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		TaskCode:  req.TaskCode,
		MemberId:  c.GetInt64("memberId"),
		Content:   req.Content,
		Num:       int32(req.Num),
		BeginTime: tms.ParseTime(req.BeginTime),
	}
	_, err := rpc.TaskServiceClient.SaveTaskWorkTime(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (t *HandlerTask) delTaskWorkTime(c *gin.Context) {
	result := &common.Result{}
	workTimeCode := c.PostForm("workTimeCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.DelTaskWorkTime(ctx, &task_service_v1.TaskReqMessage{
		WorkTimeCode: workTimeCode,
		MemberId:     c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) taskLike(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	like := c.PostForm("like")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.TaskServiceClient.TaskLike(ctx, &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		Like:     int32(parseInt(like)),
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp.Ok))
}

func (t *HandlerTask) taskStar(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	star := c.PostForm("star")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.TaskServiceClient.TaskStar(ctx, &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		Star:     int32(parseInt(star)),
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp.Ok))
}

func (t *HandlerTask) setPrivate(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	private := c.PostForm("private")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.TaskServiceClient.SetPrivate(ctx, &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		Private:  int32(parseInt(private)),
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp.Ok))
}

func (t *HandlerTask) recycleBatch(c *gin.Context) {
	result := &common.Result{}
	var taskCodes []string
	c.ShouldBind(&struct{ TaskCodes *[]string `form:"taskCodes" json:"taskCodes"` }{&taskCodes})
	// 尝试从 JSON body 读取
	type req struct {
		TaskCodes []string `form:"taskCodes" json:"taskCodes"`
	}
	var r req
	c.ShouldBind(&r)
	if len(r.TaskCodes) > 0 {
		taskCodes = r.TaskCodes
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.RecycleBatch(ctx, &task_service_v1.TaskReqMessage{
		TaskCodes: taskCodes,
		MemberId:  c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) taskToTags(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := rpc.TaskServiceClient.TaskToTags(ctx, &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	list := resp.List
	if list == nil {
		list = []*task_service_v1.TaskToTagItem{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}

func (t *HandlerTask) setTag(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	tagCode := c.PostForm("tagCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := rpc.TaskServiceClient.SetTag(ctx, &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		TagCode:  tagCode,
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp.Ok))
}

func NewTask() *HandlerTask {
	return &HandlerTask{}
}

func (t *HandlerTask) taskListByProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	page := c.PostForm("page")
	pageSize := c.PostForm("pageSize")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task_service_v1.TaskReqMessage{
		ProjectCode: projectCode,
		Page:        int64(parseInt(page)),
		PageSize:    int64(parseInt(pageSize)),
		MemberId:    c.GetInt64("memberId"),
	}
	resp, err := rpc.TaskServiceClient.TaskList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(resp))
}

func (t *HandlerTask) taskRecycle(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.TaskRecycle(ctx, &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) taskDelete(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.TaskDelete(ctx, &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) taskRecovery(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.TaskRecovery(ctx, &task_service_v1.TaskReqMessage{
		TaskCode: taskCode,
		MemberId: c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) assignTask(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	executorCode := c.PostForm("executorCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.BatchAssignTask(ctx, &task_service_v1.TaskReqMessage{
		TaskCodes:     []string{taskCode},
		BatchAssignTo: executorCode,
		MemberId:      c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) inviteMember(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	memberCode := c.PostForm("memberCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.InviteMember(ctx, &task_service_v1.TaskReqMessage{
		TaskCode:  taskCode,
		MemberCode: memberCode,
		MemberId:  c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) inviteMemberBatch(c *gin.Context) {
	result := &common.Result{}
	taskCode := c.PostForm("taskCode")
	type req struct {
		MemberCodes []string `form:"memberCodes" json:"memberCodes"`
	}
	var r req
	c.ShouldBind(&r)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.TaskServiceClient.InviteMemberBatch(ctx, &task_service_v1.TaskReqMessage{
		TaskCode:  taskCode,
		TaskCodes: r.MemberCodes,
		MemberId:  c.GetInt64("memberId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (t *HandlerTask) projectStats(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.TaskServiceClient.ProjectStats(ctx, &task_service_v1.TaskReqMessage{
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

func (t *HandlerTask) getProjectReport(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.TaskServiceClient.ProjectReport(ctx, &task_service_v1.TaskReqMessage{
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

func (t *HandlerTask) dateTotalForProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := rpc.TaskServiceClient.DateTotalForProject(ctx, &task_service_v1.TaskReqMessage{
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

// getStr 从 map 中安全取字符串
func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getInt 从 map 中安全取整数（JSON number 默认是 float64）
func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return 0
}
