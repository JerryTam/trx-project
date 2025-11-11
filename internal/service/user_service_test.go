package service

import (
	"context"
	"errors"
	"testing"
	"trx-project/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
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
	// Setup
	mockRepo := new(MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, nil, logger)

	ctx := context.Background()
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	// Test case 1: Successful registration
	t.Run("Successful registration", func(t *testing.T) {
		mockRepo.On("GetByUsername", ctx, username).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("GetByEmail", ctx, email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil).Once()

		user, err := service.Register(ctx, username, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, email, user.Email)
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Username already exists
	t.Run("Username already exists", func(t *testing.T) {
		existingUser := &model.User{Username: username}
		mockRepo.On("GetByUsername", ctx, username).Return(existingUser, nil).Once()

		user, err := service.Register(ctx, username, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "username already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Email already exists
	t.Run("Email already exists", func(t *testing.T) {
		existingUser := &model.User{Email: email}
		mockRepo.On("GetByUsername", ctx, username).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("GetByEmail", ctx, email).Return(existingUser, nil).Once()

		user, err := service.Register(ctx, username, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	// Test case 4: Database error
	t.Run("Database error", func(t *testing.T) {
		dbError := errors.New("database error")
		mockRepo.On("GetByUsername", ctx, username).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("GetByEmail", ctx, email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(dbError).Once()

		user, err := service.Register(ctx, username, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, dbError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, nil, logger)

	ctx := context.Background()
	username := "testuser"
	password := "password123"

	// Test case 1: User not found
	t.Run("User not found", func(t *testing.T) {
		mockRepo.On("GetByUsername", ctx, username).Return(nil, gorm.ErrRecordNotFound).Once()

		user, err := service.Login(ctx, username, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invalid username or password", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, nil, logger)

	ctx := context.Background()
	userID := uint(1)

	// Test case 1: Successful retrieval
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

	// Test case 2: User not found
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
	// Setup
	mockRepo := new(MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, nil, logger)

	ctx := context.Background()

	// Test case 1: Successful list
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

	// Test case 2: Invalid page
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

