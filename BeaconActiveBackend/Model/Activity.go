package Model

import (
	"fmt"
	"time"
)

type Activity struct {
	Id int
	Location Location
	UserId UserId `json:"user_id"`
	Checkin *Message
	Checkout *Message
}

const (
	MIN_INTERVAL = 10*time.Minute
	MAX_INTERVAL = 60*time.Minute
)

type InvalidCheckInError struct {
	*Message
}

func(err *InvalidCheckInError) Error() string {
	return fmt.Sprintf("Invalid checkin for message location %s at time %s", err.Location, err.Time)
}

func CheckIn(activity *Activity, msg *Message) (err error) {
	if activity.Checkin != nil {
		err = &InvalidCheckInError{msg}
		return
	}
	activity.Location = msg.Location
	activity.UserId = msg.UserId
	activity.Checkin = msg
	return
}

type InvalidCheckOutError struct {
	*Message
}

func (err *InvalidCheckOutError) Error() string {
	return fmt.Sprintf("Invalid checkout for message location %s at time %s", err.Location, err.Time)
}

func CheckOut(activity *Activity, msg *Message) (err error) {
	t0 := activity.Checkin.Time
	t1 := msg.Time
	//diff := t1.Sub(t0)
	if t0.Before(t1) && /*diff>MIN_INTERVAL && diff<MAX_INTERVAL &&*/ activity.UserId == msg.UserId && activity.Location == msg.Location && activity.Checkout == nil {
		activity.Checkout = msg
		err = activity.AddActivity()
	} else {
		err = &InvalidCheckOutError{msg}
	}
	return
}

func (activity *Activity) AddActivity() (err error) {
	statement := "insert into beaconactivedb.public.activities (userid, location, checkin, checkout) values ($1, $2, $3, $4) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(activity.UserId, activity.Location, activity.Checkin.Time, activity.Checkout.Time).Scan(&activity.Id)
	return
}

func GetAllActivities() (activities []*Activity) {
	rows, err := Db.Query("select id, userid, location, checkin, checkout from beaconactivedb.public.activities limit $1", 10)
	if err != nil {
		return
	}
	for rows.Next() {
		activity := Activity{}
		userid := UserId(0)
		var t0, t1 time.Time

		err = rows.Scan(&activity.Id, &activity.UserId, &activity.Location, &t0, &t1)
		if err != nil {
			return
		}

		m0 := &Message{}
		m0, err = GetMessage(userid, t0)
		if err != nil {
			return
		}
		m1 := &Message{}
		m1, err = GetMessage(userid, t1)
		if err != nil {
			return
		}
		activity.Checkin = m0
		activity.Checkout = m1
		activities = append(activities, &activity)
	}
	err = rows.Close()
	return
}
