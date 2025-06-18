package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"todo01/gen/db"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TodoRepository struct {
	db *sql.DB
	logger *slog.Logger
}

func NewTodoRepository(db *sql.DB, logger *slog.Logger) *TodoRepository {
	return &TodoRepository{db: db, logger: logger}
}

func (r *TodoRepository) Create(ctx context.Context, title, description string, dueDate *time.Time) (*db.Todo, error) {
	todo := &db.Todo{
		Title:       title,
		Description: null.StringFrom(description),
		Status:      "pending",
	}
	if dueDate != nil {
		todo.DueDate = null.TimeFrom(*dueDate)
	}

	err := todo.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		r.logger.Error("DB error on Create", "title", title, "error", err)
		return nil, err
	}

	return todo, nil
}

func (r *TodoRepository) Update(ctx context.Context, id int64, title, description string, dueDate *time.Time, status string) (*db.Todo, error) {
	todo, err := db.Todos(
		db.TodoWhere.ID.EQ(id),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("DB error on Update (fetch)", "id", id, "error", err)
		return nil, err
	}

	todo.Title = title
	todo.Description = null.StringFrom(description)
	todo.Status = status
	if dueDate != nil {
		todo.DueDate = null.TimeFrom(*dueDate)
	}

	_, err = todo.Update(ctx, r.db, boil.Infer())
	if err != nil {
		r.logger.Error("DB error on Update (save)", "id", id, "error", err)
		return nil, err
	}

	return todo, nil
}

func (r *TodoRepository) Delete(ctx context.Context, id int64) error {
	todo, err := db.Todos(
		db.TodoWhere.ID.EQ(id),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("DB error on Delete (fetch)", "id", id, "error", err)
		return err
	}

	todo.DeletedAt = null.TimeFrom(time.Now())
	_, err = todo.Update(ctx, r.db, boil.Infer())
	if err != nil {
		r.logger.Error("DB error on Delete (save)", "id", id, "error", err)
	}
	return err
}

func (r *TodoRepository) List(ctx context.Context, status string) ([]*db.Todo, error) {
	mods := []qm.QueryMod{
		db.TodoWhere.DeletedAt.IsNull(),
	}
	if status != "" {
		mods = append(mods, db.TodoWhere.Status.EQ(status))
	}

	todos, err := db.Todos(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("DB error on List", "status", status, "error", err)
		return nil, err
	}
	return todos, nil
} 