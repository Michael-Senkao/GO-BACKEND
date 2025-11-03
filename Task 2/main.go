package main

import (
	"fmt"
	"regexp"
	"strings"
)

//===================== WORD FREQUENCY COUNT ==========================

// getWords extracts words from text (case-insensitive)
func getWords(text string) []string {
    text = strings.ToLower(text)
    re := regexp.MustCompile(`[a-z]+`)
    return re.FindAllString(text, -1)
}

// wordFreqCount returns a map of word frequencies
func wordFreqCount(text string) map[string]int {
    words := getWords(text)
    freq := make(map[string]int)
    for _, word := range words {
        freq[word]++
    }
    return freq
}

//=========================== PALINDROME CHECK ===============================
func isPalindrome(word string) bool{
    word = strings.ToLower(word)
    left := 0
    right := len(word) - 1
    
    for {
        if left >= right{
            break
        }
        if !isAlpha(word[left]){
            left += 1
        }else if !isAlpha(word[right]){
            right -= 1
        }else if word[left] != word[right]{
            return false
        }else{
            left += 1
            right -= 1
        }
    }

    return true
}

// Check if character is an alphabet
func isAlpha(ch byte) bool{
    return ch >= 'a' && ch <= 'z'
}

// Test function
func main() {
    test1 := []struct {
        input  string
        output map[string]int
    }{
        {
            "Hello, world! Welcome to the club.",
            map[string]int{"hello": 1, "world": 1, "welcome": 1, "to": 1, "the": 1, "club": 1},
        },
    }

    test2 := []struct {
        input  string
        output bool
    }{
        {
            "Hello, world! Welcome to the club.",
            false,
        },
        {
            "HellolleH",
            true,
        },
    }

  // Tests for counting word frequency
    for i, test := range test1 {
        result := wordFreqCount(test.input)
        pass := true
        if len(result) != len(test.output) {
            pass = false
        } else {
            for k, v := range test.output {
                if result[k] != v {
                    pass = false
                    break
                }
            }
        }
        if pass {
            fmt.Printf("Test %d passes\n", i+1)
        } else {
            fmt.Printf("Test %d failed: got %v\n", i+1, result)
        }
    }

    // Tests for checking palindrome
    for i, test := range test2 {
        result := isPalindrome(test.input)
        pass := true
        if result != test.output {
            pass = false
        }
        if pass {
            fmt.Printf("Test %d passes\n", i+ len(test1) + 1)
        } else {
            fmt.Printf("Test %d failed: got %v\n", i + len(test1) + 1, result)
        }
    }
}
