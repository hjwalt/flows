package runtime_bun

import "github.com/uptrace/bun"

type BunConnection interface {
	Start() error
	Stop()
	Db() bun.IDB
}
