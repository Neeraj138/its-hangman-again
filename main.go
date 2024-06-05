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

func getWordFromLine(line []byte) string {
	trimLine := strings.TrimSpace(string(line))
	word := strings.Split(trimLine, " ")[0]
	return word
}

// hangman2.txt contains words for hangman, but needs a little preprocessing
// if you are using a different dataset, modify the below function accordingly
// this function should always return the slice of words for hangman
func getHangmanWords() []string {
	fileInput, err := os.Open("./words/hangman2.txt")

	if err != nil {
		panic(err)
	}

	defer closeFile(fileInput)

	reader := bufio.NewReader(fileInput)
	currWord := ""
	words := make([]string, 0)

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

func getRandomWord(words []string) string {
	randIndex := randGen.Intn(len(words))
	return words[randIndex]
}

func isAlreadyGuessed(guessedLetters map[rune]bool, currGuess rune) bool {
	_, isAlreadyGuessed := guessedLetters[currGuess]
	return isAlreadyGuessed
}

func updateGameState(wordToGuess string,
	currGuess rune,
	currHangmanState, noOfMatches int8) (int8, int8) {
	isCorrectGuess := strings.ContainsRune(wordToGuess, currGuess)
	newHangmanState := currHangmanState
	if isCorrectGuess {
		noOfMatches++
		fmt.Println("Yay! A correct guess...continue the streak.")
	} else {
		fmt.Println("Wrong guess :( , moving closer to death.")
		newHangmanState++
	}
	printHangman(newHangmanState)
	return newHangmanState, noOfMatches

}

func printHangman(hangmanState int8) {
	if hangmanState == -1 {
		return
	}
	assetFilePath := fmt.Sprintf("./assets/hangman%d.txt", hangmanState)
	assetFile, err := os.ReadFile(assetFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(assetFile))
}

func getUnqCharsCount(wordToGuess string) int {
	unqWord := make(map[rune]struct{}) // no internal set in golang, so using this as empty struch occupies no memory
	for _, char := range wordToGuess {
		_, isPresent := unqWord[char]
		if !isPresent {
			unqWord[char] = struct{}{}
		}
	}
	return len(unqWord)
}

func printWordState(word string, guessedLetters map[rune]bool) {
	var isPresent bool
	for _, ch := range word {
		_, isPresent = guessedLetters[ch]
		if isPresent {
			fmt.Printf("%c ", ch)
		} else {
			fmt.Printf("_ ")
		}
	}
	fmt.Printf("\n\n\n")
}

func main() {
	hangmanWords := getHangmanWords()
	wordToGuess := getRandomWord(hangmanWords)
	guessedLetters := map[rune]bool{}
	userInputReader := bufio.NewReader(os.Stdin)
	var hangmanState int8 = -1
	var noOfMatches int8 = 0
	noOfUnqChars := getUnqCharsCount(wordToGuess)
	guessesOrdered := []string{}

	printWordState(wordToGuess, guessedLetters)

	for {
		fmt.Print("Enter an English alphabet: ")
		userInput, err := userInputReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		userInput = strings.TrimSpace(userInput)

		if len(userInput) > 1 {
			fmt.Printf("Invalid input entered. Please enter an English alphabet.\n\n")
			continue
		}

		currGuess := unicode.ToLower(rune(userInput[0]))
		if currGuess < 'a' || currGuess > 'z' {
			fmt.Printf("Invalid input entered. Please enter an English alphabet.\n\n")
			continue
		}

		if isAlreadyGuessed(guessedLetters, currGuess) {
			fmt.Printf("Already guessed %c\n", currGuess)
			fmt.Printf("Try a new letter. \n\n")
		} else {
			guessedLetters[currGuess] = true
			guessesOrdered = append(guessesOrdered, string(currGuess))
			hangmanState, noOfMatches = updateGameState(wordToGuess, currGuess, hangmanState, noOfMatches)
		}

		fmt.Println("Guesses till now:", guessesOrdered)
		printWordState(wordToGuess, guessedLetters)

		if hangmanState >= 8 {
			fmt.Println("Game Lost :(")
			fmt.Println("The word is", wordToGuess)
			break
		}
		if int(noOfMatches) == noOfUnqChars {
			fmt.Println("Game won! :)")
			printHangman(9)
			break
		}
	}

}

/*
when a user enters something ,following things are possible
-> invalid input -> anything other than english alphabet should be invalid
-> already guessed letter -> for this, need to store each new guess
-> new letter
-> game end
*/
