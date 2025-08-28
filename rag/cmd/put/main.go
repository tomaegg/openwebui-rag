package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"rag/milvus"
	"rag/utils/embed"
	"rag/utils/splitter"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Process struct {
	Dir *os.Root
	Cli *milvus.RagCli
}

func NewProcess(dir string) *Process {
	root, err := os.OpenRoot(dir)
	if err != nil {
		log.Fatal(err)
	}
	return &Process{
		Cli: milvus.NewRagCli(context.TODO()),
		Dir: root,
	}
}

func (p *Process) Accept(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".txt", ".md":
		return true
	default:
		return false
	}
}

func (p *Process) Do(ctx context.Context, path string) error {
	if !p.Accept(path) {
		log.WithField("file", path).Info("skip unsupported type")
		return nil
	}

	f, err := p.Dir.Open(path)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	// 1. split the passage
	words, err := splitter.SplitContent(string(content))
	if err != nil {
		return err
	}

	// 2. embed the passage
	vecs, err := embed.Embed(words...)
	if err != nil {
		return err
	}

	// 3. insert to Milvus DB
	req := milvus.UpsertReq{
		Vector:  vecs,
		Content: words,
	}

	count, err := p.Cli.Upsert(ctx, req)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"passage":      path,
		"upsert_count": count,
	}).Info("process passage done")

	return nil
}

func (p *Process) Done() {
	p.Cli.Close(context.TODO())
}

func entry() {
	// 解析命令行参数
	dirFlag := flag.String("dir", "", "Root directory for files (default current working directory)")
	flag.Parse()

	// 获取根目录
	root := *dirFlag
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("cannot get current working directory:", err)
		}
		root = cwd
	}

	log.Infof("working directory: %s", root)

	// 剩下参数是文件路径
	files := flag.Args()
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "must provide at least one file path")
		os.Exit(1)
	}

	work(root, files)
}

func work(root string, files []string) {
	// work
	p := NewProcess(root)
	const queueSize = 20
	workChan := make(chan string, queueSize)
	go func() {
		for _, f := range files {
			workChan <- f
		}
		close(workChan)
	}()

	const numWorkers = 4
	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range workChan {
				err := p.Do(context.TODO(), f)
				if err != nil {
					log.WithError(err).Errorf("error process: %s", f)
				}
			}
		}()
	}

	wg.Wait()
	p.Done()

	log.WithField("count", len(files)).Info("embed done")
}

func main() {
	log.SetLevel(log.InfoLevel)
	entry()
}
