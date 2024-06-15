package app

import (
	"log"
	"net/http"

	c "de.amplifonx/app/controller"
	"de.amplifonx/app/db"
	m "de.amplifonx/app/model"
	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

func NewApp() {
	router := chi.NewRouter()
	db, err := db.ConnectToDB(m.MysqlConfig)
	if err != nil {
		panic(err)
	}

	userController := c.NewController[m.User](db, router, "user")
	c.NewUserRoutes(userController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
