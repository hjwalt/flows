package runtime

type Runtime interface {
	Start() error
	Stop()
}
