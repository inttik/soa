package postgresstorage

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"users/oas"
)

var (
	ErrUserExist = errors.New("user already exist")
	ErrNoUser    = errors.New("user not exist")
)

func (ps *postgresStorage) CreateUser(request *oas.CreateUserRequest) (uuid.UUID, error) {
	newUser := userInfo{
		Login:   string(request.Login),
		Pass:    string(request.Password),
		Root:    bool(request.Root.Or(false)),
		RegDate: time.Now(),
	}
	newUserProfile := userProfile{
		Email:      string(request.Email),
		LastModify: time.Now(),
	}

	err := ps.db.Transaction(func(tx *gorm.DB) error {
		user := userInfo{}
		result := tx.Where("login = ?", request.Login).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			result := tx.Create(&newUser)
			if result.Error != nil {
				return result.Error
			}
			newUserProfile.Id = newUser.Id
			result = tx.Create(&newUserProfile)
			if result.Error != nil {
				return result.Error
			}
			return nil
		} else {
			return ErrUserExist
		}
	})

	if err != nil {
		return uuid.UUID{}, err
	}
	return newUser.Id, nil
}

func (ps *postgresStorage) GetUserId(login oas.LoginString) (uuid.UUID, error) {
	id := uuid.UUID{}
	err := ps.db.Transaction(func(tx *gorm.DB) error {
		user := userInfo{}
		result := tx.Where("login = ?", login).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoUser
		} else {
			id = user.Id
			return nil
		}
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (ps *postgresStorage) GetProfile(user oas.UserId) (oas.ProfileInfo, error) {
	resp := oas.ProfileInfo{}
	err := ps.db.Transaction(func(tx *gorm.DB) error {
		info := userInfo{}
		profile := userProfile{}

		result := tx.Where("id = ?", uuid.UUID(user).String()).Limit(1).Find(&info)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoUser
		}

		result = tx.Where("id = ?", uuid.UUID(user).String()).Limit(1).Find(&profile)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoUser
		}

		resp.Login = oas.LoginString(info.Login)
		resp.Email = oas.EmailString(profile.Email)
		resp.Root = oas.NewOptRootFlag(oas.RootFlag(info.Root))
		resp.RegDate = oas.NewOptDateString(oas.DateString(info.RegDate))
		resp.LastModify = oas.NewOptDateString(oas.DateString(profile.LastModify))

		if profile.FirstName != "" {
			resp.FirstName = oas.NewOptNameString(oas.NameString(profile.FirstName))
		}
		if profile.LastName != "" {
			resp.LastName = oas.NewOptNameString(oas.NameString(profile.LastName))
		}
		if profile.ImageLink != "" {
			link, err := url.Parse(profile.ImageLink)
			if err != nil {
				return err
			}
			resp.ImageLink = oas.NewOptLinkString(oas.LinkString(*link))
		}
		if profile.BirthDate != (time.Time{}) {
			resp.BirthDate = oas.NewOptBirthString(oas.BirthString(profile.BirthDate))
		}
		if profile.Telephone != "" {
			resp.Telephone = oas.NewOptTelephoneString(oas.TelephoneString(profile.Telephone))
		}
		return nil
	})
	if err != nil {
		return oas.ProfileInfo{}, err
	}
	return resp, nil
}

func (ps *postgresStorage) GetPassword(user oas.UserId) (oas.PasswordString, error) {
	pass := oas.PasswordString("")
	err := ps.db.Transaction(func(tx *gorm.DB) error {
		info := userInfo{}
		result := tx.Where("id = ?", uuid.UUID(user).String()).Limit(1).Find(&info)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoUser
		} else {
			pass = oas.PasswordString(info.Pass)
			return nil
		}
	})
	if err != nil {
		return oas.PasswordString(""), err
	}
	return pass, nil
}

func (ps *postgresStorage) UpdateProfile(user oas.UserId, update *oas.ProfileUpdate) error {
	modify := userProfile{}

	if update.Email.IsSet() {
		modify.Email = string(update.Email.Value)
	}
	if update.FirstName.IsSet() {
		modify.FirstName = string(update.FirstName.Value)
	}
	if update.LastName.IsSet() {
		modify.LastName = string(update.LastName.Value)
	}
	if update.ImageLink.IsSet() {
		link := url.URL(update.ImageLink.Value)
		modify.ImageLink = (&link).String()
	}
	if update.BirthDate.IsSet() {
		modify.BirthDate = time.Time(update.BirthDate.Value)
	}
	if update.Telephone.IsSet() {
		modify.Telephone = string(update.Telephone.Value)
	}
	modify.LastModify = time.Now()

	err := ps.db.Transaction(func(tx *gorm.DB) error {
		profile := userProfile{}
		result := tx.Where("id = ?", uuid.UUID(user).String()).Limit(1).Find(&profile)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoUser
		}
		result = tx.Save(&profile)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoUser
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (ps *postgresStorage) GetFriends(user oas.UserId) (oas.FriendList, error) {
	return oas.FriendList{}, errors.New("not implemented")
}

func (ps *postgresStorage) UpdateFriend(user oas.UserId, update *oas.FriendModify) error {
	return errors.New("not implemented")
}
