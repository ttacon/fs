package fs

import (
	"fmt"
	"testing"
)

func Test_OS(t *testing.T) {
	o := DefaultOS()
	fmt.Println(o.Getenv("GOPATH"))
}

func Test_FakeOs(t *testing.T) {
	o := FakeOS()
	if val := o.Getenv("GOPATH"); val != "" {
		t.Error("val was not empty, was: ", val)
	}
}
