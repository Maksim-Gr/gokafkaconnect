package tui

type viewMode int

const (
	viewHome viewMode = iota
	viewBackup
	viewConfigure
	viewCreate
	viewDelete
	viewList
	viewShowConfig
	viewHealthcheck
	viewResult
)
