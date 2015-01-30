package staff

import hsm "github.com/hhkbp2/go-hsm"
import "github.com/nasu-tomoyuki/oda-toriko/julius"
import "github.com/nasu-tomoyuki/oda-toriko/voice"
import "fmt"

const (
	HSMTypeStaff hsm.HSMType = hsm.HSMTypeUser
)

type StaffHSM struct {
	*hsm.StdHSM
	
	Julius		*julius.Julius
	Voice		*voice.Voice

	Orders		map[string]int
	LastOrder	string

	IsConfirmed	bool
	IsEnabledInput	bool
	foo bool
}

func NewStaffHSM(top, initial hsm.State) *StaffHSM {
	s := StaffHSM{
		StdHSM: hsm.NewStdHSM(HSMTypeStaff, top, initial),
	}
	s.IsConfirmed		= false
	s.IsEnabledInput	= true
	s.Orders			= make(map[string]int)
	return &s
}

func (self *StaffHSM) Init() {
	self.StdHSM.Init2(self, hsm.StdEvents[hsm.EventInit])
}

func (self *StaffHSM) Dispatch(event hsm.Event) {
	self.StdHSM.Dispatch2(self, event)
}

func (self *StaffHSM) QTran(targetStateID string) {
	target := self.StdHSM.LookupState(targetStateID)
	self.StdHSM.QTranHSM(self, target)
}

func (self *StaffHSM) CurrentStateID() string {
	return self.StdHSM.State.ID()
}

func (self *StaffHSM) Update() bool {
	if self.IsEnabledInput {
		var recogout, _	= self.Julius.Input()
		if recogout != nil {
			fmt.Printf("Update: %+v\n", *recogout)
			event := int(EventBase) + recogout.Shypo[0].Gram
			self.Dispatch(NewIntEvent(event))
		}
	}
	self.Dispatch(NewEvent(EventUpdate))
	return true
}

func (self *StaffHSM) GetFoo() bool {
	return self.foo
}

func (self *StaffHSM) SetFoo(newFoo bool) {
	self.foo = newFoo
}
