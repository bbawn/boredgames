package set

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	SetLen       = 3
	NAxes        = 4
	FullDeckLen  = 81
	InitBoardLen = 12
)

type Shading byte

const (
	Filled Shading = iota
	Outline
	Stripe
)

// abbrToShading returns the Shading, given the abbreviation
func abbrToShading(a string) Shading {
	switch strings.ToLower(a) {
	case "f":
		return Filled
	case "o":
		return Outline
	case "s":
		return Stripe
	default:
		return math.MaxUint8
	}
}

//go:generate stringer -type=Shading

type Shape byte

const (
	Diamond Shape = iota
	Oval
	Squiggle
)

// abbrToShape returns the Shape, given the abbreviation
func abbrToShape(a string) Shape {
	switch strings.ToLower(a) {
	case "d":
		return Diamond
	case "o":
		return Oval
	case "s":
		return Squiggle
	default:
		return math.MaxUint8
	}
}

//go:generate stringer -type=Shape

type Color byte

const (
	Green Color = iota
	Purple
	Red
)

// abbrToColor returns the Color, given the abbreviation
func abbrToColor(a string) Color {
	switch strings.ToLower(a) {
	case "g":
		return Green
	case "p":
		return Purple
	case "r":
		return Red
	default:
		return math.MaxUint8
	}
}

//go:generate stringer -type=Color

type State byte

const (
	Playing State = iota
	SetClaimed
)

//go:generate stringer -type=State

// Card is a set game card
type Card struct {
	Color   Color
	Count   byte
	Shading Shading
	Shape   Shape
}

func (c *Card) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if len(s) != 4 {
		return fmt.Errorf("Card string must have len 4: %s", s)
	}
	color := abbrToColor(s[0:1])
	if color < 0 {
		return fmt.Errorf("Invalid color abbreviation: %c", s[0])
	}
	count, err := strconv.Atoi(s[1:2])
	if err != nil {
		return fmt.Errorf("Invalid count: %c", s[1])
	}
	shading := abbrToShading(s[2:3])
	if shading < 0 {
		return fmt.Errorf("Invalid shading abbreviation: %c", s[2])
	}
	shape := abbrToShape(s[3:])
	if color < 0 {
		return fmt.Errorf("Invalid color abbreviation: %c", s[3])
	}
	*c = Card{color, byte(count), shading, shape}
	return nil
}

func (c Card) MarshalJSON() ([]byte, error) {
	s := c.Color.String()[0:1] + strconv.Itoa(int(c.Count)) + c.Shading.String()[0:1] + c.Shape.String()[0:1]
	return json.Marshal(s)
}

func (c *Card) String() string {
	if c == nil {
		return "----"
	}
	return c.Color.String()[0:1] + strconv.Itoa(int(c.Count)) + c.Shading.String()[0:1] + c.Shape.String()[0:1]
}

// CardBase3 is representation of a card as a 4-digit, base-3 integer
type CardBase3 int

// CardBase3ToCard returns the Card for the given CardBase3
func CardBase3ToCard(cb3 CardBase3) *Card {
	return &Card{
		Shading: Shading(cb3 / 27),
		Shape:   Shape((cb3 % 27) / 9),
		Color:   Color((cb3 % 9) / 3),
		Count:   byte(cb3%3) + 1,
	}
}

// CardToCardBase3 returns the CardBase3 for the given Card
func CardToCardBase3(c *Card) CardBase3 {
	return CardBase3(byte(c.Shading*27) + byte(c.Shape*9) + byte(c.Color*3) + (c.Count - 1))
}

// CardTriple set of three cards that are a Potential Set
type CardTriple [SetLen]Card

// IsSet returns true if the given cards are a set, false otherwise
func IsSet(cs CardTriple) bool {
	if cs[0] == cs[1] || cs[1] == cs[2] || cs[0] == cs[2] {
		// Duplicates are not even a valid non-set, warn here?
		return false
	}
	return setMatch(byte(cs[0].Shading), byte(cs[1].Shading), byte(cs[2].Shading)) &&
		setMatch(byte(cs[0].Shape), byte(cs[1].Shape), byte(cs[2].Shape)) &&
		setMatch(byte(cs[0].Color), byte(cs[1].Color), byte(cs[2].Color)) &&
		setMatch(cs[0].Count, cs[1].Count, cs[2].Count)
}

// setMatch return true if the given bytes are a "match" by Set rules:
// meaning they are either all the same or all different
func setMatch(b1, b2, b3 byte) bool {
	ret := (b1 == b2 && b2 == b3) || (b1 != b2 && b2 != b3 && b1 != b3)
	return ret
}

// Deck is a deck of set Set cards
type Deck []*Card

// Pop removes and returns to top of the deck
func (d *Deck) Pop() *Card {
	c := (*d)[len(*d)-1]
	*d = (*d)[:len(*d)-1]
	return c
}

// Board is a layout of cards
type Board []*Card

// FindSet returns a set or non-set from the given board or nil if
// there are none
func (b Board) FindSet(set bool) *CardTriple {
	// Enumerate each set on board
	for i := 0; i < len(b)-2; i++ {
		c1 := b[i]
		if c1 != nil {
			for j := i + 1; j < len(b)-1; j++ {
				c2 := b[j]
				if c2 != nil {
					for k := j + 1; k < len(b); k++ {
						c3 := b[k]
						if c3 != nil {
							ct := CardTriple{*c1, *c2, *c3}
							if IsSet(ct) == set {
								return &ct
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// FindExpandSet return a set on the board, expanding until one is found
// or the game's deck is exhausted
func (g *Game) FindExpandSet() *CardTriple {
	for true {
		s := g.Board.FindSet(true)
		if s != nil {
			return s
		}
		if !g.expandBoard() {
			return nil
		}
	}
	panic("unreachable")
}

// FindCard returns the index of the given card on the Board or -1 if not found
func (b Board) FindCard(c Card) int {
	for i := 0; i < len(b); i++ {
		if b[i] != nil && *b[i] == c {
			return i
		}
	}
	return -1
}

// Player is a participant in a set game
type Player struct {
	Username string
	Sets     []CardTriple
}

// Game is an instance of a set game
type Game struct {
	ID              uuid.UUID
	Players         map[string]*Player
	Deck            Deck
	Board           Board
	ClaimedSet      CardTriple
	ClaimedUsername string
	// TODO(bbawn): do we need a logical timestamp field to detect stale operations?
}

// InvalidArgError indicates an argument is invalid
type InvalidArgError struct {
	Arg   string
	Value string
}

func (e InvalidArgError) Error() string {
	return fmt.Sprintf("Invalid value: %s for arg: %s", e.Value, e.Arg)
}

// InvalidStateError indicates the Method was called for an object that is not in
// the right state for it
type InvalidStateError struct {
	Method  string
	Details string
}

func (e InvalidStateError) Error() string {
	return fmt.Sprintf("Invalid method: %s detail: %s", e.Method, e.Details)
}

func NewGame(usernames ...string) (*Game, error) {
	g := new(Game)
	g.ID = uuid.New()
	g.Players = make(map[string]*Player)
	for _, u := range usernames {
		if u == "" {
			return nil, InvalidArgError{"username", "empty"}
		}
		if _, present := g.Players[u]; present {
			return nil, InvalidArgError{"username", u + " already present"}
		}
		g.Players[u] = &Player{Username: u, Sets: []CardTriple{}}
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
func (g *Game) expandBoard() bool {
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
// If the given cards are not a set or not present in the deck, nil is
// returned and (per game rules) the most recent set in the player's collection
// is returned to the Deck.
// NOTE: we need to add Game.logicalTime to avoid race where set was valid for
// an earlier logicalTime (should fail but not penalize)
//
// If the given set is valid and the cards are all still present on the board,
// the given set is copied to the Game's ClaimedSet (so that it can be displayed
// prior to the next round) and is added to the given player's collection and
// nil is returned.
func (g *Game) ClaimSet(username string, cs CardTriple) error {
	if g.GetState() != Playing {
		return InvalidStateError{"ClaimSet", "round already claimed by " + g.ClaimedUsername}
	}
	p, present := g.Players[username]
	if !present {
		return InvalidArgError{"username", username}
	}
	if !IsSet(cs) {
		g.penalty(p)
		// Illegal move, but not an error (we must update datastore)
		return nil
	}
	for _, c := range cs {
		i := g.Board.FindCard(c)
		if i < 0 {
			g.penalty(p)
			// Illegal move, but not an error (we must update datastore)
			return nil
		}
		g.Board[i] = nil
	}
	p.Sets = append(p.Sets, cs)
	g.ClaimedUsername = username
	g.ClaimedSet = cs
	return nil
}

// Expand adds the next Card triplet when no players can find a set.
// Only valid in playing state.
func (g *Game) Expand() error {
	if g.GetState() != Playing {
		return InvalidStateError{"Expand", "only valid in claim state"}
	}
	g.expandBoard()
	return nil
}

// NextRound transitions a game in Claimed Set state to the next round
func (g *Game) NextRound() error {
	if g.GetState() != SetClaimed {
		return InvalidStateError{"NextRound", "round not yet claimed"}
	}

	if len(g.Board) > InitBoardLen {
		// The board has been expanded, remove remaining empty card slots
		g.compress()
	} else {
		// Deal from deck to replace empty card slots
		for i, _ := range g.Board {
			if g.Board[i] == nil {
				if len(g.Deck) > 0 {
					g.Board[i] = g.Deck.Pop()
				}
			}
		}
	}

	if len(g.Deck) == 0 {
		// Empty Deck, remove remaining empty card slots
		g.compress()
	}

	g.ClaimedUsername = ""
	g.ClaimedSet = CardTriple{}
	return nil
}

func (g *Game) GetState() State {
	if g.ClaimedUsername == "" {
		return Playing
	} else {
		return SetClaimed
	}
}

// Compess removes empty cards from the game board
func (g *Game) compress() {
	newBoard := []*Card{}
	for i, _ := range g.Board {
		if g.Board[i] != nil {
			newBoard = append(newBoard, g.Board[i])
		}
	}
	g.Board = newBoard
}

func (g *Game) penalty(p *Player) {
	if len(p.Sets) > 0 {
		penaltySet := p.Sets[len(p.Sets)-1]
		g.Deck = append(g.Deck, &penaltySet[0], &penaltySet[1], &penaltySet[2])
		p.Sets = p.Sets[:len(p.Sets)-1]
	}
}
