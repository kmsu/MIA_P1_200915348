package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//id -> letra del disco + correlativo particion + 48
//EJ: A148, A248, A348, A448 -> se obtiene en el mount

func Rep(parametros []string) {
	fmt.Println("REP")
	var name string //obligatorio Nombre del reporte a generar
	var path string //obligatorio Nombre que tendrá el reporte
	var id string   //obligatorio sera el del disco o el de la particion
	//var ruta string //opcional para file y ls
	paramC := true //valida que todos los parametros sean correctos

	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		//divido cada parametro entre nombre del parametro y su valor # -size=25 -> -size, 25
		tmp := strings.Split(tmp2, "=")

		//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
		if len(tmp) != 2 {
			fmt.Println("REP Error: Valor desconocido del parametro ", tmp[0])
			paramC = false
			break //para finalizar el ciclo for con el error y no ejecutar lo que haga falta
		}

		if strings.ToLower(tmp[0]) == "name" {
			name = strings.ToLower(tmp[1])
		} else if strings.ToLower(tmp[0]) == "path" {
			// Eliminar comillas
			name = strings.ReplaceAll(tmp[1], "\"", "")
			path = name
		} else if strings.ToLower(tmp[0]) == "id" {
			id = strings.ToUpper(tmp[1]) //Mayusculas para tratarlo como case insensitive
		} else if strings.ToLower(tmp[0]) == "ruta" {
			//ruta = strings.ToLower(tmp[1])
		} else {
			fmt.Println("REP Error: Parametro desconocido: ", tmp[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if paramC {
		if name != "" && id != "" && path != "" {
			switch name {
			case "mbr":
				fmt.Println("reporte mbr")
				mbr(path, id)
			case "disk":
				fmt.Println("reporte disk")
				disk(path, id)
			case "inode":
				fmt.Println("reporte inode")
			case "journaling":
				fmt.Println("reporte journaling")
				journal(path, id)
			case "block":
				fmt.Println("reporte block")
			case "bm_inode":
				fmt.Println("reporte bm_inode")
				bm_inode(path, id)
			case "bm_block":
				fmt.Println("reporte bm_block")
				bm_block(path, id)
			case "tree":
				fmt.Println("reporte tree")
				tree(path, id)
			case "sb":
				fmt.Println("reporte SuperBloque (sb)")
				sb(path, id)
			case "file":
				fmt.Println("reporte file")
			case "ls":
				fmt.Println("reporte ls")
			default:
				fmt.Println("REP Error: Reporte ", name, " desconocido")
			}
		} else {
			fmt.Println("REP Error: Faltan parametros")
		}
	}
}

func mbr(path string, id string) {
	disk := id[0:1] //Nombre del disco
	tmp := strings.Split(path, "/")
	nombre := strings.Split(tmp[len(tmp)-1], ".")[0] //nombre que tendra el reporte. Nombre toma la posicion 0 del split resultante con . (elimina la extension)

	//abrir disco a reportar
	carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
	extension := ".dsk"
	rutaDisco := carpeta + disk + extension

	file, err := Herramientas.OpenFile(rutaDisco)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	//Asegurar que el id exista
	reportar := false
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	//if true { //para probar los reporte hayan o no particiones montadas
	if reportar {
		//reporte graphviz (cad es el contenido del reporte)
		//mbr
		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='SlateBlue' COLSPAN=\"2\"> Reporte MBR </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_tamano </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", mbr.MbrSize)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#AFA1D1'> mbr_fecha_creacion </td> \n  <td bgcolor='#AFA1D1'> %s </td> \n </tr> \n", string(mbr.FechaC[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_disk_signature </td> \n  <td bgcolor='Azure'> %d </td> \n </tr>  \n", mbr.Id)
		cad += Structs.RepGraphviz(mbr, file)
		cad += "</table> > ]\n}"

		//reporte requerido
		carpeta = filepath.Dir(path)
		rutaReporte := "." + carpeta + "/" + nombre + ".dot"

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("REP Error: Id no existe")
	}
}

func disk(path string, id string) {
	disk := id[0:1] //Nombre del disco
	tmp := strings.Split(path, "/")
	nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

	//abrir disco a reportar
	carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
	extension := ".dsk"
	rutaDisco := carpeta + disk + extension

	file, err := Herramientas.OpenFile(rutaDisco)
	if err != nil {
		return
	}

	var TempMBR Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	//Structs.PrintMBR(TempMBR)

	// Close bin file
	defer file.Close()

	//inicia contenido del reporte graphviz del disco
	cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n<tr> \n"
	cad += " <td bgcolor='SlateBlue'  ROWSPAN='3'> MBR </td>\n"
	cad += Structs.RepDiskGraphviz(TempMBR, file)
	cad += "\n</table> > ]\n}"

	//reporte requerido
	carpeta = filepath.Dir(path)
	rutaReporte := "." + carpeta + "/" + nombre + ".dot"

	Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
}

func sb(path string, id string) {
	disk := id[0:1] //Nombre del disco
	tmp := strings.Split(path, "/")
	nombre := strings.Split(tmp[len(tmp)-1], ".")[0] //nombre que tendra el reporte

	//abrir disco a reportar
	carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
	extension := ".dsk"
	rutaDisco := carpeta + disk + extension

	file, err := Herramientas.OpenFile(rutaDisco)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	//Asegurar que el id exista
	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	//if true { //para probar los reporte hayan o no particiones montadas
	if reportar {

		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='darkgreen' COLSPAN=\"2\"> <font color='white'> Reporte SUPERBLOQUE </font> </td> \n </tr> \n"
		cad += Structs.RepSB(mbr.Partitions[part], file)
		cad += "</table> > ]\n}"

		//reporte requerido
		carpeta = filepath.Dir(path)
		rutaReporte := "." + carpeta + "/" + nombre + ".dot"

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("REP Error: Id no existe")
	}
}

func journal(path string, id string) {
	disk := id[0:1] //Nombre del disco
	tmp := strings.Split(path, "/")
	nombre := strings.Split(tmp[len(tmp)-1], ".")[0] //nombre que tendra el reporte

	//abrir disco a reportar
	carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
	extension := ".dsk"
	rutaDisco := carpeta + disk + extension

	file, err := Herramientas.OpenFile(rutaDisco)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	//Asegurar que el id exista
	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	//if true { //para probar los reporte hayan o no particiones montadas
	if reportar {

		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += Structs.RepJournal(mbr.Partitions[part], file)
		cad += "</table> > ]\n}"

		//reporte requerido
		carpeta = filepath.Dir(path)
		rutaReporte := "." + carpeta + "/" + nombre + ".dot"

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("REP Error: Id no existe")
	}
}

func bm_inode(path string, id string) {
	disk := id[0:1] //Nombre del disco
	tmp := strings.Split(path, "/")
	nombre := strings.Split(tmp[len(tmp)-1], ".")[0] //nombre que tendra el reporte

	//abrir disco a reportar
	carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
	extension := ".dsk"
	rutaDisco := carpeta + disk + extension

	file, err := Herramientas.OpenFile(rutaDisco)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	//Asegurar que el id exista
	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if reportar {
		var superBloque Structs.Superblock
		err := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if err != nil {
			fmt.Println("REP Error. Particion sin formato")
			return
		}

		cad := ""
		inicio := superBloque.S_bm_inode_start
		fin := superBloque.S_bm_block_start
		count := 1 //para contar el numero de caracteres por linea (maximo 20)

		//objeto para leer un byte decodificado
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			//cargo el byte (struct de [1]byte) decodificado como las demas estructuras
			Herramientas.ReadObject(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += "0 "
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}

		//reporte requerido
		carpeta = filepath.Dir(path)
		rutaReporte := "." + carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)
	} else {
		fmt.Println("REP Error: Id no existe")
	}
}

func bm_block(path string, id string) {
	disk := id[0:1] //Nombre del disco
	tmp := strings.Split(path, "/")
	nombre := strings.Split(tmp[len(tmp)-1], ".")[0] //nombre que tendra el reporte

	//abrir disco a reportar
	carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
	extension := ".dsk"
	rutaDisco := carpeta + disk + extension

	file, err := Herramientas.OpenFile(rutaDisco)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	//Asegurar que el id exista
	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if reportar {
		var superBloque Structs.Superblock
		err := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if err != nil {
			fmt.Println("REP Error. Particion sin formato")
			return
		}

		cad := ""
		inicio := superBloque.S_bm_block_start
		fin := superBloque.S_inode_start
		count := 1 //para contar el numero de caracteres por linea (maximo 20)

		//objeto para leer un byte decodificado
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			//cargo el byte (struct de [1]byte) decodificado como las demas estructuras
			Herramientas.ReadObject(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += "0 "
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}

		//reporte requerido
		carpeta = filepath.Dir(path)
		rutaReporte := "." + carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)
	} else {
		fmt.Println("REP Error: Id no existe")
	}
}

func tree(path string, id string) {
	disk := id[0:1] //Nombre del disco
	tmp := strings.Split(path, "/")
	nombre := strings.Split(tmp[len(tmp)-1], ".")[0] //nombre que tendra el reporte

	//abrir disco a reportar
	carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
	extension := ".dsk"
	rutaDisco := carpeta + disk + extension

	file, err := Herramientas.OpenFile(rutaDisco)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	//Asegurar que el id exista
	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if reportar {

		var superBloque Structs.Superblock
		err := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if err != nil {
			fmt.Println("REP Error. Particion sin formato")
			return
		}

		var Inode0 Structs.Inode
		Herramientas.ReadObject(file, &Inode0, int64(superBloque.S_inode_start))

		cad := "digraph { \n graph [pad=0.5, nodesep=0.5, ranksep=1] \n node [ shape=plaintext ] \n rankdir=LR \n"

		//reportar el inodo
		cad += "\n Inodo0 [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
		cad += "    <tr> <td bgcolor='skyblue' colspan=\"2\" port='P0'> Inodo 0 </td> </tr> \n"

		for i := 0; i < 12; i++ {
			cad += fmt.Sprintf("    <tr> <td> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode0.I_block[i])
		}
		//Separo los ultimos 3 para marcarlos con color diferente por ser indirectos
		for i := 12; i < 15; i++ {
			cad += fmt.Sprintf("    <tr> <td bgcolor='pink'> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode0.I_block[i])
		}
		cad += "   </table> \n  > \n ]; \n"
		//fin primer inodo

		//llamar bloques
		for i := 0; i < 15; i++ {
			bloque := Inode0.I_block[i]
			if bloque != -1 {
				// No. bloque, tipo Inodo (carpeta/archivo), inodo padre, No port, superbloque, disco
				cad += treeBlock(bloque, string(Inode0.I_type[:]), 0, i+1, superBloque, file)
			}
		}
		//Inode0.I_block[12] -> trae un bloque indirecto antes de un bloque normal
		//Inode0.I_block[13] -> trae dos bloque indirecto antes de un bloque normal
		//Inode0.I_block[14] -> trae tres bloque indirecto antes de un bloque normal
		cad += "\n}"

		//reporte requerido
		carpeta = filepath.Dir(path)
		rutaReporte := "." + carpeta + "/" + nombre + ".dot"

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("REP Error: Id no existe")
	}
}

// Metodo recursivo del tree para buscar bloques
// .              No bloque,   tipo bloque,  inodo padre, No port,          superbloque,            disco
func treeBlock(idBloque int32, tipo string, idPadre int32, p int, superBloque Structs.Superblock, file *os.File) string {
	cad := fmt.Sprintf("\n Bloque%d [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n", idBloque)

	if tipo == "0" {
		//FolderBlock
		var folderBlock Structs.Folderblock
		Herramientas.ReadObject(file, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

		//Reporte del bloque actual
		cad += fmt.Sprintf("    <tr> <td bgcolor='orchid' colspan=\"2\" port='P0'> Bloque %d </td> </tr> \n", idBloque)
		cad += fmt.Sprintf("    <tr> <td> . </td> <td port='P1'> %d </td> </tr> \n", folderBlock.B_content[0].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> .. </td> <td port='P2'> %d </td> </tr> \n", folderBlock.B_content[1].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> %s </td> <td port='P3'> %d </td> </tr> \n", Structs.GetB_name(string(folderBlock.B_content[2].B_name[:])), folderBlock.B_content[2].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> %s </td> <td port='P4'> %d </td> </tr> \n", Structs.GetB_name(string(folderBlock.B_content[3].B_name[:])), folderBlock.B_content[3].B_inodo)
		cad += "   </table> \n  > \n ]; \n"
		//Enlazar inodo padre con bloque actual
		cad += fmt.Sprintf("\n Inodo%d:P%d -> Bloque%d:P0; \n", idPadre, p, idBloque) //p es el port del inodo que apunta al bloque actual
		//recorrero el folderblock para ver si apunta a otros inodos
		for i := 2; i < 4; i++ {
			inodo := folderBlock.B_content[i].B_inodo
			if inodo != -1 {
				//.       inodo hijo, bloque actual, No. port, superbloque, disco
				cad += treeInodo(inodo, idBloque, i+1, superBloque, file)
			}
		}
	} else {
		//Fileblock
		var fileBlock Structs.Fileblock
		Herramientas.ReadObject(file, &fileBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Fileblock{})))))
		//Reporte del bloque actual
		cad += fmt.Sprintf("    <tr> <td bgcolor='#ffff99' port='P0'> Bloque %d </td> </tr> \n", idBloque)
		cad += fmt.Sprintf("    <tr> <td> %s </td> </tr> \n", Structs.GetB_content(string(fileBlock.B_content[:])))
		cad += "   </table> \n  > \n ]; \n"
		//Enlazar inodo padre con bloque actual
		cad += fmt.Sprintf("\n Inodo%d:P%d -> Bloque%d:P0; \n", idPadre, p, idBloque) //p es el port del inodo que apunta al bloque actual
	}

	return cad
}

// Metodo recursivo del tree para buscar inodos
// .            No. inode,     No. bloque,  No. port,          superbloque,            disco
func treeInodo(idInodo int32, idPadre int32, p int, superBloque Structs.Superblock, file *os.File) string {
	//cargar el inodo a reportar
	var Inode Structs.Inode
	Herramientas.ReadObject(file, &Inode, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

	//reportar el inodo
	cad := fmt.Sprintf("\n Inodo%d [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n", idInodo)
	//color segun tipo de inodo
	if string(Inode.I_type[:]) == "0" {
		cad += fmt.Sprintf("    <tr> <td bgcolor='skyblue' colspan=\"2\" port='P0'> Inodo %d </td> </tr> \n", idInodo)
	} else {
		cad += fmt.Sprintf("    <tr> <td bgcolor='#7FC97F' colspan=\"2\" port='P0'> Inodo %d </td> </tr> \n", idInodo)
	}

	//recorrer los apuntadores
	for i := 0; i < 12; i++ {
		cad += fmt.Sprintf("    <tr> <td> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode.I_block[i])
	}
	//Separo los ultimos 3 para marcarlos con color diferente por ser indirectos
	for i := 12; i < 15; i++ {
		cad += fmt.Sprintf("    <tr> <td bgcolor='pink'> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode.I_block[i])
	}
	cad += "   </table> \n  > \n ]; \n"
	//fin inodo

	//Enlazar inodo padre con bloque actual
	cad += fmt.Sprintf("\n Bloque%d:P%d -> Inodo%d:P0; \n", idPadre, p, idInodo) //p es el port del inodo que apunta al bloque actual

	//llamar bloques
	for i := 0; i < 15; i++ {
		bloque := Inode.I_block[i]
		if bloque != -1 {
			//.          No. bloque, tipo Inodo (carpeta/archivo), inodo padre, port, superbloque, disco
			cad += treeBlock(bloque, string(Inode.I_type[:]), idInodo, i+1, superBloque, file)
		}
	}

	return cad
}
