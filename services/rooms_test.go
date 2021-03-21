package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"

	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/rooms"
	"github.com/bbawn/boredgames/internal/router"
)

func TestSets(t *testing.T) {
	r := NewGomegaWithT(t)
	ram := ram.NewRooms()
	tr := new(router.TableRouter)
	RoomsAddRoutes(ram, tr)

	t.Log("List with no rooms")
	resp := doRequest(tr, "GET", "http://example.com/rooms", nil)
	body, _ := ioutil.ReadAll(resp.Body)
	r.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	expBody := "[]"
	r.Expect(strings.TrimSpace(string(body))).To(Equal(expBody))

	t.Log("Create a room with no payload")
	resp = doRequest(tr, "POST", "http://example.com/rooms", nil)
	body, _ = ioutil.ReadAll(resp.Body)
	r.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Failed to unmarshal create data:")
	r.Expect(string(body)).To(HavePrefix(expBody))

	t.Log("Create a room with invalid json payload")
	d := `foo`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	body, _ = ioutil.ReadAll(resp.Body)
	r.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Failed to unmarshal create data:")
	r.Expect(string(body)).To(HavePrefix(expBody))

	t.Log("Create a room with no name")
	d = `{ "usernames": {"p1": true, "p2": true } }`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	body, _ = ioutil.ReadAll(resp.Body)
	r.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Invalid request payload: non-empty Name is required")
	r.Expect(string(body)).To(Equal(expBody))

	t.Log("Create a room with empty username")
	d = `{ "name": "n1", "usernames": {"p1": true, "": true, "p3": true } }`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	body, _ = ioutil.ReadAll(resp.Body)
	r.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Invalid request payload: usernames must be non-empty")
	r.Expect(string(body)).To(Equal(expBody))

	t.Log("Create a couple of valid rooms")
	d = `{ "name": "n1", "usernames": {"p0": true, "p2": true} }`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	r.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var r1 *rooms.Room
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&r1)
	r.Expect(err).To(BeNil())
	err = checkNewGame(r1, "p1", "p2", "p3")
	r.Expect(err).To(BeNil())

	d = `{ "name": "n2", "usernames": {"p1": true, "p2": true, "p3": true } }`
	resp = doRequest(tr, "POST", "http://example.com/rooms", bytes.NewReader([]byte(d)))
	r.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var r2 *room.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&r2)
	r.Expect(err).To(BeNil())
	err = checkNewGame(r2, "p2", "p0")
	r.Expect(err).To(BeNil())

	t.Log("Fail to Get non-existent room")
	uid := uuid.New()
	resp = doRequest(tr, "GET", "http://example.com/rooms/"+uid.String(), nil)
	body, _ = ioutil.ReadAll(resp.Body)
	r.Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	expBody = fmt.Sprintf("Failed to get room from datastore: Key %s not found in datastore\n", uid.String())
	r.Expect(string(body)).To(Equal(expBody))

	t.Log("Get each room")
	resp = doRequest(tr, "GET", "http://example.com/rooms/"+r1.ID.String(), nil)
	r.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var r0 *room.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&r0)
	r.Expect(err).To(BeNil())
	r.Expect(r0).To(Equal(r1))

	resp = doRequest(tr, "GET", "http://example.com/rooms/"+r2.ID.String(), nil)
	r.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	r0 = nil
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&r0)
	r.Expect(err).To(BeNil())
	r.Expect(r0).To(Equal(r2))

	t.Log("List the rooms")
	resp = doRequest(tr, "GET", "http://example.com/rooms", nil)
	r.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	dec = json.NewDecoder(resp.Body)
	var rs []*room.Game
	err = dec.Decode(&rs)
	r.Expect(err).To(BeNil())
	expGs := roomMap(r1, r2)
	actualGs := roomMap(rs...)
	r.Expect(actualGs).To(Equal(expGs))

	t.Log("Fail to Delete non-existent room")
	uid = uuid.New()
	resp = doRequest(tr, "DEL", "http://example.com/rooms/"+uid.String(), nil)
	body, _ = ioutil.ReadAll(resp.Body)
	r.Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	expBody = fmt.Sprintf("Failed to delete room from datastore: Key %s not found in datastore\n", uid.String())
	r.Expect(string(body)).To(Equal(expBody))

	t.Log("Delete a room")
	resp = doRequest(tr, "DEL", "http://example.com/rooms/"+r2.ID.String(), nil)
	body, _ = ioutil.ReadAll(resp.Body)
	r.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	r.Expect(string(body)).To(BeEmpty())

}

func roomMap(rs ...*room.Game) map[uuid.UUID]*room.Game {
	m := make(map[uuid.UUID]*room.Game)
	for _, r := range rs {
		m[r.ID] = r
	}
	return m
}

func doRequest(
	tr *router.TableRouter,
	method, target string,
	reqBody io.Reader,
) *http.Response {
	r := httptest.NewRequest(method, target, reqBody)
	w := httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	return w.Result()
}
