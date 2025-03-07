package storagemanager

import (
	"github.com/google/uuid"

	"users/oas"
)

type StorageManager interface {
	// Create user and returns its uuid. If user with that login exists, returns error.
	CreateUser(request *oas.CreateUserRequest) (uuid.UUID, error)
	// If user with 'login' exists, returns its uuid. Else returns error.
	GetUserId(login oas.LoginString) (uuid.UUID, error)
	// If user with uuid exists, returns its profile. Else returns error.
	GetProfile(user oas.UserId) (oas.ProfileInfo, error)
	// If user with uuid exists, returns its password. Else returns error.
	GetPassword(user oas.UserId) (oas.PasswordString, error)
	// If user with uuid exists, returns its friends. Else returns error.
	GetFriends(user oas.UserId) (oas.FriendList, error)
	// Updates user profile. If there is no such user, returns error.
	UpdateProfile(user oas.UserId, update *oas.ProfileUpdate) error
	// Updates friend state. If there is no such user, or update is incorrect, returns error.
	UpdateFriend(user oas.UserId, update *oas.FriendModify) error
	// Make root user
	MakeRootUser(login oas.LoginString, password oas.PasswordString) error
}
