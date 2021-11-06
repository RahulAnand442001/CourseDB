package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Model for Course - file
type Course struct {
	CourseId    string  `json:"id"`
	CourseName  string  `json:"name"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"`
}

// Model for Author - file
type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

// fake DB
var courses = map[string]Course{}

// middlewares , helpers
func (c *Course) IsEmpty() bool {
	return c.CourseName == ""
}

func main() {
	fmt.Println("API - Go Lang !")

	// create router
	r := mux.NewRouter()

	// seed data into map (fake DB)
	injectData()

	// routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", findAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", findOneCourse).Methods("GET")
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/course/{id}", deleteOneCourse).Methods("DELETE")

	// listen to a port
	log.Fatal(http.ListenAndServe(":4000", r))

}

func injectData() {
	// adding data into map
	courses["1"] = Course{"1", "ReactJs", 299, &Author{"Facebook", "https://reactjs.org/"}}
	courses["2"] = Course{"2", "AngularJs", 199, &Author{"Google", "https://angularjs.org/"}}
	courses["3"] = Course{"3", "VueJs", 299, &Author{"Vue", "https://vuejs.org/"}}
	courses["4"] = Course{"4", "Flask", 199, &Author{"Flask", "https://github.com/pallets/flask"}}
	courses["5"] = Course{"5", "Django", 299, &Author{"Facebook", "https://www.djangoproject.com/"}}

}

// serve home route
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to API by GO Lang !</h1>"))
}

func findAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching all courses ...")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func findOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching course by Id ...")
	w.Header().Set("Content-Type", "application/json")

	// extract params
	params := mux.Vars(r)
	if course, ok := courses[params["id"]]; ok {
		fmt.Println("Course founded !")
		json.NewEncoder(w).Encode(course)
		return
	}
	json.NewEncoder(w).Encode("No Course found with given Id")
	return
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create a course")
	w.Header().Set("Content-Type", "application/json")

	// check for empty requests
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
		return
	}

	var course Course

	json.NewDecoder(r.Body).Decode(&course)

	if course.IsEmpty() {
		json.NewEncoder(w).Encode("Please send some data")
		return
	}

	// generate unique id , string
	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100))
	courses[course.CourseId] = course
	json.NewEncoder(w).Encode(course)

}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one course")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	// get the id from the body
	var course Course
	if _, ok := courses[params["id"]]; ok {
		json.NewDecoder(r.Body).Decode(&course)
		courses[params["id"]] = course
		json.NewEncoder(w).Encode("Value updated successfully")
		return
	}
	json.NewEncoder(w).Encode("COurse Not Found ! Invalid Id")
	return

}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one course")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	if _, ok := courses[params["id"]]; ok {
		delete(courses, params["id"])
		json.NewEncoder(w).Encode("Course deleted successfully")
		return
	}
	json.NewEncoder(w).Encode("No course matched by Id ! Unable to delete course.")
	return

}
