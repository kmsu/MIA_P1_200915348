package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Mkgrp(parametros []string) {
	fmt.Println("Mkgrp")
	var name string //obligatorio

	//validar parametros
	//quito los espacios en blanco despues de cada parametro
	tmp2 := strings.TrimRight(parametros[1], " ")
	tmp := strings.Split(tmp2, "=") //separo para obtener su valor parametro=valor

	//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
	if len(tmp) != 2 {
		fmt.Println("MKGRP ERROR: Valor desconocido del parametro ", tmp[0])
		return
	}

	if strings.ToLower(tmp[0]) == "name" {
		name = tmp[1]

		//validar maximo 10 caracteres
		if len(name) > 10 {
			fmt.Println("MKGRP ERROR: name debe tener maximo 10 caracteres")
			return
		}
	} else {
		fmt.Println("MKGRP ERROR: Parametro desconocido: ", tmp[0])
		return
	}

	//LOGICA DEL COMANDO
	usuarioActual := Structs.UsuarioActual
	if usuarioActual.Status {
		if usuarioActual.Nombre == "root" {
			//CARGAR EL DISCO DONDE PODRÍA ESTAR LA PARTICION
			disk := usuarioActual.Id[0:1] // cargar el disco
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

			//buscar particion para cargar estructuras
			ejecutar := false
			index := -1
			for i := 0; i < 4; i++ {
				identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
				if identificador == usuarioActual.Id {
					ejecutar = true
					index = i
					break //para que ya no siga recorriendo si ya encontro la particion
				}
			}

			if ejecutar {
				//cargo el superbloque de la particion
				var superBloque Structs.Superblock
				err = Herramientas.ReadObject(disco, &superBloque, int64(mbr.Partitions[index].Start))
				if err != nil {
					fmt.Println("MKGRP ERROR. Particion sin formato")
					return
				}

				//el users.txt esta en el inodo 1
				var inodo Structs.Inode
				//le agrego una estructura inodo porque busco el inodo 1 (sabemos que aqui esta users.txt)
				Herramientas.ReadObject(disco, &inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

				//leer contenido del Inodo users.txt
				var contenido string
				var fileBlock Structs.Fileblock
				var idFB int32 //id/numero de ultimo fileblock para trabajar sobre ese
				for _, item := range inodo.I_block {
					if item != -1 {
						Herramientas.ReadObject(disco, &fileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
						contenido += string(fileBlock.B_content[:])
						idFB = item
					}
				}
				//UID, Tipo, Grupo, Usuario, contraseña

				//Buscar el urltimo ID distinto de 0
				lineaID := strings.Split(contenido, "\n")

				//Verificar si el grupo ya existe
				for _, registro := range lineaID[:len(lineaID)-1] {
					datos := strings.Split(registro, ",")
					if len(datos) == 3 {
						if datos[2] == name {
							fmt.Println("MKGRP ERROR: El grupo ya existe")
							return
						}
					}
				}

				//Buscar el ultimo ID activo desde el ultimo hasta el primero (ignorando los eliminado (0))
				//desde -2 porque siempre se crea un salto de linea al final generando una linea vacia al final del arreglo
				id := -1        //para guardar el nuevo ID
				var errId error //para la conversion a numero del ID
				for i := len(lineaID) - 2; i >= 0; i-- {
					registro := strings.Split(lineaID[i], ",")
					//valido que sea un grupo
					if registro[1] == "G" {
						//valido que el id sea distinto a 0 (eliminado)
						if registro[0] != "0" {
							//convierto el id en numero para sumarle 1 y crear el nuevo id
							id, errId = strconv.Atoi(registro[0])
							if errId != nil {
								fmt.Println("MKGRP ERROR: No se pudo obtener un nuevo id para el nuevo grupo")
								return
							}
							id++
							break
						}
					}
				}

				//valido que se haya encontrado un nuevo id
				if id != -1 {
					contenidoActual := string(fileBlock.B_content[:])
					posicionNulo := strings.IndexByte(contenidoActual, 0)
					data := fmt.Sprintf("%d,G,%s\n", id, name)
					//Aseguro que haya al menos un byte libre
					if posicionNulo != -1 {
						libre := 64 - (posicionNulo + len(data))
						if libre > 0 {
							copy(fileBlock.B_content[posicionNulo:], []byte(data))
							//Escribir el fileblock con espacio libre
							Herramientas.WriteObject(disco, fileBlock, int64(superBloque.S_block_start+(idFB*int32(binary.Size(Structs.Fileblock{})))))
						} else {
							//Si es 0 (quedó exacta), entra aqui y crea un bloque vacío que podrá usarse para el proximo registro

							data1 := data[:len(data)+libre]
							//Ingreso lo que cabe en el bloque actual
							copy(fileBlock.B_content[posicionNulo:], []byte(data1))
							Herramientas.WriteObject(disco, fileBlock, int64(superBloque.S_block_start+(idFB*int32(binary.Size(Structs.Fileblock{})))))

							//Creo otro fileblock para el resto de la informacion
							guardoInfo := true
							for i, item := range inodo.I_block {
								//i es el indice en el arreglo inodo.Iblock
								if item == -1 {
									guardoInfo = false
									//agrego el apuntador del bloque al inodo
									inodo.I_block[i] = superBloque.S_first_blo
									//actualizo el superbloque
									superBloque.S_free_blocks_count -= 1
									superBloque.S_first_blo += 1
									data2 := data[len(data)+libre:]
									//crear nuevo fileblock
									var newFileBlock Structs.Fileblock
									copy(newFileBlock.B_content[:], []byte(data2))

									//escribir las estructuras para guardar los cambios
									// Escribir el superbloque
									Herramientas.WriteObject(disco, superBloque, int64(mbr.Partitions[index].Start))

									//escribir el bitmap de bloques (se uso un bloque). inodo.I_block[i] contiene el numero de bloque que se uso
									Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_block_start+inodo.I_block[i]))

									//escribir inodes (es el inodo 1, porque es donde esta users.txt)
									Herramientas.WriteObject(disco, inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

									//Escribir bloques
									Herramientas.WriteObject(disco, newFileBlock, int64(superBloque.S_block_start+(inodo.I_block[i]*int32(binary.Size(Structs.Fileblock{})))))
									break
								}
							}

							if guardoInfo {
								fmt.Println("MKGRP ERROR: Espacio insuficiente para nuevo registro")
							}
						}
					}
				}

			} else {
				fmt.Println("MKGRP ERROR. Error inesperado relacionado el ID de la particion")
			}

		} else {
			fmt.Println("MKGRP ERROR: Usuario con permisos insuficientes")
		}
	} else {
		fmt.Println("MKGRP ERROR: Sesion no iniciada")
	}
}
