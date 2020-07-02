package main

import (
	"github.com/Floor-Gang/faq-manager/internal"
	util "github.com/Floor-Gang/utilpkg"
)

func main() {
	internal.Start("./config.yml")

	util.KeepAlive()
}
