package http

import (
	"bookService/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const DecimalBase = 10
const BitSize64 = 64

type BooksHandlerInterface interface {
	GetAll(c *gin.Context)
	Add(c *gin.Context)
	Find(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type BooksHandler struct {
	api *api
}

func NewBooksHandler(a *api) *BooksHandler {
	return &BooksHandler{
		api: a,
	}
}
func (h *BooksHandler) GetAll(c *gin.Context) {
	results, err := h.api.mongo.BooksRepository.GetAll()
	if err != nil {
		log.Println("GetAll GetAll err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}
	answer := map[string]interface{}{
		"items": results,
	}

	c.JSON(http.StatusOK, answer)
}

func (h *BooksHandler) Add(c *gin.Context) {
	token := h.api.auth.ExtractToken(c.Request)
	claims, err := h.api.auth.Validate(token)
	if err != nil {
		log.Println("Add Validate err: ", err)
		c.JSON(http.StatusUnauthorized, model.ErrUnauthorized)

		return
	}

	var item model.Book
	err = c.ShouldBindJSON(&item)
	if err != nil {
		log.Println("Add ShouldBindJSON err: ", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	err = h.api.mongo.BooksRepository.Insert(item, claims.BaseClaims.ID)
	if err != nil {
		log.Println("Add Insert err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book created successfully"})
}

func (h *BooksHandler) Find(c *gin.Context) {
	idStr := c.Param("id")
	ID, err := strconv.ParseUint(idStr, DecimalBase, BitSize64)
	if err != nil {
		log.Println("Find ParseUint err: ", err)
		c.JSON(http.StatusBadRequest, model.ErrInternalServerError)

		return
	}

	item, err := h.api.mongo.BooksRepository.Find(ID)
	if err != nil {
		log.Println("Find Find err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	answer := map[string]interface{}{
		"item": item,
	}

	c.JSON(http.StatusOK, answer)
}

func (h *BooksHandler) Update(c *gin.Context) {
	token := h.api.auth.ExtractToken(c.Request)
	claims, err := h.api.auth.Validate(token)
	if err != nil {
		log.Println("Update Validate err: ", err)
		c.JSON(http.StatusUnauthorized, model.ErrUnauthorized)

		return
	}

	var item model.Book
	err = c.ShouldBindJSON(&item)
	if err != nil {
		log.Println("Update ShouldBindJSON err: ", err)
		c.JSON(http.StatusBadRequest, model.ErrInternalServerError)

		return
	}

	idStr := c.Param("id")

	ID, err := strconv.ParseUint(idStr, DecimalBase, BitSize64)
	if err != nil {
		log.Println("Update ParseUint err: ", err)
		c.JSON(http.StatusBadRequest, model.ErrInternalServerError)

		return
	}

	existingBook, err := h.api.mongo.BooksRepository.Find(ID)
	if err != nil {
		log.Println("Update Find err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	if existingBook.AuthorID != claims.BaseClaims.ID {
		log.Println("Update AuthorID err: ", err)
		c.JSON(http.StatusForbidden, model.ErrForbidden)

		return
	}

	item.ID = ID
	item.AuthorID = claims.BaseClaims.ID
	err = h.api.mongo.BooksRepository.Update(item)
	if err != nil {
		log.Println("Update Update err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book updated successfully"})
}

func (h *BooksHandler) Delete(c *gin.Context) {
	token := h.api.auth.ExtractToken(c.Request)
	claims, err := h.api.auth.Validate(token)
	if err != nil {
		log.Println("Delete Validate err: ", err)
		c.JSON(http.StatusUnauthorized, model.ErrUnauthorized)

		return
	}

	idStr := c.Param("id")

	ID, err := strconv.ParseUint(idStr, DecimalBase, BitSize64)
	if err != nil {
		log.Println("Delete ParseUint err: ", err)
		c.JSON(http.StatusBadRequest, model.ErrInternalServerError)

		return
	}

	existingBook, err := h.api.mongo.BooksRepository.Find(ID)
	if err != nil {
		log.Println("Delete Find err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	if existingBook.AuthorID != claims.BaseClaims.ID {
		log.Println("Delete AuthorID err: ", err)
		c.JSON(http.StatusForbidden, model.ErrForbidden)

		return
	}

	err = h.api.mongo.BooksRepository.Delete(ID)
	if err != nil {
		log.Println("Delete Delete err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}
