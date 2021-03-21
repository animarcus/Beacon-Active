package Model

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres",
		"user=philip dbname=beaconactivedb password=marcus01 sslmode=disable")
	if err != nil {
		panic(err)
	}

	//beacons, err := GetAllBeacons()
	//if len(beacons) == 0 {
	//	cities := []Location{"Gy", "Vandoeuvres", "Meinier", "Vesenaz"}
	//	for _, city := range cities {
	//		pk, _, err := ed25519.GenerateKey(rand.Reader)
	//		if err != nil {
	//			panic(err)
	//		}
	//		beacon := Beacon{Location: city, PublicKey: pk}
	//		err = beacon.AddBeacon()
	//		if err != nil {
	//			fmt.Print(err)
	//		}
	//	}
	//}
}

func SetTimeZone(tz string) {
	statement := "set timezone = " + tz
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
}
