package main

import (
	"kctlswitch/lib"
)

func main() {
	lib.KctlVersionsList("<= 1.7")
	lib.DownloadKctl("v1.12.3", ".")
}
