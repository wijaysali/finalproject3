package recover

import (
	"MyGram/utils"
	"fmt"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

func defaultStackTraceHandler(c *fiber.Ctx, cfg Config, e interface{}) {
	buf := make([]byte, defaultStackTraceBufLen)
	buf = buf[:runtime.Stack(buf, false)]
	cfg.Logger.Printf(fmt.Sprintf("error : %v\n%s\n", e, buf))
	//logger.Printf(fmt.Sprintf("panic: %v\n%s\n", e, buf))
	//	_, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", e, buf))
}

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Catch panics
		defer func() {
			if r := recover(); r != nil {
				fiberError, isOk := r.(fiber.Error)
				if isOk && cfg.EnableStackTrace && (fiberError.Code == fiber.StatusInternalServerError) {
					cfg.StackTraceHandler(c, cfg, r)
				}
				if !isOk && cfg.EnableStackTrace {
					cfg.StackTraceHandler(c, cfg, r)
				}
				var ok bool
				//err,_ = r.(*fiber.Error)
				if isOk {
					err = c.Status(fiberError.Code).JSON(
						utils.ResponseFail(fiberError.Message))
				} else {

					if err, ok = r.(error); !ok {
						// Set error that will call the global error handler
						err = fmt.Errorf("%v", r)
					}
				}
			}
		}()

		// Return err if exist, else move to next handler
		return c.Next()
	}
}
