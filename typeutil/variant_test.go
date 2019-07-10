package typeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVariant(t *testing.T) {
	assert := require.New(t)

	assert.Equal(`test`, Variant{`test`}.String())
	assert.True(Variant{`True`}.Bool())
	assert.True(Variant{`true`}.Bool())
	assert.True(Variant{`TRUE`}.Bool())
	assert.True(Variant{`1`}.Bool())
	assert.True(Variant{`dennis`}.Bool())
	assert.True(Variant{`0.000000001`}.Bool())
	assert.False(Variant{`False`}.Bool())
	assert.False(Variant{`false`}.Bool())
	assert.False(Variant{`0`}.Bool())
	assert.False(Variant{`0.0`}.Bool())
	assert.Equal(int64(1), Variant{1}.Int())
	assert.Equal(int64(1), Variant{1.9}.Int())
	assert.Equal(float64(1.9), Variant{1.9}.Float())
	assert.True(time.Unix(1500000000, 0).Equal(Variant{1500000000}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`1500000000`}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`2017-07-14T02:40:00Z`}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`2017-07-14T02:40:00-00:00`}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`2017-07-13T22:40:00-04:00`}.Time()))
	assert.Equal([]byte{0x74, 0x65, 0x73, 0x74}, Variant{`test`}.Bytes())

	assert.Equal(map[Variant]Variant{
		V(`test`):  V(1),
		V(`other`): V(2.4),
	}, V(map[string]interface{}{
		`test`:  1,
		`other`: 2.4,
	}).Map())

	type vStructOne struct {
		Name    string
		Age     int
		Pi      float64
		enabled bool
	}

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`Age`):  V(42),
		V(`Pi`):   V(3.1415),
	}, V(vStructOne{
		Name:    `test`,
		Age:     42,
		Pi:      3.1415,
		enabled: true,
	}).Map())

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`Age`):  V(42),
		V(`Pi`):   V(3.1415),
	}, V(&vStructOne{
		Name:    `test`,
		Age:     42,
		Pi:      3.1415,
		enabled: true,
	}).Map())

	type vStructTagged struct {
		Name    string
		Age     int     `testaroo:"age"`
		Pi      float64 `testaroo:"pi,omitempty"`
		enabled bool
	}

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`age`):  V(42),
	}, V(vStructTagged{
		Name:    `test`,
		Age:     42,
		enabled: true,
	}).Map(`testaroo`))

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`age`):  V(42),
	}, V(&vStructTagged{
		Name:    `test`,
		Age:     42,
		enabled: true,
	}).Map(`testaroo`))

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`Age`):  V(42),
		V(`Pi`):   V(float64(0)),
	}, V(vStructTagged{
		Name:    `test`,
		Age:     42,
		enabled: true,
	}).Map())

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`Age`):  V(42),
		V(`Pi`):   V(float64(0)),
	}, V(&vStructTagged{
		Name:    `test`,
		Age:     42,
		enabled: true,
	}).Map())
}
