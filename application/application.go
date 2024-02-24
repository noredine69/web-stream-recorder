package application

import (
	"errors"
	"fmt"
	"os"
	"web-stream-recorder/api"
	"web-stream-recorder/dao"
	"web-stream-recorder/services/config"
	"web-stream-recorder/services/recorder"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	ErrCannotStartApp = errors.New("Application cannot start: Config are broken")
)

type ApplicationInterface interface {
	Run()
	Stop()
}

type Application struct {
	configuration config.Config
	api           *api.Api
	database      *dao.Database
	recorderApi   recorder.RecorderAPIInterface
	sigs          chan os.Signal
}

func New(configFilePath string, debugMode bool) (*Application, error) {
	app := &Application{
		configuration: readConfiguration(configFilePath),
	}

	if err := app.initServiceLayer(); err != nil {
		log.Error().Msgf("Error intializing service layers: %v", err)
		return nil, ErrCannotStartApp
	}
	return app, nil
}

func readConfiguration(configFilePath string) config.Config {
	var config config.Config

	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("unable to read config file, %v", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Sprintf("unable to decode into struct, %v", err))
	}
	return config
}

func (app *Application) Run() {
}

func (app *Application) Stop() {
}
func (app *Application) initServiceLayer() error {
	var err error
	app.database, err = dao.NewDatabase(app.configuration.Database)
	if err != nil {
		return err
	}
	app.recorderApi = recorder.New(app.configuration.Recorder, app.database)
	app.api = api.New(app.configuration, app.recorderApi)
	app.api.Run()
	return nil
}
