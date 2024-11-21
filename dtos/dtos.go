package dtos

type TaskDTO struct {
	Id    string `json:"id" validate:"required"`
	Title string `json:"title" validate:"required,min=4"`
	Doing bool   `json:"doing"`
	Done  bool   `json:"done"`
}