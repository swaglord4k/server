package app

import (
	"fmt"
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
	var configs m.Conf
	configs.GetConf(m.MysqlConfig)
	db, err := db.ConnectToMySQLDb(configs)
	if err != nil {
		fmt.Println("Couldn't connect to DB " + configs.Dbname)
		panic(err)
	}

	userController := c.NewController[m.User](db, router, "crud")
	c.NewUserRoutes(userController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
