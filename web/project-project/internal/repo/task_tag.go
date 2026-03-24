package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type TaskTagRepo interface {
	FindTagsByProjectCode(ctx context.Context, projectCode int64) ([]*data.TaskTag, error)
	FindTagByCode(ctx context.Context, code string) (*data.TaskTag, error)
	FindTagByNameAndProject(ctx context.Context, projectCode int64, name string) (*data.TaskTag, error)
	SaveTag(ctx context.Context, tag *data.TaskTag) error
	UpdateTag(ctx context.Context, tag *data.TaskTag) error
	DeleteTag(ctx context.Context, code string) error
	FindToTagsByTaskCode(ctx context.Context, taskCode int64) ([]*data.TaskToTag, error)
	FindToTagByTaskAndTag(ctx context.Context, taskCode int64, tagCode int64) (*data.TaskToTag, error)
	SaveToTag(ctx context.Context, tt *data.TaskToTag) error
	DeleteToTag(ctx context.Context, id int64) error
	FindTaskCodesByTagCode(ctx context.Context, tagCode int64) ([]int64, error)
}
