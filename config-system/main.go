package main

import (
	"fmt"
	"learning/config-system/config"
)

func main() {
	cfg := config.Load()

	fmt.Println("Environment", cfg.AppEnv)
	fmt.Println("Port:", cfg.Port)
	fmt.Println("DB URL", cfg.DBUrl)

}
