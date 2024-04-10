package main

import (
	Comandos "MIA_P1_200915348/Comandos"
	DFPM "MIA_P1_200915348/Comandos/AdminPermisosPaths" //DFPM -> Directory, File, Permision Management (Administrador de carpetas, archivos y permisos)
	DM "MIA_P1_200915348/Comandos/AdministradorDiscos"  //DM -> DiskManagement (Administrador de discos)
	FS "MIA_P1_200915348/Comandos/SistemaDeArchivos"    //FS -> FileSystem (sistema de archivos)
	US "MIA_P1_200915348/Comandos/Users"                //US -> UserS
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// MENSAJES DE INICIO
	Ms_inicio := "Bienvenido escriba un comando..."
	Ms_info := "(si desea salir escriba el comando: exit)"
	fmt.Println(Ms_inicio)
	fmt.Println(Ms_info)
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n$: ")
		reader.Scan()

		//entrada := strings.TrimRightFunc(reader.Text(), func(r rune) bool { return r == ' ' })
		entrada := strings.TrimRight(reader.Text(), " ")
		linea := strings.Split(entrada, "#") //para ignorar comentarios desde la consola manual
		//entrada := execute -path=script.txt
		if strings.ToLower(linea[0]) != "exit" {
			//fmt.Println("execute -path=script.txt")
			//entrada = "execute -path=script.txt"
			//analizar(strings.ToLower(entrada)) //--------------------OJO REVISAR SI EL TOLOWER NO AFECTA SI SE HACE DESDE AQUI
			analizar(linea[0]) // Usuario para la parte de comandos de usuario (deben mantenerse durante toda la ejecucion)
		} else {
			fmt.Println("Salida exitosa")
			break
		}
	}
}

func analizar(entrada string) {
	parametros := strings.Split(entrada, " -")
	//NOTA: podría intentar ignorar los espacios tmp2 := strings.TrimRight(parametros[0], " ") y pasar tmp2 en lugar de parametros[0]
	if strings.ToLower(parametros[0]) == "execute" {
		if len(parametros) == 2 {
			tmpParametro := strings.Split(parametros[1], "=")
			if strings.ToLower(tmpParametro[0]) == "path" && len(tmpParametro) == 2 {
				//println("ruta ", tmpParametro[1])
				//abrir el archivo
				archivo, err := os.Open(tmpParametro[1])
				if err != nil {
					fmt.Println("Error al leer el script: ", err)
					return
				}
				defer archivo.Close()
				//creo un lector de bufer para el archivo
				lector := bufio.NewScanner(archivo)
				//leer el archivo linea por linea
				for lector.Scan() {
					//Divido por # para ignorar todo lo que este a la derecha del mismo
					linea := strings.Split(lector.Text(), "#") //lector.Text() retorna la linea leida
					if len(linea[0]) != 0 {
						fmt.Println("\n*********************************************************************************************")
						fmt.Println("Linea en ejecucion: ", linea[0])
						analizar(linea[0])
					}
				}
			} else {
				fmt.Println("EXECUTE ERROR: parametro path no encontrado")
			}
		}

		//--------------------------------- ADMINISTRADOR DE DISCOS ------------------------------------------------
	} else if strings.ToLower(parametros[0]) == "mkdisk" {
		//MKDISK
		//crea un archivo binario que simula un disco con su respectivo MBR
		if len(parametros) > 1 {
			DM.Mkdisk(parametros)
		} else {
			fmt.Println("MKDISK ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "rmdisk" {
		//RMDISK
		if len(parametros) > 1 {
			DM.Rmdisk(parametros)
		} else {
			fmt.Println("RMDISK ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "fdisk" {
		//FDISK
		if len(parametros) > 1 {
			DM.Fdisk(parametros)
		} else {
			fmt.Println("FDISK ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "mount" {
		//MOUNT
		if len(parametros) > 1 {
			DM.Mount(parametros)
		} else {
			fmt.Println("MOUNT ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "unmount" {
		//UNMOUNT
		if len(parametros) > 1 {
			DM.Unmount(parametros)
		} else {
			fmt.Println("UNMOUNT ERROR: parametros no encontrados")
		}

		//--------------------------------- SISTEMA DE ARCHIVOS ----------------------------------------------------
	} else if strings.ToLower(parametros[0]) == "mkfs" {
		//MKFS
		if len(parametros) > 1 {
			FS.Mkfs(parametros)
		} else {
			fmt.Println("MKFS ERROR: parametros no encontrados")
		}

		//--------------------------------------- USERS ------------------------------------------------------------
	} else if strings.ToLower(parametros[0]) == "login" {
		//LOGIN
		if len(parametros) > 1 {
			US.Login(parametros)
		} else {
			fmt.Println("LOGIN ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "logout" {
		//LOGOUT
		if len(parametros) == 1 {
			US.Logout()
		} else {
			fmt.Println("LOGOUT ERROR: Este comando no requiere parametros")
		}

	} else if strings.ToLower(parametros[0]) == "mkgrp" {
		//MKGRP
		if len(parametros) > 1 {
			US.Mkgrp(parametros)
		} else {
			fmt.Println("MKGRP ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "rmgrp" {
		//RMGRP
		if len(parametros) > 1 {
			US.Rmgrp(parametros)
		} else {
			fmt.Println("RMGRP ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "mkusr" {
		//MKUSR
		if len(parametros) > 1 {
			US.Mkusr(parametros)
		} else {
			fmt.Println("MKUSR ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "rmusr" {
		//RMUSR
		if len(parametros) > 1 {
			US.Rmusr(parametros)
		} else {
			fmt.Println("RMUSR ERROR: parametros no encontrados")
		}

		// ------------------ ADMINISTRACION DE CARPETAS, ARCHIVOS Y PERMISOS --------------------------------------
	} else if strings.ToLower(parametros[0]) == "cat" {
		//CAT
		if len(parametros) > 1 {
			DFPM.Cat(parametros)
		} else {
			fmt.Println("CAT ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "mkdir" {
		//MKDIR
		if len(parametros) > 1 {
			DFPM.Mkdir(parametros)
		} else {
			fmt.Println("MKDIR ERROR: parametros no encontrados")
		}

		//--------------------------------------- OTROS ------------------------------------------------------------
	} else if strings.ToLower(parametros[0]) == "rep" {
		//REP
		if len(parametros) > 1 {
			Comandos.Rep(parametros)
		} else {
			fmt.Println("REP ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "pause" {
		fmt.Println("Presione enter para continuar...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

	} else if strings.ToLower(parametros[0]) == "exit" {
		fmt.Println("Salida exitosa")
		os.Exit(0)

	} else if strings.ToLower(parametros[0]) == "" {
		//para agregar lineas con cada enter sin tomarlo como error
	} else {
		fmt.Println("Comando no reconocible")
	}

}
