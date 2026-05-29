package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/darkness/green_api/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Env string

var (
	Local Env = "local"
	Dev   Env = "dev"
	Prod  Env = "prod"
)

type Config struct {
	*jsonConfig
	*envConfig
}

type envConfig struct {
	Env Env `env:"ENV" env-required:"true"`
}

type jsonConfig struct {
	Server  ServerConfig  `json:"server" mapstructure:"server" validate:"required"`
	Handler HandlerConfig `json:"handler" mapstructure:"handler" validate:"required"`
}

type ServerConfig struct {
	Port           int           `json:"port" mapstructure:"port" validate:"required"`
	ReadTimeout    time.Duration `json:"read_timeout" mapstructure:"read_timeout" validate:"required"`
	WriteTimeout   time.Duration `json:"write_timeout" mapstructure:"write_timeout" validate:"required"`
	MaxHeaderBytes int           `json:"max_header_bytes" mapstructure:"max_header_bytes" validate:"required"`
}

type HandlerConfig struct {
	AllowedCORSOrigins string `json:"allowed_cors_origins" mapstructure:"allowed_cors_origins" validate:"required"`
}

func MustConfig(log logger.Logger) *Config {
	if err := loadDotEnvUpwards(); err != nil {
		log.Panic(fmt.Sprintf("failed to load .env: %v", err))
	}

	path := fetchConfigPath()
	if path == "" {
		log.Panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Panic("config file does not exist: " + path)
	}

	viper.SetConfigFile(path)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("error reading config file: %v", err)
	}

	var jsonCfg jsonConfig
	if err := viper.Unmarshal(&jsonCfg, viper.DecodeHook(
		mapstructure.StringToTimeDurationHookFunc(),
	)); err != nil {
		log.Panicf("unable to decode config into struct: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(jsonCfg); err != nil {
		log.Panicf("unable to validate config: %v", err)
	}

	var envCfg envConfig
	if err := cleanenv.ReadEnv(&envCfg); err != nil {
		log.Panic("failed to read env config: " + err.Error())
	}

	return &Config{
		jsonConfig: &jsonCfg,
		envConfig:  &envCfg,
	}
}

func loadDotEnvUpwards() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir := wd
	for i := 0; i < 10; i++ {
		envPath := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(envPath); statErr == nil {
			return godotenv.Load(envPath)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return godotenv.Load()
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
