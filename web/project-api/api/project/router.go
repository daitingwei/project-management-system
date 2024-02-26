package project

import (
	"log"

	"github.com/gin-gonic/gin"

	"test.com/project-api/api/midd"
	"test.com/project-api/api/rpc"
	"test.com/project-api/middleware"
	"test.com/project-api/router"
)

type RouterProject struct {
}

func init() {
	log.Println("init project router")
	ru := &RouterProject{}
	router.Register(ru)
}

func (*RouterProject) Route(r *gin.Engine) {
	// 初始化grpc的客户端连接
	rpc.InitRpcProjectClient()
	h := New()
	group := r.Group("/project")
	group.Use(midd.TokenVerify())
	group.Use(middleware.Auth())
	group.Use(ProjectAuth())
	group.POST("/index", h.index)
	group.POST("/index/uploadAvatar", NewIndex().uploadAvatar)
	group.POST("/index/uploadImg", NewIndex().uploadImg)
	group.POST("/index/editPersonal", NewIndex().editPersonal)
	group.POST("/project/selfList", h.myProjectList)
	group.POST("/project", h.myProjectList)
	group.POST("/project_template", h.projectTemplate)
	group.POST("/project_template/save", h.projectTemplateSave)
	group.POST("/project_template/edit", h.projectTemplateEdit)
	group.POST("/project_template/delete", h.projectTemplateDelete)
	group.POST("/project/save", h.projectSave)
	group.POST("/project/read", h.readProject)
	group.POST("/project/recycle", h.recycleProject)
	group.POST("/project/recovery", h.recoveryProject)
	group.POST("/project_collect/collect", h.collectProject)
	group.POST("/project/edit", h.editProject)
	group.POST("/project/getLogBySelfProject", h.getLogBySelfProject)
	group.POST("/node", h.nodeList)
	t := NewTask()
	group.POST("/task_stages", t.taskStages)
	// 新增任务阶段管理路由
	// 修改说明：新增任务阶段的创建、编辑、删除API路由
	// 背景：需要支持用户在项目下管理任务阶段（列）
	group.POST("/task_stages/save", t.saveTaskStages)
	group.POST("/task_stages/edit", t.editTaskStages)
	group.POST("/task_stages/delete", t.deleteTaskStages)
	group.POST("/task_stages/sort", t.taskStagesSort)
	group.POST("/project_member/index", t.memberProjectList)
	group.POST("/task_stages/tasks", t.taskList)
	group.POST("/task/save", t.saveTask)
	group.POST("/task", t.taskListByProject)
	group.POST("/task/edit", t.editTask)
	group.POST("/task/sort", t.taskSort)
	group.POST("/task/taskDone", t.taskDone)
	group.POST("/task/selfList", t.myTaskList)
	group.POST("/task/read", t.readTask)
	group.POST("/task_member", t.listTaskMember)
	group.POST("/task/taskLog", t.taskLog)
	group.POST("/task/_taskWorkTimeList", t.taskWorkTimeList)
	group.POST("/task/saveTaskWorkTime", t.saveTaskWorkTime)
	group.POST("/file/uploadFiles", t.uploadFiles)
	// 新增文件管理路由
	// 修改说明：新增项目文件的管理API路由（查询、编辑、回收站、恢复、删除）
	// 背景：需要支持项目文件的完整生命周期管理
	group.POST("/file", t.fileIndex)
	group.POST("/file/read", t.fileRead)
	group.POST("/file/edit", t.fileEdit)
	group.POST("/file/recycle", t.fileRecycle)
	group.POST("/file/recovery", t.fileRecovery)
	group.POST("/file/delete", t.fileDelete)
	group.POST("/file/download", t.fileDownload)
	group.POST("/task/taskSources", t.taskSources)
	group.POST("/task/createComment", t.createComment)
	group.POST("/task/recycle", t.taskRecycle)
	group.POST("/task/delete", t.taskDelete)
	group.POST("/task/recovery", t.taskRecovery)
	group.POST("/task/assignTask", t.assignTask)
	group.POST("/task/batchAssignTask", t.batchAssignTask)
	group.POST("/task/editTaskWorkTime", t.editTaskWorkTime)
	group.POST("/task/delTaskWorkTime", t.delTaskWorkTime)
	group.POST("/task/like", t.taskLike)
	group.POST("/task/star", t.taskStar)
	group.POST("/task/setPrivate", t.setPrivate)
	group.POST("/task/recycleBatch", t.recycleBatch)
	group.POST("/task_member/inviteMember", t.inviteMember)
	group.POST("/task_member/inviteMemberBatch", t.inviteMemberBatch)
	group.POST("/project/_projectStats", t.projectStats)
	group.POST("/project/_getProjectReport", t.getProjectReport)
	group.POST("/task/dateTotalForProject", t.dateTotalForProject)

	group.POST("/project/quit", h.quit)
	group.POST("/project/archive", h.archive)
	group.POST("/project/recoveryArchive", h.recoveryArchive)
	group.POST("/project/delete", h.deleteProject)
	group.POST("/project/analysis", h.analysis)
	group.POST("/project_member/searchInviteMember", h.searchInviteMember)
	group.POST("/project_member/inviteMember", h.inviteMember)
	group.POST("/project_member/_joinByInviteLink", h.joinByInviteLink)
	group.POST("/project_member/removeMember", h.removeMember)
	group.POST("/project_member/_listForInvite", h.listForInvite)
	group.POST("/source_link/delete", h.sourceLinkDelete)

	a := NewAccount()
	group.POST("/account", a.account)
	group.POST("/account/_allList", a.allList)
	group.POST("/account/forbid", a.forbid)
	group.POST("/account/resume", a.resume)
	group.POST("/account/add", a.add)
	group.POST("/account/edit", a.edit)
	group.POST("/account/auth", a.auth)
	group.POST("/account/del", a.del)
	group.POST("/account/read", a.read)
	group.POST("/account/_syncDetail", a.syncDetail)
	group.POST("/account/_joinByInviteLink", a.joinByInviteLink)
	d := NewDepartment()
	group.POST("/department", d.department)
	group.POST("/department/save", d.save)
	group.POST("/department/read", d.read)
	group.POST("/department/edit", d.edit)
	group.POST("/department/delete", d.delete)
	dm := NewDepartmentMember()
	group.POST("/department_member/index", dm.list)
	group.POST("/department_member/searchInviteMember", dm.searchInviteMember)
	group.POST("/department_member/inviteMember", dm.inviteMember)
	group.POST("/department_member/removeMember", dm.removeMember)
	group.POST("/department_member/detail", dm.detail)
	auth := NewAuth()
	group.POST("/auth", auth.authList)
	group.POST("/auth/apply", auth.apply)
	group.POST("/auth/add", auth.add)
	group.POST("/auth/edit", auth.edit)
	group.POST("/auth/forbid", auth.forbid)
	group.POST("/auth/resume", auth.resume)
	group.POST("/auth/setDefault", auth.setDefault)
	group.POST("/auth/del", auth.del)
	menu := NewMenu()
	group.POST("/menu/menu", menu.menuList)
	tt := NewTaskTag()
	group.POST("/task_tag", tt.list)
	group.POST("/task_tag/save", tt.save)
	group.POST("/task_tag/edit", tt.edit)
	group.POST("/task_tag/delete", tt.del)
}
