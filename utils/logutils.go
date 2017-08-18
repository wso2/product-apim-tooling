package utils

import (
	"fmt"
	"os"
)

var loglnFunc = doNothinglnFunc

var logfFunc = doNothingfFunc

func verbosePrintlnFunc(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func doNothinglnFunc(v ...interface{}) {
}

func verbosePrintfFunc(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func doNothingfFunc(format string, a ...interface{}) {
}

func Logln(a ...interface{}) {
	loglnFunc(a...)
}

func Logf(format string, a ...interface{}) {
	logfFunc(format, a...)
}

func EnableVerboseMode() {
	loglnFunc = verbosePrintlnFunc
	logfFunc = verbosePrintfFunc
}
