package config

import (
	"github.com/joho/godotenv"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	envPrefix = "MTPROTO"
)

const (
	configPath = "./config/"
)

const (
	sessionFile configKey = "session_file"
	sessionDir  configKey = "session_directory"
	appID                 = "APP_ID"
	appHash               = "app_hash"
)

type configKey string

type configValue struct {
	Value any `mapstructure:"value"`
}

type Config struct {
	appID   int    `mapstructure:"app_id"`
	appHash string `mapstructure:"app_hash"`

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

		rtConfig.cfg = tCfg

		log.Println("config updated")
	})

	v.WatchConfig()

	return rtConfig
}

func (c *Realtime) GetSessionFile() string {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[sessionFile].String()
}

func (c *Realtime) GetSessionDirectory() string {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[sessionDir].String()
}

func (c *Realtime) GetAppID() int {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.appID
}

func (c *Realtime) GetAppHash() string {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.appHash
}

func (c configValue) String() string {
	return c.Value.(string)
}
