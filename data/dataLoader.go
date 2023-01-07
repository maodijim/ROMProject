package data

import (
	_ "embed"
)

//go:embed buff.json
var BuffJson []byte

//go:embed exchangeItems.json
var ExchangeJson []byte

//go:embed items.json
var ItemsJson []byte

//go:embed skills.json
var SkillsJson []byte

//go:embed monsters.json
var MonstersJson []byte
