package main

import (
	"fmt"
	_ "log"
	_ "github.com/vsteffen/42_api/constants"
)

func main() {
	init_keys()
	fmt.Println(g_keys)
}
