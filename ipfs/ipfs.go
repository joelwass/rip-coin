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
	cmdArgs := []string{"pubsub", "sub", "rip-coin-tx"}

	cmd := exec.Command(cmdName, cmdArgs...)
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
			s = append(s, c)
			if c == '死' {
				sub <- string(s)
				s = []rune{}
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}
}
