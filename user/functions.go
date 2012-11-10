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
func EnsureExists(username, email, password string, admin bool, resultRef interface{}) (exists bool, err error) {
	exists, err = Config.Collection.Filter("Username", username).TryOneDocument(resultRef)
	if err != nil {
		return false, err
	}
	user := FromDocument(resultRef)

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

	if exists {
		err = Config.Collection.UpdateSubDocumentWithID(user.ID, "", user)
	} else {
		err = Config.Collection.InitAndSaveDocument(resultRef.(mongo.Document))
	}

	if err != nil {
		return false, err
	}
	return exists, nil
}

func WithID(id string, resultRef interface{}) (found bool, err error) {
	return Config.Collection.TryDocumentWithID(bson.ObjectIdHex(id), resultRef)
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

func WithEmail(addr string, resultRef interface{}) (found bool, err error) {
	query := Config.Collection.FilterEqualCaseInsensitive("Email.Address", strings.TrimSpace(addr))
	return query.TryOneDocument(resultRef)
}

func ConfirmEmail(confirmationCode string) (userID, email string, confirmed bool, err error) {
	query := Config.Collection.Filter("Email.ConfirmationCode", confirmationCode)
	var user User
	found, err := query.TryOneDocument(&user)
	if !found {
		return "", "", false, err
	}

	for i := range user.Email {
		if user.Email[i].ConfirmationCode.Get() == confirmationCode {
			user.Email[i].Confirmed.SetNowUTC()
			email = user.Email[i].Address.Get()
		}
	}
	errs.Assert(email != "", "ConfirmationCode has to be found")

	err = Config.Collection.UpdateSubDocumentWithID(user.ID, "", &user)
	if err != nil {
		return "", "", false, err
	}

	return user.ID.Hex(), email, true, nil
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

func LoggedIn(session *view.Session) bool {
	return session.ID() != ""
}

func LoginEmailPassword(session *view.Session, email, password string) (emailPasswdMatch bool, err error) {
	// todo: throttle login attempts from same IP or same user to 1 per second, double delay after a few attempts
	// todo: log bad login frequency to detect attacks
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

func OfSession(session *view.Session, resultRef interface{}) (found bool, err error) {
	id := session.ID()
	if id == "" {
		return false, nil
	}
	return WithID(id, resultRef)
}
