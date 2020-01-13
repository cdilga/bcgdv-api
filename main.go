package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/jinzhu/gorm"
	//"github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

type Customer struct {
	Id      int    `gorm:"column:id; AUTO_INCREMENT" json:Id`
	Name    string `gorm:"column:name; not null" json:"Name"`
	Address string `gorm:"column:address" json:"Address"`
	Email   string `gorm:"column:email; not null" json:"Email"`
	Phone   string `gorm:"column:phone" json:"Phone"`
}

type Customers []Customer

func CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	db := InitDb()
	defer db.Close()

	var user Customers
	r.Bind(&user)

	if user.Firstname != "" && user.Lastname != "" {
		// INSERT INTO "users" (name) VALUES (user.Name);
		db.Create(&user)
		// Display error
		c.JSON(201, w.WriteHeader{"success": user})
	} else {
		// Display error
		c.JSON(422, w.WriteHeader{"error": "Fields are empty"})
	}
}

func GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	var users []Users
	// SELECT * FROM users
	db.Find(&users)

	// Display JSON result
	r.JSON(200, users)
}

func InitDb() *gorm.DB {
	// Openning file
	db, err := gorm.Open("sqlite3", ":memory:")
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}
	// Creating the table
	if !db.HasTable(&Customers{}) {
		db.CreateTable(&Customers{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Customers{})
	}

	return db
}


func NotifyHandler(w http.ResponseWriter, r *http.Request) {
	customers := Customers{
		Customer{Name: "John Doe", Address: "123 Seasame St", Email: "email@example.com", Phone: "04216117783"},
	}

	fmt.Println("Endpoint Hit: Notify")
	json.NewEncoder(w).Encode(customers)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func main() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*5, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/customer", CreateCustomerHandler).Methods("POST")
	r.HandleFunc("/customer", GetCustomerHandler).Methods("GET")
	r.HandleFunc("/notify", NotifyHandler).Methods("POST")
	http.Handle("/", r)
	// Add your routes as needed

	srv := &http.Server{
		Addr: "0.0.0.0:8081",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
