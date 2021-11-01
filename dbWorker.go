package main

import (
	"net/http"
)

type dbWorker interface {
	PullKicks(w http.ResponseWriter, r *http.Request)
	PullKick(w http.ResponseWriter, r *http.Request)
	PostKick(w http.ResponseWriter, r *http.Request)
	PostKicks(w http.ResponseWriter, r *http.Request)
	DeleteKick(w http.ResponseWriter, r *http.Request)
	DeleteKicks(w http.ResponseWriter, r *http.Request)
	SetupConnection() error
}
