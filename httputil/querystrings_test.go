package httputil

import (
	"net/url"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestQueryStringModify(t *testing.T) {
	assert := require.New(t)

	u, err := url.Parse(`https://example.com`)
	assert.NoError(err)
	assert.NotNil(u)

	SetQ(u, `test`, false)
	SetQ(u, `test`, true)

	AddQ(u, `test2`, 1)
	AddQ(u, `test2`, 3)

	SetQ(u, `nope`, true)
	DelQ(u, `nope`)

	assert.Equal(u.String(), `https://example.com?test=true&test2=1&test2=3`)
}

func TestQueryStringStringModify(t *testing.T) {
	assert := require.New(t)

	assert.Equal(
		`https://example.com?test=false`,
		SetQString(`https://example.com`, `test`, false),
	)

	assert.Equal(
		`https://example.com?test=true`,
		SetQString(`https://example.com`, `test`, true),
	)

	x := `https://example.com`
	x = AddQString(x, `test2`, 1)
	x = AddQString(x, `test2`, 3)

	assert.Equal(`https://example.com?test2=1&test2=3`, x)

	assert.Equal(`https://example.com`, DelQString(`https://example.com?nope=lol`, `nope`))
}
