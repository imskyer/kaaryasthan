package item_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	. "github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/route"
)

func TestItemCreateHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM items")
	ts := httptest.NewServer(route.RT)
	defer ts.Close()
	n := []byte(`{
  "data": {
    "type": "items",
    "attributes": {
      "title": "Some Title",
      "description": "Some description"
    }
  }
}`)
	reqPayload := make(map[string]Data)
	decoder1 := json.NewDecoder(bytes.NewReader(n))
	err := decoder1.Decode(&reqPayload)
	if err != nil {
		log.Println("Unable to decode body: ", err)
		return
	}

	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/items", bytes.NewReader(n))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respPayload := make(map[string]Data)
	decoder2 := json.NewDecoder(resp.Body)
	err = decoder2.Decode(&respPayload)
	if err != nil {
		t.Fatal("Unable to decode body: ", err)
	}
	respData := respPayload["data"]
	reqData := reqPayload["data"]
	reqData.ID = respData.ID
	id, err := strconv.Atoi(reqData.ID)
	if err != nil {
		t.Fatal("Wrong ID", err)
	}
	if id <= 0 {
		t.Errorf("ID is not 1 or above: %#v", id)
	}

	if !reflect.DeepEqual(reqData, respData) {
		t.Errorf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqData, respData)
	}
}
