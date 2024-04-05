package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
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

			//Crear el directorio
			ruta := strings.TrimPrefix(path, "/")
			crearCarpeta(0, ruta, superBloque, int64(mbr.Partitions[part].Start), disco, r)
			//crearCarpeta(0, rutaInit, superBloque, iSuperBloque, file, r)
		}
	} else {
		fmt.Println("MKDIR ERROR: falta el parametro path")
		fmt.Println("R ", r)
	}
}

// Metodo para buscar y crear los directorios
func crearCarpeta(inod int32, path string, superBloque Structs.Superblock, iSuperBloque int64, file *os.File, r bool) {
	//NOTA: TODOS LOS INODOS APUNTAN A UN BLOQUE
	ruta := strings.Split(path, "/")
	//la busqueda inicia en el inodo 0
	var Inode0 Structs.Inode
	//Herramientas.ReadObject(file, &Inode0, int64(superBloque.S_inode_start))
	Herramientas.ReadObject(file, &Inode0, int64(superBloque.S_inode_start+(inod*int32(binary.Size(Structs.Inode{})))))

	//Si tamaño de ruta es 1 estoy en la carpeta (inodo) padre
	if len(ruta) == 1 {
		fmt.Println("Creando carpeta")
		//ruta[0] es la carpeta que voy a crear
		//Recorrer los bloques directos del inodo para ver si hay espacio libre
		for i := 0; i < 12; i++ {
			idBloque := Inode0.I_block[i]
			if idBloque != -1 {
				//Existe un folderblock con idBloque que se debe revisar si tiene espacio para la nueva carpeta
				var folderBlock Structs.Folderblock
				Herramientas.ReadObject(file, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
				//Recorrer el bloque para ver si hay espacio
				for j := 2; j < 4; j++ {
					apuntador := folderBlock.B_content[j].B_inodo
					//Hay espacio en el bloque
					if apuntador == -1 {
						//modifico el bloque actual
						copy(folderBlock.B_content[j].B_name[:], ruta[0])
						ino := superBloque.S_first_ino //primer inodo libre
						folderBlock.B_content[j].B_inodo = ino
						//ACTUALIZAR EL FOLDERBLOCK ACTUAL (idBloque) EN EL ARCHIVO
						Herramientas.WriteObject(file, folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

						//creo el nuevo inodo /ruta
						var newInodo Structs.Inode
						newInodo.I_uid = Structs.UsuarioActual.IdUsr
						newInodo.I_gid = Structs.UsuarioActual.IdGrp
						newInodo.I_size = 0 //es carpeta
						//Agrego las fechas
						ahora := time.Now()
						date := ahora.Format("02/01/2006 15:04")
						copy(newInodo.I_atime[:], date)
						copy(newInodo.I_ctime[:], date)
						copy(newInodo.I_mtime[:], date)
						copy(newInodo.I_type[:], "0") //es carpeta
						copy(newInodo.I_mtime[:], "664")

						//apuntadores iniciales
						for i := int32(0); i < 15; i++ {
							newInodo.I_block[i] = -1
						}
						//El apuntador a su primer bloque (el primero disponible)
						block := superBloque.S_first_blo
						newInodo.I_block[0] = block
						//escribo el nuevo inodo (ino)
						Herramientas.WriteObject(file, newInodo, int64(superBloque.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

						//crear el nuevo bloque
						var newFolderBlock Structs.Folderblock
						newFolderBlock.B_content[0].B_inodo = ino //idInodo actual
						copy(newFolderBlock.B_content[0].B_name[:], ".")
						newFolderBlock.B_content[1].B_inodo = folderBlock.B_content[0].B_inodo //el padre es el bloque anterior
						copy(newFolderBlock.B_content[1].B_name[:], "..")
						newFolderBlock.B_content[2].B_inodo = -1
						newFolderBlock.B_content[3].B_inodo = -1
						//escribo el nuevo bloque (block)
						Herramientas.WriteObject(file, newFolderBlock, int64(superBloque.S_block_start+(block*int32(binary.Size(Structs.Folderblock{})))))

						//modifico el superbloque
						superBloque.S_free_inodes_count -= 1
						superBloque.S_free_blocks_count -= 1
						superBloque.S_first_blo += 1
						superBloque.S_first_ino += 1
						//Escribir en el archivo los cambios del superBloque
						Herramientas.WriteObject(file, superBloque, iSuperBloque)

						//escribir el bitmap de bloques (se uso un bloque).
						Herramientas.WriteObject(file, byte(1), int64(superBloque.S_bm_block_start+block))

						//escribir el bitmap de inodos (se uso un inodo).
						Herramientas.WriteObject(file, byte(1), int64(superBloque.S_bm_inode_start+ino))
						//fmt.Println(iSuperBloque)
						return
					}
				} //fin de for de buscar espacio en el bloque actual (existente)
			} else {
				//No hay bloques con espacio disponible
				//modificar el inodo actual (por el nuevo apuntador)
				block := superBloque.S_first_blo //primer bloque libre
				Inode0.I_block[i] = block
				//Escribir los cambios del inodo inicial
				Herramientas.WriteObject(file, &Inode0, int64(superBloque.S_inode_start))

				//creo el primer bloque que va a apuntar a la carpeta
				var newFolderBlock1 Structs.Folderblock
				newFolderBlock1.B_content[0].B_inodo = 0 //estoy en inodo0
				copy(newFolderBlock1.B_content[0].B_name[:], ".")
				newFolderBlock1.B_content[1].B_inodo = 0 //el padre es 0
				copy(newFolderBlock1.B_content[1].B_name[:], "..")
				ino := superBloque.S_first_ino                        //primer inodo libre
				newFolderBlock1.B_content[2].B_inodo = ino            //apuntador al inodo nuevo
				copy(newFolderBlock1.B_content[2].B_name[:], ruta[0]) //nombre del inodo nuevo
				newFolderBlock1.B_content[3].B_inodo = -1
				//escribo el nuevo bloque (block)
				Herramientas.WriteObject(file, newFolderBlock1, int64(superBloque.S_block_start+(block*int32(binary.Size(Structs.Folderblock{})))))

				//creo el nuevo inodo /ruta
				var newInodo Structs.Inode
				newInodo.I_uid = Structs.UsuarioActual.IdUsr
				newInodo.I_gid = Structs.UsuarioActual.IdGrp
				newInodo.I_size = 0 //es carpeta
				//Agrego las fechas
				ahora := time.Now()
				date := ahora.Format("02/01/2006 15:04")
				copy(newInodo.I_atime[:], date)
				copy(newInodo.I_ctime[:], date)
				copy(newInodo.I_mtime[:], date)
				copy(newInodo.I_type[:], "0") //es carpeta
				copy(newInodo.I_mtime[:], "664")

				//apuntadores iniciales
				for i := int32(0); i < 15; i++ {
					newInodo.I_block[i] = -1
				}
				//El apuntador a su primer bloque (el primero disponible)
				block2 := superBloque.S_first_blo + 1
				newInodo.I_block[0] = block2
				//escribo el nuevo inodo (ino) creado en newFolderBlock1
				Herramientas.WriteObject(file, newInodo, int64(superBloque.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

				//crear nuevo bloque del inodo
				var newFolderBlock2 Structs.Folderblock
				newFolderBlock2.B_content[0].B_inodo = ino //idInodo actual
				copy(newFolderBlock2.B_content[0].B_name[:], ".")
				newFolderBlock2.B_content[1].B_inodo = newFolderBlock1.B_content[0].B_inodo //el padre es el bloque anterior
				copy(newFolderBlock2.B_content[1].B_name[:], "..")
				newFolderBlock2.B_content[2].B_inodo = -1
				newFolderBlock2.B_content[3].B_inodo = -1
				//escribo el nuevo bloque
				Herramientas.WriteObject(file, newFolderBlock2, int64(superBloque.S_block_start+(block2*int32(binary.Size(Structs.Folderblock{})))))

				//modifico el superbloque
				superBloque.S_free_inodes_count -= 1
				superBloque.S_free_blocks_count -= 2
				superBloque.S_first_blo += 2
				superBloque.S_first_ino += 1
				Herramientas.WriteObject(file, superBloque, iSuperBloque)

				//escribir el bitmap de bloques (se uso dos bloques: block y block2).
				Herramientas.WriteObject(file, byte(1), int64(superBloque.S_bm_block_start+block))
				Herramientas.WriteObject(file, byte(1), int64(superBloque.S_bm_block_start+block2))

				//escribir el bitmap de inodos (se uso un inodo: ino).
				Herramientas.WriteObject(file, byte(1), int64(superBloque.S_bm_inode_start+ino))
				return
			}
		}
		//Si termino ambos for revisar los apuntadores indirectos
		//i=12 1 bloque indirecto
		//i=13 2 bloques indirectos
		//i=14 3 bloques indirectos
		//Si el inodo es tipo archivo los indirectos apuntaran a fileblocks

	} else {
		//si tengo el metodo que me retorne el ultimo inodo
		//crear las nuevas rutas en lugar de buscar porque ya sabre que
		//estoy en la carpeta padre de las que falten para crear la solicitada
		//que basicamente sería llamar el metodo crear carpeta de nuevo
		//EL PATH TRAE MAS DE UNA CARPETA
		//recorro el inodo buscado en los folderblock la carpeta
		for i := 0; i < 12; i++ {
			idBloque := Inode0.I_block[i]
			if idBloque != -1 {
				var folderBlock Structs.Folderblock
				Herramientas.ReadObject(file, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
				//Recorrer el bloque buscando la carpeta actual
				for j := 2; j < 4; j++ {
					apuntador := folderBlock.B_content[j].B_inodo
					if apuntador != -1 {
						pathActual := Structs.GetB_name(string(folderBlock.B_content[j].B_name[:]))
						if ruta[0] == pathActual {
							buscarInodo(apuntador, ruta[1:], path, superBloque, iSuperBloque, file, r)
							return
						}
					}
				}
			}
		}
		fmt.Println("MKDIR ERROR: No se encontro la carpeta")
	}
}

// Metodo recursivo del tree para buscar bloques
// .              No bloque,   id inodo,  ruta/archivo,       superbloque,            disco
func buscarInodo(idInodo int32, path []string, rutaInit string, superBloque Structs.Superblock, iSuperBloque int64, file *os.File, r bool) {
	var inodo Structs.Inode
	Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

	//recorro el inodo buscando la siguiente carpeta
	for i := 0; i < 12; i++ {
		idBloque := inodo.I_block[i]
		if idBloque != -1 {
			var folderBlock Structs.Folderblock
			Herramientas.ReadObject(file, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
			//Recorrer el bloque buscando la carpeta actua
			for j := 2; j < 4; j++ {
				apuntador := folderBlock.B_content[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GetB_name(string(folderBlock.B_content[j].B_name[:]))
					if path[0] == pathActual {
						buscarInodo(apuntador, path[1:], rutaInit, superBloque, iSuperBloque, file, r)
						return
					}
				}
			}
		}
	}

	//No existe carpeta path[0]
	if len(path) > 1 {
		//crear padre
		if r {
			//se tiene permiso de crear padre
			fmt.Println("Crear la carpeta padre ", path[0])
			//Creo la carpeta padre
			crearCarpeta(idInodo, path[0], superBloque, iSuperBloque, file, r)
			//mando a crear las demas rutas
			//crearCarpeta(0, rutaInit, superBloque, iSuperBloque, file, r)
			fmt.Println("tamaño path ", len(path))
			fmt.Println("path init ", rutaInit)
		} else {
			fmt.Println("Sin permiso de crear carpetas padre")
		}
	} else if len(path) == 1 {
		fmt.Println("Cree carpeta cuando venia solo un path que no existia ", path)
		crearCarpeta(idInodo, path[0], superBloque, iSuperBloque, file, r)
	}
}
