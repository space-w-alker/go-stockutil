package fileutil

import (
	"io/ioutil"
	"os"

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
