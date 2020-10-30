package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
)

var logs = logrus.New()

func findTasks(inputID ResponseID) (output []Task, err error) {
	for _, user := range users {
		if user.ID == inputID.UserID {
			for _, task := range user.Tasks {
				output = append(output, task)
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return output, nil
}

func outputTask(w http.ResponseWriter, r *http.Request) (error, []Task) {
	var inputID ResponseID
	err := json.NewDecoder(r.Body).Decode(&inputID)
	if err != nil {
		errors.Wrap(err, "Cant decode file.")
		return err, nil
	}
	output, findErr := findTasks(inputID)
	if findErr != nil {
		errors.Wrap(findErr, "Cant find tasks for current ID.")
		return findErr, nil
	}
	return nil, output
}

//Get all tasks - our /list
func getTasks(w http.ResponseWriter, r *http.Request) {
	err, output := outputTask(w, r)
	if err != nil {
		logs.Error(err, "Cant output tasks")
	}
	errEncode := json.NewEncoder(w).Encode(output)
	if errEncode != nil {
		errEncode := errors.Wrap(errEncode, "Cant encode tasks and send to bot.")
		logs.Warn(errEncode)
	}
	//TODO: Logs will be here soon
}

func FindUser(id string) User {
	for _, us := range users {
		if us.ID == id {
			return us
		}
	}
	users = append(users, User{
		ID:    id,
		Tasks: nil,
	})
	return users[len(users)-1]
}

func addTask(w http.ResponseWriter, r *http.Request) (User, error) {
	var task ResponseTask
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		//TODO: HANDLE ERROR AND MAKE RESPONSE
		errors.Wrap(err, "Cant decode current task.")
		logs.Warn(err)
		return User{ID: "", Tasks: nil}, err
	}
	//TODO: Return only user
	ourUser := FindUser(task.UserID)
	result := Task{
		ID:    strconv.Itoa(len(ourUser.Tasks) + 1),
		Title: task.UserTask,
	}
	ourUser.Tasks = append(ourUser.Tasks, result)
	return ourUser, nil
}

//createTask - должна только вызывать методы компонетнов бизнес логики и выводит ошибки
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	ourUser, err := addTask(w, r)
	if err != nil {
		errors.Wrap(err, "Cant decode current task.")
		logs.Warn(err)
	}
	//TODO: Write id generator or use libs for this
	//users[id] = ourUser
}

//Task struct
type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ResponseTask struct {
	UserID   string `json:"user_id"`
	UserTask string `json:"user_task"`
}

type ResponseID struct {
	UserID string `json:"user_id"`
}

type User struct {
	ID    string `json:"id"`
	Tasks []Task `json:"tasks"`
}

//Init users var as a slice user struct
var users []User

func main() {
	//Init router
	r := mux.NewRouter()
	//Mock data
	users = append(users, User{
		ID: "1",
		Tasks: []Task{{
			ID:    "1",
			Title: "Test title",
		}},
	})
	r.HandleFunc("/api/tasks/", getTasks).Methods("POST")
	r.HandleFunc("/api/tasks/create/", createTask).Methods("POST")
	http.ListenAndServe(":8000", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}
