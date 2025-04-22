package repository

import (
	"database/sql"

	"go.uber.org/zap"
)

func rollbackTx(logger *zap.Logger, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		logger.Error("Rollback failed", zap.Error(err))
	}
}
