package model

const (
	Normal         = 1
	Personal int32 = 1
)

const AESKey = "sdfgyrhgbxcdgryfhgywertd"

const (
	NoDeleted = iota
	Deleted
	// TODO: 删除相关功能未完全实现，Deleted 常量未被使用
)

const (
	NoArchive = iota
	Archive
	// TODO: 归档功能未完全实现，Archive 常量未被使用
)

const (
	Open = iota // TODO: 访问控制类型，代码中直接写数字1，未使用此常量
	Private
	Custom // TODO: 未被使用
)

const (
	Default = "default"
	Simple  = "simple"
)

const (
	NoCollected = iota // TODO: NoCollected 未被使用
	Collected
)

const (
	NoOwner = iota // TODO: NoOwner 未被使用
	Owner
)

const (
	NoExecutor = iota // TODO: NoExecutor 未被使用
	Executor
)
const (
	NoCanRead = iota // TODO: NoCanRead 未被使用
	CanRead // TODO: CanRead 未被使用
)

const (
	UnDone = iota // TODO: 未被使用，代码中直接写数字判断
	Done
)

const (
	NoComment = iota // TODO: NoComment 未被使用
	Comment
)