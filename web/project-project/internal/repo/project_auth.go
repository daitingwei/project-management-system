package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type ProjectAuthRepo interface {
	FindAuthList(ctx context.Context, orgCode int64) (list []*data.ProjectAuth, err error)
	FindAuthListPage(ctx context.Context, orgCode int64, page int64, pageSize int64) (list []*data.ProjectAuth, total int64, err error)
	FindDefaultAuth(ctx context.Context, orgCode int64) (*data.ProjectAuth, error)
	SaveAuth(ctx context.Context, pa *data.ProjectAuth) error
	UpdateAuth(ctx context.Context, pa *data.ProjectAuth) error
	SetDefault(ctx context.Context, orgCode int64, authId int64) error
	DeleteAuth(ctx context.Context, authId int64) error
}
