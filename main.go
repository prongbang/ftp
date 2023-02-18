package main

import (
	"github.com/prongbang/ftp/goftp"
	"github.com/prongbang/ftp/gosftp"
)

func main() {
	if false {
		goftp.Run()
	} else {
		gosftp.Run()
	}
}
