package set

import (
	"fmt"
	"math/rand"
	"testing"
)

const (
	nTestGames = 256
)

func getUsernames() []string {
	return []string{"Joe", "Natasha", "Maria", "Frank"}
}

func TestCardString(t *testing.T) {
	c1 := &Card{Filled, Triangle, Red, 1}
	s := c1.String()
	if s != "FTR1" {
		t.Errorf("expected: c1.String() to be %s, got %s", "FTR1", s)
	}
}

func TestIsSet(t *testing.T) {
	var (
		c1, c2, c3 *Card
	)

	// duplicate card, not even a non-set
	c1 = &Card{Filled, Triangle, Red, 1}
	c2 = &Card{Filled, Triangle, Red, 1}
	c3 = &Card{Filled, Triangle, Purple, 1}
	if IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v not to be a set", c1, c2, c3)
	}

	// not a set
	c1 = &Card{Filled, Triangle, Red, 1}
	c2 = &Card{Filled, Squiggle, Red, 1}
	c3 = &Card{Stripe, Squiggle, Purple, 1}
	if IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v not to be a set", c1, c2, c3)
	}

	// Shading same, shape different, color different, count different
	c1 = &Card{Filled, Triangle, Purple, 1}
	c2 = &Card{Filled, Squiggle, Red, 2}
	c3 = &Card{Filled, Oval, Green, 3}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
	// Shading different, shape same, color different, count different
	c1 = &Card{Filled, Squiggle, Purple, 1}
	c2 = &Card{Stripe, Squiggle, Red, 2}
	c3 = &Card{Outline, Squiggle, Green, 3}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
	// Shading different, shape different, color same, count different
	c1 = &Card{Filled, Triangle, Green, 1}
	c2 = &Card{Stripe, Squiggle, Green, 2}
	c3 = &Card{Outline, Oval, Green, 3}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
	// Shading different, shape different, color different, count same
	c1 = &Card{Filled, Triangle, Purple, 1}
	c2 = &Card{Stripe, Squiggle, Red, 1}
	c3 = &Card{Outline, Oval, Green, 1}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}

	// Shading same, shape same, color different, count different
	c1 = &Card{Filled, Oval, Purple, 1}
	c2 = &Card{Filled, Oval, Red, 2}
	c3 = &Card{Filled, Oval, Green, 3}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
	// Shading different, shape same, color same, count different
	c1 = &Card{Filled, Oval, Red, 1}
	c2 = &Card{Stripe, Oval, Red, 2}
	c3 = &Card{Outline, Oval, Red, 3}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
	// Shading different, shape different, color same, count same
	c1 = &Card{Filled, Triangle, Purple, 2}
	c2 = &Card{Stripe, Squiggle, Purple, 2}
	c3 = &Card{Outline, Oval, Purple, 2}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}

	// Shading different, shape same, color same, count same
	c1 = &Card{Filled, Oval, Red, 3}
	c2 = &Card{Stripe, Oval, Red, 3}
	c3 = &Card{Outline, Oval, Red, 3}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}

	// Shading same, shape different, color same, count same
	c1 = &Card{Stripe, Oval, Red, 3}
	c2 = &Card{Stripe, Triangle, Red, 3}
	c3 = &Card{Stripe, Squiggle, Red, 3}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}

	// Shading same, shape same, color different, count same
	c1 = &Card{Filled, Oval, Red, 3}
	c2 = &Card{Filled, Oval, Green, 3}
	c3 = &Card{Filled, Oval, Purple, 3}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}

	// Shading same, shape same, color same, count different
	c1 = &Card{Filled, Oval, Red, 1}
	c2 = &Card{Filled, Oval, Red, 2}
	c3 = &Card{Filled, Oval, Red, 3}
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
}

func TestInvalidGames(t *testing.T) {
	usernames := getUsernames()
	g, err := NewGame(usernames...)
	if err != nil {
		t.Errorf("expected err: %v to be nil", err)
	}
	if len(usernames) != len(g.Players) {
		t.Errorf("expected %d usernames, got %d", len(usernames), len(g.Players))
	}
	for _, u := range usernames {
		if _, present := g.Players[u]; !present {
			t.Errorf("expected username: %v missing", u)
		}
	}
	if len(g.Deck) != FullDeckLen-InitBoardLen {
		t.Errorf("expected len(p.Deck): %v got: %v", FullDeckLen-InitBoardLen, len(g.Deck))
	}

	// Board and remaining Deck should comprise full deck
	allCards := make(map[CardBase3]*Card)
	for _, c := range g.Board {
		allCards[CardToCardBase3(c)] = c
	}

	for _, c := range g.Deck {
		allCards[CardToCardBase3(c)] = c
	}
	if len(allCards) != FullDeckLen {
		t.Errorf("expected len(allCards): %v got: %v", FullDeckLen, len(allCards))
	}

	var (
		invStateErr InvalidStateError
		invArgErr   InvalidArgError
		s           []*Card
		ok          bool
	)
	// Can only go to next round if there is a claimed set
	err = g.NextRound()
	if invStateErr, ok = err.(InvalidStateError); !ok {
		t.Errorf("expected InvalidStateError, got: %v", err)
	}
	if invStateErr.Method != "NextRound" {
		t.Errorf("expected Method ClaimSet, got: %v", invStateErr.Method)
	}
	if invStateErr.Details != "round not yet claimed" {
		t.Errorf("expected Details: round not yet claimed, got: %v",
			invStateErr.Details)
	}

	// Claim with invalid username
	s = g.FindExpandSet()
	err = g.ClaimSet("Jane", s[0], s[1], s[2])
	if invArgErr, ok = err.(InvalidArgError); !ok {
		t.Errorf("expected InvalidArgError, got: %v", err)
	}
	if invArgErr.Arg != "username" {
		t.Errorf("expected Arg: username, got: %v", invArgErr.Arg)
	}
	if invArgErr.Value != "Jane" {
		t.Errorf("expected Value: Jane, got: %v", invArgErr.Value)
	}

	// Claim with non-set
	s = g.Board.FindSet(false)
	err = g.ClaimSet("Joe", s[0], s[1], s[2])
	if invArgErr, ok = err.(InvalidArgError); !ok {
		t.Errorf("expected InvalidArgError, got: %v", err)
	}
	if invArgErr.Arg != "set" {
		t.Errorf("expected Arg: set, got: %v", invArgErr.Arg)
	}
	if invArgErr.Value != fmt.Sprintf("%v %v %v", s[0], s[1], s[2]) {
		t.Errorf("expected Value: Jane, got: %v", invArgErr.Value)
	}
}

func TestValidGames(t *testing.T) {
	for i := 0; i < nTestGames; i++ {
		t.Log("Game:", i)
		usernames := getUsernames()
		g, err := NewGame(usernames...)
		if err != nil {
			t.Errorf("expected NewGame() to succeed, got: %v", err)
		}
		for s := g.FindExpandSet(); len(s) > 0; s = g.FindExpandSet() {
			t.Log("Board:", g.Board)
			t.Log("Deck len:", len(g.Deck))
			u := usernames[rand.Intn(len(usernames))]
			t.Logf("Username: %s found set: %v %v %v", u, s[0], s[1], s[2])

			uOldScore := len(g.Players[u].Sets)
			err = g.ClaimSet(u, s[0], s[1], s[2])
			if err != nil {
				t.Errorf("expected ClaimSet to succeed, got: %v", err)
			}
			if g.ClaimedUsername != u {
				t.Errorf("expected ClaimUsername to be %s, got: %s", u, g.ClaimedUsername)
			}
			uNewScore := len(g.Players[u].Sets)
			if uNewScore != uOldScore+1 {
				t.Errorf("expected uNewScore to be %d, got: %d", uOldScore+1, uNewScore)
			}
			err = g.NextRound()
			if err != nil {
				t.Errorf("expected ClaimSet to succeed, got: %v", err)
			}
			if g.ClaimedUsername != "" {
				t.Errorf("expected ClaimUsername to be empty, got: %s", g.ClaimedUsername)
			}
		}
		if len(g.Deck) != 0 {
			t.Errorf("expected Deck to be empty, got: %d", len(g.Deck))
		}
	}
}
