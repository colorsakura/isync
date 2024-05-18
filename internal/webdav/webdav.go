package webdav

import "github.com/studio-b12/gowebdav"

func NewClient(address, account, password string) *gowebdav.Client {
	c := gowebdav.NewClient(address, account, password)
	return c
}
