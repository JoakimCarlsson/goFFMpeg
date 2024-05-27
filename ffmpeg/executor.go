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
	"sync"
)

type Executor struct {
	OnProgressDataUpdated func(progressData ProgressData)
	OnOutputDataUpdated   func(output string)
}

func (e *Executor) Execute(command *ArgumentBuilder, ctx context.Context, runInParallel bool) error {
	if runInParallel {
		return e.executeInParallel(command, ctx)
	}
	return e.executeSequentially(command, ctx)
}

func (e *Executor) executeInParallel(command *ArgumentBuilder, ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	for _, argument := range command.Build() {
		wg.Add(1)
		go func(arg string) {
			defer wg.Done()
			if err := e.executeCommand(arg, ctx); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}(argument)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err, ok := <-errChan; ok {
		return err
	}

	return nil
}

func (e *Executor) executeSequentially(command *ArgumentBuilder, ctx context.Context) error {
	for _, argument := range command.Build() {
		if err := e.executeCommand(argument, ctx); err != nil {
			return err
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

func (e *Executor) executeCommand(argument string, ctx context.Context) error {
	args := splitArgs(argument)

	processCmd := exec.CommandContext(ctx, getFFMpegPath(), args...)
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

	return nil
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
