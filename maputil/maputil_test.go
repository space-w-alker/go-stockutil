package maputil

import (
	"testing"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/timeutil"
	"github.com/ghetzel/testify/assert"
	"github.com/ghetzel/testify/require"
)

var testJsonPathData = map[string]interface{}{
	"expensive": 10,
	"store": map[string]interface{}{
		"book": []map[string]interface{}{
			{
				"category": "reference",
				"author":   "Nigel Rees",
				"title":    "Sayings of the Century",
				"price":    8.95,
			},
			{
				"category": "fiction",
				"author":   "Evelyn Waugh",
				"title":    "Sword of Honour",
				"price":    12.99,
			},
			{
				"category": "fiction",
				"author":   "Herman Melville",
				"title":    "Moby Dick",
				"isbn":     "0-553-21311-3",
				"price":    8.99,
			},
			{
				"category": "fiction",
				"author":   "J. R. R. Tolkien",
				"title":    "The Lord of the Rings",
				"isbn":     "0-395-19395-8",
				"price":    22.99,
			},
		},
		"bicycle": map[string]interface{}{
			"color": "red",
			"price": 19.95,
		},
	},
}

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

func TestJSONPath(t *testing.T) {
	var fn = func(query string) interface{} {
		var out, err = JSONPath(testJsonPathData, query)

		assert.NoError(t, err, query)
		return out
	}

	for query, wanted := range map[string]interface{}{
		`$.store.book[*].author`: []interface{}{
			"Nigel Rees",
			"Evelyn Waugh",
			"Herman Melville",
			"J. R. R. Tolkien",
		},
		`$..author`: []interface{}{
			"Nigel Rees",
			"Evelyn Waugh",
			"Herman Melville",
			"J. R. R. Tolkien",
		},
		`$..price`: []interface{}{
			8.95,
			12.99,
			8.99,
			22.99,
			19.95,
		},
		`$..book[?(.price <= 8.99)].title`: []interface{}{
			"Moby Dick",
			"Sayings of the Century",
		},
		`$..book[?(.price > 10.0)].title`: []interface{}{
			"Sword of Honour",
			"The Lord of the Rings",
		},
	} {
		assert.ElementsMatch(t, wanted, fn(query), query)
	}
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
