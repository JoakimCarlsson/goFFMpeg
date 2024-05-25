package main

import (
	"context"
	"goFFMpeg/ffmpeg"
	"goFFMpeg/ffmpeg/audio"
	"goFFMpeg/ffmpeg/video"
	"log"
	"time"
)

func main() {
	inputPath := "input.mp4"
	outputPath := "output.mp4"

	encode := ffmpeg.Build().
		AddFileInput(inputPath).
		AddFileOutput(outputPath, true, func(options *ffmpeg.Argument) {
			options.
				WithVideoSettings(func(v *video.Video) {
					v.SetCodec(video.H264).
						SetCrf(21)
				}).
				WithAudioSettings(func(a *audio.Audio) {
					a.SetCodec(audio.AAC).
						SetBitrate(128)
				}).
				SetPreset(ffmpeg.Faster).
				WithOverwrite()
		})

	executor := &ffmpeg.Executor{
		OnProgressDataUpdated: func(progressData ffmpeg.ProgressData) {
			log.Printf("Progress: %+v\n", progressData)
		},
		OnOutputDataUpdated: func(output string) {
			log.Printf("Output: %s\n", output)
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if err := executor.Execute(encode, ctx); err != nil {
		log.Fatalf("Error executing FFmpeg command: %v", err)
	}
}
