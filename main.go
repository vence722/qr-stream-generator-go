package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sunshineplan/imgconv"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

const (
	DefaultChunkSize = 512
	DefaultDelay     = 10
	StageDir         = "stage"
)

var (
	DefaultNumThreads = runtime.NumCPU()
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("Sorry we're encountering fatal error:", err)
		}
	}()

	if len(os.Args) < 3 {
		fmt.Println("Usage: qr-stream-gen [options] <input_file> <output_file_in_gif>")
		fmt.Println("\tOptions: ")
		fmt.Println("\t\t--chunk_size <INT default 512>")
		fmt.Println("\t\t--delay <INT default 10>")
		fmt.Println("\t\t--num_threads <INT default runtime.NumCPU()>")
		os.Exit(1)
	}

	var chunkSize, delay, numThreads int
	flag.IntVar(&chunkSize, "chunk_size", DefaultChunkSize, "QR stream chunk size")
	flag.IntVar(&delay, "delay", DefaultDelay, "frame delay")
	flag.IntVar(&numThreads, "num_threads", DefaultNumThreads, "number of QR code generation threads running in parallel")
	flag.Parse()

	sourceFileName := flag.Arg(0)
	outputFlieName := flag.Arg(1)

	fmt.Printf("Settings: chunk_size=%d, delay=%d, num_threads=%d\n", chunkSize, delay, numThreads)
	fmt.Printf("\rGenerating QR stream [  0%%][>          ]")
	fileBytes, err := os.ReadFile(sourceFileName)
	if err != nil {
		panic(err)
	}
	fileB64 := base64.StdEncoding.EncodeToString(fileBytes)

	var chunks []string
	for i := 0; i < len(fileB64); i += chunkSize {
		end := i + chunkSize
		if end > len(fileB64) {
			end = len(fileB64)
		}
		chunks = append(chunks, fileB64[i:end])
	}

	os.RemoveAll(StageDir)
	os.MkdirAll(StageDir, 0755)
	totalSize := int(math.Round(float64(len(fileB64)/chunkSize) + 0.5))

	wg := &sync.WaitGroup{}
	chLimit := make(chan struct{}, numThreads)
	for i := 0; i < numThreads; i++ {
		chLimit <- struct{}{}
	}

	var finishedChunks int64
	for i, chunk := range chunks {
		<-chLimit
		wg.Add(1)
		go func(i int, chunk string) {
			hashedFileNameBytes := md5.Sum([]byte(sourceFileName))
			hashedFileNameHex := hex.EncodeToString(hashedFileNameBytes[:])
			header := fmt.Sprintf("[%s:%s:%d:%d]",
				sourceFileName, hashedFileNameHex, i+1, totalSize)
			generateQRCode(header+chunk, StageDir, fmt.Sprintf("stg-%d", i+1))
			fc := atomic.AddInt64(&finishedChunks, 1)
			progress := float64(fc) / float64(len(chunks)) * 100
			barLen := int(progress / 10)
			bar := strings.Repeat("=", barLen) + ">" + strings.Repeat(" ", 10-barLen)
			if int(progress) > 0 {
				fmt.Printf("\rGenerating QR stream [%3d%%][%s]", int(progress), bar)
			}
			chLimit <- struct{}{}
			wg.Done()
		}(i, chunk)
	}
	wg.Wait()
	generateGIF(StageDir, outputFlieName, delay)
	os.RemoveAll(StageDir)

	fmt.Println("\nQR stream generation is done! Output file: " + outputFlieName)
}

func generateQRCode(content string, outputDir string, outputFileName string) {
	qrc, err := qrcode.New(content)
	if err != nil {
		panic(err)
	}
	w, err := standard.New(path.Join(outputDir, outputFileName+".png"), standard.WithBuiltinImageEncoder(standard.PNG_FORMAT))
	if err != nil {
		panic(err)
	}
	if err = qrc.Save(w); err != nil {
		panic(err)
	}
	// convert to gif format
	pf, err := os.Open(path.Join(outputDir, outputFileName+".png"))
	if err != nil {
		panic(err)
	}
	pImg, err := png.Decode(pf)
	if err != nil {
		panic(err)
	}
	gf, err := os.Create(path.Join(outputDir, outputFileName+".gif"))
	if err != nil {
		panic(err)
	}
	err = imgconv.Write(gf, pImg, &imgconv.FormatOption{Format: imgconv.GIF})
	if err != nil {
		panic(err)
	}
	pf.Close()
	os.Remove(path.Join(outputDir, outputFileName+".png"))
}

func generateGIF(inputDir string, outputFilePath string, delay int) {
	outGif := &gif.GIF{}
	frameFiles, err := os.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}
	for _, frameFileEntry := range frameFiles {
		frameFile, err := os.Open(path.Join(inputDir, frameFileEntry.Name()))
		if err != nil {
			panic(err)
		}
		frame, err := gif.Decode(frameFile)
		if err != nil {
			panic(err)
		}
		frameFile.Close()

		outGif.Image = append(outGif.Image, frame.(*image.Paletted))
		outGif.Delay = append(outGif.Delay, delay)
	}
	of, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	defer of.Close()
	err = gif.EncodeAll(of, outGif)
	if err != nil {
		panic(err)
	}
}
