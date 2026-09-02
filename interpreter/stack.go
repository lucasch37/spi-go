package interpreter

import (
	"fmt"
	"strings"
)

type CallStack struct {
	Records []*ActivationRecord
}

func NewCallStack() *CallStack {
	return &CallStack{
		Records: make([]*ActivationRecord, 0),
	}
}

func (cs *CallStack) Push(record *ActivationRecord) {
	cs.Records = append(cs.Records, record)
}

func (cs *CallStack) Pop() {
	cs.Records = cs.Records[:len(cs.Records)-1]
}

func (cs *CallStack) Peek() *ActivationRecord {
	return cs.Records[len(cs.Records)-1]
}

func (cs *CallStack) String() string {
	var lines []string

	for i := len(cs.Records) - 1; i >= 0; i-- {
		lines = append(lines, cs.Records[i].String())
	}

	return "CALL STACK\n" + strings.Join(lines, "\n") + "\n"
}

type ARType int

const (
	PROGRAM ARType = iota
)

var arTypeNames = [...]string{
	"PROGRAM",
}

func (art ARType) String() string {
	if int(art) >= 0 && int(art) < len(arTypeNames) {
		return arTypeNames[art]
	}

	return "UNKNOWN"
}

type ActivationRecord struct {
	Name         string
	Type         ARType
	NestingLevel int
	Members      map[string]Object
}

func NewActivationRecord(name string, recordType ARType, nestingLevel int) *ActivationRecord {
	return &ActivationRecord{
		Name:         name,
		Type:         recordType,
		NestingLevel: nestingLevel,
		Members:      make(map[string]Object),
	}
}

func (ar *ActivationRecord) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("%d: %s %s", ar.NestingLevel, ar.Type.String(), ar.Name))

	for k, v := range ar.Members {
		lines = append(lines, fmt.Sprintf("%s: %v", k, v))
	}

	return strings.Join(lines, "\n")
}

func (ar *ActivationRecord) Set(key string, value Object) {
	ar.Members[key] = value
}

func (ar *ActivationRecord) Get(key string) Object {
	return ar.Members[key]
}
