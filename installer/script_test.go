package installer

import (
	"os"
	"testing"
)

func TestOS(t *testing.T) {
	s, _ := os.Getwd()
	t.Log(s)
}
