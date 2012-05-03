package user

import (
	"github.com/ungerik/go-mail"
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/view"
	"launchpad.net/mgo/bson"
	"reflect"
)

var (
	UserPtrType = reflect.TypeOf((*User)(nil))
	UserType    = UserPtrType.Elem()
)

func From(document interface{}) (user *User) {
	v := reflect.ValueOf(document)

	// Dereference pointers until we have a struct value
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		if v.Type() == UserType {
			return v.Addr().Interface().(*User)
		}

		count := v.NumField()
		for i := 0; i < count; i++ {
			field := v.Field(i)
			if field.Type() == UserType {
				return field.Addr().Interface().(*User)
			}
		}
	}

	panic(errs.Format("user.From(): invalid document type %T", document))
}

func New(email, password string) (user *User, doc interface{}, err error) {
	doc = Config.Collection.NewDocument()
	user = From(doc)
	err = user.AddEmail(email, "via signup")
	if err != nil {
		return nil, nil, err
	}
	user.Username.Set(user.Email[0].Address.Get())
	user.Password.SetHashed(password)
	return user, doc, nil
}

// Will use user.User.Username for search
// Sets the username also as first name
func EnsureExists(username, email, password string, admin bool) (doc interface{}, exists bool, err error) {
	doc, exists, err = Config.Collection.Filter("Username", username).GetOrCreateOne()
	if err != nil {
		return nil, false, err
	}
	user := From(doc)

	user.Username.Set(username)
	user.Name.First.Set(username)
	user.Name.Last.Set(username)

	if user.PrimaryEmail() == "" {
		err = user.AddEmail(email, "via gostart/user.EnsureExists()")
	} else {
		err = user.Email[0].Address.Set(email)
		user.Email[0].Description.Set("via gostart/user.EnsureExists()")
	}
	if err != nil {
		return nil, false, err
	}
	if user.Email[0].Confirmed.IsEmpty() {
		user.Email[0].Confirmed.SetNowUTC()
	}

	user.Password.SetHashed(password)

	user.Blocked.Set(false)
	user.Admin.Set(admin)

	err = user.Save()
	if err != nil {
		return nil, false, err
	}
	return doc, exists, nil
}

func FindByID(id string) (userDoc interface{}, found bool, err error) {
	return Config.Collection.TryDocumentWithID(bson.ObjectIdHex(id))
}

func IsConfirmedUserID(id string) (confirmed bool, err error) {
	userDoc, found, err := FindByID(id)
	if !found {
		return false, err
	}
	return From(userDoc).IdentityConfirmed(), nil
}

///////////////////////////////////////////////////////////////////////////////
// Email functions

func FindByEmail(addr string) (doc interface{}, found bool, err error) {
	addr = email.NormalizeAddress(addr)
	query := Config.Collection.Filter("Email.Address", addr)
	return query.TryOne()
}

func ConfirmEmail(confirmationCode string) (userDoc interface{}, email string, confirmed bool, err error) {
	query := Config.Collection.Filter("Email.ConfirmationCode", confirmationCode)
	userDoc, found, err := query.TryOne()
	if !found {
		return nil, "", false, err
	}
	user := From(userDoc)

	for i := range user.Email {
		if user.Email[i].ConfirmationCode.Get() == confirmationCode {
			user.Email[i].Confirmed.SetNowUTC()
			email = user.Email[i].Address.Get()
		}
	}
	errs.Assert(email != "", "ConfirmationCode has to be found")

	err = user.Save()
	if err != nil {
		return nil, "", false, err
	}

	return userDoc, email, true, nil
}

///////////////////////////////////////////////////////////////////////////////
// Session functions

const ContextCacheKey = "github.com/ungerik/go-start/user.ContextCacheKey"

func Login(session *view.Session, userDoc interface{}) {
	session.SetID(userDoc.(mongo.Document).ObjectId().Hex())
	session.User = userDoc
}

func Logout(session *view.Session) {
	session.DeleteID()
	session.User = nil
}

func LoginEmailPassword(session *view.Session, email, password string) (emailPasswdMatch bool, err error) {
	userDoc, found, err := FindByEmail(email)
	if !found {
		return false, err
	}
	if !From(userDoc).EmailPasswordMatch(email, password) {
		return false, nil
	}
	Login(session, userDoc)
	return true, nil
}

// Returns nil if there is no session user
func OfSession(session *view.Session) (userDoc interface{}) {
	if session.User != nil {
		return session.User
	}
	id, ok := session.ID()
	if !ok {
		return nil
	}
	userDoc, _, _ = FindByID(id)
	session.User = userDoc
	return userDoc
}
