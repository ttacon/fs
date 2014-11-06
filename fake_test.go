package fs

import (
	"os"
	"testing"
)

func Test_FakeOs_Getenv(t *testing.T) {
	f := FakeOS()
	v := f.Getenv("GOPATH")
	if v != "" {
		t.Errorf("expected $GOPATH to be empty, was: %v", v)
	}
}

func Test_FakeOs_Chdir_Doesnt_Exist(t *testing.T) {
	f := FakeOS()
	err := f.Chdir("/i/dont/exist")
	if err == nil {
		t.Errorf("expected err, there was none")
	} else if v, ok := err.(*os.PathError); !ok {
		t.Errorf("expected error of type *os.PathError, was: %v", v)
	}
}

func Test_FakeOs_Chdir(t *testing.T) {
	f := FakeOS()
	if err := f.Chdir(f.TempDir()); err != nil {
		t.Errorf("expected to be able to Chdir to the temp dir, err: %v", err)
	}
}

func Test_FakeOs_Setenv(t *testing.T) {
	f := FakeOS()
	err := f.Setenv("GOPATH", "/mars/valley")
	if err != nil {
		t.Errorf("expected to be able to Setenv() w/o err, err: %v", err)
	}

	val := f.Getenv("GOPATH")
	if val != "/mars/valley" {
		t.Errorf("expected val to be \"/mars/valley\", was: %q", val)
	}
}

func Test_FakeOs_Environ(t *testing.T) {
	f := FakeOS()
	envVars := f.Environ()
	if len(envVars) != 0 {
		t.Errorf("Expected Environ() to be empty, was %v", envVars)
	}

	err := f.Setenv("GOPATH", "/awesome/path")
	if err != nil {
		t.Errorf("there's no way there's an error")
	}

	envVars = f.Environ()
	if len(envVars) != 1 {
		t.Errorf(
			"expected Environ() to be GOPATH=/awesome/path, was: %v",
			envVars)
	}

	if envVars[0] != "GOPATH=/awesome/path" {
		t.Errorf(
			"expected Environ() to be GOPATH=/awesome/path, was: %v",
			envVars)
	}
}

func Test_FakeOs_Clearenv(t *testing.T) {
	f := FakeOS()
	fr, _ := f.(*fakeOS)
	if len(fr.envVars) > 0 {
		t.Errorf(
			"envVars should be empty on default FakeOS startup, was: %v",
			fr.envVars)
	}

	fr.envVars["GOPATH"] = "/best/path/ever"

	f.Clearenv()
	if len(fr.envVars) > 0 {
		t.Errorf(
			"envVars should be empty on default FakeOS startup, was: %v",
			fr.envVars)
	}
}
