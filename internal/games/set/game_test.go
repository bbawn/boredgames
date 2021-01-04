package set

import (
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

	// while s :=
}
