package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/scrypt"
	"golang.org/x/term"
)

const (
	keyValSep    = "|"
	dataFilePerm = 0644
	dataFilePath = "/home/nurtai/personal/journal/data"
)

var (
	notes    map[string]string
	password []byte
)

func main() {
	if len(os.Args) == 1 {
		printUsage()
		return
	}

	fmt.Println("Reading notes...")
	readNotes()
	fmt.Println()
	defer writeNotes()

	switch os.Args[1] {
	case "list":
		listNotes()
	case "add":
		addNote()
	case "delete":
		deleteNote()
	default:
		printUsage()
	}

	fmt.Println()
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("    list    list all notes")
	fmt.Println("    add     add a new note")
}

func listNotes() {
	fmt.Println("All notes: ")
	i := 0
	for k, v := range notes {
		i++
		fmt.Printf("%d: %s\n", i, k)
		fmt.Println(v)
		fmt.Println()
	}
}

func addNote() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println()

	fmt.Print("Content: ")
	content, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	notes[removeNewLine(title)] = removeNewLine(content)
}

func deleteNote() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println()

	delete(notes, removeNewLine(title))
}

func removeNewLine(input string) string {
	return input[:len(input)-1]
}

func readNotes() {
	notes = make(map[string]string)

	f, err := os.OpenFile(dataFilePath, os.O_RDONLY|os.O_CREATE, dataFilePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if len(data) == 0 {
		return
	}

	plainData := decrypt(data)
	scanner := bufio.NewScanner(strings.NewReader(plainData))
	for scanner.Scan() {
		line := scanner.Text()
		keyVal := strings.SplitN(line, keyValSep, 2)
		notes[keyVal[0]] = keyVal[1]
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func writeNotes() {
	fmt.Println("Writing notes...")
	data := ""
	for k, v := range notes {
		data += k + keyValSep + v + "\n"
	}
	cipherText := encrypt(data)
	err := os.WriteFile(dataFilePath, cipherText, dataFilePerm)
	if err != nil {
		panic(err)
	}
}

func encrypt(plainText string) []byte {
	aesgcm := getAesGcm()
	return aesgcm.Seal(nil, nil, []byte(plainText), nil)
}

func decrypt(cipherText []byte) string {
	aesgcm := getAesGcm()
	plainText, err := aesgcm.Open(nil, nil, cipherText, nil)
	if err != nil {
		panic(err)
	}
	return string(plainText)
}

func getAesGcm() cipher.AEAD {
	if password == nil {
		fmt.Print("Password: ")
		passKey, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		fmt.Print("\n")
		password = passKey
	}
	key, err := scrypt.Key(password, nil, 32768, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		panic(err.Error())
	}
	return aesgcm
}
