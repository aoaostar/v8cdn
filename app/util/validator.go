package util

import (
	"github.com/aoaostar/v8cdn_panel/pkg/Validator"
	"github.com/go-playground/validator/v10"
)

func FomateValidateError(err error) (string, interface{}) {
	// 获取validator.ValidationErrors类型的errors
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		if err != nil {

			return err.Error(), errs
		}
	}
	// validator.ValidationErrors类型错误则进行翻译
	for _, v := range errs.Translate(Validator.Trans) {
		return v, errs.Translate(Validator.Trans)
	}
	return err.Error(), errs

}
