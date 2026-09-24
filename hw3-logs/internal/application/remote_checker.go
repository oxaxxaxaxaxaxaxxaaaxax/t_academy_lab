package application

import (
	"net/http"
	"time"
)

const TimeoutValSeconds int = 5

type RemoteChecker struct {
}

func NewRemoteChecker() RemoteChecker {
	return RemoteChecker{}
}

func (RemoteChecker) Check(filePath []string) error {
	client := http.Client{Timeout: time.Duration(TimeoutValSeconds) * time.Second}
	httpPath := filePath[0]

	request, err := http.NewRequest(http.MethodGet, httpPath, nil)
	if err != nil {
		return ErrInvalidPathToFile
	}

	response, err := client.Do(request)
	if err != nil {
		return ErrInvalidPathToFile
	}

	if response.StatusCode == http.StatusNotFound {
		return ErrFileNotFound
	}
	return nil
}
