package ipfs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/nathanjohnson320/rip-coin/rip"
)

// TODO:
// Hit this URL to grab the latest block
// https://ipfs.io/ipfs/api/v0/get?arg=<OLD HASH>&archive=false&compress=false&compression-level=-1

// Hit this URL to get the hash for a document:
// https://ipfs.io/docs/api/#api-v0-add
// use the only-hash param to get the hash

func init() {
	// This only checks if ipfs is installed
	cmd := "ipfs"
	args := []string{}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		// Start the daemon for them
		fmt.Println("ipfs not available, please start with `ipfs daemon --enable-pubsub-experiment`")
		os.Exit(1)
	}
}

// Online checks if IPFS daemon is up
func Online() {
	_, err := http.Get("http://localhost:5001")
	if err != nil {
		fmt.Println("IPFS not running, starting manually...")
		args := []string{"daemon", "--enable-pubsub-experiment"}
		cmd := exec.Command("ipfs", args...)

		// Pipe the output to our terminal
		stdOut, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("Error creating StdoutPipe for Cmd", err)
			os.Exit(1)
		}

		// Log output to the Terminal
		scanner := bufio.NewScanner(stdOut)
		go func() {
			for scanner.Scan() {
				fmt.Printf("IPFS: %s\n", scanner.Text())
			}
		}()

		err = cmd.Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("IPFS connected...")
	}
}

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
				fmt.Println(err)
				break
			}

			// We terminate on the chinese character for dead
			// or if we're at 1mb of buffer
			if c == '死' || len(s) > 1048576 {
				str := string(s)
				s = []rune{}
				sub <- str
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
func Publish(t, msg string) {
	cmd := "ipfs"
	args := []string{"pubsub", "pub", t, msg + string('死')}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// AddToIPFS publishes the block to IPFS
func AddToIPFS(b *rip.Block) *http.Response {
	jsonBlock, err := json.Marshal(b)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error marshalling block for ipfs add", err)
		os.Exit(1)
	}
	// create ipfs req to upload the json
	req, _ := http.NewRequest("POST", "http://localhost:5001/api/v0/add?recursive=false", bytes.NewBuffer(jsonBlock))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error pushing to ipfs chain", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return resp
}

// Get returns a file for a given hash
func Get(hash string) ([]byte, error) {
	return exec.Command("ipfs", "cat", hash).Output()
}
