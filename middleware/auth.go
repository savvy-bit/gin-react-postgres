package middleware

import (
	"log"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/savvy-bit/gin-react-postgres/config"
	"github.com/savvy-bit/gin-react-postgres/model"
)

var authMiddleware *ginjwt.GinJWTMiddleware
var identityKey = "email"

// Login struct
type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
}

// Auth middleware
func Auth() *ginjwt.GinJWTMiddleware {
	return authMiddleware
}

func init() {
	var err error
	authMiddleware, err = ginjwt.New(&ginjwt.GinJWTMiddleware{
		Realm:       "gin-skeleton",
		Key:         []byte(config.Global.Server.SecurityKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		SendCookie:  true,
		PayloadFunc: func(data any) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
					"name":      v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) any {
			claims := ginjwt.ExtractClaims(c)
			return &model.User{
				Email: claims[identityKey].(string),
				Name:  claims["name"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (any, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", ginjwt.ErrMissingLoginValues
			}
			email := loginVals.Email
			password := loginVals.Password

			return model.LoginByEmailAndPassword(email, password)
		},
		Authorizator: func(data any, c *gin.Context) bool {
			if _, ok := data.(*model.User); ok {
				return true
			}

			return false
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
}
