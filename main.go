package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

type FormSuccessResponse struct {
	Success   bool
	File      []byte
	Filename  string
	Extension string
}

type FormErrorResponse struct {
	Success  bool
	ErrorMsg string
}

func main() {
	http.HandleFunc("/", renderTemplate)
	port := 8080
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")

	t, _ := template.ParseFiles("templates/index.html")
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	filename := r.FormValue("filename")
	if filename == "" {
		renderWithError(w, t, "Rilename is required")
	}
	extension := r.FormValue("extension")
	if extension == "" {
		renderWithError(w, t, "Extension is required")
	}
	size := r.FormValue("size")
	if size == "" {
		renderWithError(w, t, "Size is required")
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		renderWithError(w, t, "Failed to parse the size to float")
	}
	byteType := r.FormValue("byteType")
	if byteType == "" {
		renderWithError(w, t, "Byte type is required")
	}
	var divider int
	switch byteType {
	case "MB":
		divider = 1 << 20
	default:
		divider = 1 << 10
	}

	buf := make([]byte, sizeInt/2*divider)

	t.Execute(w, FormSuccessResponse{
		true,
		buf,
		filename,
		extension,
	})
}

func renderWithError(w http.ResponseWriter, t *template.Template, message string) {
	log.Fatal(message)
	t.Execute(w, FormErrorResponse{
		true,
		message,
	})
}
