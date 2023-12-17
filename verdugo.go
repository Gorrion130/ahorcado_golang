package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
)

func cliente(so net.Conn, palabra string, vidas int) {
	adivinadas := 0
	nombre := make([]byte, 1024)
	length, _ := so.Read(nombre)

	fmt.Println(string(nombre[:length]) + " se ha conectado.")
	so.Write([]byte(strconv.Itoa(len(palabra))))

	for adivinadas < len(palabra) && vidas > 0 {
		time.Sleep(100)

		letra := make([]byte, 4)

		so.Write([]byte("continue"))

		time.Sleep(100)

		so.Write([]byte(strconv.Itoa(vidas)))

		so.Read(letra)
		mantenerVidas := false

		for pos, letr := range palabra {
			letraDec, _ := utf8.DecodeRune(letra)
			if letr == letraDec {
				so.Write([]byte(strconv.Itoa(pos)))
				adivinadas++
				mantenerVidas = true
				time.Sleep(100)
			}
		}
		if !mantenerVidas {
			vidas--
			time.Sleep(100)
		}
		so.Write([]byte("end"))
		time.Sleep(100)
	}
	if vidas <= 0 {
		so.Write([]byte("endgame"))
		fmt.Println(string(nombre[:length]) + " ha perdido :(")
	} else {
		so.Write([]byte("wingame"))
		fmt.Println(string(nombre[:length]) + " ha ganado :)")
	}
	so.Close()
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	fmt.Printf("Introduce tu palabra: ")

	s.Scan()
	palabra := s.Text()

	fmt.Printf("Introduce las vidas disponibles: ")

	s.Scan()
	vidas, _ := strconv.Atoi(s.Text())

	n, _ := net.Listen("tcp", ":4444")

	for {
		s, _ := n.Accept()
		go cliente(s, palabra, vidas)
	}
}
