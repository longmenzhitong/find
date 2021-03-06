// Package reminder implements methods for reminding user of the to-do notes.
package reminder

import (
	"find/internal/config"
	"find/internal/logs"
	"find/internal/note"
	"fmt"
	"github.com/go-toast/toast"
	"github.com/jordan-wright/email"
	"github.com/robfig/cron"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

const needRemind = "remind@"
const reminded = "reminded@"
const reminderTypeWindows = "win"
const reminderTypeEmail = "email"

// mutex is used to ensure that the reminder is checking to-do notes serially.
var mutex sync.Mutex

// Start is used to start a reminder, checking to-do notes with specified interval seconds,
// sending notifications for necessary. It'll modify 'remind@' to 'reminded@'
// for notified notes to avoid duplicate notifications.
func Start() error {
	logs.Info("reminder: start start")
	c := cron.New()
	spec := fmt.Sprintf("*/%d * * * * ?", config.Conf.Reminder.IntervalSeconds)
	err := c.AddFunc(spec, func() {
		mutex.Lock()
		logs.Debug("reminder: check start")
		notes, err := note.Find("todo", true, false, false)
		if err != nil {
			logs.Error("find todo error: %s\n", err.Error())
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
				logs.Error("parse remind time of %s error: %s\n", timeStr, err.Error())
				continue
			}

			if time.Now().Unix() > remindTime {
				remindSucceed := false
				if strings.Contains(config.Conf.Reminder.Type, reminderTypeWindows) {
					err = remindByWindows(key, val)
					if err != nil {
						logs.Error("remind %s by windows error: %s\n", key, err.Error())
					} else {
						remindSucceed = true
					}
				}

				if strings.Contains(config.Conf.Reminder.Type, reminderTypeEmail) {
					err = remindByEmail(key, val)
					if err != nil {
						logs.Error("remind %s by email error: %s\n", key, err.Error())
					} else {
						remindSucceed = true
					}
				}

				if remindSucceed {
					newNote := strings.ReplaceAll(_note, needRemind, reminded)
					err = note.Modify(newNote)
					if err != nil {
						logs.Error("modify %s error: %s\n", newNote, err.Error())
					}
				}
			}
		}
		logs.Debug("reminder: check finished")
		mutex.Unlock()
	})
	if err != nil {
		return fmt.Errorf("add func to cron error: %v", err)
	}
	c.Start()
	logs.Info("reminder: start finished")
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

// remindByWindows is used to send a windows notification.
func remindByWindows(title, message string) error {
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

// remindByEmail is used to send a email.
func remindByEmail(title, message string) error {
	em := email.NewEmail()
	from := config.Conf.Reminder.Email.From
	em.From = fmt.Sprintf("FIND <%s>", from)
	em.To = config.Conf.Reminder.Email.To
	em.Subject = title
	em.Text = []byte(message)

	addr := config.Conf.Reminder.Email.Server
	err := em.Send(addr, smtp.PlainAuth("", from, config.Conf.Reminder.Email.AuthCode, strings.Split(addr, ":")[0]))
	if err != nil {
		return fmt.Errorf("send email error: %v", err)
	}
	return nil
}
