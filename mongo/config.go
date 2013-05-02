package mongo

import (
	"fmt"
	"labix.org/v2/mgo"
)

var Config = Configuration{
	Safe: mgo.Safe{ // be conservative...
		W:     1,
		FSync: false,
		J:     true,
	},
	CheckQuerySelectors: true,
}

var Database *mgo.Database

var Collections = map[string]*Collection{}

type Configuration struct {
	Host                string
	Database            string
	User                string
	Password            string
	Safe                mgo.Safe
	CheckQuerySelectors bool
}

func (self *Configuration) Name() string {
	return "mongo"
}

func (self *Configuration) Init() error {
	login := ""
	if Config.User != "" {
		login = Config.User + ":" + Config.Password + "@"
	}

	host := "localhost"
	if Config.Host != "" {
		host = Config.Host
	}

	// http://goneat.org/pkg/labix.org/v2/mgo/#Session.Mongo
	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	url := fmt.Sprintf("mongodb://%s%s/%s", login, host, Config.Database)

	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	session.SetSafe(&Config.Safe)

	Database = session.DB(Config.Database)

	for _, collection := range Collections {
		collection.collection = Database.C(collection.Name)
	}

	return nil
}

func (self *Configuration) Close() error {
	if Database != nil && Database.Session != nil {
		Database.Session.Close()
		Database.Session = nil
		Database = nil
	}
	return nil
}

func InitLocalhost(database, user, password string) (err error) {
	Config.Database = database
	Config.User = user
	Config.Password = password
	return Config.Init()
}
