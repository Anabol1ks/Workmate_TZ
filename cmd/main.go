package main

import "workmate_tz/internal/logger"

func main() {
	if err := logger.Init(); err != nil {
		panic(err)
	}

	defer logger.Sync()

	log := logger.L()

	log.Info("Application started")
}
