package store

import (
	"bookService/model"
	"log"

	ai "github.com/night-codes/mgo-ai"
)

const (
	collectionBooks = "books"
)

type (
	BooksRepository struct {
		store          *MongoStore
		collectionName string
	}
)

func NewBooksRepository(store *MongoStore) *BooksRepository {
	return &BooksRepository{
		store:          store,
		collectionName: collectionBooks,
	}
}

func (r *BooksRepository) GetAll() ([]model.Book, error) {
	results := []model.Book{}
	err := r.store.conn.C(collectionBooks).Find(obj{}).All(&results)
	if err != nil {
		log.Println("GetAll Find err: ", err)
	}

	return results, err
}

func (r *BooksRepository) Find(bookID uint64) (model.Book, error) {
	result := model.Book{}
	err := r.store.conn.C(collectionBooks).FindId(bookID).One(&result)
	if err != nil {
		log.Println("Find FindId err: ", err)

		return model.Book{}, err
	}

	return result, nil
}

func (r *BooksRepository) Insert(item model.Book, authorID uint64) error {
	ai.Connect(r.store.conn.C("ai"))
	item.ID = ai.Next(collectionBooks)
	item.AuthorID = authorID
	err := r.store.conn.C(collectionBooks).Insert(item)
	if err != nil {
		log.Println("Insert Insert err: ", err)
	}

	return err
}

func (r *BooksRepository) Update(item model.Book) error {
	err := r.store.conn.C(collectionBooks).UpdateId(item.ID, obj{"$set": item})
	if err != nil {
		log.Println("Update UpdateId err: ", err)
	}

	return err
}

func (r *BooksRepository) Delete(ID uint64) error {
	err := r.store.conn.C(collectionBooks).Remove(obj{"_id": ID})
	if err != nil {
		log.Println("Delete Remove err: ", err)
	}

	return err
}
