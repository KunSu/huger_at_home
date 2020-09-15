package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World")

	r := gin.Default()
	//TODO: URL design and catagorize them into different /
	r.GET("/", homePage)
	r.GET("api/signup", querySignUp)
	r.GET("api/login", queryLogin)
	r.GET("api/reset", queryResetPassword)
	r.GET("api/companySignup", queryCompanySignUp)
	r.POST("api/postTest", postTest)
	r.Run()
}
