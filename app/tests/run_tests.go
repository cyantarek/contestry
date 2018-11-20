package main

import "fmt"

func capitalize(s string) string {
	//write your code here
	var result string
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			ss := s[i]+32
			result += string(ss)
		} else {
			result += string(s[i])
		}
	}
	return result
}

func main() {
	var input string
	fmt.Scanln(&input)

	result := capitalize(input)
	fmt.Println(result)
}
