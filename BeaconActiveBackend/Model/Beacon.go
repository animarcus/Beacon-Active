package Model

import (
	"crypto/ed25519"
	"encoding/hex"
)

type Location string

type Beacon struct {
	Id int
	Location Location `json:"location"`
	PublicKey ed25519.PublicKey `json:"pk"`
}

func (beacon *Beacon) AddBeacon() (err error) {
	statement := "insert into beaconactivedb.public.beacons (publickey, location) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(hex.EncodeToString(beacon.PublicKey), beacon.Location).Scan(&beacon.Id)
	return
}

func GetAllBeacons() (beacons map[Location]*Beacon, err error) {
	beacons = map[Location]*Beacon{}
	rows, err := Db.Query("select id, publickey, location from beaconactivedb.public.beacons limit $1", 10000)
	if err != nil {
		return
	}
	for rows.Next() {
		beacon := Beacon{}
		pk := ""
		err = rows.Scan(&beacon.Id, &pk, &beacon.Location)
		if err != nil {
			return
		}
		beacon.PublicKey, err = hex.DecodeString(pk)
		if err != nil {
			return
		}
		beacons[beacon.Location] = &beacon
	}
	err = rows.Close()
	return
}

