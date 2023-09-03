package flows

import (
	"github.com/hjwalt/runway/inverse"
)

type Prebuilt interface {
	Register(ci inverse.Container)
}
