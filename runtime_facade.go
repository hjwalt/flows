package flows

import (
	"time"

	"github.com/hjwalt/runway/runtime"
)

type RuntimeFacade struct {
	Runtimes []runtime.Runtime
}

func (r *RuntimeFacade) Start() error {
	err := runtime.Start(r.Runtimes, time.Second)
	if err != nil {
		return err
	}
	runtime.Wait()
	return nil
}

func (r *RuntimeFacade) Stop() {
	runtime.Stop()
}
