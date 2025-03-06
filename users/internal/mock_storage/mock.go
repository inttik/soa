package mockstorage

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"users/oas"
)

type mockStorage struct {
	ui userInfo
	up userProfile
	f  friends
}

func NewMockStorage() (mockStorage, error) {
	return mockStorage{
		ui: userInfo{
			data:  make(map[uuid.UUID]userInfoData),
			index: make(map[oas.LoginString]uuid.UUID),
		},
		up: userProfile{data: make(map[uuid.UUID]userProfileData)},
		f:  friends{data: make(map[uuid.UUID]map[uuid.UUID]friendsData)},
	}, nil
}

func (s *mockStorage) MakeRootUser(login oas.LoginString, password oas.PasswordString) error {
	createRequest := oas.CreateUserRequest{
		Login:    login,
		Password: password,
		Email:    oas.EmailString(string(login) + "@root"),
		Root:     oas.NewOptRootFlag(true),
	}
	_, err := s.CreateUser(&createRequest)
	return err
}

func (ms *mockStorage) CreateUser(request *oas.CreateUserRequest) (uuid.UUID, error) {
	ms.ui.mx.Lock()
	ms.up.mx.Lock()
	ms.f.mx.Lock()
	defer ms.ui.mx.Unlock()
	defer ms.up.mx.Unlock()
	defer ms.f.mx.Unlock()

	_, ok := ms.ui.index[request.Login]
	if ok {
		return uuid.UUID{}, errors.New("user with login already exist")
	}

	currentUUID := uuid.New()

	ms.ui.index[request.Login] = currentUUID
	ms.ui.data[currentUUID] = userInfoData{
		login:   request.Login,
		pass:    request.Password,
		root:    request.Root.Or(false),
		regDate: oas.DateString(time.Now()),
	}
	ms.up.data[currentUUID] = userProfileData{
		email:      request.Email,
		lastModify: oas.DateString(time.Now()),
	}
	ms.f.data[currentUUID] = make(map[uuid.UUID]friendsData)

	return currentUUID, nil
}

func (ms *mockStorage) GetUserId(login oas.LoginString) (uuid.UUID, error) {
	ms.ui.mx.Lock()
	defer ms.ui.mx.Unlock()

	userUuid, ok := ms.ui.index[login]
	if !ok {
		return uuid.UUID{}, errors.New("user with login not exists")
	}
	return userUuid, nil
}

func (ms *mockStorage) GetProfile(user oas.UserId) (oas.ProfileInfo, error) {
	ms.ui.mx.Lock()
	ms.up.mx.Lock()
	defer ms.ui.mx.Unlock()
	defer ms.up.mx.Unlock()

	info, ok := ms.ui.data[uuid.UUID(user)]
	if !ok {
		return oas.ProfileInfo{}, errors.New("no user with uuid")
	}

	profile, ok := ms.up.data[uuid.UUID(user)]
	if !ok {
		log.Fatal("no user profile with uuid (INVARIANT BREAK)")
		return oas.ProfileInfo{}, errors.New("no user profile with uuid (INVARIANT BREAK)")
	}

	return oas.ProfileInfo{
		Login:      info.login,
		Email:      profile.email,
		Root:       oas.NewOptRootFlag(info.root),
		FirstName:  profile.firstName,
		LastName:   profile.lastName,
		ImageLink:  profile.imageLink,
		BirthDate:  profile.birthDate,
		Telephone:  profile.telephone,
		RegDate:    oas.NewOptDateString(info.regDate),
		LastModify: oas.NewOptDateString(profile.lastModify),
	}, nil
}

func (ms *mockStorage) GetPassword(user oas.UserId) (oas.PasswordString, error) {
	ms.ui.mx.Lock()
	defer ms.ui.mx.Unlock()

	info, ok := ms.ui.data[uuid.UUID(user)]
	if !ok {
		return oas.PasswordString(""), errors.New("no user with uuid")
	}

	return info.pass, nil
}

func (ms *mockStorage) GetFriends(user oas.UserId) (oas.FriendList, error) {
	ms.f.mx.Lock()
	defer ms.f.mx.Unlock()

	friends, ok := ms.f.data[uuid.UUID(user)]
	if !ok {
		return oas.FriendList{}, errors.New("no user with uuid")
	}

	ans := make([]oas.FriendObject, 0, len(friends))
	for friendId, data := range friends {
		fobject := oas.FriendObject{
			FriendID:   oas.UserId(friendId),
			Alias:      data.friendAlias,
			Subscibed:  data.subscribed,
			Hidden:     data.hidden,
			Paired:     data.paired,
			FriendedAt: oas.NewOptDateString(data.friendedAt),
			LastModify: oas.NewOptDateString(data.lastModify),
		}
		ans = append(ans, fobject)
	}
	return ans, nil
}

func (ms *mockStorage) UpdateProfile(user oas.UserId, update *oas.ProfileUpdate) error {
	ms.up.mx.Lock()
	defer ms.up.mx.Unlock()

	newValue, ok := ms.up.data[uuid.UUID(user)]
	if !ok {
		return errors.New("no user profile with uuid")
	}

	if update.Email.IsSet() {
		newValue.email = update.Email.Value
	}
	if update.FirstName.IsSet() {
		newValue.firstName.SetTo(update.FirstName.Value)
	}
	if update.LastName.IsSet() {
		newValue.lastName.SetTo(update.LastName.Value)
	}
	if update.ImageLink.IsSet() {
		newValue.imageLink.SetTo(update.ImageLink.Value)
	}
	if update.BirthDate.IsSet() {
		newValue.birthDate.SetTo(update.BirthDate.Value)
	}
	if update.Telephone.IsSet() {
		newValue.telephone.SetTo(update.Telephone.Value)
	}
	newValue.lastModify = oas.DateString(time.Now())

	ms.up.data[uuid.UUID(user)] = newValue
	return nil
}

func (ms *mockStorage) UpdateFriend(user oas.UserId, update *oas.FriendModify) error {
	ms.f.mx.Lock()
	defer ms.f.mx.Unlock()

	friends, ok := ms.f.data[uuid.UUID(user)]

	if !ok {
		return errors.New("no user with uuid")
	}

	friend, ok := friends[uuid.UUID(update.FriendID)]
	if !ok {
		if update.Delete.Or(false) {
			return errors.New("try to delete deleted user")
		}
		_, paired := ms.f.data[uuid.UUID(update.FriendID)][uuid.UUID(user)]

		ms.f.data[uuid.UUID(user)][uuid.UUID(update.FriendID)] = friendsData{
			friendAlias: update.Alias,
			subscribed:  update.Subscibed,
			hidden:      update.Hidden,
			paired:      oas.NewOptFriendPairedFlag(oas.FriendPairedFlag(paired)),
			friendedAt:  oas.DateString(time.Now()),
			lastModify:  oas.DateString(time.Now()),
		}
		return nil
	}
	if update.Delete.Or(false) {
		newValue, paired := ms.f.data[uuid.UUID(update.FriendID)][uuid.UUID(user)]
		if paired {
			newValue.paired.SetTo(false)
			ms.f.data[uuid.UUID(update.FriendID)][uuid.UUID(user)] = newValue
		}
		delete(ms.f.data[uuid.UUID(user)], uuid.UUID(update.FriendID))
		return nil
	}

	if update.Alias.IsSet() {
		friend.friendAlias.SetTo(update.Alias.Value)
	}
	if update.Subscibed.IsSet() {
		friend.subscribed.SetTo(update.Subscibed.Value)
	}
	if update.Hidden.IsSet() {
		friend.hidden.SetTo(update.Hidden.Value)
	}
	friend.lastModify = oas.DateString(time.Now())
	ms.f.data[uuid.UUID(user)][uuid.UUID(update.FriendID)] = friend
	return nil
}
