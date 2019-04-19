package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/globalsign/mgo/bson"
	"github.com/go-bongo/bongo"
)

var connection *bongo.Connection

func Init() error {
	config := &bongo.Config{
		ConnectionString: getConnectionString(),
		Database:         "payword-backend",
	}

	pConnection, err := bongo.Connect(config)

	if err != nil {
		return err
	}

	connection = pConnection
	return nil
}

func getConnectionString() string {
	if os.Getenv("APP_ENV") == "production" {
		return "mongo"
	}
	return "localhost"

}

// User repräsentiert den Nutzer in der Datenbank
type User struct {
	bongo.DocumentBase `bson:",inline"` // Metadata
	Username           string
	Password           string
}

type Counter struct {
	Username string
	Value    int
}

type Category struct {
	Name        string
	Usercounter []Counter
}

type Game struct {
	bongo.DocumentBase `bson:",inline"` // Metadata
	Name               string
	Admin              string
	Categories         []Category
	Members            []string
}

type Mail struct {
	bongo.DocumentBase `bson:",inline"` // Metadata
	Username           string
	Mail               string // Mail adress of the user
}

func (mail *Mail) Save() error {
	fmt.Println(mail.GetId().String())
	return connection.Collection("mails").Save(mail)
}

func GetMailAdress(username string) (*Mail, error) {
	mail := &Mail{}
	err := connection.Collection("mails").FindOne(bson.M{"username": username}, mail)

	if err != nil {
		return nil, err
	}

	return mail, nil
}

func (game *Game) Save() error {
	return connection.Collection("games").Save(game)
}

func (game *Game) CategoryExists(categoryname string) bool {
	for _, category := range game.Categories {
		if category.Name == categoryname {
			return true
		}
	}

	return false
}

func (game *Game) AddCategory(name string) error {
	// create Counter for each user
	var values []Counter
	for _, username := range game.Members {
		values = append(values, Counter{username, 0})
	}

	game.Categories = append(game.Categories, Category{name, values})

	return game.Save()
}

func (game *Game) RemoveCategory(name string) error {
	for key, category := range game.Categories {
		if category.Name == name {
			game.Categories = append(game.Categories[:key], game.Categories[key+1:]...)
			return game.Save()
		}
	}

	return errors.New("Not Found")
}

func (game *Game) Increment(categoryName string, username string) error {
	for keyCategory, category := range game.Categories {
		if category.Name == categoryName {
			for keyCounter, counter := range category.Usercounter {
				if counter.Username == username {
					game.Categories[keyCategory].Usercounter[keyCounter].Value++
					return game.Save()
				}
			}
		}
	}

	return errors.New("Not Found")
}

func (game *Game) Decrement(categoryName string, username string) error {
	for keyCategory, category := range game.Categories {
		if category.Name == categoryName {
			for keyCounter, counter := range category.Usercounter {
				if counter.Username == username {
					game.Categories[keyCategory].Usercounter[keyCounter].Value--
					return game.Save()
				}
			}
		}
	}

	return errors.New("Not Found")
}

func (game *Game) AddUser(username string) error {
	game.Members = append(game.Members, username)

	for key := range game.Categories {
		game.Categories[key].Usercounter = append(
			game.Categories[key].Usercounter,
			Counter{username, 0})
	}

	return game.Save()
}

func (game *Game) RemoveUser(username string) error {
	// Delete User from Members
	for key, member := range game.Members {
		if member == username {
			game.Members = append(game.Members[:key], game.Members[key+1:]...)
			break
		}
	}

	// Delete User from Categories
	for keyCategory, category := range game.Categories {
		for keyCounter, usercounter := range category.Usercounter {
			if usercounter.Username == username {
				// Delete user from counter array for current category
				game.Categories[keyCategory].Usercounter =
					append(
						game.Categories[keyCategory].Usercounter[:keyCounter],
						game.Categories[keyCategory].Usercounter[keyCounter+1:]...)
			}
		}
	}

	return game.Save()
}

func (game *Game) FindNewAdmin() string {
	for _, member := range game.Members {
		if member != game.Admin {
			return member
		}
	}

	return ""
}

func (game *Game) Delete() error {
	return connection.Collection("games").DeleteDocument(game)
}

func GameExists(name string) (bool, error) {
	_, err := FindGameByName(name)
	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			return false, nil
		}

		return false, err

	}

	return true, nil
}

func FindGameByName(name string) (*Game, error) {
	game := &Game{}
	err := connection.Collection("games").FindOne(bson.M{"name": name}, game)

	if err != nil {
		return &Game{}, err
	}

	return game, nil
}

func FindAllGames() []Game {
	var results *bongo.ResultSet
	results = connection.Collection("games").Find(bson.M{})

	var games []Game

	game := &Game{}

	for results.Next(game) {
		games = append(games, *game)
	}

	return games
}

func UserExists(username string) (bool, error) {
	_, err := FindUserByName(username)
	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			return false, nil
		}

		return false, err

	}

	return true, nil
}

func SaveUser(user *User) error {
	return connection.Collection("user").Save(user)
}

func FindAllUsers() []User {
	var results *bongo.ResultSet
	results = connection.Collection("user").Find(bson.M{})

	var users []User

	user := &User{}

	for results.Next(user) {
		users = append(users, *user)
	}

	return users
}

func FindUserByName(username string) (*User, error) {
	user := &User{}
	err := connection.Collection("user").FindOne(bson.M{"username": username}, user)

	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func DeleteUser(user *User) error {
	return connection.Collection("user").DeleteDocument(user)
}

func DeleteMail(mail *Mail) error {
	return connection.Collection("mails").DeleteDocument(mail)
}

func DeleteGame(game *Game) error {
	return connection.Collection("games").DeleteDocument(game)
}
