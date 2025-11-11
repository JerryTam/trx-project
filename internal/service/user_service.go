package service

import (
	"context"
	"errors"
	"fmt"
	"trx-project/internal/model"
	"trx-project/internal/repository"
	"trx-project/pkg/jwt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) (*model.User, string, error)
	Login(ctx context.Context, username, password string) (*model.User, string, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
}

type userService struct {
	repo      repository.UserRepository
	redis     *redis.Client
	logger    *zap.Logger
	jwtConfig jwt.Config
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, redis *redis.Client, logger *zap.Logger, jwtConfig jwt.Config) UserService {
	return &userService{
		repo:      repo,
		redis:     redis,
		logger:    logger,
		jwtConfig: jwtConfig,
	}
}

func (s *userService) Register(ctx context.Context, username, email, password string) (*model.User, string, error) {
	// Check if user already exists
	if _, err := s.repo.GetByUsername(ctx, username); err == nil {
		return nil, "", errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to check username", zap.Error(err))
		return nil, "", err
	}

	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, "", errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to check email", zap.Error(err))
		return nil, "", err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, "", err
	}

	// Create user
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Status:   1,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, "", err
	}

	// Generate JWT Token
	token, err := jwt.GenerateToken(user.ID, user.Username, "user", s.jwtConfig)
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return nil, "", err
	}

	s.logger.Info("User registered successfully", zap.String("username", username))
	return user, token, nil
}

func (s *userService) Login(ctx context.Context, username, password string) (*model.User, string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid username or password")
		}
		s.logger.Error("Failed to get user", zap.Error(err))
		return nil, "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// Check if user is active
	if user.Status != 1 {
		return nil, "", errors.New("user account is inactive")
	}

	// Generate JWT Token
	token, err := jwt.GenerateToken(user.ID, user.Username, "user", s.jwtConfig)
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return nil, "", err
	}

	s.logger.Info("User logged in successfully", zap.String("username", username))
	return user, token, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%d", id)
	// In production, implement proper cache logic here

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		s.logger.Error("Failed to get user", zap.Error(err))
		return nil, err
	}

	// Cache the result
	_ = cacheKey // TODO: Implement caching

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	_ = cacheKey // TODO: Implement cache invalidation

	s.logger.Info("User updated successfully", zap.Uint("user_id", user.ID))
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete user", zap.Error(err))
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = cacheKey // TODO: Implement cache invalidation

	s.logger.Info("User deleted successfully", zap.Uint("user_id", id))
	return nil
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, total, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		s.logger.Error("Failed to list users", zap.Error(err))
		return nil, 0, err
	}

	return users, total, nil
}
