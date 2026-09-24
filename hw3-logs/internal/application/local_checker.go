package application

import "strings"

type LocalChecker struct {
}

func NewLocalChecker() LocalChecker {
	return LocalChecker{}
}

func (LocalChecker) Check(filePath []string) error {
	for _, file := range filePath {
		if !strings.HasSuffix(file, ".txt") && !strings.HasSuffix(file, ".log") {
			return ErrInvalidPathToFile
		}

		exist, err := Exists(file)
		if !exist {
			return ErrFileNotFound
		}
		if err != nil {
			return err
		}
	}
	return nil
}
