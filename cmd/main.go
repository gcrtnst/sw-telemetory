package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	port := flag.Int("port", DefaultPort, "listen port")
	root := flag.String("root", DefaultRoot, "where to write files")
	title := flag.String("title", DefaultTitle, "default title of data")
	ext := flag.String("ext", DefaultExt, "file extension")
	flag.Parse()

	cfg := NewServerConfig()
	cfg.Root = *root
	cfg.Title = *title
	cfg.Ext = *ext
	srv := NewServer(cfg)
	err := srv.ListenAndServe(*port)
	log.Println(err)
	os.Exit(1)
}

const (
	DefaultPort     = 58592
	DefaultTimeout  = 1 * time.Second
	DefaultRoot     = "."
	DefaultTitle    = "telemetory"
	DefaultExt      = ".csv"
	DefaultDirMode  = os.FileMode(0777)
	DefaultFileMode = os.FileMode(0666)
)
