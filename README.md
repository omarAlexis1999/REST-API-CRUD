This is a mini project in which an API is generated
We consume data from MySQL
We use the symtem1 database, where we create the task table with 3 attributes
- id        int
- name      string
- content   string 

Prerequisites
- Golang
- Create a mysql database with its corresponding table
- Configure the .env file with the corresponding values for the database

Additional packages used in this project
- go get github.com/go-sql-driver/mysql
- go get github.com/gorilla/mux
- go get github.com/joho/godotenv
- go install -mod=mod github.com/githubnemo/CompileDaemon , this is to detect the changes

A CRUD is created
The routes for the API are as follows

Create (POST)
- http://localhost:3000/tasks

Read (GET)
- http://localhost:3000/tasks

Update (PUT)
- http://localhost:3000/tasks/{id}

Remove (DELETE)
- http://localhost:3000/tasks/{id}

To run the project:
- go run main.go
To run and detect changes use:
1. CompileDaemon
2. CompileDaemon -command="./restAPI.exe"
