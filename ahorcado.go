package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	letras := make([]byte, 1024)
	vidas := make([]byte, 1024)
	var lenVidas int
	var letra string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Intruduce tu nombre: ")
	scanner.Scan()
	nombre := scanner.Text()
	fmt.Printf("Introduce la ip del server: ")
	scanner.Scan()
	ip := scanner.Text()
	fmt.Println(ip)
	so, _ := net.Dial("tcp", ip+":4444")

	so.Write([]byte(nombre))

	lenLetras, _ := so.Read(letras)
	letrasConv, _ := strconv.Atoi(string(letras[:lenLetras]))
	ahorcado := strings.Repeat("_", letrasConv)
	var letrasEscritas []string
	for {
		datos := make([]byte, 1024)
		lenDatos, _ := so.Read(datos)

		if string(datos[:lenDatos]) != "endgame" {
			if string(datos[:lenDatos]) == "wingame" {
				clear()
				fmt.Println(ahorcado)
				fmt.Println("Vidas: " + string(vidas[:lenVidas]))
				fmt.Printf("Letras usadas: %v\n", letrasEscritas)
				fmt.Println("Has ganado :)")
				so.Close()
				break
			}
			clear()

			fmt.Println(ahorcado)

			lenVidas, _ := so.Read(vidas)

			fmt.Println("Vidas: " + string(vidas[:lenVidas]))
			fmt.Printf("Letras usadas: %v\n", letrasEscritas)

			fmt.Printf("\n\nIntroduce una letra: ")
			scanner.Scan()
			letra = scanner.Text()
			for is(letra, letrasEscritas) || len(letra) != 1 {
				clear()
				fmt.Println(ahorcado)
				fmt.Println("Vidas: " + string(vidas[:lenVidas]))
				fmt.Printf("Letras usadas: %v\n", letrasEscritas)

				fmt.Printf("\n\nIntroduce una letra: ")
				scanner.Scan()
				letra = scanner.Text()
			}
			so.Write([]byte(letra))
			letrasEscritas = append(letrasEscritas, letra)
			for {
				pos := make([]byte, 1024)
				lenPos, _ := so.Read(pos)
				if string(pos[:lenPos]) == "end" {
					break
				}
				posConv, _ := strconv.Atoi(string(pos[:lenPos]))
				ahorcado2 := []rune(ahorcado)
				ahorcado2[posConv] = rune(letra[0])
				ahorcado = string(ahorcado2)
			}
		} else {
			clear()
			fmt.Println(ahorcado)
			fmt.Println("Vidas: " + string(vidas[:lenVidas]))
			fmt.Printf("Letras usadas: %v\n", letrasEscritas)
			fmt.Println("Has perdido :(")
			so.Close()

			break
		}
	}
}

func clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func is(s string, sa []string) (is bool) {
	for _, s2 := range sa {
		if s2 == s {
			is = true
			return
		}
	}
	is = false
	return
}
