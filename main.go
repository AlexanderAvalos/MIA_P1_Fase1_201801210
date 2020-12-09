package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//estructuras

type Particion struct {
	PartStatus byte
	PartType   byte
	PartFit    [2]byte
	PartStart  int64
	PartSize   int64
	PartName   [16]byte
}

type MBR struct {
	Mbrtamano int64
	Mbrfecha  [20]byte
	Mbrdisk   int64
	Diskfit   [2]byte
	Particion [4]Particion
	Part      [4]bool
	Extend    bool
}

type EBR struct {
	PartStatusE byte
	PartFitE    byte
	PartStartE  int64
	PartSizeE   int64
	PartNextE   int64
	PartNameE   [16]byte
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
		fmt.Println("-------------comando exec---------------")
		comandoexec(linea)
	case "pause":
		fmt.Println("-------------comando pause--------------")
		fmt.Print("Presione enter para continuar")
		lector := bufio.NewReader(os.Stdin)
		entrada, _ := lector.ReadString('\n')
		fmt.Print(entrada)
	case "mkdisk":
		fmt.Println("------------comando mkdisk-------------")
		comando_mksdisk(linea)
	case "rmdisk":
		fmt.Println("------------comando rmdisk--------------")
		comando_rmdisk(linea)
	case "fdisk":
		fmt.Println("------------comando fdisk--------------")
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

	var sizeparticion string = "vacio"
	var rutaparticion string = "vacio"
	var unitparticion string = "vacio"
	var fitparticion string = "vacio"
	var typeparticion string = "vacio"
	var deleteparticion string = "vacio"
	var nameparticion string = "vacio"
	var addparticion string = "vacio"

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
			unitparticion = "kb"
			fmt.Println("unidad KB")
		}
		if banderas[3] == false {
			typeparticion = "p"
			fmt.Println("tipo Primaria")
		}
		if banderas[4] == false {
			fitparticion = "wf"
			fmt.Println("fit peor")
		}
	}

	fdisk(sizeparticion, rutaparticion, unitparticion, typeparticion, fitparticion, deleteparticion, nameparticion, addparticion)
}

func fdisk(size string, ruta string, unit string, tipo string, fit string, eliminar string, nombre string, add string) {

	var sizeparticion int64
	var rutaparticion string = ruta
	var unitparticion string = unit
	var typeparticion string = tipo
	var nameparticion string = nombre
	var fitparticion string = fit
	/*
		var deleteparticion string

		var addparticion int64*/

	sizeparticion, err := strconv.ParseInt(size, 10, 64)
	files, err := os.OpenFile(rutaparticion, os.O_RDWR, 0777)
	defer files.Close()
	if err != nil {
		panic(err)
	}
	if strings.ToLower(unitparticion) == "k" {
		sizeparticion = sizeparticion * 1024
	} else if strings.ToLower(unitparticion) == "m" {
		sizeparticion = sizeparticion * 1024 * 1024
	} else {
		sizeparticion = sizeparticion * 8
	}

	mbrauxiliar := MBR{}
	//var espacio_particion = false
	mbrauxiliar = obtenerMBR(files)

	var contador_particion = 0
	for i := 0; i < 4; i++ {
		if mbrauxiliar.Part[i] == false {
			contador_particion++
		}
	}
	if eliminar == "vacio" || add == "vacio" {
		if mbrauxiliar.Extend == false {
			fmt.Println("Quedan %d Primarias, 1 extendida", (contador_particion - 1))
		} else {
			fmt.Println("Quedan %d Primarias, 0 extendida", contador_particion)
		}
		if contador_particion >= 1 {
			fmt.Println(typeparticion)
			if strings.ToLower(typeparticion) == "p" {
				fmt.Println(sizeparticion)
				ParticionPrimaria(rutaparticion, nameparticion, typeparticion, fitparticion, unitparticion, sizeparticion)
			} else if strings.ToLower(typeparticion) == "e" && mbrauxiliar.Extend == false {

			} else if strings.ToLower(typeparticion) == "l" {

			}

		}
	}

	fmt.Println("size mbr", mbrauxiliar.Mbrtamano)
}

func ParticionExtendida(ruta string, name string, tipo string, fit string, unit string, size int64) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	check(err)
	defer file.Close()
	var auxName [16]byte
	for i, j := range []byte(name) {
		auxName[i] = byte(j)
	}
	mbraux := obtenerMBR(file)
	for i := 0; i < 4; i++ {
		if mbraux.Particion[i].PartType == 'e' || mbraux.Particion[i].PartType == 'E' {
			fmt.Println("solo puede crear una particion extendida")
			Menu()
		}
	}
	var contador int64 = 0

	for i := 0; i < 4; i++ {
		if mbraux.Particion[i].PartStatus != '1' {
			contador += mbraux.Particion[i].PartSize
		}
	}

}

func ParticionPrimaria(ruta string, name string, tipo string, fit string, unit string, size int64) {

	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	check(err)
	defer file.Close()
	var auxName [16]byte
	for i, j := range []byte(name) {
		auxName[i] = byte(j)
	}
	mbraux := obtenerMBR(file)

	var contador int64 = 0
	for i := 0; i < 4; i++ {
		if mbraux.Particion[i].PartStatus != '1' {
			contador += mbraux.Particion[i].PartSize
		}
	}

	if (mbraux.Mbrtamano - contador) >= size {
		var verficar bool = false
		for i := 0; i < 4; i++ {
			if mbraux.Particion[i].PartName == auxName {
				verficar = true
				break
			} else if mbraux.Particion[i].PartType == 'e' || mbraux.Particion[i].PartType == 'E' {
				file.Seek(mbraux.Particion[i].PartStart, 0)
				ebraux := obtenerEBR(file)
				tam := binary.Size(ebraux)
				pos, _ := file.Seek(0, os.SEEK_CUR)
				fmt.Println(pos)
				for tam != 0 && pos < int64(mbraux.Particion[i].PartSize)+int64(mbraux.Particion[i].PartStart) {
					if ebraux.PartNameE == auxName {
						verficar = true
					} else if ebraux.PartNextE == -1 {
						break
					} else {
						file.Seek(ebraux.PartNextE, 0)
						ebraux = obtenerEBR(file)
					}
				}
			}
		}

		if !verficar {
			var indice int = 0
			if strings.ToLower(fit) == "bf" {
				primerAjuste(mbraux, size)
			} else if strings.ToLower(fit) == "ff" {
				primerAjuste(mbraux, size)
			} else if strings.ToLower(fit) == "wf" {
				primerAjuste(mbraux, size)
			}

			if indice != -1 {
				var auxfit [2]byte
				for i, j := range []byte(fit) {
					auxfit[i] = byte(j)
				}
				mbraux.Particion[indice].PartFit = auxfit
				mbraux.Particion[indice].PartType = convertirstring(tipo)
				mbraux.Particion[indice].PartSize = size
				mbraux.Particion[indice].PartStatus = '0'

				if indice == 0 {
					mbraux.Particion[indice].PartStart = int64(binary.Size(mbraux)) + 1
				} else {
					mbraux.Particion[indice].PartStart = mbraux.Particion[indice-1].PartStart + mbraux.Particion[indice-1].PartSize
				}

				file.Seek(0, 0)
				var binario bytes.Buffer
				binary.Write(&binario, binary.BigEndian, &mbraux)
				escribirbinario(file, binario.Bytes())
				fmt.Println("se creo la particion")

			} else {
				fmt.Println("Ya se ah creado el maximo de particiones")
			}
		} else {
			fmt.Println("ya se ah creado un partacion con este nombre")
		}
	} else {
		fmt.Println("ya no tiene espacio en el disco ")
	}
}
func convertirstring(texto string) byte {
	var auxfit byte
	for j := range []byte(texto) {
		auxfit = byte(j)
	}
	return auxfit
}
func primerAjuste(mbraux MBR, size int64) int {
	var verificar bool = false
	var indice int = 0
	for indice < 4 {
		if mbraux.Particion[indice].PartStart == -1 || (mbraux.Particion[indice].PartStart == '1' && mbraux.Particion[indice].PartSize >= size) {
			verificar = true
			break
		}
		indice++
	}
	if verificar {
		return indice
	}
	return -1
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
	check(err)
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
			size = int64(KB)
		} else if strings.ToLower(unit) == "m" {
			for i := 0; i < MB; i++ {
				escribirbinario(file, bin.Bytes())
			}
			size = int64(MB)
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
	temporal.Extend = false
	copy(temporal.Mbrfecha[:], fecha.Format("2006-01-02 15:04:05"))
	copy(temporal.Diskfit[:], ajuste)

	EscribirMBR(file, temporal)

	/*var nuevobuffer bytes.Buffer
	enc := gob.NewEncoder(&nuevobuffer)
	enc.Encode(temporal)
	files, err := os.OpenFile(ruta, os.O_RDWR, 0777)

	files.Seek(0, 0)
	escribirbinario(files, nuevobuffer.Bytes())
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	*/
}

func EscribirMBR(file *os.File, mbraux MBR) {
	file.Seek(0, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &mbraux)
	escribirbinario(file, binario.Bytes())
}

func escribirbinario(file *os.File, binario []byte) {
	_, err := file.Write(binario)
	if err != nil {
		log.Fatal(err)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func obtenerMBR(file *os.File) MBR {
	mbrActual := MBR{}
	var size int = int(unsafe.Sizeof(mbrActual))
	data := leersiguientebyte(file, size)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &mbrActual)
	check(err)
	return mbrActual
}
func obtenerEBR(file *os.File) EBR {
	ebraux := EBR{}
	sizeRead := binary.Size(ebraux)
	data := leersiguientebyte(file, sizeRead)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &ebraux)
	check(err)
	return ebraux
}
func leersiguientebyte(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		fmt.Println(err)
	}
	return bytes
}
