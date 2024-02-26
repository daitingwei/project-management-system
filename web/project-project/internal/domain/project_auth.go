package domain

import (
	"context"
	"go.uber.org/zap"
	"strconv"
	"test.com/project-common/errs"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/repo"
	"test.com/project-project/pkg/model"
	"time"
)

type ProjectAuthDomain struct {
	projectAuthRepo       repo.ProjectAuthRepo
	userRpcDomain         *UserRpcDomain
	projectNodeDomain     *ProjectNodeDomain
	projectAuthNodeDomain *ProjectAuthNodeDomain
	accountDomain         *AccountDomain
}

func (d *ProjectAuthDomain) AuthList(orgCode int64) ([]*data.ProjectAuthDisplay, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list, err := d.projectAuthRepo.FindAuthList(c, orgCode)
	if err != nil {
		zap.L().Error("project AuthList projectAuthRepo.FindAuthList error", zap.Error(err))
		return nil, model.DBError
	}
	var pdList []*data.ProjectAuthDisplay
	for _, v := range list {
		display := v.ToDisplay()
		pdList = append(pdList, display)
	}
	return pdList, nil
}

func (d *ProjectAuthDomain) AuthListPage(orgCode int64, page int64, pageSize int64) ([]*data.ProjectAuthDisplay, int64, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list, total, err := d.projectAuthRepo.FindAuthListPage(c, orgCode, page, pageSize)
	if err != nil {
		zap.L().Error("project AuthList projectAuthRepo.FindAuthList error", zap.Error(err))
		return nil, 0, model.DBError
	}
	var pdList []*data.ProjectAuthDisplay
	for _, v := range list {
		display := v.ToDisplay()
		pdList = append(pdList, display)
	}
	return pdList, total, nil
}

func (d *ProjectAuthDomain) AllNodeAndAuth(authId int64) ([]*data.ProjectNodeAuthTree, []string, *errs.BError) {
	nodeList, err := d.projectNodeDomain.NodeList()
	if err != nil {
		return nil, nil, err
	}
	checkedList, err := d.projectAuthNodeDomain.AuthNodeList(authId)
	if err != nil {
		return nil, nil, err
	}
	list := data.ToAuthNodeTreeList(nodeList, checkedList)
	return list, checkedList, nil
}

func (d *ProjectAuthDomain) Save(conn database.DbConn, authId int64, nodes []string) *errs.BError {
	err := d.projectAuthNodeDomain.Save(conn, authId, nodes)
	if err != nil {
		return err
	}
	return nil
}

func (d *ProjectAuthDomain) AuthNodes(memberId int64) ([]string, *errs.BError) {
	account, err := d.accountDomain.FindAccount(memberId)
	if err != nil {
		zap.L().Error("project AuthNodes accountDomain.FindAccount error", zap.Error(err))
		return nil, err
	}
	if account == nil {
		return nil, model.ParamsError
	}
	// 如果是 owner，直接返回所有权限节点（拥有所有权限）
	if account.IsOwner == 1 {
		allNodes, err := d.projectAuthNodeDomain.AllNodes()
		if err != nil {
			zap.L().Error("project AuthNodes projectAuthNodeDomain.AllNodes error", zap.Error(err))
			return nil, model.DBError
		}
		return allNodes, nil
	}
	// 如果没有设置 authorize，返回空（需要通过其他方式赋予权限）
	authorize := account.Authorize
	if authorize == "" {
		return []string{}, nil
	}
	authId, parseErr := strconv.ParseInt(authorize, 10, 64)
	if parseErr != nil {
		zap.L().Error("project AuthNodes strconv.ParseInt error", zap.Error(parseErr))
		return nil, model.ParamsError
	}
	authNodeList, dbErr := d.projectAuthNodeDomain.AuthNodeList(authId)
	if dbErr != nil {
		zap.L().Error("project AuthNodes projectAuthNodeDomain.AuthNodeList error", zap.Error(dbErr))
		return nil, model.DBError
	}
	return authNodeList, nil
}

func NewProjectAuthDomain() *ProjectAuthDomain {
	return &ProjectAuthDomain{
		projectAuthRepo:       dao.NewProjectAuthDao(),
		userRpcDomain:         NewUserRpcDomain(),
		projectNodeDomain:     NewProjectNodeDomain(),
		projectAuthNodeDomain: NewProjectAuthNodeDomain(),
		accountDomain:         NewAccountDomain(),
	}
}

func (d *ProjectAuthDomain) SaveAuth(pa *data.ProjectAuth) *errs.BError {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := d.projectAuthRepo.SaveAuth(c, pa); err != nil {
		return model.DBError
	}
	return nil
}

func (d *ProjectAuthDomain) UpdateAuth(pa *data.ProjectAuth) *errs.BError {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := d.projectAuthRepo.UpdateAuth(c, pa); err != nil {
		return model.DBError
	}
	return nil
}

func (d *ProjectAuthDomain) SetDefault(orgCode int64, authId int64) *errs.BError {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := d.projectAuthRepo.SetDefault(c, orgCode, authId); err != nil {
		return model.DBError
	}
	return nil
}

func (d *ProjectAuthDomain) FindDefaultAuth(orgCode int64) (*data.ProjectAuth, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	auth, err := d.projectAuthRepo.FindDefaultAuth(c, orgCode)
	if err != nil {
		return nil, model.DBError
	}
	return auth, nil
}

func (d *ProjectAuthDomain) DeleteAuth(authId int64) *errs.BError {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := d.projectAuthRepo.DeleteAuth(c, authId); err != nil {
		return model.DBError
	}
	return nil
}
