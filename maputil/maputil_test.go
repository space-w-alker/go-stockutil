package maputil

import (
	"testing"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/timeutil"
	"github.com/ghetzel/testify/require"
)

func TestRxMapFmt(t *testing.T) {
	assert := require.New(t)

	m := rxutil.Match(rxMapFmt, `${testing.the.thing}`)
	assert.NotNil(m)
	assert.Equal(map[string]string{
		`key`:      `testing.the.thing`,
		`fallback`: ``,
		`fmt`:      ``,
	}, m.NamedCaptures())

	m = rxutil.Match(rxMapFmt, `${testing.the.thing:%48s}`)
	assert.NotNil(m)
	assert.Equal(map[string]string{
		`key`:      `testing.the.thing`,
		`fallback`: ``,
		`fmt`:      `%48s`,
	}, m.NamedCaptures())

	m = rxutil.Match(rxMapFmt, `${testing.the.thing|fallback.value}`)
	assert.NotNil(m)
	assert.Equal(map[string]string{
		`key`:      `testing.the.thing`,
		`fallback`: `fallback.value`,
		`fmt`:      ``,
	}, m.NamedCaptures())

	m = rxutil.Match(rxMapFmt, `${testing|the|thing|fallback.value}`)
	assert.NotNil(m)
	assert.Equal(map[string]string{
		`key`:      `testing`,
		`fallback`: `the|thing|fallback.value`,
		`fmt`:      ``,
	}, m.NamedCaptures())

	m = rxutil.Match(rxMapFmt, `${testing|the|thing|fallback.value:%48s}`)
	assert.NotNil(m)
	assert.Equal(map[string]string{
		`key`:      `testing`,
		`fallback`: `the|thing|fallback.value`,
		`fmt`:      `%48s`,
	}, m.NamedCaptures())
}

func TestSprintf(t *testing.T) {
	assert := require.New(t)

	assert.Equal(
		`Hello guest! Your IP is: (unknown)`,
		Sprintf("Hello ${username|guest}! Your IP is: ${ipaddress|(unknown)}"),
	)

	assert.Equal(
		`Hello guest! Your IP is: 127.0.0.1`,
		Sprintf("Hello ${username|guest}! Your IP is: ${ipaddress|(unknown)}", map[string]interface{}{
			`ipaddress`: `127.0.0.1`,
		}),
	)

	assert.Equal(
		`Hello guest! Your IP is: 127.0.0.1`,
		Sprintf("Hello ${username|guest}! Your IP is: ${ipaddress|(unknown)}", map[string]interface{}{
			`username`:  ``,
			`ipaddress`: `127.0.0.1`,
		}),
	)

	assert.Equal(
		`Hello friend! Your IP is: (unknown)`,
		Sprintf("Hello ${username|guest}! Your IP is: ${ipaddress|(unknown)}", map[string]interface{}{
			`username`: `friend`,
		}),
	)

	assert.Equal(
		`Hello friend! Your IP is: (unknown)`,
		Sprintf("Hello ${username|guest}! Your IP is: ${ipaddress|(unknown)}", map[string]interface{}{
			`username`:  `friend`,
			`ipaddress`: ``,
		}),
	)

	assert.Equal(
		`Hello friend! Your IP is: 127.0.0.1`,
		Sprintf("Hello ${username|guest}! Your IP is: ${ipaddress|(unknown)}", map[string]interface{}{
			`username`:  `friend`,
			`ipaddress`: `127.0.0.1`,
		}),
	)
}

func TestSprintfFormatting(t *testing.T) {
	assert := require.New(t)

	assert.Equal(
		`Hello guest     ! Your IP is:       (unknown)`,
		Sprintf("Hello ${username|guest:%-10s}! Your IP is: ${ipaddress|(unknown):%15s}"),
	)

	assert.Equal(
		`Hello guest     ! Your IP is:       127.0.0.1`,
		Sprintf("Hello ${username|guest:%-10s}! Your IP is: ${ipaddress|(unknown):%15s}", map[string]interface{}{
			`ipaddress`: `127.0.0.1`,
		}),
	)

	assert.Equal(
		`Hello guest     ! Your IP is:       127.0.0.1`,
		Sprintf("Hello ${username|guest:%-10s}! Your IP is: ${ipaddress|(unknown):%15s}", map[string]interface{}{
			`username`:  ``,
			`ipaddress`: `127.0.0.1`,
		}),
	)

	assert.Equal(
		`Hello friend    ! Your IP is:       (unknown)`,
		Sprintf("Hello ${username|guest:%-10s}! Your IP is: ${ipaddress|(unknown):%15s}", map[string]interface{}{
			`username`: `friend`,
		}),
	)

	assert.Equal(
		`Hello friend    ! Your IP is:       (unknown)`,
		Sprintf("Hello ${username|guest:%-10s}! Your IP is: ${ipaddress|(unknown):%15s}", map[string]interface{}{
			`username`:  `friend`,
			`ipaddress`: ``,
		}),
	)

	assert.Equal(
		`Hello friend    ! Your IP is:       127.0.0.1`,
		Sprintf("Hello ${username|guest:%-10s}! Your IP is: ${ipaddress|(unknown):%15s}", map[string]interface{}{
			`username`:  `friend`,
			`ipaddress`: `127.0.0.1`,
		}),
	)
}

func TestSprintfFormatTime(t *testing.T) {
	assert := require.New(t)

	assert.Equal(
		`the time is: 2006-01-02T15:04:05-07:00`,
		Sprintf("the time is: ${now}", map[string]interface{}{
			`now`: timeutil.ReferenceTime(),
		}),
	)

	assert.Equal(
		`the time is: January 2, 2006 (3:04pm)`,
		Sprintf("the time is: ${now:%January 2, 2006 (3:04pm)}", map[string]interface{}{
			`now`: timeutil.ReferenceTime(),
		}),
	)
}

func ExamplePrintf_usingDefaultValues() {
	Printf("Hello ${username|guest}! Your IP is: ${ipaddress|(unknown)}")
	// Output: Hello guest! Your IP is: (unknown)
}

func ExamplePrintf_suppliedWithData() {
	Printf("Hello ${username|guest}! Your IP is: ${ipaddress|(unknown)}", map[string]interface{}{
		`username`:  `friend`,
		`ipaddress`: `127.0.0.1`,
	})

	// Output: Hello friend! Your IP is: 127.0.0.1
}

func ExamplePrintf_deeplyNestedKeys() {
	Printf("Hello ${details.0.value|guest}! Your IP is: ${details.1.value|(unknown)}", map[string]interface{}{
		`details`: []map[string]interface{}{
			{
				`key`:   `username`,
				`value`: `friend`,
			}, {
				`key`:   `ipaddress`,
				`value`: `127.0.0.1`,
			},
		},
	})

	// Output: Hello friend! Your IP is: 127.0.0.1
}
