package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"encoding/xml"
	"encoding/json"
	_"fmt"
	"strconv"
)

type xmlRow struct {
	Id int `xml:"id"`
	Guid string `xml:"guid"`
	IsActive bool `xml:"isActive"`
	Balance string `xml:"balance"`
	Picture string `xml:"picture"`
	Age int `xml:"age"`
	EyeColor string `xml:"eyeColor"`
	FirstName string `xml:"first_name"`
	LastName string `xml:"last_name"`
	Gender string `xml:"gender"`
	Company string `xml:"company"`
	Email string `xml:"email"`
	Phone string `xml:"phone"`
	Address string `xml:"address"`
	About string `xml:"about"`
}

type xmlStructure struct {
	Version string `xml:"version"`
	Row []xmlRow `xml:"row"`

}

const pageSize = 25


func SearchServerSuccess(w http.ResponseWriter, r *http.Request)  {
	dataFile, err := ioutil.ReadFile("dataset.xml")
	checkError(err)

	usersXml := &xmlStructure{}
	xml.Unmarshal(dataFile, &usersXml)

	var users []User

	for _, user := range usersXml.Row {
		users = append(users, User{
			Id: user.Id,
			Name: user.FirstName,
			Age: user.Age,
			About: user.About,
			Gender: user.Gender,
		})
	}

	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	var startRow int
	if offset > 0 {
		startRow = offset * pageSize
	}

	endRow := startRow + limit
	users = users[ startRow: endRow ]

	jsonResponse, err := json.Marshal(users)
	checkError(err)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}


func TestErrorResponse(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServerSuccess))

	searchClient := &SearchClient{
		URL: ts.URL,
	}

	searchRequest := SearchRequest{
		Limit: 5,
		Offset: 0,
	}

	_, err := searchClient.FindUsers(searchRequest)

	if err != nil {
		t.Error("Dosn't work success request")
	}

	searchRequest.Limit = -1

	_, err = searchClient.FindUsers(searchRequest)
	if err.Error() != "limit must be > 0" {
		t.Error("limit must be > 0")
	}

	searchRequest.Limit = 1
	searchRequest.Offset = -1
	_, err = searchClient.FindUsers(searchRequest)
	if err.Error() != "offset must be > 0" {
		t.Error("offset must be > 0")
	}

	ts.Close()
}



func TestLimitFailed(t *testing.T)  {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerSuccess))

	searchClient := &SearchClient{
		URL: ts.URL,
	}

	_, err := searchClient.FindUsers(SearchRequest{
		Offset: 1,
		Limit: 7,
	})

	if err != nil {
		t.Error(err.Error())
	}

}

func checkError(err error)  {
	if err != nil {
		panic(err)
	}
}