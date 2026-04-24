package models

type RegisterUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserRequestBody struct {
	ID string `json:"id"`
}

type LoginUser struct {
	ID       string `json:"id" bson:"_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type GetUserResponseBody struct {
	ID    string `json:"id" bson:"_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
