package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hraban/opus"
)

type packet []byte

func main() {
	// Handle the flags arguments
	var channels = flag.Int("channel", 1, "the number of chanels the opus was recorded with")
	var sampleRate = flag.Int("sr", 16000, "the sample rate the audio was recorded with")
	var fileName = flag.String("f", "", "the opus file to be decoded")
	var outFileName = flag.String("o", "out.pcm", "the name of the file to write to")
	flag.Parse()

	if *fileName == "" {
		fmt.Println("You must define a file to read from with the -f flag.")
		flag.Usage()
		return
	}

	pkts := readDataBase64(*fileName)

	// Create the decoder
	dec, err := opus.NewDecoder(*sampleRate, *channels)
	if err != nil {
		panic(err)
	}

	// Create the final data aggregator
	pcmData := bytes.Buffer{}

	// Iterate through the packets, decoding themm and adding them to the aggregator
	for _, p := range pkts {
		pcm := make([]int16, 2000)
		n, err := dec.Decode([]byte(p), pcm)
		if err != nil {
			continue
		}
		pcmData.Write(int16ToByteSlice(pcm[:n]))
	}

	// Open the file we are going to write to
	outF, err := os.Create(*outFileName)
	if err != nil {
		fmt.Println("could not write to out file")
		return
	}
	defer outF.Close()

	// Write collected data to the file
	outF.Write(pcmData.Bytes())
}

// readDataBase64 reads newline delimited base64 encoded packets
// and returns it's raw data: a slice of byte slices(packet)
func readDataBase64(fname string) []packet {
	pkts := []packet{}

	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	total := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		decoded, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			panic(err)
		}
		total += len(decoded)
		pkts = append(pkts, packet(decoded))
	}

	return pkts
}

// int16ToByteSlice takes a int16 slice and divides it into two seperate bytes
// in little endian form
func int16ToByteSlice(s []int16) []byte {
	ret := []byte{}
	for i := 0; i < len(s); i++ {
		b1 := byte(s[i])
		b2 := byte(s[i] >> 8)
		ret = append(ret, b1)
		ret = append(ret, b2)
	}
	return ret
}
