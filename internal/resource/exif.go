package resource

import (
	"context"
	"io"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

type ExifResource struct {
	cfg entities.Config
}

func NewExifResource(cfg entities.Config) *ExifResource {
	return &ExifResource{
		cfg: cfg,
	}
}

func (e *ExifResource) Embed(ctx context.Context, payload map[string]interface{}) (func(io.Writer) error, error) {
	mediaCtx, err := jpegstructure.NewJpegMediaParser().ParseFile(e.cfg.HiddenImageFilePath)
	if err != nil {
		return nil, err
	}

	segmentList := mediaCtx.(*jpegstructure.SegmentList)
	builder, err := segmentList.ConstructExifBuilder()
	if err != nil {
		return nil, err
	}

	ifd0Ib, err := exif.GetOrCreateIbFromRootIb(builder, e.cfg.FqIfdPath)
	if err != nil {
		return nil, err
	}

	for key, val := range payload {
		err = ifd0Ib.SetStandardWithName(key, val)
		if err != nil {
			return nil, err
		}
	}

	err = segmentList.SetExif(builder)
	if err != nil {
		return nil, err
	}

	return segmentList.Write, nil
}
