package main

import (
	"context"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/colorsakura/isync/internal/config"
	"github.com/colorsakura/isync/internal/webdav"
	"github.com/spf13/cobra"
)

const MAX_BUFFER_SIZE = 8388616

func init() {}

var rootCmd = &cobra.Command{
	Use:   "isync",
	Short: "iSync is a sync service supported webdav.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cfg, err := config.UmarshalConfig(loadfile("/home/iFlygo/Documents/Projects/isync/config.toml"))
		if err != nil {
			cancel()
			log.Print("failed to load config.toml", err)
			os.Exit(-1)
		}
		log.Print(cfg)

		c := make(chan os.Signal, 1)
		// Trigger graceful shutdown on SIGINT or SIGTERM.
		// The default signal sent by the `kill` command is SIGTERM,
		// which is taken as the graceful shutdown signal for many systems, eg., Kubernetes, Gunicorn.
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-c
			cancel()
		}()

		go servStart(ctx, cfg)

		<-ctx.Done()
	},
}

func loadfile(f string) []byte {
	buf, err := os.ReadFile(f)
	if err != nil {
		log.Print("failed to read file", err)
		return nil
	}
	return buf
}

func servStart(ctx context.Context, cfg *config.Config) {
	log.Print("Start run server")
	go func() {
		singleServ(ctx, cfg)
	}()
}

func singleServ(ctx context.Context, cfg *config.Config) {
	log.Printf("%s sign in %s", cfg.Account, cfg.Address)
	c := webdav.NewClient(cfg.Address, cfg.Account, cfg.Password)
	if err := c.Connect(); err != nil {
		log.Printf("failed to connect %s", cfg.Address)
	}

	buf := make([]byte, MAX_BUFFER_SIZE) // 1M

	f := os.DirFS(cfg.Target)
	fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			r, _ := f.Open(path)
			defer r.Close()

			n, err := r.Read(buf)
			if err != nil {
				log.Printf("failed to read %s to buffer", path)
			}

			// FIXME: 无法上传大于 buf 的文件
			log.Printf("upload %s to %s", d.Name(), cfg.Directory+"/"+path)
			err = c.Write(cfg.Directory+"/"+path, buf[:n], 0664)
			if err != nil {
				log.Printf("failed to upload %s", d.Name())
			}
		}

		return nil
	})

	log.Println("finish to upload files")
}

func main() {
	rootCmd.Execute()
}
