package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Mkdisk(parametros []string) {
	fmt.Println("MKDISK")
	//valida entrada de parametros del comando leido
	//PARAMETROS: -size -unit
	var size int      //obligatorio
	fit := "F"        //por defecto es ff por eso se inicializa con ese valor (valores para fit: f, w, b pero de entrada se recibe FF, WF o BF)
	unit := 1048576   //1024*1024 Por defecto es M por eso se inicializa con este valor en bytes
	paramC := true    //valida que todos los parametros sean correctos
	sizeInit := false //Para saber si entro el parametro size (obligatorio) false -> no inicializado

	//_ sería el indice pero se omite y con [1:] indicamos que inicie el indice 1 en lugar del 0
	//esto porque en [0] esta el comando mkdisk que estamos ejecutando
	//recorro parametros del mkdisk asignando sus valores segun sea el caso
	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		//divido cada parametro entre nombre del parametro y su valor # -size=25 -> -size, 25
		tmp := strings.Split(tmp2, "=")

		//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
		if len(tmp) != 2 {
			fmt.Println("MKDISK Error: Valor desconocido del parametro ", tmp[0])
			paramC = false
			break //para finalizar el ciclo for con el error y no ejecutar lo que haga falta
		}

		//Debería pasar todos strings.ToLower(tmp[0]) ---------------------------------------------------------------------
		//en tmp valido que parametro viene en su primera posicion y que tenga un valor
		//SIZE
		if strings.ToLower(tmp[0]) == "size" {
			sizeInit = true
			var err error
			size, err = strconv.Atoi(tmp[1]) //se convierte el valor en un entero
			//if err != nil || size <= 0 { //Se manejaria como un solo error
			if err != nil {
				fmt.Println("MKDISK Error: -size debe ser un valor numerico. se leyo ", tmp[1])
				paramC = false
				break
			} else if size <= 0 { //se valida que sea mayor a 0 (positivo)
				fmt.Println("MKDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", tmp[1])
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
			} else if strings.ToLower(tmp[1]) == "wf" {
				//asigno el valor del parametro en su respectiva variable
				fit = "W"
				//Si el ajuste es ff ya esta definido por lo que si es distinto es un error
			} else if strings.ToLower(tmp[1]) != "ff" {
				fmt.Println("MKDISK Error en -fit. Valores aceptados: BF, FF o WF. ingreso: ", tmp[1])
				paramC = false
				break
			}
			//UNIT
		} else if strings.ToLower(tmp[0]) == "unit" {
			//si la unidad es k
			if strings.ToLower(tmp[1]) == "k" {
				//asigno el valor del parametro en su respectiva variable
				unit = 1024
				//si la unidad no es k ni m es error (si fuera m toma el valor con el que se inicializo unit al inicio del metodo)
			} else if strings.ToLower(tmp[1]) != "m" {
				fmt.Println("MKDISK Error en -unit. Valores aceptados: k, m. ingreso: ", tmp[1])
				paramC = false
				break
			}
			//ERROR EN LOS PARAMETROS LEIDOS
		} else {
			fmt.Println("MKDISK Error: Parametro desconocido: ", tmp[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if paramC {
		//Verificar que si se haya inicializado el parametro size
		if sizeInit {
			tam := size * unit
			var path string
			var disco string
			carpeta := "./MIA/P1/" //Ruta (carpeta donde se guardara el disco)
			extension := ".dsk"
			//_, err := os.Stat(filepath.Join(carpeta, nombreNuevo)) usando un join si generara una lista de los a archivos de la carpeta
			//Recorremos las letras aceptadas
			for letra := 'A'; letra <= 'Z'; letra++ {
				path = carpeta + string(letra) + extension
				// Verificar si el archivo con el nombre nuevo ya existe
				_, err := os.Stat(path)
				if os.IsNotExist(err) {
					// Terminar el bucle porque encontramos un nombre único
					disco = string(letra) + extension
					break
				}
			}

			//Debería haber un if por si se acaban las letras pero no sera necesario en este proyecto
			// Create file
			//fmt.Println("Archivo ", path)
			err := Herramientas.CrearDisco(path)
			if err != nil {
				fmt.Println("MKDISK Error:: ", err)
			}
			// Open bin file
			file, err := Herramientas.OpenFile(path)
			if err != nil {
				return
			}

			// create array of byte(0)
			/*
				for i := 0; i < tam; i++ {
					err := Herramientas.WriteObject(file, byte(0), int64(i))
					if err != nil {
						fmt.Println("Error: ", err)
					}
				}*/

			datos := make([]byte, tam)
			newErr := Herramientas.WriteObject(file, datos, 0)
			if newErr != nil {
				fmt.Println("MKDISK Error: ", newErr)
				return
			}

			//obtener hora para el id
			ahora := time.Now()
			//obtener los segundos y minutos
			segundos := ahora.Second()
			minutos := ahora.Minute()
			//concateno los segundos y minutos como una cadena (de 4 digitos)
			cad := fmt.Sprintf("%02d%02d", segundos, minutos)
			//convierto la cadena a numero en un id temporal
			idTmp, err := strconv.Atoi(cad)
			if err != nil {
				fmt.Println("MKDISK Error: no converti fecha en entero para id")
			}
			//fmt.Println("id guardado actual ", idTmp)
			// Create a new instance of MBR
			var newMBR Structs.MBR
			newMBR.MbrSize = int32(tam)
			newMBR.Id = int32(idTmp)
			copy(newMBR.Fit[:], fit)
			copy(newMBR.FechaC[:], ahora.Format("02/01/2006 15:04"))
			// Write object in bin file
			if err := Herramientas.WriteObject(file, newMBR, 0); err != nil {
				return
			}

			// Close bin file
			defer file.Close()

			fmt.Println("\n Se creo el disco ", disco, " de forma exitosa")

			//imprimir el disco creado para validar que todo este correcto
			var TempMBR Structs.MBR
			if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
				return
			}
			Structs.PrintMBR(TempMBR)

			fmt.Println("\n======End MKDISK======")
		} else {
			fmt.Println("MKDISK Error: Falta parametro -size")
		}
	}
}
