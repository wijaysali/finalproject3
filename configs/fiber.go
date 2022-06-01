package configs

import (
	"time"

	"MyGram/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func FiberConfigFromEnv() fiber.Config {

	fiberConfig := fiber.Config{
		AppName:         viper.GetString("APP_NAME"),
		CaseSensitive:   viper.GetBool("CASE_SENSITIVE_ROUTE"),
		Immutable:       viper.GetBool("IMMUTABLE_VALUE_CONTEXT"),
		BodyLimit:       viper.GetInt("MAX_BODY_LIMIT_SIZE") * 1024 * 1024,
		Concurrency:     viper.GetInt("MAX_CONCURRENCY"),
		ReadBufferSize:  viper.GetInt("MAX_READ_BUFFER_SIZE") * 1024,
		WriteBufferSize: viper.GetInt("MAX_WRITE_BUFFER_SIZE") * 1024,
		ReadTimeout:     time.Second * 1,
		WriteTimeout:    time.Second * 1,
		IdleTimeout:     time.Second * 1,

		EnablePrintRoutes: true, ErrorHandler: errorHandler,
	}
	if viper.IsSet("MAX_READ_TIMEOUT") {
		fiberConfig.ReadTimeout = time.Duration(viper.GetInt64("MAX_READ_TIMEOUT")) * time.Second
	}

	if viper.IsSet("MAX_RESPONSE_TIMEOUT") {
		fiberConfig.WriteTimeout = time.Duration(viper.GetInt64("MAX_RESPONSE_TIMEOUT")) * time.Second
	}
	return fiberConfig
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	if err != nil {
		// In case the SendFile fails
		return ctx.Status(code).JSON(utils.ResponseFail(err.Error()))
	}

	// Return from handler
	return nil
}
