package states


type Transition struct {
	Event string
	From  State
	To    State
}


func (self *Transition) Veto() bool {
	if vetoer, ok := self.From.(TransitionVetoer); ok {
		if vetoer.TransitionVeto(self) {
			return true
		}
	}
	if vetoer, ok := self.To.(TransitionVetoer); ok {
		if vetoer.TransitionVeto(self) {
			return true
		}
	}
	return false
}


func (self *Transition) Do() {
	if self.From != nil {
		stateMachine := self.From.StateMachine()

		self.From.Leave(self)
		if stateMachine != nil {
			stateMachine.currentState = nil
		}

		if self.To == nil {
			if stateMachine != nil {
				stateMachine.Leave(self)
			}
			return
		}
	}

	self.To.Enter(self)
	if stateMachine := self.To.StateMachine(); stateMachine != nil {
		stateMachine.currentState = self.To
	}
}
