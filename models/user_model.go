package models

type RegisterUserRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type LoginUserRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginUser struct {
	ID       string `json:"id" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
}

type GetUserResponseBody struct {
	ID    string `json:"id" bson:"_id"`
	Email string `json:"email" bson:"email"`
	Name  string `json:"name" bson:"name"`
}
