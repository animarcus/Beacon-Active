package Server

import (
	"BeaconActive/Model"
	"log"
	"net/http"
)

var OpenActivities []*Model.Activity

func StartServer() {
	server := http.Server{Addr: ":8080"}
	http.HandleFunc("/", homepage)
	http.HandleFunc("/advertise", advertise)
	http.HandleFunc("/checkin", checkin)
	http.HandleFunc("/checkout", checkout)
	http.HandleFunc("/actvities", activities)
	log.Fatal(server.ListenAndServe())
}
