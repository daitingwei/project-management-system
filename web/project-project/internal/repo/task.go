package repo

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
)

type TaskStagesTemplateRepo interface {
	Save(ctx context.Context, tst *data.MsTaskStagesTemplate) error
	FindInProTemIds(ctx context.Context, ids []int) ([]data.MsTaskStagesTemplate, error)
	FindByProjectTemplateId(ctx context.Context, projectTemplateCode int) (list []*data.MsTaskStagesTemplate, err error)
}
type TaskStagesRepo interface {
	SaveTaskStages(ctx context.Context, conn database.DbConn, ts *data.TaskStages) error
	FindStagesByProjectId(ctx context.Context, projectCode int64, page int64, pageSize int64) (list []*data.TaskStages, total int64, err error)
	FindById(ctx context.Context, id int) (ts *data.TaskStages, err error)
	FindMaxSortByProjectCode(ctx context.Context, projectCode int64) (maxSort int, err error)
	FindByProjectId(ctx context.Context, projectCode int64) (list []*data.TaskStages, err error)
	FindByStageCodeLtSort(ctx context.Context, projectCode int64, sort int) (ts *data.TaskStages, err error)
	UpdateTaskStagesSort(ctx context.Context, conn database.DbConn, ts *data.TaskStages) error
	DeleteTaskStages(ctx context.Context, conn database.DbConn, id int) error
}

type TaskRepo interface {
	FindTaskByStageCode(ctx context.Context, stageCode int) (list []*data.Task, err error)
	FindTaskMemberByTaskId(ctx context.Context, taskCode int64, memberId int64) (task *data.TaskMember, err error)
	FindTaskMemberByTaskAndMember(ctx context.Context, taskCode int64, memberCode int64) (task *data.TaskMember, err error)
	FindTaskMaxIdNum(ctx context.Context, projectCode int64) (v *int, err error)
	FindTaskSort(ctx context.Context, projectCode int64, stageCode int64) (v *int, err error)
	SaveTask(ctx context.Context, conn database.DbConn, ts *data.Task) error
	SaveTaskMember(ctx context.Context, conn database.DbConn, tm *data.TaskMember) error
	UpdateTaskMember(ctx context.Context, tm *data.TaskMember) error
	FindTaskById(ctx context.Context, taskCode int64) (ts *data.Task, err error)
	UpdateTaskSort(ctx context.Context, conn database.DbConn, ts *data.Task) error
	UpdateTaskDone(ctx context.Context, taskCode int64, done int) error
	UpdateTask(ctx context.Context, conn database.DbConn, ts *data.Task) error
	FindTaskByStageCodeLtSort(ctx context.Context, stageCode int, sort int) (ts *data.Task, err error)
	FindTaskByAssignTo(ctx context.Context, memberId int64, done int, page int64, size int64) ([]*data.Task, int64, error)
	FindTaskByMemberCode(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*data.Task, total int64, err error)
	FindTaskByCreateBy(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*data.Task, total int64, err error)
	FindTaskMemberPage(ctx context.Context, taskCode int64, page int64, size int64) (list []*data.TaskMember, total int64, err error)
	FindTaskByIds(ctx context.Context, taskIdList []int64) (list []*data.Task, err error)
	RecoverTask(ctx context.Context, taskCode int64, deleted int) error
	DeleteTask(ctx context.Context, taskCode int64) error
	DeleteTaskByPcode(ctx context.Context, pcode int64) error
	FindTaskByProjectCodeAndDeleted(ctx context.Context, projectCode int64, deleted int, page int64, size int64) (list []*data.Task, total int64, err error)
	BatchUpdateDeleted(ctx context.Context, taskCodes []int64, deleted int) error
	BatchUpdateAssignTo(ctx context.Context, taskCodes []int64, assignTo int64) error
	BatchSaveTaskMember(ctx context.Context, members []*data.TaskMember) error
	UpdateTaskFields(ctx context.Context, taskCode int64, updates map[string]interface{}) error
}
