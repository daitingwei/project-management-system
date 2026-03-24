package dao

import (
	"context"
	"time"

	"test.com/project-common/encrypts"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type TaskTagDao struct {
	conn *gorms.GormConn
}

func NewTaskTagDao() *TaskTagDao {
	return &TaskTagDao{conn: gorms.New()}
}

func (d *TaskTagDao) FindTagsByProjectCode(ctx context.Context, projectCode int64) ([]*data.TaskTag, error) {
	session := d.conn.Session(ctx)
	var list []*data.TaskTag
	err := session.Where("project_code=?", projectCode).Find(&list).Error
	return list, err
}

func (d *TaskTagDao) FindTagByCode(ctx context.Context, code string) (*data.TaskTag, error) {
	session := d.conn.Session(ctx)
	tagId := encrypts.DecryptNoErr(code)
	var tag data.TaskTag
	err := session.Where("id=?", tagId).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (d *TaskTagDao) FindTagByNameAndProject(ctx context.Context, projectCode int64, name string) (*data.TaskTag, error) {
	session := d.conn.Session(ctx)
	var tag data.TaskTag
	err := session.Where("project_code=? and name=?", projectCode, name).First(&tag).Error
	if err != nil {
		return nil, nil // not found is ok
	}
	return &tag, nil
}

func (d *TaskTagDao) SaveTag(ctx context.Context, tag *data.TaskTag) error {
	session := d.conn.Session(ctx)
	tag.CreateTime = time.Now().UnixMilli()
	return session.Save(tag).Error
}

func (d *TaskTagDao) UpdateTag(ctx context.Context, tag *data.TaskTag) error {
	session := d.conn.Session(ctx)
	return session.Model(&data.TaskTag{}).Where("id=?", tag.Id).Updates(map[string]interface{}{
		"name":  tag.Name,
		"color": tag.Color,
	}).Error
}

func (d *TaskTagDao) DeleteTag(ctx context.Context, code string) error {
	session := d.conn.Session(ctx)
	tagId := encrypts.DecryptNoErr(code)
	return session.Where("id=?", tagId).Delete(&data.TaskTag{}).Error
}

func (d *TaskTagDao) FindToTagsByTaskCode(ctx context.Context, taskCode int64) ([]*data.TaskToTag, error) {
	session := d.conn.Session(ctx)
	var list []*data.TaskToTag
	err := session.Where("task_code=?", taskCode).Find(&list).Error
	return list, err
}

func (d *TaskTagDao) FindToTagByTaskAndTag(ctx context.Context, taskCode int64, tagCode int64) (*data.TaskToTag, error) {
	session := d.conn.Session(ctx)
	var tt data.TaskToTag
	err := session.Where("task_code=? and tag_code=?", taskCode, tagCode).First(&tt).Error
	if err != nil {
		return nil, nil
	}
	return &tt, nil
}

func (d *TaskTagDao) SaveToTag(ctx context.Context, tt *data.TaskToTag) error {
	session := d.conn.Session(ctx)
	tt.CreateTime = time.Now().UnixMilli()
	return session.Save(tt).Error
}

func (d *TaskTagDao) DeleteToTag(ctx context.Context, id int64) error {
	session := d.conn.Session(ctx)
	return session.Where("id=?", id).Delete(&data.TaskToTag{}).Error
}

func (d *TaskTagDao) FindTaskCodesByTagCode(ctx context.Context, tagCode int64) ([]int64, error) {
	session := d.conn.Session(ctx)
	var list []*data.TaskToTag
	err := session.Where("tag_code=?", tagCode).Find(&list).Error
	if err != nil {
		return nil, err
	}
	var codes []int64
	for _, v := range list {
		codes = append(codes, v.TaskCode)
	}
	return codes, nil
}
