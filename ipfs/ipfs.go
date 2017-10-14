package ipfs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// TODO:
// Hit this URL to grab the latest block
// https://ipfs.io/ipfs/api/v0/get?arg=<OLD HASH>&archive=false&compress=false&compression-level=-1

// Hit this URL to get the hash for a document:
// https://ipfs.io/docs/api/#api-v0-add
// use the only-hash param to get the hash

// Subscribe listens on a message for new data
func Subscribe(sub chan string) {
	cmdName := "ipfs"
	args := []string{"pubsub", "sub", "rip-coin-tx"}

	cmd := exec.Command(cmdName, args...)
	r, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stdin = r
	reader := bufio.NewReader(cmd.Stdin)

	var err error

	go func() {
		var s []rune
		for {
			c, _, err := reader.ReadRune()
			if err != nil {
				break
			}

			// We terminate on the chinese character for dead
			if c == '死' {
				sub <- string(s)
				s = []rune{}
			} else {
				s = append(s, c)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}
}

// Publish pushes a msg to the rip-coin-tx topic
func Publish(msg string) {
	cmd := "ipfs"
	args := []string{"pubsub", "pub", "rip-coin-tx", msg + "死"}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
