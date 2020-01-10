package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/EliasPlesiotis/spaceadvent/website/cookie"
	"github.com/EliasPlesiotis/spaceadvent/website/database"
	md "github.com/EliasPlesiotis/spaceadvent/website/middleware"
)

var (
	id = 0
)

// Index Page
func Index(w http.ResponseWriter, r *http.Request) {
	value, err := cookie.ReadCookie(r, "session")
	_, found := database.Find(database.User{UserName: value["username"]})

	if err != nil || !found {
		md.Render(w, "index", database.Data{
			Logged:   false,
			Messages: []database.Message{},
		})
		return
	}

	msgs := database.GetMsg(value["username"])
	md.Render(w, "welcome", database.Data{
		Logged:   true,
		Username: value["username"],
		Messages: msgs,
	})
}

// Login Page
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		value, err := cookie.ReadCookie(r, "session")
		_, found := database.Find(database.User{UserName: value["username"]})

		if err == nil && found {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		md.Render(w, "login", database.Data{
			Logged:   false,
			Errors:   nil,
			Messages: []database.Message{},
		})
		return
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		u, found := database.Find(database.User{UserName: username})

		if found {
			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err == nil {
				hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				pwd := string(hash)

				cookie.CreateSession(w, username, pwd)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		md.Render(w, "login", database.Data{
			Logged:   false,
			Errors:   []string{"Wrong username or password"},
			Messages: []database.Message{},
		})
		return
	}
}

// Register Page
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		value, err := cookie.ReadCookie(r, "session")
		_, found := database.Find(database.User{UserName: value["username"]})

		if err == nil && found {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		md.Render(w, "register", database.Data{
			Logged:   false,
			Messages: []database.Message{},
		})
		return

	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		cpassword := r.FormValue("cpassword")

		if password != cpassword {
			md.Render(w, "register", database.Data{
				Logged:   false,
				Errors:   []string{"Could not confirm the password"},
				Messages: []database.Message{},
			})
			return
		}

		_, found := database.Find(database.User{UserName: username})

		if found {
			md.Render(w, "register", database.Data{
				Logged:   false,
				Errors:   []string{"Username Already Taken"},
				Messages: []database.Message{},
			})
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		pwd := string(hash)

		database.Create(username, pwd, email)
		cookie.CreateSession(w, username, pwd)
		database.FirstMsg(username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Logout Page
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie.ClearCookie(w, "session")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Settings Page
func Settings(w http.ResponseWriter, r *http.Request) {
	value, err := cookie.ReadCookie(r, "session")
	_, found := database.Find(database.User{UserName: value["username"]})

	if err != nil || !found {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	msgs := database.GetMsg(value["username"])
	md.Render(w, "settings", database.Data{Logged: true, Messages: msgs})
}

// Messages Page
func Messages(w http.ResponseWriter, r *http.Request) {
	value, err := cookie.ReadCookie(r, "session")
	_, found := database.Find(database.User{UserName: value["username"]})

	if err != nil || !found {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	msgs := database.GetMsg(value["username"])
	md.Render(w, "messages", database.Data{
		Logged:   true,
		Messages: msgs,
	})
}

// Replace Info of User
func Replace(w http.ResponseWriter, r *http.Request) {
	value, err := cookie.ReadCookie(r, "session")
	if err != nil {
		md.Render(w, "index", database.Data{
			Logged:   false,
			Messages: []database.Message{},
		})
		return
	}

	r.ParseForm()

	msgs := database.GetMsg(value["username"])

	oldUsername := value["username"]
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	cpassword := r.FormValue("cpassword")

	if password != cpassword {
		md.Render(w, "register", database.Data{
			Logged:   false,
			Errors:   []string{"Could not confirm the password"},
			Messages: msgs,
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	pwd := string(hash)

	database.Update(oldUsername, username, pwd, email)

	cookie.CreateSession(w, username, pwd)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Sent message to Kuiper
func Sent(w http.ResponseWriter, r *http.Request) {
	value, err := cookie.ReadCookie(r, "session")
	u, found := database.Find(database.User{UserName: value["username"]})

	if err != nil || !found {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	database.StoreMsg(u.UserName, r.FormValue("msg"))

	http.Redirect(w, r, "/messages", http.StatusSeeOther)
}

// DeleteMe deletes the User
func DeleteMe(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		value, err := cookie.ReadCookie(r, "session")
		_, found := database.Find(database.User{UserName: value["username"]})

		if err != nil || !found {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		msgs := database.GetMsg(value["username"])

		md.Render(w, "delete", database.Data{
			Logged:   true,
			Messages: msgs,
		})

	} else if r.Method == "POST" {
		value, err := cookie.ReadCookie(r, "session")
		u, found := database.Find(database.User{UserName: value["username"]})

		if err != nil || !found {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		r.ParseForm()
		username := value["username"]
		password := r.FormValue("password")

		if found {
			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err == nil {
				database.DeleteUser(username)
				database.DeleteMsg(username)
				cookie.ClearCookie(w, "session")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		md.Render(w, "delete", database.Data{
			Logged:   false,
			Errors:   []string{"Wrong password"},
			Messages: []database.Message{},
		})
		return
	}
}

// ResetMsg delete the existing messages with Kuiper
func ResetMsg(w http.ResponseWriter, r *http.Request) {
	value, err := cookie.ReadCookie(r, "session")
	u, found := database.Find(database.User{UserName: value["username"]})

	if err != nil || !found {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	username := u.UserName

	fmt.Println(username)

	database.DeleteMsg(username)
	database.FirstMsg(username)
	http.Redirect(w, r, "/messages", http.StatusSeeOther)
}
