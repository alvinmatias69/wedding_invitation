package controller

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/alvinmatias69/wedding_invitation/internal/constant"
	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/google/uuid"
)

var (
	steamTokenSuccessMsg  = "Congratulations, you've conquered the challenge. If it's before 12th of October 2024, I'd be delighted if you can come to my wedding!"
	steamTokenUnavailable = "Congratulations, you've conquered the challenge. Unfortunately, all prizes already redeemed. Hit me up and say the secret word 'fufufafa kaskus legend', our next meal will be on me. If it's before 12th of October 2024, I'd be delighted if you can come to my wedding!"
)

type Controller struct {
	cfg               entities.Config
	jwtResource       jwtResource
	exifResource      exifResource
	tokenRepository   tokenRepository
	messageRepository messageRepository
}

func New(cfg entities.Config, jwtResource jwtResource, exifResource exifResource, tokenRepository tokenRepository, messageRepository messageRepository) *Controller {
	return &Controller{
		cfg:               cfg,
		jwtResource:       jwtResource,
		exifResource:      exifResource,
		tokenRepository:   tokenRepository,
		messageRepository: messageRepository,
	}
}

func (c *Controller) GetHiddenImage(ctx context.Context, w io.Writer) error {
	token, err := c.jwtResource.GenerateToken(ctx, entities.JwtPayload{
		IssuedAt: time.Now(),
		TokenId:  uuid.NewString(),
	})
	if err != nil {
		return err
	}

	lambda, err := c.exifResource.Embed(ctx, map[string]interface{}{
		c.cfg.HiddenImageUrlTag:   c.cfg.SteamTokenPath,
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
			Message: steamTokenSuccessMsg,
		}, nil
	}

	if err != nil && !errors.Is(err, constant.ErrNotFound) {
		return entities.SteamTokenResponse{}, err
	}

	if jwtPayload.IssuedAt.Add(time.Minute * time.Duration(c.cfg.JwtExpiryMinute)).Before(time.Now()) {
		return entities.SteamTokenResponse{}, constant.ErrTokenExp
	}

	trx, err := c.tokenRepository.BeginTrx(ctx)
	if err != nil {
		return entities.SteamTokenResponse{}, err
	}

	defer trx.Rollback(ctx)

	tokenData, err = c.tokenRepository.FindOneUnclaimed(ctx, trx)
	if errors.Is(err, constant.ErrNotFound) {
		return entities.SteamTokenResponse{
			Message: steamTokenUnavailable,
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
		Message: steamTokenSuccessMsg,
	}, nil
}

func (c *Controller) GetMessages(ctx context.Context, page uint64) ([]entities.Message, error) {
	var offset = page * c.cfg.MessageLimit
	messages, err := c.messageRepository.Get(ctx, c.cfg.MessageLimit, offset)
	if errors.Is(err, constant.ErrNotFound) {
		return make([]entities.Message, 0), nil
	}

	return messages, err
}

func (c *Controller) PostMessage(ctx context.Context, msg entities.Message) error {
	return c.messageRepository.Insert(ctx, msg)
}
