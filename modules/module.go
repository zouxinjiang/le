package modules

type Module interface {
	Init() error
	Install() error
	Start(params map[string]string) error
	Stop(params map[string]string) error
}
