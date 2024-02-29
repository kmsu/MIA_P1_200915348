package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Fdisk(parametros []string) {
	fmt.Println("FDISK")
	//PARAMETROS: -type -size -path -name -unit
	var size int          //obligatorio si es creacion
	var letter string     //obligatorio (es el "path", es una letra nombre de la particion, path ya esta fijado)
	var name string       //obligatorio Nombre de la particion
	unit := 1024          //opcional /valor por defecto en KB por eso es 1024
	typee := "P"          //opcional Valores: P, E, L
	fit := "W"            //opcional valores para fit: f, w, b
	var add int           //opcional (para aumentar o reducir el tamaño de una particion)
	var opcion int        // 0 -> crear; 1 -> add; 2 -> delete (por defecto es 0 = CREAR)
	paramC := true        //Para validar que los parametros cumplen con los requisitos
	sizeInit := false     //Sirve para saber si se inicializo size (por si no viniera el parametro por ser opcional) false -> no inicializado
	var sizeValErr string //Para reportar el error si no se pudo convertir a entero el size

	//mismo proceso que el mkdisk para manejar parametros
	for _, parametro := range parametros[1:] {
		tmp := strings.Split(parametro, "=")

		//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
		if len(tmp) != 2 {
			fmt.Println("FDISK Error: Valor desconocido del parametro ", tmp[0])
			paramC = false
			break
		}

		//SIZE
		if strings.ToLower(tmp[0]) == "-size" {
			sizeInit = true
			var err error
			size, err = strconv.Atoi(tmp[1]) //se convierte el valor en un entero
			if err != nil {
				sizeValErr = tmp[1] //guarda para el reporte del error si es necesario validar size
			}

			//PATH
		} else if strings.ToLower(tmp[0]) == "-driveletter" {
			//homonimo al path
			letter = strings.ToUpper(tmp[1]) //Debe estar en mayusculas
			//Se valida si existe el disco ingresado
			carpeta := "./MIA/P1/" //Ruta (carpeta donde se guardara el disco)
			extension := ".dsk"
			path := carpeta + string(letter) + extension
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("FDISK Error: El disco ", letter, " no existe")
				paramC = false
				break // Terminar el bucle porque encontramos un nombre único
			}

			//NAME
		} else if strings.ToLower(tmp[0]) == "-name" {
			// Eliminar comillas
			name = strings.ReplaceAll(tmp[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)

			//UNIT
		} else if strings.ToLower(tmp[0]) == "-unit" {
			//k ya esta predeterminado
			if strings.ToLower(tmp[1]) == "b" {
				//asigno el valor del parametro en su respectiva variable
				unit = 1
			} else if strings.ToLower(tmp[1]) == "m" {
				unit = 1048576 //1024*1024
			} else if strings.ToLower(tmp[1]) != "k" {
				fmt.Println("FDISK Error en -unit. Valores aceptados: b, k, m. ingreso: ", tmp[1])
				paramC = false
				break
			}

			//TYPE
		} else if strings.ToLower(tmp[0]) == "-type" {
			//p esta predeterminado
			if strings.ToLower(tmp[1]) == "e" {
				typee = "E"
			} else if strings.ToLower(tmp[1]) == "l" {
				typee = "L"
			} else if strings.ToLower(tmp[1]) != "p" {
				fmt.Println("FDISK Error en -type. Valores aceptados: e, l, p. ingreso: ", tmp[1])
				paramC = false
				break
			}

			//FIT
		} else if strings.ToLower(tmp[0]) == "-fit" {
			//Si el ajuste es BF (best fit)
			if strings.ToLower(tmp[1]) == "bf" {
				//asigno el valor del parametro en su respectiva variable
				fit = "B"
				//Si el ajuste es WF (worst fit)
			} else if strings.ToLower(tmp[1]) == "ff" {
				//asigno el valor del parametro en su respectiva variable
				fit = "F"
				//Si el ajuste es ff ya esta definido por lo que si es distinto es un error
			} else if strings.ToLower(tmp[1]) != "wf" {
				fmt.Println("FDISK Error en -fit. Valores aceptados: BF, FF o WF. ingreso: ", tmp[1])
				paramC = false
				break
			}

			//DELETE
		} else if strings.ToLower(tmp[0]) == "-delete" {
			if strings.ToLower(tmp[1]) == "full" {
				if opcion == 0 {
					opcion = 2 // 2 es delete
				}
			} else {
				fmt.Println("FDISK Error. Valor de delete desconocido")
				paramC = false
				break
			}

			//ADD
		} else if strings.ToLower(tmp[0]) == "-add" {
			var err error
			add, err = strconv.Atoi(tmp[1]) //se convierte el valor en un entero
			if err != nil {
				fmt.Println("FDISK Error: El valor de \"add\" debe ser un valor numerico. se leyo ", tmp[1])
				paramC = false
				break
			} else {
				if opcion == 0 {
					opcion = 1
				}
			}

			//ERROR EN LOS PARAMETROS LEIDOS
		} else {
			fmt.Println("FDISK Error: Parametro desconocido ", tmp[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	//Si va a crear una particion verificar el size
	if opcion == 0 && paramC {
		if sizeInit { //Si viene el parametro size
			if sizeValErr == "" { //Si es un numero (si es numero la variable sizeValErr sera una cadena vacia)
				if size <= 0 { //se valida que sea mayor a 0 (positivo)
					fmt.Println("FDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", size)
					paramC = false
				}
			} else { //Si sizeValErr es una cadena (por lo que no se pudo dar valor a size)
				fmt.Println("FDISK Error: -size debe ser un valor numerico. se leyo ", sizeValErr)
				paramC = false
			}
		} else { //Si no viene el parametro size
			fmt.Println("FDISK Error: No se encuentra el parametro -size")
			paramC = false
		}
	}

	//si todos los parametros son correctos
	if paramC {
		if letter != "" && name != "" {
			// Abrir y cargar el disco
			filepath := "./MIA/P1/" + letter + ".dsk"
			disco, err := Herramientas.OpenFile(filepath)
			if err != nil {
				fmt.Println("FDisk Error: No se pudo leer el disco")
				return
			}

			//Se crea un mbr para cargar el mbr del disco
			var mbr Structs.MBR
			//Guardo el mbr leido
			if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
				return
			}

			//CREAR (opcion: 0 -> crear; 1 -> add; 2 -> delete)
			if opcion == 0 {
				sizeNewPart := size * unit //Tamaño de la nueva particion (tamaño * unidades)
				isPartExtend := false      //Indica si la particion es extendida
				isName := true             //Valida si el nombre no se repite (true no se repite)
				//guardar := false           //Indica si se debe guardar la particion, es decir, escribir en el disco

				//Si la particion es tipo extendida validar que no exista alguna extendida
				if typee == "E" {
					for i := 0; i < 4; i++ {
						tipo := string(mbr.Partitions[i].Type[:])
						//VER QUE IMPRIME EN TIPO PARA MANEJAR EL ELSE. PARA VER QUE DEVUELVE SI NO EXISTE AUN LA PARTICION
						if tipo != "E" {
							isPartExtend = true
						} else {
							fmt.Println("FDISK Error. Ya existe una particion extendida")
							isPartExtend = false
							break
						}
					}
				}

				//verificar si  el nombre existe en las particiones primarias o extendida
				for i := 0; i < 4; i++ {
					nombre := Structs.GetNamePart(mbr, i)
					if nombre == name {
						fmt.Println("FDISK Error. Ya existe la particion : ", name)
						isName = false
						break
					}
				}

				//INGRESO DE PARTICIONES PRIMARIAS Y/O EXTENDIDA (SIN LOGICAS)
				if (typee == "P" || isPartExtend) && isName { //para que  isPartExtend sea true, typee tendra que ser "E"
					//obtener el tamaño del mbr (el que ocupa fisicamente )
					sizeMBR := int32(binary.Size(mbr))
					fmt.Println("Tamaño fisico de mbr", sizeMBR)
					//Para manejar los demas ajustes hacer un if del fit para llamar a la funcion adecuada
					//F = primer ajuste; B = mejor ajuste; else -> peor ajuste

					//INSERTAR PARTICION
					//if mbr.Partitions[0].Size == 0 {
					mbr.Partitions[0].Size = int32(sizeNewPart)
					mbr.Partitions[0].Start = sizeMBR
					mbr.Partitions[0].Correlative = int32(1)
					copy(mbr.Partitions[0].Name[:], name)
					copy(mbr.Partitions[0].Fit[:], fit)
					copy(mbr.Partitions[0].Status[:], "I")
					copy(mbr.Partitions[0].Type[:], typee)
					//}
					//sobreescribir el mbr
					if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {
						return
					}

					//para verificar que lo guardo
					var TempMBR2 Structs.MBR
					// Read object from bin file
					if err := Herramientas.ReadObject(disco, &TempMBR2, 0); err != nil {
						return
					}
					Structs.PrintMBR(TempMBR2)

					// INGRESO DE PARTICIONES LOGICAS
				} else if typee == "L" && isName {
					fmt.Println("Crear particion logica")
					//validar que el nombre no exista en la logicas si el tipo es "L"
				}
				//a esta altura sigue abierto el archivo
			} else if opcion == 1 {
				//ADD
				//validar que venga unit
				add = add * unit
				if add < 0 {
					fmt.Println("Reducir espacio")
				} else if add > 0 {
					fmt.Println("aumentar espacio")
				} else {
					fmt.Println("FDISK Error. 0 no es un valor valido para aumentar o disminuir particiones")
				}
			} else if opcion == 2 {
				//validar que venga name y driveletter
				fmt.Println("eliminar particion")
			} else {
				//Creo se puede quitar porque nunca va a entrar aqui
				fmt.Println("FDISK Error. Operación desconocida (operaciones aceptadas: crear, modificar o eliminar)")
			}

			// Cierro el disco
			defer disco.Close()
			fmt.Println("======End FDISK======")
		} else {
			fmt.Println("FDISK Error. No se encontro parametro letter y/o name")
		}
	} //Fin if paramC
} //Fin FDisk
