package mockstorage

import (
	"sync"

	"github.com/google/uuid"

	"users/oas"
)

type userInfoData struct {
	login   oas.LoginString
	pass    oas.PasswordString
	root    oas.RootFlag
	regDate oas.DateString
}

type userInfo struct {
	data  map[uuid.UUID]userInfoData
	index map[oas.LoginString]uuid.UUID // index login -> uuid
	mx    sync.Mutex
}

type userProfileData struct {
	firstName  oas.OptNameString
	lastName   oas.OptNameString
	imageLink  oas.OptLinkString
	birthDate  oas.OptBirthString
	telephone  oas.OptTelephoneString
	email      oas.EmailString
	lastModify oas.DateString
}

type userProfile struct {
	data map[uuid.UUID]userProfileData
	mx   sync.Mutex
}

type friendsData struct {
	friendAlias oas.OptFriendAliasString
	subscribed  oas.OptFriendSubscribedFlag
	hidden      oas.OptFriendHiddenFlag
	paired      oas.OptFriendPairedFlag
	friendedAt  oas.DateString
	lastModify  oas.DateString
}

type friends struct {
	data map[uuid.UUID]map[uuid.UUID]friendsData
	mx   sync.Mutex
}
