package stringutil

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanInterceptorNothing(t *testing.T) {
	assert := require.New(t)
	var lines []string

	splitter := NewScanInterceptor(bufio.ScanLines)
	data := bytes.NewBuffer([]byte("first\nsecond\nthird\n"))

	scanner := bufio.NewScanner(data)
	scanner.Split(splitter.Scan)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	assert.NoError(scanner.Err())
	assert.Equal([]string{
		`first`,
		`second`,
		`third`,
	}, lines)
}

// test single subsequence
// ---------------------------------------------------------------------------------------------
func TestScanInterceptorSingle(t *testing.T) {
	assert := require.New(t)
	errors := 0
	prompts := 0
	var lines []string

	splitter := NewScanInterceptor(bufio.ScanLines, map[string]InterceptFunc{
		`[error] `: func(seq []byte) {
			errors += 1
		},
		` password: `: func(seq []byte) {
			prompts += 1
		},
		`Password: `: func(seq []byte) {
			prompts += 1
		},
	})

	data := bytes.NewBuffer([]byte(
		"Warning: Permanently added '[127.0.0.1]:2200' (ECDSA) to the list of known hosts.\n" +
			"test@127.0.0.1's password: ",
	))

	scanner := bufio.NewScanner(data)
	scanner.Split(splitter.Scan)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	assert.NoError(scanner.Err())
	assert.Equal(0, errors)
	assert.Equal(1, prompts)
	assert.Equal([]string{
		`Warning: Permanently added '[127.0.0.1]:2200' (ECDSA) to the list of known hosts.`,
		`test@127.0.0.1's password: `,
	}, lines)
}

// test multiple subsequences
// ---------------------------------------------------------------------------------------------
func TestScanInterceptorMultiple(t *testing.T) {
	assert := require.New(t)
	errors := 0
	prompts := 0
	var lines []string

	splitter := NewScanInterceptor(bufio.ScanLines, map[string]InterceptFunc{
		`[error] `: func(seq []byte) {
			errors += 1
		},
		` password: `: func(seq []byte) {
			prompts += 1
		},
		`Password: `: func(seq []byte) {
			prompts += 1
		},
	})

	data := bytes.NewBuffer([]byte(
		"Password: [error] something cool went wrong\n" +
			"test@127.0.0.1's password: ",
	))

	scanner := bufio.NewScanner(data)
	scanner.Split(splitter.Scan)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	assert.NoError(scanner.Err())
	assert.Equal(1, errors)
	assert.Equal(2, prompts)
	assert.Equal([]string{
		`Password: [error] something cool went wrong`,
		`test@127.0.0.1's password: `,
	}, lines)
}

// test add intercept after the fact
// ---------------------------------------------------------------------------------------------
func TestScanInterceptorAddIntercept(t *testing.T) {
	assert := require.New(t)
	errors := 0
	warnings := 0
	var lines []string

	splitter := NewScanInterceptor(bufio.ScanLines, map[string]InterceptFunc{
		`[error] `: func(seq []byte) {
			errors += 1
		},
	})

	data := bytes.NewBuffer([]byte(
		"Warning: Permanently added '[127.0.0.1]:2200' (ECDSA) to the list of known hosts.\n" +
			"[error] something cool went wrong\n",
	))

	scanner := bufio.NewScanner(data)
	scanner.Split(splitter.Scan)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	assert.NoError(scanner.Err())
	assert.Equal(1, errors)
	assert.Equal(0, warnings)
	assert.Equal([]string{
		`Warning: Permanently added '[127.0.0.1]:2200' (ECDSA) to the list of known hosts.`,
		`[error] something cool went wrong`,
	}, lines)

	// new scanner, same interceptor, add new data

	scanner = bufio.NewScanner(data)
	scanner.Split(splitter.Scan)

	splitter.Intercept(`Warning:`, func(seq []byte) {
		warnings += 1
	})

	lines = nil
	data.WriteString("some cool stuff going on OH NOOOO Warning: NOOOOOOO\n")

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	assert.NoError(scanner.Err())
	assert.Equal(1, warnings)
	assert.Equal([]string{
		`some cool stuff going on OH NOOOO Warning: NOOOOOOO`,
	}, lines)
}

func TestScanInterceptorBinarySubsequence(t *testing.T) {
	assert := require.New(t)
	terminators := 0

	splitter := NewScanInterceptor(bufio.ScanBytes)
	data := bytes.NewBuffer([]byte{
		0x71, 0x00, 0x5d, 0x13, 0xfe, 0x05, 0xff, 0xff,
		0xe7, 0xfe, 0x00, 0x16, 0x20, 0x02, 0x07, 0x5d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xaa, 0x55,
	})

	splitter.Intercept(string([]byte{0xAA, 0x55}), func(seq []byte) {
		terminators += 1
	})

	scanner := bufio.NewScanner(data)
	scanner.Split(splitter.Scan)

	for scanner.Scan() {
		continue
	}

	assert.NoError(scanner.Err())
	assert.Equal(1, terminators)
}
