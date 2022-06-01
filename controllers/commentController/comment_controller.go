package commentcontroller

import (
	"MyGram/app"
	"MyGram/models"
	"MyGram/utils"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
)

func AddComment(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	comment := models.Comment{}
	if err := c.BodyParser(&comment); err != nil {
		return err
	}

	comment.UserID = userIdFromClaim
	app.Validate(c, &comment)

	//validate photo id harus ada
	err := app.Db.First(&models.Photo{}, "id = ?", comment.PhotoID).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, "photo not found")
	}

	err = app.Db.Create(&comment).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": comment.ID,
		"created_at": comment.CreatedAt.Format("2006-01-02"),
		"message":    comment.Message,
		"photo_id":   comment.PhotoID,
		"user_id":    comment.UserID})
}

func EditComment(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	commentIdParam, err := c.ParamsInt("commentId")
	if err != nil {
		return err
	}

	input := map[string]interface{}{}
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	if val, ok := input["message"]; !(ok && val != "") {
		utils.ResponsePanic(fiber.StatusBadRequest, "message is required")
	}

	commentTobeEdit := models.Comment{}
	err = app.Db.First(&commentTobeEdit, "id=?", commentIdParam).Error

	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	if commentTobeEdit.UserID != userIdFromClaim {
		utils.ResponsePanic(fiber.StatusUnauthorized, "You are not permit to edit this comment")
	}

	commentTobeEdit.Message = input["message"].(string)
	err = app.Db.Save(&commentTobeEdit).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	photo := models.Photo{}

	err = app.Db.First(&photo, "id = ?", commentTobeEdit.PhotoID).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"id": commentTobeEdit.ID,
		"updated_at": commentTobeEdit.UpdatedAt.Format("2006-01-02"),
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoUrl,
		"user_id":    commentTobeEdit.UserID})
}

func DeleteComment(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	commentIdParam, err := c.ParamsInt("commentId")
	if err != nil {
		return err
	}

	commentTobeDelete := models.Comment{}
	err = app.Db.First(&commentTobeDelete, "id=?", commentIdParam).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	if commentTobeDelete.UserID != userIdFromClaim {
		utils.ResponsePanic(fiber.StatusUnauthorized, "You are not permit to delete this comment")
	}

	err = app.Db.Delete(&models.Comment{}, commentTobeDelete).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"message": "Your comment has been successfully deleted"})
}

func GetComments(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	var query strings.Builder
	resultQuery := []map[string]interface{}{}

	query.WriteString("SELECT A.id, A.message, A.photo_id, A.user_id, DATE_FORMAT(A.created_at, '%Y-%m-%d') AS created_at, DATE_FORMAT(A.updated_at, '%Y-%m-%d') AS updated_at, ")
	query.WriteString("       B.email, B.username, ")
	query.WriteString("       C.title, C.caption, C.photo_url, C.user_id AS photo_user_id ")
	query.WriteString(" FROM comments A ")
	query.WriteString(" INNER JOIN users B ON A.user_id = B.id ")
	query.WriteString(" INNER JOIN photos C ON A.photo_id = C.id ")
	query.WriteString("WHERE A.user_id = ? ")

	err := app.Db.Raw(query.String(), userIdFromClaim).Scan(&resultQuery).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	for _, row := range resultQuery {
		row["User"] = map[string]interface{}{
			"id":       row["user_id"],
			"email":    row["email"],
			"username": row["username"]}

		row["Photo"] = map[string]interface{}{
			"id":        row["photo_id"],
			"title":     row["title"],
			"caption":   row["caption"],
			"photo_url": row["photo_url"],
			"user_id":   row["photo_user_id"]}
		delete(row, "email")
		delete(row, "username")
		delete(row, "title")
		delete(row, "caption")
		delete(row, "photo_url")
		delete(row, "photo_user_id")
	}

	return c.JSON(resultQuery)
}
