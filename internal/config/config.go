// Package config gathers all configs other modules may need.
package config

import (
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Config map to program config yaml.
type Config struct {
	Find struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"find"`
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
	Reminder struct {
		Enabled         bool `yaml:"enabled"`
		IntervalSeconds int  `yaml:"interval-seconds"`
	} `yaml:"reminder"`
}

// all configs
var ConfPath string
var NotePath string
var Conf Config

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// default unchangeable config path
	ConfPath = homedir + "\\FIND.yml"
	// default changeable note path
	NotePath = homedir + "\\FIND.txt"

	// If config file not exists, then create.
	if _, err = os.Stat(ConfPath); err != nil {
		_, err = os.Create(ConfPath)
		if err != nil {
			fmt.Printf("create config file error: %s\n", err.Error())
			return
		}
		return
	}

	// parse config from yaml
	file, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		fmt.Printf("read file from %s error: %s\n", ConfPath, err.Error())
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
