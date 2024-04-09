package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/HerramientasInodos"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Cat(parametros []string) {
	fmt.Println("Cat")
	var file []string //obligatorio

	//validar que haya un usuario logeado
	if !Structs.UsuarioActual.Status {
		fmt.Println("CAT ERROR: No existe una sesion iniciada")
		return
	}

	//llenar una lista con los archivos a concatenar (si solo fuera uno, es lista de tamaño 1)
	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		tmp := strings.Split(tmp2, "=") //separo para obtener su valor parametro=valor

		//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
		if len(tmp) != 2 {
			fmt.Println("CAT ERROR: Valor desconocido del parametro ", tmp[0])
			return
		}

		if strings.ToLower(tmp[0]) == "file" {
			//si viene solo file
			tmp1 := strings.ReplaceAll(tmp[1], "\"", "")
			file = append(file, tmp1)
		} else {
			//si viene fileN
			comando := strings.Split(strings.ToLower(tmp[0]), "file")
			if comando[0] == "file" {
				_, errId := strconv.Atoi(comando[1])
				if errId != nil {
					fmt.Println("CAT ERROR: No se pudo obtener un numero de fichero")
					return
				}
				//eliminar comillas
				tmp1 := strings.ReplaceAll(tmp[1], "\"", "")
				file = append(file, tmp1)
			} else {
				fmt.Println("CAT ERROR: Parametro desconocido: ", tmp[0])
				return
			}
		}
	}

	//Cargar disco
	//Para buscar el/los inodos
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
		var contenido string
		var fileBlock Structs.Fileblock

		err := Herramientas.ReadObject(disco, &superBloque, int64(mbr.Partitions[part].Start))
		if err != nil {
			fmt.Println("CAT Error. Particion sin formato")
		}

		//buscar el contenido del archivo especificado
		for _, item := range file {
			//buscar el inodo que contiene el archivo buscado
			idInodo := HerramientasInodos.BuscarInodo(0, item, superBloque, disco)
			var inodo Structs.Inode

			//idInodo: solo puede existir archivos desde el inodo 1 en adelante (-1 no existe, 0 es raiz)
			if idInodo > 0 {
				Herramientas.ReadObject(disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
				//recorrer los fileblocks del inodo para obtener toda su informacion
				for _, idBlock := range inodo.I_block {
					if idBlock != -1 {
						Herramientas.ReadObject(disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
						contenido += string(fileBlock.B_content[:])
					}
				}
				contenido += "\n"
			} else {
				fmt.Println("CAT ERROR: No se encontro el archivo ", item)
				return
			}
			fmt.Println(contenido)
		}
	}
}
