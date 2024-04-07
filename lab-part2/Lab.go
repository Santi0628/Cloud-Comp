package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// Definir la ruta del directorio de imágenes

	var args []string

	for _, arg := range os.Args {
		args = append(args, arg)
	}

	directorio := args[2] + "\\" + args[3]
	puerto := args[1]
	tema := args[3]

	// Leer los archivos del directorio
	archivos, err := os.ReadDir(directorio)
	if err != nil {
		fmt.Println("Error al leer la carpeta:", err)
		return
	}

	// Filtrar las imágenes con extensiones válidas
	imagenes := make([]string, 0)
	formatosAceptados := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	for _, archivo := range archivos {
		extension := strings.ToLower(archivo.Name()[len(archivo.Name())-4:])
		if _, ok := formatosAceptados[extension]; ok {
			imagenes = append(imagenes, archivo.Name())
		}
	}

	// Verificar si hay suficientes imágenes para mostrar
	if len(imagenes) < 4 {
		fmt.Println("No hay suficientes imágenes para mostrar.")
		return
	}

	// Aleatorizar el orden de las imágenes
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(imagenes), func(i, j int) {
		imagenes[i], imagenes[j] = imagenes[j], imagenes[i]
	})

	// Tomar las primeras tres imágenes aleatorizadas
	imagenesMostrar := imagenes[:4]

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error al obtener el nombre del host:", err)
		return
	}

	// Manejador HTTP que sirve el contenido HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Galería de Fotos</title>
    <style>
        /* Estilos para la barra de navegación */
        .navbar {
            background-color: #333;
            overflow: hidden;
        }

        .navbar a {
            float: left;
            display: block;
            color: #f2f2f2;
            text-align: center;
            padding: 14px 20px;
            text-decoration: none;
        }

        .navbar a:hover {
            background-color: #ddd;
            color: black;
        }

        /* Estilos para la galería de fotos */
        .gallery {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            padding: 20px;
        }

        .gallery .image-container {
            position: relative;
            border-radius: 8px;
            overflow: hidden;
        }

        .gallery img {
            width: 100%;
            height: auto;
        }

        .image-caption {
            position: absolute;
            bottom: 0;
            left: 0;
            width: 100%;
            background-color: rgba(0, 0, 0, 0.7);
            color: #fff;
            padding: 5px;
            box-sizing: border-box;
            text-align: center;
        }

		footer {
            background-color: #333;
            color: #fff;
            text-align: center;
            padding: 10px;
            position: fixed;
            bottom: 0;
            width: 100%;
        }

		#hostname {
            position: absolute;
            top: 20px;
            right: 20px;
            font-size: 24px;
            font-weight: bold;
            color: #FFFFFF;
        }

		h1 {
			text-align: center;
		  }
    </style>
</head>
<body>

<div id="hostname">`+hostname+`</div>

<!-- Barra de navegación -->
<div class="navbar">
    <a href="https://www.facebook.com/profile.php?id=61556918033342">Santiago Garcia Cañas</a>
    <a href="https://www.facebook.com/profile.php?id=100078462336757">Sebastian Carmona Tapasco</a>
    <a href="https://www.facebook.com/nodier.alzatesolano">Nodier Alberto Alzate Solano</a>
</div>

<h1>Tema: `+tema+`</h1>

<!-- Galería de fotos -->
<div class="gallery">
`)
		// Insertar las imágenes en el HTML
		for _, imagen := range imagenesMostrar {
			// Construir la URL completa de la imagen
			imagenURL := "/imagenes/" + imagen
			// Agregar la etiqueta <div> para contener la imagen y el nombre
			fmt.Fprintf(w, `<div class="image-container"><img src="%s" alt="%s"><div class="image-caption">%s</div></div>`, imagenURL, imagen, imagen)
		}

		// Cerrar el HTML
		fmt.Fprint(w, `
</div>

</body>

<footer>
    Universidad del Quindío 2024 - Computación en la nube
</footer>

</html>
`)
	})

	// Servir las imágenes estáticas
	http.Handle("/imagenes/", http.StripPrefix("/imagenes/", http.FileServer(http.Dir(directorio))))

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor escuchando en http://localhost:" + puerto)
	http.ListenAndServe(":"+puerto, nil)
}

func path() string {
	dir, err := os.Getwd()
	if err != nil {
		return "Error: " + err.Error()
	}
	return dir
}
