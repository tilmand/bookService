package store

import (
	"bookService/model"
	"fmt"
	"log"

	ai "github.com/night-codes/mgo-ai"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionUsers = "users"
)

type (
	UsersRepository struct {
		store          *MongoStore
		collectionName string
	}
)

func NewUsersRepository(store *MongoStore) *UsersRepository {
	return &UsersRepository{
		store:          store,
		collectionName: collectionUsers,
	}
}

func (r *UsersRepository) GetAll() ([]model.User, error) {
	results := []model.User{}
	err := r.store.conn.C(collectionUsers).Find(obj{}).All(&results)
	if err != nil {
		log.Println("GetAll Find err: ", err)
	}

	return results, err
}

func (r *UsersRepository) Find(userID uint64) (model.User, error) {
	result := model.User{}
	err := r.store.conn.C(collectionUsers).FindId(userID).One(&result)
	if err != nil {
		log.Println("Find FindId err: ", err)

		return model.User{}, err
	}

	return result, nil
}

func (r *UsersRepository) Insert(item model.User) error {
	ai.Connect(r.store.conn.C("ai"))
	item.ID = ai.Next(collectionUsers)
	err := r.store.conn.C(collectionUsers).Insert(item)
	if err != nil {
		log.Println("Insert Insert err: ", err)
	}

	return err
}

func (r *UsersRepository) Update(item model.User) error {
	err := r.store.conn.C(collectionUsers).UpdateId(item.ID, obj{"$set": item})
	if err != nil {
		log.Println("Update UpdateId err: ", err)
	}

	return err
}

func (r *UsersRepository) Delete(ID uint64) error {
	err := r.store.conn.C(collectionUsers).Remove(obj{"_id": ID})
	if err != nil {
		log.Println("Delete Remove err: ", err)
	}

	return err
}

func (r *UsersRepository) GetByLogin(login string) (*model.User, error) {
	result := &model.User{}
	err := r.store.conn.C(collectionUsers).Find(obj{"login": login}).One(result)
	if err != nil {
		log.Println("GetByLogin Find err: ", err)
	}

	return result, err
}

func (r *UsersRepository) SaveRecoveryToken(userID uint64, recoveryToken string) error {
	user, err := r.Find(userID)
	if err != nil {
		log.Println("SaveRecoveryToken Find err: ", err)

		return err
	}

	user.RecoveryToken = recoveryToken

	if err := r.Update(user); err != nil {
		log.Println("SaveRecoveryToken Update err: ", err)

		return err
	}

	return nil
}

func (r *UsersRepository) SetPassword(userID uint64, hashedPassword string) error {
	user, err := r.Find(userID)
	if err != nil {
		log.Println("SetPassword Find err: ", err)

		return err
	}

	user.Password = hashedPassword

	if err := r.Update(user); err != nil {
		log.Println("SetPassword Update err: ", err)

		return err
	}

	return nil
}

func (r *UsersRepository) VerifyRecoveryToken(recoveryToken string) (uint64, error) {
	user, err := r.getUserByRecoveryToken(recoveryToken)
	if err != nil {
		log.Println("VerifyRecoveryToken getUserByRecoveryToken err: ", err)

		return 0, err
	}

	return user.ID, nil
}

func (r *UsersRepository) getUserByRecoveryToken(recoveryToken string) (*model.User, error) {
	user := &model.User{}
	err := r.store.conn.C(collectionUsers).Find(bson.M{"recoveryToken": recoveryToken}).One(user)
	if err != nil {
		log.Println("getUserByRecoveryToken Find err: ", err)
		if err == mgo.ErrNotFound {
			return nil, fmt.Errorf("user not found")
		}

		return nil, err
	}

	return user, nil
}
