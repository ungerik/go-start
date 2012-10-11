package user

import (
	"reflect"
	"strings"

	"labix.org/v2/mgo/bson"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/view"
)

var UserPtrType = reflect.TypeOf((*User)(nil))
var UserType = UserPtrType.Elem()

func FromDocument(document interface{}) (user *User) {
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

// Will use user.User.Username for search
// Sets the username also as first name
func EnsureExists(username, email, password string, admin bool, resultPtr interface{}) (exists bool, err error) {
	exists, err = Config.Collection.Filter("Username", username).TryOneDocument(resultPtr)
	if err != nil {
		return false, err
	}
	user := FromDocument(resultPtr)

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
		return false, err
	}
	if user.Email[0].Confirmed.IsEmpty() {
		user.Email[0].Confirmed.SetNowUTC()
	}

	user.Password.SetHashed(password)

	user.Blocked.Set(false)
	user.Admin.Set(admin)

	err = user.Save()
	if err != nil {
		return false, err
	}
	return exists, nil
}

func WithID(id string, resultPtr interface{}) (found bool, err error) {
	return Config.Collection.TryDocumentWithID(bson.ObjectIdHex(id), resultPtr)
}

func IsConfirmedUserID(id string) (confirmed bool, err error) {
	var user User
	found, err := WithID(id, &user)
	if !found {
		return false, err
	}
	return user.IdentityConfirmed(), nil
}

///////////////////////////////////////////////////////////////////////////////
// Email functions

func WithEmail(addr string, resultPtr interface{}) (found bool, err error) {
	query := Config.Collection.FilterEqualCaseInsensitive("Email.Address", strings.TrimSpace(addr))
	return query.TryOneDocument(resultPtr)
}

func ConfirmEmail(confirmationCode string, resultPtr interface{}) (email string, confirmed bool, err error) {
	query := Config.Collection.Filter("Email.ConfirmationCode", confirmationCode)
	found, err := query.TryOneDocument(resultPtr)
	if !found {
		return "", false, err
	}
	user := FromDocument(resultPtr)

	for i := range user.Email {
		if user.Email[i].ConfirmationCode.Get() == confirmationCode {
			user.Email[i].Confirmed.SetNowUTC()
			email = user.Email[i].Address.Get()
		}
	}
	errs.Assert(email != "", "ConfirmationCode has to be found")

	err = user.Save()
	if err != nil {
		return "", false, err
	}

	return email, true, nil
}

///////////////////////////////////////////////////////////////////////////////
// Session functions

const ContextCacheKey = "github.com/ungerik/go-start/user.ContextCacheKey"

func Login(session *view.Session, userDoc mongo.Document) {
	LoginID(session, userDoc.ObjectId().Hex())
}

func Logout(session *view.Session) {
	session.DeleteID()
}

func LoginID(session *view.Session, id string) {
	session.SetID(id)
}

func LoginEmailPassword(session *view.Session, email, password string) (emailPasswdMatch bool, err error) {
	email = strings.TrimSpace(email)
	var user User
	found, err := WithEmail(email, &user)
	if !found {
		return false, err
	}
	if !user.EmailPasswordMatch(email, password) {
		return false, nil
	}
	Login(session, &user)
	return true, nil
}

func OfSession(session *view.Session, resultPtr interface{}) (found bool, err error) {
	id := session.ID()
	if id == "" {
		return false, nil
	}
	return WithID(id, resultPtr)
}
