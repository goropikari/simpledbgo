package infra

import (
	"errors"
	"log"

	"github.com/goropikari/simpledb_go/lib/directio"
)

var ErrInvalidBlockSize = errors.New("invalid block size")

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
		IsDirectIO:  false,
	}
}

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
