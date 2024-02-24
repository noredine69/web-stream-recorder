package config

type Config struct {
	Api       ApiConfig
	Recorder  RecorderConfig
	Database  DatabaseConfig
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

type DatabaseConfig struct {
	InMem  bool   `json:"inMem"`
	DbPath string `json:"dbPath"`
	DbName string `json:"dbName"`
}
