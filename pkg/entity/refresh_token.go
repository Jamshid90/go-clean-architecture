package entity

import (
	"context"
	"time"
)

type RefreshToken struct {
	UserID    string
	Token     string
	CreatedAt time.Time
}

type RefreshTokenUsecase interface {
	Store(ctx context.Context, user *RefreshToken) error
	Delete(ctx context.Context, token string) error
	DeleteByUserId(ctx context.Context, id string) error
	Find(ctx context.Context, token string) (*RefreshToken, error)
}

type RefreshTokenRepository interface {
	Store(ctx context.Context, user *RefreshToken) error
	Delete(ctx context.Context, token string) error
	DeleteByUserId(ctx context.Context, id string) error
	Find(ctx context.Context, token string) (*RefreshToken, error)
}
