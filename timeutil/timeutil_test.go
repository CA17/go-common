package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	fmt.Println(time.Now().Format("20060102"))
}

func TestFormatLenTime(t *testing.T) {
	fmt.Println(FmtDatetime14String(time.Now()))
	fmt.Println(FmtDatetime8String(time.Now()))
	fmt.Println(FmtDatetime6String(time.Now()))
	fmt.Println(FmtDateString(time.Now()))
	fmt.Println(FmtDatetimeString(time.Now()))
	fmt.Println(FmtDatetimeMString(time.Now()))

}

func TestFmtCstTime(t *testing.T) {
	fmt.Println(FmtCstDatetime(time.Now()))
}
