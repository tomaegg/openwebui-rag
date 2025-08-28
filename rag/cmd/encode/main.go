package main

import (
	"os"
	"rag/utils/embed"

	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must provide one chunk")
	}

	words := os.Args[1:]

	vecs, err := embed.Embed(words...)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range vecs {
		log.WithField("dim", len(v)).Infof("%+v", v[:10])
	}
}
