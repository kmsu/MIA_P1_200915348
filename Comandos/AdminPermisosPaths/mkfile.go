package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/HerramientasInodos"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Mkfile(parametros []string) {
	fmt.Println("MKFILE")
	var path string
	var cont string
	size := 0 //opcional, si no viene toma valor 0
	r := false

	//validar que haya un usuario logeado
	if !Structs.UsuarioActual.Status {
		fmt.Println("MKFILE ERROR: No existe una sesion iniciada")
		return
	}

	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		tmp := strings.Split(tmp2, "=") //separo para obtener su valor parametro=valor

		//Capturar valores de los parametros
		if strings.ToLower(tmp[0]) == "path" {
			//Path
			//Si falta el valor del path es de esta forma porque el parametro r tiene tamaño 1
			if len(tmp) != 2 {
				fmt.Println("MKFILE Error: Valor desconocido del parametro ", tmp[0])
				return
			}
			tmp1 := strings.ReplaceAll(tmp[1], "\"", "")
			path = tmp1

			//SIZE
		} else if strings.ToLower(tmp[0]) == "size" {
			//valido que traiga los dos parametros
			if len(tmp) != 2 {
				fmt.Println("MKFILE Error: Valor desconocido del parametro ", tmp[0])
				return
			}

			//convierto a tipo int
			var err error
			size, err = strconv.Atoi(tmp[1]) //se convierte el valor en un entero
			if err != nil {
				fmt.Println("MKFILE Error: Size solo acepta valores enteros. Ingreso: ", tmp[1])
				return
			}

			//valido que sea mayor a 0
			if size < 0 {
				fmt.Println("MKFILE Error: Size solo acepta valores positivos. Ingreso: ", tmp[1])
				return
			}

			//CONT
		} else if strings.ToLower(tmp[0]) == "cont" {
			if len(tmp) != 2 {
				fmt.Println("MKFILE Error: Valor desconocido del parametro ", tmp[0])
				return
			}
			tmp1 := strings.ReplaceAll(tmp[1], "\"", "")
			cont = tmp1

			//R
		} else if strings.ToLower(tmp[0]) == "r" {
			if len(tmp) != 1 {
				fmt.Println("MKDIR Error: Valor desconocido del parametro ", tmp[0])
				return
			}
			r = true

			//ERROR
		} else {
			fmt.Println("MKFILE ERROR: Parametro desconocido: ", tmp[0])
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
			//validar que la particion este formateada
			var superBloque Structs.Superblock
			err := Herramientas.ReadObject(disco, &superBloque, int64(mbr.Partitions[part].Start))
			if err != nil {
				fmt.Println("MKDIR Error. Particion sin formato")
			}

			//Validar que exista la ruta
			stepPath := strings.Split(path, "/")
			finRuta := len(stepPath) - 1 //es el archivo -> stepPath[finRuta] = archivoNuevo.txt
			idInicial := int32(0)
			idActual := int32(0)
			crear := -1
			//No incluye a finRuta, es decir, se queda en el aterior. EJ: Tamaño=5, finRuta=4. El ultimo que evalua es stepPath[3]
			for i, itemPath := range stepPath[1:finRuta] {
				idActual = HerramientasInodos.BuscarInodo(idInicial, "/"+itemPath, superBloque, disco)
				//si el actual y el inicial son iguales significa que no existe la carpeta
				if idInicial != idActual {
					idInicial = idActual
				} else {
					crear = i + 1 //porque estoy iniciando desde 1 e i inicia en 0
					break
				}
			}

			//crear carpetas padre si se tiene permiso
			if crear != -1 {
				if r {
					for _, item := range stepPath[crear:finRuta] {
						idInicial = HerramientasInodos.CrearCarpeta(idInicial, item, int64(mbr.Partitions[part].Start), disco)
						if idInicial == 0 {
							fmt.Println("MKDIR ERROR: No se pudo crear carpeta")
							return
						}
					}
				} else {
					fmt.Println("MKDIR ERROR: Carpeta ", stepPath[crear], " no existe. Sin permiso de crear carpetas padre")
					return
				}

			}

			//verificar que no exista el archivo (recordar que BuscarInodo busca de la forma /nombreBuscar)
			idNuevo := HerramientasInodos.BuscarInodo(idInicial, "/"+stepPath[finRuta], superBloque, disco)
			if idNuevo == idInicial {
				fmt.Println("El archivo no existe")
				//si el parametro cont no viene crear con digitos 0-9
				if cont == "" {
					crearArchivo(idInicial, stepPath[finRuta], size, "", int64(mbr.Partitions[part].Start), disco)
				} //else cargar el contenido de la ruta en cont y mandarlo como parametro en crear archivo
			} else {
				fmt.Println("El archivo ya existe")
			}

		}
	} else {
		fmt.Println("MKFILE ERROR: falta el parametro path")
		fmt.Println("R ", r)
		fmt.Println("Cont ", cont)
	}
}

func crearArchivo(idInodo int32, file string, size int, contenido string, initSuperBloque int64, disco *os.File) {
	//cargar el superBloque actual
	var superB Structs.Superblock
	Herramientas.ReadObject(disco, &superB, initSuperBloque)
	// cargo el inodo de la carpeta que contendra el archivo
	var inodoFile Structs.Inode
	Herramientas.ReadObject(disco, &inodoFile, int64(superB.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

	//recorro el inodo de la carpeta para ver donde guardar el archivo (si hay espacio)
	for i := 0; i < 12; i++ {
		idBloque := inodoFile.I_block[i]
		if idBloque != -1 {
			//Existe un folderblock con idBloque que se debe revisar si tiene espacio para el nuevo archivo
			var folderBlock Structs.Folderblock
			Herramientas.ReadObject(disco, &folderBlock, int64(superB.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

			//Recorrer el bloque para ver si hay espacio y si hay crear el archivo
			for j := 2; j < 4; j++ {
				apuntador := folderBlock.B_content[j].B_inodo
				//Hay espacio en el bloque
				if apuntador == -1 {
					//modifico el bloque actual
					copy(folderBlock.B_content[j].B_name[:], file)
					ino := superB.S_first_ino //primer inodo libre
					folderBlock.B_content[j].B_inodo = ino
					//ACTUALIZAR EL FOLDERBLOCK ACTUAL (idBloque) EN EL ARCHIVO
					fmt.Println("Bloque ", idBloque)
					fmt.Println(Structs.GetB_name(string(folderBlock.B_content[0].B_name[:])), "         | ", folderBlock.B_content[0].B_inodo)
					fmt.Println(Structs.GetB_name(string(folderBlock.B_content[1].B_name[:])), "        | ", folderBlock.B_content[1].B_inodo)
					fmt.Println(Structs.GetB_name(string(folderBlock.B_content[2].B_name[:])), " | ", folderBlock.B_content[2].B_inodo)
					fmt.Println(Structs.GetB_name(string(folderBlock.B_content[3].B_name[:])), "         | ", folderBlock.B_content[3].B_inodo)
					Herramientas.WriteObject(disco, folderBlock, int64(superB.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

					//creo el nuevo inodo archivo
					var newInodo Structs.Inode
					newInodo.I_uid = Structs.UsuarioActual.IdUsr
					newInodo.I_gid = Structs.UsuarioActual.IdGrp
					newInodo.I_size = int32(size) //Size es el tamaño del archivo
					//Agrego las fechas
					ahora := time.Now()
					date := ahora.Format("02/01/2006 15:04")
					copy(newInodo.I_atime[:], date)
					copy(newInodo.I_ctime[:], date)
					copy(newInodo.I_mtime[:], date)
					copy(newInodo.I_type[:], "1") //es archivo
					copy(newInodo.I_mtime[:], "664")

					//apuntadores iniciales
					for i := int32(0); i < 15; i++ {
						newInodo.I_block[i] = -1
					}

					fmt.Println("")
					fmt.Println("bloques libres ", superB.S_free_blocks_count)
					fmt.Println("inodos libres ", superB.S_free_inodes_count)
					fmt.Println("Primer bloque libre ", superB.S_first_blo)
					fmt.Println("Primer inodo libre ", superB.S_first_ino)

					fmt.Println("")
					fmt.Println("Inodo ", ino)

					//Cargar contenido si no viene ningun contenido en el parametro cont
					guardarContenido := ""
					if contenido == "" {
						digito := 0
						for i := 0; i < size; i++ {
							if digito == 10 {
								digito = 0
							}
							guardarContenido += strconv.Itoa(digito)
							digito++
						}
					} //Si contenido trae una ruta entonces tomar lo que contenga el archivo en dicha ruta
					//y la cantidad de caracteres sera el tamaño del archivo (size)

					//El apuntador a su primer bloque (el primero disponible)
					fileblock := superB.S_first_blo

					//division del contenido en los fileblocks de 64 bytes
					inicio := 0
					fin := 0
					sizeContenido := len(guardarContenido)
					if sizeContenido < 64 {
						fin = len(guardarContenido)
					} else {
						fin = 64
					}

					//crear el/los fileblocks con el contenido del archivo0
					for i := int32(0); i < 12; i++ {
						newInodo.I_block[i] = fileblock
						//Guardar la informacion del bloque
						data := guardarContenido[inicio:fin]
						var newFileBlock Structs.Fileblock
						copy(newFileBlock.B_content[:], []byte(data))
						//escribo el nuevo bloque (fileblock)
						Herramientas.WriteObject(disco, newFileBlock, int64(superB.S_block_start+(fileblock*int32(binary.Size(Structs.Fileblock{})))))

						//modifico el superbloque (solo el bloque usado por iteracion)
						superB.S_free_blocks_count -= 1
						superB.S_first_blo += 1

						//escribir el bitmap de bloques (se usa un bloque por iteracion).
						Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_block_start+fileblock))

						fmt.Println("\nNUEVO FILEBLOCK")
						fmt.Println("contenido ", Structs.GetB_content(string(newFileBlock.B_content[:])))

						//validar si queda data que agregar al archivo para continuar con el ciclo o detenerlo
						calculo := len(guardarContenido[fin:])
						if calculo > 64 {
							inicio = fin
							fin += 64
						} else if calculo > 0 {
							inicio = fin
							fin += calculo
						} else {
							//detener el ciclo de creacion de fileblocks
							break
						}
						//Aumento el fileblock
						fileblock++
					}

					//escribo el nuevo inodo (ino)
					fmt.Println("AD ", newInodo.I_block)
					Herramientas.WriteObject(disco, newInodo, int64(superB.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

					//modifico el superbloque por el inodo usado
					superB.S_free_inodes_count -= 1
					superB.S_first_ino += 1
					//Escribir en el archivo los cambios del superBloque
					fmt.Println("bloques libres ", superB.S_free_blocks_count)
					fmt.Println("inodos libres ", superB.S_free_inodes_count)
					fmt.Println("Primer bloque libre ", superB.S_first_blo)
					fmt.Println("Primer inodo libre ", superB.S_first_ino)
					Herramientas.WriteObject(disco, superB, initSuperBloque)

					//escribir el bitmap de inodos (se uso un inodo).
					Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_inode_start+ino))

					return
				}
			} //fin de for de buscar espacio en el bloque actual (existente)
		} else {
			fmt.Println("Crear nuevo bloque")

			//No hay bloques con espacio disponible
			//modificar el inodo actual (por el nuevo apuntador)
			block := superB.S_first_blo //primer bloque libre
			inodoFile.I_block[i] = block
			//Escribir los cambios del inodo inicial
			Herramientas.WriteObject(disco, &inodoFile, int64(superB.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

			//cargo el primer bloque del inodo actual para tomar los datos de actual y padre (son los mismos para el nuevo)
			var folderBlock Structs.Folderblock
			bloque := inodoFile.I_block[0] //cargo el primer folderblock para obtener los datos del actual y su padre
			Herramientas.ReadObject(disco, &folderBlock, int64(superB.S_block_start+(bloque*int32(binary.Size(Structs.Folderblock{})))))

			//creo el primer bloque que va a apuntar al nuevo archivo
			var newFolderBlock1 Structs.Folderblock
			newFolderBlock1.B_content[0].B_inodo = folderBlock.B_content[0].B_inodo //actual
			copy(newFolderBlock1.B_content[0].B_name[:], ".")
			newFolderBlock1.B_content[1].B_inodo = folderBlock.B_content[1].B_inodo //padre
			copy(newFolderBlock1.B_content[1].B_name[:], "..")
			ino := superB.S_first_ino                          //primer inodo libre
			newFolderBlock1.B_content[2].B_inodo = ino         //apuntador al inodo nuevo
			copy(newFolderBlock1.B_content[2].B_name[:], file) //nombre del inodo nuevo
			newFolderBlock1.B_content[3].B_inodo = -1
			//escribo el nuevo bloque (block)
			Herramientas.WriteObject(disco, newFolderBlock1, int64(superB.S_block_start+(block*int32(binary.Size(Structs.Folderblock{})))))

			//escribir el bitmap de bloques
			Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_block_start+block))

			//modifico el superbloque porque mas adelante lo necesito con estos cambios
			superB.S_first_blo += 1
			superB.S_free_blocks_count -= 1

			//creo el nuevo inodo archivo
			var newInodo Structs.Inode
			newInodo.I_uid = Structs.UsuarioActual.IdUsr
			newInodo.I_gid = Structs.UsuarioActual.IdGrp
			newInodo.I_size = int32(size) //Size es el tamaño del archivo
			//Agrego las fechas
			ahora := time.Now()
			date := ahora.Format("02/01/2006 15:04")
			copy(newInodo.I_atime[:], date)
			copy(newInodo.I_ctime[:], date)
			copy(newInodo.I_mtime[:], date)
			copy(newInodo.I_type[:], "1") //es archivo
			copy(newInodo.I_mtime[:], "664")

			//apuntadores iniciales
			for i := int32(0); i < 15; i++ {
				newInodo.I_block[i] = -1
			}

			fmt.Println("")
			fmt.Println("bloques libres ", superB.S_free_blocks_count)
			fmt.Println("inodos libres ", superB.S_free_inodes_count)
			fmt.Println("Primer bloque libre ", superB.S_first_blo)
			fmt.Println("Primer inodo libre ", superB.S_first_ino)

			fmt.Println("")
			fmt.Println("Inodo ", ino)

			//Cargar contenido si no viene ningun contenido en el parametro cont
			guardarContenido := ""
			if contenido == "" {
				digito := 0
				for i := 0; i < size; i++ {
					if digito == 10 {
						digito = 0
					}
					guardarContenido += strconv.Itoa(digito)
					digito++
				}
			} //Si contenido trae una ruta entonces tomar lo que contenga el archivo en dicha ruta
			//y la cantidad de caracteres sera el tamaño del archivo (size)

			//El apuntador a su primer bloque (el primero disponible)
			fileblock := superB.S_first_blo

			//division del contenido en los fileblocks de 64 bytes
			inicio := 0
			fin := 0
			sizeContenido := len(guardarContenido)
			if sizeContenido < 64 {
				fin = len(guardarContenido)
			} else {
				fin = 64
			}

			//crear el/los fileblocks con el contenido del archivo0
			for i := int32(0); i < 12; i++ {
				newInodo.I_block[i] = fileblock
				//Guardar la informacion del bloque
				data := guardarContenido[inicio:fin]
				var newFileBlock Structs.Fileblock
				copy(newFileBlock.B_content[:], []byte(data))
				//escribo el nuevo bloque (fileblock)
				Herramientas.WriteObject(disco, newFileBlock, int64(superB.S_block_start+(fileblock*int32(binary.Size(Structs.Fileblock{})))))

				//modifico el superbloque (solo el bloque usado por iteracion)
				superB.S_free_blocks_count -= 1
				superB.S_first_blo += 1

				//escribir el bitmap de bloques (se usa un bloque por iteracion).
				Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_block_start+fileblock))

				fmt.Println("\nNUEVO FILEBLOCK")
				fmt.Println("contenido ", Structs.GetB_content(string(newFileBlock.B_content[:])))

				//validar si queda data que agregar al archivo para continuar con el ciclo o detenerlo
				calculo := len(guardarContenido[fin:])
				if calculo > 64 {
					inicio = fin
					fin += 64
				} else if calculo > 0 {
					inicio = fin
					fin += calculo
				} else {
					//detener el ciclo de creacion de fileblocks
					break
				}
				//Aumento el fileblock
				fileblock++
			}

			//escribo el nuevo inodo (ino)
			fmt.Println("AD ", newInodo.I_block)
			Herramientas.WriteObject(disco, newInodo, int64(superB.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

			//modifico el superbloque por el inodo usado
			superB.S_free_inodes_count -= 1
			superB.S_first_ino += 1
			//Escribir en el archivo los cambios del superBloque
			fmt.Println("bloques libres ", superB.S_free_blocks_count)
			fmt.Println("inodos libres ", superB.S_free_inodes_count)
			fmt.Println("Primer bloque libre ", superB.S_first_blo)
			fmt.Println("Primer inodo libre ", superB.S_first_ino)
			Herramientas.WriteObject(disco, superB, initSuperBloque)

			//escribir el bitmap de inodos (se uso un inodo).
			Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_inode_start+ino))

			return

		}
	}

}
