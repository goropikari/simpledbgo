package infra

import (
	"errors"
	"log"

	"github.com/goropikari/simpledb_go/lib/directio"
)

// ErrInvalidBlockSize is an error type that means given block size is invalid.
var ErrInvalidBlockSize = errors.New("invalid block size")

// Config is config of database server.
type Config struct {
	DBPath      string
	LogFileName string
	BlockSize   int
	IsDirectIO  bool
}

// NewConfig is a constructor of Config.
func NewConfig(dbPath string, blockSize int, logFileName string) Config {
	return Config{
		DBPath:      dbPath,
		BlockSize:   blockSize,
		LogFileName: logFileName,
		IsDirectIO:  false,
	}
}

// NewDirectIOConfig is a constructor of config for DirectIO.
func NewDirectIOConfig(dbPath string, blockSize int, logFileName string) Config {
	if blockSize%directio.BlockSize != 0 {
		log.Fatal(ErrInvalidBlockSize)
	}

	return Config{
		DBPath:      dbPath,
		BlockSize:   blockSize,
		LogFileName: logFileName,
		IsDirectIO:  true,
	}
}

// SetDefaults sets default configuration value.
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
