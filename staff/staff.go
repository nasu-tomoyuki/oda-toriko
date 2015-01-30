package staff

import hsm "github.com/hhkbp2/go-hsm"
import "github.com/nasu-tomoyuki/oda-toriko/julius"
import "github.com/nasu-tomoyuki/oda-toriko/voice"

func NewStaff(julius *julius.Julius, voice *voice.Voice) *StaffHSM {
	top					:= hsm.NewTop()
	initial				:= hsm.NewInitial(top, StateGlobalID)
	sGlobal				:= NewGlobalState(top)
	NewIdlingState(sGlobal)
	sInput				:= NewInputState(sGlobal)
	NewInputHearingState(sInput)
	NewInputConfirmationState(sInput)
	NewConfirmationState(sGlobal)
	NewByeState(sGlobal)

	sm := NewStaffHSM(top, initial)
	sm.Julius	= julius
	sm.Voice	= voice
	sm.Init()
	return sm
}
