package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/resend/resend-go/v2"
	"github.com/restingdemon/thaparEvents/utils"
)

type Feedbacks struct {
	Name    string `json:"name" bson:"name"`
	Email   string `json:"email" bson:"email"`
	Contact string `json:"contact" bson:"contact"`
	Message string `json:"message" bson:"message"`
}

func Feedback(w http.ResponseWriter, r *http.Request) {
	var feedback = &Feedbacks{}

	utils.ParseBody(r, feedback)

	if feedback.Contact == "" || feedback.Email == "" || feedback.Message == "" || feedback.Name == "" {
		http.Error(w, "Not a valid request ", http.StatusBadRequest)
		return
	}
	message := fmt.Sprintf("Feedback Recieved From %s with Contact No %s and email %s", feedback.Name, feedback.Contact, feedback.Email)
	apiKey := os.Getenv("Feedback_KEY")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Thapar Events <onboarding@resend.dev>",
		To:      []string{"agarg8_be21@thapar.edu"},
		Text:    feedback.Message,
		Subject: message,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(sent.Id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Feedback sent succesfully!"))
	
}
