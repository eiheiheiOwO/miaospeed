package main

import (
	"github.com/miaokobot/miaospeed/utils"
)

var COMPILATIONTIME string
var BUILDCOUNT string
var COMMIT string
var BRAND string
var VERSION string

func main() {
	utils.COMPILATIONTIME = COMPILATIONTIME
	utils.BUILDCOUNT = BUILDCOUNT
	utils.COMMIT = COMMIT
	utils.BRAND = BRAND
	utils.VERSION = VERSION
	RunCli()
}
