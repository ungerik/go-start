package states

import (
	"bytes"
)


const (
	Transient = iota
	SessionPersistent
	UserPersistent
)


type Persistence uint


type State interface {
	SetStateMachine(*StateMachine)
	StateMachine() *StateMachine
	Name() string
	Enter(*Transition)
	Leave(*Transition)
}


type TransitionVetoer interface {
	TransitionVeto(*Transition) bool
}


func writeStatePathRecursive(state State, buf *bytes.Buffer) *bytes.Buffer {
	parent := state.StateMachine()
	if parent != nil {
		writeStatePathRecursive(parent, buf)
	}
	buf.WriteByte('/')
	buf.WriteString(state.Name())
	return buf
}


func StatePath(state State) string {
	return writeStatePathRecursive(state, &bytes.Buffer{}).String()
}
