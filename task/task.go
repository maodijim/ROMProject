package task

type Template interface {
	SetLogger()
	Start()
	Stop()
}
