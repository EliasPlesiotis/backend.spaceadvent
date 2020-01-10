package cookie

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
)

var (
	key1 []byte
	key2 []byte

	s *securecookie.SecureCookie
)

//Cred DO NOT TOUCH
type Cred struct {
	Key1 string
	Key2 string
}

func init() {
	data := Cred{}

	jsonFile, err := os.Open("./cred.json")
	if err != nil {
		panic(err)
	}

	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, &data)

	key1 = []byte(data.Key1)
	key2 = []byte(data.Key2)

	s = securecookie.New(key1, key2)
}

//ReadCookie get the Cookie Value
func ReadCookie(r *http.Request, name string) (map[string]string, error) {
	var cookie *http.Cookie
	var err error

	if cookie, err = r.Cookie(name); err == nil {
		value := make(map[string]string)
		if err = s.Decode(name, cookie.Value, &value); err == nil {
			return value, err
		}
	}
	return nil, err
}

// CreateSession for Login
func CreateSession(w http.ResponseWriter, username string, password string) {
	var value = make(map[string]string)
	value["username"] = username
	value["password"] = password

	encoded, err := s.Encode("session", value)
	if err == nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: encoded, Path: "/"})
	}
}

//ClearCookie delete the cookie
func ClearCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "0",
		MaxAge: -1}

	http.SetCookie(w, cookie)
}
