package models

type Task struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Doing bool   `json:"doing"`
	Done  bool   `json:"done"`
}
