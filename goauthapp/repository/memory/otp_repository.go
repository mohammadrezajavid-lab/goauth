package memory

import (
	"context"
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/service/goauth"
	"github.com/mohammadrezajavid-lab/goauth/pkg/logger"
	"log/slog"
	"sync"
	"time"
)

type otpItem struct {
	code      string
	expiresAt time.Time
}

type OtpRepository struct {
	store map[string]otpItem
	mu    sync.RWMutex
}

func NewOtpRepository() goauth.OTPRepository {
	repo := &OtpRepository{
		store: make(map[string]otpItem),
	}
	go repo.cleanupExpiredOTPs()
	return repo
}

func (r *OtpRepository) Save(ctx context.Context, phoneNumber, otp string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[phoneNumber] = otpItem{
		code:      otp,
		expiresAt: time.Now().Add(2 * time.Minute),
	}
	logger.L().Debug("OTP saved in memory", slog.String("phone_number", phoneNumber))
	return nil
}

func (r *OtpRepository) Find(ctx context.Context, phoneNumber string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, found := r.store[phoneNumber]
	log := logger.L().With(slog.String("phone_number", phoneNumber))

	if !found {
		log.Warn("OTP not found in store")
		return "", goauth.ErrOTPNotFound
	}

	if time.Now().After(item.expiresAt) {
		r.mu.RUnlock() // Release read lock
		r.mu.Lock()    // Acquire write lock
		delete(r.store, phoneNumber)
		r.mu.Unlock() // Release write lock
		r.mu.RLock()  // Re-acquire read lock for defer

		log.Warn("OTP has expired")
		return "", goauth.ErrOTPNotFound
	}

	log.Debug("OTP found successfully")
	return item.code, nil
}

func (r *OtpRepository) cleanupExpiredOTPs() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		cleanedCount := 0
		for phone, item := range r.store {
			if time.Now().After(item.expiresAt) {
				delete(r.store, phone)
				cleanedCount++
			}
		}
		r.mu.Unlock()
		if cleanedCount > 0 {
			logger.L().Debug("Cleaned up expired OTPs", slog.Int("count", cleanedCount))
		}
	}
}
