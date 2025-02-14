package main

import (
	"fmt"
	"iwogo/Models"
	"iwogo/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//---------------------conect db postgre------------------------------
	//connStr := os.Getenv("DATABASE_URL")
	// db, err := gorm.Open(postgres.Open(Config.DbURL(Config.BuildDBConfig())), &gorm.Config{})

	//----------------------conect neon postgre------------------------------------
	connStr := os.Getenv("DATABASE_URL") + "&options=endpoint%3Dep-hidden-lab-a4ykgi4e"
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	//----------------rabitmq-----------------
	// conn, errs := amqp.Dial(os.Getenv("RABBITMQ_HOST"))
	// errorWrapper(errs, "Failed to connect rabbitmq")
	// defer conn.Close()

	// ch, errs := conn.Channel()
	// errorWrapper(errs, "Failed to open a channel")
	// defer ch.Close()

	//---------------------------------------------

	if err != nil {
		fmt.Println("Status DB:", err)
	}

	//----------------Migration DB-----------------
	db.AutoMigrate(
		&Models.User{},
	)

	//Routes.SubscribeMessage(ch, "go-queue_order")

	router := routes.SetupRouter(db)
	router.Run(":" + os.Getenv("APP_PORT"))

}
