package ginzero

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	gin.DefaultWriter = NewWriter(zerolog.InfoLevel)
	gin.DefaultErrorWriter = NewWriter(zerolog.ErrorLevel)
	gn := gin.New()
	require.NotNil(t, gn)
}
