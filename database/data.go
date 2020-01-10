package database

import (
	"fmt"

	kuiper "github.com/EliasPlesiotis/spaceadvent/website/Kuiper"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	// DB the database
	DB *gorm.DB
)

func init() {
	var err error
	DB, err = gorm.Open("sqlite3", "./data.db")

	if err != nil {
		fmt.Println(err)
	}

	// Migrate the schema
	DB.AutoMigrate(&User{}, &Message{})
}

// Find user in database
func Find(u User) (User, bool) {
	var user User
	isthere := false
	DB.Where(&u).First(&user)
	if user.UserName == u.UserName {
		isthere = true
	}

	return user, isthere
}

// Create new user
func Create(username, password, email string) {
	DB.Create(&User{UserName: username, Password: password, Email: email})
}

// Update the User
func Update(oldUsername string, username string, password string, email string) {
	oldu := User{UserName: oldUsername}
	u, _ := Find(oldu)

	u.UserName = username
	u.Password = password
	u.Email = email
	u.Messages = oldu.Messages

	DB.Save(&u)
}

// DeleteUser deletes users
func DeleteUser(username string) {
	DB.Delete(&User{UserName: username})
}

// FirstMsg to user
func FirstMsg(username string) {
	msg1 := Message{
		UserID: username,
		Mine:   false,
		Text:   "Hi " + username + ".I am Kuiper an artificial inteligence to help you colonise the universe.",
	}

	msg2 := Message{
		UserID: username,
		Mine:   false,
		Text:   "In order to understand your messages use the keywords [shoot, gravity, move, bullets, shield]",
	}

	DB.Create(&msg1)
	DB.Create(&msg2)

}

// GetMsg for a User
func GetMsg(username string) []Message {
	var msgs []Message
	DB.Raw("SELECT * FROM messages WHERE user_id = ?;", username).Scan(&msgs)
	return msgs
}

// DeleteMsg deletes the messages of a user
func DeleteMsg(username string) {
	DB.Exec("delete from messages where user_id = ?;", username)
}

// StoreMsg to database
func StoreMsg(username, text string) {
	msg1 := Message{
		UserID: username,
		Mine:   true,
		Text:   text,
	}

	resp := kuiper.Responde(text)

	msg2 := Message{
		UserID: username,
		Mine:   false,
		Text:   resp,
	}

	DB.Create(&msg1)
	DB.Create(&msg2)
}

// Data for the templates
type Data struct {
	Logged   bool
	Username string
	Messages []Message
	Errors   []string
}
