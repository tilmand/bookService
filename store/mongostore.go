package store

import (
	"bookService/config"
	"bookService/model"
	"fmt"
	"log"

	ai "github.com/night-codes/mgo-ai"
	"gopkg.in/mgo.v2"
)

type obj map[string]interface{}

type MongoStore struct {
	conn            *mgo.Database
	BooksRepository *BooksRepository
	UsersRepository *UsersRepository
}

type Database interface {
	GetByLogin(login string) (*model.User, error)
	GetByID(ID uint64) (*model.User, error)
	Insert(user model.User) error
	SaveRecoveryToken(userID uint64, token string) error
	VerifyRecoveryToken(token string) (uint64, error)
}

func NewMongoStore(cfg *config.Config) (*MongoStore, error) {
	dsn := fmt.Sprintf("mongodb://%s:%d", cfg.Mongo.Host, cfg.Mongo.Port)
	session, err := mgo.Dial(dsn)
	if err != nil {
		return nil, err
	}
	db := session.DB(cfg.Mongo.Database)

	store := &MongoStore{
		conn: db,
	}

	initStore(db)
	store.BooksRepository = store.Books()
	store.UsersRepository = store.Users()

	return store, nil
}

func initStore(db *mgo.Database) {
	if err := createCollections(db); err != nil {
		log.Printf("Err: %v", err)
	}

	if err := insertDocuments(db); err != nil {
		log.Printf("Err: %v", err)
	}

	log.Println("Data successfully inserted into MongoDB collection.")
}

func createCollections(db *mgo.Database) error {
	index := mgo.Index{
		Key:    []string{"_id"},
		Unique: true,
	}
	err := db.C("books").EnsureIndex(index)
	if err != nil {
		return err
	}
	index2 := mgo.Index{
		Key:    []string{"_id"},
		Unique: true,
	}
	err = db.C("users").EnsureIndex(index2)

	return err
}

func insertDocuments(db *mgo.Database) error {
	var err error
	books := []model.Book{
		{Name: "Book 1", AuthorID: 4},
		{Name: "Book 2", AuthorID: 4},
		{Name: "Book 3", AuthorID: 3},
	}

	for _, book := range books {
		ai.Connect(db.C("ai"))
		book.ID = ai.Next("books")
		err = db.C("books").Insert(book)
		if err != nil {
			log.Printf("Err: %v", err)

			return err
		}
	}
	users := []model.User{
		{Login: "Book 1", Role: "Author"},
		{Login: "Book 2", Role: "Author"},
		{Login: "Book 3", Role: "Author"},
	}

	for _, user := range users {
		ai.Connect(db.C("ai"))
		user.ID = ai.Next("users")
		err = db.C("users").Insert(user)
		if err != nil {
			log.Printf("Err: %v", err)

			return err
		}
	}

	index := mgo.Index{
		Key:    []string{"login"},
		Unique: true,
	}
	err = db.C("users").EnsureIndex(index)
	if err != nil {
		log.Printf("Err: %v", err)

		return err
	}

	return err
}

func (s *MongoStore) Books() *BooksRepository {
	if s.BooksRepository == nil {
		s.BooksRepository = NewBooksRepository(s)
	}

	return s.BooksRepository
}

func (s *MongoStore) Users() *UsersRepository {
	if s.UsersRepository == nil {
		s.UsersRepository = NewUsersRepository(s)
	}

	return s.UsersRepository
}
