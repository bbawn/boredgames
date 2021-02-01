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
		c1, c2, c3 *Card
	)
	g := NewGomegaWithT(t)

	// duplicate card, not even a non-set
	c1 = &Card{Red, 1, Filled, Diamond}
	c2 = &Card{Red, 1, Filled, Diamond}
	c3 = &Card{Purple, 1, Filled, Diamond}
	g.Expect(IsSet(c1, c2, c3)).NotTo(BeTrue())

	// not a set
	c1 = &Card{Red, 1, Filled, Diamond}
	c2 = &Card{Red, 1, Filled, Squiggle}
	c3 = &Card{Purple, 1, Stripe, Squiggle}
	g.Expect(IsSet(c1, c2, c3)).NotTo(BeTrue())

	// Shading same, shape different, color different, count different
	c1 = &Card{Purple, 1, Filled, Diamond}
	c2 = &Card{Red, 2, Filled, Squiggle}
	c3 = &Card{Green, 3, Filled, Oval}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())
	// Shading different, shape same, color different, count different
	c1 = &Card{Purple, 1, Filled, Squiggle}
	c2 = &Card{Red, 2, Stripe, Squiggle}
	c3 = &Card{Green, 3, Outline, Squiggle}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
	// Shading different, shape different, color same, count different
	c1 = &Card{Green, 1, Filled, Diamond}
	c2 = &Card{Green, 2, Stripe, Squiggle}
	c3 = &Card{Green, 3, Outline, Oval}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}
	// Shading different, shape different, color different, count same
	c1 = &Card{Purple, 1, Filled, Diamond}
	c2 = &Card{Red, 1, Stripe, Squiggle}
	c3 = &Card{Green, 1, Outline, Oval}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())
	if !IsSet(c1, c2, c3) {
		t.Errorf("expected: %v, %v, %v to be a set", c1, c2, c3)
	}

	// Shading same, shape same, color different, count different
	c1 = &Card{Purple, 1, Filled, Oval}
	c2 = &Card{Red, 2, Filled, Oval}
	c3 = &Card{Green, 3, Filled, Oval}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())

	// Shading different, shape same, color same, count different
	c1 = &Card{Red, 1, Filled, Oval}
	c2 = &Card{Red, 2, Stripe, Oval}
	c3 = &Card{Red, 3, Outline, Oval}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())

	// Shading different, shape different, color same, count same
	c1 = &Card{Red, 2, Filled, Diamond}
	c2 = &Card{Red, 2, Stripe, Squiggle}
	c3 = &Card{Red, 2, Outline, Oval}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())

	// Shading different, shape same, color same, count same
	c1 = &Card{Red, 2, Filled, Oval}
	c2 = &Card{Red, 2, Stripe, Oval}
	c3 = &Card{Red, 2, Outline, Oval}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())

	// Shading same, shape different, color same, count same
	c1 = &Card{Red, 2, Stripe, Diamond}
	c2 = &Card{Red, 2, Stripe, Oval}
	c3 = &Card{Red, 2, Stripe, Squiggle}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())

	// Shading same, shape same, color different, count same
	c1 = &Card{Red, 2, Stripe, Oval}
	c2 = &Card{Green, 2, Stripe, Oval}
	c3 = &Card{Purple, 2, Stripe, Oval}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())

	// Shading same, shape same, color same, count different
	c1 = &Card{Red, 1, Stripe, Oval}
	c2 = &Card{Red, 2, Stripe, Oval}
	c3 = &Card{Red, 3, Stripe, Oval}
	g.Expect(IsSet(c1, c2, c3)).To(BeTrue())
}

func TestInvalidGames(t *testing.T) {
	g := NewGomegaWithT(t)
	usernames := getUsernames()
	game, err := NewGame(usernames...)
	g.Expect(err).To(Succeed())
	g.Expect(len(game.Players)).To(Equal(len(usernames)))
	for _, u := range usernames {
		g.Expect(game.Players).To(HaveKey(u))
	}
	g.Expect(len(game.Deck)).To(Equal(FullDeckLen - InitBoardLen))

	// Board and remaining Deck should comprise full deck
	allCards := make(map[CardBase3]*Card)
	for _, c := range game.Board {
		allCards[CardToCardBase3(c)] = c
	}

	for _, c := range game.Deck {
		allCards[CardToCardBase3(c)] = c
	}
	g.Expect(len(allCards)).To(Equal(FullDeckLen))

	var (
		s []*Card
	)
	// Can only go to next round if there is a claimed set
	err = game.NextRound()
	g.Expect(err).To(MatchError(InvalidStateError{"NextRound", "round not yet claimed"}))

	// Claim with invalid username
	s = game.FindExpandSet()
	err = game.ClaimSet("Jane", s[0], s[1], s[2])
	g.Expect(err).To(MatchError(InvalidArgError{"username", "Jane"}))

	// Claim with non-set
	s = game.Board.FindSet(false)
	err = game.ClaimSet("Joe", s[0], s[1], s[2])
	g.Expect(err).To(Succeed())
}

func TestValidGames(t *testing.T) {
	g := NewGomegaWithT(t)
	for i := 0; i < nTestGames; i++ {
		t.Log("Game:", i)
		usernames := getUsernames()
		game, err := NewGame(usernames...)
		g.Expect(err).To(Succeed())
		for s := game.FindExpandSet(); len(s) > 0; s = game.FindExpandSet() {
			t.Log("Board:", game.Board)
			t.Log("Deck len:", len(game.Deck))
			u := usernames[rand.Intn(len(usernames))]
			t.Logf("Username: %s found set: %v %v %v", u, s[0], s[1], s[2])

			uOldScore := len(game.Players[u].Sets)
			err = game.ClaimSet(u, s[0], s[1], s[2])
			g.Expect(err).To(Succeed())
			g.Expect(game.ClaimedUsername).To(Equal(u))
			uNewScore := len(game.Players[u].Sets)
			g.Expect(uNewScore).To(Equal(uOldScore + 1))
			err = game.NextRound()
			g.Expect(err).To(Succeed())
			g.Expect(game.ClaimedUsername).To(BeEmpty())
		}
		g.Expect(game.Deck).To(BeEmpty())
	}
}
