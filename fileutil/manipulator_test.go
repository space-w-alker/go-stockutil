package fileutil

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testTextPre = "// this is a cool test\n"
const testTextBody = ("\n" +
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut \n" +
	"labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco \n" +
	"laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in \n" +
	"voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat \n" +
	"cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." +
	"\n" +
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut \n" +
	"labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco \n" +
	"laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in \n" +
	"voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat \n" +
	"cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

const testText = testTextPre + testTextBody

func tt() io.Reader {
	return bytes.NewBufferString(testText)
}

func TestReadManipulatorNoOp(t *testing.T) {
	assert := require.New(t)

	rm := NewReadManipulator(tt())
	out, err := ioutil.ReadAll(rm)
	assert.NoError(err)
	assert.Equal(testText, string(out))
}

func TestReadManipulatorDoNothingFunction(t *testing.T) {
	assert := require.New(t)

	rm := NewReadManipulator(tt(), func(in []byte) ([]byte, error) {
		return in, nil
	})

	out, err := ioutil.ReadAll(rm)
	assert.NoError(err)
	assert.Equal(testText, string(out))
}

func TestReadManipulatorRemoveComments(t *testing.T) {
	assert := require.New(t)

	rm := NewReadManipulator(tt(), func(in []byte) ([]byte, error) {
		line := string(in)
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, `//`) {
			return nil, nil
		} else {
			return in, nil
		}
	})

	out, err := ioutil.ReadAll(rm)
	assert.NoError(err)
	assert.Equal("\n"+testTextBody, string(out))
}

func TestReadManipulatorRemoveBlankLines(t *testing.T) {
	assert := require.New(t)

	rm := NewReadManipulator(tt(), func(in []byte) ([]byte, error) {
		if line := strings.TrimSpace(string(in)); len(line) == 0 {
			return nil, SkipToken
		} else {
			return in, nil
		}
	})

	out, err := ioutil.ReadAll(rm)
	assert.NoError(err)
	assert.Equal(strings.TrimSpace(testTextPre)+testTextBody, string(out))
}

func TestReadManipulatorDolorToBacon(t *testing.T) {
	assert := require.New(t)

	rm := NewReadManipulator(tt(), func(in []byte) ([]byte, error) {
		line := strings.Replace(string(in), `dolor`, `bacon`, -1)
		return []byte(line), nil
	})

	out, err := ioutil.ReadAll(rm)
	assert.NoError(err)
	assert.Equal(
		strings.Replace(testText, `dolor`, `bacon`, -1),
		string(out),
	)
}
