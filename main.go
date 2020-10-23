package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

//Get all tasks - our /list
func getTasks(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "application/json")
	//json.NewEncoder(w).Encode(tasks)
	r.ParseForm()
	s:=r.Form
	a:=s["id"][0]
	for _, user := range users {
		if user.ID == a {
			for _, task := range user.Tasks {
			fmt.Println(user.ID)
			w.Write([]byte(task.ID +". " + task.Title + "\n"))
			}
		}
	}
}


//Get single task
/*func getTask(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r) //Get params
	//Loop through tasks and find with id
	for _, item := range users {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Task{})
}*/


func FindUser(id string) (User,int){
	for i,us := range users{
		if(us.ID==id) {
			return us, i
		}
	}
	fmt.Println("uyfhdfhhghgd")
	users=append(users, User{
		ID:    id,
		Tasks: nil,
	})
	return users[len(users)-1],len(users)-1
}



//Create a new task - our /add
func createTask(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "application/json")
	var task ResponseTask
	_ = json.NewDecoder(r.Body).Decode(&task)
	var ourUser, index = FindUser(task.UserID)
	var result Task
	result = Task{
		ID:     strconv.Itoa(len(ourUser.Tasks) + 1),
		Title:  task.UserTask,
	}
	ourUser.Tasks = append(ourUser.Tasks, result)
	users[index]=ourUser
	fmt.Println(users)
	}


//Update task
/*func updateTask(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "application/json")
	r.ParseForm()
	s:=r.Form
	a:=s["id"][0]
	for index, item := range tasks {
		if item.ID == a {
			tasks = append(tasks[:index], tasks[index+1:]...)
			var task Task
			_ = json.NewDecoder(r.Body).Decode(&task)
			task.ID = a //Mock ID
			tasks = append(tasks, task)
			json.NewEncoder(w).Encode(task)
			return
		}
	}
}*/


//Delete task - our /do
/*func deleteTask(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "application/json")
	r.ParseForm()
	s:=r.Form
	a:=s["id"][0]
	for index, item := range tasks {
		if item.ID == string(a) {
			tasks = append(tasks[:index], tasks[index+1:]...)
			break

		}
	}
}*/

//Task struct
type Task struct {
	ID 		string 	`json:"id"`
	Title 	string 	`json:"title"`
}

type ResponseTask struct {
	UserID 	 string `json:"user_id"`
	UserTask string	`json:"user_task"`
}

type User struct {
	ID 	string 	 `json:"id"`
	Tasks []Task `json:"tasks"`
}




//Init tasks var as a slice user struct
var users []User

func main() {

	//Init router
	r := mux.NewRouter()

	//Mock data
	users = append(users, User{
		ID:     "1",
		Tasks: []Task{{
			ID:    "1",
			Title: "Test title",
		}},
	})


	//Route handlers / endpoints
	r.HandleFunc("/api/tasks/", getTasks).Methods("GET")
	//r.HandleFunc("/api/tasks/{id}", getTask).Methods("GET")
	r.HandleFunc("/api/tasks/create/", createTask).Methods("POST")
	//r.HandleFunc("/api/tasks/{id}", updateTask).Methods("PUT")
	//r.HandleFunc("/api/tasks/do/", deleteTask).Methods("DELETE")
	http.ListenAndServe(":8000", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}



