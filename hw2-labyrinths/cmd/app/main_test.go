package main

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/adapter"
)

func TestHelpOutput(t *testing.T) {
	cmd := exec.Command("./app", "--help")
	var errBuf bytes.Buffer
	cmd.Stderr = bufio.NewWriter(&errBuf)

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	//--help return data in stderr
	errBytes := errBuf.Bytes()
	errStr := string(errBytes[:])

	tokens := []string{"generate", "solve", "generate_algorithm", "solve_algorithm", "file",
		"output", "enable_unicode_walls", "height", "width", "start", "goal"}
	for _, token := range tokens {
		if !strings.Contains(errStr, token) {
			t.Fatalf("Output does not contain '%s'", token)

		}
	}
}

func TestValidStartPoint(t *testing.T) {
	//In my realization maze_properties file contains start point
	cmd := exec.Command("./app", "-solve", "-file=maze_pr_wrong_start.txt")
	var errBuf bytes.Buffer
	cmd.Stderr = bufio.NewWriter(&errBuf)
	var outBuf bytes.Buffer
	cmd.Stdout = bufio.NewWriter(&outBuf)

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	errBytes := errBuf.Bytes()
	errStr := strings.Trim(string(errBytes[:]), "\n")

	if errStr != "Error:invalid config" {
		t.Fatalf("Output does not contain '%s'", errStr)

	}
}

func TestMaze1x1(t *testing.T) {
	cmd := exec.Command("./app", "-generate=t",
		"-height=1", "-width=1", "-output=maze_1x1.txt")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	file, err := os.Open("../../mazes/maze_1x1.txt")
	if err != nil {
		t.Fatalf("Failed to open maze file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			adapter.ShowError(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	lineCounter := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 3 {
			t.Fatalf("Maze must contain row length = 3, but got %d ", len(line))
		}
		lineCounter++
	}
	if lineCounter > 3 {
		t.Fatalf("Maze must contain col length = 3, but got %d ", lineCounter)
	}
}

func TestCorrectSolve(t *testing.T) {
	cmd := exec.Command("./app", "-solve", "-file=maze_pr_a_star.txt",
		"-enable_unicode_walls=t", "-output=maze_a_star.txt")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	file, err := os.Open("../../mazes/maze_a_star.txt")
	if err != nil {
		t.Fatalf("Failed to open maze file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			adapter.ShowError(err)
		}
	}()

	reader := bufio.NewReader(file)

	b := make([]byte, 1024)
	n, err := reader.Read(b)
	if err != nil {
		t.Fatalf("Failed to read b: %v", err)
	}

	realMaze := string(b[:n])

	expectedMaze := `┌───────────┬─┬───┬─┐
│O     ░    │¤│░  │¤│
│*╶───┬─┐ ╷ ╵ │ ╷ │ │
│*    │ │ │   │ │ │ │
│*╶───┘ └─┤ ╶─┤ ├─┤ │
│***      │   │ │ │ │
├─╴*╶───┬─┤ ┌─┘ ╵ ╵ │
│  *****│ │ │       │
│ ╷ ╶─┐*╵ ├─┘ ╷ ┌─╴ │
│¤│   │***│***│ │   │
│ └─┐ └─┐*╵*╷*└─┼───┤
│   │  ░│***│***│  ░│
├───┘ ╷ │ ┌─┼─╴*╵ ╶─┤
│     │ │ │ │  *****│
│ ┌─╴ ├─┴─┘ ├─╴ ┌─╴*│
│¤│   │     │  ░│  *│
│ │ ╷ ╵ ╶───┼─╴ ├─╴*│
│¤│ │       │¤  │***│
│ │ │ ╷ ╷ ╶─┤ ╶─┤*╶─┤
│ │ │ │ │   │   │**E│
└─┴─┴─┴─┴───┴───┴───┘`

	if strings.TrimSpace(realMaze) != strings.TrimSpace(expectedMaze) {
		t.Log(realMaze)
		t.Log(expectedMaze)
		t.Fatalf("Got maze isn't expected")
	}
}
