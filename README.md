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
--num_threads <INT default runtime.NumCPU()>
```

Example:
```
# basic usage
./qr-stream-gen in.zip out.gif

# specify options
./qr-stream-gen --chunk_size 1024 --delay 30 --num_threads in.zip out.gif
```

## Download binaries
See the [release page](https://github.com/vence722/qr-stream-generator-go/releases/latest).

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

## Limitations
Due to the low encoding efficiency of the QR code (a single QR code can only bring around 1~3 KB data), it is not suggested to use this tool to encode files larger than 100KB. Normally it is used for sharing codes or other small text files. Also please make sure you're using it legally - DO NOT USE IT TO STEAL SENSITIVE DATA.