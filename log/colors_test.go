package log

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestCSprintf(t *testing.T) {
	assert := require.New(t)

	assert.Equal("this \x1b[30mblack\x1b[0m word", CSprintf("this ${black}black${reset} word"))
	assert.Equal("this \x1b[31mred\x1b[0m word", CSprintf("this ${red}red${reset} word"))
	assert.Equal("this \x1b[32mgreen\x1b[0m word", CSprintf("this ${green}green${reset} word"))
	assert.Equal("this \x1b[33myellow\x1b[0m word", CSprintf("this ${yellow}yellow${reset} word"))
	assert.Equal("this \x1b[34mblue\x1b[0m word", CSprintf("this ${blue}blue${reset} word"))
	assert.Equal("this \x1b[35mmagenta\x1b[0m word", CSprintf("this ${magenta}magenta${reset} word"))
	assert.Equal("this \x1b[36mcyan\x1b[0m word", CSprintf("this ${cyan}cyan${reset} word"))
	assert.Equal("this \x1b[37mwhite\x1b[0m word", CSprintf("this ${white}white${reset} word"))

	assert.Equal("this \x1b[90mblack\x1b[0m word", CSprintf("this ${black+h}black${reset} word"))
	assert.Equal("this \x1b[91mred\x1b[0m word", CSprintf("this ${red+h}red${reset} word"))
	assert.Equal("this \x1b[92mgreen\x1b[0m word", CSprintf("this ${green+h}green${reset} word"))
	assert.Equal("this \x1b[93myellow\x1b[0m word", CSprintf("this ${yellow+h}yellow${reset} word"))
	assert.Equal("this \x1b[94mblue\x1b[0m word", CSprintf("this ${blue+h}blue${reset} word"))
	assert.Equal("this \x1b[95mmagenta\x1b[0m word", CSprintf("this ${magenta+h}magenta${reset} word"))
	assert.Equal("this \x1b[96mcyan\x1b[0m word", CSprintf("this ${cyan+h}cyan${reset} word"))
	assert.Equal("this \x1b[97mwhite\x1b[0m word", CSprintf("this ${white+h}white${reset} word"))

	assert.Equal("this \\[\x1b[90m\\]black\\[\x1b[0m\\] word", TermSprintf("this ${black+h}black${reset} word"))
	assert.Equal("this \\[\x1b[91m\\]red\\[\x1b[0m\\] word", TermSprintf("this ${red+h}red${reset} word"))
	assert.Equal("this \\[\x1b[92m\\]green\\[\x1b[0m\\] word", TermSprintf("this ${green+h}green${reset} word"))
	assert.Equal("this \\[\x1b[93m\\]yellow\\[\x1b[0m\\] word", TermSprintf("this ${yellow+h}yellow${reset} word"))
	assert.Equal("this \\[\x1b[94m\\]blue\\[\x1b[0m\\] word", TermSprintf("this ${blue+h}blue${reset} word"))
	assert.Equal("this \\[\x1b[95m\\]magenta\\[\x1b[0m\\] word", TermSprintf("this ${magenta+h}magenta${reset} word"))
	assert.Equal("this \\[\x1b[96m\\]cyan\\[\x1b[0m\\] word", TermSprintf("this ${cyan+h}cyan${reset} word"))
	assert.Equal("this \\[\x1b[97m\\]white\\[\x1b[0m\\] word", TermSprintf("this ${white+h}white${reset} word"))

	// test color expression stripping
	assert.Equal("this black word", CStripf("this ${black+h}black${reset} word"))
	assert.Equal("this red word", CStripf("this ${red+h}red${reset} word"))
	assert.Equal("this green word", CStripf("this ${green+h}green${reset} word"))
	assert.Equal("this yellow word", CStripf("this ${yellow+h}yellow${reset} word"))
	assert.Equal("this blue word", CStripf("this ${blue+h}blue${reset} word"))
	assert.Equal("this magenta word", CStripf("this ${magenta+h}magenta${reset} word"))
	assert.Equal("this cyan word", CStripf("this ${cyan+h}cyan${reset} word"))
	assert.Equal("this white word", CStripf("this ${white+h}white${reset} word"))
}
