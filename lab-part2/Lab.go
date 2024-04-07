package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Estructura para pasar datos a la plantilla HTML
type PageData struct {
	Hostname string
	Tema     string
	TemaDesc string
	Imagenes []ImageInfo
}

// Estructura para almacenar información de imagen junto con su texto
type ImageInfo struct {
	Nombre    string
	Texto     string
	Link      string
	NombreImg string
	Base64    string
}

func main() {

	enlaces := map[string]string{
		"React JS.png": "https://react.dev/",
		"Spring.png":   "https://spring.io/",
		"OnRails.png":  "https://rubyonrails.org/",
		"Angular.png":  "https://angularjs.org/",
		"Django.jpg":   "https://www.djangoproject.com/",
		"Express.png":  "https://expressjs.com/es/",
		"Flask.png":    "https://flask.palletsprojects.com/en/3.0.x/",
		"Laravel.jpg":  "https://laravel.com/",
		"Symfony.png":  "https://symfony.es/",
		"Vue JS.jpg":   "https://vuejs.org/",

		"C#.png":         "https://learn.microsoft.com/es-es/dotnet/csharp/",
		"Golang.png":     "https://go.dev/",
		"Java.png":       "https://www.java.com/es/",
		"JavaScript.png": "https://developer.mozilla.org/es/docs/Web/JavaScript",
		"Kotlin.jpg":     "https://kotlinlang.org/",
		"Php.jpg":        "https://www.php.net/",
		"Python.png":     "https://www.python.org/",
		"R.png":          "https://www.r-project.org/",
		"SQL.png":        "https://es.wikipedia.org/wiki/SQL",
		"TypeScript.png": "https://www.typescriptlang.org/",

		"Debian.jpg":     "https://www.debian.org/",
		"Linux mint.png": "https://linuxmint.com/",
		"Red hat.png":    "https://www.redhat.com/es",
		"Ubuntu.png":     "https://ubuntu.com/",
	}

	var nombresImagenes []string

	// Obtener argumentos de la línea de comandos
	args := os.Args

	if len(args) < 4 {
		fmt.Println("Uso: <puerto> <directorio> <tema>")
		return
	}

	puerto := args[1]
	directorio := args[2] + "\\" + args[3]
	tema := args[3]
	temaDesc := "default"

	if tema == "Lenguajes de programacion" {
		temaDesc = "Los lenguajes de programación son herramientas que permiten a los programadores comunicarse con las computadoras, dándoles instrucciones para realizar tareas específicas mediante un conjunto de reglas y sintaxis definidas.		"
	} else if tema == "Frameworks" {
		temaDesc = "Son conjuntos de herramientas y librerías predefinidas que proporcionan una estructura y funcionalidades comunes para facilitar el desarrollo de software, permitiendo a los programadores enfocarse en la lógica específica de sus aplicaciones en lugar de reinventar soluciones genéricas."
	} else if tema == "Linux" {
		temaDesc = "Es un sistema operativo de código abierto y gratuito basado en el núcleo Linux, desarrollado por una comunidad de colaboradores en todo el mundo. Es conocido por su estabilidad, seguridad y flexibilidad, y es ampliamente utilizado en una variedad de dispositivos, desde servidores hasta dispositivos móviles y sistemas integrados. Linux ofrece una amplia gama de distribuciones, cada una adaptada para diferentes propósitos y preferencias de los usuarios, lo que lo convierte en una opción popular tanto para usuarios domésticos como para empresas."
	}

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

	// Crear una lista de ImageInfo para almacenar el nombre de la imagen y su texto correspondiente
	var imagenesInfo []ImageInfo

	for _, archivo := range archivos {
		nombreArchivo := archivo.Name()
		extension := strings.ToLower(filepath.Ext(nombreArchivo))
		if _, ok := formatosAceptados[extension]; ok {
			imagenes = append(imagenes, nombreArchivo)
			nombresImagenes = append(nombresImagenes, nombreArchivo)

			// Construir el nombre del archivo de texto correspondiente
			nombreArchivoTxt := strings.TrimSuffix(nombreArchivo, extension) + ".txt"
			rutaArchivoTxt := filepath.Join(directorio, nombreArchivoTxt)

			// Leer el contenido del archivo de texto
			texto, err := ioutil.ReadFile(rutaArchivoTxt)
			if err != nil {
				fmt.Printf("Error al leer el archivo %s: %v\n", nombreArchivoTxt, err)
				continue
			}

			// Obtener el enlace correspondiente al nombre de la imagen
			enlace, ok := enlaces[nombreArchivo]
			if !ok {
				fmt.Printf("No se encontró enlace para la imagen: %s\n", nombreArchivo)
				continue
			}

			rutaImagen := filepath.Join(directorio, nombreArchivo)
			imagenBase64, err := leerImagenBase64(rutaImagen)
			if err != nil {
				fmt.Printf("Error al leer la imagen %s como base64: %v\n", nombreArchivo, err)
				continue
			}

			// Agregar la información de la imagen (nombre, texto y enlace) a la lista de ImageInfo
			imagenesInfo = append(imagenesInfo, ImageInfo{
				Nombre:    nombreArchivo,
				Texto:     string(texto),
				Link:      enlace,
				NombreImg: nombreArchivo,
				Base64:    imagenBase64,
			})
		}
	}

	// Verificar si hay suficientes imágenes para mostrar
	if len(imagenesInfo) < 4 {
		fmt.Println("No hay suficientes imágenes con texto para mostrar.")
		return
	}

	// Aleatorizar el orden de las imágenes
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(imagenesInfo), func(i, j int) {
		imagenesInfo[i], imagenesInfo[j] = imagenesInfo[j], imagenesInfo[i]
	})

	// Tomar las primeras cuatro imágenes aleatorizadas
	imagenesMostrar := imagenesInfo[:4]

	// Obtener el nombre del host
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error al obtener el nombre del host:", err)
		return
	}

	// Definir la estructura de datos para pasar a la plantilla
	data := PageData{
		Hostname: hostname,
		Tema:     tema,
		TemaDesc: temaDesc,
		Imagenes: imagenesMostrar,
	}

	// Manejador HTTP que sirve el contenido HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Cargar la plantilla HTML
		tmpl := template.Must(template.ParseFiles("template.html"))

		// Ejecutar la plantilla y pasar los datos
		err := tmpl.Execute(w, data)
		if err != nil {
			fmt.Println("Error al ejecutar la plantilla:", err)
			return
		}
	})

	// Servir las imágenes estáticas
	http.Handle("/imagenes/", http.StripPrefix("/imagenes/", http.FileServer(http.Dir(directorio))))

	// Iniciar el servidor en el puerto especificado
	fmt.Println("Servidor escuchando en http://localhost:" + puerto)
	http.ListenAndServe(":"+puerto, nil)
}

func leerImagenBase64(ruta string) (string, error) {
	// Leer la imagen como un array de bytes
	imgBytes, err := os.ReadFile(ruta)
	if err != nil {
		return "", err
	}

	// Codificar los bytes de la imagen en base64
	encodedImg := base64.StdEncoding.EncodeToString(imgBytes)

	return encodedImg, nil
}
