package router

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/orashus/rtodo/config"
	"github.com/orashus/rtodo/constants"
	"github.com/orashus/rtodo/store"
	"github.com/orashus/rtodo/types"
	"github.com/orashus/rtodo/utils"
)

func ListHandler(tags []string) {
	todoList, err := store.LoadTodos(config.Config.FILE_PATH)
	if err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error loading todos ->", err)
		return
	}

	if slices.Contains(tags, constants.TAGS.COMPLETED) {
		res := []types.Todo{}

		for _, todo := range todoList {
			if todo.Completed {
				res = append(res, todo)
			}
		}

		utils.PrintTodos(&res)
		return
	}

	utils.PrintTodos(&todoList)
}

func RemoveHandler(ids []string, shouldPrint bool) {
	todoList, err := store.LoadTodos(config.Config.FILE_PATH)
	if err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error loading todos ->", err)
		return
	}

	hasSome := false

	todoList = slices.DeleteFunc(todoList, func(todo types.Todo) bool {
		if slices.Contains(ids, todo.ID) {
			hasSome = true
		}
		return slices.Contains(ids, todo.ID)
	})

	if !hasSome {
		utils.PrintDelimiter()
		fmt.Println("Todo(s) not found")
		if shouldPrint {
			utils.PrintTodos(&todoList)
		}
		return
	}

	if err := store.SaveTodos(config.Config.FILE_PATH, &todoList); err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error saving todos:", err)
		return
	}

	utils.PrintDelimiter()
	fmt.Printf("Todo(s) with ID(s) '%+v' removed successfully\n", ids)
	if shouldPrint {
		utils.PrintTodos(&todoList)
	}
}

func RemoveCompletedHandler(shouldPrint bool) {
	todoList, err := store.LoadTodos(config.Config.FILE_PATH)
	if err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error loading todos ->", err)
		return
	}

	todoList = slices.DeleteFunc(todoList, func(todo types.Todo) bool {
		return todo.Completed
	})

	if err := store.SaveTodos(config.Config.FILE_PATH, &todoList); err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error saving todos:", err)
		return
	}

	utils.PrintDelimiter()
	fmt.Println("Completed todos removed successfully")
	if shouldPrint {
		utils.PrintTodos(&todoList)
	}
}

func ClearHandler(shouldPrint bool) {
	err := os.Remove(config.Config.FILE_PATH)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) { // or I could do os.IsNotExist(err)
			/*
				This just checks if the error is known to report that a file or directory does not exist.
				os.IsNotExist() does the same thing but is restricted to only errors returned by the os package.
					It is satisfied by [ErrNotExist] as well as some syscall errors.
			*/
			utils.PrintDelimiter()
			fmt.Println("No todos to clear")
			return
		}

		utils.PrintDelimiter()
		fmt.Println("Error clearing todos:", err)
		return
	}

	utils.PrintDelimiter()
	fmt.Println("Todos cleared successfully")
	if shouldPrint {
		utils.PrintTodos(&[]types.Todo{})
	}
}

func MarkAsCompleteHandler(ids []string, shouldPrint bool) {
	todoList, err := store.LoadTodos(config.Config.FILE_PATH)
	if err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error loading todos ->", err)
		return
	}

	hasSome := false

	todoList = utils.Map(todoList, func(todo types.Todo, _ int) types.Todo {
		if slices.Contains(ids, todo.ID) {
			hasSome = true
			todo.Completed = true
		}
		return todo
	})

	if !hasSome {
		utils.PrintDelimiter()
		fmt.Println("Todo(s) not found")
		if shouldPrint {
			utils.PrintTodos(&todoList)
		}
		return
	}

	if err := store.SaveTodos(config.Config.FILE_PATH, &todoList); err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error saving todos:", err)
		return
	}

	utils.PrintDelimiter()
	fmt.Printf("Todo(s) with ID(s) '%+v' marked as complete\n", ids)
	if shouldPrint {
		utils.PrintTodos(&todoList)
	}
}

func AddHandler(titles []string, shouldPrint bool) {
	todoList, err := store.LoadTodos(config.Config.FILE_PATH)
	if err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error loading todos ->", err)
		return
	}

	existingIds := utils.Map(todoList, func(todo types.Todo, _ int) string {
		return todo.ID
	})

	for _, title := range titles {
		id, err := utils.GenerateUniqueId(existingIds)

		if err != nil {
			utils.PrintDelimiter()
			fmt.Println("Error generating unique ID:", err)
			return
		}

		newTodo := types.Todo{
			ID:        id,
			Title:     title,
			Completed: false,
			CreatedAt: time.Now(),
		}

		todoList = append(todoList, newTodo)
	}

	if err := store.SaveTodos(config.Config.FILE_PATH, &todoList); err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error saving todos:", err)
		return
	}

	utils.PrintDelimiter()
	fmt.Println("Todo(s) added successfully")
	if shouldPrint {
		utils.PrintTodos(&todoList)
	}
}

func UpdateHandler(id string, title string, shouldPrint bool) {
	todoList, err := store.LoadTodos(config.Config.FILE_PATH)
	if err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error loading todos ->", err)
		return
	}

	hasSome := false

	todoList = utils.Map(todoList, func(todo types.Todo, _ int) types.Todo {
		if todo.ID == id {
			hasSome = true
			todo.Title = title
		}
		return todo
	})

	if !hasSome {
		utils.PrintDelimiter()
		fmt.Println("Todo not found")
		if shouldPrint {
			utils.PrintTodos(&todoList)
		}
		return
	}

	if err := store.SaveTodos(config.Config.FILE_PATH, &todoList); err != nil {
		utils.PrintDelimiter()
		fmt.Println("Error saving todos:", err)
		return
	}

	utils.PrintDelimiter()
	fmt.Printf("Todo with ID '%s' updated successfully\n", id)
	if shouldPrint {
		utils.PrintTodos(&todoList)
	}
	return
}
