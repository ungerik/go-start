package user

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/modelext"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/view"
)

func NewCollection(name string) *mongo.Collection {
	return mongo.NewCollection(name, modelext.NameDocLabelSelectors...)
}

type User struct {
	mongo.DocumentBase `bson:",inline"`
	Name               modelext.Name
	Username           model.String `view:"size=20"`
	Password           model.Password
	Blocked            model.Bool
	Admin              model.Bool
	PostalAddress      modelext.PostalAddress `view:"label=Postal Address"`
	Phone              []PhoneNumber
	Web                []Website
	Email              []EmailIdentity
	Facebook           []FacebookIdentity
	Twitter            []TwitterIdentity
	LinkedIn           []LinkedInIdentity
	Xing               []XingIdentity
	GitHub             []GitHubIdentity
	Skype              []SkypeIdentity
}

func (self *User) SetEmailPassword(email, password string) error {
	err := self.AddEmail(email, "via signup")
	if err != nil {
		return err
	}
	self.Username.Set(self.Email[0].Address.Get())
	self.Password.SetHashed(password)
	return nil
}

func (self *User) IdentityConfirmed() bool {
	return self.Blocked == false && (self.EmailPasswordConfirmed() ||
		self.FacebookIdentityConfirmed() ||
		self.TwitterIdentityConfirmed() ||
		self.LinkedInIdentityConfirmed())
}

func (self *User) EmailPasswordMatch(email, password string) bool {
	if !self.EmailPasswordConfirmed() {
		return false
	}
	return self.Password.EqualsHashed(password)
}

// One email address has to be confirmed and the password must not be empty
func (self *User) EmailPasswordConfirmed() bool {
	if self.Password != "" {
		for i := range self.Email {
			if self.Email[i].Confirmed != "" {
				return true
			}
		}
	}
	return false
}

func (self *User) EmailConfirmed(email string) bool {
	for i := range self.Email {
		if self.Email[i].Address.Get() == email {
			return !self.Email[i].Confirmed.IsEmpty()
		}
	}
	return false
}

func (self *User) ConfirmEmailPassword() {
	if len(self.Email) > 0 {
		self.Email[0].Confirmed.SetNowUTC()
	}
}

func (self *User) FacebookIdentityConfirmed() bool {
	for i := range self.Facebook {
		if self.Facebook[i].Confirmed != "" {
			return true
		}
	}
	return false
}

func (self *User) TwitterIdentityConfirmed() bool {
	for i := range self.Twitter {
		if self.Twitter[i].Confirmed != "" {
			return true
		}
	}
	return false
}

func (self *User) LinkedInIdentityConfirmed() bool {
	for i := range self.Twitter {
		if self.Twitter[i].Confirmed != "" {
			return true
		}
	}
	return false
}

func (self *User) PrimaryEmail() string {
	if len(self.Email) == 0 {
		return ""
	}
	return self.Email[0].Address.Get()
}

func (self *User) HasEmail(address string) bool {
	for _, email := range self.Email {
		if email.Address.EqualsCaseinsensitive(address) {
			return true
		}
	}
	return false
}

func (self *User) AddEmail(address, description string) (err error) {
	var email EmailIdentity
	err = email.Address.Set(address)
	if err != nil {
		return err
	}
	email.Description.Set(description)
	self.Email = append(self.Email, email)

	/// self.Email[0] = email

	return nil
}

func (self *User) PrimaryPhone() string {
	if len(self.Phone) == 0 {
		return ""
	}
	return self.Phone[0].Number.Get()
}

func (self *User) HasPhone(number string) bool {
	if number == "" {
		return false
	}
	number = model.NormalizePhoneNumber(number)
	for i := range self.Phone {
		if self.Phone[i].Number.Get() == number {
			return true
		}
	}
	return false
}

func (self *User) AddPhone(number, description string) {
	var phone PhoneNumber
	phone.Number.Set(number)
	phone.Description.Set(description)
	self.Phone = append(self.Phone, phone)
}

func (self *User) PrimaryWeb() string {
	if len(self.Web) == 0 {
		return ""
	}
	return self.Web[0].Url.Get()
}

func (self *User) PrimaryEmailIdentity() *EmailIdentity {
	if len(self.Email) == 0 {
		return nil
	}
	return &self.Email[0]
}

func (self *User) PrimaryFacebookIdentity() *FacebookIdentity {
	if len(self.Facebook) == 0 {
		return nil
	}
	return &self.Facebook[0]
}

func (self *User) PrimaryTwitterIdentity() *TwitterIdentity {
	if len(self.Twitter) == 0 {
		return nil
	}
	return &self.Twitter[0]
}

func (self *User) PrimaryLinkedInIdentity() *LinkedInIdentity {
	if len(self.LinkedIn) == 0 {
		return nil
	}
	return &self.LinkedIn[0]
}

func (self *User) PrimaryXingIdentity() *XingIdentity {
	if len(self.Xing) == 0 {
		return nil
	}
	return &self.Xing[0]
}

func (self *User) PrimaryGitHubIdentity() *GitHubIdentity {
	if len(self.GitHub) == 0 {
		return nil
	}
	return &self.GitHub[0]
}

func (self *User) PrimarySkypeIdentity() *SkypeIdentity {
	if len(self.Skype) == 0 {
		return nil
	}
	return &self.Skype[0]
}

///////////////////////////////////////////////////////////////////////////////
// PhoneNumber

type PhoneNumber struct {
	Number      model.Phone
	Description model.String
}

///////////////////////////////////////////////////////////////////////////////
// Website

// Implements view.LinkModel
type Website struct {
	Url         model.Url `view:"label=URL"`
	Title       model.String
	Description model.String
}

func (self *Website) URL(ctx *view.Context) string {
	return self.Url.Get()
}

func (self *Website) LinkContent(ctx *view.Context) view.View {
	return view.Escape(self.LinkTitle(ctx))
}

func (self *Website) LinkTitle(ctx *view.Context) string {
	if self.Title.IsEmpty() {
		return self.Description.Get()
	}
	return self.Title.Get()
}

func (self *Website) LinkRel(ctx *view.Context) string {
	return ""
}
