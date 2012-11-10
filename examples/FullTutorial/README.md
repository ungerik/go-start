## Tutorial with login and user administration

Download, build and run example:

	go get github.com/ungerik/go-start/examples/FullTutorial
	go install github.com/ungerik/go-start/examples/FullTutorial && FullTutorial


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

