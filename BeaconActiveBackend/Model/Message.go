package Model

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type InvalidMessageError struct {
	*Message
}

func (err InvalidMessageError) Error() string {
	return fmt.Sprintf("Message from user %d at location %s at time %s cannot be verified", err.UserId, err.Location, err.Time.String())
}

type Message struct {
	Id int
	UserId UserId `json:"user_id"`
	Time time.Time `json:"time"`
	Location Location `json:"location"`
	Signature []byte `json:"sig"`
}

func (msg *Message) AddMessage() error {
	if ok := msg.Verify(); !ok {
		return InvalidMessageError{Message: msg}
	}
	statement := "insert into beaconactivedb.public.messages (userid, time, signature, location) values ($1, $2, $3, $4) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		fmt.Println("bad philip")
		fmt.Println("bad philip")
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(msg.UserId, msg.Time, hex.EncodeToString(msg.Signature), msg.Location).Scan(&msg.Id)
	return nil
}

func GetAllMessages() (messages map[UserId]*Message, err error) {
	messages = map[UserId]*Message{}
	rows, err := Db.Query("select id, userid, time, signature, location from beaconactivedb.public.messages limit $1", 100000)
	if err != nil {
		return
	}
	for rows.Next() {
		msg := Message{}
		sig := ""
		err = rows.Scan(&msg.Id, &msg.UserId, &msg.Time, &sig, &msg.Location)
		if err != nil {
			return
		}
		msg.Signature, err = hex.DecodeString(sig)
		if err != nil {
			return
		}
		messages[msg.UserId] = &msg
	}
	err = rows.Close()
	return
}

func GetMessageFrom(user *User) (messages []*Message, err error) {
	msgs := map[UserId]*Message{}
	msgs, err = GetAllMessages()
	for uid, msg := range msgs {
		if uid == user.Id {
			messages = append(messages, msg)
		}
	}
	return
}

func GetMessage(id UserId, time time.Time)  (msg *Message, err error) {
	messages := map[UserId]*Message{}
	messages, err = GetAllMessages()
	msg = messages[id]
	return
}

func (msg *Message) Verify() bool {
	beacons, err := GetAllBeacons()
	if err != nil {
		panic(err)
	}
	beacon, ok := beacons[msg.Location]
	if !ok {
		fmt.Println("bad philip")
		fmt.Println("bad philip", msg.Location)
		//return false
	}
	data := msg.Bytes()
	fmt.Println(data)
	return ed25519.Verify(beacon.PublicKey, data, msg.Signature)
}

func (msg *Message) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 32+64+len(msg.Location)))
	binary.Write(buf, binary.BigEndian, msg.UserId)
	binary.Write(buf, binary.BigEndian, msg.Time.Unix())
	buf.Write([]byte(msg.Location))
	return buf.Bytes()
}

