package converter

import (
	"git.circuitco.de/self/bellbox"
)

type Translation interface {
	Handle(input []byte) bellbox.Message
}

type ConvertRule struct {
	Target     string
	Path       string
	Ruleset    string
	SenderName string
}

type TranslationConfig struct {
	Bellbox string
	Default ConvertRule
	Convert []ConvertRule
}
