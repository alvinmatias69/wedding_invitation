package service

import (
	"context"
	"io"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

type ExifService struct {
	cfg entities.Config
}

func NewExifService(cfg entities.Config) *ExifService {
	return &ExifService{
		cfg: cfg,
	}
}

func (e *ExifService) EmbedAndWrite(ctx context.Context, token string, w io.Writer) error {
	mediaCtx, err := jpegstructure.NewJpegMediaParser().ParseFile(e.cfg.HiddenImageFilePath)
	if err != nil {
		return err
	}

	segmentList := mediaCtx.(*jpegstructure.SegmentList)
	builder, err := segmentList.ConstructExifBuilder()
	if err != nil {
		return err
	}

	ifd0Ib, err := exif.GetOrCreateIbFromRootIb(builder, e.cfg.FqIfdPath)
	if err != nil {
		return err
	}

	err = ifd0Ib.SetStandardWithName(e.cfg.HiddenImageUrlTag, e.cfg.HiddenImageUrlValue)
	if err != nil {
		return err
	}

	err = ifd0Ib.SetStandardWithName(e.cfg.HiddenImageTokenTag, token)
	if err != nil {
		return err
	}

	err = segmentList.SetExif(builder)
	if err != nil {
		return err
	}

	return segmentList.Write(w)
}
