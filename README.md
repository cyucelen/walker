<p align="center">
  <img height="400px" src="assets/logo.png">
</p>
<p align="center">
    <b>Seamlessly fetch paginated data from any source!</b>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/cyucelen/walker?tab=doc">
    <img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white" alt="godoc" title="godoc"/>
  </a>
  <a href="https://github.com/cyucelen/walker/tags">
    <img src="https://img.shields.io/github/v/tag/cyucelen/walker" alt="semver tag" title="semver tag"/>
  </a>
  <a href="https://github.com/cyucelen/walker/actions/workflows/go.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/cyucelen/walker/go.yml?branch=master" />
  </a>
  <a href="https://codecov.io/gh/cyucelen/walker">
    <img src="https://codecov.io/gh/cyucelen/walker/branch/master/graph/badge.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/cyucelen/walker">
    <img src="https://goreportcard.com/badge/github.com/cyucelen/walker" />
  </a>
  <a href="https://github.com/cyucelen/walker/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/cyucelen/walker.svg">
  </a>
</p>

# walker

Walker simplifies the process of fetching paginated data from any data source. With Walker, you can easily configure the start position and count of documents to fetch, depending on your needs. Additionally, Walker supports parallel processing, allowing you to fetch data more efficiently and at a faster rate.

The real purpose of the library is to provide a solution for walking through the pagination of API endpoints. With the `NewApiWalker`, you can easily fetch data from any paginated API endpoint and process the data concurrently. You can also create your own custom walker to fit your specific use case.

## Features

* Provides a walker to paginate through the pagination of API endpoint. This is for scraping an API, if such a term exists.
* `cursor` and `offset` pagination strategies.
* Fetching and processing data concurrently without any effort.
* Total fetch count limiting
* Rate limiting

## Examples

### Basic Usage

```go
func source(start, fetchCount int) ([]int, error) {
	return []int{start, fetchCount}, nil
}

func sink(result []int, stop func()) error {
	fmt.Println(result)
	return nil
}

func main() {
	walker.New(source,sink).Walk()
}
```
**Output:**
```
[0 10]
[1 10]
[4 10]
[2 10]
[3 10]
[5 10]
[8 10]
[9 10]
[7 10]
[6 10]
...
to Infinity
```

* `source` function will receive `start` as the page number and `count` as the number of documents. Use this values to fetch data from your source.
* `sink` function will receive the result you returned from `source` and a `stop` function. You can save the results in this function and decide to stop sourcing any further pages depending on your results by calling `stop` function, otherwise it will continue to forever unless [a limit provided](#configuration).
* Beware of order is not ensured since source and sink functions called concurrently.

### Walking through the pagination of API endpoints 

**Fetching all the breweries from `Open Brewery DB`:**

```go
func buildRequest(start, fetchCount int) (*http.Request, error) {
	url := fmt.Sprintf("https://api.openbrewerydb.org/breweries?page=%d&per_page=%d", start, fetchCount)
	return http.NewRequest(http.MethodGet, url, http.NoBody)
}

func sink(res *http.Response, stop func()) error {
	var payload []map[string]any
	json.NewDecoder(res.Body).Decode(&payload)

	if len(payload) == 0 {
		stop()
		return nil
	}

	return saveBreweries(payload)
}

func main() {
	walker.NewApiWalker(http.DefaultClient, buildRequest, sink).Walk()
}
```

To create API walker you just need to provide: 
* `RequestBuilder` function to create http request using provided values
* `sink` function to process the http response

Check [examples](/example/) for more usecases.

## Configuration

| Option           | Description                                            | Default                     | Available Values                                          |
| ---------------- | ------------------------------------------------------ | --------------------------- | --------------------------------------------------------- |
| WithPagination   | Defines the pagination strategy                        | `walker.OffsetPagination{}` | `walker.OffsetPagination{}`, `walker.CursorPagination{}`  |
| WithMaxBatchSize | Defines limit for document count to stop after reached | `10`                        | `int`                                                     |
| WithParallelism  | Defines number of workers to run provided source       | `runtime.NumCPU()`          | `int`                                                     |
| WithLimiter      | Defines limit for document count to stop after reached | `walker.InfiniteLimiter()`  | `walker.InfiniteLimiter()`, `walker.ConstantLimiter(int)` |
| WithRateLimit    | Defines rate limit by **count** and per **duration**   | `unlimited`                 | `(int, time.Duration)`                                    |
| WithContext      | Defines context                                        | `context.Background()`      | `context.Context`                                         |


## Contribution

I would like to accept any contributions to make `walker` better and feature rich. Feel free to contribute with your usecase!