package main

import (
	"find/internal/note"
	"find/internal/order"
	"find/internal/stdin"
	"fmt"
	"os"
	"strings"
)

func init() {
	err := note.Check()
	if err != nil {
		fmt.Printf("check note error: %s\n", err.Error())
	}
}

func main() {
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

		_order := order.Parse(input)
		keyword := strings.TrimSpace(strings.TrimPrefix(input, _order))

		switch _order {
		case order.Find:
			_, err = note.Find(keyword, true, false, true)
			if err != nil {
				fmt.Printf("find %s error: %s\n", keyword, err.Error())
				continue
			}
		case order.Add:
			same, err := note.Find(note.GetKey(keyword), true, true, false)
			if err != nil {
				fmt.Printf("find %s error: %s\n", note.GetKey(keyword), err.Error())
				continue
			}
			if len(same) > 0 {
				fmt.Println("Duplicate key.")
				continue
			}
			err = note.Write(&[]string{keyword}, os.O_APPEND)
			if err != nil {
				fmt.Printf("append note error: %s\n", err.Error())
				continue
			}
			succeed()
		case order.Delete:
			err = note.Delete(keyword, true)
			if err != nil {
				fmt.Printf("delete note error: %s\n", err.Error())
				continue
			}
			succeed()
		case order.FastDelete:
			err = note.Delete(keyword, false)
			if err != nil {
				fmt.Printf("delete note error: %s\n", err.Error())
				continue
			}
			succeed()
		case order.Modify:
			err = note.Delete(note.GetKey(keyword), false)
			if err != nil {
				fmt.Printf("delete %s error: %s\n", note.GetKey(keyword), err.Error())
				continue
			}
			err = note.Write(&[]string{keyword}, os.O_APPEND)
			if err != nil {
				fmt.Printf("append note error: %s\n", err.Error())
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
