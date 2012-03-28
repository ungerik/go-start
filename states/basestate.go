package states


type BaseState struct {
	name         string
	stateMachine *StateMachine
}


func (self *BaseState) SetStateMachine(parent *StateMachine) {
	self.stateMachine = parent
}


func (self *BaseState) StateMachine() *StateMachine {
	return self.stateMachine
}


func (self *BaseState) Name() string {
	return self.name
}


func (self *BaseState) Enter(*Transition) {
	panic("Enter() of state not implemented")
}


func (self *BaseState) Leave(*Transition) {
	panic("Leave() of state not implemented")
}
