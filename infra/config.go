package infra

type Config struct {
	DBPath      string
	LogFileName string
	BlockSize   int
	IsDirectIO  bool
}

func NewConfig(dbPath string, blockSize int, logFileName string) Config {
	return Config{
		DBPath:      dbPath,
		BlockSize:   blockSize,
		LogFileName: logFileName,
	}
}

func (config *Config) SetDefaults() {
	if config.DBPath == "" {
		config.DBPath = "simpledb"
	}

	if config.BlockSize == 0 {
		config.BlockSize = 4096
	}

	if config.LogFileName == "" {
		config.LogFileName = "logfile"
	}
}
