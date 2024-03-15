package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Mount(parametros []string) {
	fmt.Println("MOUNT")
	//PARAMETROS: -driveletter -name
	var letter string //obligatorio (es el "path", es una letra nombre del disco donde esta la particion, path ya esta fijado)
	var name string
	paramC := true //Para validar que los parametros cumplen con los requisitos

	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		tmp := strings.Split(tmp2, "=")

		//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
		if len(tmp) != 2 {
			fmt.Println("MOUNT Error: Valor desconocido del parametro ", tmp[0])
			paramC = false
			return
		}

		//PATH
		if strings.ToLower(tmp[0]) == "driveletter" {
			//homonimo al path
			letter = strings.ToUpper(tmp[1]) //Debe estar en mayusculas
			//Se valida si existe el disco ingresado
			carpeta := "./MIA/P1/" //Ruta (carpeta donde se guardara el disco)
			extension := ".dsk"
			path := carpeta + string(letter) + extension
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("MOUNT Error: El disco ", letter, " no existe")
				paramC = false
				break // Terminar el bucle porque encontramos un nombre único
			}

			//NAME
		} else if strings.ToLower(tmp[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(tmp[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)

			//ERROR EN LOS PARAMETROS LEIDOS
		} else {
			fmt.Println("MOUNT Error: Parametro desconocido ", tmp[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	//Si todos los parametros son correctos
	if paramC {
		if letter != "" && name != "" {
			// Abrir y cargar el disco
			filepath := "./MIA/P1/" + letter + ".dsk"
			disco, err := Herramientas.OpenFile(filepath)
			if err != nil {
				fmt.Println("MOUNT Error: No se pudo leer el disco")
				return
			}

			//Se crea un mbr para cargar el mbr del disco
			var mbr Structs.MBR
			//Guardo el mbr leido
			if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
				return
			}

			montar := true //usar si se van a montar logicas
			reportar := false
			for i := 0; i < 4; i++ {
				nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
				if nombre == name {
					montar = false
					if string(mbr.Partitions[i].Status[:]) != "A" {
						if string(mbr.Partitions[i].Type[:]) != "E" {
							//id = letter + correlativo + 48 (48 -> ultimos dos digitos de 200915348)
							id := strings.ToUpper(letter) + strconv.Itoa(i+1) + "48"

							//modificar la particion que se va a montar
							copy(mbr.Partitions[i].Status[:], "A")
							copy(mbr.Partitions[i].Id[:], id)

							//sobreescribir el mbr para guardar los cambios
							if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
								return
							}
							reportar = true
							fmt.Println("Particion con nombre ", name, " montada correctamente")
						} else {
							fmt.Println("MOUNT Error. No se puede montar una particion extendida")
						}
					} else {
						fmt.Println("MOUNT Error. Particion ya montada")
					}
					break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
				}
			}

			if montar {
				fmt.Println("MOUNT Error. No se pudo montar la particion ", name)
				fmt.Println("MOUNT Error. No se encontro la particion")
			}

			//MOSTRAR PARTICIONES MONTADAS
			if reportar {
				fmt.Println("\nLISTA DE PARTICIONES MONTADAS\n ")
				for i := 0; i < 4; i++ {
					estado := string(mbr.Partitions[i].Status[:])
					if estado == "A" {
						fmt.Printf("Partition %d: name: %s, status: %s, id: %s, tipo: %s, start: %d, size: %d, fit: %s, correlativo: %d\n", i, string(mbr.Partitions[i].Name[:]), string(mbr.Partitions[i].Status[:]), string(mbr.Partitions[i].Id[:]), string(mbr.Partitions[i].Type[:]), mbr.Partitions[i].Start, mbr.Partitions[i].Size, string(mbr.Partitions[i].Fit[:]), mbr.Partitions[i].Correlative)
					}
				}
			}

		} else {
			fmt.Println("MOUNT Error. No se encontro parametro letter y/o name")
		}
	}
}
