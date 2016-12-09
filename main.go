package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Users table model ...
type Users struct {
	UserUUID string `gorm:"column:user_uuid;primary_key;type:varchar(36);not null"`
	Key      string `gorm:"type:varchar(100);not null"`
	Name     string `gorm:"type:varchar(100);not null"`
}

// Requests table model ...
type Requests struct {
	RequestUUID string    `gorm:"column:request_uuid;primary_key"`
	UserUUID    string    `gorm:"column:user_uuid"`
	Users       Users     `gorm:"ForeignKey:UserUUID"`
	Type        string    `gorm:"type:varchar(10)"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	Status      string    `gorm:"type:varchar(10);not null"`
	Params      string    `gorm:"type:json;"`
}

// Results table model ...
type Results struct {
	ResultID    uint      `gorm:"column:result_id;primary_key;not null"`
	RequestUUID string    `gorm:"column:request_uuid"`
	Requests    Requests  `gorm:"ForeignKey:RequestUUID"`
	id          string    `gorm:"column:user_uuid"`
	Users       Users     `gorm:"ForeignKey:UserUUID"`
	AddedAt     time.Time `gorm:"column:added_at"`
}

func main() {

	i := Impl{}
	i.InitDB()
	i.InitSchema()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/tsa/members_intersect", i.MembersIntersect),
		rest.Post("/tsa/get_result", i.GetResult),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

// ReqStuct request stuct
type ReqStuct struct {
	Auth      string   `json:"auth"`
	Groups    []string `json:"groups"`
	MemberMin int      `json:"member_min"`
	RequestID string   `json:"request_id"`
	Offset    int      `json:"offset"`
}

type OrderParams struct {
	Groups     []string `json:"groups"`
	MembersMin int      `json:"members_min"`
}

type OfsetResults struct {
	IDS      []string `json:"ids"`
	Finished bool     `json:"finished"`
}

type Impl struct {
	DB *gorm.DB
}

func (i *Impl) InitDB() {
	var err error
	i.DB, err = gorm.Open("postgres", "host=localhost user=postgres dbname=find_users sslmode=disable password=020407")
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&Users{}, &Requests{}, &Results{})
}

func (i *Impl) CheckAuth(w rest.ResponseWriter, r *rest.Request, req ReqStuct) Users {
	currentUser := Users{}
	if req.Auth == "" {
		rest.Error(w, "authorization key required", 400)
		return currentUser
	}

	if i.DB.Where("key = ?", req.Auth).First(&currentUser).Error != nil {
		rest.Error(w, "user not found", 400)
		return currentUser
	}
	return currentUser
}

func (i *Impl) MembersIntersect(w rest.ResponseWriter, r *rest.Request) {
	reqStuct := ReqStuct{}

	err := r.DecodeJsonPayload(&reqStuct)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser := i.CheckAuth(w, r, reqStuct)
	if err != nil || currentUser.UserUUID == "" {
		return
	}

	orderParams := OrderParams{}
	orderParams.Groups = reqStuct.Groups
	orderParams.MembersMin = reqStuct.MemberMin

	uuid, err := newUUID()
	if err != nil {
		rest.Error(w, "can't generate new uuid", 400)
		return
	}

	jsonParams, err := json.Marshal(&orderParams)
	if err != nil {
		rest.Error(w, "error marshal OrderParams", 400)
		return
	}

	requests := Requests{uuid, currentUser.UserUUID, currentUser, "пересечение сообществ", time.Now(), "PROCESSING", string(jsonParams)}

	if err := i.DB.Save(&requests).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&requests.RequestUUID)
}

func (i *Impl) GetResult(w rest.ResponseWriter, r *rest.Request) {
	reqStuct := ReqStuct{}

	err := r.DecodeJsonPayload(&reqStuct)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser := i.CheckAuth(w, r, reqStuct)
	if err != nil || currentUser.UserUUID == "" {
		return
	}
	result := Results{}
	if i.DB.Select("result_id").Where("request_uuid = ?", reqStuct.RequestID).Offset(reqStuct.Offset).Find(&result).Error != nil {
		rest.Error(w, "user not found", 400)
		return
	}
	fmt.Println(result)
}
