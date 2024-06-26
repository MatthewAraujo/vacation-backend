package post

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MatthewAraujo/vacation-backend/r2"
	"github.com/MatthewAraujo/vacation-backend/types"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type location struct {
	latitude  string
	longitude string
}

var r2url = "https://pub-573aa5de8ec347b1b4edc5be19999bfc.r2.dev/"

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
	file, err := folder.Readdir(1) // Obter apenas um arquivo
	if err != nil {
		fmt.Println("Erro ao ler conteúdo da pasta:", err)
		return types.PhotoInfo{}, err
	}

	// Verificar se há arquivos na pasta
	if len(file) == 0 {
		fmt.Println("Pasta vazia.")
		return types.PhotoInfo{}, nil
	}

	// Obter o nome do primeiro arquivo
	firstName := file[0].Name()

	s3Service, err := r2.NewR2Service()

	if err != nil {
		log.Fatal("Erro ao criar o serviço R2:", err)
		return types.PhotoInfo{}, err
	}

	// Upload do arquivo para o R2
	image := folderPath + "/" + firstName

	imagePath := filepath.Join(folderPath, firstName)

	fileReader, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Erro ao abrir o arquivo:", err)
		return types.PhotoInfo{}, err
	}

	defer fileReader.Close()

	err = s3Service.UploadFileToR2(context.TODO(), firstName, fileReader)
	if err != nil {
		log.Fatal("Erro ao fazer upload do arquivo para o R2:", err)
		return types.PhotoInfo{}, err
	}

	log.Println("Arquivo enviado para o R2 com sucesso.")
	location := getLocation(image)
	photoUrl := r2url + firstName
	return types.PhotoInfo{
		PhotoURL: photoUrl,
		Location: location,
	}, nil
}

func getLocation(fp string) types.Location {
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

	latitude := parseCoord(loc.latitude)
	longitude := parseCoord(loc.longitude)

	lat := convertCoord(latitude, "S")
	lon := convertCoord(longitude, "W")

	return types.Location{
		Latitude:  lat,
		Longitude: lon,
	}
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
func convertCoord(coord []string, ref string) float64 {

	degress := coord[0]
	minutes := coord[1]
	seconds := coord[2]

	degress = strings.Split(degress, "/")[0]
	minutes = strings.Split(minutes, "/")[0]
	seconds = strings.Split(seconds, "/")[0]

	// Convert degrees to float
	degrees, err := strconv.ParseFloat(degress, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Convert minutes to float
	minutesFloat, err := strconv.ParseFloat(minutes, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Convert seconds to float
	secondsFloat, err := strconv.ParseFloat(seconds, 64)
	if err != nil {
		log.Fatal(err)
	}
	secondsDecimal := secondsFloat / 100000000

	// Calculate decimal degrees
	decimal := degrees + (minutesFloat / 60) + (secondsDecimal / 3600)

	// If the reference indicates south or west, make the decimal negative
	if ref == "S" || ref == "W" {
		decimal = -decimal
	}

	return decimal
}

func parseCoord(coordStr string) []string {
	// Dividir a string da coordenada em partes
	parts := strings.Split(coordStr, ",")

	d := parts[0]
	m := parts[1]
	s := parts[2]

	// Extrair os valores de cada parte
	dVal := extractValue(d)
	mVal := extractValue(m)
	sVal := extractValue(s)
	if dVal != "" && mVal != "" && sVal != "" {
		return []string{dVal, mVal, sVal}
	}
	return nil
}

func extractValue(str string) string {
	start := strings.Index(str, `"`) + 1
	end := strings.LastIndex(str, `"`)

	if start == -1 || end == -1 {
		return ""
	}

	return str[start:end]
}
