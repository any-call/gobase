package mycmd

import (
	"os/exec"
	"testing"
)

func TestCombinedOutput(t *testing.T) {
	if output, err := Execbash("docker images", func(c *exec.Cmd) {
	}, true); err != nil {
		t.Error("err is :", err)
	} else {
		t.Log("output is :\n", output)
	}
}
