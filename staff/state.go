package staff

import hsm "github.com/hhkbp2/go-hsm"

import (
	"log"
	"fmt"
	"strconv"
	"github.com/nasu-tomoyuki/oda-toriko/julius"
)

const (
	StateGlobalID			string	= "global"				// 全体
	StateIdlingID					= "idling"				// 注文受け付け待ち
	StateInputID					= "inputing"			// 注文受け付け中
	StateInputHearingID				= "inputing:hearing"	// 注文を聞いている
	StateInputConfirmationID		= "inputing:confirmation"	// 注文を復唱
	StateConfirmationID				= "confirmation"		// すべての注文の確認
	StateByeID						= "bye"					// 完了
)

const (
	MenuItem				int	= iota
	MenuNum
	MenuCancel
)


var NumToCountStr = map[int]string {
	0:		"なし",
	1:		"おひとつ",
	2:		"おふたつ",
	3:		"みっつ",
	4:		"よっつ",
	5:		"いつつ",
	6:		"むっつ",
	7:		"ななつ",
	8:		"やっつ",
	9:		"ここのつ",
	10:		"じゅっこ",
}
func getCountString(n int) string {
	return NumToCountStr[n]
}

func Logln(v ...interface{}) {
	log.Println(v...)
}

type VerboseStateHead struct {
	*hsm.StateHead
	ID string
}

func NewVerboseStateHead(super hsm.State) *VerboseStateHead {
	return &VerboseStateHead{
		StateHead: hsm.NewStateHead(super),
	}
}

func (self *VerboseStateHead) Init(hsm hsm.HSM, event hsm.Event) (state hsm.State) {
	Logln(self.ID, "- Init")
	return nil
}

func (self *VerboseStateHead) Entry(hsm hsm.HSM, event hsm.Event) (state hsm.State) {
	Logln(self.ID, "- Entry")
	return nil
}

func (self *VerboseStateHead) Exit(hsm hsm.HSM, event hsm.Event) (state hsm.State) {
	Logln(self.ID, "- Exit")
	return nil
}


/*
 * 状態
 */
type GlobalState struct {
	*VerboseStateHead
}
func NewGlobalState(super hsm.State) *GlobalState {
	object := &GlobalState {
		VerboseStateHead:	NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}
func (_ *GlobalState) ID() string {
	return StateGlobalID
}
func (self *GlobalState) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	sm.QInit(StateIdlingID)
	return nil
}
func (self *GlobalState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
//	annotatedHSM, ok := sm.(*AnnotatedHSM)
//	hsm.AssertTrue(ok)
	/*
	switch event.Type() {
	case EventE:
		sm.QTran(StateS211ID)
		return nil
	}
	*/
	return self.Super()
}

type IdlingState struct {
	*VerboseStateHead
}
func NewIdlingState(super hsm.State) *IdlingState {
	object := &IdlingState {
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}
func (_ *IdlingState) ID() string {
	return StateIdlingID
}
func (self *IdlingState) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	return self.Super()
}
func (self *IdlingState) Entry(sm hsm.HSM, event hsm.Event) hsm.State {
	staffHSM, ok := sm.(*StaffHSM)
	hsm.AssertTrue(ok)
	staffHSM.Orders		= make(map[string]int)
	staffHSM.LastOrder	= ""
	staffHSM.IsConfirmed	= false
	staffHSM.IsEnabledInput	= true

	j := staffHSM.Julius
	j.Terminate()
	j.ActivateGrams(julius.GramStartOrder)
	j.DeactivateGrams(julius.GramFinishOrder, julius.GramYes, julius.GramCancel, julius.GramMenu)
	j.Resume()
	return nil
}
func (self *IdlingState) Exit(sm hsm.HSM, event hsm.Event) hsm.State {
	return self.Super()
}
func (self *IdlingState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	switch event.Type() {
	case EventStartOrder:
		staffHSM, ok := sm.(*StaffHSM)
		hsm.AssertTrue(ok)
		j := staffHSM.Julius
		v := staffHSM.Voice
		j.Terminate()
		j.ActivateGrams(julius.GramFinishOrder, julius.GramCancel, julius.GramMenu)
		j.DeactivateGrams(julius.GramStartOrder)
		v.Play("call.wav")
		v.Say("いらっしゃいませ。ご注文をおうかがいします。")
		v.Play("come_in.wav")
		j.Resume()
		sm.QTran(StateInputID)
		return nil
	}
	return self.Super()
}

type InputState struct {
	*VerboseStateHead
}
func NewInputState(super hsm.State) *InputState {
	object := &InputState {
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}
func (_ *InputState) ID() string {
	return StateInputID
}
func (self *InputState) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	sm.QInit(StateInputHearingID)
	return nil
}
func (self *InputState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	return self.Super()
}

type InputHearingState struct {
	*VerboseStateHead
}
func NewInputHearingState(super hsm.State) *InputHearingState {
	object := &InputHearingState {
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}
func (_ *InputHearingState) ID() string {
	return StateInputHearingID
}
func (self *InputHearingState) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	return self.Super()
}
func (self *InputHearingState) Entry(sm hsm.HSM, event hsm.Event) hsm.State {
	staffHSM, ok := sm.(*StaffHSM)
	hsm.AssertTrue(ok)

	if staffHSM.IsConfirmed {
		j	:= staffHSM.Julius
		j.Terminate()
		j.ActivateGrams(julius.GramYes)
		j.Resume()
	}
	return nil
}
func (self *InputHearingState) Exit(sm hsm.HSM, event hsm.Event) hsm.State {
	staffHSM, ok := sm.(*StaffHSM)
	hsm.AssertTrue(ok)
	if staffHSM.IsConfirmed {
		j	:= staffHSM.Julius
		j.Terminate()
		j.DeactivateGrams(julius.GramYes)
		j.Resume()
		staffHSM.IsConfirmed	= false
	}
	return nil
}
func (self *InputHearingState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	staffHSM, ok := sm.(*StaffHSM)
	hsm.AssertTrue(ok)

	switch event.Type() {
	case EventMenu:
		sm.QTran(StateInputConfirmationID)
		return nil
	case EventYes:
		if staffHSM.IsConfirmed == true {
			sm.QTran(StateByeID)
			return nil
		}
		return self.Super()
	case EventFinishOrder:
		if staffHSM.IsConfirmed == true {
			sm.QTran(StateByeID)
			return nil
		}
		sm.QTran(StateConfirmationID)
		return nil
	}
	return self.Super()
}

type InputConfirmationState struct {
	*VerboseStateHead
}
func NewInputConfirmationState(super hsm.State) *InputConfirmationState {
	object := &InputConfirmationState {
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}
func (_ *InputConfirmationState) ID() string {
	return StateInputConfirmationID
}
func (self *InputConfirmationState) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	return self.Super()
}
func (self *InputConfirmationState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	switch event.Type() {
	case EventUpdate:
		staffHSM, ok := sm.(*StaffHSM)
		hsm.AssertTrue(ok)

		nextState	:= StateInputHearingID
		j := staffHSM.Julius
		v := staffHSM.Voice
		j.Terminate()
		v.Play("accepted.wav")
		isCanceled		:= false
		item := ""
		num := 1
		var count	string
		var text	string
		for _, v := range j.Recogout.Shypo[0].Whypo {
			switch v.Classid {
			case MenuItem:
				item	= v.Word
			case MenuNum:
				num, _		= strconv.Atoi(v.Word)
			case MenuCancel:
				isCanceled	= true
			}
		}
		// アイテムの指定がなければ最後の注文
		if item == "" {
			item	= staffHSM.LastOrder
			staffHSM.LastOrder	= ""
		}
		// キャンセル処理
		if isCanceled {
			// 一つも注文していなければ終了
			if len(staffHSM.Orders) == 0 {
				nextState	= StateByeID
			} else {
				if item != "" {
					var ok bool
					_, ok = staffHSM.Orders[item]
					if ok {
						delete(staffHSM.Orders, item)
						text		= fmt.Sprintf("%sをキャンセルしました\n", item)
					} else {
						text		= fmt.Sprintf("%sのご注文はありません\n", item)
					}
					v.Say(text)
				}
			}
		} else {
			if item != "" {
				count		= getCountString(num)
				text		= fmt.Sprintf("%sを%s\n", item, count)
				staffHSM.Orders[item]	= num
				staffHSM.LastOrder		= item
				v.Say(text)
			}
		}

		j.Resume()
		v.Play("come_in.wav")
		
		sm.QTran(nextState)
		return nil
	}
	return self.Super()
}

type ConfirmationState struct {
	*VerboseStateHead
}
func NewConfirmationState(super hsm.State) *ConfirmationState {
	object := &ConfirmationState {
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}
func (_ *ConfirmationState) ID() string {
	return StateConfirmationID
}
func (self *ConfirmationState) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	return self.Super()
}
func (self *ConfirmationState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	switch event.Type() {
	case EventUpdate:
		staffHSM, ok := sm.(*StaffHSM)
		hsm.AssertTrue(ok)

		j := staffHSM.Julius
		v := staffHSM.Voice
		if len(staffHSM.Orders) == 0 {
			sm.QTran(StateByeID)
		} else {
			staffHSM.IsConfirmed	= true
			j.Terminate()
			v.Play("accepted.wav")
			orders		:= ""
			for k, v := range staffHSM.Orders {
				orders += fmt.Sprintf("%sを%s、", k, getCountString(v))
			}
			v.Say("ご注文を確認します。" + orders + "　以上でよろしいですか？")
			j.Resume()
			v.Play("come_in.wav")
			sm.QTran(StateInputID)
		}
		return nil
	}
	return self.Super()
}

type ByeState struct {
	*VerboseStateHead
}
func NewByeState(super hsm.State) *ByeState {
	object := &ByeState {
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}
func (_ *ByeState) ID() string {
	return StateByeID
}
func (self *ByeState) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	return self.Super()
}
func (self *ByeState) Entry(sm hsm.HSM, event hsm.Event) hsm.State {
	staffHSM, ok := sm.(*StaffHSM)
	hsm.AssertTrue(ok)

	staffHSM.IsEnabledInput	= false
	return nil
}
func (self *ByeState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))

	switch event.Type() {
	case EventUpdate:
		staffHSM, ok := sm.(*StaffHSM)
		hsm.AssertTrue(ok)

		staffHSM.IsConfirmed	= true
		staffHSM.IsEnabledInput	= true

		j := staffHSM.Julius
		v := staffHSM.Voice
		j.Terminate()
		v.Play("accepted.wav")
		if len(staffHSM.Orders) == 0 {
			v.Say("また、のちほどおうかがいします。")
		} else {
			v.Say("ご注文、ありがとうございました。")
		}
		j.Resume()
		sm.QTran(StateIdlingID)
		return nil
	}
	return self.Super()
}





