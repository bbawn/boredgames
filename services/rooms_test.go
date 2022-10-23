package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"

	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/rooms"
	"github.com/bbawn/boredgames/internal/router"
)

func TestRooms(t *testing.T) {
	g := NewGomegaWithT(t)
	ram := ram.NewRooms()
	tr := new(router.TableRouter)
	RoomsAddRoutes(ram, tr)

	t.Log("List with no rooms")
	resp := doRequest(tr, "GET", "http://example.com/rooms", nil)
	body, _ := ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	expBody := "[]"
	g.Expect(strings.TrimSpace(string(body))).To(Equal(expBody))

	t.Log("Create a room with no payload")
	resp = doRequest(tr, "POST", "http://example.com/rooms", nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	g.Expect(string(body)).To(HavePrefix("Failed to unmarshal create data:"))

	t.Log("Create a room with invalid json payload")
	d := `foo`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	g.Expect(string(body)).To(HavePrefix("Failed to unmarshal create data:"))

	t.Log("Create a room with no name")
	d = `{ "usernames": {"p1": true, "p2": true } }`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	g.Expect(string(body)).To(Equal("Invalid request payload err: non-empty Name is required\n"))

	t.Log("Create a room with empty username")
	d = `{ "name": "n1", "usernames": {"p1": true, "": true, "p3": true } }`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	g.Expect(string(body)).To(Equal("Invalid request payload err: usernames must be non-empty\n"))

	t.Log("Create a couple of valid rooms")
	d = `{ "name": "n1", "usernames": {"p0": true, "p2": true} }`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	body, _ = ioutil.ReadAll(resp.Body)
	var expRoom, r1 *rooms.Room
	err := json.Unmarshal(body, &r1)
	g.Expect(err).To(BeNil())
	err = json.Unmarshal([]byte(d), &expRoom)
	g.Expect(err).To(BeNil())
	if !reflect.DeepEqual(expRoom, r1) {
		t.Errorf("Post returned %#v, expected %#v", r1, expRoom)
	}

	d = `{ "name": "n2", "usernames": {"p1": true, "p2": true, "p3": true } }`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	body, _ = ioutil.ReadAll(resp.Body)
	var r2 *rooms.Room
	err = json.Unmarshal(body, &r2)
	g.Expect(err).To(BeNil())
	expRoom = nil
	err = json.Unmarshal([]byte(d), &expRoom)
	g.Expect(err).To(BeNil())
	if !reflect.DeepEqual(expRoom, r2) {
		t.Errorf("Post returned %#v, expected %#v", r2, expRoom)
	}

	t.Log("Fail to Get non-existent room")
	uid := uuid.New()
	resp = doRequest(tr, "GET", "http://example.com/rooms/"+uid.String(), nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	expBody = fmt.Sprintf("Failed to get room from datastore: Key %s not found in datastore\n", uid.String())
	g.Expect(string(body)).To(Equal(expBody))

	t.Log("Get each room")
	resp = doRequest(tr, "GET", "http://example.com/rooms/"+r1.Name, nil)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var r0 *rooms.Room
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&r0)
	g.Expect(err).To(BeNil())
	g.Expect(r0).To(Equal(r1))

	resp = doRequest(tr, "GET", "http://example.com/rooms/"+r2.Name, nil)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	r0 = nil
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&r0)
	g.Expect(err).To(BeNil())
	g.Expect(r0).To(Equal(r2))

	t.Log("List the rooms")
	resp = doRequest(tr, "GET", "http://example.com/rooms", nil)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	dec = json.NewDecoder(resp.Body)
	var rs []*rooms.Room
	err = dec.Decode(&rs)
	g.Expect(err).To(BeNil())
	expGs := roomMap(r1, r2)
	actualGs := roomMap(rs...)
	g.Expect(actualGs).To(Equal(expGs))

	t.Log("Fail to Delete non-existent room")
	uid = uuid.New()
	resp = doRequest(tr, "DEL", "http://example.com/rooms/"+uid.String(), nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	expBody = fmt.Sprintf("Failed to delete room from datastore: Key %s not found in datastore\n", uid.String())
	g.Expect(string(body)).To(Equal(expBody))

	t.Log("Delete a room")
	resp = doRequest(tr, "DEL", "http://example.com/rooms/"+r2.Name, nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	g.Expect(string(body)).To(BeEmpty())
}

func roomMap(rs ...*rooms.Room) map[string]*rooms.Room {
	m := make(map[string]*rooms.Room)
	for _, r := range rs {
		m[r.Name] = r
	}
	return m
}
