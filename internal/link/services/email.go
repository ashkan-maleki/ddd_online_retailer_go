package services

import (
	"fmt"
	"log"
)

func SendEmail(addr, text string) {
	log.Println(fmt.Sprintf("Email was sent for %s with content: \n%s ", addr, text))
}
