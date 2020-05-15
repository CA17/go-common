package validutil

import (
	"regexp"
	"unicode"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/ca17/go-common"
)

const (
	// 匹配大陆电话
	cnPhoneStr = `((\d{3,4})-?)?` + // 区号
		`\d{7,8}` + // 号码
		`(-\d{1,4})?` // 分机号，分机号的连接符号不能省略。

	// 匹配大陆手机号码
	cnMobileStr = `(0|\+?86)?` + // 匹配 0,86,+86
		`(13[0-9]|` + // 130-139
		`14[57]|` + // 145,147
		`15[0-35-9]|` + // 150-153,155-159
		`17[0678]|` + // 170,176,177,17u
		`18[0-9])` + // 180-189
		`[0-9]{8}`

	// 匹配邮箱
	emailStr = `[\w.-]+@[\w_-]+\w{1,}[\.\w-]+`

	// 匹配IP4
	ip4Str = `((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`

	// 匹配IP6，参考以下网页内容：
	// http://blog.csdn.net/jiangfeng08/article/details/7642018
	ip6Str = `(([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|` +
		`(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|` +
		`(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|` +
		`(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))`

		// 同时匹配IP4和IP6
	ipStr = "(" + ip4Str + ")|(" + ip6Str + ")"

	// 匹配域名
	domainStr = `[a-zA-Z0-9][a-zA-Z0-9_-]{0,62}(\.[a-zA-Z0-9][a-zA-Z0-9_-]{0,62})*(\.[a-zA-Z][a-zA-Z0-9]{0,10}){1}`

	// 匹配URL
	urlStr = `((https|http|ftp|rtsp|mms)?://)?` + // 协议
		`(([0-9a-zA-Z]+:)?[0-9a-zA-Z_-]+@)?` + // pwd:user@
		"(" + ipStr + "|(" + domainStr + "))" + // IP或域名
		`(:\d{1,4})?` + // 端口
		`(/+[a-zA-Z0-9][a-zA-Z0-9_.-]*/*)*` + // path
		`(\?([a-zA-Z0-9_-]+(=[a-zA-Z0-9_-]*)*)*)*` // query
)

var (
	Validtool  *validator.Validate
	ValidTrans *ut.Translator
)

func NewValidatorTrans() (*validator.Validate, *ut.Translator, error) {
	var uni *ut.UniversalTranslator
	var validate *validator.Validate
	zh := zh.New()
	uni = ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	validate = validator.New()
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, nil, err
	}
	return validate, &trans, nil
}

func init() {
	var err error
	Validtool, ValidTrans, err = NewValidatorTrans()
	common.Must(err)
}

func regexpCompile(str string) *regexp.Regexp {
	return regexp.MustCompile("^" + str + "$")
}

var (
	email    = regexpCompile(emailStr)
	ip4      = regexpCompile(ip4Str)
	ip6      = regexpCompile(ip6Str)
	ip       = regexpCompile(ipStr)
	url      = regexpCompile(urlStr)
	cnPhone  = regexpCompile(cnPhoneStr)
	cnMobile = regexpCompile(cnMobileStr)
)

// 判断val是否能正确匹配exp中的正则表达式。
// val可以是[]byte, []rune, string类型。
func isMatch(exp *regexp.Regexp, val interface{}) bool {
	switch v := val.(type) {
	case []rune:
		return exp.MatchString(string(v))
	case []byte:
		return exp.Match(v)
	case string:
		return exp.MatchString(v)
	default:
		return false
	}
}

// 验证中国大陆的电话号码。支持如下格式：
//  0578-12345678-1234
//  057812345678-1234
// 若存在分机号，则分机号的连接符不能省略。
func IsCnPhone(val interface{}) bool {
	return isMatch(cnPhone, val)
}

// 验证中国大陆的手机号码
func IsCnMobile(val interface{}) bool {
	return isMatch(cnMobile, val)
}

// 验证一个值是否标准的URL格式。支持IP和域名等格式
func IsURL(val interface{}) bool {
	return isMatch(url, val)
}

// 验证一个值是否为IP，可验证IP4和IP6
func IsIP(val interface{}) bool {
	return isMatch(ip, val)
}

// 验证一个值是否为IP6
func IsIP6(val interface{}) bool {
	return isMatch(ip6, val)
}

// 验证一个值是滞为IP4
func IsIP4(val interface{}) bool {
	return isMatch(ip4, val)
}

// 验证一个值是否匹配一个邮箱。
func IsEmail(val interface{}) bool {
	return isMatch(email, val)
}

func CheckPassword(s string) bool {
	var hasNumber, hasLetter bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsLetter(c):
			hasLetter = true
		}
	}
	return hasNumber && hasLetter
}

func CheckStrongPassword(s string) bool {
	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case c == '#' || c == '|':
			return false
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasNumber && hasUpperCase && hasLowercase && hasSpecial
}
