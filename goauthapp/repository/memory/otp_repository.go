package memory

import (
	"context"
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/service/goauth"
	"github.com/patrickmn/go-cache"
	"time"
)

type OTPCacheConfig struct {
	Expiration      time.Duration `koanf:"expiration"`
	CleanupInterval time.Duration `koanf:"cleanup_interval"`
}
type GoCacheOTPRepository struct {
	c      *cache.Cache
	config OTPCacheConfig
}

func NewGoCacheOTPRepository(config OTPCacheConfig) goauth.OTPRepository {
	c := cache.New(config.Expiration, config.CleanupInterval)

	return &GoCacheOTPRepository{c: c, config: config}
}

func (r *GoCacheOTPRepository) Save(ctx context.Context, phoneNumber, otp string) error {
	r.c.Set(phoneNumber, otp, cache.DefaultExpiration)
	return nil
}

func (r *GoCacheOTPRepository) Find(ctx context.Context, phoneNumber string) (string, error) {
	if otp, found := r.c.Get(phoneNumber); found {
		return otp.(string), nil
	}
	return "", goauth.ErrOTPNotFound
}
