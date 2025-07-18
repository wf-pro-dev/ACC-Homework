package ai

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/gonfva/docxlib"
)

func ExtractTextFromPPTX(filePath string) (string, error) {
	cmd := exec.Command("pptx2txt.sh", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run script: %v", err)
	}
	return out.String(), nil
}

func ExtractTextFromPDF(filePath string) (string, error) {
	var res_text string
	doc, err := fitz.New(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer doc.Close()
	for n := 0; n < doc.NumPage(); n++ {
		text, err := doc.Text(n)
		if err != nil {
			log.Fatal(err)
		}
		res_text += text
	}
	fmt.Println(res_text)
	return res_text, nil
}

func ExtractTextFromDocx(filePath string) (string, error) {
	var res_text string
	readFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	fileinfo, err := readFile.Stat()
	if err != nil {
		panic(err)
	}
	size := fileinfo.Size()
	doc, err := docxlib.Parse(readFile, int64(size))
	if err != nil {
		panic(err)
	}
	for _, para := range doc.Paragraphs() {
		for _, child := range para.Children() {
			if child.Run != nil {
				res_text += child.Run.Text.Text
			}

		}
		res_text += "\n"
	}
	return strings.TrimSpace(res_text), nil
}
