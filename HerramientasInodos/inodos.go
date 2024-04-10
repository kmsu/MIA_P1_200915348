package HerramientasInodos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func BuscarInodo(idInodo int32, path string, superBloque Structs.Superblock, file *os.File) int32 {
	//Dividir la ruta por cada /
	stepsPath := strings.Split(path, "/")
	//el arreglo vendra [ ,val1, val2] por lo que me corro una posicion
	tmpPath := stepsPath[1:]
	//fmt.Println("Ruta actual ", tmpPath)

	//cargo el inodo a partir del cual voy a buscar
	var Inode0 Structs.Inode
	Herramientas.ReadObject(file, &Inode0, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
	//Recorrer los bloques directos (carpetas/archivos) en la raiz
	var folderBlock Structs.Folderblock
	for i := 0; i < 12; i++ {
		idBloque := Inode0.I_block[i]
		if idBloque != -1 {
			Herramientas.ReadObject(file, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
			//Recorrer el bloque actual buscando la carpeta/archivo en la raiz
			for j := 2; j < 4; j++ {
				//apuntador es el apuntador del bloque al inodo (carpeta/archivo), si existe es distinto a -1
				apuntador := folderBlock.B_content[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GetB_name(string(folderBlock.B_content[j].B_name[:]))
					if tmpPath[0] == pathActual {
						//buscarInodo(apuntador, ruta[1:], path, superBloque, iSuperBloque, file, r)
						if len(tmpPath) > 1 {
							return buscarIrecursivo(apuntador, tmpPath[1:], superBloque.S_inode_start, superBloque.S_block_start, file)
						} else {
							return apuntador
						}
					}
				}
			}
		}
	}
	//agregar busqueda en los apuntadores indirectos
	//i=12 -> simple; i=13 -> doble; i=14 -> triple
	//Si no encontro nada retornar 0 (la raiz)
	return idInodo
}

// Buscar inodo de forma recursiva
func buscarIrecursivo(idInodo int32, path []string, iStart int32, bStart int32, file *os.File) int32 {
	//cargo el inodo actual
	var inodo Structs.Inode
	Herramientas.ReadObject(file, &inodo, int64(iStart+(idInodo*int32(binary.Size(Structs.Inode{})))))

	//Nota: el inodo tiene tipo. No es necesario pero se podria validar que sea carpeta
	//recorro el inodo buscando la siguiente carpeta
	var folderBlock Structs.Folderblock
	for i := 0; i < 12; i++ {
		idBloque := inodo.I_block[i]
		if idBloque != -1 {
			Herramientas.ReadObject(file, &folderBlock, int64(bStart+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
			//Recorrer el bloque buscando la carpeta actua
			for j := 2; j < 4; j++ {
				apuntador := folderBlock.B_content[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GetB_name(string(folderBlock.B_content[j].B_name[:]))
					if path[0] == pathActual {
						if len(path) > 1 {
							//sin este if path[1:] termina en un arreglo de tamaño 0 y retornaria -1
							return buscarIrecursivo(apuntador, path[1:], iStart, bStart, file)
						} else {
							//cuando el arreglo path tiene tamaño 1 esta en la carpeta que busca
							return apuntador
						}
					}
				}
			}
		}
	}
	//agregar busqueda en los apuntadores indirectos
	//i=12 -> simple; i=13 -> doble; i=14 -> triple
	return -1
}

func CrearCarpeta(idInode int32, carpeta string, initSuperBlocque int64, disco *os.File) int32 {
	var superBloque Structs.Superblock
	Herramientas.ReadObject(disco, &superBloque, initSuperBlocque)

	var inodo Structs.Inode
	Herramientas.ReadObject(disco, &inodo, int64(superBloque.S_inode_start+(idInode*int32(binary.Size(Structs.Inode{})))))

	fmt.Println("Creando carpeta ", carpeta)
	//Recorrer los bloques directos del inodo para ver si hay espacio libre
	for i := 0; i < 12; i++ {
		idBloque := inodo.I_block[i]
		if idBloque != -1 {
			//Existe un folderblock con idBloque que se debe revisar si tiene espacio para la nueva carpeta
			var folderBlock Structs.Folderblock
			Herramientas.ReadObject(disco, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

			//Recorrer el bloque para ver si hay espacio
			for j := 2; j < 4; j++ {
				apuntador := folderBlock.B_content[j].B_inodo
				//Hay espacio en el bloque
				if apuntador == -1 {
					//modifico el bloque actual
					copy(folderBlock.B_content[j].B_name[:], carpeta)
					ino := superBloque.S_first_ino //primer inodo libre
					folderBlock.B_content[j].B_inodo = ino
					//ACTUALIZAR EL FOLDERBLOCK ACTUAL (idBloque) EN EL ARCHIVO
					Herramientas.WriteObject(disco, folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

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
					Herramientas.WriteObject(disco, newInodo, int64(superBloque.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

					//crear el nuevo bloque
					var newFolderBlock Structs.Folderblock
					newFolderBlock.B_content[0].B_inodo = ino //idInodo actual
					copy(newFolderBlock.B_content[0].B_name[:], ".")
					newFolderBlock.B_content[1].B_inodo = folderBlock.B_content[0].B_inodo //el padre es el bloque anterior
					copy(newFolderBlock.B_content[1].B_name[:], "..")
					newFolderBlock.B_content[2].B_inodo = -1
					newFolderBlock.B_content[3].B_inodo = -1
					//escribo el nuevo bloque (block)
					Herramientas.WriteObject(disco, newFolderBlock, int64(superBloque.S_block_start+(block*int32(binary.Size(Structs.Folderblock{})))))

					//modifico el superbloque
					superBloque.S_free_inodes_count -= 1
					superBloque.S_free_blocks_count -= 1
					superBloque.S_first_blo += 1
					superBloque.S_first_ino += 1
					//Escribir en el archivo los cambios del superBloque
					Herramientas.WriteObject(disco, superBloque, initSuperBlocque)

					//escribir el bitmap de bloques (se uso un bloque).
					Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_block_start+block))

					//escribir el bitmap de inodos (se uso un inodo).
					Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_inode_start+ino))
					//fmt.Println(iSuperBloque)
					return ino
				}
			} //fin de for de buscar espacio en el bloque actual (existente)
		} else {
			//No hay bloques con espacio disponible
			//modificar el inodo actual (por el nuevo apuntador)
			block := superBloque.S_first_blo //primer bloque libre
			inodo.I_block[i] = block
			//Escribir los cambios del inodo inicial
			Herramientas.WriteObject(disco, &inodo, int64(superBloque.S_inode_start+(idInode*int32(binary.Size(Structs.Inode{})))))

			//creo el primer bloque que va a apuntar a la carpeta
			var newFolderBlock1 Structs.Folderblock
			newFolderBlock1.B_content[0].B_inodo = 0 //estoy en inodo0
			copy(newFolderBlock1.B_content[0].B_name[:], ".")
			newFolderBlock1.B_content[1].B_inodo = 0 //el padre es 0
			copy(newFolderBlock1.B_content[1].B_name[:], "..")
			ino := superBloque.S_first_ino                        //primer inodo libre
			newFolderBlock1.B_content[2].B_inodo = ino            //apuntador al inodo nuevo
			copy(newFolderBlock1.B_content[2].B_name[:], carpeta) //nombre del inodo nuevo
			newFolderBlock1.B_content[3].B_inodo = -1
			//escribo el nuevo bloque (block)
			Herramientas.WriteObject(disco, newFolderBlock1, int64(superBloque.S_block_start+(block*int32(binary.Size(Structs.Folderblock{})))))

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
			Herramientas.WriteObject(disco, newInodo, int64(superBloque.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

			//crear nuevo bloque del inodo
			var newFolderBlock2 Structs.Folderblock
			newFolderBlock2.B_content[0].B_inodo = ino //idInodo actual
			copy(newFolderBlock2.B_content[0].B_name[:], ".")
			newFolderBlock2.B_content[1].B_inodo = newFolderBlock1.B_content[0].B_inodo //el padre es el bloque anterior
			copy(newFolderBlock2.B_content[1].B_name[:], "..")
			newFolderBlock2.B_content[2].B_inodo = -1
			newFolderBlock2.B_content[3].B_inodo = -1
			//escribo el nuevo bloque
			Herramientas.WriteObject(disco, newFolderBlock2, int64(superBloque.S_block_start+(block2*int32(binary.Size(Structs.Folderblock{})))))

			//modifico el superbloque
			superBloque.S_free_inodes_count -= 1
			superBloque.S_free_blocks_count -= 2
			superBloque.S_first_blo += 2
			superBloque.S_first_ino += 1
			Herramientas.WriteObject(disco, superBloque, initSuperBlocque)

			//escribir el bitmap de bloques (se uso dos bloques: block y block2).
			Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_block_start+block))
			Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_block_start+block2))

			//escribir el bitmap de inodos (se uso un inodo: ino).
			Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_inode_start+ino))
			return ino
		}
	}
	return 0
}
