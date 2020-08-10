package refreshtoken

import (
	"github.com/Jamshid90/go-clean-architecture/pkg/errors"
	"time"
	"context"
	"github.com/Jamshid90/go-clean-architecture/pkg/entity"
)

type refreshTokenUsecase struct {
	refreshTokenRepo entity.RefreshTokenRepository
	contextTimeout time.Duration
}
// New refresh token usecase
func NewRefreshTokenUsecase(repo entity.RefreshTokenRepository, timeout time.Duration) refreshTokenUsecase  {
	return refreshTokenUsecase{
		refreshTokenRepo: repo,
		contextTimeout: timeout,
	}
}

// Before Store
func (r *refreshTokenUsecase) BeforeStore(m *entity.RefreshToken) error {
	m.CreatedAt = time.Now().UTC()
	return nil
}

// Store
func (r *refreshTokenUsecase) Store(ctx context.Context, m *entity.RefreshToken) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	if err := r.BeforeStore(m); err != nil {
		return err
	}

	return r.refreshTokenRepo.Store(ctx, m)
}

// Delete
func (r *refreshTokenUsecase) Delete(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	existedToken, err := r.refreshTokenRepo.Find(ctx, token)
	if err != nil {
		return err
	}

	if *existedToken == (entity.RefreshToken{}) {
		return errors.NewErrNotFound("user")
	}

	return r.refreshTokenRepo.Delete(ctx, token)
}

// Delete by user id
func (r *refreshTokenUsecase) DeleteByUserId(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	return r.refreshTokenRepo.DeleteByUserId(ctx, id)
}

// Find
func (r *refreshTokenUsecase) Find(ctx context.Context, token string) (*entity.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	return r.refreshTokenRepo.Find(ctx, token)
}
