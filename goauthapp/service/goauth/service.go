package goauth

import (
	"context"
	"errors"
	"fmt"
	"github.com/mohammadrezajavid-lab/goauth/pkg/logger"
	"github.com/mohammadrezajavid-lab/goauth/pkg/otpgenerator"
	"github.com/mohammadrezajavid-lab/goauth/pkg/token"
	"log/slog"
	"math"
)

// --- Repository Interfaces (Contracts) ---

type ListUsersParams struct {
	Page     int
	PageSize int
	Search   string
}

type UserRepository interface {
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id int64) (*User, error)
	List(ctx context.Context, params ListUsersParams) ([]User, int, error)
}

type OTPRepository interface {
	Save(ctx context.Context, phoneNumber, otp string) error
	Find(ctx context.Context, phoneNumber string) (string, error)
}

// --- Service Implementation ---

// Service handles the business logic for authentication.
type Service struct {
	userRepo   UserRepository
	otpRepo    OTPRepository
	tokenMaker token.Maker
}

// NewService creates a new auth service.
func NewService(userRepo UserRepository, otpRepo OTPRepository, tokenMaker token.Maker) Service {
	return Service{
		userRepo:   userRepo,
		otpRepo:    otpRepo,
		tokenMaker: tokenMaker,
	}
}

// --- OTP and Login Methods ---

func (s *Service) GenerateOTP(ctx context.Context, request *GenerateOTPRequest) (*GenerateOTPResponse, error) {
	log := logger.L().With(slog.String("phone_number", request.PhoneNumber))
	log.Info("OTP generation requested")

	otp, err := otpgenerator.GenerateOTP(6)
	if err != nil {
		log.Error("Could not generate OTP", slog.String("error", err.Error()))
		return nil, fmt.Errorf("could not generate OTP: %w", err)
	}

	err = s.otpRepo.Save(ctx, request.PhoneNumber, otp)
	if err != nil {
		log.Error("Could not save OTP", slog.String("error", err.Error()))
		return nil, fmt.Errorf("could not save OTP: %w", err)
	}

	log.Info("OTP code generated successfully", slog.String("otp_code", otp))
	return &GenerateOTPResponse{OTP: otp}, nil
}

func (s *Service) VerifyAndLogin(ctx context.Context, request *VerifyOTPRequest) (*VerifyOTPResponse, error) {
	log := logger.L().With(slog.String("phone_number", request.PhoneNumber))
	log.Info("Login verification requested")

	storedOTP, err := s.otpRepo.Find(ctx, request.PhoneNumber)
	if err != nil {
		log.Warn("Failed to find OTP", slog.String("error", err.Error()))
		return nil, err
	}

	if storedOTP != request.OTP {
		log.Warn("Invalid OTP code provided")
		return nil, ErrInvalidOTPCode
	}

	log.Info("OTP verified successfully, checking user status")

	user, err := s.userRepo.FindByPhoneNumber(ctx, request.PhoneNumber)
	isNewUser := false

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			log.Info("User not found, proceeding with registration")
			isNewUser = true
			newUser := &User{PhoneNumber: request.PhoneNumber}
			if createErr := s.userRepo.Create(ctx, newUser); createErr != nil {
				log.Error("Failed to register new user", slog.String("error", createErr.Error()))
				return nil, createErr
			}
			user = newUser
		} else {
			log.Error("Database error while finding user", slog.String("error", err.Error()))
			return nil, errors.Join(ErrInternalService, err)
		}
	} else {
		log.Debug("Existing user found", slog.String("phone_number", user.PhoneNumber))
	}

	accessToken, err := s.tokenMaker.CreateToken(user.ID)
	if err != nil {
		log.Error("Failed to create access token", slog.String("error", err.Error()))
		return nil, ErrInternalService
	}

	log.Debug("JWT token generated for user", slog.Int64("user_id", user.ID))

	return &VerifyOTPResponse{
		Token: accessToken,
		IsNew: isNewUser,
	}, nil
}

// --- User Management Methods ---

func (s *Service) GetUser(ctx context.Context, request *GetUserRequest) (*User, error) {
	log := logger.L().With(slog.Int64("user_id", request.ID))
	log.Info("GetUser service method called")

	user, err := s.userRepo.FindByID(ctx, request.ID)
	if err != nil {
		log.Warn("Failed to get user by ID", slog.String("error", err.Error()))
		return nil, err
	}

	return user, nil
}

func (s *Service) ListUsers(ctx context.Context, request *ListUsersRequest) (*ListUsersResponse, error) {
	log := logger.L().With(slog.Any("request", request))
	log.Info("ListUsers service method called")

	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	params := ListUsersParams{
		Page:     request.Page,
		PageSize: request.PageSize,
		Search:   request.Search,
	}

	users, totalRecords, err := s.userRepo.List(ctx, params)
	if err != nil {
		log.Error("Failed to list users from repository", slog.String("error", err.Error()))
		return nil, ErrInternalService
	}

	totalPages := 0
	if totalRecords > 0 {
		totalPages = int(math.Ceil(float64(totalRecords) / float64(request.PageSize)))
	}

	metadata := PaginationMetadata{
		CurrentPage:  request.Page,
		PageSize:     request.PageSize,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
	}

	return &ListUsersResponse{
		Users:    users,
		Metadata: metadata,
	}, nil
}
