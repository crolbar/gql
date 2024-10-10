package util

import (
	"fmt"
	"os"
)

func Logg(s string) {
    file, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    _, err = file.WriteString(s + "\n")
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }
}
