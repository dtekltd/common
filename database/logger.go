package database

import (
	"context"
	"log"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

type FilteredLogger struct {
	logger.Interface
}

func (l *FilteredLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()

	// Only log INSERT and UPDATE statements
	if strings.HasPrefix(strings.ToUpper(sql), "INSERT") ||
		strings.HasPrefix(strings.ToUpper(sql), "UPDATE") {
		log.Printf("SQL: %s", sql)
	}
}
