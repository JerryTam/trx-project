package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"trx-project/internal/model"
	"trx-project/pkg/jwt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MockUserRepository 是 UserRepository 的 mock 实现
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func TestUserService_Register(t *testing.T) {
	// 配置
	mockRepo := new(MockUserRepository)
	logger, _ := zap.NewDevelopment()
	jwtConfig := jwt.Config{
		Secret:     "test-secret",
		Issuer:     "test",
		ExpireTime: 24 * time.Hour,
	}
	service := NewUserService(mockRepo, nil, logger, jwtConfig)

	ctx := context.Background()
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	// 测试用例 1: 成功注册
	t.Run("Successful registration", func(t *testing.T) {
		mockRepo.On("GetByUsername", ctx, username).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("GetByEmail", ctx, email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil).Once()

		user, token, err := service.Register(ctx, username, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, token)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, email, user.Email)
		mockRepo.AssertExpectations(t)
	})

	// 测试用例 2: 用户名已存在
	t.Run("Username already exists", func(t *testing.T) {
		existingUser := &model.User{Username: username}
		mockRepo.On("GetByUsername", ctx, username).Return(existingUser, nil).Once()

		user, token, err := service.Register(ctx, username, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
		assert.Equal(t, "username already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	// 测试用例 3: 邮箱已存在
	t.Run("Email already exists", func(t *testing.T) {
		existingUser := &model.User{Email: email}
		mockRepo.On("GetByUsername", ctx, username).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("GetByEmail", ctx, email).Return(existingUser, nil).Once()

		user, token, err := service.Register(ctx, username, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
		assert.Equal(t, "email already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	// 测试用例 4: 数据库错误
	t.Run("Database error", func(t *testing.T) {
		dbError := errors.New("database error")
		mockRepo.On("GetByUsername", ctx, username).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("GetByEmail", ctx, email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(dbError).Once()

		user, token, err := service.Register(ctx, username, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
		assert.Equal(t, dbError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	// 配置
	mockRepo := new(MockUserRepository)
	logger, _ := zap.NewDevelopment()
	jwtConfig := jwt.Config{
		Secret:     "test-secret",
		Issuer:     "test",
		ExpireTime: 24 * time.Hour,
	}
	service := NewUserService(mockRepo, nil, logger, jwtConfig)

	ctx := context.Background()
	username := "testuser"
	password := "password123"

	// 测试用例 1: 用户不存在
	t.Run("User not found", func(t *testing.T) {
		mockRepo.On("GetByUsername", ctx, username).Return(nil, gorm.ErrRecordNotFound).Once()

		user, token, err := service.Login(ctx, username, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
		assert.Equal(t, "invalid username or password", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	// 配置
	mockRepo := new(MockUserRepository)
	logger, _ := zap.NewDevelopment()
	jwtConfig := jwt.Config{
		Secret:     "test-secret",
		Issuer:     "test",
		ExpireTime: 24 * time.Hour,
	}
	service := NewUserService(mockRepo, nil, logger, jwtConfig)

	ctx := context.Background()
	userID := uint(1)

	// 测试用例 1: 成功获取
	t.Run("Successful retrieval", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
		}
		mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil).Once()

		user, err := service.GetUserByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		mockRepo.AssertExpectations(t)
	})

	// 测试用例 2: 用户不存在
	t.Run("User not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, userID).Return(nil, gorm.ErrRecordNotFound).Once()

		user, err := service.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_ListUsers(t *testing.T) {
	// 配置
	mockRepo := new(MockUserRepository)
	logger, _ := zap.NewDevelopment()
	jwtConfig := jwt.Config{
		Secret:     "test-secret",
		Issuer:     "test",
		ExpireTime: 24 * time.Hour,
	}
	service := NewUserService(mockRepo, nil, logger, jwtConfig)

	ctx := context.Background()

	// 测试用例 1: 成功获取列表
	t.Run("Successful list", func(t *testing.T) {
		expectedUsers := []*model.User{
			{ID: 1, Username: "user1"},
			{ID: 2, Username: "user2"},
		}
		mockRepo.On("List", ctx, 0, 10).Return(expectedUsers, int64(2), nil).Once()

		users, total, err := service.ListUsers(ctx, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(users))
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	// 测试用例 2: 无效页码 - 应默认为 1
	t.Run("Invalid page - should default to 1", func(t *testing.T) {
		expectedUsers := []*model.User{}
		mockRepo.On("List", ctx, 0, 10).Return(expectedUsers, int64(0), nil).Once()

		users, total, err := service.ListUsers(ctx, 0, 10)

		assert.NoError(t, err)
		assert.Equal(t, 0, len(users))
		assert.Equal(t, int64(0), total)
		mockRepo.AssertExpectations(t)
	})
}
