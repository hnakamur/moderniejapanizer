package main

import (
	"fmt"
	"os"

	"github.com/hnakamur/moderniejapanizer"
	"github.com/hnakamur/w32syscall"
)

func main() {
	err := moderniejapanizer.Reboot(w32syscall.SHTDN_REASON_MINOR_SECURITYFIX)
	if err != nil {
		fmt.Printf("Error while rebooting: %v\n", err)
		os.Exit(1)
	}
}
