package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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


type RequestCheckin struct {
	UserID   uint32    `json:"user_id"`
}

type Beacon struct {
	PublicKey  ed25519.PublicKey  `json:"pk"`
	signingKey ed25519.PrivateKey
	Location   string           `json:"location"`
}

func (b *Beacon) Values() url.Values {
	return url.Values{
		KeyLocation: {b.Location},
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
		Signature: ed25519.Sign(b.signingKey, msg.Bytes()),
	}
}

func (b *Beacon) Advertise() error  {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	_, err = http.Post(URLAdvertise, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	return nil
}

func (b *Beacon) Save() error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	if err = os.WriteFile(string(b.Location+".json"), data, 0664); err != nil {
		return err
	}
	return nil
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
	}
}

func (b *Beacon) Load() error {
	buf := make([]byte, 1024)
	f, err := os.Open(b.Location + ".json")
	if err == nil {
		n, err := f.Read(buf)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(buf[:n], &b); err != nil {
			return err
		}
	} else if errors.Is(err, fs.ErrNotExist){
		b.PublicKey, b.signingKey, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		b.Save()
	} else {
		return err
	}
	return nil
}

func main() {
	var beacon Beacon
	flag.StringVar(&beacon.Location, "location", "Central Park", "the name of the ")
	flag.Parse()
	if err := beacon.Load(); err != nil {
		log.Panicln(err)
	}
	if err := beacon.Advertise(); err != nil {
		log.Panicln(err)
	}
	fmt.Println(beacon.Location)

	http.HandleFunc("/checkin", checkinHandler(&beacon))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
