package video

type VideoCodec string

const (
	H264 VideoCodec = "libx264"
	H265 VideoCodec = "libx265"
	VP8  VideoCodec = "libvpx"
	VP9  VideoCodec = "libvpx-vp9"
)
