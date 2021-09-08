package util

import (
	"github.com/aoaostar/v8cdn_panel/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func PrintError(err error, c *gin.Context) {
	// 获取validator.ValidationErrors类型的errors
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		if err != nil {
			c.JSON(http.StatusOK, Msg(
				"error",
				err.Error(), nil,
			))
			return
		}
	}
	// validator.ValidationErrors类型错误则进行翻译
	for _, v := range errs.Translate(config.Trans) {
		c.JSON(http.StatusOK, Msg(
			"error",
			v, errs.Translate(config.Trans),
		))
		return
	}

}
