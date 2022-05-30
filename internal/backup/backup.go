// Package backup implements methods for handling backup located on redis.
package backup

import (
	"encoding/json"
	"find/internal/config"
	"find/internal/files"
	"find/internal/logs"
	"find/internal/redish"
	"fmt"
	"github.com/go-redis/redis"
	"os"
)

// rds is a pointer of redis client.
var rds *redis.Client

// rdsKey is a redis key for representing backup.
var rdsKey string

func init() {
	rds = redish.Client
	rdsKey = config.RedisKey()
}

// Sync is used to pull or push redis backup, decided by different cases.
func Sync(isNewNote bool, file *os.File, fileInfo os.FileInfo) error {
	cmd := rds.ZCard(rdsKey)
	size, err := cmd.Result()
	if err != nil {
		return fmt.Errorf("get result of %v error: %v", cmd, err)
	}

	// case1: If file is new, and redis backup is not empty, then pull.
	if isNewNote && size > 0 {
		err = pull(file)
		if err != nil {
			return fmt.Errorf("pull backup to %v error: %v", file, err)
		}
		return nil
	}

	lastModTime := float64(fileInfo.ModTime().Unix())
	// case2: If file is not new, and redis backup is empty, then push.
	if !isNewNote && size == 0 {
		err = push(file, lastModTime)
		if err != nil {
			return fmt.Errorf("push backup from %v error: %v", file, err)
		}
		return nil
	}

	// case3: If file is not new, and redis backup is not empty, then compare
	// file's last modify time and redis' last backup time to decide pull or push.
	if !isNewNote && size > 0 {
		latestBak, err := getLatest()
		if err != nil {
			return fmt.Errorf("get latest backup error: %v", err)
		}
		lastBakTime := latestBak.Score
		if lastBakTime > lastModTime {
			err = pull(file)
			if err != nil {
				return fmt.Errorf("pull backup to %v error: %v", file, err)
			}
			return nil
		}
		if lastBakTime < lastModTime {
			err = push(file, lastModTime)
			if err != nil {
				return fmt.Errorf("push backup from %v error: %v", file, err)
			}
			return nil
		}
	}
	return nil
}

// pull is used to sync newest backup from redis to file.
func pull(file *os.File) error {
	logs.Info("backup: pull start")
	bak, err := getLatest()
	if err != nil {
		return fmt.Errorf("get latest bak error: %v", err)
	}

	var notes []string
	jsonNotes := fmt.Sprintf("%s", bak.Member)
	err = json.Unmarshal([]byte(jsonNotes), &notes)
	if err != nil {
		return fmt.Errorf("json unmarshal of %v error: %v", jsonNotes, err)
	}

	if len(notes) > 0 {
		err = files.WriteLinesToFile(file, &notes)
		if err != nil {
			return fmt.Errorf("write %v to %v error: %v", notes, file, err)
		}
	}

	logs.Info("backup: pull finished")
	return nil
}

// getLatest is used to fetch newest backup from redis, returning a pointer of redis zset and error.
func getLatest() (*redis.Z, error) {
	cmd := rds.ZRangeWithScores(rdsKey, -1, -1)
	bak, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("get result of %v error: %v", cmd, err)
	}
	return &bak[0], nil
}

// push is used to sync newest backup from file to redis.
func push(file *os.File, lastModTime float64) error {
	logs.Info("backup: push start")
	notes, err := files.ReadLinesFromFile(file)
	if err != nil {
		return fmt.Errorf("read lines from %s error: %v", config.Conf.Find.NotePath, err)
	}

	if len(notes) == 0 {
		return nil
	}

	jsonNotes, err := json.Marshal(notes)
	if err != nil {
		return fmt.Errorf("json marshal of %v error: %v", notes, err)
	}

	jsonStr := string(jsonNotes)
	rds.ZAdd(rdsKey, redis.Z{
		// file's last modify time
		Score: lastModTime,
		// json of file's data
		Member: jsonStr,
	})
	logs.Info("backup: push finished")
	return nil
}
