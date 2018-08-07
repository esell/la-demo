package main

type SubmitterList struct {
	Submissions []Submitter `json:"submissions"`
}

type Submitter struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Topic  string `json:"topic"`
	Status string `json:"status"`
}
