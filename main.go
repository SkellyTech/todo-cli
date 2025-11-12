package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Task struct {
	ID        int
	Title     string
	Completed bool
}

const dataFile = "tasks.json"

// Load tasks from file
func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Save tasks to file
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo [add|list|done] [task]")
		return
	}

	command := os.Args[1]
	tasks, _ := loadTasks()

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo add [task name]")
			return
		}
		title := os.Args[2]
		task := Task{ID: len(tasks) + 1, Title: title}
		tasks = append(tasks, task)
		saveTasks(tasks)
		fmt.Println("✅ Added:", title)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return
		}
		fmt.Println("Your tasks:")
		for _, t := range tasks {
			status := " "
			if t.Completed {
				status = "✓"
			}
			fmt.Printf("%d. [%s] %s\n", t.ID, status, t.Title)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done [task number]")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task number")
			return
		}
		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Completed = true
				saveTasks(tasks)
				fmt.Println("🎉 Marked as done:", tasks[i].Title)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Task not found.")
		}

	default:
		fmt.Println("Unknown command:", command)
	}
}
