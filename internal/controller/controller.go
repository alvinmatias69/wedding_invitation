package controller

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/google/uuid"
)

type Controller struct {
	cfg             entities.Config
	jwtResource     jwtResource
	exifResource    exifResource
	tokenRepository tokenRepository
}

func New(cfg entities.Config, jwtResource jwtResource, exifResource exifResource, tokenRepository tokenRepository) *Controller {
	return &Controller{
		cfg:             cfg,
		jwtResource:     jwtResource,
		exifResource:    exifResource,
		tokenRepository: tokenRepository,
	}
}

func (c *Controller) GetHiddenImage(ctx context.Context, w io.Writer) error {
	token, err := c.jwtResource.GenerateToken(ctx, map[string]interface{}{
		"iat":      time.Now().Unix(),
		"token_id": uuid.NewString(),
	})
	if err != nil {
		return err
	}

	lambda, err := c.exifResource.Embed(ctx, map[string]interface{}{
		c.cfg.HiddenImageUrlTag:   c.cfg.HiddenImageUrlValue,
		c.cfg.HiddenImageTokenTag: token,
	})
	if err != nil {
		return err
	}

	return lambda(w)
}

func (c *Controller) GetSteamToken(ctx context.Context, token string) (entities.SteamTokenResponse, error) {
	jwtPayload, err := c.jwtResource.ParseToken(ctx, token)
	if err != nil {
		return entities.SteamTokenResponse{}, err
	}

	tokenData, err := c.tokenRepository.GetByJwtToken(ctx, jwtPayload.TokenId)
	if err == nil {
		return entities.SteamTokenResponse{
			TokenId: tokenData.SteamToken,
			Message: "lorem ipsum",
		}, nil
	}

	if err != nil && !errors.Is(err, errors.New("not found")) {
		return entities.SteamTokenResponse{}, err
	}

	if jwtPayload.IssuedAt.Add(time.Millisecond * time.Duration(c.cfg.JwtExpiryMs)).After(time.Now()) {
		return entities.SteamTokenResponse{}, errors.New("token expired")
	}

	trx, err := c.tokenRepository.BeginTrx(ctx)
	if err != nil {
		return entities.SteamTokenResponse{}, err
	}

	defer trx.Rollback(ctx)

	tokenData, err = c.tokenRepository.FindOneUnclaimed(ctx, trx)
	if errors.Is(err, errors.New("not found")) {
		return entities.SteamTokenResponse{
			Message: "lorem ipsum too late",
		}, nil
	}

	if err != nil {
		return entities.SteamTokenResponse{}, err
	}

	err = c.tokenRepository.Claim(ctx, trx, tokenData.Id, jwtPayload.TokenId)
	if err != nil {
		return entities.SteamTokenResponse{}, err
	}

	trx.Commit(ctx)

	return entities.SteamTokenResponse{
		TokenId: tokenData.SteamToken,
		Message: "lorem ipsum",
	}, nil
}
