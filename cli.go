package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/miaokobot/miaospeed/utils"
)

var cmdName string = "miaospeed"

const logo string = " __  __ _            ____                      _ \n|  \\/  (_) __ _  ___/ ___| _ __   ___  ___  __| |\n| |\\/| | |/ _` |/ _ \\___ \\| '_ \\ / _ \\/ _ \\/ _` |\n| |  | | | (_| | (_) |__) | |_) |  __/  __/ (_| |\n|_|  |_|_|\\__,_|\\___/____/| .__/ \\___|\\___|\\__,_|\n                          |_|                    "

type SubCliType string

const (
	SCTMisc       SubCliType = "misc"
	SCTServer     SubCliType = "server"
	SCTScriptTest SubCliType = "script"
)

func RunCli() {
	subCmd := SubCliType("")
	if len(os.Args) >= 2 {
		subCmd = SubCliType(os.Args[1])
	}

	cmdName = path.Base(os.Args[0])
	switch subCmd {
	case SCTMisc:
		RunCliMisc()
	case SCTServer:
		RunCliServer()
	case SCTScriptTest:
		RunCliScriptTest()
	default:
		RunCliDefault()
	}
}

func RunCliDefault() {
	sflag := flag.NewFlagSet(cmdName, flag.ExitOnError)

	versionOnly := sflag.Bool("version", false, "display version and exit")
	sflag.Parse(os.Args[1:])

	if *versionOnly {
		fmt.Println(logo)
		fmt.Printf("version: %s\n", utils.VERSION)
		fmt.Printf("commit: %s\n", utils.COMMIT)
		fmt.Printf("compilation time: %s\n", utils.COMPILATIONTIME)
		os.Exit(0)
	}

	sflag.Usage()

	fmt.Printf("\n")
	fmt.Printf("Subcommands of %s:\n", cmdName)
	fmt.Printf("  server\n")
	fmt.Printf("        start the miaospeed backend as a server.\n")
	fmt.Printf("  script\n")
	fmt.Printf("        run a temporary script test to test the correctness of your script.\n")
	fmt.Printf("  misc\n")
	fmt.Printf("        other utility toolkit provided by miaospeed.\n")
	fmt.Printf("Run this command to see the usage of the server option: \n")
	fmt.Printf("  %s server -help\n", cmdName)
	os.Exit(0)
}

func parseFlag(sflag *flag.FlagSet) {
	verboseMode := sflag.Bool("verbose", false, "whether to print out systems log")

	sflag.Parse(os.Args[2:])

	if *verboseMode {
		utils.VerboseLevel = utils.LTLog
	}
}
