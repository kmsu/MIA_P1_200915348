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
		if strings.ToLower(tmp[0]) == "size" {
			sizeInit = true
			var err error
			size, err = strconv.Atoi(tmp[1]) //se convierte el valor en un entero
			if err != nil {
				sizeValErr = tmp[1] //guarda para el reporte del error si es necesario validar size
			}

			//PATH
		} else if strings.ToLower(tmp[0]) == "driveletter" {
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
		} else if strings.ToLower(tmp[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(tmp[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)

			//UNIT
		} else if strings.ToLower(tmp[0]) == "unit" {
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
		} else if strings.ToLower(tmp[0]) == "type" {
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
		} else if strings.ToLower(tmp[0]) == "fit" {
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
		} else if strings.ToLower(tmp[0]) == "delete" {
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
		} else if strings.ToLower(tmp[0]) == "add" {
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

				//Si la particion es tipo extendida validar que no exista alguna extendida
				isPartExtend := true //Indica si aun esta disponible la particion extendida
				if typee == "E" {
					for i := 0; i < 4; i++ {
						tipo := string(mbr.Partitions[i].Type[:])
						if tipo == "E" {
							isPartExtend = false
							fmt.Println("FDISK Error. Ya existe una particion extendida")
							break
						}
					}
				}

				//verificar si  el nombre existe en las particiones primarias o extendida
				isName := true //Valida si el nombre no se repite (true no se repite)
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name {
						isName = false
						fmt.Println("FDISK Error. Ya existe la particion : ", name)
						break
					}
				}

				if isName {
					//validar que tampoco exista en las logicas
					//fmt.Println("Verificar que no exista una logica con el nombre que se quiere agregar")
					fmt.Println("")
				}

				//INGRESO DE PARTICIONES PRIMARIAS Y/O EXTENDIDA (SIN LOGICAS)
				sizeNewPart := size * unit //Tamaño de la nueva particion (tamaño * unidades)
				guardar := false           //Indica si se debe guardar la particion, es decir, escribir en el disco
				var newPart Structs.Partition
				if (typee == "P" || isPartExtend) && isName { //para que  isPartExtend sea true, typee tendra que ser "E"
					sizeMBR := int32(binary.Size(mbr)) //obtener el tamaño del mbr (el que ocupa fisicamente: 165)
					//Para manejar los demas ajustes hacer un if del fit para llamar a la funcion adecuada
					//F = primer ajuste; B = mejor ajuste; else -> peor ajuste

					//INSERTAR PARTICION (Primer ajuste)
					mbr, newPart = primerAjuste(mbr, typee, sizeMBR, int32(sizeNewPart), name, fit) //int32(sizeNewPart) es para castear el int a int32 que es el tipo que tiene el atributo en el struct Partition
					guardar = newPart.Size != 0

					//escribimos el MBR en el archivo. Lo que no se llegue a escribir en el archivo (aqui) se pierde, es decir, los cambios no se guardan
					if guardar {
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
						fmt.Println("\nParticion con nombre " + name + " creada exitosamente")
						Structs.PrintMBR(TempMBR2)
					} else {
						//Lo podría eliminar pero tendria que modificar en el metodo del ajuste todos los errores para que aparezca el nombre que se intento ingresar como nueva particion
						fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
					}

					// INGRESO DE PARTICIONES LOGICAS
				} else if typee == "L" && isName {
					fmt.Println("Crear particion logica")
					//validar que el nombre no exista en la logicas si el tipo es "L"
				}
				//a esta altura sigue abierto el archivo

				//------------------------------ADD---------------------
			} else if opcion == 1 {
				add = add * unit
				//-------------------------si se quita espacio----------------------------------------------------------------------
				//Particiones extendida o primarias
				if add < 0 {
					fmt.Println("Reducir espacio")
					reducir := true //Si cambia a false es que redujo una de las primarias o la extendida
					for i := 0; i < 4; i++ {
						nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
						if nombre == name {
							reducir = false
							newSize := mbr.Partitions[i].Size + int32(add)
							if newSize > 0 {
								mbr.Partitions[i].Size += int32(add)
								if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
									return
								}
								fmt.Println("Particion con nombre ", name, " se redujo correctamente")
							} else {
								fmt.Println("FDISK Error. El tamaño que intenta eliminar es demasiado grande")
							}
							break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
						}
					}

					//particiones logicas
					if reducir {
						//Reducir logica, si se reduce alguna logica cambiar reducir a false
						fmt.Println("Reducir logicas")
					}

					if reducir {
						fmt.Println("FDISK Error. No existe la particion a reducir")
					}

					//------------------Si se aumenta espacio-----------------------------------------------------------------------
				} else if add > 0 {
					fmt.Println("aumentar espacio")
					//Primarias y/o extendida
					evaluar := 0
					//Si el aumento es en particion 1
					if Structs.GetName(string(mbr.Partitions[0].Name[:])) == name {
						if mbr.Partitions[1].Start == 0 {
							if mbr.Partitions[2].Start == 0 {
								if mbr.Partitions[3].Start == 0 {
									evaluar = int(mbr.MbrSize - mbr.Partitions[0].GetEnd())
								} else {
									evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[0].GetEnd())
								}
							} else {
								evaluar = int(mbr.Partitions[2].Start - mbr.Partitions[0].GetEnd())
							}
						} else {
							evaluar = int(mbr.Partitions[1].Start - mbr.Partitions[0].GetEnd())
						}

						//evaluar > 0 -> si hay espacio para aumentar. add <= evaluar -> si lo que quiero aumentar cabe en el espacio disponible
						if evaluar > 0 && add <= evaluar {
							//aumenta el tamaño de 1
							mbr.Partitions[0].Size += int32(add)
							if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
								return
							}
							fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
						} else {
							fmt.Println("FDISK Error. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
						}
						//Particion 2
					} else if Structs.GetName(string(mbr.Partitions[1].Name[:])) == name {
						if mbr.Partitions[2].Start == 0 {
							if mbr.Partitions[3].Start == 0 {
								evaluar = int(mbr.MbrSize - mbr.Partitions[1].GetEnd())
							} else {
								evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[1].GetEnd())
							}
						} else {
							evaluar = int(mbr.Partitions[2].Start - mbr.Partitions[1].GetEnd())
						}
						//aumenta el tamaño de 2
						if evaluar > 0 && add <= evaluar {
							mbr.Partitions[1].Size += int32(add)
							if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
								return
							}
							fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
						} else {
							fmt.Println("FDISK Error. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
						}
						//Particion 3
					} else if Structs.GetName(string(mbr.Partitions[2].Name[:])) == name {
						if mbr.Partitions[3].Start == 0 {
							evaluar = int(mbr.MbrSize - mbr.Partitions[2].GetEnd())
						} else {
							evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[2].GetEnd())
						}
						//aumenta el tamaño de 3
						if evaluar > 0 && add <= evaluar {
							mbr.Partitions[2].Size += int32(add)
							if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
								return
							}
							fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
						} else {
							fmt.Println("FDISK Error. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
						}
						//Particion 4
					} else if Structs.GetName(string(mbr.Partitions[3].Name[:])) == name {
						evaluar = int(mbr.MbrSize - mbr.Partitions[3].GetEnd())
						//aumenta el tamaño de 4
						if evaluar > 0 && add <= evaluar {
							mbr.Partitions[3].Size += int32(add)
							if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
								return
							}
							fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
						} else {
							fmt.Println("FDISK Error. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
						}
					} else {
						//Verificar si fueran las logicas
						fmt.Println("aumentar logicas")
						//si despues de intentar aumentar las logicas no encontro nada, reportar el error que no existe la particion
					}
				} else {
					fmt.Println("FDISK Error. 0 no es un valor valido para aumentar o disminuir particiones")
				}

				//--------------------- Eliminar particiones -----------------------------------------------------
			} else if opcion == 2 {
				fmt.Println("eliminar particion")
				//-------- primarias o extendida-----------------------------------------------------
				del := true //para saber si se elimino la particion (true es que no se elimino, esto para facilitar el if que valida esta varible)
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name {
						var newPart Structs.Partition
						mbr.Partitions[i] = newPart
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
							return
						}
						del = false
						fmt.Println("particion con nombre ", name, " eliminada")
						break
					}
				}

				//Eliminar logicas
				if del {
					//si elimina una logica cambia a false
					fmt.Println("Eliminar logicas")
				}

				//Valido si no se elimino nada para reportar el error
				if del {
					fmt.Println("FDISK Error. No se encontro la particion a eliminar")
				}

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

func primerAjuste(mbr Structs.MBR, typee string, sizeMBR int32, sizeNewPart int32, name string, fit string) (Structs.MBR, Structs.Partition) {
	var newPart Structs.Partition
	var noPart Structs.Partition //para revertir el set info (simula volverla null)

	//PARTICION 1 (libre) - (size = 0 no se ha creado)
	if mbr.Partitions[0].Size == 0 {
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if mbr.Partitions[1].Size == 0 {
			if mbr.Partitions[2].Size == 0 {
				//caso particion 4 (no existe)
				if mbr.Partitions[3].Size == 0 {
					//859 <= 1024 - 165
					if sizeNewPart <= mbr.MbrSize-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} else {
					//particion 4 existe
					// 600 < 765 - 165 (600 maximo aceptado)
					if sizeNewPart <= mbr.Partitions[3].Start-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						//Si cabe despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Reordeno el correlativo para que coincida con el numero de particion en que se guardo
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					}
				}
				//Fin no existe particion 4
			} else {
				// 3 existe
				//entre mbr y 3 -> 300 <= 465 -165
				if sizeNewPart <= mbr.Partitions[2].Start-sizeMBR {
					mbr.Partitions[0] = newPart
				} else {
					//si no cabe entre el mbr y 3 debe ser despues de 3, es decir, en 4
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} else {
						//4 existe
						//hay espacio entre 3 y 4
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Reordenando los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3 //new part traia 4 y quedo en la tercer particion por eso tambien se modifica aqui
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//Hay espacio despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//reconfiguro los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} //fin si hay espacio entre 3 y 4
				} //fin si no cabe antes de 3
			} //fin 3 existe
		} else {
			//2 existe
			//Si la nueva particion se puede guardar antes de 2
			if sizeNewPart <= mbr.Partitions[1].Start-sizeMBR {
				mbr.Partitions[0] = newPart
			} else {
				//Si no cabe entre mbr y 2
				//Validar si existen 3 y 4
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[2].Size == 0 {
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} else {
						//4 existe (estamos entre 2 y 4)
						//62 < 69-6 (62 maximo aceptado)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							//Si no cabe entre 2 y 4, ver si cabe despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							if sizeNewPart <= mbr.MbrSize-newPart.Start { //1 <= 100-99
								mbr.Partitions[2] = mbr.Partitions[3]
								mbr.Partitions[3] = newPart
								//reordeno correlativos
								mbr.Partitions[2].Correlative = 3
							} else {
								newPart = noPart
								fmt.Println("FDISK Error. Espacio insuficiente")
							}
						} //Fin si cabe antes o despues de 4
					} //fin de 4 existe o no existe
				} else {
					//3 existe
					//entre 2 y 3
					if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
						mbr.Partitions[0] = mbr.Partitions[1]
						mbr.Partitions[1] = newPart
						//Reordeno correlativos
						mbr.Partitions[0].Correlative = 1
						mbr.Partitions[1].Correlative = 2
					} else if mbr.Partitions[3].Size == 0 {
						//entre 3 y el final
						//cambiamos el inicio de la nueva particion porque 3 existe y no cabe antes de 3
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} else {
						//si 4 existe
						//hay espacio entre 3 y 4
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Reordeno correlativos
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//entre 4 y el final
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Reordeno correlativos
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} //Fin si 4 existe o no (3 activa)
				} //Fin 3 existe o no existe
			} //Fin entre 2 y final (antes de 2 o depues de 2)
		} //Fin 2 existe o no existe
		//Fin de 1 no existe

		//PARTICION 2 (no existe)
	} else if mbr.Partitions[1].Size == 0 {
		//Si hay espacio entre el mbr y particion 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start { //particion 1 ya existe (debe existir para entrar a este bloque)
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno correlativo
			mbr.Partitions[1].Correlative = 2
		} else {
			//Si no hay espacio antes de particion 1
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2) //el nuevo inicio es donde termina 1
			if mbr.Partitions[2].Size == 0 {
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[1] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} else {
					//4 existe
					//entre 1 y 4
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[1] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						//despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						//Reordeno correlativo
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} //Fin 4 existe o no existe
			} else {
				//3 Activa
				//entre 1 y 3
				if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
					mbr.Partitions[1] = newPart
				} else {
					//despues de 3
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
							//corrijo el correlativo
							mbr.Partitions[3].Correlative = 4
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} else {
						//4 existe
						//entre 3 y 4
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Corrijo el correlativo
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//Despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Corrijo los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} //fin 4 existe o no existe
				} //Fin para entre 1 y 3, y despues de 3
			} //Fin 3 existe o no existe
		} //Fin antes o despues de particion 1
		//Fin particion 2 no existe

		//PARTICION 3
	} else if mbr.Partitions[2].Size == 0 {
		//antes de 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno los correlativos
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			//entre 1 y 2
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				//Reordeno correlativo
				mbr.Partitions[2].Correlative = 3
			} else {
				//despues de 2
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[2] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} else {
					//4 existe
					//entre 2 y 4
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[2] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						//despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						//Reordeno correlativo
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} //Fin de 4 existe o no existe
			} //Fin espacio entre 1 y 2 o despues de 2
		} //Fin espacio antes de 1
		//Fin particion 3

		//PARTICION 4
	} else if mbr.Partitions[3].Size == 0 {
		//antes de 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[3] = mbr.Partitions[2]
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno los correlativos
			mbr.Partitions[3].Correlative = 4
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			//si no cabe antes de 1
			//entre 1 y 2
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				//Reordeno correlativos
				mbr.Partitions[3].Correlative = 4
				mbr.Partitions[2].Correlative = 3
			} else if sizeNewPart <= mbr.Partitions[2].Start-mbr.Partitions[1].GetEnd() {
				//entre 2 y 3
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = newPart
				//Reordeno correlativo
				mbr.Partitions[3].Correlative = 4
			} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[2].GetEnd() {
				//despues de 3
				newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
				mbr.Partitions[3] = newPart
			} else {
				newPart = noPart
				fmt.Println("FDISK Error. Espacio insuficiente")
			}
		} //Fin antes y despues de 1
		//Fin particion 4
	} else {
		newPart = noPart
		fmt.Println("FDISK Error. Particiones primarias y/o extendidas ya no disponibles")
	}

	return mbr, newPart
}
