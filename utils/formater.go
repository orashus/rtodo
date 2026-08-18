package utils

import (
	"fmt"
	"time"

	"github.com/orashus/rtodo/types"
)

func PrintTodos(todos *[]types.Todo) {
	fmt.Printf("%-10s %-28s %-10s %-50s\n", "ID", "Created At", "Completed", "Title")
	fmt.Println(Delimiter)
	for _, todo := range *todos {
		fmt.Printf("%-10s %-28s %-10v %-50s\n", todo.ID, todo.CreatedAt.Format(time.RFC3339), todo.Completed, todo.Title)
	}
}

func PrintDelimiter() {
	fmt.Println(Delimiter)
}
