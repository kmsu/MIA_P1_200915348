package HerramientasInodos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"os"
	"strings"
)

func BuscarInodo(nombre string, superBloque Structs.Superblock, file *os.File) Structs.Inode {
	var inodo Structs.Inode

	//la busqueda inicia en el inodo 0
	var Inode0 Structs.Inode
	Herramientas.ReadObject(file, &Inode0, int64(superBloque.S_inode_start))

	//Recorrer los bloques del inodo
	for i := 0; i < 15; i++ {
		bloque := Inode0.I_block[i]
		if bloque != -1 {
			//.             No. bloque, tipo Inodo, ruta/archivo buscada, superbloque, disco
			inodo = buscarBlock(bloque, string(Inode0.I_type[:]), nombre, superBloque, file)
			//si el inodo tiene valores es porque ya lo encontro y detine la busqueda
			if inodo.I_size != 0 {
				return inodo
			}
		}
	}
	//separar de i=12 a i=14
	//i=12 1 bloque indirecto
	//i=13 2 bloques indirectos
	//i=14 3 bloques indirectos
	//Si el inodo es tipo archivo los indirectos apuntaran a fileblocks

	return inodo
}

// Metodo recursivo del tree para buscar bloques
// .              No bloque,   tipo inodo,  ruta/archivo,       superbloque,            disco
func buscarBlock(idBloque int32, tipo string, nombre string, superBloque Structs.Superblock, file *os.File) Structs.Inode {
	var inodo Structs.Inode
	ruta := strings.Split(nombre, "/")

	if tipo == "0" {
		// FolderBlock
		var folderBlock Structs.Folderblock
		Herramientas.ReadObject(file, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

		//buscar en el folderblock la siguiente carpeta o el archivo buscado
		for i := 2; i < 4; i++ {
			text := Structs.GetB_name(string(folderBlock.B_content[i].B_name[:]))
			//comparo el atributo B_name (nombre del inodo) con el primer valor de la ruta buscada
			if text == ruta[0] {
				//Si encontro el nombre de carpeta/archivo actual buscado
				if len(ruta) == 1 {
					//ya encontro lo que busca
					idInodo := folderBlock.B_content[i].B_inodo
					if idInodo != -1 {
						Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
					}
				}
				//si el tamaño de ruta es mayor que 1 significa que debe seguir buscando dentro de mas directorios. enviar ruta[1] al buscador de inodo
				break
			}
		}
	}
	return inodo
}
