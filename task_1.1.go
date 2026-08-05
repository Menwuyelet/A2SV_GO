// Task:  Word Frequency Count
// Write a Go function that takes a string as input 
// and returns a dictionary containing the frequency 
// of each word in the string. Treat words in a 
// case-insensitive manner and ignore punctuation marks.
package main

import (
	"fmt"
	"unicode"
	"strings"
)

func wordFrequency(input string) map[string]int {

	// Convert the input string to lowercase to avoid case sensitivity 
	input = strings.ToLower(input)
	
	// create new string from our input removing all the punctuation marks
	// using unicode.IsPunct method 
	cleaned := ""

	for _, chr := range input{
		if !unicode.IsPunct(chr){
			cleaned += string(chr)
		}
	}

	// split the cleaned input by space as separator and iterate through it
	// counting the frequency of the words using a map
	freq := make(map[string] int)

	splitedWord := strings.Fields(cleaned)

	for _, word := range splitedWord {
	
		freq[word] += 1
		
	}
	
	return freq
}



func main() {
	// example usage
	input := "Hello, world! Hello Go. Go is great; test, test , test , test, test"
	freq := wordFrequency(input)
	fmt.Println(freq)
}
