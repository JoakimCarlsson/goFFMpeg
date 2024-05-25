package video

import "fmt"

type Video struct {
	Arguments []string
}

func (v *Video) SetCodec(codec VideoCodec) *Video {
	return v.addArgument(fmt.Sprintf("-c:v %s", codec))
}

func (v *Video) SetCrf(crf int) *Video {
	return v.addArgument(fmt.Sprintf("-crf %d", crf))
}

func (v *Video) SetBitrate(bitrate int) *Video {
	return v.addArgument(fmt.Sprintf("-b:v %dk", bitrate))
}

func (v *Video) SetMaxRate(maxRate int) *Video {
	return v.addArgument(fmt.Sprintf("-maxrate %dk", maxRate))
}

func (v *Video) SetBufSize(bufSize int) *Video {
	return v.addArgument(fmt.Sprintf("-bufsize %dk", bufSize))
}

func (v *Video) addArgument(argument string) *Video {
	v.Arguments = append(v.Arguments, argument)
	return v
}
