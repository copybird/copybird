package core

type ModuleGroup string

const (
	ModuleGroupBackup  ModuleGroup = "backup"
	ModuleGroupRestore ModuleGroup = "restore"
	ModuleGroupGlobal  ModuleGroup = "global"
)
