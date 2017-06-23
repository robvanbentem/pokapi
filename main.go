package main

import (
	"fmt"
	"github.com/loganjspears/joker/hand"
	"github.com/loganjspears/joker/jokertest"
	"math/rand"
	"time"
)

type YatesCards struct {
	Cards []*hand.Card
	Ceil  int
}

func (y *YatesCards) Init() {
	rand.Seed(time.Now().UTC().UnixNano())
	y.Ceil = len(y.Cards) - 1
}

func (y *YatesCards) Take() *hand.Card {
	if y.Ceil == 0 {
		y.Ceil = len(y.Cards) - 1
		return y.Cards[0]
	}

	p := rand.Intn(y.Ceil)

	c := y.Cards[y.Ceil]
	picked := *y.Cards[p]
	y.Cards[y.Ceil] = y.Cards[p]
	y.Cards[p] = c

	y.Ceil--

	return &picked
}

func Calc(a *[]*hand.Card, times int, c chan string) {
	avail := make([]*hand.Card, len(*a))
	copy(avail, *a)
	y := YatesCards{avail, 0}
	y.Init()

	for n := 0; n < times; n++ {

		table := []*hand.Card{y.Take(), y.Take(), y.Take(), y.Take()}

		h1 := hand.New(append(hand1, table[0], table[1], table[2], table[3]))
		h2 := hand.New(append(hand2, table[0], table[1], table[2], table[3]))

		result := h1.CompareTo(h2)

		if result > 0 {
			c <- "hand1"
		} else if result < 0 {
			c <- "hand2"
		} else {
			c <- "tie"
		}
	}
}

var hand1, hand2 []*hand.Card

func main() {

	cards := hand.Cards()
	available := make([]*hand.Card, 0)
	inhands := make([]*hand.Card, 0)

	hand1 = jokertest.Cards("2s", "Ah", "2d")
	hand2 = jokertest.Cards("Ks", "Jh", "4c")

	for _, c := range hand1 {
		inhands = append(inhands, c)
	}

	for _, c := range hand2 {
		inhands = append(inhands, c)
	}

	taken := false
	for _, card := range cards {
		for _, inhand := range inhands {

			if inhand == card {
				taken = true
				break
			} else {
			}
		}

		if !taken {
			available = append(available, card)
		}
		taken = false
	}

	fmt.Print("Hnd1: ")
	fmt.Println(hand1)
	fmt.Print("Hnd2: ")
	fmt.Println(hand2)
	fmt.Print("Deck: ")
	fmt.Println(available)

	y := YatesCards{available, 0}
	y.Init()

	c := make(chan string)

	threads := 8
	threadscount := 12500

	m := map[string]int{
		"tie":   0,
		"hand1": 0,
		"hand2": 0,
		"total": 0,
	}

	total := threads * threadscount

	rand.Seed(time.Now().UTC().UnixNano())
	for n := 0; n < threads; n++ {
		go Calc(&available, threadscount, c)
	}

	for {
		select {
		case r := <-c:
			m[r]++
			m["total"]++
		}

		if m["total"] == total {
			break
		}
	}

	f1 := float64(m["hand1"])
	f2 := float64(m["hand2"])

	t := f1 + f2 + float64(m["tie"])

	fmt.Printf("Games: %d, Hnd1: %d, Hnd2: %d, Tie: %d, %.2f/%.2f\n", m["total"], m["hand1"], m["hand2"], m["tie"], f1/t*100, f2/t*100)
}
