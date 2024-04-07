package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

func main2() {

	var args []string

	for _, arg := range os.Args {
		args = append(args, arg)
	}

	directorio := path() + "\\..\\" + args[1] + "\\" + args[2]

	archivos, err := os.ReadDir(directorio)
	if err != nil {
		fmt.Println("Error al leer la carpeta:", err)
		return
	}

	imagenes := make([]string, 0)

	formatosAceptados := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	for _, archivo := range archivos {
		extension := archivo.Name()[len(archivo.Name())-4:]

		if _, ok := formatosAceptados[extension]; ok {
			imagenes = append(imagenes, archivo.Name())
		}
	}

	encondingbase64(directorio + "\\" + imagenes[0])

	fmt.Printf("Cantidad de fotos: %d\n", len(imagenes))

	for _, imagen := range imagenes {
		fmt.Println(imagen)
	}
}

// func path() string {
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		return "Error: " + err.Error()
// 	}
// 	return dir
// }

func encondingbase64(ruta string) {
	// Leer la imagen como un array de bytes
	imgBytes, err := os.ReadFile(ruta)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Codificar los bytes de la imagen en base64
	encodedImg := base64.StdEncoding.EncodeToString(imgBytes)

	// Imprimir la imagen codificada en base64 en la consola
	fmt.Println(encodedImg)
}
