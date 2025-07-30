package def

import "golang.org/x/exp/constraints"

type Opt = int64

const (
	OptAdd    Opt = 0 //增加
	OptModify Opt = 1 //修改
	OptDel    Opt = 2 //删除
)

const (
	StatusRunning    = 1 //正常运行
	StatusWaitRun    = 2 //等待运行中
	StatusWaitStop   = 3 //等待暂停中
	StatusStopped    = 4 //已停用
	StatusWaitDelete = 5 //等待删除(该状态会先删除定时任务再删除表里的任务)
)

const Unknown = 0

const (
	True  = 1 //是
	False = 2 //否
)

type Bool = int64

const (
	Enable  Bool = 1 //启用
	Disable Bool = 2 //禁用
)

const (
	Male   = 1 //男性
	Female = 2 //女鞋
)

const (
	FilterTrue    = 1 //是
	FilterFalse   = 2 //否
	FilterNoAdmin = 3 // admin不过滤
)

func ToBool[boolType constraints.Integer](in boolType, defaultVal ...bool) bool {
	if in == True {
		return true
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return false
}
func ToIntBool[boolType constraints.Integer](in bool) boolType {
	if in == true {
		return True
	}
	return False
}
