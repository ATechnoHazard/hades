package entity

import (
	"github.com/dgrijalva/jwt-go"
	"reflect"
)

type Participant struct {
	Name               string `json:"name"`
	RegistrationNumber string `json:"registrationNumber" gorm:"primary_key"`
	Email              string `json:"email"`
	PhoneNumber        string `json:"phoneNumber"`
	Gender             string `json:"gender"`
}

type Attendee struct {
	Name               string `json:"name"`
	RegistrationNumber string `json:"registrationNumber" gorm:"primary_key"`
	Email              string `json:"email"`
	PhoneNumber        string `json:"phoneNumber"`
	Gender             string `json:"gender"`
	EventName          string `json:"eventName"`
}

type Event struct {
	ClubName              string      `json:"clubName"`
	Name                  string      `json:"name"`
	ToDate                string      `json:"toDate"`
	FromDate              string      `json:"fromDate"`
	ToTime                string      `json:"toTime"`
	FromTime              string      `json:"fromTime"`
	Budget                string      `json:"budget"`
	Description           string      `json:"description"`
	Category              string      `json:"category"`
	Venue                 string      `json:"venue"`
	Attendance            string      `json:"attendance"`
	ExpectedParticipants  string      `json:"expectedParticipants"`
	FacultyCoordinator    Participant `json:"facultyCoordinator"`
	StudentCoordinator    Participant `json:"studentCoordinator"`
	PROrequest            string      `json:"PROrequest"`
	CampusEngineerRequest string      `json:"campusEngineerRequest"`
	Duration              string      `json:"duration"`
	Status                string      `json:"status"`
	// MainSponsor           Participant `json:"mainSponsor"`
}

type Guest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	Gender         string `json:"gender"`
	Stake          string `json:"stake"`
	LocationOfStay string `json:"locationOfStay"`
}

type Attendance struct {
	EventName  string `json:"eventName"`
	Email      string `json:"email"`
	Day        int    `json:"day"`
	CouponName string `json:"couponName"`
}

type Coupon struct {
	Name string `json:"name"`
	Desc string `json:"description"`
	Day  int    `json:"day"`
}

type User struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	LinkedIn    string `json:"linkedIn"`
	Facebook    string `json:"facebook"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	DeviceToken string `json:"deviceToken"`
}

type Organization struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	CreatedAt   string `json:"createdAt"`
	Website     string `json:"website"`
}

type Token struct {
	Email        string `json:"email"`
	Role         string `json:"role"`
	Organization string `json:"organization"`
	jwt.StandardClaims
}

func (v Event) GetField(field string, value string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.FieldByName(value).String()
}
