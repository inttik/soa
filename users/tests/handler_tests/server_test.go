package handlers_test

import (
	"net/url"
	"testing"
	"time"
	"users/oas"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name string
		rh   registerHelper
	}{
		{
			name: "simple",
			rh: registerHelper{
				req: oas.CreateUserRequest{
					Login:    "login",
					Password: "pass",
					Email:    "email@example",
				},
				expect201: true,
			},
		},
		{
			name: "verbose simple",
			rh: registerHelper{
				req: oas.CreateUserRequest{
					Login:    "login",
					Password: "pass",
					Email:    "email@example",
					Root:     oas.NewOptRootFlag(false),
				},
				expect201: true,
			},
		},
		{
			name: "bad format",
			rh: registerHelper{
				req: oas.CreateUserRequest{
					Login:    "login",
					Password: "pass",
					Email:    "email_at_example",
				},
				expect400: true,
			},
		},
		{
			name: "already created",
			rh: registerHelper{
				req: oas.CreateUserRequest{
					Login:    "admin",
					Password: "admin",
					Email:    "email@example",
				},
				expect400: true,
			},
		},
		{
			name: "bad access",
			rh: registerHelper{
				req: oas.CreateUserRequest{
					Login:    "login",
					Password: "pass",
					Email:    "email@example",
					Root:     oas.NewOptRootFlag(true),
				},
				expect403: true,
			},
		},
		{
			name: "valid access",
			rh: registerHelper{
				req: oas.CreateUserRequest{
					Login:    "login",
					Password: "pass",
					Email:    "email@example",
					Root:     oas.NewOptRootFlag(true),
				},
				admin:     true,
				expect201: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := state{}
			s.setup(t)

			s.applyRegister(test.rh, t)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name          string
		lh            loginHelper
		registerFirst bool
		admin         bool
	}{
		{
			name: "admin login",
			lh: loginHelper{
				req: oas.LoginUserRequest{
					Login:    "admin",
					Password: "admin",
				},
				expect200: true,
			},
		},
		{
			name: "user login",
			lh: loginHelper{
				req: oas.LoginUserRequest{
					Login:    "user",
					Password: "pass",
				},
				expect200: true,
			},
			registerFirst: true,
		},
		{
			name: "non default admin login",
			lh: loginHelper{
				req: oas.LoginUserRequest{
					Login:    "admin2",
					Password: "pass2",
				},
				expect200: true,
			},
			registerFirst: true,
			admin:         true,
		},
		{
			name: "admin bad pass",
			lh: loginHelper{
				req: oas.LoginUserRequest{
					Login:    "admin",
					Password: "incorrect",
				},
				expect400: true,
			},
		},
		{
			name: "not found",
			lh: loginHelper{
				req: oas.LoginUserRequest{
					Login:    "user",
					Password: "pass",
				},
				expect404: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := state{}
			s.setup(t)
			if test.registerFirst {
				s.applyRegister(registerHelper{
					req: oas.CreateUserRequest{
						Login:    test.lh.req.Login,
						Password: test.lh.req.Password,
						Email:    "default@email.com",
					},
					admin:     test.admin,
					expect201: true,
				}, t)
			}
			s.applyLogin(test.lh, t)
		})
	}
}

func TestNewAdminIsAdmin(t *testing.T) {
	s := state{}
	s.setup(t)

	UUID, err := s.applyRegister(registerHelper{
		req: oas.CreateUserRequest{
			Login:    "admin2",
			Password: "pass2",
			Email:    "default@email.com",
		},
		admin:     true,
		expect201: true,
	}, t)
	assert.NoError(t, err, "201 is expected")

	resp, err := s.applyLogin(loginHelper{
		req: oas.LoginUserRequest{
			Login:    "admin2",
			Password: "pass2",
		},
		expect200: true,
	}, t)
	assert.NoError(t, err, "200 is expected")
	assert.Equalf(t, *UUID, resp.UserID, "UUID not change")

	s.adminToken.Token = string(resp.Token)
	_, err = s.applyRegister(registerHelper{
		req: oas.CreateUserRequest{
			Login:    "admin3",
			Password: "pass3",
			Email:    "default@email.com",
		},
		admin:     true,
		expect201: true,
	}, t)
	assert.NoError(t, err, "New admin is admin")
}

func TestProfileGet(t *testing.T) {
	user1 := registerHelper{
		req: oas.CreateUserRequest{
			Login:    "user1",
			Password: "pass1",
			Email:    "email1@gm.co",
		},
		expect201: true,
	}
	user2 := registerHelper{
		req: oas.CreateUserRequest{
			Login:    "user2",
			Password: "pass2",
			Email:    "email2@gm.co",
		},
		expect201: true,
	}
	s := state{}
	s.setup(t)

	id1, err := s.applyRegister(user1, t)
	assert.NoErrorf(t, err, "correct registration")
	id2, err := s.applyRegister(user2, t)
	assert.NoErrorf(t, err, "correct registation")

	invalid_id := oas.UserId(uuid.New())
	id3 := &invalid_id

	t.Run("read self", func(t *testing.T) {
		profile, err := s.applyProfileGet(profileGetHelper{
			req: *id1,
			logined: &oas.LoginUserRequest{
				Login:    "user1",
				Password: "pass1",
			},
			expect200: true,
		}, t)
		assert.NoErrorf(t, err, "expect 200")

		assert.Equal(t, "user1", string(profile.Login))
		assert.Equal(t, "email1@gm.co", string(profile.Email))
		assert.Equalf(t, false, bool(profile.Root.Or(true)), "extended info")

		profile, err = s.applyProfileGet(profileGetHelper{
			req: *id2,
			logined: &oas.LoginUserRequest{
				Login:    "user2",
				Password: "pass2",
			},
			expect200: true,
		}, t)
		assert.NoErrorf(t, err, "expect 200")

		assert.Equal(t, "user2", string(profile.Login))
		assert.Equal(t, "email2@gm.co", string(profile.Email))
		assert.Equalf(t, false, bool(profile.Root.Or(true)), "extended info")
	})

	t.Run("read other", func(t *testing.T) {
		profile, err := s.applyProfileGet(profileGetHelper{
			req: *id2,
			logined: &oas.LoginUserRequest{
				Login:    "user1",
				Password: "pass1",
			},
			expect200: true,
		}, t)
		assert.NoErrorf(t, err, "expect 200")

		assert.Equal(t, "user2", string(profile.Login))
		assert.Equal(t, "email2@gm.co", string(profile.Email))
		assert.Equalf(t, true, bool(profile.Root.Or(true)), "no extended info")

		profile, err = s.applyProfileGet(profileGetHelper{
			req:       *id2,
			expect200: true,
		}, t)
		assert.NoErrorf(t, err, "expect 200")

		assert.Equal(t, "user2", string(profile.Login))
		assert.Equal(t, "email2@gm.co", string(profile.Email))
		assert.Equalf(t, true, bool(profile.Root.Or(true)), "no extended info")

		profile, err = s.applyProfileGet(profileGetHelper{
			req:       *id2,
			admin:     true,
			expect200: true,
		}, t)
		assert.NoErrorf(t, err, "expect 200")

		assert.Equal(t, "user2", string(profile.Login))
		assert.Equal(t, "email2@gm.co", string(profile.Email))
		assert.Equalf(t, false, bool(profile.Root.Or(true)), "extended info")
	})

	t.Run("read no exist", func(t *testing.T) {
		s.applyProfileGet(profileGetHelper{
			req: *id3,
			logined: &oas.LoginUserRequest{
				Login:    "user1",
				Password: "pass1",
			},
			expect404: true,
		}, t)

		s.applyProfileGet(profileGetHelper{
			req:       *id3,
			expect404: true,
		}, t)

		s.applyProfileGet(profileGetHelper{
			req:       *id3,
			admin:     true,
			expect404: true,
		}, t)
	})
}

func TestProfileUpdate(t *testing.T) {
	imageLink, err := url.Parse("http://link.com")
	assert.NoErrorf(t, err, "regular link")
	link := oas.LinkString(*imageLink)

	birthTime, err := time.Parse("2006-01-02", "2020-12-31")
	assert.NoErrorf(t, err, "regular date")
	birth := oas.BirthString(birthTime)

	badUUID := uuid.New()
	badUserId := oas.UserId(badUUID)

	user := registerHelper{
		req: oas.CreateUserRequest{
			Login:    "user",
			Password: "pass",
			Email:    "email@gm.co",
		},
		expect201: true,
	}
	user2 := registerHelper{
		req: oas.CreateUserRequest{
			Login:    "user2",
			Password: "pass",
			Email:    "email@gm.co",
		},
		expect201: true,
	}
	tests := []struct {
		name      string
		pph       profilePostHelper
		check_set bool
	}{
		{
			name: "fully updated",
			pph: profilePostHelper{
				register: &user,
				req: oas.ProfileUpdate{
					Email:     oas.NewOptEmailString("new@gmail.com"),
					FirstName: oas.NewOptNameString("Ivan"),
					LastName:  oas.NewOptNameString("Ivanov"),
					ImageLink: oas.NewOptLinkString(link),
					BirthDate: oas.NewOptBirthString(birth),
					Telephone: oas.NewOptTelephoneString("+32 (325) 123-23-32"),
				},
				admin:     true,
				expect200: true,
			},
			check_set: true,
		},
		{
			name: "partial updated 1",
			pph: profilePostHelper{
				register: &user,
				req: oas.ProfileUpdate{
					Email:     oas.NewOptEmailString("new@gmail.com"),
					FirstName: oas.NewOptNameString("Ivan"),
				},
				admin:     true,
				expect200: true,
			},
			check_set: true,
		},
		{
			name: "partial updated 2",
			pph: profilePostHelper{
				register: &user,
				req: oas.ProfileUpdate{
					LastName: oas.NewOptNameString("Ivanov"),
				},
				admin:     true,
				expect200: true,
			},
			check_set: true,
		},
		{
			name: "self update",
			pph: profilePostHelper{
				register: &user,
				req: oas.ProfileUpdate{
					LastName: oas.NewOptNameString("Ivanov"),
				},
				logined: &oas.LoginUserRequest{
					Login:    "user",
					Password: "pass",
				},
				expect200: true,
			},
			check_set: true,
		},
		{
			name: "bad format 1",
			pph: profilePostHelper{
				register: &user,
				req: oas.ProfileUpdate{
					Email: oas.NewOptEmailString("new_at_gmail.com"),
				},
				admin:     true,
				expect400: true,
			},
		},
		{
			name: "bad format 2",
			pph: profilePostHelper{
				register: &user,
				req: oas.ProfileUpdate{
					Telephone: oas.NewOptTelephoneString("88005553535"),
				},
				admin:     true,
				expect400: true,
			},
		},
		{
			name: "not auth",
			pph: profilePostHelper{
				register: &user,
				req: oas.ProfileUpdate{
					FirstName: oas.NewOptNameString("Ivan"),
				},
				expect401: true,
			},
		},
		{
			name: "bad access",
			pph: profilePostHelper{
				register: &user,
				logined: &oas.LoginUserRequest{
					Login:    "user2",
					Password: "pass",
				},
				req: oas.ProfileUpdate{
					FirstName: oas.NewOptNameString("Ivan"),
				},
				expect403: true,
			},
		},
		{
			name: "not found",
			pph: profilePostHelper{
				target: &badUserId,
				admin:  true,
				req: oas.ProfileUpdate{
					FirstName: oas.NewOptNameString("Ivan"),
				},
				expect404: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := state{}
			s.setup(t)

			_, err := s.applyRegister(user2, t)
			assert.NoErrorf(t, err, "regular registration")

			profile, err := s.applyProfilePost(test.pph, t)
			if !test.check_set {
				return
			}
			assert.NoError(t, err, "expect 200")

			if test.pph.req.Email.IsSet() {
				assert.Equal(t, test.pph.req.Email.Value, profile.Email)
			}
			if test.pph.req.FirstName.IsSet() {
				assert.Equal(t, test.pph.req.FirstName.Value, profile.FirstName.Value)
			} else {
				assert.Equal(t, false, profile.FirstName.IsSet())
			}
			if test.pph.req.LastName.IsSet() {
				assert.Equal(t, test.pph.req.LastName.Value, profile.LastName.Value)
			} else {
				assert.Equal(t, false, profile.LastName.IsSet())
			}
			if test.pph.req.ImageLink.IsSet() {
				assert.Equal(t, test.pph.req.ImageLink.Value, profile.ImageLink.Value)
			} else {
				assert.Equal(t, false, profile.ImageLink.IsSet())
			}
			if test.pph.req.BirthDate.IsSet() {
				assert.Equal(t, test.pph.req.BirthDate.Value, profile.BirthDate.Value)
			} else {
				assert.Equal(t, false, profile.BirthDate.IsSet())
			}
			if test.pph.req.Telephone.IsSet() {
				assert.Equal(t, test.pph.req.Telephone.Value, profile.Telephone.Value)
			} else {
				assert.Equal(t, false, profile.Telephone.IsSet())
			}
		})
	}
}
