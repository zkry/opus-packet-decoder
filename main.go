package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hraban/opus"
)

type packet []byte

func main() {
	// Handle the flags arguments
	var channels = flag.Int("channel", 1, "the number of chanels the opus was recorded with")
	var sampleRate = flag.Int("sr", 16000, "the sample rate the audio was recorded with")
	var fileName = flag.String("f", "", "the opus file to be decoded")
	var outFileName = flag.String("o", "out.pcm", "the name of the file to write to. If writing to stdout as base64, still use this flag to specify type")
	var isBase64 = flag.Bool("b64", false, "prints the output data as a base64 encoded string to stdout")

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

	// Check if the out format is raw pcm
	format := strings.Split(*outFileName, ".")[1]

	// If we are not generating raw data we need to go through a conversion process
	// from pcm to mp3, wav, opus, etc. using ffmpeg
	if format != "pcm" {
		// First we create a temporary directory to write the raw data to
		tmpF, err := ioutil.TempFile("", "raw_recording_data")
		if err != nil {
			fmt.Println("could not write to out file")
			log.Fatal(err)
		}
		defer os.Remove(tmpF.Name())
		tmpF.Write(pcmData.Bytes())
		tmpF.Close()

		tmpOutF, err := ioutil.TempFile("", "conversion_data")
		defer os.Remove(tmpOutF.Name())

		stderr := bytes.Buffer{}
		stdout := bytes.Buffer{}
		cmd := exec.Command("ffmpeg", "-f", "s16le", "-ar", strconv.Itoa(*sampleRate/1000)+"k", "-i", tmpF.Name(), "-y", "-f", format, tmpOutF.Name())
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println(stdout.String())
			fmt.Println(stderr.String())
			fmt.Println("could not convert file to type ", format, ":", err)
			return
		}
		tmpOutF.Close()

		data, err := ioutil.ReadFile(tmpOutF.Name())
		if err != nil {
			fmt.Println("trouble converting file to type ", format, " and saving file")
		}
		pcmData = bytes.Buffer{}
		pcmData.Write(data)
	}
	// At this point pcmData has the data of the file that we want. We will now either
	// write that data to a file or to stdout as base64

	if *isBase64 {
		b64str := base64.StdEncoding.EncodeToString(pcmData.Bytes())
		fmt.Println(b64str)
		return
	}

	outF, err := os.Create(*outFileName)
	if err != nil {
		fmt.Println("could not save to file ", *outFileName, ":", err)
		return
	}
	defer outF.Close()

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
