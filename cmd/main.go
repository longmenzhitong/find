package main

import (
	"find/internal/config"
	"find/internal/logs"
	"find/internal/note"
	"find/internal/order"
	"find/internal/reminder"
	"find/internal/stdin"
	"find/internal/weather"
	"fmt"
	"os"
)

func init() {
	err := note.Check()
	if err != nil {
		logs.Error("check note error: %s\n", err.Error())
	}
	if config.Conf.Reminder.Enabled {
		err := reminder.Start()
		if err != nil {
			logs.Error("start reminder error: %s\n", err.Error())
		}
	}
	logs.Info("initialization finished")
}

func main() {
	logs.Info("FIND started")

	fmt.Println("=================")
	fmt.Println("Welcome to FIND!")
	fmt.Println("=================")
	for true {
		fmt.Print("[FIND]# ")
		input, err := stdin.ReadString()
		if err != nil {
			logs.Error("read input error: %s\n", err.Error())
			continue
		}

		var fast bool
		var all bool

		param := order.Param(input)

		switch order.Order(input) {
		case order.Find:
			_, err = note.Find(param, true, false, true)
			if err != nil {
				logs.Error("find %s error: %s\n", param, err.Error())
				continue
			}
		case order.Add:
			same, err := note.Find(note.GetKey(param), true, true, false)
			if err != nil {
				logs.Error("find %s before add error: %s\n", note.GetKey(param), err.Error())
				continue
			}
			if len(same) > 0 {
				fmt.Println("Duplicate key.")
				continue
			}
			err = note.Write(&[]string{param}, os.O_APPEND)
			if err != nil {
				logs.Error("add %s error: %s\n", param, err.Error())
				continue
			}
			succeed()
		case order.Delete:
			fast, param = order.Fast(param)
			all, param = order.All(param)
			err = note.Delete(param, !fast, !all)
			if err != nil {
				logs.Error("delete %s error: %s\n", param, err.Error())
				continue
			}
			succeed()
		case order.Modify:
			err = note.Modify(param)
			if err != nil {
				logs.Error("modify %s error: %s\n", param, err.Error())
				continue
			}
			succeed()
		case order.Weather:
			all, param = order.All(param)
			if param == "" {
				fmt.Println("Need address.")
				continue
			}
			err = weather.Search(param, all)
			if err != nil {
				logs.Error("search weather of %s error: %s\n", param, err.Error())
				continue
			}
		case order.Exit:
			os.Exit(1)
		}
	}
}

func succeed() {
	fmt.Println("Succeed.")
}
