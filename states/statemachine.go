// A state machine - not ready yet.
package states


type StateMachine struct {
	name         string
	States       []State
	Transitions  []*Transition
	currentState State
	parent       *StateMachine
}


func (self *StateMachine) SetStateMachine(parent *StateMachine) {
	self.parent = parent
}


func (self *StateMachine) StateMachine() *StateMachine {
	return self.parent
}


func (self *StateMachine) Name() string {
	return self.name
}


func (self *StateMachine) Enter(transition *Transition) {
	if self.currentState != nil {
		panic("Current state machine already entered")
	}

	transition = &Transition{}
	for _, state := range self.States {
		transition.To = state
		if !transition.Veto() {
			transition.Do()
			return
		}
	}

	transition.From = self
	transition.To = nil
	transition.Do()
}


func (self *StateMachine) Leave(transition *Transition) {
}


func (self *StateMachine) CurrentState() State {
	return self.currentState
}


func (self *StateMachine) Event(event string) {
	if self.currentState == nil {
		panic("State machine can't handle event because it's not active")
	}

	for _, tr := range self.Transitions {
		if tr.From == self.currentState && tr.Event == event && !tr.Veto() {
			tr.Do()
			return
		}
	}
}
