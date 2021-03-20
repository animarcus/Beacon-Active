package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mdp/qrterminal/v3"
)

const (
	URLAdvertise = "http://philips-macbook.local:8080/advertise"
)

const (
	KeyLocation     = "location"
	KeyBeaconPublic = "pk"
	KeyUserID       = "user_id"
	KeyTime         = "time"
	KeySignature    = "sig"
)

type (
	Location string
)

type RequestCheckin struct {
	UserID   uint32    `json:"user_id"`
}

type Beacon struct {
	PublicKey  ed25519.PublicKey  `json:"pk"`
	SigningKey ed25519.PrivateKey `json:"sk"`
	Location   Location           `json:"location"`
}

func (b *Beacon) Values() url.Values {
	return url.Values{
		KeyLocation: {string(b.Location)},
		KeyBeaconPublic: {hex.EncodeToString(b.PublicKey)},
	}
}
func (b *Beacon) NewMessage(userID uint32) *Message {
	return &Message{
		UserID:   userID,
		Time:     time.Now(),
		Location: b.Location,
	}
}
func (b *Beacon) NewSignedMessage(msg *Message) *SignedMessage {
	return &SignedMessage{
		Message:   *msg,
		Signature: ed25519.Sign(b.SigningKey, msg.Bytes()),
	}
}

func (b *Beacon) Advertise() {
	data, err := json.Marshal(b)
	if err != nil {
		log.Panicln(err)
		return
	}
	_, err = http.Post(URLAdvertise, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Panicln(err)
		return
	}
}

func (b *Beacon) Save() {
	data, err := json.Marshal(b)
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile(string(b.Location+".json"), data, 0664)
	if err != nil {
		log.Println(err)
	}
}

func checkinHandler(beacon *Beacon) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Declare a new Person struct.
		var checkin RequestCheckin

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&checkin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		msg := beacon.NewSignedMessage(beacon.NewMessage(checkin.UserID))
		msgStr, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		qrterminal.Generate(string(msgStr), qrterminal.L, w)
		//qrterminal.Generate(msg.Values().Encode(), qrterminal.L, os.Stdout)
	}
}

func main() {
	var beacon Beacon

	location := *flag.String("location", "Central Park", "the name of the ")
	buf := make([]byte, 1024)
	f, err := os.Open(location + ".json")
	if err == nil {
		n, err := f.Read(buf)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(buf[:n], &beacon)
		if err != nil {
			log.Println(err)
		}
	} else if errors.Is(err, fs.ErrNotExist){
		beacon.PublicKey, beacon.SigningKey, _ = ed25519.GenerateKey(rand.Reader)
		beacon.Location = Location(location)
		beacon.Save()
	} else {
		log.Panicln(err)
	}
	http.HandleFunc("/checkin", checkinHandler(&beacon))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
