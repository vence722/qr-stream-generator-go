# QR Stream Generator
The Go implementation of QR stream generator

## Demo QR stream scanner site: 
https://violet-shaylah-82.tiiny.site

## Sample QR stream
![demo-qr-stream](https://violet-shaylah-82.tiiny.site/demo-qr-stream.gif)

## Prerequisite
Go version >= 1.22

## Usage
```
qr-stream-gen <input_file> <output_file_in_gif>
```

Options: 
```
--chunk_size <INT default 512>
--delay <INT default 10>
```

## Build from source code
Clone this repository, then run the following build command:

For Windows(64bit): 
```
GOOS=windows GOARCH=amd64 go build -o qr-stream-gen
```

For Linux(64bit):
```
GOOS=linux GOARCH=amd64 go build -o qr-stream-gen
```

For MacOS(64bit, Intel CPU)
```
GOOS=darwin GOARCH=amd64 go build -o qr-stream-gen
```

For MacOS(64bit, M1/M2/M3 CPU)
```
GOOS=darwin GOARCH=arm64 go build -o qr-stream-gen
```
