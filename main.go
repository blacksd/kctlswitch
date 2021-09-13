package main

import (
	"kctlswitch/lib"
)

func main() {
	lib.BuildKctlList("https://github.com/kubernetes/kubectl.git")
}
