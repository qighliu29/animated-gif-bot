package main

import (
	"fmt"

	ct "github.com/daviddengcn/go-colortext"
)

func failed(format string, args ...interface{}) {
	ct.Foreground(ct.Red, false)
	fmt.Printf(format, args...)
	ct.ResetColor()
}

func success(format string, args ...interface{}) {
	ct.Foreground(ct.Green, false)
	fmt.Printf(format, args...)
	ct.ResetColor()
}
