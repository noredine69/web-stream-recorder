package application

import "github.com/rs/zerolog/log"

type ApplicationMock struct {
	RunFunc  func()
	StopFunc func()
}

func NewMock() *ApplicationMock {
	return &ApplicationMock{}
}

func (fs *ApplicationMock) Run() {
	if fs.RunFunc != nil {
		fs.RunFunc()
		return
	}
	log.Warn().Msgf("No mocked provided for Run function")
}

func (fs *ApplicationMock) Stop() {
	if fs.StopFunc != nil {
		fs.StopFunc()
		return
	}
	log.Warn().Msgf("No mocked provided for Stop function")
}
