package main

import (
	//"fmt"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func main() {

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello Golang")
// 	})

// 	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
// 		fmt.Println(err.Error())

// 	}
// }

func main() {
	router := gin.Default()
	router.GET("/batteries", getBattery)
	router.GET("/batteries/:serialid", batteryByID)
	router.PATCH("/rent", rentBattery)
	router.PATCH("/return", returnBattery)
	router.POST("/batteries", createBatteries)
	router.POST("/register", createUsers)
	router.Run("localhost:8080")
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------//
//
//	Battery
//
// ----------------------------------------------------------------------------------------------------------------------------------------------------------------//
type battery struct {
	SerialId  string `json:"serialid"`
	ModelName string `json:"modelname"`
	Specs     string `json:"specs"`
	Status    string `json:"status"`
}

var batteries = []battery{
	//Available = Can be rent
	//Not Available = Cannot be rent malfunctioned or to be replaced
	//Reserved = Rented to someone
	{SerialId: "00001", ModelName: "Lithium-ion", Specs: "50.26V-26.1Ah", Status: "Reserved"},
	{SerialId: "00023", ModelName: "Lithium-ionV2", Specs: "50.26V-26.1Ah", Status: "Available"},
	{SerialId: "00100", ModelName: "Lithium-ionV3", Specs: "50.26V-26.1Ah", Status: "Not Available"},
}

func getBattery(b *gin.Context) {
	b.IndentedJSON(http.StatusOK, batteries) //200
}

func createBatteries(b *gin.Context) {
	var newBattery battery
	if err := b.BindJSON(&newBattery); err != nil {
		return
	}

	batteries = append(batteries, newBattery)
	b.IndentedJSON(http.StatusCreated, newBattery) //201
}

func getBatteryBySerialID(serialid string) (*battery, error) {
	for i, b := range batteries {
		if b.SerialId == serialid {
			return &batteries[i], nil
		}
	}
	return nil, errors.New("result not found")
}

func batteryByID(b *gin.Context) {
	serialid := b.Param("serialid")
	battery, err := getBatteryBySerialID(serialid)

	if err != nil {
		b.IndentedJSON(http.StatusNotFound, gin.H{"message": "Battery not found."})
		return
	}
	b.IndentedJSON(http.StatusOK, battery)
}

func rentBattery(b *gin.Context) {
	serialid, ok := b.GetQuery("serialid")

	if !ok {
		b.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing serial id "})
		return
	}

	battery, err := getBatteryBySerialID(serialid)

	if err != nil {
		b.IndentedJSON(http.StatusNotFound, gin.H{"message": "Battery not found."})
		return
	}

	switch battery.Status {
	case "Available":
		battery.Status = "Reserved"
		b.IndentedJSON(http.StatusOK, gin.H{"msg": "Successfully rent battery", "Serial": battery.SerialId})

		return
	case "Reserved":
		b.IndentedJSON(http.StatusOK, gin.H{"msg": "Battery not available.", "Serial": battery.SerialId})
		return
	case "Not Available":
		b.IndentedJSON(http.StatusOK, gin.H{"msg": "Battery needs maintenance.", "Serial": battery.SerialId})
		return
	}
}

func returnBattery(b *gin.Context) {
	serialid, ok := b.GetQuery("serialid")

	if !ok {
		b.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing serial id "})
		return
	}

	battery, err := getBatteryBySerialID(serialid)

	if err != nil {
		b.IndentedJSON(http.StatusNotFound, gin.H{"message": "Battery not found."})
		return
	}

	switch battery.Status {
	case "Available":
		b.IndentedJSON(http.StatusOK, gin.H{"msg": "Battery is available", "Serial": battery.SerialId})

		return
	case "Reserved":
		battery.Status = "Available"
		b.IndentedJSON(http.StatusOK, gin.H{"msg": "Battery now returned.", "Serial": battery.SerialId})
		return
	case "Not Available":
		battery.Status = "Available"
		b.IndentedJSON(http.StatusOK, gin.H{"msg": "Battery now available.", "Serial": battery.SerialId})
		return
	}
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------//
//
//	User
//
// ----------------------------------------------------------------------------------------------------------------------------------------------------------------//
type user struct {
	UID        string `json:"uid"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	ReservedBt string `json:"-"`
}

var users = []user{
	//Available = Can be rent
	//Not Available = Cannot be rent malfunctioned or to be replaced
	//Reserved = Rented to someone
	{UID: "Test0001", FirstName: "Zero", LastName: "tester", Password: "-", Email: "-", ReservedBt: "-"},
	{UID: "Test0002", FirstName: "One", LastName: "tester", Password: "-", Email: "-", ReservedBt: "-"},
	{UID: "Test0003", FirstName: "Two", LastName: "tester", Password: "-", Email: "-", ReservedBt: "-"},
}

func createUsers(u *gin.Context) {
	var newUser user
	if err := u.BindJSON(&newUser); err != nil {
		return
	}

	users = append(users, newUser)
	u.IndentedJSON(http.StatusCreated, newUser) //201
}

//---------------------------------------------------------------------------------------------------------------------------------------------------------------//
