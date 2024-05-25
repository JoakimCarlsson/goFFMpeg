package audio

import "fmt"

type Audio struct {
	Arguments []string
}

func (a *Audio) SetCodec(codec AudioCodec) *Audio {
	return a.addArgument(fmt.Sprintf("-c:a %s", codec))
}

func (a *Audio) SetBitrate(bitrate int) *Audio {
	return a.addArgument(fmt.Sprintf("-b:a %dk", bitrate))
}

func (a *Audio) SetChannels(channels Channels) *Audio {
	a.Arguments = append(a.Arguments, fmt.Sprintf("-ac %s", channels))
	return a
}

func (a *Audio) SetSamplingRate(rate int) *Audio {
	return a.addArgument(fmt.Sprintf("-ar %d", rate))
}

func (a *Audio) addArgument(argument string) *Audio {
	a.Arguments = append(a.Arguments, argument)
	return a
}
