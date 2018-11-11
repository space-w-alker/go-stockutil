package fileutil

import (
	"io/ioutil"
	"os"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

func IsTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

func ReadAll(filename string) ([]byte, error) {
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		return ioutil.ReadAll(file)
	} else {
		return nil, err
	}
}

func ReadAllString(filename string) (string, error) {
	if data, err := ReadAll(filename); err == nil {
		return string(data), nil
	} else {
		return ``, err
	}
}

func ReadAllLines(filename string) ([]string, error) {
	if data, err := ReadAllString(filename); err == nil {
		return strings.Split(data, "\n"), nil
	} else {
		return nil, err
	}
}

func ReadFirstLine(filename string) (string, error) {
	if lines, err := ReadAllLines(filename); err == nil {
		return lines[0], nil
	} else {
		return ``, err
	}
}

func MustReadAll(filename string) []byte {
	if data, err := ReadAll(filename); err == nil {
		return data
	} else {
		panic(err.Error())
	}
}

func MustReadAllString(filename string) string {
	if data, err := ReadAllString(filename); err == nil {
		return data
	} else {
		panic(err.Error())
	}
}
