// new Deck -> Create a list of playing cards. Essentially an array of strings
// print -> Log out the contents of a deck of cards
// shuffle -> shuffles all the cards in a deck
// deal -> create a hand of cards
// saveToFile ->  Save a list a cards to a file on local
// newDechFromFile -> Load a list of cards from the local machine

// Dynamic Types -> Javascript, Ruby, Python
// Static Types -> C++, Goland, Java
// Basic Go Types -> bool, string, float, int
// We have array vs slice
// array fixed length list but slice an array can grow or shrink
// all element same type
package main

func main() {
	//cards := newDeck()
	//
	//hand, reaminingDeck := deal(cards, 5)
	//hand.print()
	//reaminingDeck.print()

	//cards := newDeck()
	//fmt.Println(cards.toString())
	//cards.saveToFile("my_cards.txt")
	cards := newDeckFromFile("my_cards.txt")
	//fmt.Println(cards)
	cards.shuffle()
	cards.print()
}
