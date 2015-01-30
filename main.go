/*
 * ファミレスごっこが出来ます
 * 小田鳥子
 */
package main

import (
	"flag"
	"fmt"
	"github.com/nasu-tomoyuki/oda-toriko/julius"
	"github.com/nasu-tomoyuki/oda-toriko/voice"
	"github.com/nasu-tomoyuki/oda-toriko/staff"
)


func main() {
	addr		:= flag.String("addr", "localhost", "julius address")
	port		:= flag.String("port", "10500", "julius port")
	assetsPath	:= flag.String("assets", "./assets/", "assets path")
	tmpPath		:= flag.String("tmp", "./tmp/", "tmp path")
	flag.Parse()

	v		:= voice.NewVoice(*assetsPath + "sfx/", *tmpPath)
	j		:= julius.NewJulius()

	if !j.Connect(*addr, *port) {
		fmt.Printf("Failed to connect to %s:%s\n", *addr, *port)
		return
	}
	defer j.Disconnect()

	staff := staff.NewStaff(j, v)
	fmt.Println(staff.CurrentStateID())

	for {
		staff.Update()
	}
}
