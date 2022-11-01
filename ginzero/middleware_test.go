package ginzero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

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
