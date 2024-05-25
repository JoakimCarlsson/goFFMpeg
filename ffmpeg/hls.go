package ffmpeg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type HlsArgument struct {
	Argument
	HlsArguments []string
}

func (a *HlsArgument) AddProfile(configure func(*HlsArgument)) *HlsArgument {
	profile := &HlsArgument{}
	configure(profile)
	profile.Arguments = append(profile.Arguments, "-f hls", strings.Join(profile.HlsArguments, " "))
	a.Arguments = append(a.Arguments, strings.Join(profile.Arguments, " "))
	return a
}

func (a *HlsArgument) SetManifestOutput(filePath string, createDirectory bool) *HlsArgument {
	if !createDirectory {
		return a.addArgument(filePath)
	}
	directory := filepath.Dir(filePath)
	if _, err := os.Stat(directory); os.IsNotExist(err) && directory != "" {
		os.MkdirAll(directory, os.ModePerm)
	}
	return a.addArgument(fmt.Sprintf("\"%s\"", filePath))
}

func (a *HlsArgument) SetSegmentTime(segmentTime int) *HlsArgument {
	return a.addArgument(fmt.Sprintf("-segment_time %d", segmentTime))
}

func (a *HlsArgument) SetSegmentListSize(segmentListSize int) *HlsArgument {
	return a.addArgument(fmt.Sprintf("-segment_list_size %d", segmentListSize))
}

func (a *HlsArgument) SetSegmentFormat(segmentFormat string) *HlsArgument {
	return a.addArgument(fmt.Sprintf("-segment_format %s", segmentFormat))
}

func (a *HlsArgument) SetHlsTime(hlsTime int) *HlsArgument {
	return a.addArgument(fmt.Sprintf("-hls_time %d", hlsTime))
}

func (a *HlsArgument) SetHlsAllowCache(allowCache bool) *HlsArgument {
	return a.addArgument(fmt.Sprintf("-hls_allow_cache %d", boolToInt(allowCache)))
}

func (a *HlsArgument) SetHlsPlaylistType(playlistType string) *HlsArgument {
	return a.addArgument(fmt.Sprintf("-hls_playlist_type %s", strings.ToLower(playlistType)))
}

func (a *HlsArgument) SetHlsSegmentFilename(segmentFilename string, createDirectory bool) *HlsArgument {
	if !createDirectory {
		return a.addArgument(fmt.Sprintf("-hls_segment_filename %s", segmentFilename))
	}
	directory := filepath.Dir(segmentFilename)
	if _, err := os.Stat(directory); os.IsNotExist(err) && directory != "" {
		os.MkdirAll(directory, os.ModePerm)
	}
	return a.addArgument(fmt.Sprintf("-hls_segment_filename %s", segmentFilename))
}

func (a *HlsArgument) addArgument(argument string) *HlsArgument {
	if !strings.HasSuffix(argument, " ") {
		argument += " "
	}
	a.HlsArguments = append(a.HlsArguments, argument)
	return a
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
