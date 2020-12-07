package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

//estructuras

type Particion struct {
	Partstatus byte
	Tipo       byte
	Partfit    byte
	Inicio     int64
	Partsize   int64
	Partname   [16]byte
}

type MBR struct {
	Mbrtamano int64
	Mbrfecha  time.Time
	Mbrdisk   int64
	Diskfit   [2]byte
	Part      [4]bool
	Particion [4]Particion
	Extend    bool
}

type EBR struct {
	Partstatus byte
	Partfit    byte
	Inicio     int64
	Partsize   int64
	Partname   [16]byte
	Partnext   int64
}

//inicio	fmt.Println(buffer.Bytes)
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
		comando_mksdisk(linea)
	case "rmdisk":
		fmt.Println("comando rmdisk")
		comando_rmdisk(linea)
	case "fdisk":
		fmt.Println("comando fdisk")
		comando_fsdisk(linea)
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

func comando_fsdisk(linea string) {
	var banderas = [8]bool{false, false, false, false, false, false, false, false}

	var sizeparticion string
	var rutaparticion string
	var unitparticion string
	var fitparticion string
	var typeparticion string
	var deleteparticion string
	var nameparticion string
	var addparticion string

	comando := strings.Split(linea, " ")
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-size->") {
			banderas[0] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-unit->") {
			banderas[1] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-path->") {
			banderas[2] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-type->") {
			banderas[3] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-fit->") {
			banderas[4] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-delete->") {
			banderas[5] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-name->") {
			banderas[6] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-add->") {
			banderas[7] = true
		}
	}
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-size->") {
			aux := strings.Split(comando[i], "->")
			sizeparticion = aux[1]
			fmt.Println("tam", sizeparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-unit->") {
			aux := strings.Split(comando[i], "->")
			unitparticion = aux[1]
			fmt.Println("unit", unitparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-path->") {
			aux := strings.Split(comando[i], "->")
			rutaparticion = aux[1]
			fmt.Println("ruta", rutaparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-type->") {
			aux := strings.Split(comando[i], "->")
			typeparticion = aux[1]
			fmt.Println("tipo", typeparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-fit->") {
			aux := strings.Split(comando[i], "->")
			fitparticion = aux[1]
			fmt.Println("fit", fitparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-delete->") {
			aux := strings.Split(comando[i], "->")
			deleteparticion = aux[1]
			fmt.Println("eliminar", deleteparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-name->") {
			aux := strings.Split(comando[i], "->")
			nameparticion = aux[1]
			fmt.Println("nombre", nameparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-add->") {
			aux := strings.Split(comando[i], "->")
			addparticion = aux[1]
			fmt.Println("agregar", addparticion)
		}
	}
	if banderas[5] == false {
		if banderas[1] == false {
			unitparticion = "KB"
			fmt.Println("unidad KB")
		}
		if banderas[3] == false {
			typeparticion = "P"
			fmt.Println("tipo Primaria")
		}
		if banderas[4] == false {
			typeparticion = "WF"
			fmt.Println("fit peor")
		}
	}
}
func comando_mksdisk(linea string) {
	var banderas = [4]bool{false, false, false, false}
	var sizearchivo string
	var rutaarchivo string
	var unitarchivo string
	var fitarchivo string
	comando := strings.Split(linea, " ")
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-size->") {
			banderas[0] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-path->") {
			banderas[2] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-fit->") {
			banderas[3] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-unit->") {
			banderas[1] = true
		}
	}
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-size->") && banderas[0] == true {
			aux := strings.Split(comando[i], "->")
			sizearchivo = aux[1]
			fmt.Println("tam", sizearchivo)
		} else if strings.Contains(strings.ToLower(comando[i]), "-path->") && banderas[2] == true {
			aux := strings.Split(comando[i], "->")
			rutaarchivo = aux[1]
			fmt.Println("ruta", rutaarchivo)
		} else if strings.Contains(strings.ToLower(comando[i]), "-fit->") && banderas[3] == true {
			aux := strings.Split(comando[i], "->")
			fitarchivo = aux[1]
			fmt.Println("fit", fitarchivo)
		} else if strings.Contains(strings.ToLower(comando[i]), "-unit->") && banderas[1] == true {
			aux := strings.Split(comando[i], "->")
			unitarchivo = aux[1]
			fmt.Println("unit", unitarchivo)
		}
	}
	if banderas[1] == false {
		unitarchivo = "MB"
		fmt.Println("unidad MB")
	}
	if banderas[3] == false {
		fitarchivo = "FF"
		fmt.Println("fit FF")
	}

	mkdisk(sizearchivo, rutaarchivo, unitarchivo, fitarchivo)

}
func comando_rmdisk(linea string) {

	comando := strings.Split(linea, " ")
	var rutaarchivo string = ""

	if strings.Contains(strings.ToLower(comando[1]), "-path->") {
		aux := strings.Split(comando[1], "->")
		rutaarchivo = aux[1]
		fmt.Println("ruta", rutaarchivo)
	}

	err := os.Remove(rutaarchivo)
	if err != nil {
		log.Fatal("Error", err)
	} else {
		fmt.Println("eliminado correctamente")
	}
	//rmdisk -path->/home/alex/disco/Disco1.dsk
}

//Mkdisk -Size->1000 -unit->K -path->/home/alex/disco/Disco2.dsk -fit->BF
func mkdisk(size string, ruta string, unidad string, ajuste string) {
	sizeArchivo, err := strconv.ParseInt(size, 10, 64)
	var rutaArchivo string = ruta
	var unitArchivo string = unidad
	var fitArchivo string = ajuste
	if err != nil {
		log.Fatal(err)
	}

	if sizeArchivo > 0 {
		crearArchivo(sizeArchivo, unitArchivo, rutaArchivo, fitArchivo)
	}

}

func crearArchivo(size int64, unit string, ruta string, ajuste string) {
	var path = ""
	nombre := strings.Split(ruta, "/")
	for i := 0; i < (len(nombre) - 1); i++ {
		path = path + nombre[i] + "/"
	}
	err := os.MkdirAll(path, 0777)
	file, err := os.Create(ruta)

	var BT int = int(size)
	var KB int = BT * 1024
	var MB int = BT * 1024 * 1024

	var binario int8 = 0
	aux := &binario
	var bin bytes.Buffer
	binary.Write(&bin, binary.BigEndian, aux)

	if file != nil {
		if strings.ToLower(unit) == "k" {
			for i := 0; i < KB; i++ {
				escribirbinario(file, bin.Bytes())
			}
		} else if strings.ToLower(unit) == "m" {
			for i := 0; i < MB; i++ {
				escribirbinario(file, bin.Bytes())
			}
		} else {
			fmt.Println("No se pudo crear el disco")
		}
	} else {
		fmt.Println("no se pudo crear el Archivo")
	}

	numero := rand.NewSource(time.Now().UnixNano())
	random := rand.New(numero)
	temporal := MBR{Mbrtamano: size, Mbrdisk: int64(random.Intn(100))}
	fecha := time.Now()
	temporal.Part[0] = false
	temporal.Part[1] = false
	temporal.Part[2] = false
	temporal.Part[3] = false
	temporal.Mbrfecha = fecha
	temporal.Extend = false
	copy(temporal.Diskfit[:], ajuste)

	var nuevobuffer bytes.Buffer
	enc := gob.NewEncoder(&nuevobuffer)
	enc.Encode(temporal)
	files, err := os.OpenFile(ruta, os.O_RDWR, 0777)

	files.Seek(0, 0)
	escribirbinario(files, nuevobuffer.Bytes())
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}

}

func escribirbinario(file *os.File, binario []byte) {
	_, err := file.Write(binario)
	if err != nil {
		log.Fatal(err)
	}
}
