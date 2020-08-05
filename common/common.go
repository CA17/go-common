package common

import (
	"crypto/rand"
	sha256_ "crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
)

var (
	TIME_LAYOUT  = "2006-01-02 15:04:05"
	EmptyList    = []interface{}{}
	EmptyTime, _ = time.Parse("2006-01-02 15:04:05 Z0700 MST", "1979-11-30 00:00:00 +0000 GMT")
)

const (
	NA       = "N/A"
	ENABLED  = "enabled"
	DISABLED = "disabled"
)

// print usage
func Usage(str string) {
	fmt.Fprintf(os.Stderr, str)
	flag.PrintDefaults()
}

// 创建目录
func MakeDir(path string) {
	f, err := os.Stat(path)
	if err != nil || f.IsDir() == false {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			log.Println("create dir fail！", err)
			return
		}
	}
}

// 判断文件是否存在
func FileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

// 判断文件目录是否存在
func DirExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && info.IsDir()
}

// panic error
func Must(err error) {
	if err != nil {
		panic(errors.WithStack(err))
	}
}

func MustCallBefore(err error, callbefore func()) {
	if err != nil {
		callbefore()
		panic(errors.WithStack(err))
	}
}

func Must2(v interface{}, err error) interface{} {
	Must(err)
	return v
}

func IgnoreError(v interface{}, err error) interface{} {
	return v
}

// 生成 UUID
func UUID() string {
	unix32bits := uint32(time.Now().UTC().Unix())
	buff := make([]byte, 12)
	numRead, err := rand.Read(buff)
	if numRead != len(buff) || err != nil {
		Must(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x-%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:])
}

// 生成  int64
func UUIDint64() int64 {
	node, err := snowflake.NewNode(1000)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return node.Generate().Int64()
}

func UUIDBase32() (string, error) {
	node, err := snowflake.NewNode(mrand.Int63n(64))
	if err != nil {
		return "", err
	}
	id := node.Generate()
	// Print out the ID in a few different ways.
	return id.Base32(), nil
}

// 验证Email格式
func ValidateEmail(email string) (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return
}

// 转换到大驼峰格式
func ToCamelCase(str string) string {
	temp := strings.Split(str, "_")
	for i, r := range temp {
		temp[i] = strings.Title(r)
	}
	return strings.Join(temp, "")
}

var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// 转换到下划线格式
func ToSnakeCase(str string) string {
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

// hash 方法
func Sha256Hash(src string) string {
	h := sha256_.New()
	h.Write([]byte(src))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// 判断字符串是否在列表中
func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfEmptyStr(src string, defval string) string {
	if src == "" {
		return defval
	}
	return src
}

func split(s string, size int) []string {
	ss := make([]string, 0, len(s)/size+1)
	for len(s) > 0 {
		if len(s) < size {
			size = len(s)
		}
		ss, s = append(ss, s[:size]), s[size:]

	}
	return ss
}

func File2Base64(file string) string {
	data := Must2(ioutil.ReadFile(file))
	return base64.StdEncoding.EncodeToString(data.([]byte))
}

func Base642file(b64str string, file string) error {
	data, err := base64.StdEncoding.DecodeString(b64str)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 777)
}

func CopyFile(dstName, srcName string, perm os.FileMode) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func IfNA(src string, defval string) string {
	if src == "N/A" || src == "" {
		return defval
	}
	return src
}

func EmptyToNA(src string) string {
	if src == "" {
		return NA
	}
	return src
}

func SetEmptyStrToNA(t interface{}) {
	d := reflect.TypeOf(t).Elem()
	for j := 0; j < d.NumField(); j++ {
		ctype := d.Field(j).Type.String()
		if ctype == "string" {
			val := reflect.ValueOf(t).Elem().Field(j)
			if val.String() == "" {
				val.SetString(NA)
			}
		}
	}
}

// 解析表单时间
// t 表单时间字符串
// hms 时分秒
func ParseFormTime(t, hms string) (time.Time, error) {
	if len(t) < 10 {
		return EmptyTime, errors.New("时间格式不正确， 必须是yyyy-MM-dd")
	}
	timestr := t[:10] + " " + hms
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return EmptyTime, err
	}
	return time.ParseInLocation(TIME_LAYOUT, timestr, loc)
}

func UrlJoin(hurl string, elm ...string) string {
	u, err := url.Parse(hurl)
	Must(err)
	u.Path = path.Join(u.Path, path.Join(elm...))
	return u.String()
}

func Interface2String(inter interface{}) string {
	switch inter.(type) {
	case string:
		return inter.(string)
	case int64:
		return strconv.FormatInt(inter.(int64), 10)
	case int:
		return strconv.Itoa(inter.(int))
	case float32:
		return fmt.Sprintf("%f", inter.(float32))
	case float64:
		return fmt.Sprintf("%f", inter.(float64))
	default:
		return fmt.Sprintf("%v", inter)
	}
}

func GetValuesFromMapByKey(valmap map[string]string, keys []string) string {
	var result = make([]string, 0)
	for _, key := range keys {
		val, ok := valmap[key]
		if ok && val != "" {
			result = append(result, val)
		}
	}
	return strings.Join(result, ",")
}

func FmtSecondDesc(secs int64) string {

	if secs > 60 && secs < 3600 {
		m := secs / 60
		return fmt.Sprintf("%d分钟", m)
	}

	if secs > 3600 && secs < 86400 {
		h := secs / 3600
		m := secs % 3600 / 60
		return fmt.Sprintf("%d小时%d分钟", h, m)
	}

	if secs > 86400 {
		d := secs / 86400
		h := secs % 86400 / 3600
		m := secs % 86400 % 3600 / 60
		return fmt.Sprintf("%d天%d小时%d分钟", d, h, m)
	}

	return fmt.Sprintf("%d秒", secs)
}

