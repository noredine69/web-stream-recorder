package config

type Config struct {
	Api       ApiConfig
	Recorder  RecorderConfig
	DebugMode bool
}

type RecorderConfig struct {
	BuffSize int    `json:"buffSize"`
	Path     string `json:"path"`
	Provider ProviderConfig
}
type ProviderConfig struct {
	Url      string `json:"url"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ApiConfig struct {
	Port      int    `json:"port"`
	SecretKey string `json:"secretKey"`
}
