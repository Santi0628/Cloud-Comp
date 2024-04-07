package main

import (
	"fmt"
	"os"
)

func main() {
	directorio := "C:\\Users\\santi\\Desktop\\Imagenes"

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

	fmt.Printf("Cantidad de fotos: %d\n", len(imagenes))

	for _, imagen := range imagenes {
		fmt.Println(imagen)
	}
}
