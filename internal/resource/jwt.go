package resource

import (
	"context"
	"errors"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/golang-jwt/jwt/v5"
)

var (
	defaultSigningMethod = jwt.SigningMethodHS512
)

type JwtResource struct {
	cfg entities.Config
}

func NewJwtResource(cfg entities.Config) *JwtResource {
	return &JwtResource{
		cfg: cfg,
	}
}

func (j *JwtResource) GenerateToken(ctx context.Context, payload entities.JwtPayload) (string, error) {
	token := jwt.NewWithClaims(defaultSigningMethod, jwt.MapClaims{
		"iat":      payload.IssuedAt.Unix(),
		"token_id": payload.TokenId,
	})

	return token.SignedString([]byte(j.cfg.JwtKey))
}

func (j *JwtResource) ParseToken(ctx context.Context, token string) (entities.JwtPayload, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Error token parsing")
		}
		return j.cfg.JwtKey, nil
	})
	if err != nil {
		return entities.JwtPayload{}, err
	}

	if jwtToken == nil {
		return entities.JwtPayload{}, errors.New("Invalid token")
	}

	issuedAt, err := jwtToken.Claims.GetIssuedAt()
	if err != nil {
		return entities.JwtPayload{}, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return entities.JwtPayload{}, errors.New("Error claims")
	}

	tokenId, ok := claims["token_id"]
	if !ok {
		return entities.JwtPayload{}, errors.New("token id not found")
	}

	return entities.JwtPayload{
		IssuedAt: issuedAt.Time,
		TokenId:  tokenId.(string),
	}, nil
}
