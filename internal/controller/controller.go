package controller

import (
	"context"
	"io"
)

type Controller struct {
	jwtService  jwtService
	exifService exifService
}

func New(jwtService jwtService, exifService exifService) *Controller {
	return &Controller{
		jwtService:  jwtService,
		exifService: exifService,
	}
}

func (c *Controller) GetHiddenImage(ctx context.Context, w io.Writer) error {
	token, err := c.jwtService.GetToken(ctx)
	if err != nil {
		return err
	}

	return c.exifService.EmbedAndWrite(ctx, token, w)
}
