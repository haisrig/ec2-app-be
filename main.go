package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"./ec2data"
	"github.com/gorilla/mux"
)

func login(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user map[string]interface{}
	json.Unmarshal([]byte(string(reqBody)), &user)
	if user["username"] == "srikanth" && user["password"] == "srikanth" {
		w.WriteHeader(200)
		fmt.Println("Done")
		return
	}
	w.WriteHeader(401)
	return
}

func getData(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	owner := query["owner"][0]
	json.NewEncoder(w).Encode(ec2data.GetEc2Data(owner))
}

func startInstance(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	instanceId := query["instanceId"][0]
	statusCode := ec2data.StartInstance(instanceId)
	w.WriteHeader(statusCode)
	return
}

func stopInstance(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	instanceId := query["instanceId"][0]
	statusCode := ec2data.StopInstance(instanceId)
	w.WriteHeader(statusCode)
	return
}

func terminateInstance(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	instanceId := query["instanceId"][0]
	statusCode := ec2data.TerminateInstance(instanceId)
	w.WriteHeader(statusCode)
	return
}

func main() {
	fmt.Println("Listening on port 8000")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/login", login).Methods("POST")
	myRouter.HandleFunc("/get", getData)
	myRouter.HandleFunc("/start", startInstance).Methods("POST")
	myRouter.HandleFunc("/stop", stopInstance).Methods("POST")
	myRouter.HandleFunc("/terminate", terminateInstance).Methods("POST")
	fmt.Println(&myRouter)
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
