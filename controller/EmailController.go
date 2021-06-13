package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gitlab.jit.team/smurfs/email-service/entity"
	"gitlab.jit.team/smurfs/email-service/repository"
	gomail "gopkg.in/mail.v2"
)

type EmailController struct {
	emailDialer     *gomail.Dialer
	emailRepository *repository.EmailRepository
}

func NewEmailController(repo *repository.EmailRepository) *EmailController {
	emailController := new(EmailController)
	emailController.emailDialer = gomail.NewDialer("smtp.office365.com", 587, "timeports2@jit.team", os.Getenv("EMAIL_PASSWORD"))
	emailController.emailRepository = repo
	return emailController
}

func (controller EmailController) Register(r *mux.Router) {
	r.HandleFunc("/email", controller.sendEmail).Methods(http.MethodPost)
	r.HandleFunc("/email", controller.findAll).Methods(http.MethodGet)
	fmt.Println("Registered Email controller")
}

func (controller *EmailController) sendEmail(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var request entity.Email
	json.Unmarshal(reqBody, &request)

	email := gomail.NewMessage()
	email.SetHeader("From", "timeports2@jit.team")
	email.SetHeader("To", request.Receiver)
	email.SetHeader("Subject", request.Title)
	email.SetBody("text/plain", request.Content)

	log.Println(fmt.Sprintf(`Sending mail '%s' to '%s'`, request.Title, request.Receiver))
	if err := controller.emailDialer.DialAndSend(email); err != nil {
		fmt.Println(err)
		panic(err)
	}

	controller.emailRepository.Save(request)
	w.WriteHeader(http.StatusCreated)
}

func (controller *EmailController) findAll(w http.ResponseWriter, r *http.Request) {
	all := controller.emailRepository.FindAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(all)
}

type EmailRequest struct {
	Subject string `json:"subject"`
	To      string `json:"to"`
	Content string `json:"content"`
}
