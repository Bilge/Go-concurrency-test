package main

import "strings"

func pad(str string, n int) string {
	if len(str) < n {
		return strings.Repeat(`0`, n-len(str)) + str
	}

	return str
}
