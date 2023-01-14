package runtime

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewMulti(configurations ...Configuration[*MultiRuntime]) Runtime {
	consumer := &MultiRuntime{}
	for _, configuration := range configurations {
		consumer = configuration(consumer)
	}
	return consumer
}

// configuration
func WithController(controller Controller) Configuration[*MultiRuntime] {
	return func(c *MultiRuntime) *MultiRuntime {
		c.controller = controller
		return c
	}
}

func WithRuntime(runtime Runtime) Configuration[*MultiRuntime] {
	return func(c *MultiRuntime) *MultiRuntime {
		if c.runtimes == nil {
			c.runtimes = make([]Runtime, 0)
		}
		c.runtimes = append(c.runtimes, runtime)
		return c
	}
}

// implementation
type MultiRuntime struct {
	controller Controller
	runtimes   []Runtime
}

func (r *MultiRuntime) Start() error {
	for i := 0; i < len(r.runtimes); i++ {
		if err := r.runtimes[i].Start(); err != nil {
			panic(err)
		}
	}

	go r.Interrupt()

	r.controller.Started()
	r.controller.Wait()

	return nil
}

func (r *MultiRuntime) Stop() {
	for i := len(r.runtimes); i > 0; i-- {
		r.runtimes[i-1].Stop()
	}

	r.controller.Stopped()
}

func (runtime *MultiRuntime) Interrupt() {
	interruptSignal := make(chan os.Signal, 10)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGTERM)
	<-interruptSignal
	runtime.Stop()
	// os.Exit(0) -- don't need os.Exit if everything is cleaned up properly
}

func (runtime *MultiRuntime) Panic() {
	logger.Infof("panicking")
	if x := recover(); x != nil {
		logger.Error("runtime panic", zap.Error(x.(error)))
		runtime.Stop()
		os.Exit(1)
	}
}
