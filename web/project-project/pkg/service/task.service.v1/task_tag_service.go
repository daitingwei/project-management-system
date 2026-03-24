package task_service_v1

import (
	"context"
	"time"

	"go.uber.org/zap"

	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-grpc/task"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/pkg/model"
)

// TaskTagList 查询项目下的标签列表
func (t *TaskService) TaskTagList(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskTagListResponse, error) {
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tags, err := t.taskTagRepo.FindTagsByProjectCode(c, projectCode)
	if err != nil {
		zap.L().Error("TaskTagList FindTagsByProjectCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var list []*task.TaskTagMessage
	for _, v := range tags {
		list = append(list, &task.TaskTagMessage{
			Code:        encrypts.EncryptNoErr(v.Id),
			Name:        v.Name,
			Color:       v.Color,
			ProjectCode: msg.ProjectCode,
		})
	}
	return &task.TaskTagListResponse{List: list}, nil
}

// SaveTaskTag 创建标签
func (t *TaskService) SaveTaskTag(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskTagResponse, error) {
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 检查同名
	exist, err := t.taskTagRepo.FindTagByNameAndProject(c, projectCode, msg.Name)
	if err != nil {
		zap.L().Error("SaveTaskTag FindTagByNameAndProject error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist != nil {
		return nil, errs.GrpcError(model.TaskTagNameExist)
	}

	tag := &data.TaskTag{
		ProjectCode: projectCode,
		Name:        msg.Name,
		Color:       msg.Color,
	}
	if err := t.taskTagRepo.SaveTag(c, tag); err != nil {
		zap.L().Error("SaveTaskTag SaveTag error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskTagResponse{
		Tag: &task.TaskTagMessage{
			Code:        encrypts.EncryptNoErr(tag.Id),
			Name:        tag.Name,
			Color:       tag.Color,
			ProjectCode: msg.ProjectCode,
		},
	}, nil
}

// EditTaskTag 编辑标签
func (t *TaskService) EditTaskTag(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskTagResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	old, err := t.taskTagRepo.FindTagByCode(c, msg.TagCode)
	if err != nil || old == nil {
		return nil, errs.GrpcError(model.DBError)
	}
	old.Name = msg.Name
	old.Color = msg.Color
	if err := t.taskTagRepo.UpdateTag(c, old); err != nil {
		zap.L().Error("EditTaskTag UpdateTag error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskTagResponse{
		Tag: &task.TaskTagMessage{
			Code:  msg.TagCode,
			Name:  old.Name,
			Color: old.Color,
		},
	}, nil
}

// DeleteTaskTag 删除标签
func (t *TaskService) DeleteTaskTag(ctx context.Context, msg *task.TaskReqMessage) (*task.DeleteTaskTagResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := t.taskTagRepo.DeleteTag(c, msg.TagCode); err != nil {
		zap.L().Error("DeleteTaskTag DeleteTag error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.DeleteTaskTagResponse{}, nil
}

// TaskToTags 查询任务关联的标签列表
func (t *TaskService) TaskToTags(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskToTagsResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	toTags, err := t.taskTagRepo.FindToTagsByTaskCode(c, taskCode)
	if err != nil {
		zap.L().Error("TaskToTags FindToTagsByTaskCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var list []*task.TaskToTagItem
	for _, v := range toTags {
		list = append(list, &task.TaskToTagItem{
			TagCode: encrypts.EncryptNoErr(v.TagCode),
		})
	}
	return &task.TaskToTagsResponse{List: list}, nil
}

// SetTag 切换任务标签（存在则删除，不存在则插入）
func (t *TaskService) SetTag(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	tagCode := encrypts.DecryptNoErr(msg.TagCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	exist, err := t.taskTagRepo.FindToTagByTaskAndTag(c, taskCode, tagCode)
	if err != nil {
		zap.L().Error("SetTag FindToTagByTaskAndTag error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist != nil {
		// 已存在，删除
		if err := t.taskTagRepo.DeleteToTag(c, exist.Id); err != nil {
			zap.L().Error("SetTag DeleteToTag error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		return &task.CommonBoolResponse{Ok: false}, nil
	}
	// 不存在，插入
	tt := &data.TaskToTag{
		TaskCode: taskCode,
		TagCode:  tagCode,
	}
	if err := t.taskTagRepo.SaveToTag(c, tt); err != nil {
		zap.L().Error("SetTag SaveToTag error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.CommonBoolResponse{Ok: true}, nil
}

// ProjectStats 项目任务统计（stub）
func (t *TaskService) ProjectStats(ctx context.Context, msg *task.TaskReqMessage) (*task.ProjectStatsResponse, error) {
	return &task.ProjectStatsResponse{}, nil
}

// ProjectReport 项目任务报告（stub）
func (t *TaskService) ProjectReport(ctx context.Context, msg *task.TaskReqMessage) (*task.ProjectReportResponse, error) {
	return &task.ProjectReportResponse{}, nil
}

// DateTotalForProject 项目按日期任务统计（stub）
func (t *TaskService) DateTotalForProject(ctx context.Context, msg *task.TaskReqMessage) (*task.DateTotalListResponse, error) {
	return &task.DateTotalListResponse{}, nil
}

// DelTaskWorkTime 删除工时记录
func (t *TaskService) DelTaskWorkTime(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	id := encrypts.DecryptNoErr(msg.WorkTimeCode)
	if berr := t.taskWorkTimeDomain.DelTaskWorkTime(id); berr != nil {
		return nil, errs.GrpcError(berr)
	}
	return &task.CommonBoolResponse{Ok: true}, nil
}

// TaskLike 点赞/取消点赞任务
func (t *TaskService) TaskLike(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ts, err := t.taskRepo.FindTaskById(c, taskCode)
	if err != nil || ts == nil {
		return nil, errs.GrpcError(model.DBError)
	}
	ts.Like = int(msg.Like)
	if err := t.transaction.Action(func(conn database.DbConn) error {
		return t.taskRepo.UpdateTask(c, conn, ts)
	}); err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.CommonBoolResponse{Ok: ts.Like == 1}, nil
}

// TaskStar 收藏/取消收藏任务
func (t *TaskService) TaskStar(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ts, err := t.taskRepo.FindTaskById(c, taskCode)
	if err != nil || ts == nil {
		return nil, errs.GrpcError(model.DBError)
	}
	ts.Star = int(msg.Star)
	if err := t.transaction.Action(func(conn database.DbConn) error {
		return t.taskRepo.UpdateTask(c, conn, ts)
	}); err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.CommonBoolResponse{Ok: ts.Star == 1}, nil
}

// SetPrivate 设置任务隐私模式
func (t *TaskService) SetPrivate(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ts, err := t.taskRepo.FindTaskById(c, taskCode)
	if err != nil || ts == nil {
		return nil, errs.GrpcError(model.DBError)
	}
	ts.Private = int(msg.Private)
	if err := t.transaction.Action(func(conn database.DbConn) error {
		return t.taskRepo.UpdateTask(c, conn, ts)
	}); err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.CommonBoolResponse{Ok: true}, nil
}

// RecycleBatch 批量回收任务（软删除）
func (t *TaskService) RecycleBatch(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var taskCodes []int64
	for _, code := range msg.TaskCodes {
		taskCodes = append(taskCodes, encrypts.DecryptNoErr(code))
	}
	if len(taskCodes) == 0 {
		return &task.CommonBoolResponse{Ok: true}, nil
	}
	if err := t.taskRepo.BatchUpdateDeleted(c, taskCodes, 1); err != nil {
		zap.L().Error("RecycleBatch BatchUpdateDeleted error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.CommonBoolResponse{Ok: true}, nil
}

// BatchAssignTask 批量分配任务负责人
func (t *TaskService) BatchAssignTask(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var taskCodes []int64
	for _, code := range msg.TaskCodes {
		taskCodes = append(taskCodes, encrypts.DecryptNoErr(code))
	}
	if len(taskCodes) == 0 {
		return &task.CommonBoolResponse{Ok: true}, nil
	}
	assignTo := encrypts.DecryptNoErr(msg.BatchAssignTo)
	if err := t.taskRepo.BatchUpdateAssignTo(c, taskCodes, assignTo); err != nil {
		zap.L().Error("BatchAssignTask BatchUpdateAssignTo error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.CommonBoolResponse{Ok: true}, nil
}

// InviteMember 单个邀请成员加入任务
func (t *TaskService) InviteMember(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	memberCode := encrypts.DecryptNoErr(msg.MemberCode)

	taskInfo, err := t.taskRepo.FindTaskById(c, taskCode)
	if err != nil {
		zap.L().Error("InviteMember FindTaskById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if taskInfo == nil {
		return nil, errs.GrpcError(model.ParamsError)
	}

	existingMember, err := t.taskRepo.FindTaskMemberByTaskAndMember(c, taskCode, memberCode)
	if err != nil {
		zap.L().Error("InviteMember FindTaskMemberByTaskAndMember error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	if existingMember != nil {
		taskInfo.AssignTo = memberCode
		existingMember.IsExecutor = 1
		err = t.transaction.Action(func(conn database.DbConn) error {
			if err := t.taskRepo.UpdateTaskMember(c, existingMember); err != nil {
				zap.L().Error("InviteMember UpdateTaskMember error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			if err := t.taskRepo.UpdateTask(c, conn, taskInfo); err != nil {
				zap.L().Error("InviteMember UpdateTask error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		newMember := &data.TaskMember{
			TaskCode:   taskCode,
			MemberCode: memberCode,
			IsExecutor: 1,
			IsOwner:    0,
		}
		taskInfo.AssignTo = memberCode
		err = t.transaction.Action(func(conn database.DbConn) error {
			if err := t.taskRepo.SaveTaskMember(c, conn, newMember); err != nil {
				zap.L().Error("InviteMember SaveTaskMember error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			if err := t.taskRepo.UpdateTask(c, conn, taskInfo); err != nil {
				zap.L().Error("InviteMember UpdateTask error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return &task.CommonBoolResponse{Ok: true}, nil
}

// InviteMemberBatch 批量邀请成员加入任务
func (t *TaskService) InviteMemberBatch(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	var members []*data.TaskMember
	for _, code := range msg.TaskCodes {
		memberCode := encrypts.DecryptNoErr(code)
		members = append(members, &data.TaskMember{
			TaskCode:   taskCode,
			MemberCode: memberCode,
			IsExecutor: 0,
			IsOwner:    0,
		})
	}
	if len(members) == 0 {
		return &task.CommonBoolResponse{Ok: true}, nil
	}
	if err := t.taskRepo.BatchSaveTaskMember(c, members); err != nil {
		zap.L().Error("InviteMemberBatch BatchSaveTaskMember error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.CommonBoolResponse{Ok: true}, nil
}

// GetListByTaskTag 按标签查询任务列表
func (t *TaskService) GetListByTaskTag(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskListResponse, error) {
	tagCode := encrypts.DecryptNoErr(msg.TagCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	taskCodes, err := t.taskTagRepo.FindTaskCodesByTagCode(c, tagCode)
	if err != nil {
		zap.L().Error("GetListByTaskTag FindTaskCodesByTagCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if len(taskCodes) == 0 {
		return &task.TaskListResponse{}, nil
	}
	tasks, err := t.taskRepo.FindTaskByIds(c, taskCodes)
	if err != nil {
		zap.L().Error("GetListByTaskTag FindTaskByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var list []*task.TaskMessage
	for _, v := range tasks {
		list = append(list, &task.TaskMessage{
			Id:          v.Id,
			ProjectCode: encrypts.EncryptNoErr(v.ProjectCode),
			Name:        v.Name,
			Done:        int32(v.Done),
			Deleted:     int32(v.Deleted),
			StageCode:   encrypts.EncryptNoErr(int64(v.StageCode)),
			Code:        encrypts.EncryptNoErr(v.Id),
		})
	}
	return &task.TaskListResponse{List: list}, nil
}

// DeleteSourceLink 删除资源链接
func (t *TaskService) DeleteSourceLink(ctx context.Context, msg *task.TaskReqMessage) (*task.CommonBoolResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := t.sourceLinkRepo.DeleteByCode(c, msg.FileCode); err != nil {
		zap.L().Error("DeleteSourceLink DeleteByCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.CommonBoolResponse{Ok: true}, nil
}
