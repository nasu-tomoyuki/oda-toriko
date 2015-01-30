/*
 * julius の制御
 */
package julius

import (
	"net"
	"time"
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
)

type GramType uint32
const (
	GramDummy			GramType = iota
	GramStartOrder
	GramFinishOrder
	GramYes
	GramCancel
	GramMenu
)


type Whypo struct {
	XMLName xml.Name "WHYPO"
	Word    string   `xml:"WORD,attr"`
	Classid int      `xml:"CLASSID,attr"`
	Phone   string   `xml:"PHONE,attr"`
	Cm      float32  `xml:"CM,attr"`
}

type Shypo struct {
	XMLName xml.Name "SHYPO"
	Rank    int      `xml:"RANK,attr"`
	Score   float32  `xml:"SCORE,attr"`
	Gram    int      `xml:"GRAM,attr"`
	Whypo   []Whypo  `xml:"WHYPO"`
}

type Recogout struct {
	XMLName xml.Name "RECOGOUT"
	Shypo   []Shypo  `xml:"SHYPO"`
}

func (self *Recogout) Clear() {
	self.Shypo = make([]Shypo, 0)
}


type Julius struct {
	TelnetAddr	string
	TelnetPort	string
	Connection	net.Conn
	Reader		*bufio.Reader
	Recogout	Recogout
}

func NewJulius() *Julius {
	s := Julius {
		TelnetAddr: "",
		TelnetPort: "",
		Connection: nil,
		Reader: nil,
	}
	return &s
}

// telnet 接続
func (self *Julius)Connect(addr string, port string) bool {
	self.TelnetAddr		= addr
	self.TelnetPort		= port
	var err error
	self.Connection, err = net.DialTimeout("tcp", addr + ":" + port, 8 * time.Second)
	if err != nil {
		return false
	}

	fmt.Printf("Connected to %s:%s\n", addr, port)
	self.Reader = bufio.NewReaderSize(self.Connection, 4096)
	return true
}

// telnet を切断
func (self *Julius)Disconnect() bool {
	if self.Connection == nil {
		return false
	}
	self.Connection.Close()
	self.Connection = nil
	self.Reader = nil
	return true
}

func (self *Julius)Input() (*Recogout, error) {
	for {
		str, _, err := self.Reader.ReadLine()
		if len(str) > 0 {
			r, err := self.waitRecogout(string(str))
			if err != nil {
				return nil, err
			}
			if r == false {
				continue
			}
			if len(self.Recogout.Shypo) > 0 {
				return &self.Recogout, nil
			}
		}
		if err == io.EOF {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
	}
}

func (self *Julius)send(cmd string) {
	self.Connection.Write([]byte(cmd + "\n"))
	// 適当にウェイトを入れる
	// 本来は応答を待ちたいけれど、TERMINATE 中は無反応なのでウェイトで対応
	time.Sleep(10 * time.Millisecond)
}

func (self *Julius)read() []byte {
	str, _, err := self.Reader.ReadLine()
	if nil != err {
		return nil
	}
	return str
}

func (self *Julius)Terminate() {
	self.send("TERMINATE")
}

func (self *Julius)Resume() {
	self.send("RESUME")
}

func (self *Julius)DeactivateGrams(grams ...GramType) {
	s := ""
	for _, v := range grams {
		s += fmt.Sprint(v) + " "
	}
	self.send("DEACTIVATEGRAM\n" + s)
//	self.read()
}

func (self *Julius)ActivateGrams(grams ...GramType) {
	s := ""
	for _, v := range grams {
		s += " " + fmt.Sprint(v)
	}
	self.send("ACTIVATEGRAM\n" + s)
//	self.read()
}



func (self *Julius)waitRecogout(startString string) (bool, error) {
	if strings.Contains(startString, "<RECOGOUT>") == false {
//		fmt.Printf("%+v\n", startString)
		return false, nil
	}
	var recogout string
	recogout = startString
	for {
		str, _, err := self.Reader.ReadLine()
		if len(str) > 0 {
			s := string(str)
			recogout += s
			if strings.Contains(s, "</RECOGOUT>") == true {
				err = self.parseRecogout(recogout)
				if nil != err {
					return false, err
				}
				return true, nil
			}
		}
		if err == io.EOF {
//			fmt.Println("[EOF]")
			return false, err
		}
		if err != nil {
			return false, err
		}
	}
}

func (self *Julius)parseRecogout(str string) error {
	self.Recogout.Clear()
	err := xml.Unmarshal([]byte(str), &self.Recogout)
	if err != nil {
		log.Fatal(err)
		return err
	}
//	fmt.Printf("parseRecogout: %+v\n", self.Recogout)
	return err
}

