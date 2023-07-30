package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/runtime"
)

// constructor
func NewRouteFlow(configurations ...runtime.Configuration[*RouteFlow]) *RouteFlow {
	rp := &RouteFlow{
		entries: FlowEntries{
			Entries: make([]FlowEntry, 0),
		},
	}
	for _, configuration := range configurations {
		rp = configuration(rp)
	}
	return rp
}

// configurations
func WithFlow(input []FlowEntryTopic, output []FlowEntryTopic, stateTable string) runtime.Configuration[*RouteFlow] {
	return func(psf *RouteFlow) *RouteFlow {
		psf.entries.Entries = append(psf.entries.Entries, FlowEntry{Input: input, Output: output})
		return psf
	}
}

func WithFlowStatelessOneToOne[
	IK any,
	IV any,
	OK any,
	OV any,
](
	input topic.Topic[IK, IV],
	output topic.Topic[OK, OV],
) runtime.Configuration[*RouteFlow] {
	return func(psf *RouteFlow) *RouteFlow {
		psf.entries.Entries = append(psf.entries.Entries, FlowEntry{
			Input:  []FlowEntryTopic{FlowEntryFromTopic[IK, IV](input)},
			Output: []FlowEntryTopic{FlowEntryFromTopic[OK, OV](output)},
		})
		return psf
	}
}

func WithFlowStatelessOneToTwo[
	IK any,
	IV any,
	OK1 any,
	OV1 any,
	OK2 any,
	OV2 any,
](
	input topic.Topic[IK, IV],
	output1 topic.Topic[OK1, OV1],
	output2 topic.Topic[OK2, OV2],
) runtime.Configuration[*RouteFlow] {
	return func(psf *RouteFlow) *RouteFlow {
		psf.entries.Entries = append(psf.entries.Entries, FlowEntry{
			Input: []FlowEntryTopic{FlowEntryFromTopic[IK, IV](input)},
			Output: []FlowEntryTopic{
				FlowEntryFromTopic[OK1, OV1](output1),
				FlowEntryFromTopic[OK2, OV2](output2),
			},
		})
		return psf
	}
}

func WithFlowStatefulOneToOne[
	IK any,
	IV any,
	OK any,
	OV any,
](
	input topic.Topic[IK, IV],
	output topic.Topic[OK, OV],
	stateTable string,
) runtime.Configuration[*RouteFlow] {
	return func(psf *RouteFlow) *RouteFlow {
		psf.entries.Entries = append(psf.entries.Entries, FlowEntry{
			Input:      []FlowEntryTopic{FlowEntryFromTopic[IK, IV](input)},
			Output:     []FlowEntryTopic{FlowEntryFromTopic[OK, OV](output)},
			StateTable: stateTable,
		})
		return psf
	}
}

func WithFlowStatefulOneToTwo[
	IK any,
	IV any,
	OK1 any,
	OV1 any,
	OK2 any,
	OV2 any,
](
	input topic.Topic[IK, IV],
	output1 topic.Topic[OK1, OV1],
	output2 topic.Topic[OK2, OV2],
	stateTable string,
) runtime.Configuration[*RouteFlow] {
	return func(psf *RouteFlow) *RouteFlow {
		psf.entries.Entries = append(psf.entries.Entries, FlowEntry{
			Input: []FlowEntryTopic{FlowEntryFromTopic[IK, IV](input)},
			Output: []FlowEntryTopic{
				FlowEntryFromTopic[OK1, OV1](output1),
				FlowEntryFromTopic[OK2, OV2](output2),
			},
			StateTable: stateTable,
		})
		return psf
	}
}

func WithFlowMaterialiseOneToOne[
	IK any,
	IV any,
](
	input topic.Topic[IK, IV],
) runtime.Configuration[*RouteFlow] {
	return func(psf *RouteFlow) *RouteFlow {
		psf.entries.Entries = append(psf.entries.Entries, FlowEntry{
			Input:  []FlowEntryTopic{FlowEntryFromTopic[IK, IV](input)},
			Output: []FlowEntryTopic{},
		})
		return psf
	}
}

func FlowEntryFromTopic[IK any, IV any](topic topic.Topic[IK, IV]) FlowEntryTopic {
	return FlowEntryTopic{
		Topic:     topic.Topic(),
		KeyType:   fmt.Sprintf("%T", topic.KeyFormat().Default()),
		ValueType: fmt.Sprintf("%T", topic.ValueFormat().Default()),
	}
}

// implementation
type FlowEntryTopic struct {
	Topic     string `json:"topic"`
	KeyType   string `json:"key_type"`
	ValueType string `json:"value_type"`
}

type FlowEntry struct {
	Input      []FlowEntryTopic `json:"inputs"`
	Output     []FlowEntryTopic `json:"outputs"`
	StateTable string           `json:"state_table"`
}

type FlowEntries struct {
	Entries []FlowEntry `json:"entries"`
}

type RouteFlow struct {
	entries FlowEntries
}

func (rp *RouteFlow) Handle(w http.ResponseWriter, req *http.Request) {
	err := WriteJson[FlowEntries](w, 200, rp.entries, format.Json[FlowEntries]())
	if err != nil {
		WriteError(w, 500, errors.New("server error"))
	}
}
