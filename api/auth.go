package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type GoogleVerified struct {
	UserName string `json:"username"`
}


func UserInfo(c *gin.Context) {
	fmt.Println("User Information.")
	user, _ := c.Get("id")
	fmt.Println("user", user.(*GoogleVerified).UserName)
	c.JSON(200, gin.H{
		"userID":   "Test",
		"userName": user.(*GoogleVerified).UserName,
		"text":     "Hello World.",
	})
}
