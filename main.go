package main

import (
	"MIA_P1_200915348/Comandos"
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// MENSAJES DE INICIO
	Ms_inicio := "Bienvenido escriba un comando..."
	Ms_info := "(si desea salir escriba el comando: exit)"
	fmt.Println(Ms_inicio)
	fmt.Println(Ms_info)
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n$: ")
		reader.Scan()
		//entrada := strings.TrimRightFunc(reader.Text(), func(r rune) bool { return r == ' ' })
		entrada := strings.TrimRight(reader.Text(), " ")
		//entrada := execute -path=script.txt
		if strings.ToLower(entrada) != "exit" {
			//fmt.Println("execute -path=script.txt")
			//entrada = "execute -path=script.txt"
			//analizar(strings.ToLower(entrada)) //--------------------OJO REVISAR SI EL TOLOWER NO AFECTA SI SE HACE DESDE AQUI
			analizar(entrada) //--------------------OJO REVISAR SI EL TOLOWER NO AFECTA SI SE HACE DESDE AQUI
		} else {
			fmt.Println("Salida exitosa")
			break
		}
	}
}

func analizar(entrada string) {
	parametros := strings.Split(entrada, " ")
	//Debería pasar todos strings.ToLower(parametros[0])
	if strings.ToLower(parametros[0]) == "execute" {
		if len(parametros) == 2 {
			tmpParametro := strings.Split(parametros[1], "=")
			if strings.ToLower(tmpParametro[0]) == "-path" && len(tmpParametro) == 2 {
				//println("ruta ", tmpParametro[1])
				//abrir el archivo
				archivo, err := os.Open(tmpParametro[1])
				if err != nil {
					fmt.Println("Error al leer el script: ", err)
					return
				}
				defer archivo.Close()
				//creo un lector de bufer para el archivo
				lector := bufio.NewScanner(archivo)
				//leer el archivo linea por linea
				for lector.Scan() {
					//Divido por # para ignorar todo lo que este a la derecha del mismo
					linea := strings.Split(lector.Text(), "#") //lector.Text() retorna la linea leida
					if len(linea[0]) != 0 {
						fmt.Println("*********************************************************************************************")
						fmt.Println("Linea en ejecucion: ", linea[0])
						analizar(linea[0])
					}
				}
			}
		}

	} else if strings.ToLower(parametros[0]) == "mkdisk" {
		//crea un archivo binario que simula un disco con su respectivo MBR
		if len(parametros) > 1 {
			Comandos.Mkdisk(parametros)
		} else {
			fmt.Println("MKDISK ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "fdisk" {
		if len(parametros) > 1 {
			Comandos.Fdisk(parametros)
		} else {
			fmt.Println("FDISK ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "rep" {
		fmt.Println("reportes")
		// Open bin file
		file, err := Herramientas.OpenFile("./Discos/A.dsk")
		if err != nil {
			return
		}
		var TempMBR Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
			return
		}
		// Print object
		Structs.PrintMBR(TempMBR)

		// Close bin file
		defer file.Close()

		//reporte graphviz
		rep := "digraph { node [ shape=none ] nodo0 [ label = < <table border=\"1\"> <tr>"
		rep += Structs.RepGraphviz(TempMBR)
		rep += "</tr> </table> > ] }"
		Herramientas.RepGraphizMBR("Mbr.dot", rep)

	} else if strings.ToLower(parametros[0]) == "pause" {
		fmt.Println("Presione enter para continuar...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

	} else if strings.ToLower(parametros[0]) == "exit" {
		fmt.Println("Salida exitosa")
		os.Exit(0)

	} else {
		fmt.Println("Comando no reconocible")
	}

}
