package utils

import (
	"encoding/json"
	"os"
	"os/exec"
)

func ReadMessage(b []byte, dest interface{}) error {
	err := json.Unmarshal(b, dest)
	return err
}

func RunCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
