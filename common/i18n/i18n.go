package i18n

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/ypli0629/yoyo/common/logger"
)

var Trans ut.Translator

func Setup() {
	setupValidatorI18n()
}

func setupValidatorI18n() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zh := zh.New()
		uts := ut.New(zh, zh)
		var err bool
		Trans, err = uts.GetTranslator("zh")
		zhTranslations.RegisterDefaultTranslations(v, Trans)
		if !err {
			logger.Panicf("Faile to setup i18n: %s", err)
		}
	}
}

func Translate(err error) error {
	errs := err.(validator.ValidationErrors)
	return errors.New(fmt.Sprint(errs.Translate(Trans)))
}
