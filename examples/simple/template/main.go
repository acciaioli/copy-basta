package main

import (
	"fmt"

    "{{.goModule}}/source"
)

func main() {
    fmt.Println(source.NewUserBio())
}
