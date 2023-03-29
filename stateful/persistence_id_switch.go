package stateful

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
)

// constructor
func NewSinglePersistenceIdSwitch(configurations ...runtime.Configuration[*PersistenceIdSwitch]) PersistenceIdFunction {
	persistenceIdFunction := &PersistenceIdSwitch{
		functions: make(map[string]PersistenceIdFunction),
	}
	for _, configuration := range configurations {
		persistenceIdFunction = configuration(persistenceIdFunction)
	}
	return persistenceIdFunction.Apply
}

// configuration
func WithPersistenceIdSwitchPersistenceIdFunction(topic string, f PersistenceIdFunction) runtime.Configuration[*PersistenceIdSwitch] {
	return func(pis *PersistenceIdSwitch) *PersistenceIdSwitch {
		pis.functions[topic] = f
		return pis
	}
}

// implementation
type PersistenceIdSwitch struct {
	functions map[string]PersistenceIdFunction
}

func (r *PersistenceIdSwitch) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes]) (string, error) {
	fn, fnExists := r.functions[m.Topic]
	if !fnExists {
		return "", TopicMissingError(m.Topic)
	}
	return fn(c, m)
}
