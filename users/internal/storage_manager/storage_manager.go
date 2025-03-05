package storagemanager

import (
	"users/oas"

	"github.com/google/uuid"
)

type StorageManager interface {
	// Create user and returns its uuid. If user with that login exists, returns error.
	CreateUser(request oas.CreateUserRequest) (uuid.UUID, error)
	// If user with 'login' exists, returns its uuid. Else returns error.
	GetUserId(login oas.LoginString) (uuid.UUID, error)
	// If user with uuid exists, returns its profile. Else returns error.
	GetProfile(uuid oas.UserId) (oas.ProfileInfo, error)
	// If user with uuid exists, returns its password. Else returns error.
	GetPassword(uuid oas.UserId) (oas.PasswordString, error)
	// If user with uuid exists, returns its friends. Else returns error.
	GetFriends(uuid oas.UserId) (oas.FriendList, error)
	// Updates user profile. If there is no such user, or update is incorrect, returns error.
	UpdateProfile(uuid oas.UserId, update oas.ProfileUpdate) error
	// Updates friend state. If there is no such user, or update is incorrect, returns error.
	UpdateFriend(uuidd oas.UserId, update oas.FriendModify) error
}
