package audio

type AudioCodec string

const (
	AAC    AudioCodec = "aac"
	MP3    AudioCodec = "libmp3lame"
	Opus   AudioCodec = "libopus"
	Vorbis AudioCodec = "libvorbis"
)
