package task_service_v1

import (
	"context"
	"time"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/min"
	"test.com/project-common/tms"
	"test.com/project-grpc/task"
	"test.com/project-grpc/user/login"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/domain"
	"test.com/project-project/internal/interceptor"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/rpc"
	"test.com/project-project/pkg/model"
)

type TaskService struct {
	task.UnimplementedTaskServiceServer
	cache                  repo.Cache
	transaction            tran.Transaction
	projectRepo            repo.ProjectRepo
	projectTemplateRepo    repo.ProjectTemplateRepo
	taskStagesTemplateRepo repo.TaskStagesTemplateRepo
	taskStagesRepo         repo.TaskStagesRepo
	taskRepo               repo.TaskRepo
	projectLogRepo         repo.ProjectLogRepo
	taskWorkTimeRepo       repo.TaskWorkTimeRepo
	fileRepo               repo.FileRepo
	sourceLinkRepo         repo.SourceLinkRepo
	taskWorkTimeDomain     *domain.TaskWorkTimeDomain
	taskTagRepo            repo.TaskTagRepo
	departmentRepo         repo.DepartmentRepo
}

func New() *TaskService {
	return &TaskService{
		cache:                  dao.Rc,
		transaction:            dao.NewTransaction(),
		projectRepo:            dao.NewProjectDao(),
		projectTemplateRepo:    dao.NewProjectTemplateDao(),
		taskStagesTemplateRepo: dao.NewTaskStagesTemplateDao(),
		taskStagesRepo:         dao.NewTaskStagesDao(),
		taskRepo:               dao.NewTaskDao(),
		projectLogRepo:         dao.NewProjectLogDao(),
		taskWorkTimeRepo:       dao.NewTaskWorkTimeDao(),
		fileRepo:               dao.NewFileDao(),
		sourceLinkRepo:         dao.NewSourceLinkDao(),
		taskWorkTimeDomain:     domain.NewTaskWorkTimeDomain(),
		taskTagRepo:            dao.NewTaskTagDao(),
		departmentRepo:         dao.NewDepartmentDao(),
	}
}

func (t *TaskService) TaskStages(co context.Context, msg *task.TaskReqMessage) (*task.TaskStagesResponse, error) {
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	page := msg.Page
	pageSize := msg.PageSize
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stages, total, err := t.taskStagesRepo.FindStagesByProjectId(ctx, projectCode, page, pageSize)
	if err != nil {
		zap.L().Error("project SaveProject taskStagesRepo.FindStagesByProjectId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	var tsMessages []*task.TaskStagesMessage
	copier.Copy(&tsMessages, stages)
	if tsMessages == nil {
		return &task.TaskStagesResponse{List: tsMessages, Total: 0}, nil
	}
	stagesMap := data.ToTaskStagesMap(stages)
	for _, v := range tsMessages {
		taskStages := stagesMap[int(v.Id)]
		v.Code = encrypts.EncryptNoErr(int64(v.Id))
		v.CreateTime = tms.FormatByMill(taskStages.CreateTime)
		v.ProjectCode = msg.ProjectCode
	}
	return &task.TaskStagesResponse{List: tsMessages, Total: total}, nil
}

func (t *TaskService) SaveTaskStages(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskStagesResponse, error) {
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)

	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	maxSort, err := t.taskStagesRepo.FindMaxSortByProjectCode(c, projectCode)
	if err != nil {
		zap.L().Error("SaveTaskStages FindMaxSortByProjectCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	ts := &data.TaskStages{
		Name:        msg.Name,
		ProjectCode: projectCode,
		Sort:        maxSort + 1,
		Description: "",
		CreateTime:  time.Now().UnixMilli(),
		Deleted:     model.NoDeleted,
	}
	err = t.transaction.Action(func(conn database.DbConn) error {
		return t.taskStagesRepo.SaveTaskStages(ctx, conn, ts)
	})
	if err != nil {
		zap.L().Error("SaveTaskStages error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskStagesResponse{}, nil
}

func (t *TaskService) MemberProjectList(co context.Context, msg *task.TaskReqMessage) (*task.MemberProjectResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	projectMembers, total, err := t.projectRepo.FindProjectMemberByPid(ctx, projectCode)
	if err != nil {
		zap.L().Error("project MemberProjectList projectRepo.FindProjectMemberByPid error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	// 2.拿上用户id列表 去请求用户信息
	if projectMembers == nil || len(projectMembers) <= 0 {
		return &task.MemberProjectResponse{List: nil, Total: 0}, nil
	}
	var mIds []int64
	pmMap := make(map[int64]*data.ProjectMember)
	for _, v := range projectMembers {
		mIds = append(mIds, v.MemberCode)
		pmMap[v.MemberCode] = v
	}
	userMsg := &login.UserMessage{
		MIds: mIds,
	}
	zap.L().Info("MemberProjectList rpc call", zap.Any("userMsg", userMsg))
	memberMessageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, userMsg)
	zap.L().Info("MemberProjectList rpc result", zap.Any("memberMessageList", memberMessageList), zap.Error(err))
	if err != nil {
		zap.L().Error("project MemberProjectList LoginServiceClient.FindMemInfoByIds error", zap.Error(err))
	}
	
	memberMap := make(map[int64]*login.MemberMessage)
	if memberMessageList != nil && len(memberMessageList.List) > 0 {
		for _, v := range memberMessageList.List {
			memberMap[v.Id] = v
		}
	}
	
	if len(memberMap) < len(mIds) {
		accounts, err := t.departmentRepo.FindMemberAccountsByIds(ctx, mIds)
		if err == nil && len(accounts) > 0 {
			for _, acc := range accounts {
				if _, ok := memberMap[acc.MemberCode]; !ok {
					memberMap[acc.MemberCode] = &login.MemberMessage{
						Id:       acc.MemberCode,
						Name:     acc.Name,
						Avatar:   acc.Avatar,
						Email:    acc.Email,
						Code:     encrypts.EncryptNoErr(acc.MemberCode),
						Status:   int32(acc.Status),
					}
				}
			}
		}
	}

	var list []*task.MemberProjectMessage
	for _, pm := range projectMembers {
		v, ok := memberMap[pm.MemberCode]
		if !ok {
			continue
		}
		mpm := &task.MemberProjectMessage{
			MemberCode: pm.MemberCode,
			Name:       v.Name,
			Avatar:     v.Avatar,
			Email:      v.Email,
			Code:       v.Code,
			IsOwner:    int32(pm.IsOwner),
		}
		list = append(list, mpm)
	}
	return &task.MemberProjectResponse{List: list, Total: total}, nil
}
func (t *TaskService) TaskList(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskListResponse, error) {
	stageCode := encrypts.DecryptNoErr(msg.StageCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	taskList, err := t.taskRepo.FindTaskByStageCode(c, int(stageCode))
	if err != nil {
		zap.L().Error("project task TaskList FindTaskByStageCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var taskDisplayList []*data.TaskDisplay
	var mIds []int64
	for _, v := range taskList {
		display := v.ToTaskDisplay()
		if v.Private == 1 {
			// 代表隐私模式
			taskMember, err := t.taskRepo.FindTaskMemberByTaskId(ctx, v.Id, msg.MemberId)
			if err != nil {
				zap.L().Error("project task TaskList taskRepo.FindTaskMemberByTaskId error", zap.Error(err))
				return nil, errs.GrpcError(model.DBError)
			}
			if taskMember != nil {
				display.CanRead = model.CanRead
			} else {
				display.CanRead = model.NoCanRead
			}
		}
		taskDisplayList = append(taskDisplayList, display)
		mIds = append(mIds, v.AssignTo)
	}
	if mIds == nil || len(mIds) <= 0 {
		return &task.TaskListResponse{List: nil}, nil
	}
	// in ()
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mIds})
	if err != nil {
		zap.L().Error("project task TaskList LoginServiceClient.FindMemInfoByIds error", zap.Error(err))
		return nil, err
	}
	memberMap := make(map[int64]*login.MemberMessage)
	for _, v := range messageList.List {
		memberMap[v.Id] = v
	}
	for _, v := range taskDisplayList {
		message := memberMap[encrypts.DecryptNoErr(v.AssignTo)]
		e := data.Executor{}
		if message != nil {
			e.Name = message.Name
			e.Avatar = message.Avatar
		}
		v.Executor = e
	}
	var taskMessageList []*task.TaskMessage
	copier.Copy(&taskMessageList, taskDisplayList)
	return &task.TaskListResponse{List: taskMessageList}, nil
}
func (t *TaskService) SaveTask(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMessage, error) {
	// 1. 检查业务逻辑
	if msg.Name == "" {
		return nil, errs.GrpcError(model.TaskNameNotNull)
	}
	stageCode := encrypts.DecryptNoErr(msg.StageCode)
	taskStages, err := t.taskStagesRepo.FindById(ctx, int(stageCode))
	if err != nil {
		zap.L().Error("project task SaveTask taskStagesRepo.FindById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if taskStages == nil {
		return nil, errs.GrpcError(model.TaskStagesNotNull)
	}
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	project, err := t.projectRepo.FindProjectById(ctx, projectCode)
	if err != nil {
		zap.L().Error("project task SaveTask projectRepo.FindProjectById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if project == nil || project.Deleted == model.Deleted {
		return nil, errs.GrpcError(model.ProjectAlreadyDeleted)
	}
	maxIdNum, err := t.taskRepo.FindTaskMaxIdNum(ctx, projectCode)
	if err != nil {
		zap.L().Error("project task SaveTask taskRepo.FindTaskMaxIdNum error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if maxIdNum == nil {
		a := 0
		maxIdNum = &a
	}
	maxSort, err := t.taskRepo.FindTaskSort(ctx, projectCode, stageCode)
	if err != nil {
		zap.L().Error("project task SaveTask taskRepo.FindTaskSort error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if maxSort == nil {
		a := 0
		maxSort = &a
	}
	assignTo := encrypts.DecryptNoErr(msg.AssignTo)
	ts := &data.Task{
		Name:        msg.Name,
		CreateTime:  time.Now().UnixMilli(),
		CreateBy:    msg.MemberId,
		AssignTo:    assignTo,
		ProjectCode: projectCode,
		StageCode:   int(stageCode),
		IdNum:       *maxIdNum + 1,
		Private:     project.OpenTaskPrivate,
		Sort:        *maxSort + 65536,
		BeginTime:   time.Now().UnixMilli(),
		EndTime:     time.Now().Add(2 * 24 * time.Hour).UnixMilli(),
	}
	err = t.transaction.Action(func(conn database.DbConn) error {
		err = t.taskRepo.SaveTask(ctx, conn, ts)
		if err != nil {
			zap.L().Error("project task SaveTask taskRepo.SaveTask error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}

		tm := &data.TaskMember{
			MemberCode: assignTo,
			TaskCode:   ts.Id,
			JoinTime:   time.Now().UnixMilli(),
			IsOwner:    model.Owner,
		}
		if assignTo == msg.MemberId {
			tm.IsExecutor = model.Executor
		}
		err = t.taskRepo.SaveTaskMember(ctx, conn, tm)
		if err != nil {
			zap.L().Error("project task SaveTask taskRepo.SaveTaskMember error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	display := ts.ToTaskDisplay()
	member, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: assignTo})
	if err != nil {
		return nil, err
	}
	display.Executor = data.Executor{
		Name:   member.Name,
		Avatar: member.Avatar,
		Code:   member.Code,
	}
	// 添加任务动态
	// TODO: 这里"task"写死了，需要改为动态传入以支持其他模块（如需求、文档等）
	createProjectLog(t.projectLogRepo, ts.ProjectCode, ts.Id, ts.Name, ts.AssignTo, "create", "task")

	tm := &task.TaskMessage{}
	copier.Copy(tm, display)
	interceptor.ClearTaskListCache()
	return tm, nil
}

func createProjectLog(
	logRepo repo.ProjectLogRepo,
	projectCode int64,
	taskCode int64,
	taskName string,
	toMemberCode int64,
	logType string,
	actionType string) {
	remark := ""
	if logType == "create" {
		remark = "创建了任务"
	}
	pl := &data.ProjectLog{
		MemberCode:  toMemberCode,
		SourceCode:  taskCode,
		Content:     taskName,
		Remark:      remark,
		ProjectCode: projectCode,
		CreateTime:  time.Now().UnixMilli(),
		Type:        logType,
		ActionType:  actionType,
		Icon:        "plus",
		IsComment:   0,
		IsRobot:     0,
	}
	logRepo.SaveProjectLog(pl)
}

func (t *TaskService) TaskSort(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskSortResponse, error) {
	preTaskCode := encrypts.DecryptNoErr(msg.PreTaskCode)
	toStageCode := encrypts.DecryptNoErr(msg.ToStageCode)
	if msg.PreTaskCode == msg.NextTaskCode {
		return &task.TaskSortResponse{}, nil
	}
	err := t.sortTask(preTaskCode, msg.NextTaskCode, toStageCode)
	if err != nil {
		return nil, err
	}
	zap.L().Info("TaskSort 完成，清除所有任务缓存")
	interceptor.ClearTaskListCache()
	return &task.TaskSortResponse{}, nil
}

func (t *TaskService) TaskDone(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskSortResponse, error) {
	zap.L().Info("TaskDone 被调用", zap.String("taskCode", msg.TaskCode), zap.Int32("done", msg.Done))
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	done := msg.Done
	zap.L().Info("解密后的 taskCode", zap.Int64("taskCode", taskCode), zap.Int("done", int(done)))
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := t.taskRepo.UpdateTaskDone(c, taskCode, int(done))
	if err != nil {
		zap.L().Error("project task TaskDone error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	taskInfo, _ := t.taskRepo.FindTaskById(c, taskCode)
	if taskInfo != nil {
		interceptor.ClearTaskListCache()
		createProjectLog(t.projectLogRepo, taskInfo.ProjectCode, taskCode, taskInfo.Name, msg.MemberId, "update", "taskDone")
	}
	return &task.TaskSortResponse{}, nil
}

func (t *TaskService) sortTask(preTaskCode int64, nextTaskCode string, toStageCode int64) error {
	// 1. 从小到大排
	// 2. 原有的顺序  比如 1 2 3 4 5 4排到2前面去 4的序号在1和2 之间 如果4是最后一个 保证 4比所有的序号都打 如果 排到第一位 直接置为0

	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ts, err := t.taskRepo.FindTaskById(c, preTaskCode)
	if err != nil {
		zap.L().Error("project task TaskSort taskRepo.FindTaskById error", zap.Error(err))
		return errs.GrpcError(model.DBError)
	}
	err = t.transaction.Action(func(conn database.DbConn) error {
		// 如果相等是不需要进行改变的
		ts.StageCode = int(toStageCode)
		if nextTaskCode != "" {
			// 意味着要进行排序的替换
			nextTaskCode := encrypts.DecryptNoErr(nextTaskCode)
			next, err := t.taskRepo.FindTaskById(c, nextTaskCode)
			if err != nil {
				zap.L().Error("project task TaskSort taskRepo.FindTaskById error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			// next.Sort 要找到比它小的那个任务
			prepre, err := t.taskRepo.FindTaskByStageCodeLtSort(c, next.StageCode, next.Sort)
			if err != nil {
				zap.L().Error("project task TaskSort taskRepo.FindTaskByStageCodeLtSort error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			if prepre != nil {
				ts.Sort = (prepre.Sort + next.Sort) / 2
			}
			if prepre == nil {
				ts.Sort = 0
			}
			// sort := ts.Sort
			// ts.Sort = next.Sort
			// next.Sort = sort
			// err = t.taskRepo.UpdateTaskSort(c, conn, next)
			// if err != nil {
			//	zap.L().Error("project task TaskSort taskRepo.UpdateTaskSort error", zap.Error(err))
			//	return errs.GrpcError(model.DBError)
			// }
		} else {
			maxSort, err := t.taskRepo.FindTaskSort(c, ts.ProjectCode, int64(ts.StageCode))
			if err != nil {
				zap.L().Error("project task TaskSort taskRepo.FindTaskSort error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			if maxSort == nil {
				a := 0
				maxSort = &a
			}
			ts.Sort = *maxSort + 65536
		}
		if ts.Sort < 50 {
			// 重置排序
			err = t.resetSort(toStageCode)
			if err != nil {
				zap.L().Error("project task TaskSort resetSort error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			return t.sortTask(preTaskCode, nextTaskCode, toStageCode)
		}
		err = t.taskRepo.UpdateTaskSort(c, conn, ts)
		if err != nil {
			zap.L().Error("project task TaskSort taskRepo.UpdateTaskSort error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})
	return err
}

func (t *TaskService) resetSort(stageCode int64) error {
	list, err := t.taskRepo.FindTaskByStageCode(context.Background(), int(stageCode))
	if err != nil {
		return err
	}
	return t.transaction.Action(func(conn database.DbConn) error {
		iSort := 65536
		for index, v := range list {
			v.Sort = (index + 1) * iSort
			return t.taskRepo.UpdateTaskSort(context.Background(), conn, v)
		}
		return nil
	})

}

func (t *TaskService) MyTaskList(ctx context.Context, msg *task.TaskReqMessage) (*task.MyTaskListResponse, error) {
	var tsList []*data.Task
	var err error
	var total int64
	if msg.TaskType == 1 {
		// 我执行的
		tsList, total, err = t.taskRepo.FindTaskByAssignTo(ctx, msg.MemberId, int(msg.Type), msg.Page, msg.PageSize)
		if err != nil {
			zap.L().Error("project task MyTaskList taskRepo.FindTaskByAssignTo error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if msg.TaskType == 2 {
		// 我参与的
		tsList, total, err = t.taskRepo.FindTaskByMemberCode(ctx, msg.MemberId, int(msg.Type), msg.Page, msg.PageSize)
		if err != nil {
			zap.L().Error("project task MyTaskList taskRepo.FindTaskByMemberCode error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if msg.TaskType == 3 {
		// 我创建的
		tsList, total, err = t.taskRepo.FindTaskByCreateBy(ctx, msg.MemberId, int(msg.Type), msg.Page, msg.PageSize)
		if err != nil {
			zap.L().Error("project task MyTaskList taskRepo.FindTaskByCreateBy error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if tsList == nil || len(tsList) <= 0 {
		return &task.MyTaskListResponse{List: nil, Total: 0}, nil
	}
	var pids []int64
	var mids []int64
	for _, v := range tsList {
		pids = append(pids, v.ProjectCode)
		mids = append(mids, v.AssignTo)
	}
	pListChan := make(chan []*data.Project)
	defer close(pListChan)
	mListChan := make(chan *login.MemberMessageList)
	defer close(mListChan)
	// 1.
	go func() {
		pList, _ := t.projectRepo.FindProjectByIds(ctx, pids)
		pListChan <- pList
	}()
	go func() {
		mList, _ := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{
			// 2.  1,2这两个请求无关联性  go+channel
			MIds: mids,
		})
		mListChan <- mList
	}()
	pList := <-pListChan
	projectMap := data.ToProjectMap(pList)
	mList := <-mListChan
	mMap := make(map[int64]*login.MemberMessage)
	for _, v := range mList.List {
		mMap[v.Id] = v
	}
	var mtdList []*data.MyTaskDisplay
	for _, v := range tsList {
		var name, avatar string
		if memberMessage := mMap[v.AssignTo]; memberMessage != nil {
			name = memberMessage.Name
			avatar = memberMessage.Avatar
		}
		mtd := v.ToMyTaskDisplay(projectMap[v.ProjectCode], name, avatar)
		mtdList = append(mtdList, mtd)
	}
	var myMsgs []*task.MyTaskMessage
	copier.Copy(&myMsgs, mtdList)
	return &task.MyTaskListResponse{List: myMsgs, Total: total}, nil
}

func (t *TaskService) ReadTask(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMessage, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	taskInfo, err := t.taskRepo.FindTaskById(c, taskCode)
	if err != nil {
		zap.L().Error("project task ReadTask taskRepo FindTaskById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if taskInfo == nil {
		return &task.TaskMessage{}, nil
	}
	display := taskInfo.ToTaskDisplay()
	if taskInfo.Private == 1 {
		taskMember, err := t.taskRepo.FindTaskMemberByTaskId(ctx, taskInfo.Id, msg.MemberId)
		if err != nil {
			zap.L().Error("project task TaskList taskRepo.FindTaskMemberByTaskId error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		if taskMember != nil {
			display.CanRead = model.CanRead
		} else {
			display.CanRead = model.NoCanRead
		}
	}
	pj, err := t.projectRepo.FindProjectById(c, taskInfo.ProjectCode)
	if err == nil && pj != nil {
		display.ProjectName = pj.Name
	}
	taskStages, err := t.taskStagesRepo.FindById(c, taskInfo.StageCode)
	if err == nil && taskStages != nil {
		display.StageName = taskStages.Name
	}
	if taskInfo.AssignTo != 0 {
		memberMessage, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: taskInfo.AssignTo})
		if err != nil {
			zap.L().Error("project task TaskList LoginServiceClient.FindMemInfoById error", zap.Error(err))
		} else {
			display.Executor = data.Executor{
				Name:   memberMessage.Name,
				Avatar: memberMessage.Avatar,
			}
		}
	}
	var taskMessage = &task.TaskMessage{}
	copier.Copy(taskMessage, display)
	return taskMessage, nil
}
func (t *TaskService) ListTaskMember(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMemberList, error) {
	// 查询 task member表 根据memberCode去查询用户信息
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	taskMemberPage, total, err := t.taskRepo.FindTaskMemberPage(c, taskCode, msg.Page, msg.PageSize)
	if err != nil {
		zap.L().Error("project task TaskList taskRepo.FindTaskMemberPage error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if len(taskMemberPage) == 0 {
		return &task.TaskMemberList{List: []*task.TaskMemberMessage{}, Total: 0}, nil
	}

	// 获取项目信息以获取组织code
	taskInfo, err := t.taskRepo.FindTaskById(c, taskCode)
	if err != nil || taskInfo == nil {
		return &task.TaskMemberList{List: []*task.TaskMemberMessage{}, Total: 0}, nil
	}
	projectInfo, err := t.projectRepo.FindProjectById(c, taskInfo.ProjectCode)
	if err != nil || projectInfo == nil {
		return &task.TaskMemberList{List: []*task.TaskMemberMessage{}, Total: 0}, nil
	}
	organizationCode := projectInfo.OrganizationCode

	var mids []int64
	for _, v := range taskMemberPage {
		mids = append(mids, v.MemberCode)
	}
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mids})
	if err != nil {
		zap.L().Error("project task ListTaskMember FindMemInfoByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	mMap := make(map[int64]*login.MemberMessage, len(messageList.List))
	for _, v := range messageList.List {
		mMap[v.Id] = v
	}

	// 获取组织内所有memberAccount，用于查询membarAccountCode
	accounts, err := t.departmentRepo.FindMemberAccountsByOrg(c, organizationCode)
	if err != nil {
		zap.L().Error("project task ListTaskMember FindMemberAccountsByOrg error", zap.Error(err))
	}
	accountMap := make(map[int64]*data.MemberAccount)
	for _, a := range accounts {
		accountMap[a.MemberCode] = a
	}

	var taskMemeberMemssages []*task.TaskMemberMessage
	for _, v := range taskMemberPage {
		tm := &task.TaskMemberMessage{}
		tm.Code = encrypts.EncryptNoErr(v.MemberCode)
		tm.Id = v.Id
		if message, ok := mMap[v.MemberCode]; ok {
			tm.Name = message.Name
			tm.Avatar = message.Avatar
		}
		// 设置 membarAccountCode
		if account, ok := accountMap[v.MemberCode]; ok {
			tm.MembarAccountCode = encrypts.EncryptNoErr(account.Id)
		}
		tm.IsExecutor = int32(v.IsExecutor)
		tm.IsOwner = int32(v.IsOwner)
		taskMemeberMemssages = append(taskMemeberMemssages, tm)
	}
	return &task.TaskMemberList{List: taskMemeberMemssages, Total: total}, nil
}
func (t *TaskService) TaskLog(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskLogList, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	all := msg.All
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var list []*data.ProjectLog
	var total int64
	var err error
	if all == 1 {
		// 显示全部
		list, total, err = t.projectLogRepo.FindLogByTaskCode(c, taskCode, int(msg.Comment))
	}
	if all == 0 {
		// 分页
		list, total, err = t.projectLogRepo.FindLogByTaskCodePage(c, taskCode, int(msg.Comment), int(msg.Page), int(msg.PageSize))
	}
	if err != nil {
		zap.L().Error("project task TaskLog projectLogRepo.FindLogByTaskCodePage error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if total == 0 {
		return &task.TaskLogList{}, nil
	}
	var displayList []*data.ProjectLogDisplay
	var mIdList []int64
	for _, v := range list {
		mIdList = append(mIdList, v.MemberCode)
	}
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(c, &login.UserMessage{MIds: mIdList})
	mMap := make(map[int64]*login.MemberMessage)
	if err == nil && messageList != nil {
		for _, v := range messageList.List {
			mMap[v.Id] = v
		}
	}
	for _, v := range list {
		display := v.ToDisplay()
		m := data.Member{}
		if message, ok := mMap[v.MemberCode]; ok {
			m.Name = message.Name
			m.Id = message.Id
			m.Avatar = message.Avatar
			m.Code = message.Code
		}
		display.Member = m
		displayList = append(displayList, display)
	}
	var l []*task.TaskLog
	copier.Copy(&l, displayList)
	return &task.TaskLogList{List: l, Total: total}, nil
}

func (t *TaskService) TaskWorkTimeList(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskWorkTimeResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	list, err := t.taskWorkTimeDomain.TaskWorkTimeList(taskCode)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var l []*task.TaskWorkTime
	copier.Copy(&l, list)
	return &task.TaskWorkTimeResponse{List: l, Total: int64(len(l))}, nil
}

func (t *TaskService) SaveTaskWorkTime(ctx context.Context, msg *task.TaskReqMessage) (*task.SaveTaskWorkTimeResponse, error) {
	tmt := &data.TaskWorkTime{}
	tmt.BeginTime = msg.BeginTime
	tmt.Num = int(msg.Num)
	tmt.Content = msg.Content
	tmt.TaskCode = encrypts.DecryptNoErr(msg.TaskCode)
	tmt.MemberCode = msg.MemberId
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := t.taskWorkTimeRepo.Save(c, tmt)
	if err != nil {
		zap.L().Error("project task SaveTaskWorkTime taskWorkTimeRepo.Save error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.SaveTaskWorkTimeResponse{}, nil
}

func (t *TaskService) SaveTaskFile(ctx context.Context, msg *task.TaskFileReqMessage) (*task.TaskFileResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	linkType := msg.LinkType
	if linkType == "" {
		linkType = "task"
	}
	f := &data.File{
		PathName:  msg.PathName,
		Title:     msg.FileName,
		Extension: msg.Extension,
		Size:      int(msg.Size),
		ObjectType:       "",
		OrganizationCode: encrypts.DecryptNoErr(msg.OrganizationCode),
		TaskCode:         encrypts.DecryptNoErr(msg.TaskCode),
		ProjectCode:      encrypts.DecryptNoErr(msg.ProjectCode),
		CreateBy:         msg.MemberId,
		CreateTime:       time.Now().UnixMilli(),
		Downloads:        0,
		Extra:            "",
		Deleted:          model.NoDeleted,
		FileType:         msg.FileType,
		FileUrl:          msg.FileUrl,
		DeletedTime:      0,
	}
	sl := &data.SourceLink{
		SourceType: "file",
		SourceCode: f.Id,
		LinkType:         linkType,
		LinkCode:         taskCode,
		OrganizationCode: encrypts.DecryptNoErr(msg.OrganizationCode),
		CreateBy:         msg.MemberId,
		CreateTime:       time.Now().UnixMilli(),
		Sort:             0,
	}
	err := t.transaction.Action(func(conn database.DbConn) error {
		var err error
		err = t.fileRepo.SaveTx(ctx, conn, f)
		if err != nil {
			zap.L().Error("project task SaveTaskFile fileRepo.SaveTx error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		sl.SourceCode = f.Id
		err = t.sourceLinkRepo.SaveTx(ctx, conn, sl)
		if err != nil {
			zap.L().Error("project task SaveTaskFile sourceLinkRepo.SaveTx error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &task.TaskFileResponse{}, nil
}

func (t *TaskService) TaskSources(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskSourceResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	sourceLinks, err := t.sourceLinkRepo.FindByTaskCode(context.Background(), taskCode)
	if err != nil {
		zap.L().Error("project task TaskSources sourceLinkRepo.FindByTaskCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if len(sourceLinks) == 0 {
		return &task.TaskSourceResponse{}, nil
	}
	var fIdList []int64
	for _, v := range sourceLinks {
		fIdList = append(fIdList, v.SourceCode)
	}
	files, err := t.fileRepo.FindByIds(context.Background(), fIdList)
	if err != nil {
		zap.L().Error("project task SaveTaskFile fileRepo.FindByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	fMap := make(map[int64]*data.File)
	for _, v := range files {
		fMap[v.Id] = v
	}
	var list []*data.SourceLinkDisplay
	for _, v := range sourceLinks {
		list = append(list, v.ToDisplay(fMap[v.SourceCode]))
	}
	var slMsg []*task.TaskSourceMessage
	copier.Copy(&slMsg, list)
	return &task.TaskSourceResponse{List: slMsg}, nil
}

func (t *TaskService) CreateComment(ctx context.Context, msg *task.TaskReqMessage) (*task.CreateCommentResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	taskById, err := t.taskRepo.FindTaskById(context.Background(), taskCode)
	if err != nil {
		zap.L().Error("project task CreateComment fileRepo.FindTaskById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	pl := &data.ProjectLog{
		MemberCode: msg.MemberId,
		Content:    msg.CommentContent,
		Remark:     msg.CommentContent,
		// TODO: 这里"createComment"写死了，需要改为动态传入以支持其他类型的日志
		// 可改为从 msg 中获取日志类型，支持 createComment、updateTask、deleteTask 等
		Type:       "createComment",
		CreateTime: time.Now().UnixMilli(),
		SourceCode: taskCode,
		// TODO: 这里"task"写死了，需要改为动态传入以支持其他模块（如需求、文档等）的评论
		// 可添加 sourceType 参数，根据传入的模块类型动态设置（如 task、requirement、doc 等）
		ActionType:   "task",
		ToMemberCode: 0,
		IsComment:    model.Comment,
		ProjectCode:  taskById.ProjectCode,
		Icon:         "plus",
		IsRobot:      0,
	}
	t.projectLogRepo.SaveProjectLog(pl)
	return &task.CreateCommentResponse{}, nil
}

// EditTask 更新任务字段（只更新 msg 中非零值的字段），然后返回最新任务详情
func (t *TaskService) EditTask(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMessage, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	zap.L().Info("EditTask called", zap.Int64("taskCode", taskCode), zap.Int32("status", msg.Status), zap.Int32("pri", msg.Pri), zap.String("name", msg.Name))
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 构建只包含需要更新字段的 map，避免 GORM Updates(struct) 跳过零值
	updates := make(map[string]interface{})
	if msg.Name != "" {
		updates["name"] = msg.Name
	}
	if msg.Content != "" {
		updates["description"] = msg.Content
	}
	if msg.Pri >= 0 {
		updates["pri"] = msg.Pri
	}
	if msg.Status >= 0 {
		updates["status"] = msg.Status
	}
	if msg.Num != 0 {
		updates["work_time"] = msg.Num
	}
	// BeginTime == math.MinInt64：未传，跳过；== -1：清除置0；> 0：正常设置
	const timeNotSet = int64(-9223372036854775808)
	if msg.BeginTime != timeNotSet {
		if msg.BeginTime > 0 {
			updates["begin_time"] = msg.BeginTime
		} else {
			updates["begin_time"] = 0
		}
	}
	if msg.EndTime != timeNotSet {
		if msg.EndTime > 0 {
			updates["end_time"] = msg.EndTime
		} else {
			updates["end_time"] = 0
		}
	}

	if len(updates) > 0 {
		zap.L().Info("EditTask updates", zap.Any("updates", updates))
		if err := t.taskRepo.UpdateTaskFields(c, taskCode, updates); err != nil {
			zap.L().Error("EditTask UpdateTaskFields error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		interceptor.ClearTaskListCache()
	}

	return t.ReadTask(ctx, msg)
}

// TaskStagesSort 对阶段进行排序（调整 pre 阶段到 next 阶段前面）
func (t *TaskService) TaskStagesSort(ctx context.Context, msg *task.TaskStagesSortReqMessage) (*task.TaskStagesSortResponse, error) {
	preStageCode := encrypts.DecryptNoErr(msg.PreStageCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	preStage, err := t.taskStagesRepo.FindById(c, int(preStageCode))
	if err != nil || preStage == nil {
		return nil, errs.GrpcError(model.DBError)
	}

	err = t.transaction.Action(func(conn database.DbConn) error {
		if msg.NextStageCode != "" {
			nextStageCode := encrypts.DecryptNoErr(msg.NextStageCode)
			nextStage, err := t.taskStagesRepo.FindById(c, int(nextStageCode))
			if err != nil || nextStage == nil {
				return errs.GrpcError(model.DBError)
			}
			prepre, err := t.taskStagesRepo.FindByStageCodeLtSort(c, nextStage.ProjectCode, nextStage.Sort)
			if err != nil {
				return errs.GrpcError(model.DBError)
			}
			if prepre != nil {
				preStage.Sort = (prepre.Sort + nextStage.Sort) / 2
			} else {
				preStage.Sort = 0
			}
		} else {
			maxSort, err := t.taskStagesRepo.FindMaxSortByProjectCode(c, preStage.ProjectCode)
			if err != nil {
				return errs.GrpcError(model.DBError)
			}
			preStage.Sort = maxSort + 65536
		}
		return t.taskStagesRepo.UpdateTaskStagesSort(c, conn, preStage)
	})
	if err != nil {
		return nil, err
	}
	return &task.TaskStagesSortResponse{}, nil
}

// EditTaskStages 修改阶段名称
func (t *TaskService) EditTaskStages(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskStagesResponse, error) {
	stageCode := encrypts.DecryptNoErr(msg.StageCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ts, err := t.taskStagesRepo.FindById(c, int(stageCode))
	if err != nil || ts == nil {
		return nil, errs.GrpcError(model.DBError)
	}
	ts.Name = msg.Name
	err = t.transaction.Action(func(conn database.DbConn) error {
		return t.taskStagesRepo.UpdateTaskStagesSort(c, conn, ts)
	})
	if err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskStagesResponse{}, nil
}

// DeleteTaskStages 删除阶段
func (t *TaskService) DeleteTaskStages(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskStagesResponse, error) {
	stageCode := encrypts.DecryptNoErr(msg.StageCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := t.transaction.Action(func(conn database.DbConn) error {
		return t.taskStagesRepo.DeleteTaskStages(c, conn, int(stageCode))
	})
	if err != nil {
		zap.L().Error("DeleteTaskStages error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskStagesResponse{}, nil
}

// fileToMessage 将 data.File 转换为 task.TaskFileMessage
func fileToMessage(f *data.File, creatorName string) *task.TaskFileMessage {
	return &task.TaskFileMessage{
		Id:               f.Id,
		Code:             encrypts.EncryptNoErr(f.Id),
		PathName:         f.PathName,
		FileName:         f.Title,
		Extension:        f.Extension,
		Size:             int64(f.Size),
		ProjectCode:      encrypts.EncryptNoErr(f.ProjectCode),
		TaskCode:         encrypts.EncryptNoErr(f.TaskCode),
		OrganizationCode: encrypts.EncryptNoErr(f.OrganizationCode),
		FileUrl:          f.FileUrl,
		FileType:         f.FileType,
		CreateBy:         encrypts.EncryptNoErr(f.CreateBy),
		CreateTime:       tms.FormatByMill(f.CreateTime),
		Downloads:        int32(f.Downloads),
		Extra:            f.Extra,
		Deleted:          int32(f.Deleted),
		DeletedTime:      tms.FormatByMill(f.DeletedTime),
		CreatorName:      creatorName,
	}
}

// FileIndex 查询项目文件列表
func (t *TaskService) FileIndex(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskFileListResponse, error) {
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	list, total, err := t.fileRepo.FindByProjectCode(c, projectCode, int(msg.Deleted), msg.Page, msg.PageSize)
	if err != nil {
		zap.L().Error("FileIndex FindByProjectCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	var mIds []int64
	for _, f := range list {
		if f.CreateBy > 0 {
			mIds = append(mIds, f.CreateBy)
		}
	}

	mMap := make(map[int64]string)
	if len(mIds) > 0 {
		memberMessageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mIds})
		if err != nil {
			zap.L().Error("FileIndex FindMemInfoByIds error", zap.Error(err))
		} else if memberMessageList != nil && memberMessageList.List != nil {
			for _, m := range memberMessageList.List {
				mMap[m.Id] = m.Name
			}
		}
	}

	var msgs []*task.TaskFileMessage
	for _, f := range list {
		msgs = append(msgs, fileToMessage(f, mMap[f.CreateBy]))
	}
	return &task.TaskFileListResponse{List: msgs, Total: total}, nil
}

// FileRead 查询单个文件详情
func (t *TaskService) FileRead(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskFileResponse, error) {
	fileCode := encrypts.DecryptNoErr(msg.FileCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	file, err := t.fileRepo.FindById(c, fileCode)
	if err != nil {
		zap.L().Error("FileRead FindById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if file == nil {
		return &task.TaskFileResponse{}, nil
	}

	creatorName := ""
	if file.CreateBy > 0 {
		memberMessage, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: file.CreateBy})
		if err != nil {
			zap.L().Error("FileRead FindMemInfoById error", zap.Error(err))
		} else if memberMessage != nil {
			creatorName = memberMessage.Name
		}
	}

	fileMsg := fileToMessage(file, creatorName)
	return &task.TaskFileResponse{File: fileMsg}, nil
}

// FileEdit 修改文件标题
func (t *TaskService) FileEdit(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskFileResponse, error) {
	fileCode := encrypts.DecryptNoErr(msg.FileCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := t.fileRepo.UpdateTitle(c, fileCode, msg.Title); err != nil {
		zap.L().Error("FileEdit UpdateTitle error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskFileResponse{}, nil
}

// FileRecycle 软删除文件（放入回收站）
func (t *TaskService) FileRecycle(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskFileResponse, error) {
	fileCode := encrypts.DecryptNoErr(msg.FileCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := t.fileRepo.UpdateDeleted(c, fileCode, 1); err != nil {
		zap.L().Error("FileRecycle UpdateDeleted error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskFileResponse{}, nil
}

// FileRecovery 从回收站恢复文件
func (t *TaskService) FileRecovery(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskFileResponse, error) {
	fileCode := encrypts.DecryptNoErr(msg.FileCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := t.fileRepo.UpdateDeleted(c, fileCode, 0); err != nil {
		zap.L().Error("FileRecovery UpdateDeleted error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskFileResponse{}, nil
}

// FileDelete 永久删除文件
func (t *TaskService) FileDelete(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskFileResponse, error) {
	fileCode := encrypts.DecryptNoErr(msg.FileCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	file, err := t.fileRepo.FindById(c, fileCode)
	if err != nil {
		zap.L().Error("FileDelete FindById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if file != nil && file.PathName != "" {
		minioClient, err := min.GetMinioClient()
		if err == nil && minioClient != nil {
			bucket := min.GetBucket()
			err = minioClient.Delete(ctx, bucket, file.PathName)
			if err != nil {
				zap.L().Error("FileDelete MinIO Delete error", zap.Error(err))
			}
		}
	}

	if err := t.fileRepo.DeleteById(c, fileCode); err != nil {
		zap.L().Error("FileDelete DeleteById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskFileResponse{}, nil
}

// TaskRecycle 软删除任务（放入回收站）
func (t *TaskService) TaskRecycle(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskRecycleResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := t.taskRepo.RecoverTask(c, taskCode, 1); err != nil {
		zap.L().Error("TaskRecycle RecoverTask error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	interceptor.ClearTaskListCache()
	return &task.TaskRecycleResponse{}, nil
}

// TaskDelete 永久删除任务
func (t *TaskService) TaskDelete(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskDeleteResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := t.taskRepo.DeleteTask(c, taskCode); err != nil {
		zap.L().Error("TaskDelete DeleteTask error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	interceptor.ClearTaskListCache()
	return &task.TaskDeleteResponse{}, nil
}

// TaskRecovery 从回收站恢复任务
func (t *TaskService) TaskRecovery(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskRecoveryResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := t.taskRepo.RecoverTask(c, taskCode, 0); err != nil {
		zap.L().Error("TaskRecovery RecoverTask error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	interceptor.ClearTaskListCache()
	return &task.TaskRecoveryResponse{}, nil
}
