package repo

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
)

type FileRepo interface {
	// Save(ctx context.Context, file *data.File) error
	SaveTx(ctx context.Context, conn database.DbConn, file *data.File) error
	FindByIds(background context.Context, ids []int64) (list []*data.File, err error)
	FindByProjectCode(ctx context.Context, projectCode int64, deleted int, page int64, pageSize int64) (list []*data.File, total int64, err error)
	FindById(ctx context.Context, id int64) (*data.File, error)
	UpdateTitle(ctx context.Context, id int64, title string) error
	UpdateDeleted(ctx context.Context, id int64, deleted int) error
	DeleteById(ctx context.Context, id int64) error
}
