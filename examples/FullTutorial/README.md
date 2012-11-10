## Tutorial with login and user administration

Download, build and run example:

	go get github.com/ungerik/go-start/examples/FullTutorial
	go install github.com/ungerik/go-start/examples/FullTutorial && FullTutorial

Project directory structure and default file names:

* models
* static
	* images
	* css
* templates
* views
	* admin
	* root
	* paths.go
* config.json
* main.go


Load JSON config file and initialize packages:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/main.go#L30)

	config.Load("config.json",
		&email.Config,
		&mongo.Config,
		&user.Config,
		&view.Config,
		&media.Config,
		&mongomedia.Config,
	)

User model and mongo collection:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/models/user.go)

	var Users = mongo.NewCollection("users")

	type User struct {
		user.User `bson:",inline"`

		Image  media.ImageRef
		Gender model.Choice `model:"options=,Male,Female"`
	}

Querying the user of the session and authentication:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/admin/auth.go)

	type Admin_Authenticator struct{}

	func (self *Admin_Authenticator) Authenticate(ctx *Context) (ok bool, err error) {
		var u models.User
		found, err := user.OfSession(ctx.Session, &u)
		if !found {
			return false, err
		}
		return u.Admin.Get(), nil
	}

	func init() {
		Admin_Auth = new(Admin_Authenticator)
	}

View paths (URLS):
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/paths.go)

	var (
		CSS                View
		Homepage           ViewWithURL
		...
	)

	func Paths() *ViewPath {
		return &ViewPath{View: Homepage, Sub: []ViewPath{
			media.ViewPath("media"),
			{Name: "style.css", View: CSS},
			...
		}}
	}

Views:

	CSS = NewHTML5BoilerplateCSSTemplate(
		TemplateContext(models.NewColorScheme()),
		"css/common.css",
		"css/style.css",
	)
