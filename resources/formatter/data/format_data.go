package main

import (
	"bufio"
	"os"

	"github.com/jeromelesaux/martine/log"
)

func main() {
	log.Default()
	scanner := bufio.NewScanner(os.Stdin)
	count := 0
	for scanner.Scan() {
		if count%8 == 0 {
			log.GetLogger().Info("\ndb ")
		}
		log.GetLogger().Info("%s", scanner.Text())
		if (count+1)%8 != 0 {
			log.GetLogger().Info(", ")
		}
		count++
	}
}
