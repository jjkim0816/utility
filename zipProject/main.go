package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/ulikunitz/xz"
)

func compressXZ() {
	fmt.Println("compressXZ start : ", time.Now())
	inputPath := "./18.log"
	outputPath := "./18.log.xz"

	// compress file
	compressFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer compressFile.Close()

	w, err := xz.NewWriter(compressFile)
	if err != nil {
		log.Fatal(err)
	}

	original, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	scanner := bufio.NewScanner(original)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if _, err := io.WriteString(w, scanner.Text()+"\n"); err != nil {
			log.Fatal(err)
		}
	}

	if err := w.Close(); err != nil {
		log.Fatal(err)
	}

	// os.Remove(inputPath)

	fmt.Println("compressXZ end : ", time.Now())
}

func depressXZ() {
	fmt.Println("depressXZ start : ", time.Now())
	inputPath := "./18.log.xz"
	outputPath := "./test_18.log"

	compress, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer compress.Close()

	r, err := xz.NewReader(compress)
	if err != nil {
		log.Fatal(err)
	}

	decompress, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(decompress, r); err != nil {
		log.Fatal(err)
	}
	fmt.Println("depressXZ start : ", time.Now())
}

func compressTAR() {
	var inputPath = "./18.log"
	var outputPath = "./test_18.tar.gz"

	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	tw := tar.NewWriter(writer)
	defer tw.Close()

	body, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	if body != nil {
		hdr := &tar.Header{
			Name: path.Base(inputPath),
			Mode: int64(0644),
			Size: int64(len(body)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write(body); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	switch os.Args[1] {
	case "-c":
		compressXZ()
	case "-d":
		depressXZ()
	default:
		compressTAR()
	}
}
