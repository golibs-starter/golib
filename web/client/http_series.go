package client

const (
	SeriesInformational HttpSeries = 1
	SeriesSuccessful    HttpSeries = 2
	SeriesRedirection   HttpSeries = 3
	SeriesClientError   HttpSeries = 4
	SeriesServerError   HttpSeries = 5
)

type HttpSeries int

func NewHttpSeries(statusCode int) HttpSeries {
	return HttpSeries(statusCode / 100)
}

func (h HttpSeries) Is(series HttpSeries) bool {
	return h == series
}

func (h HttpSeries) IsError() bool {
	return h.Is(SeriesClientError) || h.Is(SeriesServerError)
}
