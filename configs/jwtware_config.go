package configs

import (
	"MyGram/configs/rsaconfig"
	"MyGram/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

var JwtWareConfig jwtware.Config

func ConfigJwtware() {
	JwtWareConfig = jwtware.Config{
		SigningMethod:  "RS512",
		SigningKey:     rsaconfig.PrivateKey.Public(),
		ErrorHandler:   JwtWareErrorHandler,
		ContextKey:     "user",
		TokenLookup:    "header:Authorization,cookie:token",
		SuccessHandler: jwtWareSuccessHandler,
	}
}

func JwtWareErrorHandler(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return errors.New("please login first")
	}
	return c.Status(fiber.StatusUnauthorized).JSON(utils.ResponseFail("Invalid or expired JSON Web Token"))
}

func jwtWareSuccessHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	//userId != claims["user_id"].(float64)
	c.Locals("user", claims)

	return c.Next()

}
