// This is the main package for the qk command.
// It is responsible for parsing the command line arguments and executing the appropriate command.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"qk/commands"
	"qk/internal"
	"qk/internal/types/definitions"
	"slices"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	args := os.Args[1:]

	if *debug {
		debug_i := slices.Index(args, "--debug")
		copy(args[debug_i:], args[debug_i+1:])
		args[len(args)-1] = ""
		args = args[:len(args)-1]
		fmt.Println("Arguments:", args)
		fmt.Println("Debug mode:", *debug)
	}

	var customCommands []definitions.CustomCommand
	customCommands = append(customCommands, commands.GetDefaultCommands()...)
	customCommands = append(customCommands, commands.GetUserCommands()...)
	customCommands = append(customCommands, commands.GetWorkCommands()...)

	ret := internal.Qk(args, customCommands, debug)
	// ret = strings.ReplaceAll(ret, " ", "_")
	if ret != "" {
		if *debug {
			fmt.Println("\nGenerated Command:")
		}
		fmt.Println(ret)
	} else {
		log.Fatal("No matching command found.")
	}
}
