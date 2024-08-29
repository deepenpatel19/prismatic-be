package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/deepenpatel19/prismatic-be/schemas"
)

var Config schemas.ProjectConfiguration

func ReadEnvFile() {
	environment := os.Getenv("ENVIRONMENT")
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Error while getting current working directory :: ", err)
		os.Exit(1)
	}
	fmt.Println(currentWorkingDirectory)
	fileName := currentWorkingDirectory + "/.env"

	if environment == "production" {
		err := os.WriteFile(fileName, []byte(os.Getenv("DATA")), 0755)
		if err != nil {
			fmt.Println("unable to write file: %w", err)
			os.Exit(1)
		}
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error while reading env file :: ", err)
		os.Exit(1)
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		fmt.Println("Error while converting env file data into struct :: ", err)
		os.Exit(1)
	}

	Config.DBString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		Config.DBConfig.DBUser,
		Config.DBConfig.DBPassword,
		Config.DBConfig.DBHost,
		Config.DBConfig.DBPort,
		Config.DBConfig.DBName,
		Config.DBConfig.DBSSLMode,
	)

	Config.DBQueryTimeout = 60

}
