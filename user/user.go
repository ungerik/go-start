package user

import (
	"github.com/ungerik/go-mail"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/modelext"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// User

type User struct {
	mongo.DocumentBase `bson:",inline"`
	Name               modelext.Name
	Username           model.String
	Password           model.Password
	Blocked            model.Bool
	Admin              model.Bool
	PostalAddress      modelext.PostalAddress `gostart:"label=Postal Address"`
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

func (self *User) AddMissingIdentityFields() {
	if len(self.Phone) == 0 {
		self.Phone = make([]PhoneNumber, 1)
	}
	if len(self.Web) == 0 {
		self.Web = make([]Website, 1)
	}
	if len(self.Email) == 0 {
		self.Email = make([]EmailIdentity, 1)
	}
	if len(self.Facebook) == 0 {
		self.Facebook = make([]FacebookIdentity, 1)
	}
	if len(self.Twitter) == 0 {
		self.Twitter = make([]TwitterIdentity, 1)
	}
	if len(self.LinkedIn) == 0 {
		self.LinkedIn = make([]LinkedInIdentity, 1)
	}
	if len(self.Xing) == 0 {
		self.Xing = make([]XingIdentity, 1)
	}
	if len(self.GitHub) == 0 {
		self.GitHub = make([]GitHubIdentity, 1)
	}
	if len(self.Skype) == 0 {
		self.Skype = make([]SkypeIdentity, 1)
	}
}

func (self *User) RemoveEmptyIdentityFields() {
	self.Phone = utils.DeleteEmptySliceElements(self.Phone).([]PhoneNumber)
	self.Web = utils.DeleteEmptySliceElements(self.Web).([]Website)
	self.Email = utils.DeleteEmptySliceElements(self.Email).([]EmailIdentity)
	self.Facebook = utils.DeleteEmptySliceElements(self.Facebook).([]FacebookIdentity)
	self.Twitter = utils.DeleteEmptySliceElements(self.Twitter).([]TwitterIdentity)
	self.LinkedIn = utils.DeleteEmptySliceElements(self.LinkedIn).([]LinkedInIdentity)
	self.Xing = utils.DeleteEmptySliceElements(self.Xing).([]XingIdentity)
	self.GitHub = utils.DeleteEmptySliceElements(self.GitHub).([]GitHubIdentity)
	self.Skype = utils.DeleteEmptySliceElements(self.Skype).([]SkypeIdentity)
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
	address = email.NormalizeAddress(address)
	for _, email := range self.Email {
		if email.Address.Get() == address {
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
	Url         model.Url `gostart:"label=URL"`
	Title       model.String
	Description model.String
}

func (self *Website) URL(context *view.Context, args ...string) string {
	return self.Url.Get()
}

func (self *Website) LinkContent(context *view.Context) view.View {
	return view.Escape(self.LinkTitle(context))
}

func (self *Website) LinkTitle(context *view.Context) string {
	if self.Title.IsEmpty() {
		return self.Description.Get()
	}
	return self.Title.Get()
}

func (self *Website) LinkRel(context *view.Context) string {
	return ""
}
