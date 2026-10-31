package config

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	envPrefix = "MTPROTO"
)

const (
	configPath = "./config/"
)

type configKey string

type configValue struct {
	Value any `mapstructure:"value"`
}

type Config struct {
	appID       int
	appHash     string
	phoneNumber string
	password    string

	values map[configKey]configValue
}

type Realtime struct {
	m sync.RWMutex

	cfg *Config
}

func MustLoad() *Realtime {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	config := &Config{
		values: make(map[configKey]configValue),
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("can't load .env: %v", err)
	}

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	config.appID = v.GetInt(appID)
	config.appHash = v.GetString(appHash)
	config.phoneNumber = v.GetString(phoneNumber)
	config.password = v.GetString(password)

	err = v.ReadInConfig()
	if err != nil {
		log.Fatalf("MustLoad.ReadInConfig: %v", err)
	}

	err = v.Unmarshal(&config.values)
	if err != nil {
		log.Fatalf("MustLoad.Unmarshal: %v", err)
	}

	rtConfig := &Realtime{
		cfg: config,
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		tCfg := &Config{
			values: make(map[configKey]configValue),
		}

		errCh := v.Unmarshal(&tCfg.values)
		if errCh != nil {
			log.Printf("can't update config: %v\n", errCh)

			return
		}

		rtConfig.m.Lock()
		defer rtConfig.m.Unlock()

		tCfg.appHash = config.appHash
		tCfg.appID = config.appID
		tCfg.phoneNumber = config.phoneNumber
		tCfg.password = config.password

		rtConfig.cfg = tCfg

		log.Println("config updated")
	})

	v.WatchConfig()

	return rtConfig
}

func (c configValue) String() string {
	return c.Value.(string)
}

func (c configValue) UInt8() uint8 {
	return c.Value.(uint8)
}

func (c configValue) UInt32() uint32 {
	return c.Value.(uint32)
}
