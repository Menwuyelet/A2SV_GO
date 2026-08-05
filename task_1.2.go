// Write a Go function that takes a string as input 
// and checks whether it is a palindrome or not. 
// A palindrome is a word, phrase, number, or other 
// sequence of characters that reads the same forward 
// and backward (ignoring spaces, punctuation, and capitalization).
package main

import(
	"fmt"
	"unicode"
	"strings"
)


func check(s string) bool{
	// normalize the input to lower case
	input := strings.ToLower(s)

	// remove all space and punctuations
	cleaned := ""

	for _, chr := range input{
		if !unicode.IsPunct(chr) && !unicode.IsSpace(chr){
			cleaned += string(chr)
		}
	}

	// iterate over the cleaned version and use two ptr 
	// to verify if it is palindrome or not
	// we change it to slice of runes to iterate over since 
	// iterating over string is not safe
	temp := []rune(cleaned)

	for i, j := 0, len(temp)-1; i<j; i, j = i+1, j-1{
		if temp[i] != temp[j] {
			return false
		}
	}

	return true
}

func main() {
	fmt.Println(check("makam"))
}