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
	"github.com/bbawn/boredgames/internal/games/set"
	"github.com/bbawn/boredgames/internal/router"
)

func TestSets(t *testing.T) {
	g := NewGomegaWithT(t)
	ram := ram.NewSets()
	tr := new(router.TableRouter)
	SetsAddRoutes(ram, tr)

	t.Log("List with no games")
	resp := doRequest(tr, "GET", "http://example.com/sets", nil)
	body, _ := ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	expBody := "[]"
	g.Expect(strings.TrimSpace(string(body))).To(Equal(expBody))

	t.Log("Create a game with no payload")
	resp = doRequest(tr, "POST", "http://example.com/sets", nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Failed to unmarshal create data:")
	g.Expect(string(body)).To(HavePrefix(expBody))

	t.Log("Create a game with invalid json payload")
	d := `foo`
	resp = doRequest(tr, "POST", "http://example.com/sets", bytes.NewReader([]byte(d)))
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Failed to unmarshal create data:")
	g.Expect(string(body)).To(HavePrefix(expBody))

	t.Log("Create a game with empty username")
	d = `{ "usernames": [ "p1", "", "p3" ] }`
	resp = doRequest(tr, "POST", "http://example.com/sets", bytes.NewReader([]byte(d)))
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Failed to create new game: Invalid value: empty for arg: username\n")
	g.Expect(string(body)).To(Equal(expBody))

	t.Log("Create a couple of games")
	d = `{ "usernames": [ "p1", "p2", "p3" ] }`
	resp = doRequest(tr, "POST", "http://example.com/sets", bytes.NewReader([]byte(d)))
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var g1 *set.Game
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&g1)
	g.Expect(err).To(BeNil())
	err = checkNewGame(g1, "p1", "p2", "p3")
	g.Expect(err).To(BeNil())

	d = `{ "usernames": [ "p2", "p0" ] }`
	resp = doRequest(tr, "POST", "http://example.com/sets", bytes.NewReader([]byte(d)))
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var g2 *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g2)
	g.Expect(err).To(BeNil())
	err = checkNewGame(g2, "p2", "p0")
	g.Expect(err).To(BeNil())

	t.Log("Fail to Get non-existent game")
	uid := uuid.New()
	resp = doRequest(tr, "GET", "http://example.com/sets/"+uid.String(), nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	expBody = fmt.Sprintf("Failed to get game from datastore: Key %s not found in datastore\n", uid.String())
	g.Expect(string(body)).To(Equal(expBody))

	t.Log("Get each game")
	resp = doRequest(tr, "GET", "http://example.com/sets/"+g1.ID.String(), nil)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var g0 *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g0)
	g.Expect(err).To(BeNil())
	g.Expect(g0).To(Equal(g1))

	resp = doRequest(tr, "GET", "http://example.com/sets/"+g2.ID.String(), nil)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	g0 = nil
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g0)
	g.Expect(err).To(BeNil())
	g.Expect(g0).To(Equal(g2))

	t.Log("List the games")
	resp = doRequest(tr, "GET", "http://example.com/sets", nil)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	dec = json.NewDecoder(resp.Body)
	var gs []*set.Game
	err = dec.Decode(&gs)
	g.Expect(err).To(BeNil())
	expGs := gameMap(g1, g2)
	actualGs := gameMap(gs...)
	g.Expect(actualGs).To(Equal(expGs))

	t.Log("Fail to Delete non-existent game")
	uid = uuid.New()
	resp = doRequest(tr, "DEL", "http://example.com/sets/"+uid.String(), nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	expBody = fmt.Sprintf("Failed to delete game from datastore: Key %s not found in datastore\n", uid.String())
	g.Expect(string(body)).To(Equal(expBody))

	t.Log("Delete a game")
	resp = doRequest(tr, "DEL", "http://example.com/sets/"+g2.ID.String(), nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	g.Expect(string(body)).To(BeEmpty())

	t.Log("Next move in invalid state")
	resp = doRequest(tr, "POST", "http://example.com/sets/"+g1.ID.String()+"/next", nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusConflict))
	expBody = fmt.Sprintf("Failed to advance game to next round: Invalid method: NextRound detail: round not yet claimed\n")
	g.Expect(string(body)).To(Equal(expBody))

	t.Log("Claim a set with no payload")
	resp = doRequest(tr, "POST", "http://example.com/sets/"+g1.ID.String()+"/claim", nil)
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Failed to claim set in game: Invalid value:  for arg: username\n")
	g.Expect(string(body)).To(Equal(expBody))

	t.Log("Claim a set with invalid json payload")
	payload := []byte(`foo`)
	resp = doRequest(tr, "POST", "http://example.com/sets/"+g1.ID.String()+"/claim", bytes.NewReader(payload))
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Failed to claim set in game: Invalid value:  for arg: username\n")
	g.Expect(string(body)).To(Equal(expBody))

	s1 := g1.FindExpandSet()
	t.Log("Claim a set with invalid username in payload")
	payload = claimPayload("nonplayer", *s1)
	resp = doRequest(tr, "POST", "http://example.com/sets/"+g1.ID.String()+"/claim", bytes.NewReader(payload))
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	expBody = fmt.Sprintf("Failed to claim set in game: Invalid value: nonplayer for arg: username\n")
	g.Expect(string(body)).To(Equal(expBody))

	t.Log("Claim a set with non-set in payload (penalty)")
	nonset := g1.Board.FindSet(false)
	payload = claimPayload("p1", *nonset)
	resp = doRequest(tr, "POST", "http://example.com/sets/"+g1.ID.String()+"/claim", bytes.NewReader(payload))
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var g1ClaimFail *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g1ClaimFail)
	g.Expect(err).To(BeNil())
	g.Expect(g1ClaimFail.ID).To(Equal(g1.ID))
	g.Expect(g1ClaimFail.GetState()).To(Equal(set.Playing))

	t.Log("Claim a set")
	payload = claimPayload("p1", *s1)
	resp = doRequest(tr, "POST", "http://example.com/sets/"+g1.ID.String()+"/claim", bytes.NewReader(payload))
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var g1Claimed *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g1Claimed)
	g.Expect(err).To(BeNil())
	err = checkNextGame(g1, g1Claimed)
	g.Expect(err).To(BeNil())
	g.Expect(len(g1Claimed.Players["p1"].Sets)).To(Equal(1))
	// TODO: DeepEqual sets

	t.Log("Claim a set in invalid game state")
	resp = doRequest(tr, "POST", "http://example.com/sets/"+g1.ID.String()+"/claim", bytes.NewReader(payload))
	body, _ = ioutil.ReadAll(resp.Body)
	g.Expect(resp.StatusCode).To(Equal(http.StatusConflict))
	expBody = fmt.Sprintf("Failed to claim set in game: Invalid method: ClaimSet detail: round already claimed by p1\n")
	g.Expect(string(body)).To(Equal(expBody))

	t.Log("Valid Next round request")
	resp = doRequest(tr, "POST", "http://example.com/sets/"+g1.ID.String()+"/next", nil)
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var g1Next *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g1Next)
	g.Expect(err).To(BeNil())
	err = checkNextGame(g1Claimed, g1Next)
	g.Expect(err).To(BeNil())
}

// checkNewGame validates that the given game is in a valid initial state
func checkNewGame(g0 *set.Game, usernames ...string) error {
	if g0.ID.URN() == "" {
		return fmt.Errorf("Invalid ID: %s", g0.ID)
	}
	if len(g0.Players) != len(usernames) {
		return fmt.Errorf("Expected %d players, got %d", len(g0.Players), len(usernames))
	}
	for _, u := range usernames {
		var (
			p  *set.Player
			ok bool
		)
		if p, ok = g0.Players[u]; !ok {
			return fmt.Errorf("Player with username %s not found", u)
		}
		if p.Username != u {
			return fmt.Errorf("Expected player Username %s, got %s", u, p.Username)
		}
		if len(p.Sets) != 0 {
			return fmt.Errorf("Expected empty Sets for player %s", u)
		}
	}
	return nil
}

// checkNextGame validates that g1 is a valid next state of g0
func checkNextGame(g0, g1 *set.Game) error {
	if g0.ID != g1.ID {
		return fmt.Errorf("Expected g0 ID %s to equal g1 ID %s", g0.ID, g1.ID)
	}
	if g0.GetState() == set.Playing {
		if g1.GetState() != set.SetClaimed {
			return fmt.Errorf("Expected g1 State to be SetClaimed")
		}
	} else {
		if g1.GetState() != set.Playing {
			return fmt.Errorf("Expected g1 State to be Playing")
		}
	}
	// TODO? we could check invariants on players, scores, deck size, etc...
	return nil
}

func gameMap(gs ...*set.Game) map[uuid.UUID]*set.Game {
	m := make(map[uuid.UUID]*set.Game)
	for _, g := range gs {
		m[g.ID] = g
	}
	return m
}

func claimPayload(username string, cs set.CardTriple) []byte {
	cd := claimData{
		Username: username,
		Cards:    cs,
	}
	payload, err := json.Marshal(&cd)
	if err != nil {
		panic(fmt.Sprintf("Unexpected Marshal err: %s", err))
	}
	return payload
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
