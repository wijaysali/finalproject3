package app

import (
	"MyGram/utils"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var Db *gorm.DB
var Validator *validator.Validate
var Port strings.Builder
var Trans ut.Translator

func Validate(c *fiber.Ctx, data interface{}) error {
	err := Validator.Struct(data)
	if err != nil {

		// translate all error at once
		errs := err.(validator.ValidationErrors)

		strBuilder := strings.Builder{}
		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned and you'll see something surprising
		// translations are i18n aware!!!!
		// eg. '10 characters' vs '1 character'
		for idx, e := range errs {
			strBuilder.WriteString(e.Translate(Trans))
			if len(errs) != idx+1 {
				strBuilder.WriteString(", ")
			}
		}
		utils.ResponsePanic(fiber.StatusBadRequest, strBuilder.String())

	}
	return nil

}
