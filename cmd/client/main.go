package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"todo01/gen/go/proto/todo/v1"
	"todo01/gen/go/proto/todo/v1/todov1connect"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/bufbuild/connect-go"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	action := flag.String("action", "", "Action to perform (create, update, delete, list)")
	title := flag.String("title", "", "Todo title")
	description := flag.String("description", "", "Todo description")
	dueDate := flag.String("due-date", "", "Due date (YYYY-MM-DD)")
	status := flag.String("status", "", "Todo status (pending, completed)")
	id := flag.Int64("id", 0, "Todo ID")
	flag.Parse()

	client := todov1connect.NewTodoServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	ctx := context.Background()

	switch *action {
	case "create":
		createTodo(ctx, client, *title, *description, *dueDate)
	case "update":
		updateTodo(ctx, client, *id, *title, *description, *dueDate, *status)
	case "delete":
		deleteTodo(ctx, client, *id)
	case "list":
		listTodos(ctx, client, *status)
	default:
		fmt.Println("Please specify an action: create, update, delete, or list")
		os.Exit(1)
	}
}

func createTodo(ctx context.Context, client todov1connect.TodoServiceClient, title, description, dueDate string) {
	var dueDatePb *timestamppb.Timestamp
	if dueDate != "" {
		t, err := time.Parse("2006-01-02", dueDate)
		if err != nil {
			slog.Error("Invalid due date format", "error", err)
			os.Exit(1)
		}
		dueDatePb = timestamppb.New(t)
	}

	resp, err := client.CreateTodo(ctx, connect.NewRequest(&todov1.CreateTodoRequest{
		Title:       title,
		Description: description,
		DueDate:     dueDatePb,
	}))
	if err != nil {
		slog.Error("Failed to create todo", "error", err)
		os.Exit(1)
	}

	printTodo(resp.Msg.Todo)
}

func updateTodo(ctx context.Context, client todov1connect.TodoServiceClient, id int64, title, description, dueDate, status string) {
	var dueDatePb *timestamppb.Timestamp
	if dueDate != "" {
		t, err := time.Parse("2006-01-02", dueDate)
		if err != nil {
			slog.Error("Invalid due date format", "error", err)
			os.Exit(1)
		}
		dueDatePb = timestamppb.New(t)
	}

	var statusEnum todov1.TodoStatus
	switch status {
	case "pending":
		statusEnum = todov1.TodoStatus_TODO_STATUS_PENDING
	case "completed":
		statusEnum = todov1.TodoStatus_TODO_STATUS_COMPLETED
	}

	resp, err := client.UpdateTodo(ctx, connect.NewRequest(&todov1.UpdateTodoRequest{
		Id:          id,
		Title:       title,
		Description: description,
		DueDate:     dueDatePb,
		Status:      statusEnum,
	}))
	if err != nil {
		slog.Error("Failed to update todo", "error", err)
		os.Exit(1)
	}

	printTodo(resp.Msg.Todo)
}

func deleteTodo(ctx context.Context, client todov1connect.TodoServiceClient, id int64) {
	_, err := client.DeleteTodo(ctx, connect.NewRequest(&todov1.DeleteTodoRequest{
		Id: id,
	}))
	if err != nil {
		slog.Error("Failed to delete todo", "error", err)
		os.Exit(1)
	}

	fmt.Println("Todo deleted successfully")
}

func listTodos(ctx context.Context, client todov1connect.TodoServiceClient, status string) {
	var statusEnum todov1.TodoStatus
	switch status {
	case "pending":
		statusEnum = todov1.TodoStatus_TODO_STATUS_PENDING
	case "completed":
		statusEnum = todov1.TodoStatus_TODO_STATUS_COMPLETED
	}

	resp, err := client.ListTodos(ctx, connect.NewRequest(&todov1.ListTodosRequest{
		Status: statusEnum,
	}))
	if err != nil {
		slog.Error("Failed to list todos", "error", err)
		os.Exit(1)
	}

	for _, todo := range resp.Msg.Todos {
		printTodo(todo)
		fmt.Println("---")
	}
}

func printTodo(todo *todov1.Todo) {
	fmt.Printf("ID: %d\n", todo.Id)
	fmt.Printf("Title: %s\n", todo.Title)
	fmt.Printf("Description: %s\n", todo.Description)
	if todo.DueDate != nil {
		fmt.Printf("Due Date: %s\n", todo.DueDate.AsTime().Format("2006-01-02"))
	}
	fmt.Printf("Status: %s\n", todo.Status)
	fmt.Printf("Created: %s\n", todo.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", todo.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))
}