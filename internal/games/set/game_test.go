package set

import (
	"fmt"
	"testing"
)

func TestGame(t *testing.T) {
	usernames := []string{"Joe", "Natasha", "Maria", "Frank"}
	g, err := NewGame(usernames)
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
		invStateErr *InvalidStateError
		invArgErr   *InvalidArgError
		set         []*Card
		ok          bool
	)
	// Can only go to next round if there is a claimed set
	fmt.Printf("Calling NextRound\n")
	err = g.NextRound()
	if invStateErr, ok = err.(*InvalidStateError); !ok {
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
	set = g.testFindSet()
	err = g.ClaimSet("Jane", set[0], set[1], set[2])
	if invArgErr, ok = err.(*InvalidArgError); !ok {
		t.Errorf("expected InvalidArgError, got: %v", err)
	}
	if invArgErr.Arg != "username" {
		t.Errorf("expected Arg: username, got: %v", invArgErr.Arg)
	}
	if invArgErr.Value != "Jane" {
		t.Errorf("expected Value: Jane, got: %v", invArgErr.Value)
	}

	// Claim with non-set
	/*
		err = g.ClaimSet("Joe", s[0], s[1], s[2])
		if invStateErr, ok := err.(*InvalidArgError); !ok {
			t.Errorf("expected InvalidArgError, got: %v", err)
		}
		if invArgErr.Arg != "username" {
			t.Errorf("expected Arg: username, got: %v", invArgErr.Arg)
		}
		if invArgErr.Value != "Jane" {
			t.Errorf("expected Value: Jane, got: %v", invArgErr.Value)
		}

		// while s :=
	*/
}

// Return a set on the board, expanding until one is found
func (g *Game) testFindSet() []*Card {
	fmt.Printf("FindSet\n")
	for true {
		s := g.Board.FindSet()
		fmt.Printf("s: %#v\n", s)
		if len(s) == SetLen {
			return s
		}
		if !g.ExpandBoard() {
			return []*Card{}
		}
	}
	panic("unreachable")
}
