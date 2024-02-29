package Structs

import (
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

// Metodos de Partition
func GetNamePart(mbr MBR, i int) string {
	nombre := string(mbr.Partitions[i].Name[:])
	posicionNulo := strings.IndexByte(nombre, 0)
	nombre = nombre[:posicionNulo]
	return nombre
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

// Reportes de los Structs
func PrintMBR(data MBR) {
	fmt.Println(fmt.Sprintf("CreationDate: %s, fit: %s, size: %d", string(data.FechaC[:]), string(data.Fit[:]), data.MbrSize))
	for i := 0; i < 4; i++ {
		fmt.Println(fmt.Sprintf("Partition %d: %s, %s, %d, %d, %s", i, string(data.Partitions[i].Name[:]), string(data.Partitions[i].Type[:]), data.Partitions[i].Start, data.Partitions[i].Size, string(data.Partitions[i].Fit[:])))
	}
}

func RepGraphviz(data MBR) string {
	cad := fmt.Sprintf("\n<td bgcolor=\"yellow\"> Signature: %d <br/> Tamaño MBR: %d <br/> Fecha Creado: %s}</td>", data.Id, data.MbrSize, data.FechaC)
	for i := 0; i < 4; i++ {
		if data.Partitions[i].Size > 0 {
			cad += fmt.Sprintf(" \n<td bgcolor=\"green\">Nombre: %s <br/>Tipo: %s  <br/>Tamaño: %d </td>", GetNamePart(data, i), data.Partitions[i].Type, data.Partitions[i].Size)
		} else {
			cad += " \n<td bgcolor=\"red\">Nombre: unknow   <br/>Tipo: unknow <br/>Tamaño: 0 </td>"
		}
	}
	return cad
}
