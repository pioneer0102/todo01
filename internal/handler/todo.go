package handler

import (
	"context"
	"time"
	"log/slog"

	"todo01/gen/go/proto/todo/v1"
	"todo01/gen/db"
	"todo01/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/bufbuild/connect-go"
)

type TodoHandler struct {
	repo *repository.TodoRepository
    logger *slog.Logger
}

func NewTodoHandler(repo *repository.TodoRepository, logger *slog.Logger) *TodoHandler {
	return &TodoHandler{repo: repo, logger: logger}
}

func (h *TodoHandler) CreateTodo(ctx context.Context, req *connect.Request[todov1.CreateTodoRequest]) (*connect.Response[todov1.CreateTodoResponse], error) {
	h.logger.Info("CreateTodo called", "title", req.Msg.Title)
	var dueDate *time.Time
	if req.Msg.DueDate != nil {
		t := req.Msg.DueDate.AsTime()
		dueDate = &t
	}

	todo, err := h.repo.Create(ctx, req.Msg.Title, req.Msg.Description, dueDate)
	if err != nil {
		h.logger.Error("Failed to create todo", "error", err)
		return nil, err
	}

	h.logger.Info("Todo created", "id", todo.ID)
	return connect.NewResponse(&todov1.CreateTodoResponse{
		Todo: convertToProtoTodo(todo),
	}), nil
}

func (h *TodoHandler) UpdateTodo(ctx context.Context, req *connect.Request[todov1.UpdateTodoRequest]) (*connect.Response[todov1.UpdateTodoResponse], error) {
	h.logger.Info("UpdateTodo called", "id", req.Msg.Id)
	var dueDate *time.Time
	if req.Msg.DueDate != nil {
		t := req.Msg.DueDate.AsTime()
		dueDate = &t
	}

	todo, err := h.repo.Update(ctx, req.Msg.Id, req.Msg.Title, req.Msg.Description, dueDate, toDBStatus(req.Msg.Status))
	if err != nil {
		h.logger.Error("Failed to update todo", "id", req.Msg.Id, "error", err)
		return nil, err
	}

	h.logger.Info("Todo updated", "id", todo.ID)
	return connect.NewResponse(&todov1.UpdateTodoResponse{
		Todo: convertToProtoTodo(todo),
	}), nil
}

func (h *TodoHandler) DeleteTodo(ctx context.Context, req *connect.Request[todov1.DeleteTodoRequest]) (*connect.Response[todov1.DeleteTodoResponse], error) {
	h.logger.Info("DeleteTodo called", "id", req.Msg.Id)
	err := h.repo.Delete(ctx, req.Msg.Id)
	if err != nil {
		h.logger.Error("Failed to delete todo", "id", req.Msg.Id, "error", err)
		return nil, err
	}

	h.logger.Info("Todo deleted", "id", req.Msg.Id)
	return connect.NewResponse(&todov1.DeleteTodoResponse{}), nil
}

func (h *TodoHandler) ListTodos(ctx context.Context, req *connect.Request[todov1.ListTodosRequest]) (*connect.Response[todov1.ListTodosResponse], error) {
	h.logger.Info("ListTodos called")
	var status string
	if req.Msg.Status != todov1.TodoStatus_TODO_STATUS_UNSPECIFIED {
		status = req.Msg.Status.String()
	}

	todos, err := h.repo.List(ctx, status)
	if err != nil {
		h.logger.Error("Failed to list todos", "error", err)
		return nil, err
	}

	protoTodos := make([]*todov1.Todo, len(todos))
	for i, todo := range todos {
		protoTodos[i] = convertToProtoTodo(todo)
	}

	h.logger.Info("ListTodos completed", "count", len(protoTodos))
	return connect.NewResponse(&todov1.ListTodosResponse{
		Todos: protoTodos,
	}), nil
}

func convertToProtoTodo(todo *db.Todo) *todov1.Todo {
	protoTodo := &todov1.Todo{
		Id:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description.String,
		Status:      toProtoStatus(todo.Status),
		CreatedAt:   timestamppb.New(todo.CreatedAt),
		UpdatedAt:   timestamppb.New(todo.UpdatedAt),
	}

	if todo.DueDate.Valid {
		protoTodo.DueDate = timestamppb.New(todo.DueDate.Time)
	}

	return protoTodo
}

func toProtoStatus(status string) todov1.TodoStatus {
	switch status {
	case "pending":
		return todov1.TodoStatus_TODO_STATUS_PENDING
	case "completed":
		return todov1.TodoStatus_TODO_STATUS_COMPLETED
	default:
		return todov1.TodoStatus_TODO_STATUS_UNSPECIFIED
	}
}

func toDBStatus(status todov1.TodoStatus) string {
	switch status {
	case todov1.TodoStatus_TODO_STATUS_PENDING:
		return "pending"
	case todov1.TodoStatus_TODO_STATUS_COMPLETED:
		return "completed"
	default:
		return "pending"
	}
}