package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"fmt"
	"strings"
)

func Rmgrp(parametros []string) {
	fmt.Println("RMGRP")
	var name string //obligatorio

	//Validar parametros
	//quito los espacios en blanco despues de cada parametro
	tmp2 := strings.TrimRight(parametros[1], " ")
	tmp := strings.Split(tmp2, "=") //separo para obtener su valor parametro=valor

	//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
	if len(tmp) != 2 {
		fmt.Println("RMGRP ERROR: Valor desconocido del parametro ", tmp[0])
		return
	}

	if strings.ToLower(tmp[0]) == "name" {
		name = tmp[1]
	} else {
		fmt.Println("RMGRP ERROR: Parametro desconocido: ", tmp[0])
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
					fmt.Println("RMGRP ERROR. Particion sin formato")
					return
				}

				//el users.txt esta en el inodo 1
				var inodo Structs.Inode
				//le agrego una estructura inodo porque busco el inodo 1 (sabemos que aqui esta users.txt)
				Herramientas.ReadObject(disco, &inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

				//leer contenido del Inodo users.txt
				var contenido string
				var fileBlock Structs.Fileblock
				for _, item := range inodo.I_block {
					if item != -1 {
						Herramientas.ReadObject(disco, &fileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
						contenido += string(fileBlock.B_content[:])
					}
				}

				//Separo el contenido wn un arreglo por lineas por lineas
				lineaID := strings.Split(contenido, "\n")

				modificar := false
				for i, registro := range lineaID[:len(lineaID)-1] {
					//lineaID[:len(lineaID)-1] toma el arreglo hasta la penultima posicion (es un salto de linea)
					datos := strings.Split(registro, ",")
					//valido que sea linea de grupo
					if len(datos) == 3 {
						if datos[2] == name {
							//por si ya estaba eliminado
							if datos[0] != "0" {
								modificar = true
								datos[0] = "0"
								mod := datos[0] + "," + datos[1] + "," + datos[2]
								lineaID[i] = mod
							} else {
								fmt.Println("RMGRP ERROR. Grupo ya eliminado")
							}
							break
						}
					}
				}

				//eliminar todos los usuarios del grupo
				if modificar {
					for i, registro := range lineaID[:len(lineaID)-1] {
						//lineaID[:len(lineaID)-1] toma el arreglo hasta la penultima posicion (es un salto de linea)
						datos := strings.Split(registro, ",")
						//valido que sea linea de usuario
						if len(datos) == 5 {
							//valido que se corresponda al mismo grupo
							if datos[2] == name {
								//por si ya estaba eliminado
								if datos[0] != "0" {
									datos[0] = "0"
									mod := datos[0] + "," + datos[1] + "," + datos[2] + "," + datos[3] + "," + datos[4]
									lineaID[i] = mod
								}
							}
						}
					}
				}

				if modificar {
					mod := ""
					for _, reg := range lineaID {
						mod += reg + "\n"
					}

					inicio := 0
					var fin int
					if len(mod) > 64 {
						//si el contenido es mayor a 64 bytes. la primera vez termina en 64
						fin = 64
					} else {
						//termina en el tamaño del contenido. Solo habra un fileblock porque ocupa menos de la capacidad de uno
						fin = len(mod)
					}

					for _, newItem := range inodo.I_block {
						if newItem != -1 {
							//tomo 64 bytes de la cadena o los bytes que queden
							data := mod[inicio:fin]
							//Modifico y guardo el bloque actual
							var newFileBlock Structs.Fileblock
							copy(newFileBlock.B_content[:], []byte(data))
							Herramientas.WriteObject(disco, newFileBlock, int64(superBloque.S_block_start+(newItem*int32(binary.Size(Structs.Fileblock{})))))
							//muevo a los siguientes 64 bytes de la cadena (o los que falten)
							inicio = fin
							calculo := len(mod[fin:]) //tamaño restante de la cadena
							//else if
							if calculo > 64 {
								fin += 64
							} else {
								fin += calculo
							}
							//Si el bloque no esta lleno ya no habra mass fileblocks
						}
					}
				}

			} else {
				fmt.Println("RMGRP ERROR. Error inesperado relacionado el ID de la particion")
			}

		} else {
			fmt.Println("RMGRP ERROR: Usuario con permisos insuficientes")
		}
	} else {
		fmt.Println("RMGRP ERROR: Sesion no iniciada")
	}
}
