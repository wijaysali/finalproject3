package configs

import (
	"MyGram/middlewares/recover"
	"io"
	"log"
	"os"
)

//ini digunakan untuk logger dari midleware recover
func RecoverConfig() recover.Config {
	mw := io.MultiWriter(os.Stdout)
	logger := log.New(mw, "\r\n", log.LstdFlags)
	config := recover.Config{
		EnableStackTrace: true,
		Logger:           logger,
	}
	return config
}
