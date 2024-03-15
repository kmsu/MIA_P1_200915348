package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"fmt"
	"strings"
)

func Unount(parametros []string) {
	fmt.Println("UNMOUNT")
	//PARAMETROS: -driveletter -name
	var id string  //obligatorio
	paramC := true //Para validar que los parametros cumplen con los requisitos

	//quito los espacios en blano despues de cada parametro
	tmp2 := strings.TrimRight(parametros[1], " ")
	tmp := strings.Split(tmp2, "=")

	//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
	if len(tmp) != 2 {
		fmt.Println("UNMOUNT Error: Valor desconocido del parametro ", tmp[0])
		paramC = false
		return
	}

	if strings.ToLower(tmp[0]) == "id" {
		id = tmp[1]
	} else {
		fmt.Println("UNMOUNT Error: Parametro desconocido ", tmp[0])
		paramC = false
	}

	//Si todos los parametros son correctos
	if paramC {
		//tomar el nombre del disco con case insensitive
		disk := strings.ToUpper(id[0:1])
		//generar ruta del disco que puede contener el id
		carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
		extension := ".dsk"
		rutaDisco := carpeta + disk + extension

		//abrir el disco que podría contener el id
		disco, err := Herramientas.OpenFile(rutaDisco)
		if err != nil {
			return
		}

		//cargar el mbr
		var mbr Structs.MBR
		if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
			return
		}

		//cerrar el archivo del disco
		defer disco.Close()

		desmontar := true //usar si se van a montar logicas
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				desmontar = false
				name := Structs.GetName(string(mbr.Partitions[i].Name[:]))
				var unmount Structs.Partition

				//Eliminar el id usando el id de la variable unmount
				mbr.Partitions[i].Id = unmount.Id
				copy(mbr.Partitions[i].Status[:], "I")

				//sobreescribir el mbr para guardar los cambios
				if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
					return
				}
				fmt.Println("Particion con nombre ", name, " desmontada correctamente")
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		if desmontar {
			fmt.Println("UNMOUNT Error. No se pudo desmontar la particion con id ", id)
			fmt.Println("UNMOUNT Error. No existe el id")
		} else {
			fmt.Println("\nLISTA DE PARTICIONES MONTADAS (muestra que se desmonto la particion)\n ")
			for i := 0; i < 4; i++ {
				estado := string(mbr.Partitions[i].Status[:])
				if estado == "A" {
					fmt.Printf("Partition %d: name: %s, status: %s, id: %s, tipo: %s, start: %d, size: %d, fit: %s, correlativo: %d\n", i, string(mbr.Partitions[i].Name[:]), string(mbr.Partitions[i].Status[:]), string(mbr.Partitions[i].Id[:]), string(mbr.Partitions[i].Type[:]), mbr.Partitions[i].Start, mbr.Partitions[i].Size, string(mbr.Partitions[i].Fit[:]), mbr.Partitions[i].Correlative)
				}
			}
		}
	}
}
