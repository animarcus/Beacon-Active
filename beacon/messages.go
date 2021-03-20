package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"net/url"
	"strconv"
	"time"
)

type Message struct {
	UserID   uint32    `json:"user_id"`
	Time     time.Time `json:"time"`
	Location string  `json:"location"`
}

type SignedMessage struct {
	Message
	Signature []byte `json:"sig"`
}

func (msg *Message) Values() url.Values {
	values := url.Values{}
	values.Set(KeyUserID, strconv.Itoa(int(msg.UserID)))
	values.Set(KeyTime, msg.Time.String())
	values.Set(KeyLocation, string(msg.Location))
	return values
}

func (msg *SignedMessage) Values() url.Values {
	values := msg.Message.Values()
	values.Set(KeySignature, hex.EncodeToString(msg.Signature))
	return values
}

func (msg *Message) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 32+64+len(msg.Location)))
	binary.Write(buf, binary.BigEndian, msg.UserID)
	binary.Write(buf, binary.BigEndian, msg.Time.Unix())
	buf.Write([]byte(msg.Location))
	return buf.Bytes()
}
