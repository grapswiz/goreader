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
	Email     string
	Admin     bool
	CreatedAt time.Time
}

func createUser(c appengine.Context, email string, admin bool) error {
	user := User{email, admin, time.Now()}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "User", nil), &user)

	return err
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

	//TODO uがnil以外の場合、Datastoreを検索して、もしuのuserエンティティが存在しなかったらエンティティを作る
	// entityのkeyにemailを設定したい（検索するときkeyだけで済むように）
	if u == nil {
		return
	}
	createUser(c, u.Email, u.Admin)
}
