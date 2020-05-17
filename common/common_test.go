package common

import (
	"fmt"
	"path"
	"strings"
	"testing"
)

func TestPath(t *testing.T) {
	aa := "ddd_css.css"
	t.Log(path.Ext(aa))
	t.Log(strings.TrimSuffix(aa, path.Ext(aa)))
}

type S struct {
	Name   string
	Age    int    `update:"age"`
	Value  string `update:"value"`
	Value1 string
	Value2 string
	Value3 string
	Value4 string
	Value5 string
}

func TestSetEmptyStrToNA(t *testing.T) {
	v := S{
		Name:   "helo",
		Value:  "",
		Value1: "",
		Value2: "",
		Value3: "",
		Value4: "",
		Value5: "",
	}
	SetEmptyStrToNA(&v)
	fmt.Println(v)
}

func BenchmarkSetEmptyStrToNA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := S{
			Name:   "name",
			Value:  "",
			Value1: "",
			Value2: "",
			Value3: "",
			Value4: "",
			Value5: "",
		}
		SetEmptyStrToNA(&v)
	}
}

func BenchmarkSetEmptyStrToNAn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := S{
			Name:   "name",
			Value:  "",
			Value1: "",
			Value2: "",
			Value3: "",
			Value4: "",
			Value5: "",
		}
		v.Value = _NA
		v.Value1 = _NA
		v.Value2 = _NA
		v.Value3 = _NA
		v.Value4 = _NA
		v.Value5 = _NA
	}
}

func TestSetMapFrom(t *testing.T) {
	m := map[string]interface{}{}
	v := S{
		Name:   "name",
		Value:  "vvv",
		Value1: "",
		Value2: "",
		Value3: "",
		Value4: "",
		Value5: "",
	}
	SetMapFrom(m, &v)
	fmt.Printf("%+v", m)

}
