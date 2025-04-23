package main

import (
	"fmt"
	"math/rand"
)

type Card struct {
	Color        string
	Value        string
	IsAce        bool
	NumericValue uint8
}

type Deck []Card

type Player struct {
	Hand  []Card
	Score uint8
}

type Game struct {
	Deck   Deck
	Player Player
	Dealer Player
}

func NewDeck() Deck {
	colors := []string{"♥️", "♦️", "♣️", "♠️"}
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}
	numericValues := []uint8{2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10, 11}

	var deck Deck
	for _, color := range colors {
		cardValueIndex := 0
		for _, value := range values {
			deck = append(deck, Card{
				Color:        color,
				Value:        value,
				IsAce:        value == "Ace",
				NumericValue: numericValues[cardValueIndex],
			})

			cardValueIndex++
		}
	}

	return deck
}

func PrintCard(card Card) {
	fmt.Printf("%s %s", card.Color, card.Value)
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(*d), func(i, j int) {
		(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
	})
}

func (g *Game) DealCard(dealer bool) {
	card := g.Deck[0]
	g.Deck = g.Deck[1:]
	if dealer {
		g.Dealer.Hand = append(g.Dealer.Hand, card)
	} else {
		g.Player.Hand = append(g.Player.Hand, card)
	}
}

func (p *Player) CalculateScore() {
	p.Score = 0
	aces := 0
	for _, card := range p.Hand {
		p.Score += card.NumericValue
		if card.IsAce {
			aces++
		}
	}

	// handle aces
	for p.Score > 21 && aces > 0 {
		p.Score -= 10
		aces--
	}
}

func (g *Game) PlayerTurn() {
	for {
		g.Player.CalculateScore()
		if g.Player.Score > 21 {
			break
		}

		if g.Player.Score < 17 {
			g.DealCard(false)
		} else {
			break
		}
	}
}

func (g *Game) DealerTurn(playerScore uint8) {
	for {
		g.Dealer.CalculateScore()
		if g.Dealer.Score >= playerScore || g.Dealer.Score >= 17 {
			break
		}

		g.DealCard(true)
	}
}

func playGame() (playerWins, dealerWins, tiesCount uint32) {
	deck := NewDeck()
	deck.Shuffle()

	game := Game{
		Deck: deck,
		Player: Player{
			Hand:  []Card{},
			Score: 0,
		},
		Dealer: Player{
			Hand:  []Card{},
			Score: 0,
		},
	}

	for range 2 {
		game.DealCard(false)
		game.DealCard(true)
	}

	game.Player.CalculateScore()
	game.Dealer.CalculateScore()

	game.PlayerTurn()
	if game.Player.Score == 21 {
		return 1, 0, 0
	} else if game.Player.Score > 21 {
		return 0, 1, 0
	}

	game.DealerTurn(game.Player.Score)
	if game.Dealer.Score > 21 {
		return 1, 0, 0
	} else if game.Dealer.Score == game.Player.Score {
		return 0, 0, 1
	} else if game.Dealer.Score > game.Player.Score {
		return 0, 1, 0
	} else {
		return 1, 0, 0
	}
}

func main() {
	var playerWins uint32 = 0
	var dealerWins uint32 = 0
	var tiesCount uint32 = 0

	for range 100000 {
		w, l, t := playGame()

		playerWins += w
		dealerWins += l
		tiesCount += t
	}

	totalGames := playerWins + dealerWins + tiesCount
	fmt.Printf("+ Total Games: %d\n", totalGames)

	playerWinRate := float64(playerWins) / float64(totalGames) * 100
	tieRate := float64(tiesCount) / float64(totalGames) * 100
	fmt.Printf("  Player Win Rate: %.2f%%\n", playerWinRate)
	fmt.Printf("  Tie Rate: %.2f%%\n", tieRate)
}
