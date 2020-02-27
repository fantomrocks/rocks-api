package common

import (
	"github.com/spf13/viper"
	"log"
)

// Structure describes configuration options for Crystal API server.
type Config struct {
	// server specific options
	AppName  string
	BindAddr string
	Cors     []string

	// logger specific options
	LogLevel  string
	LogFormat string

	// database configuration details
	DbDriver             string
	DbHost               string
	DbPort               string
	DbName               string
	DbUser               string
	DbPassword           string
	DbMaxOpenConnections int

	// RPC connection to the related block chain node
	RpcUrl string
}

// Define Context key for configuration access.
type ConfigContextKey struct{}

// default configuration options
var defaults = map[string]interface{}{
	"server.name": "FantomRocksApi",
	"server.cors": []string{"*"},

	"logger.level":  "INFO",
	"logger.format": "%{color}%{time:2019-01-01 15:04:05} [%{level:.6s}] %{shortfunc}:%{color:reset} %{message}",

	"db.driver":    "postgres",
	"db.host":      "localhost",
	"db.port":      "8084",
	"db.name":      "fantom",
	"db.user":      "default-user",
	"db.password":  "default-password",
	"db.pool_size": "10",

	"rpc.url": "~/.lachesis/data/lachesis.ipc",
}

// Function provides loaded configuration for Crystal API server.
func LoadConfig() *Config {
	cfg := viper.New()

	// what is the expected name of the common file
	cfg.SetConfigName("config")

	// where to look for common files
	cfg.AddConfigPath("$HOME/.fantomrocks")
	cfg.AddConfigPath(".")

	// set default values
	applyDefaults(cfg)

	// try to read the file
	err := cfg.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error reading configuration file: %s\n", err)
	}

	// build and return the common struct
	return &Config{
		// basics
		AppName:  cfg.GetString("server.name"),
		BindAddr: cfg.GetString("server.listen"),
		Cors:     cfg.GetStringSlice("server.cors"),

		// logger
		LogLevel:  cfg.GetString("logger.level"),
		LogFormat: cfg.GetString("logger.format"),

		// db connection details
		DbDriver:             cfg.GetString("db.driver"),
		DbHost:               cfg.GetString("db.host"),
		DbPort:               cfg.GetString("db.port"),
		DbName:               cfg.GetString("db.name"),
		DbUser:               cfg.GetString("db.user"),
		DbPassword:           cfg.GetString("db.password"),
		DbMaxOpenConnections: cfg.GetInt("db.pool_size"),

		// RPC related
		RpcUrl: cfg.GetString("rpc.url"),
	}
}

// load default/predefined values to the configuration manager
func applyDefaults(cfg *viper.Viper) {
	for key, value := range defaults {
		cfg.SetDefault(key, value)
	}
}
