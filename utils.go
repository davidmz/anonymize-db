package main

import (
	"fmt"

	"github.com/davidmz/mustbe"
)

func wrapError(errTemplate string, args ...interface{}) func(err error) {
	return func(err error) {
		mustbe.Thrown(fmt.Errorf(errTemplate, append(args, err)...))
	}
}

func mustbeDone(foo func(), errTemplate string, args ...interface{}) {
	defer mustbe.Catched(wrapError(errTemplate, args...))
	foo()
}
