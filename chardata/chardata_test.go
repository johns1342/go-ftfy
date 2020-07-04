package chardata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	assert.Equal(t, "latin-1", CharmapEncodings[0], "CharmapEncodings is missing latin-1")
}

func TestRegexASCII(t *testing.T) {
	assert := assert.New(t)
	re := buildRegexes()["ascii"]
	assert.True(re.MatchString("123\x34ABC!.|abcd\r\n"), "Failed valid ASCII match.")
	assert.False(re.MatchString("123\xF4ABC!.|abcd\r\n"), "Failed to detect non-ASCII value.")

}
