package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func config(ctx *cli.Context) (err error) {
	fmt.Println("G_HOME: " + os.Getenv("G_HOME"))
	fmt.Println("G_MIRROR: " + os.Getenv("G_MIRROR"))
	fmt.Println("GOROOT: " + os.Getenv("GOROOT"))
	return nil
}
