package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

type MenuItem struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type Review struct {
	Name     string `json:"name"`
	Dish     string `json:"dish"`
	Rating   int    `json:"rating"`
	Comments string `json:"comments"`
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleMenu(w http.ResponseWriter, r *http.Request) {

	// Here fetch the menu items from the data server
	resp, err := http.Get("http://localhost:4002/data")

	// If there is an error, return an internal server error

	// Close the response body when the function returns

	// Here read the response body

	if err != nil {
		http.Error(w, "Error getting menu items", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, "Error reading menu items", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var menuItems []MenuItem
	err = json.Unmarshal(body, &menuItems)
	if err != nil {
		http.Error(w, "Error decoding menu items", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/menu.html")
	if err != nil {
		http.Error(w, "Error loading menu page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, menuItems)
}

func handleReviewForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/review_form.html")
	if err != nil {
		http.Error(w, "Error loading review form", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleReviewSubmission(w http.ResponseWriter, r *http.Request) {

	// Here parse the form data recieved from the review form

	// Here create a new review object from the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}
	name := r.FormValue("name")
	dish := r.FormValue("dish")
	comments := r.FormValue("comments")
	rating := stringToInt(r.FormValue("rating"))
	review := Review{Name: name, Dish: dish, Rating: rating, Comments: comments}

	reviewData, err := json.Marshal(review)
	if err != nil {
		http.Error(w, "Error encoding review data", http.StatusInternalServerError)
		return
	}

	// Here post the review data to the data server
	resp, err := http.Post("http://localhost:4002/addReview", "application/json", bytes.NewBuffer(reviewData))
	if err != nil {
		http.Error(w, "Error posting review", http.StatusInternalServerError)
		return
	}
	// If there is an error, return an internal server error
	defer resp.Body.Close()

	// Close the response body when the function returns

	http.Redirect(w, r, "/reviews", http.StatusSeeOther)
}

func handleReviews(w http.ResponseWriter, r *http.Request) {

	// Here fetch the reviews from the data server

	// If there is an error, return an internal server error

	// Close the response body when the function returns

	// Here read the response body
	resp, err := http.Get("http://localhost:4002/reviews")

	if err != nil {
		http.Error(w, "Error getting reviews", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, "Error reading menu items", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var reviews []Review
	err = json.Unmarshal(body, &reviews)
	if err != nil {
		http.Error(w, "Error decoding reviews", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/reviews.html")
	if err != nil {
		http.Error(w, "Error loading reviews page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, reviews)
}

func stringToInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
