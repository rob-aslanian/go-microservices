package profile

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
)

// Profile represents user profile
type Profile struct {
	account.Account `bson:",inline"`

	Avatar   string `bson:"avatar"`
	Headline string `bson:"headline"`
	Story    string `bson:"story"`

	// Experiences     []*Experience     `bson:"experiences"`
	// Educations      []*Education      `bson:"educations"`
	// Skills          []*Skill          `bson:"skills"`
	// KnownLanguages  []*KnownLanguage  `bson:"known_languages"`
	// Interests       []*Interest       `bson:"interests"`
	// Accomplishments []*Accomplishment `bson:"accomplishments"`

	IsMe                    bool   `bson:"-"`
	IsMember                bool   `bson:"-"`
	IsFriend                bool   `bson:"-"`
	IsFollow                bool   `bson:"-"`
	IsFavorite              bool   `bson:"-"`
	IsFriendRequestSend     bool   `bson:"-"`
	IsFriendRequestRecieved bool   `bson:"-"`
	FriendshipID            string `bson:"-"`
	IsBlocked               bool   `bson:"-"`
	IsOnline                bool   `bson:"-"`
	CompletePercent         int8   `bson:"-"`
	MutualConectionsAmount  int32  `bson:"-"`

	Translation           map[string]*Translation `bson:"translation"`
	CurrentTranslation    string                  `bson:"-"`
	AvailableTranslations []string                `bson:"-"`
}


// Users ... 
type Users struct{
   UsersAmount 	int32   			`bson:"users_amount"`
   Users 		[]*account.User 	`bson:"users"`	
}


// Translation ...
type Translation struct {
	FirstName string `bson:"first_name"`
	Lastname  string `bson:"last_name"`
	Headline  string `bson:"headline"`
	Story     string `bson:"story"`
	Nickname  string `bson:"nickname"`
}

// ApplyPrivacies removes some info depending of privacy settings
func (p *Profile) ApplyPrivacies(ctx context.Context) {
	if p.Privacy == nil {
		return
	}

	if !isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.Privacy.ProfilePicture) {
		p.Avatar = ""
	}

	if !isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.Privacy.ActiveStatus) {
		p.IsOnline = false
	}

	// TODO:
	// p.Privacy.FindByEmail
	// p.Privacy.FindByPhone
	// p.Privacy.MyConnections
	// p.Privacy.ShareEdits

}

// ApplyPermissions removes some information depending on permissions
func (p *Profile) ApplyPermissions(ctx context.Context) {
	if p.Privacy == nil {
		return
	}

	if p.Birthday != nil {
		if !isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.Birthday.Permission.Type) {
			p.Birthday = nil
		}
	}

	if p.Location != nil {
		if !isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.Location.Permission.Type) {
			p.Location = nil
		}
	}

	if p.MiddleName != nil && p.MiddleName.Permission != nil {
		if !isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.MiddleName.Permission.Type) {
			p.MiddleName = nil
		}
	}

	if p.NativeName != nil && p.NativeName.Permission != nil {
		if !isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.NativeName.Permission.Type) {
			p.NativeName = nil
		}
	}

	if p.Nickname != nil && p.Nickname.Permission != nil {
		if !isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.Nickname.Permission.Type) {
			p.Nickname = nil
		}
	}

	if p.Patronymic != nil && p.Patronymic.Permission != nil {
		if !isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.Patronymic.Permission.Type) {
			p.Patronymic = nil
		}
	}

	if len(p.Emails) > 0 {
		emails := make([]*account.Email, 0, len(p.Emails))
		for i := range p.Emails {
			if isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.Emails[i].Permission.Type) {
				emails = append(emails, p.Emails[i])
			}
		}
		p.Emails = emails
	}

	if len(p.Phones) > 0 {
		phones := make([]*account.Phone, 0, len(p.Phones))
		for i := range p.Phones {
			if isAllowed(p.IsMe, p.IsFriend, p.IsMember, p.Phones[i].Permission.Type) {
				phones = append(phones, p.Phones[i])
			}
		}
		p.Phones = phones
	}
}

func isAllowed(isMe bool, isFriend bool, isMember bool, permission account.PermissionType) bool {
	if isMe || permission == account.PermissionTypeEveryone {
		return true
	} else if isFriend {
		if permission != account.PermissionTypeMe {
			return true
		}
	} else if isMember {
		if permission == account.PermissionTypeMembers {
			return true
		}
	}

	return false
}

// Translate ...
func (p *Profile) Translate(ctx context.Context, lang string) string {
	if p == nil || lang == "" {
		return "en"
	}

	if tr, isExists := p.Translation[lang]; isExists {
		p.FirstName = tr.FirstName
		p.Lastname = tr.Lastname
		p.Headline = tr.Headline
		p.Story = tr.Story

		if p.Nickname != nil {
			p.Nickname.Nickname = tr.Nickname
		}

	} else {
		return "en"
	}

	return lang
}
