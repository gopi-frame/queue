package database

import (
	"strings"

	"github.com/gopi-frame/contract/foundation"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/database"
	"gorm.io/gorm"
)

type FailedJobProviderConnector struct {
	app foundation.Application
}

func (f *FailedJobProviderConnector) Connect(config map[string]any) queue.FailedJobProvider {
	databaseManager := f.app.Get("db").(*database.DatabaseManager)
	provider := new(DatabaseFailedJobProvider)
	if connection, ok := config["connection"]; ok {
		switch connection := connection.(type) {
		case string:
			if connection := strings.TrimSpace(connection); connection != "" {
				provider.db = databaseManager.Connection(connection)
			} else {
				provider.db = databaseManager.Connection()
			}
		case *gorm.DB:
			provider.db = connection
		case gorm.Dialector:
			db, err := gorm.Open(connection)
			if err != nil {
				panic(err)
			}
			provider.db = db
		default:
			provider.db = databaseManager.Connection()
		}
	} else {
		provider.db = databaseManager.Connection()
	}
	if table, ok := config["table"]; ok {
		switch table := table.(type) {
		case string:
			provider.table = table
		default:
			provider.table = DefaultFailedJobTable
		}
	}
	return provider
}
