package handler

import (
	"fmt"
	"net/http"
	"time"
)


type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	Username  string    `json:"username" bson:"username"`
	Image     string    `json:"image" bson:"image"`
	Password  string    `json:"password,omitempty" bson:"password"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}


func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register endpoint hit")
	w.WriteHeader(http.StatusNotImplemented)
}

func(u *User) Login(w http.ResponseWriter, r *http.Request){
	fmt.Println("Login endpoint hit")
	w.WriteHeader(http.StatusNotImplemented)
}
func(u *User) Update(w http.ResponseWriter, r *http.Request){
	fmt.Println("Update endpoint hit")
	w.WriteHeader(http.StatusNotImplemented)
}