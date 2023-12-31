package service

import (
	"log"
	"mini-account-service/db"
	"mini-account-service/server"
	"mini-account-service/server/session"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

func parseParams() (
	string,
	string,
	string,
	db.Orm,
	sessions.Store,
	string,
) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	orm, err := db.NewOrm()
	if err != nil {
		log.Fatal(err)
	}

	name := os.Getenv("SERVICE_NAME")
	if len(name) == 0 {
		log.Fatal("env param SERVICE_NAME is empty")
	}
	ver := os.Getenv("VERSION")
	if len(ver) == 0 {
		log.Fatal("env param VERSION is empty")
	}
	at := os.Getenv("SERVER_LISTEN_AT")
	if len(at) == 0 {
		log.Fatal("env param SERVER_LISTEN_AT is empty")
	}
	secret := os.Getenv("SESSION_SECRET")
	if len(secret) == 0 {
		log.Fatal("env param SESSION_SECRET is empty")
	}
	gameData := os.Getenv("GAMEDATA_SERVICE_URL")
	if len(gameData) == 0 {
		log.Fatal("env param GAMEDATA_SERVICE_URL is empty")
	}

	store := session.NewSessionStroe(secret)

	return name, ver, at, orm, store, gameData
}

func Run() {
	server.Run(parseParams())
}
