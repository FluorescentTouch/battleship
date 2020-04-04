package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"

	_ "my/battleship/docs"

	"my/battleship/battlefield"
)

// @title Swagger Example API
// @version 2.0
// @description This is a battleships game server.

// @contact.name Shamil Garatuev
// @contact.email garatuev@gmail.com

// @host localhost:8080
// @BasePath /

func main() {
	log := logrus.New()

	bs := battlefield.NewService(log)
	be := battlefield.NewEndpoints(log, bs)
	bh := battlefield.NewHandlers(log, be)

	router := mux.NewRouter()

	// swagger
	router.PathPrefix("/swagger/").
		Handler(httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json"))).
		Methods("GET")

	// API
	router.HandleFunc("/create-matrix", bh.CreateBattleField).Methods("POST")

	log.Infof("listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
