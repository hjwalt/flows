package runtime_retry

import "errors"

func AlwaysTry(err error) bool {
	if err != nil && err.Error() == "stopped" {
		return false
	}
	return true
}

func Stopped() error {
	return errors.New("stopped")
}
