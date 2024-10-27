package models

type SignupUser struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type DatabaseInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type SignupAdmin struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginAdmin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Course struct {
	Title       string `json:"title" validate:"required,min=10,max=255"`
	Description string `json:"description" validate:"required,min=100,max=2500"`
	ImageLink   string `json:"image_link" validate:"omitempty,url"`
	AdminID     int    `json:"admin_id" validate:"required"`
}

type Purchase struct {
	UserID   int `json:"user_id" validate:"required"`
	CourseID int `json:"course_id" validate:"required"`
}
