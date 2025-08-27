package main

import (
	"rag/embed"

	log "github.com/sirupsen/logrus"
)

func main() {
	words := []string{
		"数字经济",
		"数字中国",
	}
	vecs, err := embed.Embed(words)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range vecs {
		log.WithField("dim", len(v)).Infof("%+v", v[:10])
	}
}
