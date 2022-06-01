package usercontroller

import (
	"MyGram/app"
	"MyGram/configs/rsaconfig"
	"MyGram/models"
	"MyGram/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"github.com/golang-jwt/jwt/v4"
)

type InputLogin struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func Register(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	app.Validate(c, &user)
	user.Password = utils.EncryptToBcrypt(&user.Password)
	err := app.Db.Create(&user).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	result := fiber.Map{
		"id":       user.ID,
		"age":      user.Age,
		"email":    user.Email,
		"username": user.Username,
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func Login(c *fiber.Ctx) error {
	user := models.User{}
	input := InputLogin{}
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	app.Validate(c, &input)

	err := app.Db.First(&user, "email=?", input.Email).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	if !utils.IsHashBcryptMatch(&input.Password, &user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "email and password not match", "success": false})
	}

	result := map[string]interface{}{"user_id": user.ID}
	if viper.GetInt("AUTH_JWT_EXPIRATION") != -99 {
		result["exp"] = time.Now().Add(time.Minute * time.Duration(viper.GetInt("AUTH_JWT_EXPIRATION"))).Unix()
	}
	claims := jwt.MapClaims(result)

	// Create token
	rawToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	// Generate encoded token and send it as response.
	signedToken, err := rawToken.SignedString(rsaconfig.PrivateKey)
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"token": signedToken})
}

func EditUser(c *fiber.Ctx) error {
	inputEditUser := struct {
		Username string `validate:"required"`
		Email    string `validate:"required,email"`
	}{}

	userIdParam, err := c.ParamsInt("userId")
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := int(userClaim["user_id"].(float64))

	if err != nil {
		return err
	}
	if err := c.BodyParser(&inputEditUser); err != nil {
		return err
	}

	app.Validate(c, &inputEditUser)

	if userIdParam != userIdFromClaim {
		utils.ResponsePanic(fiber.StatusUnauthorized, "You are not permit to edit this user")
	}

	user := models.User{}

	user.ID = uint(userIdParam)
	err = app.Db.First(&user).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	user.Username = inputEditUser.Username
	user.Email = inputEditUser.Email
	err = app.Db.Save(&user).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"id": user.ID,
		"email":     user.Email,
		"username":  user.Username,
		"age":       user.Age,
		"udated_at": user.UpdatedAt.Format("2006-01-02")})

}

func DeleteUser(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	err := app.Db.First(&models.User{}, userIdFromClaim).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	err = app.Db.Delete(&models.User{}, userIdFromClaim).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{
		"message": "Your account has been successfully deleted",
	})
}
