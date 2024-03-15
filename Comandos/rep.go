package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"fmt"
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
			path = tmp[1]
		} else if strings.ToLower(tmp[0]) == "id" {
			id = tmp[1]
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
			case "block":
				fmt.Println("reporte block")
			case "bm_inode":
				fmt.Println("reporte bm_inode")
			case "bm_block":
				fmt.Println("reporte bm_block")
			case "tree":
				fmt.Println("reporte tre")
			case "sb":
				fmt.Println("reporte sb")
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
	disk := strings.ToUpper(id[0:1]) //tomar el nombre del disco con case insensitive
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
	disk := strings.ToUpper(id[0:1]) //tomar el nombre del disco con case insensitive
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
