package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type location struct {
	latitude  string
	longitude string
}

func (p *location) Walk(name exif.FieldName, tag *tiff.Tag) error {
	switch name {
	case exif.GPSLatitude:
		p.latitude = tag.String()
	case exif.GPSLongitude:
		p.longitude = tag.String()
	}
	return nil
}

func main() {
	folderPath, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("erro ao obter o diretório de trabalho atual: %v", err)
		fmt.Println(err)
		log.Fatal(err)
	}
	folderPath += "/tmp"

	// Abrir a pasta
	folder, err := os.Open(folderPath)
	if err != nil {
		fmt.Println("Erro ao abrir a pasta:", err)
		log.Fatal(err)
	}
	defer folder.Close()

	// Ler o conteúdo da pasta
	files, err := folder.Readdir(1) // Obter apenas um arquivo
	if err != nil {
		fmt.Println("Erro ao ler conteúdo da pasta:", err)
		log.Fatal(err)
	}

	// Verificar se há arquivos na pasta
	if len(files) == 0 {
		fmt.Println("Pasta vazia.")
		log.Fatal("Pasta vazia.")
	}

	// Obter o nome do primeiro arquivo
	firstName := files[0].Name()
	fname := folderPath + "/" + firstName

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	location := location{}
	err = x.Walk(&location)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("S : %s\nW: %s\n", location.latitude, location.longitude)
}
