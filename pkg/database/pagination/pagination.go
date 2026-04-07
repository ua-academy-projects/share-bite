package pagination

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Result[T any] struct {
	Items []T
	Total int
}

type Params struct {
	Table   string
	Columns string
	Where   string
	OrderBy string
	Args    []any
	Offset  int
	Limit   int
}

func List[T any](
	ctx context.Context,
	db database.DB,
	queryName string,
	p Params,
	scanner func(pgx.Rows) (T, error),
) (Result[T], error) {
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", p.Table, p.Where)
	countQ := database.Query{
		Name: queryName + ".count",
		Sql:  countSQL,
	}

	var total int
	if err := db.QueryRowContext(ctx, countQ, p.Args...).Scan(&total); err != nil {
		return Result[T]{}, fmt.Errorf("count: %w", err)
	}

	nArgs := len(p.Args)
	itemsSQL := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d",
		p.Columns, p.Table, p.Where, p.OrderBy, nArgs+1, nArgs+2,
	)
	itemsQ := database.Query{
		Name: queryName + ".list",
		Sql:  itemsSQL,
	}

	args := append(p.Args, p.Limit, p.Offset)
	rows, err := db.QueryContext(ctx, itemsQ, args...)
	if err != nil {
		return Result[T]{}, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var items []T
	for rows.Next() {
		item, err := scanner(rows)
		if err != nil {
			return Result[T]{}, fmt.Errorf("scan: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return Result[T]{}, fmt.Errorf("rows: %w", err)
	}

	return Result[T]{
		Items: items,
		Total: total,
	}, nil
}
