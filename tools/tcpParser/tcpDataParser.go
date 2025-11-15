package main

import (
	"encoding/hex"
	"strings"

	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	// For Testing Data collected from WireShark
	inputs := [1][2]string{
		{
			"0251007815f6235e64c5711052897fa07a9fbd98cdee2b76b4054977c6045a8cd5945decb16d00db3c2d10a1dc771512c88f3bbce541b28bdbdbad19548d2ce83dc68767425beb1571b37cb28ca895f85fbdbe1bfb7da738ae1220023d00af8ef664c9030e52ace135f1b3ef7cd6c4bc94a7269d2c269ebdd2ebf4af50a4de21c1b5220ee09bd22d61f11e97c3cb4e20ab203ca8d29a0e5b50d5c8799438025100bfa3d08a53c69715ef0abdb3f0ea2824eb162f58010581c13240267619561f0a68fb5c900d0892f116983d278c2054361af63dcde591281319548d2ce83dc68767425beb1571b37cb28ca895f85fbdbe1bfb7da738ae1220",
			"Received",
		},
	}

	for _, v := range inputs {
		input, _ := hex.DecodeString(v[0])
		output := utils.ParseBody(input, utils.CipherKey)
		log.Printf("%s", strings.Repeat("-", 25))
		log.Printf("%s the following commands", v[1])
		utils.TranslateMsg(output)
	}
}
