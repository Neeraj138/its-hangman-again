package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
	// fmt.Println("...closing file")
}

func getWordFromLine(line []byte) (word string) {
	trimLine := strings.TrimSpace(string(line))
	word = strings.Split(trimLine, " ")[0]
	return word
}

// hangman2.txt contains words for hangman, but needs a little preprocessing
// if you are using a different dataset, modify the below function accordingly
// this function should always return the slice of words for hangman
func getHangmanWords() (words []string) {
	fileInput, err := os.Open("./hangman2.txt")

	if err != nil {
		panic(err)
	}

	defer closeFile(fileInput)

	reader := bufio.NewReader(fileInput)
	currWord := ""
	words = make([]string, 0)

	for {
		currLine, _, err := reader.ReadLine()

		if err != nil && err != io.EOF {
			panic(err)
		}

		if err == io.EOF {
			break
		}

		currWord = getWordFromLine(currLine)
		words = append(words, currWord)
	}

	return words
}

func getRandomWord(words []string) (word string) {
	randIndex := randGen.Intn(len(words))
	return words[randIndex]
}

func checkIfAlreadyGuessed(guessedLetters map[rune]bool, currGuess rune) (isAlreadyGuessed bool) {
	_, isAlreadyGuessed = guessedLetters[unicode.ToLower(currGuess)]
	return isAlreadyGuessed
}

func main() {
	hangmanWords := getHangmanWords()
	wordToGuess := getRandomWord(hangmanWords)
	fmt.Println(wordToGuess)
	guessedLetters := map[rune]bool{'n': true, 'e': true}
	fmt.Println(checkIfAlreadyGuessed(guessedLetters, 'n'))
}

/*
when a user enters something ,following things are possible
-> invalid input -> anything other than english alphabet should be invalid
-> already guessed letter -> for this, need to store each new guess
-> new letter
-> game end
*/
