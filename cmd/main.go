package main

import (
	"find/internal/config"
	"find/internal/note"
	"find/internal/order"
	"find/internal/reminder"
	"find/internal/stdin"
	"fmt"
	"os"
)

func init() {
	err := note.Check()
	if err != nil {
		fmt.Printf("check note error: %s\n", err.Error())
	}
}

func main() {
	if config.ReminderEnabled {
		err := reminder.Start()
		if err != nil {
			fmt.Printf("start reminder error: %s\n", err.Error())
		}
	}

	fmt.Println("=================")
	fmt.Println("Welcome to FIND!")
	fmt.Println("=================")
	for true {
		fmt.Print("[FIND]# ")
		input, err := stdin.ReadString()
		if err != nil {
			fmt.Printf("read input error: %s\n", err.Error())
			continue
		}

		param := order.Param(input)

		switch order.Order(input) {
		case order.Find:
			_, err = note.Find(param, true, false, true)
			if err != nil {
				fmt.Printf("find %s error: %s\n", param, err.Error())
				continue
			}
		case order.Add:
			same, err := note.Find(note.GetKey(param), true, true, false)
			if err != nil {
				fmt.Printf("find %s error: %s\n", note.GetKey(param), err.Error())
				continue
			}
			if len(same) > 0 {
				fmt.Println("Duplicate key.")
				continue
			}
			err = note.Write(&[]string{param}, os.O_APPEND)
			if err != nil {
				fmt.Printf("append note error: %s\n", err.Error())
				continue
			}
			succeed()
		case order.Delete:
			var fast bool
			var all bool
			fast, param = order.Fast(param)
			all, param = order.All(param)
			err = note.Delete(param, !fast, !all)
			if err != nil {
				fmt.Printf("delete note error: %s\n", err.Error())
				continue
			}
			succeed()
		case order.Modify:
			err = note.Modify(param)
			if err != nil {
				fmt.Printf("modify %s error: %s\n", param, err.Error())
				continue
			}
			succeed()
		case order.Exit:
			os.Exit(1)
		}
	}
}

func succeed() {
	fmt.Println("Succeed.")
}
