package timeutil

import (
	"errors"
	"time"
)

const (
	Datetime14Layout      = "20060102150405"
	Datetime8Layout      = "20060102"
	Datetime6Layout      = "200601"
	YYYYMMDDHHMMSS_LAYOUT = "2006-01-02 15:04:05"
	YYYYMMDDHHMM_LAYOUT   = "2006-01-02 15:04"
	YYYYMMDD_LAYOUT       = "2006-01-02"
)

var (
	ShangHaiLOC, _ = time.LoadLocation("Asia/Shanghai")
	EmptyList    = make([]interface{}, 0)
	EmptyTime, _ = time.Parse("2006-01-02 15:04:05 Z0700 MST", "1979-11-30 00:00:00 +0000 GMT")
)

func FmtCstDatetime(t time.Time) string {
	return t.In(ShangHaiLOC).String()
}

// yyyy-MM-dd hh:mm:ss 年-月-日 时:分:秒
func FmtDatetimeString(t time.Time) string {
	return t.Format(YYYYMMDDHHMMSS_LAYOUT)
}

// yyyy-MM-dd hh:mm 年-月-日 时:分
func FmtDatetimeMString(t time.Time) string {
	return t.Format(YYYYMMDDHHMM_LAYOUT)
}

// yy-MM-dd 年-月-日
func FmtDateString(t time.Time) string {
	return t.Format(YYYYMMDD_LAYOUT)
}

// yyyyMMddhhmmss 年月日时分秒
func FmtDatetime14String(t time.Time) string {
	return t.Format(Datetime14Layout)
}

// yyyyMMdd 年月日
func FmtDatetime8String(t time.Time) string {
	return t.Format(Datetime8Layout)
}

// yyyyMM  年月
func FmtDatetime6String(t time.Time) string {
	return t.Format(Datetime6Layout)
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
	return time.ParseInLocation(YYYYMMDDHHMMSS_LAYOUT, timestr, loc)
}


func parseWithLocation(locationName string, timeStr string, format string) (time.Time, error) {
	if l, err := time.LoadLocation(locationName); err != nil {
		println(err.Error())
		return time.Time{}, err
	} else {
		lt, _ := time.ParseInLocation(format, timeStr, l)
		return lt, nil
	}
}
