package internal

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"
)

type Band string

const (
	HeavyRainfall           Band = "hrp"
	Infrared                Band = "b13"
	Visible                 Band = "b03"
	WaterVapor              Band = "b08"
	ShortWaveIR             Band = "b07"
	DayMicrophysicsRGB      Band = "dms"
	NightMicrophysicsRGB    Band = "ngt"
	DustRGB                 Band = "dst"
	AirmassRGB              Band = "arm"
	DaySnowFogRGB           Band = "dsl"
	NaturalColorRGB         Band = "dnc"
	TrueColorRGBEnhanced    Band = "tre"
	TrueColorReproduction   Band = "trm"
	DayConvectiveStormRGB   Band = "cve"
	Sandwich                Band = "snd"
	VisibleAndInfrared      Band = "vis"
	VisibleAndInfraredNight Band = "irv"
)

type ImageLink struct {
	Band Band
	Link string
}

func (i ImageLink) FileName() string {
	u, _ := url.Parse(i.Link)

	return path.Base(u.Path)
}

func (i ImageLink) SequenceNum() (string, error) {
	fileName := strings.TrimSuffix(i.FileName(), ".jpg")
	parts := strings.Split(fileName, "_")
	if len(parts) != 3 {
		return "", fmt.Errorf("unknown file name: %s", i.FileName())
	}

	parsedTime, err := time.Parse("1504", parts[2])
	if err != nil {
		return "", fmt.Errorf("error unsupported file name: %s", i.FileName())
	}

	newTime := parsedTime.Add(10 * time.Minute)

	return newTime.Format("15") + newTime.Format("04"), nil
}

func NewImageLink(b Band, link string) *ImageLink {
	return &ImageLink{Band: b, Link: link}
}
