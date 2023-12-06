package mycmd

import "testing"

func TestCombinedOutput(t *testing.T) {
	if output, err := Execbash("docker images"); err != nil {
		t.Error("err is :", err)
	} else {
		t.Log("output is :\n", output)
	}
}
