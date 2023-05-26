package main

import (
	"fmt"

	"github.com/WeiAnAn/url-shortener/internal/utils"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println(utils.GenerateURL(7))
	}

	fmt.Println(utils.GenerateURL(7))
}
