package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

//go:embed words.txt
var wordLib string

//go:embed all.txt
var allWordLib string
var wordLibPath = flag.String("w", "", "Path to your custom word list")
var rounds = flag.Int("r", 6, "Number of rounds to play")
var wordArr []string
var allWordArr []string

const (
	greenBg    = "\033[42m\033[30m"
	brownBg    = "\033[47m\033[30m"
	yellowBg   = "\033[43m\033[30m"
	resetColor = "\033[0m"
)

const start = " __          __           _ _      \n \\ \\        / /          | | |     \n  \\ \\  /\\  / /__  _ __ __| | | ___ \n   \\ \\/  \\/ / _ \\| '__/ _` | |/ _ \\\n    \\  /\\  / (_) | | | (_| | |  __/\n     \\/  \\/ \\___/|_|  \\__,_|_|\\___|"

func init() {
	flag.Parse()

	if *wordLibPath != "" {
		b, err := os.ReadFile(*wordLibPath)
		if err != nil {
			log.Println(err)
		} else {
			wordLib = string(b)
		}
	}
	if wordLib == "" {
		log.Fatalln("No word list provided")
	}

	allWords := strings.Split(allWordLib, "\n")
	for _, word := range allWords {
		allWordArr = append(allWordArr, word[:5])
	}

	words := strings.Split(wordLib, "\n")
	i := 0
	defer func() {
		err := recover()
		if err != nil {
			log.Fatalf("line %d, %s", i+1, err)
		}
	}()
	for ; i < len(words); i++ {
		wordArr = append(wordArr, words[i][:5])
	}
}

func main() {
	fmt.Printf("%s\n\n", start)

	newWordle(wordArr)
}

func newWordle(words []string) {
	rand.Seed(time.Now().Unix())
	word := words[rand.Intn(len(words))]

	for r := 0; r < *rounds; r++ {
		fmt.Printf("Round %d/%d\n", r+1, *rounds)
		fmt.Printf(">")

		var guess string
		_, _ = fmt.Scanln(&guess)

		if !isLetters(guess) || len(guess) != 5 || !in(guess, allWordArr) {
			fmt.Printf("Invalid guess\n\n")
			r = r - 1
			continue
		}

		var output string

		for i := 0; i < 5; i++ {
			switch {
			case guess[i] == word[i]:
				output += greenBg + string(guess[i]) + resetColor
			case contains(guess[i], word):
				output += yellowBg + string(guess[i]) + resetColor
			default:
				output += brownBg + string(guess[i]) + resetColor
			}
		}
		fmt.Printf("%s\n\n", output)

		if guess == word {
			fmt.Println("You win!")
			return
		}
	}

	fmt.Println("You failed!")
	fmt.Println("The answer is:", word)
	return
}

// determine if a word is in the word list
func in(str string, arr []string) bool {
	sort.Strings(arr)
	i := sort.SearchStrings(arr, str)
	if i < len(arr) && arr[i] == str {
		return true
	}
	return false
}

// determine if a word has the letter in it
func contains(letter uint8, word string) bool {
	for i := 0; i < len(word); i++ {
		if word[i] == letter {
			return true
		}
	}
	return false
}

// determine if the input is combined by letters
func isLetters(str string) bool {
	match, _ := regexp.MatchString(`^[A-Za-z]+$`, str)
	return match
}
