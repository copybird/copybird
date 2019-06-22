package input

type Input interface {
	Init()
	Backup()
	Close()
}
