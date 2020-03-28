package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	
	// For Gin and JWT
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"	
	
	// For Google social account auth
	"google.golang.org/api/oauth2/v2"
	
	// Import local packages
	"github.com/deepenpatel19/go-gin-api/api"

	"github.com/deepenpatel19/go-gin-api/models"
	
)

type tokenSignIn struct {
	IdToken string `form:"idtoken" json:"idtoken" binding:"required"`
	SocialSource string `form:"socialsource" json:"socialsource" binding:"required"`
}

var identityKey = "id"

func verifyIdToken(idToken string) (*oauth2.Tokeninfo, error) {
	var httpClient = &http.Client{}
    oauth2Service, err := oauth2.New(httpClient)
    tokenInfoCall := oauth2Service.Tokeninfo()
    tokenInfoCall.IdToken(idToken)
    tokenInfo, err := tokenInfoCall.Do()
    if err != nil {
        return nil, err
    }
    return tokenInfo, nil
}

func checkGoogleID(idToken string) (email string, status bool) {
	token, err := verifyIdToken(idToken)
	if err != nil {
		fmt.Println("err", err)
	}
	return token.Email, token.VerifiedEmail
}


func main() {
	models.InitializeDB()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Users",
		})
	})

	
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*api.GoogleVerified); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &api.GoogleVerified{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var tokenVals tokenSignIn
			c.ShouldBind(&tokenVals)
			
			var idToken = tokenVals.IdToken
			var socialsource = tokenVals.SocialSource
			email, verificationStatus := checkGoogleID(idToken)		
			googleVerified := api.GoogleVerified{}
			
			googleVerified.UserName = email

			if err := c.ShouldBind(&googleVerified); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			if verificationStatus {
				// TODO: Throw error on db error for user create/exists.
				models.UserCreate(email)
				return &api.GoogleVerified{
					UserName: email,
				}, nil
			} else {
				return nil, jwt.ErrFailedAuthentication
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			v1, ok1 := data.(*api.GoogleVerified)
			userExists := models.UserExist(v1.UserName)
			if userExists && ok1 {
				return true
			} else {
				return false
			}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/tokensignin/", authMiddleware.LoginHandler)
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	auth := r.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/info/", api.UserInfo)
	}

	if err := http.ListenAndServe(":"+"8000", r); err != nil {
		log.Fatal(err)
	}	
}

