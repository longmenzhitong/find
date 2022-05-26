// Package config gathers all configs other modules may need.
package config

import (
	"find/internal/files"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"strings"
)

// all configs
var Username string
var Password string
var ConfPath string
var NotePath string
var RdsConf *redis.Options
var RdsKey string
var ReminderEnabled = true
var ReminderIntervalSeconds = 1

func init() {
	RdsConf = &redis.Options{}

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// default unchangeable config path
	ConfPath = homedir + "\\FIND.conf"
	// default changeable note path
	NotePath = homedir + "\\FIND.txt"

	if _, err = os.Stat(ConfPath); err != nil {
		// If config file not exists, then create.
		_, err = os.Create(ConfPath)
		if err != nil {
			fmt.Printf("create config file error: %s\n", err.Error())
			return
		}
	}

	// If config file exists, then read.
	lines, err := files.ReadLinesFromPath(ConfPath)
	if err != nil {
		fmt.Printf("read config file error: %s\n", err.Error())
		return
	}

	// Parse config and load into memory.
	for _, line := range lines {
		i := strings.Index(line, ":")
		key := line[:i]
		val := line[i+1:]
		switch key {
		case "find.username":
			Username = val
		case "find.password":
			Password = val
		case "rds.addr":
			RdsConf.Addr = val
		case "rds.password":
			RdsConf.Password = val
		case "rds.db":
			RdsConf.DB, err = strconv.Atoi(val)
			if err != nil {
				fmt.Printf("invalid redis db: %s\n", val)
			}
		case "notePath":
			NotePath = val
		case "reminder.enabled":
			ReminderEnabled, err = strconv.ParseBool(val)
			if err != nil {
				fmt.Printf("invalid reminder enabled: %s\n", val)
			}
		case "reminder.interval-seconds":
			ReminderIntervalSeconds, err = strconv.Atoi(val)
			if err != nil {
				fmt.Printf("invalid reminder interval seconds: %s\n", val)
			}
		default:
		}
	}

	if Username != "" && Password != "" {
		// a redis key for representing backup
		RdsKey = Username + ":" + Password
	}
}
