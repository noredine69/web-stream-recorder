package dao

/*

func TestLogsState(t *testing.T) {
	db, err := NewDatabase(models.DatabaseConfig{InMem: true})
	assert.NoError(t, err)

	logs, err := GetAllIsolatedLogs(db)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(logs))

	log, err := GetIsolatedLog(1, db)
	assert.Equal(t, 0, log.ID)
	assert.NotNil(t, err)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())

	logState := logapi.LogMetaData{ID: 1, Name: "name", Description: "desc", UserID: 12}
	err = SaveIsolatedLog(&logState, db)
	assert.NoError(t, err)

	logs, err = GetAllIsolatedLogs(db)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(logs))
	logRead := logs[0]
	assert.NotNil(t, logRead)
	assert.Equal(t, logState, logRead)

	count, err := CountWithIds([]int{1}, db)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	err = SetExtracted([]int{1}, db)
	assert.NoError(t, err)

	log, err = GetIsolatedLog(1, db)
	assert.Nil(t, err)
	assert.NotNil(t, log)
	assert.True(t, log.IsExtract)

	log.IsExtract = false
	logsToSave := []logapi.LogMetaData{log}
	err = SaveIsolatedLogs(&logsToSave, db)
	assert.Nil(t, err)

	log, err = GetIsolatedLog(1, db)
	assert.Nil(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, logState, log)

	err = DeleteLog(&log, db)
	assert.Nil(t, err)

	logs, err = GetAllIsolatedLogs(db)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(logs))

	err = DeleteAllLogs(db)
	assert.NoError(t, err)

	logs, err = GetAllIsolatedLogs(db)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(logs))

	db.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&logapi.LogMetaData{})
}
*/
