package dtos

type UserProfileDTO struct {
	Username      string
	Email         string
	NumberOfNotes int
}

func New(username, email string, numberOfNotes int) *UserProfileDTO {
	return &UserProfileDTO{
		Username:      username,
		Email:         email,
		NumberOfNotes: numberOfNotes,
	}
}
