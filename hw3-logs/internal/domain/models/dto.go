package models

import "time"

type DTO struct {
	Files               []string
	TotalRequestsCount  int64
	ResponseSizeInBytes ResponseSizeInBytes
	Resources           []ResourcesInfo
	ResponseCodesInfo   []ResponseCode
	RequestsPerDate     []RequestsPerDate
	UniqueProtocols     []string
}

type ResponseSizeInBytes struct {
	Average float64
	Max     int64
	P95     float64
}

func NewResponseSizeInBytes(average, p95 float64, max int64) ResponseSizeInBytes {
	return ResponseSizeInBytes{
		Average: average,
		Max:     max,
		P95:     p95,
	}
}

type ResourcesInfo struct {
	Resource           string
	TotalRequestsCount int64
}

func NewResourcesInfo(resource string, totalRequestsCount int64) ResourcesInfo {
	return ResourcesInfo{
		Resource:           resource,
		TotalRequestsCount: totalRequestsCount,
	}
}

type ResponseCode struct {
	Code                int
	TotalResponsesCount int64
}

func NewResponseCodesInfo(code int, totalResponsesCount int64) ResponseCode {
	return ResponseCode{code, totalResponsesCount}
}

type RequestsPerDate struct {
	Date                    time.Time
	Weekday                 time.Weekday
	TotalRequestsCount      int64
	TotalRequestsPercentage float64
}

func NewRequestsPerDate(date time.Time, weekday time.Weekday, totalRequestsCount int64, totalRequestsPercentage float64) RequestsPerDate {
	return RequestsPerDate{
		Date:                    date,
		Weekday:                 weekday,
		TotalRequestsCount:      totalRequestsCount,
		TotalRequestsPercentage: totalRequestsPercentage,
	}
}
