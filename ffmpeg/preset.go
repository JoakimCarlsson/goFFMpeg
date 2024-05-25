package ffmpeg

type Preset string

const (
	Ultrafast Preset = "ultrafast"
	Superfast Preset = "superfast"
	Veryfast  Preset = "veryfast"
	Faster    Preset = "faster"
	Fast      Preset = "fast"
	Medium    Preset = "medium"
	Slow      Preset = "slow"
	Slower    Preset = "slower"
	Veryslow  Preset = "veryslow"
	Placebo   Preset = "placebo"
)
