package Comandos

import (
	"MIA_P1_200915348/Herramientas"
	"MIA_P1_200915348/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Login(parametros []string) {
	fmt.Println("Login")
	var user string //obligatorio
	var pass string //obligatorio
	var id string   //obligatorio. Id de la particion en la que quiero iniciar sesion
	paramC := true

	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		tmp := strings.Split(tmp2, "=") //separo para obtener su valor parametro=valor

		//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
		if len(tmp) != 2 {
			fmt.Println("LOGIN Error: Valor desconocido del parametro ", tmp[0])
			paramC = false
			break
		}

		//Capturar valores de los parametros
		//ID
		if strings.ToLower(tmp[0]) == "id" {
			id = strings.ToUpper(tmp[1])

			//USER
		} else if strings.ToLower(tmp[0]) == "user" {
			user = tmp[1]

			//PASS
		} else if strings.ToLower(tmp[0]) == "pass" {
			pass = tmp[1]

			//ERROR EN LOS PARAMETROS LEIDO
		} else {
			fmt.Println("LOGIN ERROR: Parametro desconocido: ", tmp[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if paramC {
		if id != "" && user != "" && pass != "" {
			//Validar que no haya usuario logeado
			if Structs.UsuarioActual.Status {
				fmt.Println("LOGIN ERROR: Ya existe una sesion iniciada, cierre sesion para iniciar otra")
				return
			}

			//CARGAR EL DISCO DONDE PODRÍA ESTAR LA PARTICION
			disk := id[0:1] // cargar el disco
			//generar ruta del disco que puede contener el id
			carpeta := "./MIA/P1/" //Ruta (carpeta donde se leera el disco)
			extension := ".dsk"
			rutaDisco := carpeta + disk + extension

			//abrir el disco que podría contener el id
			disco, err := Herramientas.OpenFile(rutaDisco)
			if err != nil {
				return
			}

			//cargar el mbr
			var mbr Structs.MBR
			if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
				return
			}

			//cerrar el archivo del disco
			defer disco.Close()

			//Asegurar que el id exista
			logear := false
			index := -1
			for i := 0; i < 4; i++ {
				identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
				if identificador == id {
					logear = true
					index = i
					break //para que ya no siga recorriendo si ya encontro la particion
				}
			}

			if logear {
				var superBloque Structs.Superblock
				err := Herramientas.ReadObject(disco, &superBloque, int64(mbr.Partitions[index].Start))
				if err != nil {
					fmt.Println("LOGIN Error. Particion sin formato")
					return
				}

				//Se que el users.txt esta en el inodo 1
				var inodo Structs.Inode
				//le agrego una estructura inodo porque busco el inodo 1 (sabemos que aqui esta users.txt)
				Herramientas.ReadObject(disco, &inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

				//leer datos del users.txt (todos los fileblocks que esten en este inodo (archivo))
				var contenido string
				var fileBlock Structs.Fileblock
				for _, item := range inodo.I_block {
					if item != -1 {
						Herramientas.ReadObject(disco, &fileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
						contenido += string(fileBlock.B_content[:])
					}
				}

				linea := strings.Split(contenido, "\n")
				//UID, Tipo, Grupo, Usuario, contraseña

				loginFail := true //para saber si encontro el usuaio
				for _, reg := range linea {
					usuario := strings.Split(reg, ",")

					if len(usuario) == 5 {
						//que no este borrado logicamente
						if usuario[0] != "0" {
							if usuario[3] == user {
								if usuario[4] == pass {
									loginFail = false
									Structs.UsuarioActual.Id = id  //id de la particion
									buscarIdGrp(linea, usuario[2]) //id del grupo al que pertenece el usuario
									idUsr(usuario[0])              //id del usuario
									Structs.UsuarioActual.Nombre = user
									Structs.UsuarioActual.Status = true
									fmt.Println("Inicio de sesion exitoso. \nBienvenido ", user)
								} else {
									loginFail = false
									fmt.Println("LOGIN ERROR: Contraseña incorrecta")
								}
								break
							}
						}
					}
				}

				if loginFail {
					fmt.Println("LOGIN ERROR: No se encontro el usuario")
				}
			} else {
				fmt.Println("LOGIN ERROR: Id no existe")
			}
		} else {
			fmt.Println("LOGIN ERROR: Falta alguno de los siguientes parametros -> id, user o pass")
		}
	}
}

func buscarIdGrp(lineaID []string, grupo string) {
	for _, registro := range lineaID[:len(lineaID)-1] {
		datos := strings.Split(registro, ",")
		if len(datos) == 3 {
			if datos[2] == grupo {
				//convertir a numero
				id, errId := strconv.Atoi(datos[0])
				if errId != nil {
					fmt.Println("LOGIN ERROR: Error desconcocido con el idGrp")
					return
				}
				Structs.UsuarioActual.IdGrp = int32(id)
				return
			}
		}
	}
}

func idUsr(id string) {
	idU, errId := strconv.Atoi(id)
	if errId != nil {
		fmt.Println("LOGIN ERROR: Error desconcocido con el idUsr")
		return
	}
	Structs.UsuarioActual.IdUsr = int32(idU)
}
