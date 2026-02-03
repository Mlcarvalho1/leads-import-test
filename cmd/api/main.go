package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"your-app/database"
	"your-app/routes"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const banner = `
       ,_---~~~~~----._
_,,_,*^____      _____''*g*\"*,          Welcome to your app!
/ __/ /'     ^.  /      \ ^@q   f     /
[  @f | @))    |  | @))   l  0 _/    /
\'/   \~____ / __ \_____/    \      /
 |           _l__l_           I     /
 }          [______]           I   /
 ]            | | |            |
 ]             ~ ~             |
 |                            |
  |                           |`

func main() {
	app := fiber.New()
	database.ConnectDb()

	routes.SetupRoutes(app)

	log.Println(banner)

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
