## Tutorial with login and user administration

Download, build and run example:

	go get github.com/ungerik/go-start/examples/FullTutorial
	go install github.com/ungerik/go-start/examples/FullTutorial && FullTutorial

Project directory structure and default file names:

* models/
* static/
	* images/
	* css/
* templates/
* views/
	* admin/
	* root/
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

Render CSS as template with a dynamic color-scheme context together with HTML5 boilerplate CSS normalization:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/root/css.go)

	CSS = NewHTML5BoilerplateCSSTemplate(
		TemplateContext(models.NewColorScheme()),
		"css/common.css",
		"css/style.css",
	)

Static views, using own function TitleBar() to create repeating content:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/root/profile.go#L26)

	view = DIV("row",
		DIV("cell right-border",
			TitleBar("My Profile"),
			DIV("main-content",
				H3("Email: ", email),
				...
			),
		),
	)

Printf instead of template-rendering:

	views := Views{
		H2(user.Name.String()),
		Printf("Email confirmed: %v", user.EmailPasswordConfirmed()),
	}

Dynamic views return a view object for a request context:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/shared.go#L11)

	DynamicView(
		func(ctx *Context) (View, error) {
			if ctx.Request.RequestURI == "/" {
				// return nil, so nothing will be rendered
				return nil, nil
			}
			return A("/", HTML("&larr; Back to the homepage")), nil
		},
	)

Render views use a writer to render the response for a request:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/admin/export-emails.go)

	RenderView(
		func(ctx *Context) (err error) {
			// Download instead of display:
			ctx.Response.ContentDispositionAttachment("emails.csv")
			i := models.Users.Iterator()
			var u models.User
			for i.Next(&u) {
				if len(u.Email) > 0 {
					// Context.Response is a writer
					ctx.Response.Printf("%s, %s, %s \n", u.Email[0].Address, u.Name.First.String(), u.Name.Last.String())
				}
			}
			return i.Err()
		},
	),

Menu is a high level struct that creates a list with links as items:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/admin/shared.go#L34)

	&Menu{
		Class:           "menu",
		ItemClass:       "menu-item",
		ActiveItemClass: "active",
		BetweenItems:    " &nbsp;/&nbsp; ",
		Items: []LinkModel{
			NewPageLink(&Admin, "Dashboard"),
			NewPageLink(&Admin_Users, "Users"),
			NewPageLink(&Admin_Images, "Images"),
		},
	}

Simple submit button without form content:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/admin/admin.go#L11)

	&Form{
		FormID:            "clearcaches",
		SubmitButtonText:  "Clear Page Caches",
		OnSubmit: func(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
			ClearAllCaches()
			return "All caches cleared", nil, nil
		},
	}

Auto generated form for a model:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/admin/user_x.go#L55)

	var userModel models.User
	models.Users.DocumentWithID(id, &userModel)

	&Form{
		SubmitButtonText: "Save User Data",
		FormID:           "user" + userModel.ID.Hex(),
		GetModel:         FormModel(&userModel),
		OnSubmit:         OnFormSubmitSaveModelAndRedirect(Admin_Users),
	}

More form features:
(https://github.com/ungerik/go-start/blob/master/examples/FullTutorial/views/admin/user_x.go#L33)

	type passwordFormModel struct {
		Current_password model.Password
		New_password     model.String `model:"minlen=6"`
	}

	&Form{
		SubmitButtonText:  "Save password and mark email as confirmed",
		SubmitButtonClass: "button",
		FormID:            "password",
		DisabledFields:    []string{"Current_password"},
		GetModel: func(form *Form, ctx *Context) (interface{}, error) {
			return &passwordFormModel{Current_password: userModel.Password}, nil
		},
		OnSubmit: func(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
			m := formModel.(*passwordFormModel)
			userModel.Password.SetHashed(m.New_password.Get())
			userModel.ConfirmEmailPassword()
			return "", StringURL("."), userModel.Save()
		},
	}