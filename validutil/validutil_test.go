package validutil

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestTransZh(t *testing.T) {
	a, b, _ := NewValidatorTrans()
	err := a.Var("a", "max=0")
	if err != nil {
		l := err.(validator.ValidationErrors).Translate(*b)
		t.Log(l)
	}

}
