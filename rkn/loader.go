package rkn

import (
	"bufio"
	"fmt"
	"github.com/zmap/go-iptree/iptree"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func DownloadDump(url, path string) error {
	log.Printf("Downloading dump from %s to %s", url, path)

	client := http.Client{Timeout: Timeout}

	resp, err := client.Get(url)

	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		log.Printf("File %s exists, removing it", path)
		if err = os.Remove(path); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer f.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}

	defer resp.Body.Close()

	if _, err = io.Copy(f, resp.Body); err != nil {
		return err
	}

	log.Println("Dump successful downloaded")

	return nil
}

func DbLoadDump(path string) (*iptree.IPTree, error) {

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 8*1024*1024)

	t := iptree.New()
	t.AddByString("0.0.0.0/0", 0)

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ";")
		if len(fields) < 2 {
			continue
		}
		for _, ip := range strings.Split(fields[0], "|") {
			t.AddByString(strings.TrimSpace(ip), 1)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return t, err
}
