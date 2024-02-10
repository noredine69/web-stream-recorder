package main

import (
	"flag"
	"web-stream-recorder/application"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Info().Msg("Starting the Ethplorer proxy...")

	configFilePath := flag.String("config", "resources/config/config.json", "path for config file")
	debugMode := flag.Bool("debug", false, "(optional) start in debug mode")
	flag.Parse()

	if *debugMode {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	app, err := application.New(*configFilePath, *debugMode)

	if err != nil {
		log.Fatal().Msgf("Application cannot be started correctly, turning off...")
	}
	app.Run()
}
