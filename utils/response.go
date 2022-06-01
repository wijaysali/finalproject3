package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ResponseData struct {
	Success bool
	Data    interface{}
}

func ResponseSuccess(data interface{}) fiber.Map {
	return fiber.Map{
		"success": true,
		"data":    data,
	}
}

func ResponseFail(data interface{}) fiber.Map {
	return fiber.Map{
		"success": false,
		"message": data,
	}
}

func ResponsePanic(statusCode int, errMsg ...string) {
	strBuilder := strings.Builder{}

	for idx, str := range errMsg {
		strBuilder.WriteString(str)
		if len(errMsg) != idx+1 {
			strBuilder.WriteString(", ")
		}
	}
	panic(fiber.Error{Code: statusCode, Message: strBuilder.String()})
}
