// Package config gathers all configs other modules may need.
package config

import (
	"find/internal/files"
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Config map to program config yaml.
type Config struct {
	Find struct {
		NotePath string `yaml:"notePath"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"find"`
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
	Reminder struct {
		Enabled         bool   `yaml:"enabled"`
		Type            string `yaml:"type"`
		IntervalSeconds int    `yaml:"interval-seconds"`
	} `yaml:"reminder"`
	Email struct {
		From     string   `yaml:"from"`
		To       []string `yaml:"to"`
		AuthCode string   `yaml:"authCode"`
	}
}

// all configs
var Conf Config

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// default unchangeable config path
	confPath := homedir + "\\FIND.yml"

	// If config yaml not exists, then init.
	if _, err = os.Stat(confPath); err != nil {
		err = initYaml(confPath)
		if err != nil {
			fmt.Printf("init yaml error: %s\n", err)
			return
		}
	}

	// parse config from yaml
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Printf("read file from %s error: %s\n", confPath, err.Error())
		return
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		fmt.Printf("unmarshal yaml error: %s\n", err.Error())
		return
	}
}

// RedisKey is used to get a redis key for representing backup.
func RedisKey() string {
	if Conf.Find.Username != "" && Conf.Find.Password != "" {
		return Conf.Find.Username + ":" + Conf.Find.Password
	}
	return ""
}

// RedisConf is used to get redis config for backup.
func RedisConf() *redis.Options {
	return &redis.Options{
		Addr:     Conf.Redis.Address,
		Password: Conf.Redis.Password,
		DB:       Conf.Redis.Db,
	}
}

// initYaml is used to create config yaml and write initial configs.
func initYaml(confPath string) error {
	file, err := os.Create(confPath)
	if err != nil {
		return fmt.Errorf("create config file error: %v", err)
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	initialConfigs := []string{
		"find:",
		"  notePath: " + homedir + "\\FIND.txt",
		"  username:",
		"  password:",
		"redis:",
		"  address:",
		"  password:",
		"  db:",
		"reminder:",
		"  enabled: true",
		"  type: win",
		"  interval-seconds: 1",
		"email:",
		"  from:",
		"  to:",
		"  authCode:",
	}
	err = files.WriteLinesToFile(file, &initialConfigs)
	if err != nil {
		return fmt.Errorf("write initial configs to file error: %v", err)
	}

	return nil
}
