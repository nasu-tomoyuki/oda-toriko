/*
 * 
 */
package voice

import (
	"os"
	"os/exec"
	"io/ioutil"
	"fmt"
)

type Voice struct {
	TmpPath			string
	TmpWav			string
	TmpTxt			string
	SfxPath			string
	OpenJtalk		string
	OpenJtalkDic	string
	OpenJtalkVoice	string
	Player			string
}

func	NewVoice(sfxPath string, tmpPath string) *Voice {
	s	:= Voice {
		TmpPath:		tmpPath,
		TmpWav:			"tmp.wav",
		TmpTxt:			"tmp.txt",
		SfxPath:		sfxPath,
		OpenJtalk:		"open_jtalk",
		OpenJtalkDic:	"/var/lib/mecab/dic/open-jtalk/naist-jdic/",
		OpenJtalkVoice:	"/usr/share/hts-voice/mei/mei_normal.htsvoice",
		Player:			"aplay",
	}

	fi, _ := os.Stat(tmpPath)
	if nil == fi {
		os.Mkdir(tmpPath, 0777)
	}
	return &s
}


func	(self *Voice)Say(text string) {
	fmt.Println(text)
	textFile	:= self.TmpPath + self.TmpTxt
	ioutil.WriteFile(textFile, []byte(text), 0644)

	outFile	:= self.TmpPath + self.TmpWav
    cmd := exec.Command(self.OpenJtalk)
    cmd.Args = []string {
      "-x", self.OpenJtalkDic,
      "-m", self.OpenJtalkVoice,
      textFile,
      "-ow", outFile,
    }
    cmd.Run()

    cmd = exec.Command(self.Player, outFile)
    cmd.Run()
}

func	(self *Voice)Play(sfx string) {
	sfxFile	:= self.SfxPath + sfx
    cmd := exec.Command(self.Player, sfxFile)
    cmd.Run()
}

