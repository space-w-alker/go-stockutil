package timeutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDuration(t *testing.T) {
	assert := require.New(t)

	v, err := ParseDuration(``)
	assert.NoError(err)
	assert.Zero(v)

	for in, out := range map[string]time.Duration{
		`4h`:    time.Duration(4 * time.Hour),
		`4H`:    time.Duration(4 * time.Hour),
		`1d`:    time.Duration(24 * time.Hour),
		`1day`:  time.Duration(24 * time.Hour),
		`1days`: time.Duration(24 * time.Hour),
		`5 years 4 weeks 3 days 2 hours 1 minute`: time.Duration(44546*time.Hour) + time.Minute,
		`1d1h`:  time.Duration(25 * time.Hour),
		`1d 1h`: time.Duration(25 * time.Hour),
	} {
		v, err := ParseDuration(in)
		assert.NoError(err)
		assert.Equal(out, v, fmt.Sprintf("in=%v", in))
	}
}

func TestDurationHMS(t *testing.T) {
	assert := require.New(t)

	h, m, s := DurationHMS(0)
	assert.Equal(0, h)
	assert.Equal(0, m)
	assert.Equal(0, s)

	h, m, s = DurationHMS(time.Second)
	assert.Equal(0, h)
	assert.Equal(0, m)
	assert.Equal(1, s)

	h, m, s = DurationHMS(time.Minute)
	assert.Equal(0, h)
	assert.Equal(1, m)
	assert.Equal(0, s)

	h, m, s = DurationHMS(time.Hour)
	assert.Equal(1, h)
	assert.Equal(0, m)
	assert.Equal(0, s)

	h, m, s = DurationHMS(59 * time.Second)
	assert.Equal(0, h)
	assert.Equal(0, m)
	assert.Equal(59, s)

	h, m, s = DurationHMS(119 * time.Second)
	assert.Equal(0, h)
	assert.Equal(1, m)
	assert.Equal(59, s)

	h, m, s = DurationHMS(23*time.Hour + 4*time.Minute + 13*time.Second)
	assert.Equal(23, h)
	assert.Equal(4, m)
	assert.Equal(13, s)
}

func TestFormatTimer(t *testing.T) {
	assert := require.New(t)

	out := FormatTimer(0)
	assert.Equal(`0:00`, out)

	out = FormatTimer(time.Second)
	assert.Equal(`0:01`, out)

	out = FormatTimer(time.Minute)
	assert.Equal(`1:00`, out)

	out = FormatTimer(time.Hour)
	assert.Equal(`1:00:00`, out)

	out = FormatTimer(59 * time.Second)
	assert.Equal(`0:59`, out)

	out = FormatTimer(119 * time.Second)
	assert.Equal(`1:59`, out)

	out = FormatTimer(23*time.Hour + 4*time.Minute + 13*time.Second)
	assert.Equal(`23:04:13`, out)
}
