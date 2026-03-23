package dao

import (
	"context"
	"gorm.io/gorm"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorms"
)

type TaskStagesDao struct {
	conn *gorms.GormConn
}

func (t *TaskStagesDao) FindById(ctx context.Context, id int) (ts *data.TaskStages, err error) {
	err = t.conn.Session(ctx).Where("id=?", id).Find(&ts).Error
	return
}

func (t *TaskStagesDao) FindStagesByProjectId(
	ctx context.Context,
	projectCode int64,
	page int64,
	pageSize int64) (list []*data.TaskStages, total int64, err error) {
	session := t.conn.Session(ctx)
	err = session.Model(&data.TaskStages{}).
		Where("project_code=?", projectCode).
		Order("sort asc").
		Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	err = session.Model(&data.TaskStages{}).
		Where("project_code=?", projectCode).
		Count(&total).Error
	return
}

func (t *TaskStagesDao) FindByProjectId(ctx context.Context, projectCode int64) (list []*data.TaskStages, err error) {
	err = t.conn.Session(ctx).Model(&data.TaskStages{}).
		Where("project_code=? and deleted=0", projectCode).
		Order("sort asc").
		Find(&list).Error
	return
}

func (t *TaskStagesDao) FindByStageCodeLtSort(ctx context.Context, projectCode int64, sort int) (ts *data.TaskStages, err error) {
	err = t.conn.Session(ctx).Model(&data.TaskStages{}).
		Where("project_code=? and sort < ? and deleted=0", projectCode, sort).
		Order("sort desc").
		Limit(1).
		Find(&ts).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (t *TaskStagesDao) UpdateTaskStagesSort(ctx context.Context, conn database.DbConn, ts *data.TaskStages) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).
		Where("id=?", ts.Id).
		Select("sort", "name").
		Updates(ts).Error
	return err
}

func (t *TaskStagesDao) SaveTaskStages(ctx context.Context, conn database.DbConn, ts *data.TaskStages) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).Save(&ts).Error
	return err
}

func (t *TaskStagesDao) FindMaxSortByProjectCode(ctx context.Context, projectCode int64) (maxSort int, err error) {
	var result struct {
		MaxSort int
	}
	err = t.conn.Session(ctx).Model(&data.TaskStages{}).
		Where("project_code = ?", projectCode).
		Select("MAX(sort) as max_sort").
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.MaxSort, nil
}

func NewTaskStagesDao() *TaskStagesDao {
	return &TaskStagesDao{
		conn: gorms.New(),
	}
}

func (t *TaskStagesDao) DeleteTaskStages(ctx context.Context, conn database.DbConn, id int) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).Where("id=?", id).Delete(&data.TaskStages{}).Error
	return err
}
