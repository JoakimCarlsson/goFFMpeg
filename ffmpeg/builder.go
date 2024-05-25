package ffmpeg

import (
	"fmt"
	"goFFMpeg/ffmpeg/audio"
	"goFFMpeg/ffmpeg/video"
	"os"
	"path/filepath"
	"strings"
)

type ArgumentBuilder struct {
	commands  []string
	inputFile string
}

func Build() *ArgumentBuilder {
	return &ArgumentBuilder{}
}

func (b *ArgumentBuilder) AddFileInput(filePath string) *ArgumentBuilder {
	b.inputFile = fmt.Sprintf("-i \"%s\" ", filePath)
	return b
}

func (b *ArgumentBuilder) AddFileOutput(filePath string, createDirectory bool, addArguments func(*Argument)) *ArgumentBuilder {
	if createDirectory {
		directory := filepath.Dir(filePath)
		if _, err := os.Stat(directory); os.IsNotExist(err) && directory != "" {
			err := os.MkdirAll(directory, os.ModePerm)
			if err != nil {
				return nil
			}
		}
	}

	args := &Argument{}
	if addArguments != nil {
		addArguments(args)
	}

	command := b.inputFile + strings.Join(args.Arguments, " ") + fmt.Sprintf(" \"%s\"", filePath)
	b.commands = append(b.commands, command)
	return b
}

func (b *ArgumentBuilder) AddHlsOutput(addArguments func(*HlsArgument)) *ArgumentBuilder {
	args := &HlsArgument{}
	if addArguments != nil {
		addArguments(args)
	}

	for _, argument := range args.Arguments {
		command := b.inputFile + argument
		b.commands = append(b.commands, command)
	}
	return b
}

func (b *ArgumentBuilder) Build() []string {
	return b.commands
}

type Argument struct {
	ArgumentBase
}

func (a *Argument) WithNoAudio() *Argument {
	return a.addArgument("-an")
}

func (a *Argument) WithNoVideo() *Argument {
	return a.addArgument("-vn")
}

func (a *Argument) WithOverwrite() *Argument {
	return a.addArgument("-y")
}

func (a *Argument) WithThreads(threads int) *Argument {
	return a.addArgument(fmt.Sprintf("-threads %d", threads))
}

func (a *Argument) ExtractFrames(frameCount int) *Argument {
	return a.addArgument(fmt.Sprintf("-vframes %d", frameCount))
}

func (a *Argument) SetCustomArgument(argument string) *Argument {
	return a.addArgument(argument)
}

func (a *Argument) SetPreset(preset Preset) *Argument {
	return a.addArgument(fmt.Sprintf("-preset %s", preset))
}

func (a *Argument) WithAudioFilters(options func(*video.AudioFilter)) *Argument {
	args := &video.AudioFilter{}
	options(args)
	a.Arguments = append(a.Arguments, "-af", fmt.Sprintf("\"%s\"", strings.Join(args.Arguments, ", ")))
	return a
}

func (a *Argument) WithMetadata(options func(*video.MetaData)) *Argument {
	args := &video.MetaData{}
	options(args)
	a.Arguments = append(a.Arguments, args.Arguments...)
	return a
}

func (a *Argument) WithVideoFilters(options func(*video.Filter)) *Argument {
	args := &video.Filter{}
	options(args)
	a.Arguments = append(a.Arguments, "-vf", fmt.Sprintf("\"%s\"", strings.Join(args.Arguments, ", ")))
	return a
}

func (a *Argument) WithVideoSettings(options func(*video.Video)) *Argument {
	args := &video.Video{}
	options(args)
	a.Arguments = append(a.Arguments, args.Arguments...)
	return a
}

func (a *Argument) WithAudioSettings(options func(*audio.Audio)) *Argument {
	args := &audio.Audio{}
	options(args)
	a.Arguments = append(a.Arguments, args.Arguments...)
	return a
}

func (a *Argument) addArgument(argument string) *Argument {
	if !strings.HasSuffix(argument, " ") {
		argument += " "
	}
	a.Arguments = append(a.Arguments, argument)
	return a
}
