# QR Stream Generator
The Go implementation of QR stream generator. This utility is used to encode any small files into a series of QR codes, named "QR stream". Then user could use a scanner (normally a mobile phone) to decode the QR stream to get the original file. This is useful when there's a networking barrier between the user and the server inside a private network. It helps for sharing files in such situation.

## Try the demo QR stream scanner below: 
https://violet-shaylah-82.tiiny.site

## Sample QR stream
<img src="https://violet-shaylah-82.tiiny.site/demo-qr-stream.gif" width="250" height="250">

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

## Download binaries
See the <a href="https://github.com/vence722/qr-stream-generator-go/releases/tag/v1.0.0">release page</a>.

## Build from source code
Clone this repository, then run the following build command:

For Windows(64bit): 
```
GOOS=windows GOARCH=amd64 go build -o qr-stream-gen.exe
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
