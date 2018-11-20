package main

import "fmt"

func capitalize(s string) string {
    //write your code here
    return "Cyan Tarek"
}

func main() {
    var input string
    fmt.Scanln(&input)

    result := capitalize(input)
    fmt.Println(result)
}

        