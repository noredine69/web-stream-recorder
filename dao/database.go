package dao

import (
	"sync"
	"web-stream-recorder/services/config"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB    *gorm.DB
	mutex *sync.Mutex
}

func NewDatabase(databaseConfig config.DatabaseConfig) (*Database, error) {
	db := Database{mutex: &sync.Mutex{}}

	var dbNameWithDSN string
	if databaseConfig.InMem {
		dbNameWithDSN = "file::memory:"
	} else {
		dbNameWithDSN = databaseConfig.DbPath + databaseConfig.DbName
	}

	gormDB, err := gorm.Open(sqlite.Open(dbNameWithDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while trying to open or create DB")
		return nil, err
	}

	db.DB = gormDB

	err = db.autoMigrate()
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (db *Database) autoMigrate() error {
	err := db.DB.AutoMigrate(&Record{})
	if err != nil {
		log.Error().Err(err).Msg("Error while trying to automigrate Record")
		return err
	}
	return nil
}

func (db *Database) LockAccess() {
	db.mutex.Lock()
}

func (db *Database) UnlockAccess() {
	db.mutex.Unlock()
}
