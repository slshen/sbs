package boxscore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaste(t *testing.T) {
	assert := assert.New(t)
	s := paste(`hello
world`, "xxxxxx", 4, 0)
	assert.Equal(`hello    xxxxxx
world
`, s)
}

func TestFirstWord(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("Athletics Mercado", firstWord("Athletics Mercado Walling", 20))
}
