package routers

import (
	"MyGram/configs"
	commentcontroller "MyGram/controllers/commentController"
	photocontroller "MyGram/controllers/photoController"
	socialmediacontroller "MyGram/controllers/socialmediaController"
	usercontroller "MyGram/controllers/userController"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func RegisterRouters(fiberApp *fiber.App) {
	//grouping api
	//api := fiberApp.Group("/api", onlyjsonrequest.New())
	authenticationMdware := jwtware.New(configs.JwtWareConfig)

	userRoute := fiberApp.Group("/users")

	userRoute.Post("/register", usercontroller.Register)
	userRoute.Post("login", usercontroller.Login)
	userRoute.Put("/:userId", authenticationMdware, usercontroller.EditUser)
	userRoute.Delete("/", authenticationMdware, usercontroller.DeleteUser)

	photoRoute := fiberApp.Group("/photos")
	photoRoute.Use(authenticationMdware)
	photoRoute.Post("/", photocontroller.AddPhoto)
	photoRoute.Put("/:photoId", photocontroller.EditPhoto)
	photoRoute.Delete("/:photoId", photocontroller.DeletePhoto)
	photoRoute.Get("/", photocontroller.GetPhotos)

	commentRoute := fiberApp.Group("/comments")
	commentRoute.Use(authenticationMdware)
	commentRoute.Post("/", commentcontroller.AddComment)
	commentRoute.Put("/:commentId", commentcontroller.EditComment)
	commentRoute.Delete("/:commentId", commentcontroller.DeleteComment)
	commentRoute.Get("/", commentcontroller.GetComments)

	socialMediaRoute := fiberApp.Group("/socialmedias")
	socialMediaRoute.Use(authenticationMdware)
	socialMediaRoute.Post("/", socialmediacontroller.AddSocialMedia)
	socialMediaRoute.Put("/:commentId", socialmediacontroller.EditSocialMedia)
	socialMediaRoute.Delete("/:commentId", socialmediacontroller.DeleteSocialMedia)
	socialMediaRoute.Get("/", socialmediacontroller.GetSocialMedia)
	// userroute.MakeRoute(api)
	// authroute.MakeRoute(api)
	// roleroute.MakeRoute(api)
	// api.Get("/setDataRedis", func(c *fiber.Ctx) error {
	// 	gows.RedisDb.Set("coba", []byte("nilai coba"), time.Minute)
	// 	return c.JSON("sukses")
	// })

	// api.Get("/getDataRedis", func(c *fiber.Ctx) error {
	// 	gows.RedisDb.Set("coba", []byte("nilai coba"), time.Minute)
	// 	data, _ := gows.RedisDb.Get("coba")
	// 	return c.JSON(string(data))
	// })

}
