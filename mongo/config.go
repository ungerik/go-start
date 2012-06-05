package mongo

import (
	"fmt"
	"launchpad.net/mgo"
)

var database *mgo.Database

var collections map[string]*Collection = map[string]*Collection{}

type Configuration struct {
	Host                string
	Database            string
	User                string
	Password            string
	Safe                mgo.Safe
	CheckQuerySelectors bool
}

var Config Configuration = Configuration{
	Safe:                mgo.Safe{FSync: true, J: true},
	CheckQuerySelectors: true,
}

func InitLocalhost(database, user, password string) (err error) {
	Config.Database = database
	Config.User = user
	Config.Password = password
	return Init()
}

func Init() (err error) {
	login := ""
	if Config.User != "" {
		login = Config.User + ":" + Config.Password + "@"
	}

	host := "localhost"
	if Config.Host != "" {
		host = Config.Host
	}

	// http://goneat.org/pkg/launchpad.net/mgo/#Session.Mongo
	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	url := fmt.Sprintf("mongodb://%s%s/%s", login, host, Config.Database)

	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	session.SetSafe(&Config.Safe)

	database = session.DB(Config.Database)

	for _, collection := range collections {
		collection.collection = database.C(collection.Name)
	}

	return nil
}

func Close() {
	if database.Session != nil {
		database.Session.Close()
	}
}
