package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/orashus/rtodo/config"
	"github.com/orashus/rtodo/constants"
	"github.com/orashus/rtodo/router"
	"github.com/orashus/rtodo/store"
	"github.com/orashus/rtodo/utils"
)

const (
	version   = "1.0.0"
	appName   = "rtodo"
	FILE_PATH = "/tmp/r_apps_rtodo.json"
	// FILE_PATH = "test-todos.json"
)

func init() {
	/*
		In Go, the init() function is a special, reserved function used for package initialization
		that runs automatically before the main() function.
		You can have multiple init() functions in a package, and they will run in the order they are defined.
		Init functions are commonly used for:
		- Setting up package-level variables
		- Initializing package-level data structures
		- Registering package-level functions with the runtime
		- Setting up package-level configuration
		- Validating package-level dependencies
		- Initializing package-level logging
	*/
	config.Config.Configure(FILE_PATH)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s  version %s\n", appName, version)
		fmt.Println("Please provide a command")
		return
	}

	var shouldPrint bool

	command, input, tags := utils.ParseInput(os.Args[1:])

	if slices.Contains(tags, constants.TAGS.PRINT) || slices.Contains(tags, constants.TAGS.LIST) {
		shouldPrint = true
	}

	if slices.Contains(tags, constants.TAGS.VERSION) {
		fmt.Printf("%s  version %s\n", appName, version)
		return
	}

	switch command {
	case constants.COMMANDS.LIST:
		router.ListHandler(tags)
	case constants.COMMANDS.DELETE:
		fallthrough
	case constants.COMMANDS.RM:
		fallthrough
	case constants.COMMANDS.REMOVE:
		if len(input) == 0 {
			fmt.Println("Please provide an ID to remove")
			return
		}

		router.RemoveHandler(input, shouldPrint)
	case constants.COMMANDS.RMC:
		router.RemoveCompletedHandler(shouldPrint)
	case constants.COMMANDS.CLEAR:
		router.ClearHandler(shouldPrint)
	case constants.COMMANDS.COMPLETE:
		fallthrough
	case constants.COMMANDS.DONE:
		fallthrough
	case constants.COMMANDS.FINISH:
		fallthrough
	case constants.COMMANDS.CHECK:
		fallthrough
	case constants.COMMANDS.MARK:
		if len(input) == 0 {
			fmt.Println("Please provide an ID to mark as complete")
			return
		}

		router.MarkAsCompleteHandler(input, shouldPrint)
	case constants.COMMANDS.ADD:
		if len(input) == 0 {
			fmt.Println("Please provide a title to add")
			return
		}

		router.AddHandler(input, shouldPrint)
	case constants.COMMANDS.UPDATE:
		if len(input) < 2 {
			fmt.Println("Please provide an ID and a title to update")
			return
		}

		router.UpdateHandler(input[0], input[1], shouldPrint)
	default:
		fmt.Println("Invalid command")
		utils.PrintDelimiter()

		if shouldPrint {
			todoList, _ := store.LoadTodos(FILE_PATH)
			utils.PrintTodos(&todoList)
		}
		return
	}
}
