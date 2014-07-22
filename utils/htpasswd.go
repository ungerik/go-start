package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// SHA1Base64String hashes data with SHA1 and encodes the result as base64 string.
func SHA1Base64String(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// ReadHtpasswdFile returns a map of usernames to base64 encoded SHA1 hashed passwords
func ReadHtpasswdFile(filename string) (userPass map[string]string, modified time.Time, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, time.Time{}, err
	}
	modified = info.ModTime()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ':'
	csvReader.Comment = '#'
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, time.Time{}, err
	}

	userPass = make(map[string]string)

	for _, record := range records {
		username := record[0]
		password := record[1]
		if len(password) < 5 {
			return nil, time.Time{}, errors.New("Invalid password")
		}
		if password[:5] != "{SHA}" {
			return nil, time.Time{}, errors.New("Unsupported password format, must be SHA1")
		}
		userPass[username] = password[5:]
	}

	return userPass, modified, nil
}

func ReadHtpasswdFileUsers(filename string) (users []string, err error) {
	userPass, _, err := ReadHtpasswdFile(filename)
	if err != nil {
		return nil, err
	}
	users = make([]string, 0, len(userPass))
	for user := range userPass {
		users = append(users, user)
	}
	return users, nil
}

func WriteHtpasswdFile(filename string, userPass map[string]string) error {
	var buf bytes.Buffer
	for user, pass := range userPass {
		_, err := fmt.Fprintf(&buf, "%s:{SHA}%s\n", user, pass)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0660)
}
