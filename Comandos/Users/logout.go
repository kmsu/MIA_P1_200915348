package Comandos

import (
	"MIA_P1_200915348/Structs"
	"fmt"
)

func Logout() {
	fmt.Println("Logout")
	//Validar que haya usuario logeado
	if !Structs.UsuarioActual.Status {
		fmt.Println("LOGOUT ERROR: No existe una sesion iniciada")
		return
	}
	//Cierra sesion
	Structs.UsuarioActual.Status = false
	fmt.Println("Sesion cerrada correctamente. \nHasta luego ", Structs.UsuarioActual.Nombre)
	//Con el estado = false es suficiente pero limpio para evitar posibles conflictos futuros
	Structs.UsuarioActual.Id = ""
	Structs.UsuarioActual.IdGrp = 0
	Structs.UsuarioActual.IdUsr = 0
	Structs.UsuarioActual.Nombre = ""
}
