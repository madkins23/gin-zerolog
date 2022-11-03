package ginzero

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCode = float64(200)
	testFrag = "fragment"
	testIP   = "1.2.3.4"
	testLvl  = "debug"
	testMeth = "GET"
	testMsg  = "Request"
	testPath = "Never/More"
	testPort = "666"
	testQry  = "goober=snoofus"
)

func TestLogger(t *testing.T) {
	logFn := Logger()
	require.NotNil(t, logFn)
	ctxt, engine := gin.CreateTestContext(httptest.NewRecorder())
	ctxt.Request = &http.Request{
		RemoteAddr: testIP + ":" + testPort + "/" + testPath + "?" + testQry + "#" + testFrag,
		Method:     testMeth,
		URL: &url.URL{
			Host:        testIP + ":" + testPort,
			Path:        testPath,
			RawPath:     testPath + "?" + testQry,
			ForceQuery:  false,
			RawQuery:    testQry,
			Fragment:    testFrag,
			RawFragment: testFrag,
		},
	}
	require.NotNil(t, ctxt)
	require.NotNil(t, engine)

	// Trap output from running log function.
	zLog := log.Logger
	defer func() { log.Logger = zLog }()
	buffer := &bytes.Buffer{}
	log.Logger = zerolog.New(buffer)
	logFn(ctxt)

	// Process log output which is in JSON.
	var record map[string]interface{}
	require.NoError(t, json.Unmarshal(buffer.Bytes(), &record))
	assert.Equal(t, testCode, record["code"])
	assert.Equal(t, testIP, record["ip"])
	assert.Equal(t, testLvl, record["level"])
	assert.Equal(t, testMsg, record["message"])
	assert.Equal(t, testMeth, record["meth"])
	assert.Equal(t, testPath+"?"+testQry, record["path"])

	fmt.Println(buffer.String())
}

//////////////////////////////////////////////////////////////////////////

func ExampleLogger() {
	const port = "55555"

	// Switch zerolog to console mode.
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:     os.Stdout,
		NoColor: true,
		// Don't show duration or time as they mess up the Output comparison.
		FieldsExclude: []string{"dur"},
		PartsExclude:  []string{"time"},
	})

	// Get rid of standard errors from gin.
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	// Create the gin router.
	router := gin.Default()

	// Install the ginzero Logger() to re-route middleware logging to zerolog.
	router.Use(Logger())
	// Use() this object before the routing calls below.

	// Add gin routing.
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below.
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Unable to listen and serve")
		}
	}()

	// Wait for server to start to avoid connection error.
	// TODO: Is there a better way to wait for the server to start?
	time.Sleep(250 * time.Millisecond)

	// Ping this server.
	if response, err := http.Get("http://:" + port + "/ping"); err != nil {
		log.Error().Err(err).Interface("response", response).Msg("Ping response")
	}

	// Shutdown the server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	// Output:
	// DBG Request code=200 ip=127.0.0.1 meth=GET path=/ping
}
