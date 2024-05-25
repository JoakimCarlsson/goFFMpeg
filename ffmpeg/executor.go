package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Executor struct {
	OnProgressDataUpdated func(progressData ProgressData)
	OnOutputDataUpdated   func(output string)
}

func (e *Executor) Execute(command *ArgumentBuilder, ctx context.Context) error {
	for _, argument := range command.Build() {
		fullCommand := getFFMpegPath() + " " + argument
		fmt.Println("Executing command:", fullCommand)

		args := splitArgs(argument)
		fmt.Printf("args: %+v\n", args)

		processCmd := exec.CommandContext(ctx, getFFMpegPath(), args...)
		fmt.Println("Command: ", processCmd.String())
		processCmd.Stdout = os.Stdout

		stderr, err := processCmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("failed to get stderr pipe: %w", err)
		}

		if err := processCmd.Start(); err != nil {
			return fmt.Errorf("failed to start FFmpeg process: %w", err)
		}

		go e.readStream(stderr, ctx)

		if err := processCmd.Wait(); err != nil {
			return fmt.Errorf("FFMpeg process exited with error: %w", err)
		}
	}
	return nil
}

func splitArgs(cmd string) []string {
	var args []string
	var currentArg strings.Builder
	inQuotes := false

	for i := 0; i < len(cmd); i++ {
		char := cmd[i]
		switch char {
		case ' ':
			if inQuotes {
				currentArg.WriteByte(char)
			} else {
				if currentArg.Len() > 0 {
					args = append(args, currentArg.String())
					currentArg.Reset()
				}
			}
		case '"':
			inQuotes = !inQuotes
		default:
			currentArg.WriteByte(char)
		}
	}
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}
	return args
}

func (e *Executor) readStream(stream io.ReadCloser, ctx context.Context) {
	defer func(stream io.ReadCloser) {
		err := stream.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error closing stream: %v\n", err)
		}
	}(stream)

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()
		select {
		case <-ctx.Done():
			return
		default:
			if progressData, ok := parseProgressData(line); ok {
				if e.OnProgressDataUpdated != nil {
					e.OnProgressDataUpdated(progressData)
				}
			} else {
				if e.OnOutputDataUpdated != nil {
					e.OnOutputDataUpdated(line)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading from stream: %v\n", err)
	}
}

func getFFMpegPath() string {
	ffmpegExecutableName := "ffmpeg"
	if runtime.GOOS == "windows" {
		ffmpegExecutableName += ".exe"
	}

	path := os.Getenv("PATH")
	pathElements := filepath.SplitList(path)

	for _, pathElement := range pathElements {
		tempFullPath := filepath.Join(pathElement, ffmpegExecutableName)
		if _, err := os.Stat(tempFullPath); err == nil {
			return tempFullPath
		}
	}

	panic("FFMpeg executable not found in any searched location")
}
