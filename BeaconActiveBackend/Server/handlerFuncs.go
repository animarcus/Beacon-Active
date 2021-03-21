package Server

import (
	"BeaconActive/Model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)



func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Homepage \n\n")
	user := Model.User{}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Header", "Content-Type")
	w.WriteHeader(http.StatusAccepted)
	dec := json.NewDecoder(r.Body)
	_ = dec.Decode(&user)
	err := user.AddUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("New user id %d from given name %s \n", user.Id, user.UserName)
	_ = json.NewEncoder(w).Encode(user)
}

func advertise(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Advertise, receiving beacon... \n")
	dec := json.NewDecoder(r.Body)
	beacon := Model.Beacon{}
	_ = dec.Decode(&beacon)
	fmt.Printf("Received beacon from location %s\n", beacon.Location)
	err := beacon.AddBeacon()
	if err != nil {
		fmt.Printf("Erroradding beacon: %s\n", err)
		return
	}
}


func checkin(w http.ResponseWriter, r *http.Request) {



	//w.WriteHeader(200)
	fmt.Printf("Checking in... \n")
	dec := json.NewDecoder(r.Body)
	msg := Model.Message{}
	_ = dec.Decode(&msg)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Header", "Content-Type")
	w.WriteHeader(http.StatusAccepted)

	fmt.Printf("Received message from user at location %s at time %s\n", msg.Location, msg.Time)
	err := msg.AddMessage()
	fmt.Printf("Adding msg to database... \n")
	if err != nil {
		fmt.Printf("Error while adding msg to database: %s\n", err)
		return
	}

	activity := Model.Activity{}
	fmt.Printf("Checking in activity from msg... \n")
	err = Model.CheckIn(&activity, &msg)
	if err != nil {
		fmt.Printf("Error while checking in message: %s\n", err)
		return
	}
	OpenActivities = append(OpenActivities, &activity)
	fmt.Printf("All open activities: \n")
	for _, act := range OpenActivities {
		fmt.Printf("\t Activity: user %s, location %s, time %s\n", act.UserId, act.Location, act.Checkin.Time)
	}
}

func checkout(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Checking out... \n")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Header", "Content-Type")
	w.WriteHeader(http.StatusAccepted)
	dec := json.NewDecoder(r.Body)
	msg := Model.Message{}
	_ = dec.Decode(&msg)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Printf("Received message from user at location %s at time %s\n", msg.Location, msg.Time)
	err := msg.AddMessage()
	fmt.Printf("Adding msg to database... \n")
	if err != nil {
		fmt.Printf("Error while adding msg to database: %s\n", err)
		return
	}

	fmt.Printf("Checking out activity from msg... \n")

	activityIndex := 0
	for i, act := range OpenActivities {
		if act.UserId == msg.UserId  {
			err = Model.CheckOut(act, &msg)
			if err == nil {
				activityIndex = i
			}
		}
	}
	copy(OpenActivities[activityIndex:], OpenActivities[activityIndex+1:])
	OpenActivities[len(OpenActivities)-1] = &Model.Activity{}
	OpenActivities = OpenActivities[:len(OpenActivities)-1]
	fmt.Printf("All open activities: \n")
	for _, act := range OpenActivities {
		fmt.Printf("\t Activity: user %s, location %s, time %s\n", act.UserId, act.Location, act.Checkin)
	}
}

func activities(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Activities... \n")
	user := Model.User{}
	dec := json.NewDecoder(r.Body)
	_ = dec.Decode(&user)
	w.Header().Set("Content-Type", "application/json")
	allActivities := Model.GetAllActivities()
	allOutActs := []OutputActivity{}
	for _, act := range allActivities {
		if act.UserId == user.Id {
			outAct := OutputActivity{
				Location: act.Location,
				UserId:   user.Id,
				Checkin:  act.Checkin.Time,
				Checkout: act.Checkout.Time,
			}
			fmt.Printf("\t Appending ctivity: user %s, location %s, time %s\n", act.UserId, act.Location, act.Checkin)
			allOutActs = append(allOutActs, outAct)
		}
	}
	j, _ := json.Marshal(allOutActs)
	fmt.Printf("Json format: %s ", j)
	_, _ = w.Write(j)
}

type OutputActivity struct {
	Location Model.Location `json:"location"`
	UserId Model.UserId `json:"user_id"`
	Checkin time.Time `json:"checkin"`
	Checkout time.Time `json:"checkout"`
}

