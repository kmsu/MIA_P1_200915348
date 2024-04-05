package Structs

import (
	"MIA_P1_200915348/Herramientas"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// NOTA: Recordar que los atributos de los struct deben iniciar con mayuscula
type MBR struct {
	MbrSize    int32        //mbr_tamano
	FechaC     [16]byte     //mbr_fecha_creacion
	Id         int32        //mbr_dsk_signature (random de forma unica)
	Fit        [1]byte      //dsk_fit
	Partitions [4]Partition //mbr_partitions
}

type Partition struct {
	Status      [1]byte  //part_status
	Type        [1]byte  //part_type
	Fit         [1]byte  //part_fit
	Start       int32    //part_start
	Size        int32    //part_s
	Name        [16]byte //part_name
	Correlative int32    //part_correlative
	Id          [4]byte  //part_id
}

// Setear valores de la particion
func (p *Partition) SetInfo(newType string, fit string, newStart int32, newSize int32, name string, correlativo int32) {
	p.Size = newSize
	p.Start = newStart
	p.Correlative = correlativo
	copy(p.Name[:], name)
	copy(p.Fit[:], fit)
	copy(p.Status[:], "I")
	copy(p.Type[:], newType)
}

// Metodos de Partition
func GetName(nombre string) string {
	posicionNulo := strings.IndexByte(nombre, 0)
	//Si posicionNulo retorna -1 no hay bytes nulos
	if posicionNulo != -1 {
		//guarda la cadena hasta el primer byte nulo (elimina los bytes nulos)
		nombre = nombre[:posicionNulo]
	}
	return nombre
}

func GetId(nombre string) string {
	//si existe id, no contiene bytes nulos
	posicionNulo := strings.IndexByte(nombre, 0)
	//si posicionNulo  no es -1, no existe id.
	if posicionNulo != -1 {
		nombre = "-"
	}
	return nombre
}

func (p *Partition) GetEnd() int32 {
	return p.Start + p.Size
}

type EBR struct {
	Status [1]byte //part_mount (si esta montada)
	Type   [1]byte
	Fit    [1]byte  //part_fit
	Start  int32    //part_start
	Size   int32    //part_s
	Name   [16]byte //part_name
	Next   int32    //part_next
}

func (e *EBR) SetInfo(fit string, newStart int32, newSize int32, name string, newNext int32) {
	e.Size = newSize
	e.Start = newStart
	e.Next = newNext
	copy(e.Name[:], name)
	copy(e.Fit[:], fit)
	copy(e.Status[:], "I")
	copy(e.Type[:], "L")
}

func (e *EBR) GetEnd() int32 {
	return e.Start + e.Size + int32(binary.Size(e))
}

// Reportes de los Structs
func PrintMBR(data MBR) {
	fmt.Println("\n     Disco")
	fmt.Printf("CreationDate: %s, fit: %s, size: %d, id: %d\n", string(data.FechaC[:]), string(data.Fit[:]), data.MbrSize, data.Id)
	for i := 0; i < 4; i++ {
		fmt.Printf("Partition %d: %s, %s, %d, %d, %s, %d\n", i, string(data.Partitions[i].Name[:]), string(data.Partitions[i].Type[:]), data.Partitions[i].Start, data.Partitions[i].Size, string(data.Partitions[i].Fit[:]), data.Partitions[i].Correlative)
	}
}

func PrintEbr(data EBR) {
	fmt.Println("part_status ", string(data.Status[:]))
	fmt.Println("part_type ", string(data.Type[:]))
	fmt.Println("part_fit: ", string(data.Fit[:]))
	fmt.Println("part_start: ", data.Start)
	fmt.Println("part_s ", data.Size)
	fmt.Println("part_name: ", string(data.Name[:]))
	fmt.Println("next_part: ", data.Next)
}

func RepGraphviz(data MBR, disco *os.File) string {
	disponible := int32(0)
	cad := ""
	inicioLibre := int32(binary.Size(data)) //Para ir guardando desde donde hay espacio libre despues de cada particion
	for i := 0; i < 4; i++ {
		if data.Partitions[i].Size > 0 {

			disponible = data.Partitions[i].Start - inicioLibre
			inicioLibre = data.Partitions[i].Start + data.Partitions[i].Size

			//reporta si hay espacio libre antes de la particion
			if disponible > 0 {
				cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#808080' COLSPAN=\"2\"> ESPACIO LIBRE <br/> %d bytes </td> \n </tr> \n", disponible)
			}
			//Reporta el contenido de la particion
			cad += " <tr>\n  <td bgcolor='DeepSkyBlue' COLSPAN=\"2\"> PARTICION </td> \n </tr> \n"
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(data.Partitions[i].Status[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_type </td> \n  <td bgcolor='LightSkyBlue'> %s </td> \n </tr> \n", string(data.Partitions[i].Type[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(data.Partitions[i].Fit[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_start </td> \n  <td bgcolor='LightSkyBlue'> %d </td> \n </tr> \n", data.Partitions[i].Start)
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", data.Partitions[i].Size)
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_name </td> \n  <td bgcolor='LightSkyBlue'> %s </td> \n </tr> \n", GetName(string(data.Partitions[i].Name[:])))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_id </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", GetId(string(data.Partitions[i].Id[:])))
			if string(data.Partitions[i].Type[:]) == "E" {
				cad += repLogicas(data.Partitions[i], disco)
			}
		}
	}

	//si hay espacio despues de la 4ta particion
	disponible = data.MbrSize - inicioLibre
	if disponible > 0 {
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#808080' COLSPAN=\"2\"> ESPACIO LIBRE <br/> %d bytes </td> \n </tr> \n", disponible)
	}

	return cad
}

func repLogicas(particion Partition, disco *os.File) string {
	cad := ""

	var actual EBR
	if err := Herramientas.ReadObject(disco, &actual, int64(particion.Start)); err != nil {
		fmt.Println("REP ERROR: No se encontro un ebr para reportar logicas")
		return ""
	}

	//Primera logica
	if actual.Size != 0 {
		cad += " <tr>\n  <td bgcolor='SteelBlue' COLSPAN=\"2\"> PARTICION LOGICA </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Status[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_next </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Next)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Fit[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_start </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Start)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", actual.Size)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_name </td> \n  <td bgcolor='SkyBlue'> %s </td> \n </tr> \n", GetName(string(actual.Name[:])))
	}

	//resto de logicas
	for actual.Next != -1 {
		if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
			fmt.Println("REP ERROR: fallo al leer particiones logicas")
			return ""
		}
		cad += " <tr>\n  <td bgcolor='SteelBlue' COLSPAN=\"2\"> PARTICION LOGICA </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Status[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_next </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Next)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Fit[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_start </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Start)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", actual.Size)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_name </td> \n  <td bgcolor='SkyBlue'> %s </td> \n </tr> \n", GetName(string(actual.Name[:])))
	}
	return cad
}

func RepDiskGraphviz(data MBR, disco *os.File) string {
	disponible := int32(0)
	cad := ""
	cadLogicas := ""
	cant := 0
	inicioLibre := int32(binary.Size(data)) //Para ir guardando desde donde hay espacio libre despues de cada particion
	for i := 0; i < 4; i++ {
		if data.Partitions[i].Size > 0 {
			disponible = data.Partitions[i].Start - inicioLibre
			inicioLibre = data.Partitions[i].Start + data.Partitions[i].Size
			//reporta si hay espacio libre antes de la particion
			if disponible > 0 {
				porcentaje := float64(disponible) * 100 / float64(data.MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080'  ROWSPAN='3'> ESPACIO LIBRE <br/> %.2f %% </td> \n ", porcentaje)
			}
			porcentaje := float64(data.Partitions[i].Size) * 100 / float64(data.MbrSize)
			if string(data.Partitions[i].Type[:]) == "P" {
				cad += fmt.Sprintf(" <td bgcolor='LightSkyBlue' ROWSPAN='3'> PRIMARIA <br/> %.2f %% </td>\n", porcentaje)
			} else {
				cant, cadLogicas = repLogicasDisk(data.MbrSize, data.Partitions[i], disco)
				cad += fmt.Sprintf(" <td bgcolor='SteelBlue' COLSPAN='%d'> EXTENDIDA </td>\n", cant)
			}
		}
	}

	//si hay espacio despues de la 4ta particion
	disponible = data.MbrSize - inicioLibre
	if disponible > 0 {
		porcentaje := float64(disponible) * 100 / float64(data.MbrSize)
		cad += fmt.Sprintf(" <td bgcolor='#808080'  ROWSPAN='3'> ESPACIO LIBRE <br/> %.2f %% </td> \n", porcentaje)
	}
	cad += "</tr>"    //esta y la siguiente deberian estar en RepDiskGraphiz con la siguiente linea
	cad += cadLogicas //es decir junto con esta
	return cad
}

func repLogicasDisk(MbrSize int32, particion Partition, disco *os.File) (int, string) {
	cant := 0
	cad := "\n\n<tr> \n"
	porcentaje := 0.0

	var actual EBR
	sizeEBR := int32(binary.Size(actual))

	//Cargo el EBR original
	if err := Herramientas.ReadObject(disco, &actual, int64(particion.Start)); err != nil {
		fmt.Println("REP ERROR: Error al leer particiones logicas")
		porcentaje = float64(particion.Size) * 100 / float64(MbrSize)
		return 1, fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
	}

	//Primera logica (Si elimine esta particion el size sera 0 y conserva todos sus demas atributos)
	if actual.Size != 0 {
		//reporto la particion
		porcentaje = float64(actual.Size+sizeEBR) * 100 / float64(MbrSize)
		cad += " <td bgcolor='royalblue3' ROWSPAN='2'> EBR </td>\n"
		cad += fmt.Sprintf(" <td bgcolor='darkgoldenrod2' ROWSPAN='2'> LOGICA <br/> %.2f %% </td>\n", porcentaje)
		cant += 2

		//Verifico si queda espacio libre y lo reporto
		if actual.Next != -1 {
			disponible := actual.Next - actual.GetEnd()
			if disponible > 0 {
				porcentaje = float64(disponible) * 100 / float64(MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
				cant++
			}
		} else {
			disponible := particion.GetEnd() - actual.GetEnd()
			if disponible > 0 {
				porcentaje = float64(disponible) * 100 / float64(MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
				cant++
			}
		}
	} else {
		//Espacio libre
		if actual.Next == -1 {
			//Esta vacia la extendida
			porcentaje = float64(particion.Size) * 100 / float64(MbrSize)
			cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
			cant++
		} else {
			//hay un espacio libre al inicio y existe al menos una particion
			porcentaje = float64(actual.Next-particion.Start) * 100 / float64(MbrSize)
			cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
			cant++
		}
	}

	//Resto de particiones logicas
	for actual.Next != -1 {
		//actual = actual.next (Es decir me paso a la siguiente particion)
		if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
			fmt.Println("REP ERROR: Error al leer particiones logicas")
			porcentaje = float64(particion.Size) * 100 / float64(MbrSize)
			return 1, fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
		}

		//reporto la particion actual
		porcentaje = float64(actual.Size+sizeEBR) * 100 / float64(MbrSize)
		cad += " <td bgcolor='royalblue3' ROWSPAN='2'> EBR </td>\n"
		cad += fmt.Sprintf(" <td bgcolor='darkgoldenrod2' ROWSPAN='2'> LOGICA <br/> %.2f %% </td>\n", porcentaje)
		cant += 2

		//valido si hay espacio disponible en medio o al final
		if actual.Next != -1 {
			disponible := actual.Next - actual.GetEnd()
			if disponible > 0 {
				porcentaje = float64(disponible) * 100 / float64(MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
				cant++
			}
		} else {
			//Es la ultima
			disponible := particion.GetEnd() - actual.GetEnd()
			if disponible > 0 {
				porcentaje = float64(disponible) * 100 / float64(MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
				cant++
			}
		}
	}

	cad += "</tr>\n"
	return cant, cad
}
