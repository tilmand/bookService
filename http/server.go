package http

import (
	"bookService/auth"
	"bookService/store"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var wg sync.WaitGroup

type api struct {
	mongo  *store.MongoStore
	router *gin.Engine
	auth   auth.Middleware

	booksHandler *BooksHandler
	authHandler  *AuthHandler
	db           store.Database
}

func NewServer(mongo *store.MongoStore, auth *auth.Middleware) *api {
	api := &api{
		mongo: mongo,
		auth:  *auth,
	}

	api.router = configureRouter(api)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: api.router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.ListenAndServe()
		if err != nil {
			log.Println("NewServer ListenAndServe err: ", err)
		}
	}()

	return nil
}

func Wait() {
	wg.Wait()
}

func (a *api) Books() *BooksHandler {
	if a.booksHandler == nil {
		a.booksHandler = NewBooksHandler(a)
	}

	return a.booksHandler
}

func (a *api) Auth() *AuthHandler {
	if a.authHandler == nil {
		a.authHandler = NewAuthHandler(a)
	}

	return a.authHandler
}
