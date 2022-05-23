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
		case "username":
			Username = val
		case "password":
			Password = val
		case "rdsaddr":
			RdsConf.Addr = val
		case "rdspswd":
			RdsConf.Password = val
		case "rdsdb":
			RdsConf.DB, err = strconv.Atoi(val)
			if err != nil {
				fmt.Printf("invalid redis db: %s", val)
			}
		case "notePath":
			NotePath = val
		default:
		}
	}

	if Username != "" && Password != "" {
		// a redis key for representing backup
		RdsKey = Username + ":" + Password
	}
}
