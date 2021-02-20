package set

import (
	. "github.com/onsi/gomega"
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
	g := NewGomegaWithT(t)
	c1 := &Card{Red, 1, Filled, Diamond}
	s := c1.String()
	g.Expect(s).To(Equal("R1FD"))
}

func TestIsSet(t *testing.T) {
	var (
		c1, c2, c3 Card
	)
	g := NewGomegaWithT(t)

	// duplicate card, not even a non-set
	c1 = Card{Red, 1, Filled, Diamond}
	c2 = Card{Red, 1, Filled, Diamond}
	c3 = Card{Purple, 1, Filled, Diamond}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).NotTo(BeTrue())

	// not a set
	c1 = Card{Red, 1, Filled, Diamond}
	c2 = Card{Red, 1, Filled, Squiggle}
	c3 = Card{Purple, 1, Stripe, Squiggle}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).NotTo(BeTrue())

	// Shading same, shape different, color different, count different
	c1 = Card{Purple, 1, Filled, Diamond}
	c2 = Card{Red, 2, Filled, Squiggle}
	c3 = Card{Green, 3, Filled, Oval}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())
	// Shading different, shape same, color different, count different
	c1 = Card{Purple, 1, Filled, Squiggle}
	c2 = Card{Red, 2, Stripe, Squiggle}
	c3 = Card{Green, 3, Outline, Squiggle}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())
	if !IsSet(CardTriple{c1, c2, c3}) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
	// Shading different, shape different, color same, count different
	c1 = Card{Green, 1, Filled, Diamond}
	c2 = Card{Green, 2, Stripe, Squiggle}
	c3 = Card{Green, 3, Outline, Oval}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())
	if !IsSet(CardTriple{c1, c2, c3}) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
	// Shading different, shape different, color different, count same
	c1 = Card{Purple, 1, Filled, Diamond}
	c2 = Card{Red, 1, Stripe, Squiggle}
	c3 = Card{Green, 1, Outline, Oval}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())
	if !IsSet(CardTriple{c1, c2, c3}) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}

	// Shading same, shape same, color different, count different
	c1 = Card{Purple, 1, Filled, Oval}
	c2 = Card{Red, 2, Filled, Oval}
	c3 = Card{Green, 3, Filled, Oval}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())

	// Shading different, shape same, color same, count different
	c1 = Card{Red, 1, Filled, Oval}
	c2 = Card{Red, 2, Stripe, Oval}
	c3 = Card{Red, 3, Outline, Oval}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())

	// Shading different, shape different, color same, count same
	c1 = Card{Red, 2, Filled, Diamond}
	c2 = Card{Red, 2, Stripe, Squiggle}
	c3 = Card{Red, 2, Outline, Oval}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())

	// Shading different, shape same, color same, count same
	c1 = Card{Red, 2, Filled, Oval}
	c2 = Card{Red, 2, Stripe, Oval}
	c3 = Card{Red, 2, Outline, Oval}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())

	// Shading same, shape different, color same, count same
	c1 = Card{Red, 2, Stripe, Diamond}
	c2 = Card{Red, 2, Stripe, Oval}
	c3 = Card{Red, 2, Stripe, Squiggle}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())

	// Shading same, shape same, color different, count same
	c1 = Card{Red, 2, Stripe, Oval}
	c2 = Card{Green, 2, Stripe, Oval}
	c3 = Card{Purple, 2, Stripe, Oval}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())

	// Shading same, shape same, color same, count different
	c1 = Card{Red, 1, Stripe, Oval}
	c2 = Card{Red, 2, Stripe, Oval}
	c3 = Card{Red, 3, Stripe, Oval}
	g.Expect(IsSet(CardTriple{c1, c2, c3})).To(BeTrue())
}

func TestGameSequence(t *testing.T) {
	g := NewGomegaWithT(t)
	usernames := getUsernames()
	game, err := NewGame(usernames...)
	g.Expect(err).To(Succeed())
	g.Expect(len(game.Players)).To(Equal(len(usernames)))
	g.Expect(game.GetState()).To(Equal(Playing))
	for _, u := range usernames {
		g.Expect(game.Players).To(HaveKey(u))
	}
	g.Expect(len(game.Deck)).To(Equal(FullDeckLen - InitBoardLen))
	g.Expect(len(game.Board)).To(Equal(InitBoardLen))

	// Board and remaining Deck should comprise full deck
	allCards := make(map[CardBase3]*Card)
	for _, c := range game.Board {
		allCards[CardToCardBase3(c)] = c
	}

	for _, c := range game.Deck {
		allCards[CardToCardBase3(c)] = c
	}
	g.Expect(len(allCards)).To(Equal(FullDeckLen))

	// Can only go to next round if there is a claimed set
	err = game.NextRound()
	g.Expect(err).To(MatchError(InvalidStateError{"NextRound", "round not yet claimed"}))
	g.Expect(game.GetState()).To(Equal(Playing))

	// Expand is valid in Playing state
	err = game.Expand()
	g.Expect(err).To(Succeed())
	g.Expect(game.GetState()).To(Equal(Playing))
	g.Expect(len(game.Board)).To(Equal(InitBoardLen + SetLen))

	// Claim with invalid username
	s := game.FindExpandSet()
	err = game.ClaimSet("Jane", *s)
	g.Expect(err).To(MatchError(InvalidArgError{"username", "Jane"}))
	g.Expect(game.GetState()).To(Equal(Playing))

	// Claim with non-set
	s = game.Board.FindSet(false)
	err = game.ClaimSet("Joe", *s)
	g.Expect(err).To(Succeed())
	g.Expect(game.GetState()).To(Equal(Playing))

	// Claim with set
	s = game.FindExpandSet()
	err = game.ClaimSet("Joe", *s)
	g.Expect(err).To(Succeed())
	g.Expect(game.GetState()).To(Equal(SetClaimed))

	// Claim in claimed state fails
	s = game.Board.FindSet(false)
	err = game.ClaimSet("Jane", *s)
	g.Expect(err).To(MatchError(InvalidStateError{"ClaimSet", "round already claimed by Joe"}))
	g.Expect(game.GetState()).To(Equal(SetClaimed))

	// Expand in claimed state fails
	err = game.Expand()
	g.Expect(err).To(MatchError(InvalidStateError{"Expand", "only valid in claim state"}))
	g.Expect(game.GetState()).To(Equal(SetClaimed))

	// Next succeeds in claimed state, compresses expanded board
	oldLen := len(game.Board)
	err = game.NextRound()
	g.Expect(err).To(Succeed())
	g.Expect(game.GetState()).To(Equal(Playing))
	g.Expect(len(game.Board)).To(Equal(oldLen - SetLen))
}

func TestGamesLoop(t *testing.T) {
	g := NewGomegaWithT(t)
	for i := 0; i < nTestGames; i++ {
		t.Log("Game:", i)
		usernames := getUsernames()
		game, err := NewGame(usernames...)
		g.Expect(err).To(Succeed())
		for s := game.FindExpandSet(); s != nil; s = game.FindExpandSet() {
			t.Log("Board:", game.Board)
			t.Log("Deck len:", len(game.Deck))
			u := usernames[rand.Intn(len(usernames))]
			t.Logf("Username: %s found set: %v %v %v", u, s[0], s[1], s[2])

			uOldScore := len(game.Players[u].Sets)
			err = game.ClaimSet(u, *s)
			g.Expect(err).To(Succeed())
			g.Expect(game.ClaimedUsername).To(Equal(u))
			uNewScore := len(game.Players[u].Sets)
			g.Expect(uNewScore).To(Equal(uOldScore + 1))

			oldBoardLen := len(game.Board)
			oldDeckLen := len(game.Deck)
			err = game.NextRound()
			g.Expect(err).To(Succeed())
			g.Expect(game.ClaimedUsername).To(BeEmpty())
			if oldBoardLen > InitBoardLen || oldDeckLen == 0 {
				g.Expect(len(game.Board)).To(Equal(oldBoardLen - SetLen))
			} else {
				g.Expect(len(game.Board)).To(Equal(oldBoardLen))
			}
		}
		g.Expect(game.Deck).To(BeEmpty())
	}
}
