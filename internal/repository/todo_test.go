package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log/slog"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/todo?parseTime=true")
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	return db
}

func TestCreateTodo(t *testing.T) {
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodoRepository(dbConn, logger)
	ctx := context.Background()

	title := "TestCreate"
	description := "Create Desc"
	due := time.Now().Add(24 * time.Hour)
	todo, err := repo.Create(ctx, title, description, &due)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}
	if todo.Title != title {
		t.Errorf("expected title %q, got %q", title, todo.Title)
	}
}

func TestReadTodo(t *testing.T) {
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodoRepository(dbConn, logger)
	ctx := context.Background()

	title := "TestRead"
	description := "Read Desc"
	due := time.Now().Add(24 * time.Hour)
	todo, err := repo.Create(ctx, title, description, &due)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	todos, err := repo.List(ctx, "")
	if err != nil {
		t.Fatalf("failed to list todos: %v", err)
	}
	found := false
	for _, td := range todos {
		if td.ID == todo.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created todo not found in list")
	}
}

func TestUpdateTodo(t *testing.T) {
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodoRepository(dbConn, logger)
	ctx := context.Background()

	title := "TestUpdate"
	description := "Update Desc"
	due := time.Now().Add(24 * time.Hour)
	todo, err := repo.Create(ctx, title, description, &due)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	newTitle := "Updated"
	newStatus := "completed"
	updated, err := repo.Update(ctx, todo.ID, newTitle, description, &due, newStatus)
	if err != nil {
		t.Fatalf("failed to update todo: %v", err)
	}
	if updated.Title != newTitle || updated.Status != newStatus {
		t.Errorf("update did not persist changes")
	}
}

func TestDeleteTodo(t *testing.T) {
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodoRepository(dbConn, logger)
	ctx := context.Background()

	title := "TestDelete"
	description := "Delete Desc"
	due := time.Now().Add(24 * time.Hour)
	todo, err := repo.Create(ctx, title, description, &due)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	err = repo.Delete(ctx, todo.ID)
	if err != nil {
		t.Fatalf("failed to delete todo: %v", err)
	}

	todos, err := repo.List(ctx, "")
	if err != nil {
		t.Fatalf("failed to list todos after delete: %v", err)
	}
	for _, td := range todos {
		if td.ID == todo.ID {
			t.Errorf("deleted todo still found in list")
		}
	}
}
