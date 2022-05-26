// Package note implements methods for handling user's local data file which is called note.
package note

import (
	"find/internal/backup"
	"find/internal/config"
	"find/internal/constant"
	"find/internal/files"
	"find/internal/redish"
	"find/internal/stdin"
	"fmt"
	"os"
	"strings"
)

// Path is the path of note which is loaded from config.
var Path string

func init() {
	Path = config.NotePath
}

// Check is used to ensure that the note is available,
// and then synchronize if redis config is available too.
func Check() error {
	fileInfo, err := os.Stat(Path)
	isNewNote := err != nil
	var file *os.File
	if isNewNote {
		// If note not exists, then create.
		file, err = os.Create(Path)
		if err != nil {
			return fmt.Errorf("create %s error: %v", Path, err)
		}
	} else {
		// If note exists, then open.
		file, err = os.OpenFile(Path, os.O_RDWR, 0)
		if err != nil {
			return fmt.Errorf("open %s error: %v", Path, err)
		}
	}
	defer func() {
		_ = file.Close()
	}()

	if config.RdsKey != "" && redish.Client != nil {
		// If redis config is available, then sync.
		err = backup.Sync(isNewNote, file, fileInfo)
		if err != nil {
			return fmt.Errorf("sync backup error: %v", err)
		}
	}

	return nil
}

// Find is used to lookup note according to keyword from user's input and multiple options,
// returning a string slice of result and error.
func Find(keyword string, include bool, accurate bool, print bool) ([]string, error) {
	keywords := strings.Split(keyword, " ")
	notes, err := files.ReadLinesFromPath(Path)
	if err != nil {
		return nil, fmt.Errorf("read lines from %s error: %v", Path, err)
	}

	results := make([]string, 0)
	for _, note := range notes {
		var hit bool
		if accurate {
			hit = GetKey(note) == keyword
		} else {
			hit = containsAll(GetKey(note), keywords)
		}

		if !include {
			hit = !hit
		}

		if hit {
			results = append(results, note)
		}
	}

	if print {
		if len(results) == 0 {
			fmt.Println("Empty result.")
		} else {
			for _, result := range results {
				fmt.Printf("%s: %s\n", GetKey(result), GetVal(result))
			}
		}
	}

	return results, nil
}

// containsAll is used to judge if source string contains all target strings ignoring the case,
// returning true if contains all and false otherwise.
func containsAll(source string, targets []string) bool {
	for _, target := range targets {
		if !strings.Contains(strings.ToLower(source), strings.ToLower(target)) {
			return false
		}
	}
	return true
}

// Write is used to persist notes into local data file by specified mode,
// and will asynchronously update the backup if the redis config is available.
func Write(notes *[]string, mod int) error {
	file, err := os.OpenFile(Path, mod, 0)
	if err != nil {
		return fmt.Errorf("open %s error: %v", Path, err)
	}

	err = files.WriteLinesToFile(file, notes)
	if err != nil {
		return fmt.Errorf("write %v to %v error: %v", notes, file, err)
	}

	go func() {
		err = Check()
		if err != nil {
			fmt.Printf("check note error: %s", err.Error())
		}
	}()

	return nil
}

// Delete is used to remove note from local data file after optional confirming,
// and will asynchronously update the backup if the redis config is available.
func Delete(keyword string, confirm bool, accurate bool) error {
	var yesOrNo string
	if confirm {
		fmt.Println("Will delete:")
		_, err := Find(keyword, true, accurate, true)
		if err != nil {
			return fmt.Errorf("find %s error: %v", keyword, err)
		}
		fmt.Println("Sure delete? [y/n]")
		tmp, err := stdin.ReadString()
		if err != nil {
			return fmt.Errorf("read input error: %v", err)
		}
		yesOrNo = tmp
	} else {
		yesOrNo = constant.Yes
	}

	if yesOrNo == constant.Yes {
		notes, err := Find(keyword, false, accurate, false)
		if err != nil {
			return fmt.Errorf("find %s error: %v", keyword, err)
		}
		err = Write(&notes, os.O_RDWR|os.O_TRUNC)
		if err != nil {
			return fmt.Errorf("write note error: %v", err)
		}
	}

	return nil
}

// Modify is used to update note in local date file by delete and write,
// and will asynchronously update the backup if the redis config is available.
func Modify(note string) error {
	err := Delete(GetKey(note), false, true)
	if err != nil {
		return fmt.Errorf("delete %s error: %v", GetKey(note), err)
	}
	err = Write(&[]string{note}, os.O_APPEND)
	if err != nil {
		return fmt.Errorf("append %s error: %v", note, err)
	}
	return nil
}

// GetKey is used to parse key of note, returning the key.
func GetKey(note string) string {
	i := strings.Index(note, ":")
	if i == -1 {
		return ""
	}

	return note[:i]
}

// GetVal is used to parse value of note, returning the value.
func GetVal(note string) string {
	i := strings.Index(note, ":")
	if i == -1 {
		return ""
	}

	return note[i+1:]
}
