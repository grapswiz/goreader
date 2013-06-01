package goreader

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Auth struct {
	LoginUrl  string
	LogoutUrl string
	LoggedIn  bool
}

func (auth Auth) toJson() string {
	b, err := json.Marshal(auth)
	if err != nil {
		return "{}"
	}
	return string(b)
}

type User struct {
	Admin     bool
	CreatedAt time.Time
}

func createUser(c appengine.Context, key string, admin bool) error {
	user := User{admin, time.Now()}
	_, err := datastore.Put(c, datastore.NewKey(c, "User", key, 0, nil), &user)

	return err
}

func existsUser(c appengine.Context, key string) bool {
	k := datastore.NewKey(c, "User", key, 0, nil)
	q := datastore.NewQuery("User").
		Filter("__key__ =", k).
		Limit(1).
		KeysOnly()
	for t := q.Run(c); ; {
		var x User
		key, err := t.Next(&x)
		if err != nil {
			break
		}
		if key != nil {
			return true
		} else {
			break
		}
	}
	return false
}

func authHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	u := user.Current(c)
	var auth Auth
	auth.LoginUrl, _ = user.LoginURL(c, "/")
	auth.LogoutUrl, _ = user.LogoutURL(c, "/")
	if u == nil {
		auth.LoggedIn = false
	} else {
		auth.LoggedIn = true
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(rw, "%s", auth.toJson())

	if u == nil {
		return
	}
	if existsUser(c, u.Email) {
		return
	}
	createUser(c, u.Email, u.Admin)
}
