package Comandos

import (
	"fmt"
	"os"
	"strings"
)

func Rmdisk(parametros []string) {
	fmt.Println("RMDISK")
	//quito los espacios en blano despues de cada parametro
	tmp2 := strings.TrimRight(parametros[1], " ")
	tmp := strings.Split(tmp2, "=")

	if len(tmp) != 2 {
		fmt.Println("FDISK Error: Valor desconocido del parametro ", tmp[0])
		return
	}

	if strings.ToLower(tmp[0]) == "driveletter" {
		letter := strings.ToUpper(tmp[1]) //Debe estar en mayusculas
		carpeta := "./MIA/P1/"            //Ruta (carpeta donde se guardara el disco)
		extension := ".dsk"
		path := carpeta + string(letter) + extension

		//validar si existe el archivo a eliminar
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			fmt.Println("RMDISK Error: El disco ", letter, " no existe")
			return
		}

		//si existe realizar proceso de eliminacion

		//Mensaaje de confirmacion para eliminar el archivo
		fmt.Printf("¿Estás seguro de que quieres eliminar el disco %s? (s/n): ", letter)
		var respuesta string
		fmt.Scanln(&respuesta)

		//Convertir la respuesta a minúsculas para permitir respuestas tanto en mayúsculas como en minúsculas
		respuesta = strings.ToLower(respuesta)

		// Validar la respuesta
		if respuesta == "s" || respuesta == "si" {
			err2 := os.Remove(path)
			if err2 != nil {
				fmt.Println("RMDISK Error: No pudo removerse el disco ")
				return
			}
			fmt.Println("Disco ", letter, "eliminado correctamente:")
		} else {
			fmt.Println("Operación cancelada. El archivo no ha sido eliminado")
		}

	} else {
		fmt.Println("RMDISK Error: Parametro desconocido ", tmp[0])
	}
}
