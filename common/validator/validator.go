package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni   *ut.UniversalTranslator
	Trans ut.Translator
)

func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		zh := zh.New()
		uni = ut.New(en, zh)
		Trans, _ = uni.GetTranslator("zh")
		zh_trans.RegisterDefaultTranslations(v, Trans)
	}
}
