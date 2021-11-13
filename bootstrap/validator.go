package bootstrap

import (
	"github.com/aoaostar/v8cdn_panel/pkg/Validator"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		//\/:*?”<>|
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		Validator.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			log.Errorf("uni.GetTranslator(%s) failed", locale)
			return
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, Validator.Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, Validator.Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, Validator.Trans)
		}
		return
	}
	return
}
func InitValidator() {

	if err := InitTrans("zh"); err != nil {
		log.Debugf("init trans failed, err:%v\n", err)
		return
	}
}
