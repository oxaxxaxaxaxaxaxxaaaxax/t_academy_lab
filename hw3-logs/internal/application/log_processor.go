package application

import (
	"bufio"
	"io"
	"log/slog"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain/models"
)

var (
	ErrLogFormatMismatch = errors.New("Log format mismatch")
	nginxLogRegex        = regexp.MustCompile(
		`^(\S+) - (\S+) \[([^]]+)] "([^"]*)" (\d{3}) (\d+) "([^"]*)" "([^"]*)"$`,
	)
)

type LogSource struct {
	SourceType SourceType
	Path       []string
}

func NewLogSource(sourceType SourceType, path []string) LogSource {
	return LogSource{SourceType: sourceType, Path: path}
}

func NewLogStats(paths []string) LogStats {
	return LogStats{
		Files:                           paths,
		ResponseSizeByLogs:              make([]float64, 0),
		MostFrequencyRequestedResources: make([]ResourcesFrequency, 0),
		PercentageRequestByDayOfWeek:    make(map[time.Weekday]float64),
		WeekByDate:                      make(map[time.Time]time.Weekday),
		RequestByDayOfWeek:              make(map[time.Weekday]int64),
		UniqueDataTransferProtocols:     make(map[string]bool),
		FrequencyCodeResponse:           make(map[int]int64),
		FrequencyRequestedResources:     make(map[string]int64),
	}
}

type Log struct {
	RequestHTTPReferrer   string
	RequestHTTPUserAgent  string
	RequestIP             string
	RequestMethod         string
	RequestProtocol       string
	RequestURL            string
	RequestTime           time.Time
	RequestUser           string
	ResponseBodyBytesSent int64
	ResponseStatus        int
}

type LogStats struct {
	Files                           []string
	RequestCount                    int64
	ResponseAverageSize             float64
	ResponseTotalSize               int64
	ResponseSizeByLogs              []float64
	MostFrequencyRequestedResources []ResourcesFrequency
	PercentageRequestByDayOfWeek    map[time.Weekday]float64
	WeekByDate                      map[time.Time]time.Weekday
	ResponseMaxSize                 int64
	ResponseSizePercentile          float64
	FrequencyCodeResponse           map[int]int64
	FrequencyRequestedResources     map[string]int64
	RequestByDayOfWeek              map[time.Weekday]int64
	//This is set of protocols
	UniqueDataTransferProtocols map[string]bool
}

type SourceType int

const (
	Local SourceType = iota
	Remote
)

func ParsePath(path string) (LogSource, error) {
	if strings.HasPrefix(path, "https://") {
		slog.Debug("using http URL", "path", path)
		return NewLogSource(Remote, []string{path}), nil
	}
	if strings.ContainsAny(path, "*?[") {
		slog.Debug("using * in filepath")
		paths, err := filepath.Glob(path)
		if err != nil {
			return LogSource{}, err
		}
		if len(paths) == 0 {
			return LogSource{}, ErrFileNotFound
		}
		abs := make([]string, 0, len(paths))
		for _, p := range paths {
			a, err := filepath.Abs(p)
			if err != nil {
				return LogSource{}, err
			}
			abs = append(abs, a)
		}
		return NewLogSource(Local, abs), nil

	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return LogSource{}, err
	}
	return NewLogSource(Local, []string{absPath}), nil
}

func (l LogStats) GetLogStats(reader io.Reader) (LogStats, error) {
	slog.Debug("getting log stats")
	s := bufio.NewScanner(reader)
	for s.Scan() {
		logLine := s.Text()
		log, err := ParseLogLine(logLine)
		if err != nil {
			slog.Warn("error with scan log line:", "err", err)
			return LogStats{}, err
		}
		l = l.UpdateStats(log)

	}
	return l, nil
}

func (l LogStats) ConvertToDto() models.DTO {
	slog.Info("convert stats to dto")

	var resourcesInfo = make([]models.ResourcesInfo, 0)
	for resources, frequency := range l.FrequencyRequestedResources {
		resourcesInfo = append(resourcesInfo, models.NewResourcesInfo(resources, frequency))
	}

	var responseCodeInfo = make([]models.ResponseCode, 0)
	for responseCode, frequency := range l.FrequencyCodeResponse {
		responseCodeInfo = append(responseCodeInfo, models.NewResponseCodesInfo(responseCode, frequency))
	}

	var requestsPerDate = make([]models.RequestsPerDate, 0)
	for date, week := range l.WeekByDate {
		requestsPerDate = append(requestsPerDate, models.NewRequestsPerDate(date, week,
			l.RequestByDayOfWeek[week], l.PercentageRequestByDayOfWeek[week]))
	}

	var uniqueProtocols = make([]string, 0)
	for protocol := range l.UniqueDataTransferProtocols {
		uniqueProtocols = append(uniqueProtocols, protocol)
	}

	return models.DTO{
		Files:              l.Files,
		TotalRequestsCount: l.RequestCount,
		ResponseSizeInBytes: models.NewResponseSizeInBytes(l.ResponseAverageSize,
			l.ResponseSizePercentile, l.ResponseMaxSize),
		Resources:         resourcesInfo,
		ResponseCodesInfo: responseCodeInfo,
		RequestsPerDate:   requestsPerDate,
		UniqueProtocols:   uniqueProtocols,
	}
}

func ParseLogLine(line string) (Log, error) {
	matches := nginxLogRegex.FindStringSubmatch(line)
	if matches == nil {
		return Log{}, ErrLogFormatMismatch
	}

	timeLocalStr := matches[3]
	requestStr := matches[4]
	statusStr := matches[5]
	bodyBytesStr := matches[6]

	timeLocal, err := time.Parse(TimeLayout, timeLocalStr)
	if err != nil {
		return Log{}, err
	}

	status, err := strconv.Atoi(statusStr)
	if err != nil {
		return Log{}, err
	}
	bodyBytes, _ := strconv.Atoi(bodyBytesStr)
	parts := strings.SplitN(requestStr, " ", 3)
	method, url, requestProtocol := parts[0], parts[1], parts[2]

	return Log{
		RequestIP:             matches[1],
		RequestUser:           matches[2],
		RequestTime:           timeLocal,
		RequestMethod:         method,
		RequestURL:            url,
		RequestProtocol:       requestProtocol,
		ResponseStatus:        status,
		ResponseBodyBytesSent: int64(bodyBytes),
		RequestHTTPReferrer:   matches[7],
		RequestHTTPUserAgent:  matches[8],
	}, nil
}

func (l LogStats) UpdateStats(log Log) LogStats {
	l.RequestCount += 1
	if l.ResponseMaxSize < log.ResponseBodyBytesSent {
		l.ResponseMaxSize = log.ResponseBodyBytesSent
	}
	l.ResponseSizeByLogs = append(l.ResponseSizeByLogs, float64(log.ResponseBodyBytesSent))
	l.ResponseTotalSize += log.ResponseBodyBytesSent
	l.FrequencyRequestedResources[log.RequestURL]++
	l.FrequencyCodeResponse[log.ResponseStatus]++
	l.RequestByDayOfWeek[log.RequestTime.Weekday()]++
	l.WeekByDate = UpdateWeekByDayMap(l.WeekByDate, log.RequestTime)
	l.UniqueDataTransferProtocols[log.RequestProtocol] = true
	return l
}

func UpdateWeekByDayMap(weekByDate map[time.Time]time.Weekday, requestTime time.Time) map[time.Time]time.Weekday {
	year, month, day := requestTime.Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	if _, ok := weekByDate[date]; ok {
		return weekByDate
	}
	weekByDate[date] = requestTime.Weekday()

	return weekByDate
}

func (l LogStats) CalculateRemainingStats() (LogStats, error) {
	slog.Info("start calculating remaining stats")
	l.ResponseAverageSize = float64(l.ResponseTotalSize) / float64(l.RequestCount)

	percentile, err := stats.Percentile(l.ResponseSizeByLogs, 95)
	if err != nil {
		return LogStats{}, err
	}

	l.ResponseSizePercentile = percentile
	l.MostFrequencyRequestedResources = SortFrequencyResourcesDesc(l.FrequencyRequestedResources)
	l.PercentageRequestByDayOfWeek = CalculatePercentageRequest(l.RequestByDayOfWeek, l.RequestCount)

	return l, nil
}

type ResourcesFrequency struct {
	Res       string
	Frequency int64
}

func SortFrequencyResourcesDesc(frRes map[string]int64) []ResourcesFrequency {
	var ss []ResourcesFrequency

	for k, v := range frRes {
		ss = append(ss, ResourcesFrequency{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Frequency > ss[j].Frequency
	})
	return ss
}

func CalculatePercentageRequest(requestByDayOfWeek map[time.Weekday]int64, requestCount int64) map[time.Weekday]float64 {
	percReq := make(map[time.Weekday]float64)

	for k, v := range requestByDayOfWeek {
		percReq[k] = float64(v) / float64(requestCount)
	}
	return percReq
}
