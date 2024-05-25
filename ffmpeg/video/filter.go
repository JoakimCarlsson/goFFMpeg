package video

import "fmt"

type Filter struct {
	Arguments []string
}

func (vf *Filter) SetScale(width, height int) *Filter {
	return vf.addArgument(fmt.Sprintf("scale=%d:%d", width, height))
}

func (vf *Filter) SetScaleCustom(scale string) *Filter {
	return vf.addArgument(fmt.Sprintf("scale=%s", scale))
}

func (vf *Filter) SetScaleStr(scale string) *Filter {
	return vf.addArgument("scale=" + scale)
}

func (vf *Filter) SetFps(fps float64) *Filter {
	return vf.addArgument(fmt.Sprintf("fps=%f", fps))
}

func (vf *Filter) SetTranspose(transpose string) *Filter {
	return vf.addArgument("transpose=" + transpose)
}

func (vf *Filter) addArgument(argument string) *Filter {
	vf.Arguments = append(vf.Arguments, argument)
	return vf
}

type AudioFilter struct {
	Arguments []string
}

func (af *AudioFilter) addArgument(argument string) *AudioFilter {
	af.Arguments = append(af.Arguments, argument)
	return af
}

type MetaData struct {
	Arguments []string
}

func (md *MetaData) addArgument(argument string) *MetaData {
	md.Arguments = append(md.Arguments, argument)
	return md
}
