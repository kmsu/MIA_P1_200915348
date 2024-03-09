package Structs

import (
	"encoding/binary"
	"fmt"
	"strings"
)

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
	nombre = nombre[:posicionNulo]
	return nombre
}

func GetId(nombre string) string {
	var id string
	posicionNulo := strings.IndexByte(nombre, 0)
	nombre = nombre[:posicionNulo]
	if nombre == "" {
		id = "-"
	} else {
		id = nombre
	}
	return id
}

func (p *Partition) GetEnd() int32 {
	return p.Start + p.Size
}

/*
func GetEnd(part Partition) int32 {
	return part.Start + part.Size
}*/

type EBR struct {
	Status [1]byte //part_mount (si esta montada)
	Type   [1]byte
	Fit    [1]byte  //part_fit
	Start  int32    //part_start
	Size   int32    //part_s
	Name   [16]byte //part_name
	Next   int32    //part_next
}

// Reportes de los Structs
func PrintMBR(data MBR) {
	fmt.Println("\n     Disco")
	fmt.Printf("CreationDate: %s, fit: %s, size: %d, id: %d\n", string(data.FechaC[:]), string(data.Fit[:]), data.MbrSize, data.Id)
	for i := 0; i < 4; i++ {
		fmt.Printf("Partition %d: %s, %s, %d, %d, %s, %d\n", i, string(data.Partitions[i].Name[:]), string(data.Partitions[i].Type[:]), data.Partitions[i].Start, data.Partitions[i].Size, string(data.Partitions[i].Fit[:]), data.Partitions[i].Correlative)
	}
}

func RepGraphviz(data MBR) string {
	disponible := int32(0)
	cad := ""
	inicioLibre := int32(binary.Size(data)) //Para ir guardando desde donde hay espacio libre despues de cada particion
	for i := 0; i < 4; i++ {
		if data.Partitions[i].Size > 0 {

			disponible = data.Partitions[i].Start - inicioLibre
			inicioLibre = data.Partitions[i].Start + data.Partitions[i].Size

			//reporta si hay espacio libre antes de la particion
			if disponible > 0 {
				cad += fmt.Sprintf("<tr> <td bgcolor='#808080' COLSPAN=\"2\"> ESPACIO LIBRE <br/> %d bytes </td> </tr> ", disponible)
			}
			//Reporta el contenido de la particion
			cad += " <tr>\n  <td bgcolor='DeepSkyBlue' COLSPAN=\"2\"> PARTICION </td> \n </tr> \n"
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n <td bgcolor='Azure'> %s </td> \n </tr> \n", string(data.Partitions[i].Status[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_type</td> \n  <td bgcolor='LightSkyBlue'> %s </td> \n </tr> \n", string(data.Partitions[i].Type[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(data.Partitions[i].Fit[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_start</td> \n  <td bgcolor='LightSkyBlue'> %d </td> \n </tr> \n", data.Partitions[i].Start)
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", data.Partitions[i].Size)
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_name </td> \n  <td bgcolor='LightSkyBlue'> %s </td> \n </tr> \n", GetName(string(data.Partitions[i].Name[:])))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_id </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", GetId(string(data.Partitions[i].Id[:])))
		}
	}

	//si hay espacio despues de la 4ta particion
	disponible = data.MbrSize - inicioLibre
	if disponible > 0 {
		cad += fmt.Sprintf("<tr> <td bgcolor='#808080' COLSPAN=\"2\"> ESPACIO LIBRE <br/> %d bytes </td> </tr> ", disponible)
	}

	return cad
}

func RepDiskGraphviz(data MBR) string {
	disponible := int32(0)
	cad := ""
	cadLogicas := ""
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
				//cant, cadLogicas = metodo(path)
				cad += " <td bgcolor='SteelBlue' COLSPAN='1'> EXTENDIDA </td>\n"
				cadLogicas = fmt.Sprintf("\n\n<tr> \n <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td> \n</tr>\n", porcentaje)
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
