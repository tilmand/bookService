package http

import (
	"bookService/auth"
	"bookService/mocks"
	"bookService/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func (f *FakeMongoStore) GetAll() ([]model.Book, error) {
	return []model.Book{}, nil
}

func (f *FakeMongoStore) Find(ID uint64) (model.Book, error) {
	return model.Book{}, nil
}

func (f *FakeMongoStore) Update(book model.Book) error {
	return nil
}

func (f *FakeMongoStore) Delete(ID uint64) error {
	return nil
}

func (f *FakeJWTAuth) ExtractToken(req *http.Request) string {
	return "fake_token"
}

func (f *FakeJWTAuth) Validate(token string) (*auth.AccessClaims, error) {
	return &auth.AccessClaims{BaseClaims: auth.BaseClaims{ID: 123}}, nil
}

func TestGetAllHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBooksHandler := mocks.NewMockBooksHandlerInterface(ctrl)

	mockBooksHandler.EXPECT().GetAll(gomock.Any()).Return()

	req, _ := http.NewRequest("GET", "/books", nil)

	router := gin.Default()
	router.GET("/books", func(c *gin.Context) {
		mockBooksHandler.GetAll(c)
	})
	router.ServeHTTP(httptest.NewRecorder(), req)
}

func TestAddHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBooksHandler := mocks.NewMockBooksHandlerInterface(ctrl)

	mockBooksHandler.EXPECT().Add(gomock.Any()).Return()

	req, _ := http.NewRequest("POST", "/book", nil)

	router := gin.Default()
	router.POST("/book", func(c *gin.Context) {
		mockBooksHandler.Add(c)
	})
	router.ServeHTTP(httptest.NewRecorder(), req)
}

func TestFindHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBooksHandler := mocks.NewMockBooksHandlerInterface(ctrl)

	mockBooksHandler.EXPECT().Find(gomock.Any()).Return()

	req, _ := http.NewRequest("GET", "/book/123", nil)

	router := gin.Default()
	router.GET("/book/:id", func(c *gin.Context) {
		mockBooksHandler.Find(c)
	})
	router.ServeHTTP(httptest.NewRecorder(), req)
}

func TestUpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBooksHandler := mocks.NewMockBooksHandlerInterface(ctrl)

	mockBooksHandler.EXPECT().Update(gomock.Any()).Return()

	req, _ := http.NewRequest("PUT", "/book/123", nil)

	router := gin.Default()
	router.PUT("/book/:id", func(c *gin.Context) {
		mockBooksHandler.Update(c)
	})
	router.ServeHTTP(httptest.NewRecorder(), req)
}

func TestDeleteHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBooksHandler := mocks.NewMockBooksHandlerInterface(ctrl)

	mockBooksHandler.EXPECT().Delete(gomock.Any()).Return()

	req, _ := http.NewRequest("DELETE", "/book/123", nil)

	router := gin.Default()
	router.DELETE("/book/:id", func(c *gin.Context) {
		mockBooksHandler.Delete(c)
	})
	router.ServeHTTP(httptest.NewRecorder(), req)
}
