package HttpVerbs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// course struct
type Course struct {
	Id               int
	Name             string
	Duration         int     // in days
	Fee              float64 // course fee
	Features         CourseFeatures
	Content          Content
	Interaction      Interaction
	Feedback         bool
	ProgressTracking bool
}

type CourseFeatures struct {
	Quizzes     bool
	Assignments bool
	Discussions bool
}

type Content struct {
	Videos        bool
	Presentations bool
	Simulations   bool
}

type Interaction struct {
	Students    bool
	Instructors bool
}

// All listed courses
var Courses []Course

// middleware
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do stuff before the handlers
		token := r.Header.Get("token")      // jwt/ paceto token for authentication
		if len(token) == 0 || token == "" { // unauthorised
			fmt.Println("Empty token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// middle wares can be consumed to log any common information which is application for all the routes
		h.ServeHTTP(w, r)
		// do stuff after the hadlers

	})
}

// custom middleware in case we want to consume
func Middleware2(s string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do stuff
			fmt.Println(s)
			h.ServeHTTP(w, r)
		})
	}
}

func StartServer() {
	r := mux.NewRouter()
	//seeding the data
	// creating DSA course
	course := Course{}
	course.Id = 1
	course.Name = "DSA in depth"
	course.Duration = 180 // days
	course.Fee = 7999     //Rs
	features := CourseFeatures{}
	features.Quizzes = true
	features.Assignments = true
	features.Discussions = true
	course.Features = features
	content := Content{}
	content.Videos = true
	content.Presentations = true
	content.Simulations = false
	course.Content = content
	interaction := Interaction{}
	interaction.Instructors = true
	interaction.Students = false
	course.Interaction = interaction
	Courses = append(Courses, course)

	// creating DSA course
	course = Course{}
	course.Id = 2
	course.Name = "HLD in depth"
	course.Duration = 90 // days
	course.Fee = 3999    //Rs
	features = CourseFeatures{}
	features.Quizzes = true
	features.Assignments = true
	features.Discussions = true
	course.Features = features
	content = Content{}
	content.Videos = true
	content.Presentations = false
	content.Simulations = false
	course.Content = content
	interaction = Interaction{}
	interaction.Instructors = true
	interaction.Students = true
	course.Interaction = interaction
	Courses = append(Courses, course)

	// call middleware first to check authetication
	r.Use(Middleware)
	//routing
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/courses", getCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getCourse).Methods("GET")

	// listen to a port
	fmt.Println("Starting a mux server , Port = 4000")
	log.Fatal(http.ListenAndServe(":4000", r))
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	// data insertion in to database
	// success
	w.Header().Set("Content-Type", "application/json")
	var course Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	course.Id = len(Courses) + 1
	Courses = append(Courses, course)
	json.NewEncoder(w).Encode(course)
	w.WriteHeader(http.StatusOK)
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Courses)
	w.WriteHeader(http.StatusOK)
}

func getCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, course := range Courses {
		if course.Id == id {
			json.NewEncoder(w).Encode(course)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
