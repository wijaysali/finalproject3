package i18nconfig

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var EngBundle *i18n.Bundle

func Configi18n() {
	EngBundle = i18n.NewBundle(language.English)
	EngBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	EngBundle.MustLoadMessageFile("i18n/en.json")
	EngBundle.MustLoadMessageFile("i18n/id.json")

}
