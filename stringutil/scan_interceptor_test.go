package stringutil

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ghetzel/testify/require"
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

func TestScanInterceptorRepeats(t *testing.T) {
	assert := require.New(t)
	the := 0
	father := 0
	had := 0

	splitter := NewScanInterceptor(bufio.ScanLines, map[string]InterceptFunc{
		`the`: func(seq []byte) {
			the += 1
		},
		`Father`: func(seq []byte) {
			father += 1
		},
		`had`: func(seq []byte) {
			had += 1
		},
	})

	data := bytes.NewBuffer([]byte(
		"It was November. Although it was not yet late, the sky was dark when I turned into Laundress " +
			"Passage. Father had finished for the day, switched off the shop lights and closed the shutters; " +
			"but so I would not come home to darkness he had left on the light over the stairs to the flat. " +
			"Through the glass in the door it cast a foolscap rectangle of paleness onto the wet pavement, and " +
			"it was while I was standing in that rectangle, about to turn my key in the door, that I first saw " +
			"the letter. Another white rectangle, it was on the fifth step from the bottom, where I couldn't miss it.\n" +
			"\n" +
			"I closed the door and put the shop key in its usual place behind Bailey's Advanced Principles of Geometry. " +
			"Poor Bailey. No one has wanted his fat gray book for thirty years. Sometimes I wonder what he makes of his " +
			"role as guardian of the bookshop keys. I don't suppose it's the destiny he had in mind for the masterwork " +
			"that he spent two decades writing.",
	))

	scanner := bufio.NewScanner(data)
	scanner.Split(splitter.Scan)

	for scanner.Scan() {
		continue
	}

	assert.NoError(scanner.Err())
	assert.Equal(21, the)
	assert.Equal(1, father)
	assert.Equal(3, had)
}
