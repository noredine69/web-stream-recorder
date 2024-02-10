package recorder

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"web-stream-recorder/services/config"

	"github.com/rs/zerolog/log"
)

type RecorderAPIInterface interface {
	ScheduleRecording(ctx context.Context, delay int64, channel int64, duration int64) error
}

type RecorderAPI struct {
	config config.RecorderConfig
}

func New(config config.RecorderConfig) *RecorderAPI {
	return &RecorderAPI{
		config: config,
	}
}

func (recorder *RecorderAPI) ScheduleRecording(ctx context.Context, delay int64, channel int64, duration int64) error {

	go func() {
		log.Info().Msgf("Recording %d seconds for channel %d is scheduled in %d seconds from now", duration, channel, delay)
		<-time.After(time.Second * time.Duration(delay))
		err := recorder.startRecording(ctx, channel, duration)
		if err != nil {
			log.Error().Err(err).Msgf("Error while starting recording")
		}
	}()
	return nil
}

func (recorder *RecorderAPI) startRecording(ctxParent context.Context, channel int64, duration int64) error {
	requestUrl := fmt.Sprintf("%s/%s/%s/%d", recorder.config.Provider.Url, recorder.config.Provider.Login, recorder.config.Provider.Password, channel)
	//log.Error().Msgf("pplip %s", requestUrl)

	final_path := fmt.Sprintf("%s/record_%d_%s.mp4", recorder.config.Path, channel, time.Now().Format("2006_01_02_15_04_05"))
	temp_path := fmt.Sprintf("%s.tmp", final_path)

	req, _ := http.NewRequest("GET", requestUrl, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Error while closing response body")
		}
	}()

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
			if data != nil {
				//log.Error().Msgf("f.Write(data) %d", len(data))
				n, errWrite := f.Write(data)
				if errWrite != nil || n != len(data) {
					log.Error().Err(errWrite).Msgf("Error while writing to file")
				}
			}
		case <-time.After(5 * time.Second):
			log.Error().Msgf("time.After(5 * time.Second)")
			errRename := os.Rename(temp_path, final_path)
			if errRename != nil {
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
