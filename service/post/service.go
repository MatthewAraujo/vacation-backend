package post

import (
	"fmt"
	"log"
	"os"

	"github.com/MatthewAraujo/vacation-backend/types"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type location struct {
	latitude  string
	longitude string
}

func GetPhotoInfos() (types.PhotoInfo, error) {
	folderPath, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("erro ao obter o diretório de trabalho atual: %v", err)
		fmt.Println(err)
		return types.PhotoInfo{}, err
	}
	folderPath += "/tmp"

	// Abrir a pasta
	folder, err := os.Open(folderPath)
	if err != nil {
		fmt.Println("Erro ao abrir a pasta:", err)
		return types.PhotoInfo{}, err
	}
	defer folder.Close()

	// Ler o conteúdo da pasta
	files, err := folder.Readdir(1) // Obter apenas um arquivo
	if err != nil {
		fmt.Println("Erro ao ler conteúdo da pasta:", err)
		return types.PhotoInfo{}, err
	}

	// Verificar se há arquivos na pasta
	if len(files) == 0 {
		fmt.Println("Pasta vazia.")
		return types.PhotoInfo{}, nil
	}

	// Obter o nome do primeiro arquivo
	firstName := files[0].Name()
	photoUrl := folderPath + "/" + firstName

	location := getLocation(photoUrl)

	return types.PhotoInfo{
		PhotoURL: photoUrl,
		Location: location,
	}, nil
}

func getLocation(fp string) string {
	fname := fp

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	loc := location{}
	err = x.Walk(&loc)
	if err != nil {
		log.Fatal("Error walking the exif data", err)
		log.Fatal(err)
	}

	realLocation := getLatitudeAndLongitude(loc)
	return realLocation
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
func getLatitudeAndLongitude(gpsInfo location) string {
	return fmt.Sprintf("%s,%s", gpsInfo.latitude, gpsInfo.longitude)
}
