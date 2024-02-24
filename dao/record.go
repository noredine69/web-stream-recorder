package dao

import (
	"time"

	"github.com/rs/zerolog/log"
)

func (db *Database) SaveRecord(record *Record) (int, error) {
	db.LockAccess()
	defer db.UnlockAccess()

	res := db.DB.Save(record)
	return record.ID, res.Error
}

func (db *Database) FindRecord(channel int64, expectedStart time.Time, expectedEnd time.Time) ([]Record, error) {
	var recs []Record
	db.LockAccess()
	defer db.UnlockAccess()

	if err := db.DB.Where("channel = ? AND status in (0, 1) AND ((expected_start <= ? AND expected_end >= ?) OR (? <= expected_start AND ? >= expected_start))",
		channel, expectedStart, expectedStart, expectedStart, expectedEnd).Find(&recs).Error; err != nil {
		return []Record{}, err
	}

	return recs, nil
}

func (db *Database) GetRecord(key int) (Record, error) {
	var rec Record
	db.LockAccess()
	defer db.UnlockAccess()

	if err := db.DB.Where("ID = ?", key).First(&rec).Error; err != nil {
		return Record{}, err
	}

	return rec, nil
}

func (db *Database) GetAllRecords() ([]Record, error) {
	var recs []Record
	db.LockAccess()
	defer db.UnlockAccess()

	result := db.DB.Find(&recs)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error while requesting records")
		return nil, result.Error
	}
	return recs, nil
}

func (db *Database) StartsRecord(id int) error {
	db.LockAccess()
	defer db.UnlockAccess()

	res := db.DB.Model(&Record{}).Where("ID = ?", id).Updates(map[string]interface{}{"start": time.Now(), "status": 1})
	return res.Error
}

func (db *Database) RecordSucceed(id int) error {
	db.LockAccess()
	defer db.UnlockAccess()

	res := db.DB.Model(&Record{}).Where("ID = ?", id).Updates(map[string]interface{}{"end": time.Now(), "status": 2})
	return res.Error
}

func (db *Database) RecordFails(id int) error {
	db.LockAccess()
	defer db.UnlockAccess()

	res := db.DB.Model(&Record{}).Where("ID = ?", id).Updates(map[string]interface{}{"end": time.Now(), "status": 4})
	return res.Error
}
