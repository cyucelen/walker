package walker_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cyucelen/walker"
	"github.com/streetbyters/aduket"
	"github.com/stretchr/testify/assert"
)

func TestApiWalker(t *testing.T) {
	client, requestRecorder := aduket.NewServer(
		http.MethodGet, "/books",
		aduket.StatusCode(http.StatusOK),
		aduket.ByteBody([]byte{'w'}),
	)
	defer client.Close()

	requestBuilder := func(start, fetchCount int) (*http.Request, error) {
		return http.NewRequest(http.MethodGet, fmt.Sprintf("%s/books?start=%d&count=%d", client.URL, start, fetchCount), http.NoBody)
	}

	var actualResponseBody []byte
	sink := func(res *http.Response, stop func()) error {
		actualResponseBody, _ = ioutil.ReadAll(res.Body)
		return nil
	}

	apiWalker := walker.NewApiWalker(
		http.DefaultClient,
		requestBuilder,
		sink,
		walker.WithLimiter(walker.ConstantLimiter(10)),
		walker.WithMaxBatchSize(10),
		walker.WithPagination(walker.CursorPagination{}),
		walker.WithParallelism(1),
	)

	apiWalker.Walk()

	requestRecorder.AssertQueryParamEqual(t, "start", []string{"0"})
	requestRecorder.AssertQueryParamEqual(t, "count", []string{"10"})
	assert.Equal(t, []byte{'w'}, actualResponseBody)
}
