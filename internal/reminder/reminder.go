// Package reminder implements methods for reminding user of the to-do notes.
package reminder

import (
	"find/internal/config"
	"find/internal/note"
	"fmt"
	"github.com/go-toast/toast"
	"github.com/robfig/cron"
	"strings"
	"time"
)

const needRemind = "remind@"
const reminded = "reminded@"

// Start is used to start a reminder, checking to-do notes with specified interval seconds,
// sending a windows desktop notification for necessary. It'll modify 'remind@' to 'reminded@'
// for notified notes to avoid duplicate notifications.
func Start() error {
	c := cron.New()
	spec := fmt.Sprintf("*/%d * * * * ?", config.ReminderIntervalSeconds)
	err := c.AddFunc(spec, func() {
		notes, err := note.Find("todo", true, false, false)
		if err != nil {
			fmt.Printf("find todo error: %s\n", err.Error())
			return
		}

		for _, _note := range notes {
			key := note.GetKey(_note)
			val := note.GetVal(_note)

			if !strings.Contains(val, needRemind) {
				continue
			}

			timeStr := strings.TrimSpace(strings.Split(val, needRemind)[1])
			remindTime, err := parseRemindTime(timeStr)
			if err != nil {
				fmt.Printf("parse remind time of %s error: %s\n", timeStr, err.Error())
				continue
			}

			if time.Now().Unix() > remindTime {
				err = remind(key, val)
				if err != nil {
					fmt.Printf("remind %s error: %s\n", key, err.Error())
					continue
				}
				newNote := strings.ReplaceAll(_note, needRemind, reminded)
				err = note.Modify(newNote)
				if err != nil {
					fmt.Printf("modify %s error: %s\n", newNote, err.Error())
				}
			}
		}
	})
	if err != nil {
		return fmt.Errorf("add func to cron error: %v", err)
	}
	c.Start()
	return nil
}

// parseRemindTime is used to parse remindTime from string, returning remindTime(accurate to minutes) and error.
func parseRemindTime(timeStr string) (int64, error) {
	now := time.Now()
	_time, err := time.Parse("15:04", timeStr)
	if err == nil {
		return time.Date(now.Year(), now.Month(), now.Day(), _time.Hour(), _time.Minute(), 0, 0, time.Local).Unix(), nil
	}
	_time, err = time.Parse("2006-01-02 15:04", timeStr)
	if err == nil {
		return time.Date(_time.Year(), _time.Month(), _time.Day(), _time.Hour(), _time.Minute(), 0, 0, time.Local).Unix(), nil
	}
	return -1, err
}

// remind is used to send a windows desktop notification.
func remind(title, message string) error {
	n := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog",
		Title:   title,
		Message: message,
		//Icon: "go.png", // This file must exist (remove this line if it doesn't)
		Actions: []toast.Action{
			{"protocol", "OK", ""},
		},
	}
	return n.Push()
}
