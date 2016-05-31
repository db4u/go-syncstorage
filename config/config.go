package config

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vrischmann/envconfig"
)

type LogConfig struct {

	// logging level, panic, fatal, error, warn, info, debug
	Level string `envconfig:"default=info"`

	// use mozlog format
	Mozlog bool `envconfig:"default=false"`
}

var Config struct {
	Log      *LogConfig
	Hostname string `envconfig:"optional"`
	Host     string `envconfig:"default=0.0.0.0"`
	Port     int
	Secrets  []string
	DataDir  string
	TTL      int `envconfig:"default=300"` // seconds to wait before closing a user's api handler

	MaxOpenFiles int `envconfig:"default=64"`
}

// so we can use config.Port and not config.Config.Port
var (
	Hostname string
	Log      *LogConfig
	Host     string
	Port     int
	DataDir  string
	Secrets  []string
	TTL      time.Duration
)

func init() {
	if err := envconfig.Init(&Config); err != nil {
		log.Fatalf("Config Error: %s\n", err)
	}

	if Config.Port < 1 || Config.Port > 65535 {
		log.Fatal("Config.Error: PORT invalid")
	}

	if Config.DataDir != ":memory:" {
		if _, err := os.Stat(Config.DataDir); os.IsNotExist(err) {
			log.Fatal("Config Error: DATA_DIR does not exist")
		}

		stat, err := os.Stat(Config.DataDir)
		if os.IsNotExist(err) {
			log.Fatal("Config Error: DATA_DIR does not exist")
		}
		if !stat.IsDir() {
			log.Fatal("Config Error: DATA_DIR is not a directory")
		}

		Config.DataDir = filepath.Clean(Config.DataDir)
		testfile := Config.DataDir + string(os.PathSeparator) + "test.writable"
		f, err := os.Create(testfile)
		if err != nil {
			log.Fatal("Config Error: DATA_DIR is not writable")
		} else {
			f.Close()
			os.Remove(testfile)
		}
	}

	switch Config.Log.Level {
	case "panic", "fatal", "error", "warn", "info", "debug":
	default:
		log.Fatalf("Config Error: LOG_LEVEL must be [panic, fatal, error, warn, info, debug]")
	}

	if Config.Hostname == "" {
		Config.Hostname, _ = os.Hostname()
	}

	if Config.TTL <= 0 {
		log.Fatal("TTL must be > 0")
	}

	Hostname = Config.Hostname
	Log = Config.Log
	Host = Config.Host
	Port = Config.Port
	Secrets = Config.Secrets
	DataDir = Config.DataDir
	TTL = time.Duration(Config.TTL) * time.Second
}
