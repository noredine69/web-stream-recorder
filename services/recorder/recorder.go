package recorder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"web-stream-recorder/dao"
	"web-stream-recorder/services/config"

	"github.com/rs/zerolog/log"
)

var (
	errorRecordAlreayExist = errors.New("at leat one record already exists for this channel and this date")
	errorBodyNil           = errors.New("response body is nil")
)

type RecorderAPIInterface interface {
	ScheduleRecording(ctx context.Context, delay int64, channel int64, duration int64) error
	GetAllRecords() []dao.Record
}

type RecorderAPI struct {
	config   config.RecorderConfig
	database *dao.Database
}

func New(config config.RecorderConfig, database *dao.Database) *RecorderAPI {
	return &RecorderAPI{
		config:   config,
		database: database,
	}
}
func (recorder *RecorderAPI) GetAllRecords() []dao.Record {
	records, err := recorder.database.GetAllRecords()
	if err != nil {
		log.Error().Err(err).Msg("Error while getting all records")
		return []dao.Record{}
	}
	return records
}

func (recorder *RecorderAPI) ScheduleRecording(ctx context.Context, delay int64, channel int64, duration int64) error {
	now := time.Now()
	expectedStart := now.Add(time.Second * time.Duration(delay))
	expectedEnd := expectedStart.Add(time.Second * time.Duration(duration))
	records, err := recorder.database.FindRecord(channel, expectedStart, expectedEnd)
	if err != nil {
		log.Error().Err(err).Msg("Error while finding records")
		return err
	}
	if len(records) > 0 {
		log.Info().Msgf("Found %d records for channel %d", len(records), channel)
		return errorRecordAlreayExist
	}

	go func() {
		log.Info().Msgf("Recording %d seconds for channel %d is scheduled in %d seconds from now", duration, channel, delay)

		//.Format("2006_01_02_15_04_05")
		recordId, err := recorder.database.SaveRecord(&dao.Record{
			Name:          "name",
			Channel:       uint64(channel),
			Date:          &now,
			ExpectedStart: &expectedStart,
			ExpectedEnd:   &expectedEnd,
			Description:   "desc",
			//0:scheduled, 1:in progress, 2:completed, 3:canceled, 4:failed
			Status: 0,
		})
		if err != nil {
			log.Error().Err(err).Msgf("Error while saving record for channel %d", channel)
		}
		<-time.After(time.Second * time.Duration(delay))
		err = recorder.startRecording(ctx, recordId, channel, duration)
		if err != nil {
			log.Error().Err(err).Msgf("Error while starting recording")
		}
	}()
	return nil
}

func (recorder *RecorderAPI) startRecording(ctxParent context.Context, recordId int, channel int64, duration int64) error {
	requestUrl := fmt.Sprintf("%s/%s/%s/%d", recorder.config.Provider.Url, recorder.config.Provider.Login, recorder.config.Provider.Password, channel)
	//log.Error().Msgf("pplip %s", requestUrl)

	err := recorder.database.StartsRecord(recordId)
	if err != nil {
		log.Error().Err(err).Msg("Error while starting recording")
		return err
	}

	final_path := fmt.Sprintf("%s/record_%d_%s.mp4", recorder.config.Path, channel, time.Now().Format("2006_01_02_15_04_05"))
	temp_path := fmt.Sprintf("%s.tmp", final_path)

	req, _ := http.NewRequest("GET", requestUrl, nil)
	resp, errReq := http.DefaultClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			err := resp.Body.Close()
			if err != nil {
				log.Error().Err(err).Msgf("Error while closing response body")
			}
		}
	}()
	if errReq != nil {
		log.Error().Err(errReq).Msgf("Error while making request to %s", requestUrl)
		return errReq
	}

	if resp == nil || resp.Body == nil {
		err := recorder.database.RecordFails(recordId)
		if err != nil {
			log.Error().Err(err).Msgf("Error while marking record as failed")
		}
		return errorBodyNil
	}
	f, _ := os.OpenFile(temp_path, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	ctx, cancel := context.WithDeadline(ctxParent, time.Now().Add(time.Second*time.Duration(duration)))
	defer cancel()
	ch := make(chan []byte, 1)

	go func() {
		for {
			select {
			default:
				ch <- recorder.GetLastBlock(resp.Body)
			case <-ctx.Done():
				log.Error().Msgf("Canceled by timeout")
				return
			}
		}
	}()

	for {
		select {
		case data := <-ch:
			if data == nil {
				recorder.database.RecordFails(recordId)
			} else {
				//log.Error().Msgf("f.Write(data) %d", len(data))
				n, errWrite := f.Write(data)
				if errWrite != nil || n != len(data) {
					log.Error().Err(errWrite).Msgf("Error while writing to file")
				}
			}
		case <-time.After(5 * time.Second):
			log.Error().Msgf("time.After(5 * time.Second)")
			errRename := os.Rename(temp_path, final_path)
			if errRename == nil {
				recorder.database.RecordSucceed(recordId)
			} else {
				recorder.database.RecordFails(recordId)
				log.Error().Err(errRename).Msgf("Error while renaming file from %s, to %s", temp_path, final_path)
			}
			return nil
		}
	}
}

func (recorder *RecorderAPI) GetLastBlock(source io.ReadCloser) []byte {
	buf := make([]byte, 32*1024)
	n, err := source.Read(buf)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		log.Error().Err(err).Msgf("Error while downloading: %v", err)
	}
	if n > 0 {
		//log.Error().Msgf("buf[:n] %d", n)
		return buf[:n]
	}
	return nil
}
