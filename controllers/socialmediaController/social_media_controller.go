package socialmediacontroller

import (
	"MyGram/app"
	"MyGram/models"
	"MyGram/utils"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
)

func AddSocialMedia(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	socialMedia := models.SocialMedia{}
	if err := c.BodyParser(&socialMedia); err != nil {
		return err
	}

	socialMedia.UserID = userIdFromClaim
	app.Validate(c, &socialMedia)

	err := app.Db.Create(&socialMedia).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": socialMedia.ID,
		"created_at":       socialMedia.CreatedAt.Format("2006-01-02"),
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaUrl,
		"user_id":          socialMedia.UserID})
}

func EditSocialMedia(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	socialMediaIdParam, err := c.ParamsInt("socialMediaId")
	if err != nil {
		return err
	}

	inputSocialMedia := models.SocialMedia{}
	if err := c.BodyParser(&inputSocialMedia); err != nil {
		return err
	}

	app.Validate(c, &inputSocialMedia)

	socialMediaTobeEdit := models.SocialMedia{}
	err = app.Db.First(&socialMediaTobeEdit, "id=?", socialMediaIdParam).Error

	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	if socialMediaTobeEdit.UserID != userIdFromClaim {
		utils.ResponsePanic(fiber.StatusUnauthorized, "You are not permit to edit this comment")
	}

	socialMediaTobeEdit.Name = inputSocialMedia.Name
	socialMediaTobeEdit.SocialMediaUrl = inputSocialMedia.SocialMediaUrl
	err = app.Db.Save(&socialMediaTobeEdit).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"id": socialMediaTobeEdit.ID,
		"updated_at":       socialMediaTobeEdit.UpdatedAt.Format("2006-01-02"),
		"name":             socialMediaTobeEdit.Name,
		"social_media_url": socialMediaTobeEdit.SocialMediaUrl,
		"user_id":          socialMediaTobeEdit.UserID})
}

func DeleteSocialMedia(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	socialMediaIdParam, err := c.ParamsInt("socialMediaId")
	if err != nil {
		return err
	}

	socialMediaTobeDelete := models.SocialMedia{}
	err = app.Db.First(&socialMediaTobeDelete, "id=?", socialMediaIdParam).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	if socialMediaTobeDelete.UserID != userIdFromClaim {
		utils.ResponsePanic(fiber.StatusUnauthorized, "You are not permit to delete this comment")
	}

	err = app.Db.Delete(&models.SocialMedia{}, socialMediaTobeDelete).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"message": "Your social media has been successfully deleted"})
}

func GetSocialMedia(c *fiber.Ctx) error {
	userClaim := c.Locals("user").(jwt.MapClaims)
	userIdFromClaim := uint(userClaim["user_id"].(float64))

	var query strings.Builder
	resultQuery := []map[string]interface{}{}

	query.WriteString("SELECT A.id, A.name, A.social_media_url, A.user_id as userId, DATE_FORMAT(A.created_at, '%Y-%m-%d') AS createdAt, DATE_FORMAT(A.updated_at, '%Y-%m-%d') AS updatedAt, ")
	query.WriteString("       B.username, coalesce(C.photo_url,'') as photo_url ")
	//query.WriteString("       C.title, C.caption, C.photo_url, C.user_id AS photo_user_id ")
	query.WriteString(" FROM social_media A ")
	query.WriteString(" INNER JOIN users B ON A.user_id = B.id ")
	query.WriteString(" LEFT JOIN photos C ON A.user_Id = C.user_id ")
	query.WriteString("WHERE A.user_id = ? ")

	err := app.Db.Raw(query.String(), userIdFromClaim).Scan(&resultQuery).Error
	if err != nil {
		utils.ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}

	for _, row := range resultQuery {
		row["User"] = map[string]interface{}{
			"id":                row["user_id"],
			"username":          row["username"],
			"profile_photo_url": row["photo_url"]}
		delete(row, "email")
		delete(row, "username")
		delete(row, "title")
		delete(row, "caption")
		delete(row, "photo_url")
		delete(row, "photo_user_id")
	}

	return c.JSON(resultQuery)
}
