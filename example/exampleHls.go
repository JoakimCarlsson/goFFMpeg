package main

import (
	"context"
	"github.com/JoakimCarlsson/goFFMpeg/ffmpeg"
	"github.com/JoakimCarlsson/goFFMpeg/ffmpeg/audio"
	"github.com/JoakimCarlsson/goFFMpeg/ffmpeg/video"
	"log"
	"time"
)

func main() {
	inputPath := "input.mp4"

	hls := ffmpeg.Build().
		AddFileInput(inputPath).
		AddHlsOutput(func(options *ffmpeg.HlsArgument) {
			options.
				AddProfile(func(x *ffmpeg.HlsArgument) {
					x.SetHlsTime(4).
						SetHlsPlaylistType("vod").
						SetHlsSegmentFilename("hls\\audio\\meow_%03d.ts", true).
						SetManifestOutput("hls\\audio\\meow.m3u8", true).
						WithAudioSettings(func(v *audio.Audio) {
							v.SetCodec(audio.AAC).
								SetChannels(audio.Stereo).
								SetBitrate(198).
								SetSamplingRate(48000)
						}).
						WithNoVideo()
				}).
				AddProfile(func(x *ffmpeg.HlsArgument) {
					x.SetHlsTime(4).
						SetHlsPlaylistType("vod").
						SetHlsSegmentFilename("hls\\1080\\meow_%03d.ts", true).
						SetManifestOutput("hls\\1080\\meow.m3u8", true).
						WithVideoSettings(func(v *video.Video) {
							v.SetCodec(video.H264).
								SetBitrate(5000).
								SetMaxRate(5350).
								SetBufSize(7500)
						}).
						WithVideoFilters(func(v *video.Filter) {
							v.SetScaleCustom("-2:1080")
						}).
						WithNoAudio()
				}).
				AddProfile(func(x *ffmpeg.HlsArgument) {
					x.SetHlsTime(4).
						SetHlsPlaylistType("vod").
						SetHlsSegmentFilename("hls\\720\\meow_%03d.ts", true).
						SetManifestOutput("hls\\720\\meow.m3u8", true).
						WithVideoSettings(func(v *video.Video) {
							v.SetCodec(video.H264).
								SetBitrate(2800).
								SetMaxRate(2996).
								SetBufSize(4200)
						}).
						WithVideoFilters(func(vf *video.Filter) {
							vf.SetScaleCustom("-2:720")
						}).
						WithNoAudio()
				}).
				AddProfile(func(x *ffmpeg.HlsArgument) {
					x.SetHlsTime(4).
						SetHlsPlaylistType("vod").
						SetHlsSegmentFilename("hls\\360\\meow_%03d.ts", true).
						SetManifestOutput("hls\\360\\meow.m3u8", true).
						WithVideoSettings(func(v *video.Video) {
							v.SetCodec(video.H264).
								SetBitrate(800).
								SetMaxRate(856).
								SetBufSize(1200)
						}).
						WithVideoFilters(func(vf *video.Filter) {
							vf.SetScaleCustom("-2:360")
						}).
						WithNoAudio()
				})
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

	if err := executor.Execute(hls, ctx); err != nil {
		log.Fatalf("Error executing FFmpeg command: %v", err)
	}
}
