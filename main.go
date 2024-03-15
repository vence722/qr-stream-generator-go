package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"math"
	"os"
	"path"

	"github.com/sunshineplan/imgconv"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

const (
	ChunkSize = 100
	StageDir  = "stage"
	Delay     = 10
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: qr-stream-gen <input_file> <output_file_in_gif>")
		os.Exit(1)
	}

	sourceFileName := os.Args[1]
	outputFlieName := os.Args[2]

	fileBytes, err := os.ReadFile(sourceFileName)
	if err != nil {
		panic(err)
	}
	fileB64 := base64.URLEncoding.EncodeToString(fileBytes)

	var chunks []string
	for i := 0; i < len(fileB64); i += ChunkSize {
		end := i + ChunkSize
		if end > len(fileB64) {
			end = len(fileB64)
		}
		chunks = append(chunks, fileB64[i:end])
	}

	os.RemoveAll(StageDir)
	os.MkdirAll(StageDir, 0755)
	totalSize := int(math.Round(float64(len(fileB64)/ChunkSize) + 0.5))
	for i, chunk := range chunks {
		hashedFileNameBytes := md5.Sum([]byte(sourceFileName))
		hashedFileNameHex := hex.EncodeToString(hashedFileNameBytes[:])
		header := fmt.Sprintf("[%s:%s:%d:%d]",
			sourceFileName, hashedFileNameHex, i+1, totalSize)
		generateQRCode(header+chunk, StageDir, fmt.Sprintf("stg-%d", i+1))
	}
	generateGIF(StageDir, outputFlieName, Delay)
	os.RemoveAll(StageDir)
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
