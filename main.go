package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

//variables globales

//estructuras

type Particion struct {
	partstatus byte
	tipo       byte
	partfit    byte
	inicio     int64
	partsize   int64
	partname   [16]byte
}

type MBR struct {
	mbrtamano int64
	mbrfecha  string
	mbrdisk   int64
	diskfit   byte
	particion [4]Particion
	part      [4]bool
	extend    bool
}

type EBR struct {
	partstatus byte
	partfit    byte
	inicio     int64
	partsize   int64
	partname   [16]byte
	partnext   int64
}

//inicio
func main() {
	Menu()
}

func Menu() {
	fmt.Println("ingrese comando a ejecutar")
	var comandoentrante string
	leer := bufio.NewReader(os.Stdin)
	entrada, _ := leer.ReadString('\n')
	comandoentrante = strings.TrimRight(entrada, "\r\n")
	if comandoentrante == "salir" {
		fmt.Println("usted salio exitosamente")
	} else {
		leercomando(comandoentrante)
		Menu()
	}
}

func leercomando(linea string) {
	comando := strings.Split(linea, " ")
	comparador := strings.ToLower(comando[0])
	switch comparador {
	case "exec":
		fmt.Println("comando exec")
		comandoexec(linea)
	case "pause":
		fmt.Println("comando pause")
		fmt.Print("Presione enter para continuar")
		lector := bufio.NewReader(os.Stdin)
		entrada, _ := lector.ReadString('\n')
		fmt.Print(entrada)
	case "mkdisk":
		fmt.Println("comando mkdisk")
	case "rmdisk":
		fmt.Println("comando rmdisk")
	case "fdisk":
		fmt.Println("comando fdisk")
	case "mount":
		fmt.Println("comando mount")
	case "unmount":
		fmt.Println("comando unmount")
	case "rep":
		fmt.Println("comando reportes")
	}
}

//comandos

func comandoexec(linea string) {
	comando := strings.Split(linea, " ")
	ruta := strings.Split(strings.ToLower(comando[1]), "-path->")
	leerarchivo(ruta[1])
}

func leerarchivo(ruta string) {
	file, err := os.Open(ruta)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		comando := strings.Split(scanner.Text(), " ")
		for i := 0; i < len(comando); i++ {
			fmt.Println(comando[i])
			leercomando(comando[i])
		}
	}
}
