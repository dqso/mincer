package fonts

import (
	_ "embed"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
)

var (
	//go:embed LiberationMono-Regular.ttf
	liberationMonoRegular []byte

	//go:embed LiberationMono-Bold.ttf
	liberationMonoBold []byte

	//go:embed LiberationSans-Regular.ttf
	liberationSansRegular []byte

	//go:embed LiberationSans-Bold.ttf
	liberationSansBold []byte
)

var (
	Normal font.Face
	Bold   font.Face
)

const dpi = 72

func init() {
	tt, err := opentype.Parse(liberationSansRegular)
	if err != nil {
		log.Fatal(err)
	}
	Normal, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

	tt, err = opentype.Parse(liberationSansBold)
	if err != nil {
		log.Fatal(err)
	}
	Bold, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}
}
