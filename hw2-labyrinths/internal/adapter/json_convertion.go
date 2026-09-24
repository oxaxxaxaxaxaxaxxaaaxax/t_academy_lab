package adapter

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

func LoadGraphFromJSON(path string) (domain.Graph, error) {
	var g domain.Graph

	file, err := os.Open(path)
	if err != nil {
		return g, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return g, err
	}

	err = json.Unmarshal(data, &g)
	return g, err
}

func SaveGraphToJSON(path string, g domain.Graph) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	return w.Flush()
}
