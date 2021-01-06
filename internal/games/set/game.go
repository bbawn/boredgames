package set

import (
	"fmt"
	"math/rand"
)

const (
	SetLen       = 3
	NAxes        = 4
	FullDeckLen  = 81
	InitBoardLen = 12
)

type Shading byte

const (
	Solid Shading = iota
	Stripe
	Outline
)

//go:generate stringer -type=Shading

type Shape byte

const (
	Oval Shape = iota
	Triangle
	Squiggle
)

//go:generate stringer -type=Shape

type Color byte

const (
	Red Color = iota
	Purple
	Green
)

//go:generate stringer -type=Color

// Card is a set game card
type Card struct {
	Shading Shading
	Shape   Shape
	Color   Color
	Count   byte
}

// CardBase3 is representation of a card as a 4-digit, base-3 integer
type CardBase3 int

// CardBase3ToCard returns the Card for the given CardBase3
func CardBase3ToCard(cb3 CardBase3) *Card {
	return &Card{
		Shading: Shading(cb3 / 27),
		Shape:   Shape((cb3 % 27) / 9),
		Color:   Color((cb3 % 9) / 3),
		Count:   byte(cb3 % 3),
	}
}

// CardToCardBase3 returns the CardBase3 for the given Card
func CardToCardBase3(c *Card) CardBase3 {
	return CardBase3(byte(c.Shading*27) + byte(c.Shape*9) + byte(c.Color*3) + c.Count)
}

// IsSet returns true if the given cards are a set, false otherwise
func IsSet(c1, c2, c3 *Card) bool {
	return setMatch(byte(c1.Shading), byte(c2.Shading), byte(c3.Shading)) &&
		setMatch(byte(c1.Shape), byte(c2.Shape), byte(c3.Shape)) &&
		setMatch(byte(c1.Color), byte(c2.Color), byte(c3.Color)) &&
		setMatch(c1.Count, c2.Count, c3.Count)
}

// setMatch return true if the given bytes are a "match" by Set rules:
// meaning they are either all the same or all different
func setMatch(b1, b2, b3 byte) bool {
	return (b1 == b2 && b2 == b3) || (b1 != b2 && b2 != b3)
}

// Deck is a deck of set cards
type Deck []*Card

// Pop removes and returns to top of the deck
func (d *Deck) Pop() *Card {
	c := (*d)[len(*d)-1]
	*d = (*d)[:len(*d)-1]
	return c
}

// Board is a layout of cards
type Board []*Card

// FindSet returns a set from the given board or empty array if there are none
func (b Board) FindSet() []*Card {
	// Enumerate each set on board
	for i := 0; i < len(b)-2; i++ {
		c1 := b[i]
		for j := i + 1; j <= len(b)-1; j++ {
			c2 := b[j]
			for k := j + 1; k <= len(b); k++ {
				c3 := b[k]
				if IsSet(c1, c2, c3) {
					return []*Card{c1, c2, c3}
				}
			}
		}
	}
	return []*Card{}
}

// FindCard returns the index of the given card on the Board or -1 if not found
func (b Board) FindCard(c *Card) int {
	for i := 0; i < len(b); i++ {
		if *b[i] == *c {
			return i
		}
	}
	return -1
}

// Player is a participant in a set game
type Player struct {
	Username string
	Sets     [][]*Card
}

// Game is an instance of a set game
type Game struct {
	Players         map[string]*Player
	Deck            Deck
	Board           Board
	ClaimedSet      []*Card
	ClaimedUsername string
	// TODO(bbawn): do we need a logical timestamp field to detect stale operations?
}

// InvalidArgError indicates an argument is invalid
type InvalidArgError struct {
	Arg   string
	Value string
}

func (e InvalidArgError) Error() string {
	return fmt.Sprintf("Invalid value: %s for arg: %s", e.Arg, e.Value)
}

// InvalidStateError indicates the
type InvalidStateError struct {
	Method  string
	Details string
}

func (e InvalidStateError) Error() string {
	return fmt.Sprintf("Invalid method: %s detail: %s", e.Method, e.Details)
}

func NewGame(usernames []string) (*Game, error) {
	g := new(Game)
	g.Players = make(map[string]*Player)
	for _, u := range usernames {
		if _, present := g.Players[u]; present {
			return nil, &InvalidArgError{"username", u}
		}
		g.Players[u] = &Player{Username: u}
	}
	g.Deck = make([]*Card, FullDeckLen)
	for i, _ := range g.Deck {
		g.Deck[i] = CardBase3ToCard(CardBase3(i))
	}
	rand.Shuffle(len(g.Deck), func(i, j int) {
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	})
	// Deal cards from deck to board
	g.Board = make([]*Card, InitBoardLen)
	for i, _ := range g.Board {
		g.Board[i] = g.Deck.Pop()
	}
	return g, nil
}

// ExpandBoard adds a new set-length column to the Game's board
// from the Game's deck. Returns true if there were enough cards
// in the deck, otherwise false.
func (g *Game) ExpandBoard() bool {
	for i := 0; i < SetLen; i++ {
		if len(g.Deck) == 0 {
			return false
		}
		g.Board = append(g.Board, g.Deck.Pop())
	}
	return true
}

// ClaimSet validates and processes a set claim from a player.
//
// If a set has already been claimed for this round, an
// InvalidMethodError() is returned.
//
// If the given username is not a player in the Game, an
// InvalidArgError(Arg="username") is returned.
//
// If the given cards are not a set, an InvaidArgError(Arg="set") is returned
// and (per game rules) the most recent set in the player's collection is
// returned to the Deck.
//
// If the given set is valid and any of the cards are no longer present on the
// board, an InvaidArgError(Arg="card") is returned (this can happen if another
// player has claimed the set or another set containing any of the claimed set's
// cards).
//
// If the given set is valid and the cards are all still present on the board,
// the given set is copied to the Game's ClaimedSet (so that it can be displayed
// prior to the next round) and is added to the given player's collection and
// nil is returned.
func (g *Game) ClaimSet(username string, c1, c2, c3 *Card) error {
	if g.ClaimedUsername != "" {
		return &InvalidStateError{"ClaimSet", "round already claimed by " + username}
	}
	p, present := g.Players[username]
	if !present {
		return &InvalidArgError{"username", username}
	}
	if !IsSet(c1, c2, c3) {
		if len(p.Sets) > 0 {
			p.Sets = p.Sets[:len(p.Sets)-1]
		}
		return &InvalidArgError{"set", fmt.Sprintf("%v %v %v", c1, c2, c3)}
	}
	set := []*Card{c1, c2, c3}
	for _, c := range set {
		i := g.Board.FindCard(c)
		if i < 0 {
			return &InvalidArgError{"card", fmt.Sprintf("%v", c1)}
		}
		g.Board[i] = nil
	}
	return nil
}

// NextRound transitions a game in Claimed Set state to the next round
func (g *Game) NextRound() error {
	if len(g.ClaimedSet) != SetLen {
		return &InvalidStateError{"NextRound", "round not yet claimed"}
	}
	if len(g.Deck) < 3 {
		panic(fmt.Sprintf("not enough remaining cards: %v", len(g.Deck)))
	}
	for i, _ := range g.Board {
		if g.Board[i] == g.ClaimedSet[0] ||
			g.Board[i] == g.ClaimedSet[1] ||
			g.Board[i] == g.ClaimedSet[2] {
			if len(g.Deck) > 0 {
				// Card remain in Deck, replace old set card with deck card
				g.Board[i] = g.Deck.Pop()
			} else {
				// Empty Deck, just remove old set cards
				g.Board = append(g.Board[:i], g.Board[i+1:]...)
			}
		}
	}
	g.ClaimedUsername = ""
	g.ClaimedSet = []*Card{}
	return nil
}
