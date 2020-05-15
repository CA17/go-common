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
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	EmptyList = new([]interface{})
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
			log.Fatalf("create dir fail！", err)
			return
		}
	}
}

// 判断文件是否存在
func FileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

// panic error
func Must(err error) {
	if err != nil {
		panic(err)
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

func IfEmpty(src interface{}, defval interface{}) interface{} {
	if IsEmpty(src) {
		return defval
	}
	return src
}
func IfEmptyStr(src string, defval string) string {
	if src == "" {
		return defval
	}
	return src
}

// IsEmpty checks if a value is empty or not.
// A value is considered empty if
// - integer, float: zero
// - bool: false
// - string, array: len() == 0
// - slice, map: nil or len() == 0
// - interface, pointer: nil or the referenced value is empty
func IsEmpty(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Invalid:
		return true
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsEmpty(v.Elem().Interface())
	case reflect.Struct:
		v, ok := value.(time.Time)
		if ok && v.IsZero() {
			return true
		}
	}

	return false
}

func IsNotEmpty(value interface{}) bool {
	return !IsEmpty(value)
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
