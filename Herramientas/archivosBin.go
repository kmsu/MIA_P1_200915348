// Operaciones con archivos binarios
package Herramientas

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// En Go cuando manejo en paquetes la funcion debe iniciar con mayuscula para poder ser exportada

// funcion para crear un archivo binario
func CrearDisco(path string) error {
	//asegurar que exista la ruta (el directorio) creando la ruta
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error al crear el disco, path: ", err)
		return err
	}

	//crear el archivo si aun no existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		newFile, err := os.Create(path)
		if err != nil {
			fmt.Println("Error al crear el disco: ", err)
			return err
		}
		defer newFile.Close()
	}
	return nil
}

func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile ==", err)
		return nil, err
	}
	return file, nil
}

// Function to Write an object in a bin file
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0) //(posicion, desde donde) -> (5,0) significa a la posicion 5 desde el inicio del archivo
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

//Se puede optimizar el codigo haciento que WriteObject no retorne nada, en su lugar que imprima el error y listo

// Function to Read an object from a bin file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

func RepGraphizMBR(path string, contenido string, nombre string) error {
	//asegurar la ruta
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error al crear el reporte, path: ", err)
		return err
	}
	// Abrir o crear un archivo para escritura
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return err
	}
	defer file.Close()

	// Escribir en el archivo
	_, err = file.WriteString(contenido)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return err
	}

	rep2 := dir + "/" + nombre + ".png"
	cmd := exec.Command("dot", "-Tpng", path, "-o", rep2)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error al generar el reporte PNG: %v", err)
	}

	return err
}

func Reporte(path string, contenido string) error {
	//asegurar la ruta
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error al crear el reporte, path: ", err)
		return err
	}
	// Abrir o crear un archivo para escritura
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return err
	}
	defer file.Close()

	// Escribir en el archivo
	_, err = file.WriteString(contenido)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return err
	}

	return err
}

// para eliminar en el archivo una particion logica
func DelPartL(size int32) []byte {
	datos := make([]byte, size)
	return datos
}

// probar la escritura de la particion logica
func EscribirPartL(size int32) string {
	cad := strings.Repeat("L", int(size))
	return cad
}
