package main

import (
	"log"

	"github.com/bernylinville/wtf/flags"
)

func main() {
	// 日志格式 2009/01/23 01:23:23 message
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Paras and handle flags
	flags := flags.NewFlags()
	flags.Parse()

	// Load the configuration file

}
