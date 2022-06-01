package photocontroller

import (
	"MyGram/app"
	"MyGram/models"
	"MyGram/utils"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
)

func AddPhoto(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	photo := models.Photo{}
	if err := c.BodyParser(&photo); err != nil {
		return err
	}

	photo.UserID = userIdFromClaim
	app.Validate(c, &photo)

	err := app.Db.Create(&photo).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": photo.ID,
		"created_at": photo.CreatedAt.Format("2006-01-02"),
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoUrl,
		"user_id":    photo.UserID})
}

func EditPhoto(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	photosIdParam, err := c.ParamsInt("photoId")
	if err != nil {
		return err
	}

	photoInput := models.Photo{}
	if err := c.BodyParser(&photoInput); err != nil {
		return err
	}

	app.Validate(c, &photoInput)

	photoTobeEdit := models.Photo{}
	err = app.Db.First(&photoTobeEdit, "id=?", photosIdParam).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	if photoTobeEdit.UserID != userIdFromClaim {
		utils.ResponsePanic(fiber.StatusUnauthorized, "You are not permit to edit this photo")
	}

	photoTobeEdit.Caption = photoInput.Caption
	photoTobeEdit.Title = photoInput.Title
	photoTobeEdit.PhotoUrl = photoInput.PhotoUrl
	err = app.Db.Save(&photoTobeEdit).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"id": photoTobeEdit.ID,
		"updated_at": photoTobeEdit.UpdatedAt.Format("2006-01-02"),
		"title":      photoTobeEdit.Title,
		"caption":    photoTobeEdit.Caption,
		"photo_url":  photoTobeEdit.PhotoUrl,
		"user_id":    photoTobeEdit.UserID})
}

func DeletePhoto(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	photosIdParam, err := c.ParamsInt("photoId")
	if err != nil {
		return err
	}

	photoTobeDelete := models.Photo{}
	err = app.Db.First(&photoTobeDelete, "id=?", photosIdParam).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	if photoTobeDelete.UserID != userIdFromClaim {
		utils.ResponsePanic(fiber.StatusUnauthorized, "You are not permit to delete this photo")
	}

	err = app.Db.Delete(&models.Photo{}, photosIdParam).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"message": "Your photo has been successfully deleted"})
}

func GetPhotos(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	var query strings.Builder
	resultQuery := []map[string]interface{}{}

	query.WriteString("SELECT A.id, A.title, A.caption, A.photo_url, A.user_id, DATE_FORMAT(A.created_at, '%Y-%m-%d') AS created_at, DATE_FORMAT(A.updated_at, '%Y-%m-%d') AS updated_at, ")
	query.WriteString("       B.email, B.username ")
	query.WriteString(" FROM photos A INNER JOIN users B ON A.user_id = B.id ")
	query.WriteString("WHERE A.user_id = ? ")

	err := app.Db.Raw(query.String(), userIdFromClaim).Scan(&resultQuery).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	for _, row := range resultQuery {
		row["User"] = map[string]interface{}{"email": row["email"], "username": row["username"]}
		delete(row, "email")
		delete(row, "username")
	}

	return c.JSON(resultQuery)
}
