go-start is a high level web-framework for Go, like Django for Python or Rails for Ruby.

Installation:
go get github.com/ungerik/go-start

Documentation:
http://godoc.org/github.com/ungerik/go-start

Tutorial with user login and administration:
https://github.com/ungerik/go-start/tree/master/examples/FullTutorial

Presentation Slides:
http://speakerdeck.com/u/ungerik/

First real world application:
http://startuplive.in/

Copyright (c) 2012 Erik Unger
MIT License See: LICENSE file


Intro:
======

Features:

* HTML views can be defined in Go syntax
* Optional template system
* HTML5 Boilerplate page template (Mustache template, will be changed to Go v1 template)
* Unified data model for forms and databases
* Data models are simple Go structs
* MongoDB as default database
* User management/authentication out of the box
* Additional packages for
	* Email (Google Mail defaults): http://github.com/ungerik/go-mail
	* Gravatar: http://github.com/ungerik/go-gravatar
	* RSS parsing: http://github.com/ungerik/go-rss
	* Amiando event management: http://github.com/ungerik/go-amiando
		(used by http://startuplive.in)
* Sublime Text 2 Snippets
	* See sublime-snippets/


The case for Go:
================

https://gist.github.com/3731476
	

Views:
======

The philosophy for creating HTML views is (unlike Rails/Django) that you should
not have to learn yet another language to be able to write templates.
There are several very simple template languages out there that reduce program
code like logic within the template, but it’s still yet another syntax to learn.

In go-start the HTML structure of a page is represented by a structure of
type safe Go objects.
It should feel like writing HTML but using the syntax of Go.
And no, it has nothing to do with the mess of intertwined markup and code in PHP.

Example of a static view:

	view := Views{
		DIV("myclass",
			H1("Example HTML structure"),
			P("This is a paragraph"),
			P(
				HTML("Some unescaped HTML:<br/>"),
				Printf("The number of the beast: %d", 666),
				Escape("Will be escaped: 666 < 999"),
			),
			A_blank("http://go-lang.org", "A very simple link"),
		),
		HR(),
		PRE("	<- pre formated text, followed by a list:"),
		UL("red", "green", "blue"),
		&Template{
			Filename: "mytemplate.html",
			GetContext: func(requestContext *Context) (interface{}, error) {
				return map[string]string{"Key": "Value"}, nil
			},
		},
	}

Example of a dynamic view:

	view := DynamicView(
		func(context *Context) (view View, err error) {
			var names []string
			i := models.Users.Sort("Name.First").Sort("Name.Last").Iterator();
			for doc := i.Next(); doc != nil; doc = i.Next() {
				names = append(names, doc.(*models.User).Name.String())
			}
			if i.Err() != nil {
				return nil, i.Err()
			}			
			return &List{	// List = higher level abstraction, UL() = shortcut
				Class: "my-ol",
				Ordered: true,
				Model: EscapeStringsListModel(names),
			}, nil
		},
	)

Beside DynamicView there is also a ModelView. It takes a model.Iterator
and creates a dynamic view for every iterated data item:

	view := &ModelView{
		GetModelIterator: func(context *Context) model.Iterator {
			return models.Users.Sort("Name.First").Sort("Name.Last").Iterator()
		},
		GetModelView: func(model interface{}, context *Context) (view View, err error) {
			user := model.(*models.User)
			return PrintfEscape("%s, ", user.Name), nil
		},
	}


Pages and URLs:
===============

	Homepage := &Page{
		OnPreRender: func(page *Page, context *Context) (err error) {
			context.Data = &PerPageData{...} // Set global page data at request context
		},
		WriteTitle: func(context *Context, writer io.Writer) (err error) {
			writer.Write([]byte(context.Data.(*PerPageData).DynamicTitle))
			return nil
		},
		CSS:          HomepageCSS,
		WriteHeader:  RSS("go-start.org RSS Feed", &RssFeed)
		WriteScripts: PageWriters(
			Config.Page.DefaultWriteScripts,
			JQuery,   // jQuery/UI is built-in
			JQueryUI,
			JQueryUIAutocompleteFromURL(".select-username", IndirectURL(&API_Usernames), 2),
			GoogleAnalytics(GoogleAnalyticsID), // Google Analytics is built-in
		)
		Content: Views{},
	}


	Admin_Auth := NewBasicAuth("go-start.org", "admin", "password123")

	func Paths() *ViewPath {
		return &ViewPath{View: Homepage, Sub: []ViewPath{                           // /
			{Name: "style.css", View: HomepageCSS},                             // /style.css
			{Name: "feed", View: RssFeed},                                      // /feed/
			{Name: "admin", View: Admin, Auth: Admin_Auth, Sub: []ViewPath{     // /admin/
				{Name: "user", Args: 1, View: Admin_User, Auth: Admin_Auth}, // /admin/user/<USER_ID>/
			}},
			{Name: "api", Sub: []ViewPath{                                      // 404 because no view defined
				{Name: "users.json", View: API_Usernames},                  // /api/users.json
			}},
		}
	}

	view.Init("go-start.org", CookieSecret, "pkg/myproject", "pkg/gostart") // Set site name, cookie secret and static paths
	view.Config.RedirectSubdomains = []string{"www"}     // Redirect from www.
	view.Config.Page.DefaultMetaViewport = "width=960px" // Page width for mobile devices
	view.RunConfigFile(Paths(), "run.config")            // Run server with path structure and values from config file




Models:
=======

Data is abstacted as models. The same model abstraction and data validation is
used for HTML forms and for databases. So a model can be loaded from a database,
displayed as an HTML form and saved back to the database after submit.
This is not always a good practice, but it shows how easy things can be.

A model is a simple Go struct that uses gostart/model types as struct members.
Custom model wide validation is done by adding a Validate() method to the
struct type:

	type SignupFormModel struct {
		Email     model.Email    `gostart:"required"`
		Password1 model.Password `gostart:"required|label=Password|minlen=6"`
		Password2 model.Password `gostart:"label=Repeat password"`
	}

	func (self *SignupFormModel) Validate(metaData model.MetaData) []*model.ValidationError {
		if self.Password1 != self.Password2 {
			return model.NewValidationErrors(os.NewError("Passwords don't match"), metaData)
		}
		return model.NoValidationErrors
	}


Here is how a HTML form is created that displays input fields for the SignupFormModel:

	form := &Form{
		ButtonText: "Signup",
		FormID:     "user_signup",
		GetModel: func(form *Form, context *Context) (interface{}, error) {
			return &SignupFormModel{}, nil
		},
		OnSubmit: func(form *Form, formModel interface{}, context *Context) (err error) {
			m := formModel.(*SignupFormModel)
			// ... create user in db and send confirmation email ...
			return err
		},
	}


MongoDB is the default database of go-start utilizing Gustavo Niemeyer's
great lib mgo (http://labix.org/mgo).

Mongo collections and queries are encapsulated to make them compatible with the
go-start data model concept, and a little bit easier to use.

Example of a collection and document struct:

	var ExampleDocs *mongo.Collection = mongo.NewCollection("exampledocs", (*ExampleDoc)(nil))

	type ExampleDoc struct {
		mongo.DocumentBase `bson:",inline"`                 // Give it a Mongo ID
		Person             mongo.Ref  `gostart:"to=people"` // Mongo ID ref to a document in "people" collection
		LongerText         model.Text `gostart:"rows=5|cols=80|maxlen=400"`
		Integer            model.Int  `gostart:"min=1|max=100"`
		Email              model.Email    // Normalization + special treament in forms
		PhoneNumber        model.Phone    // Normalization + special treament in forms
		Password           model.Password // Hashed + special treament in forms
		SubDoc             struct {
			Day       model.Date
			Drinks    []mongo.Choice `gostart:"options=Beer,Wine,Water"` // Mongo array of strings
			RealFloat model.Float    `gostart:"valid" // Must be a real float value, not NaN or Inf
		}
	}

Example query:

	i := models.Users.Filter("Name.Last", "Smith").Sort("Name.First").Iterator();
	for doc := i.Next(); doc != nil; doc = i.Next() {
		user := doc.(*models.User)
		// ...
	}
	// Err() returns any error after Next() returned nil:
	if i.Err() != nil {
		panic(i.Err())
	}

A new mongo.Document is always created by the corresponding collection object
to initialize it with meta information about its collection.
This way it is possible to implement Save() or Remove() methods for the document.

Example for creating, modifying and saving a document:

	user := models.Users.NewDocument().(*models.User)

	user.Name.First.Set("Erik")
	user.Name.Last.Set("Unger")

	err := user.Save()


