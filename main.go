package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/signup", UserSignup).Methods("POST")
	userRouter.HandleFunc("/login", UserLogin).Methods("POST")
	userRouter.HandleFunc("/courses", SeeAllCourse).Methods("GET")
	userRouter.HandleFunc("/purchase", PurchaseCourse).Methods("POST")
	userRouter.HandleFunc("/purchases-courses", SeeAllPurchasesCourse).Methods("GET")

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/signup", AdminSignup).Methods("POST")
	adminRouter.HandleFunc("/login", AdminLogin).Methods("POST")
	adminRouter.HandleFunc("/create", CreateCourse).Methods("POST")
	adminRouter.HandleFunc("/delete", DeleteCourse).Methods("DELETE")
	adminRouter.HandleFunc("/add", AddCourse).Methods("POST")

	fmt.Println("Starting the server at port 8000")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println("Error while starting the server")
	}
}

func UserSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User Signup")
}
func UserLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User Login")
}
func SeeAllCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Courses")
}
func PurchaseCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Buy from these")
}
func SeeAllPurchasesCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Your Purchased Courses")
}

func AdminSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Admin Singup")
}
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Admin Login")
}
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create Course")
}
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Course Deleted")
}
func AddCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Course Added")
}
