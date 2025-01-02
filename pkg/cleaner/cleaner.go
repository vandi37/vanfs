package cleaner

import (
	"os"
	"os/exec"
	"runtime"
)

type Cleaner map[string]func()

func New() Cleaner {
	clear := make(map[string]func())
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["default"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	return clear
}

func (c Cleaner) Clear() {
	value, ok := c[runtime.GOOS]
	if ok {
		value()
	} else {
		c["default"]()
	}
}
