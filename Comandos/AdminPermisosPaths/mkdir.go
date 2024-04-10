package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/HerramientasInodos"
	"MIA_P1_200915348/Structs"
	"fmt"
	"strings"
)

func Mkdir(parametros []string) {
	fmt.Println("MKDIR")
	var path string
	r := false

	//validar que haya un usuario logeado
	if !Structs.UsuarioActual.Status {
		fmt.Println("MKDIR ERROR: No existe una sesion iniciada")
		return
	}

	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		tmp := strings.Split(tmp2, "=") //separo para obtener su valor parametro=valor

		//Capturar valores de los parametros
		if strings.ToLower(tmp[0]) == "path" {
			//Path
			//Si falta el valor del path
			if len(tmp) != 2 {
				fmt.Println("MKDIR Error: Valor desconocido del parametro ", tmp[0])
				return
			}
			tmp1 := strings.ReplaceAll(tmp[1], "\"", "")
			path = tmp1

			//R
		} else if strings.ToLower(tmp[0]) == "r" {
			if len(tmp) != 1 {
				fmt.Println("MKDIR Error: Valor desconocido del parametro ", tmp[0])
				return
			}
			r = true

			//ERROR
		} else {
			fmt.Println("MKDIR ERROR: Parametro desconocido: ", tmp[0])
			return
		}
	}

	if path != "" {
		//CARGA DE INFORMACION NECESARIA PARA EL COMANDO
		//Cargar disco
		id := Structs.UsuarioActual.Id
		disk := id[0:1] //Nombre del disco
		//abrir disco a reportar
		carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
		extension := ".dsk"
		rutaDisco := carpeta + disk + extension

		disco, err := Herramientas.OpenFile(rutaDisco)
		if err != nil {
			return
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
			return
		}

		// Close bin file
		defer disco.Close()

		//buscar particion con id actual
		buscar := false
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				buscar = true
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		if buscar {
			var superBloque Structs.Superblock

			err := Herramientas.ReadObject(disco, &superBloque, int64(mbr.Partitions[part].Start))
			if err != nil {
				fmt.Println("MKDIR Error. Particion sin formato")
			}

			//No existe buscar padre
			stepPath := strings.Split(path, "/")
			buscar := ""
			idInicial := int32(0)
			idActual := int32(0)
			crear := -1
			for i, itemPath := range stepPath[1:] {
				buscar = "/" + itemPath
				idActual = HerramientasInodos.BuscarInodo(idInicial, buscar, superBloque, disco)
				if idInicial != idActual {
					idInicial = idActual
				} else {
					crear = i + 1 //porque estoy iniciando desde 1 e i inicia en 0
					break
				}
			}

			if crear != -1 {
				if crear == len(stepPath)-1 {
					fmt.Println("Crear la carpeta ", stepPath[crear])
					HerramientasInodos.CrearCarpeta(idInicial, stepPath[crear], int64(mbr.Partitions[part].Start), disco)
				} else {
					if r {
						for _, item := range stepPath[crear:] {
							idInicial = HerramientasInodos.CrearCarpeta(idInicial, item, int64(mbr.Partitions[part].Start), disco)
							if idInicial == 0 {
								fmt.Println("MKDIR ERROR: No se pudo crear carpeta")
								return
							}
						}
					} else {
						fmt.Println("MKDIR ERROR: Sin permiso de crear carpetas padre")
					}
				}
			} else {
				fmt.Println("MKDIR ERROR: Carpeta ya existe")
			}
		}
	} else {
		fmt.Println("MKDIR ERROR: falta el parametro path")
		fmt.Println("R ", r)
	}
}
