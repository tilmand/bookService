package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func configureRouter(api *api) *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())

	public := router.Group("api/v1")

	public.POST("/signIn", api.Auth().SignIn)
	public.POST("/refresh", api.Auth().Refresh)
	public.POST("/signUp", api.Auth().SignUp)
	public.POST("/recover", api.Auth().Recover)
	public.POST("/setNewPassword/:token", api.Auth().SetNewPassword)

	public.GET("/books", api.Books().GetAll)
	public.GET("/book/:id", api.Books().Find)
	public.POST("/book", api.Books().Add)
	public.PUT("/book/:id", api.Books().Update)
	public.DELETE("/book/:id", api.Books().Delete)

	router.NoRoute(func(c *gin.Context) {
		log.Println("route not found")
		c.JSON(http.StatusNotFound, errors.New("record not found"))
	})

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding,"+
			"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)

			return
		}

		c.Next()
	}
}
