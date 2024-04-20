package user_and_post_service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
	"github.com/khailequang334/social_network/internal/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (uaps *UserAndPostService) CreateUser(ctx context.Context, request *user_and_post.UserDetailInfo) (*user_and_post.UserResult, error) {
	salt := GenerateSalt(16)
	hashedPassword, err := HashPassword(request.GetUserPassword(), salt)
	if err != nil {
		return nil, err
	}

	newUser := types.User{
		HashedPassword: hashedPassword,
		Salt:           string(salt),
		FirstName:      request.GetFirstName(),
		LastName:       request.GetLastName(),
		DateOfBirth:    request.Dob.AsTime(),
		Email:          request.GetEmail(),
		UserName:       request.GetUserName(),
	}
	// add new user in DB
	uaps.DB.Create(&newUser)

	return &user_and_post.UserResult{
		Status: user_and_post.UserResult_OK,
		Info: &user_and_post.UserDetailInfo{
			UserId:   int64(newUser.ID),
			UserName: newUser.UserName,
		},
	}, nil
}

func GenerateSalt(length int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt := make([]byte, length)

	// generate salt randomly
	for i := 0; i < length; i++ {
		salt[i] = alphanum[random.Intn(len(alphanum))]
	}

	return string(salt)
}

func HashPassword(password string, salt string) (string, error) {
	// Add salt
	passwordWithSalt := []byte(password + salt)

	// Generate the bcrypt hash
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordWithSalt, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// EditUser edit user request by looking up user id in mysql database and update it
func (uaps *UserAndPostService) EditUser(ctx context.Context, request *user_and_post.EditUserRequest) (*user_and_post.EditUserResponse, error) {
	var user types.User
	uaps.DB.Where(&types.User{ID: uint(request.UserId)}).First(&user)

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	// update FirstName
	if request.FirstName != nil {
		user.FirstName = request.GetFirstName()
	}
	// update LastName
	if request.LastName != nil {
		user.LastName = request.GetLastName()
	}
	// update Password
	if request.UserPassword != nil {
		salt := GenerateSalt(16)
		hashedPassword, err := HashPassword(request.GetUserPassword(), salt)
		if err != nil {
			return nil, err
		}
		user.HashedPassword = hashedPassword
		user.Salt = string(salt)
	}
	// update DOB
	if request.Dob != nil {
		user.DateOfBirth = request.Dob.AsTime()
	}

	uaps.DB.Save(&user)

	return &user_and_post.EditUserResponse{
		UserId: int64(user.ID),
	}, nil
}

func (uaps *UserAndPostService) AuthenticateUser(ctx context.Context, request *user_and_post.AuthenticateUserRequest) (*user_and_post.AuthenticateUserResponse, error) {
	// get user detail infor from DB
	var user types.User
	result := uaps.DB.Where(&types.User{UserName: request.GetUserName()}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &user_and_post.AuthenticateUserResponse{
			Status: user_and_post.AuthenticateUserResponse_USER_NOT_FOUND,
		}, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	// add salt to password to compare
	passwordWithSalt := []byte(request.UserPassword + user.Salt)
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), passwordWithSalt)
	if err != nil {
		return &user_and_post.AuthenticateUserResponse{
			Status: user_and_post.AuthenticateUserResponse_WRONG_PASSWORD,
		}, nil
	}

	return &user_and_post.AuthenticateUserResponse{
		Status: user_and_post.AuthenticateUserResponse_OK,
	}, nil
}
