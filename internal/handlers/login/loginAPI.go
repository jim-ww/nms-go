package login

import (
	"fmt"
	"net/http"

	"github.com/jim-ww/nms-go/internal/dtos"
	"github.com/jim-ww/nms-go/internal/validators"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	dto := dtos.LoginDTO{
		Username: username,
		Password: password,
	}
	fmt.Println("loginData:", dto)

	// TODO validate

	// TODO add authentication logic

}

func Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	dto := &dtos.RegisterDTO{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}
	validators.ValidateRegisterDTO(dto)
	// TODO validate

	// TODO check if username or email already exists

	// TODO add authentication logic
}
