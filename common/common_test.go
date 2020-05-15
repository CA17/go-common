package common

import (
	"path"
	"strings"
	"testing"
)

func TestPath(t *testing.T) {
	aa := "ddd_css.css"
	t.Log(path.Ext(aa))
	t.Log(strings.TrimSuffix(aa, path.Ext(aa)))
}
