package configs

import (
	"MyGram/app"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni *ut.UniversalTranslator
)

func ConfigValidator() {
	app.Validator = validator.New()
	en := en.New()
	uni = ut.New(en, en)

	app.Trans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(app.Validator, app.Trans)

}
