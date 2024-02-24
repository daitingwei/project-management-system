package repo

import (
	"context"
	"gorm.io/gorm"

	"test.com/project-user/internal/data/organization"
)

type OrganizationRepo interface {
	SaveOrganization(tx *gorm.DB, ctx context.Context, org *organization.Organization) error
	FindOrganizationByMemId(ctx context.Context, memId int64) ([]*organization.Organization, error)
	FindOrganizationById(ctx context.Context, id int64) (*organization.Organization, error)
	UpdateOrganization(ctx context.Context, org *organization.Organization) error
}
