package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

type BaseRepository struct {
	db *sql.DB
}

func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

func (br *BaseRepository) BuildWhereClause(filters map[string]string) (string, []interface{}, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	for key, value := range filters {
		if value != "" {
			condition := fmt.Sprintf("%s = $%d", key, argIndex)
			conditions = append(conditions, condition)
			args = append(args, value)
			argIndex++
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args, nil
}

func (br *BaseRepository) ExecuteTransaction(fn func(*sql.Tx) error) error {
	tx, err := br.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (br *BaseRepository) GetDB() *sql.DB {
	return br.db
}
