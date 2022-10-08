package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:name`
	Content string `json:content`
}

// conectionBD create the connection to the database
// return a pointer with the connection to our database
func conectionBD() *sql.DB {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	driver := os.Getenv("DB_DRIVER")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	conection, err := sql.Open(driver, user+":"+password+"@tcp("+host+":"+port+")/"+name)
	if err != nil {
		panic(err.Error())
	}
	return conection
}

// getTasks this method get all saved task
func getTasks(w http.ResponseWriter, r *http.Request) {

	conection := conectionBD()
	getRegisters, err := conection.Query("SELECT * from task")

	if err != nil {
		panic(err.Error())
	}

	taskAux := task{}
	arrayTasks := []task{}
	for getRegisters.Next() {
		var id int
		var name, content string
		err = getRegisters.Scan(&id, &name, &content)
		if err != nil {
			panic(err.Error())
		}
		taskAux.ID = id
		taskAux.Name = name
		taskAux.Content = content
		arrayTasks = append(arrayTasks, taskAux)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(arrayTasks)
}

// createTask this method create a new task
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}

	json.Unmarshal(reqBody, &newTask)

	conection := conectionBD()
	insertRegister, err := conection.Prepare("Insert into task (name, content) values(?,?)")

	if err != nil {
		panic(err.Error())
	}
	insertRegister.Exec(newTask.Name, newTask.Content)

	getLastRow, err := conection.Query("SELECT last_insert_id()")
	defer getLastRow.Close()

	var lastId int
	for getLastRow.Next() {
		err = getLastRow.Scan(&lastId)
		if err != nil {
			panic(err.Error())
		}
	}
	newTask.ID = lastId

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// getTask this method gets a task from the id
func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalidad ID")
		return
	}

	conection := conectionBD()
	getTask, err := conection.Query("SELECT * from task where id=?", taskID)
	defer getTask.Close()

	if err != nil {
		panic(err.Error())
	}

	taskAux := task{}
	for getTask.Next() {
		var id int
		var name, content string
		err = getTask.Scan(&id, &name, &content)
		if err != nil {
			panic(err.Error())
		}
		taskAux.ID = id
		taskAux.Name = name
		taskAux.Content = content

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(taskAux)
		return
	}
	fmt.Fprintf(w, "Task with ID %v has not been found", taskID)

}

// deleteTask this method delete a task from the id
func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalidad ID")
		return
	}

	conection := conectionBD()

	getTask, err := conection.Query("SELECT id from task where id=?", taskID)
	if err != nil {
		panic(err.Error())
	}

	if !getTask.Next() {
		fmt.Fprintf(w, "Error : ID %v not exist", taskID)
		return
	}

	deleteRegister, err := conection.Prepare("DELETE from task where id=?")
	if err != nil {
		panic(err.Error())
	}

	deleteRegister.Exec(taskID)
	fmt.Fprintf(w, "Task with ID %v has been remove succesfully", taskID)
}

// updateTask this method update a task from the id and the body request
func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updateTask task
	if err != nil {
		fmt.Fprintf(w, "Invalidad ID")
		return
	}

	conection := conectionBD()

	getTask, err := conection.Query("SELECT id from task where id=?", taskID)
	if err != nil {
		panic(err.Error())
	}

	if !getTask.Next() {
		fmt.Fprintf(w, "Error : ID %v not exist", taskID)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please Enter Valid Data")
	}
	json.Unmarshal(reqBody, &updateTask)

	editRegister, err := conection.Prepare("UPDATE task SET name=?,content=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	editRegister.Exec(updateTask.Name, updateTask.Content, taskID)

	fmt.Fprintf(w, "The task with ID %v has been updated successfully", taskID)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("APP_PORT")
	log.Fatal(http.ListenAndServe(":"+port, router))

}
