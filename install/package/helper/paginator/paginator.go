package paginator

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"math"
	"sync"
)

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {

	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int64(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func CursorPaginate(tableName string, db *gorm.DB, cursor *Cursor, ctx context.Context) error {
	var wg sync.WaitGroup

	if cursor != nil {
		if cursor.Start != 0 {
			db.Where(fmt.Sprintf("%s.id < ?", tableName), cursor.Start)
		}
		wg.Add(1)
		//go countData(db, &cursor.TotalRows, &wg, ctx)
		wg.Wait()
		db.Limit(cursor.GetLimit()).Order(cursor.
			GetSort())
		cursor.TotalPages = int(math.Ceil(float64(cursor.TotalRows) / float64(cursor.Limit)))
	}

	return db.Scan(ctx).Error
}
