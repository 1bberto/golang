package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// FileToMove represents a request to run a command.
type FileToMove struct {
	BasePath string      //TESTE
	Folder   string      //TESTE
	Name     string      //TESTE
	File     os.FileInfo //TESTE
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var FilesInSourceFolder []FileToMove
	var FilesInFolder []FileToMove
	var FilesToDelete []FileToMove
	if len(os.Args) != 4 {
		log.Fatal("INVALID PARAMETERS, \n\t[SOURCE Folder]\n [TARGET Folder]\n [DAYS TO KEEP]")
	}

	sourceFolder := os.Args[1]
	targetFolder := os.Args[2]
	diasCorte := os.Args[3]

	dias, _ := strconv.Atoi(diasCorte)

	fmt.Printf("SOURCE Folder -> %q\nTARGET Folder -> %q\nDAYS TO KEEP-> %q\n", sourceFolder, targetFolder, diasCorte)

	getFilesFromFolder(sourceFolder, "", &FilesInSourceFolder, true)
	//Clean the folders in source folder
	for _, file := range FilesInSourceFolder {
		file.CheckFolder(dias)
	}

	getFilesFromFolder(sourceFolder, "", &FilesInFolder, false)
	//move files from source to target
	for _, file := range FilesInFolder {
		sourcePath, targetPath := createFolder(targetFolder, file)
		fmt.Printf("SOURCE Folder -> %s", sourcePath)
		fmt.Printf("TARGET Folder -> %s", targetPath)
		moveFile(sourcePath, targetPath)
		file.CheckFolder(dias)
	}

	getFilesFromFolder(targetFolder, "", &FilesToDelete, true)
	//Clean the folders in target folder
	for _, file := range FilesToDelete {
		file.CheckFolder(dias)
	}
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source File: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest File: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output File failed: %s", err)
	}
	err = os.Remove(sourcePath)

	if err != nil {
		return fmt.Errorf("Failed removing original File: %s", err)
	}
	return nil
}

func getFilesFromFolder(sourceFolder string, Folder string, Files *[]FileToMove, addFolder bool) {
	if Folder != "" {
		sourceFolder = fmt.Sprintf("%s\\%s", sourceFolder, Folder)
	}

	FilesInFolder, err := ioutil.ReadDir(sourceFolder)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range FilesInFolder {
		if !file.Mode().IsDir() {
			*Files = append(*Files, FileToMove{sourceFolder, Folder, sourceFolder + "\\" + file.Name(), file})
		} else {
			if addFolder {
				*Files = append(*Files, FileToMove{sourceFolder, Folder, sourceFolder + "\\" + file.Name(), file})
			}
			getFilesFromFolder(sourceFolder, file.Name(), Files, addFolder)
		}
	}
}

func createFolder(target string, FileToMove FileToMove) (sourcePath string, targetPath string) {
	parts := []string{target, strconv.Itoa(FileToMove.File.ModTime().Year()), strconv.Itoa(int(FileToMove.File.ModTime().Month()))}
	source := []string{FileToMove.BasePath, FileToMove.File.Name()}

	if FileToMove.Folder != "" {
		parts = append(parts, FileToMove.Folder)
	}

	path := strings.Join(parts, "")
	//Checks whether the directory exists or not
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 1)
	}

	parts = append(parts, FileToMove.File.Name())

	sourcePath = strings.Join(source, "/")
	targetPath = strings.Join(parts, "/")
	return
}

// IsEmpty Check if Folder is empty
func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// CheckFolder check the folder
func (file FileToMove) CheckFolder(diasCorte int) {
	fileToTest, err := os.OpenFile(file.Name, 0, 0x0)

	if err != nil {
		return
	}

	fileToTest.Close()

	if file.File.Mode().IsDir() {
		isEmpty, err := IsEmpty(file.Name)
		if err != nil {
			log.Printf("%s\n", err)
			log.Fatal(err)
		}

		if isEmpty {
			os.Remove(file.Name)
		}
	} else {
		now := time.Now()
		before := now.AddDate(0, 0, diasCorte*-1)

		if file.File.ModTime().Before(before) {
			os.Remove(file.Name)
		}
	}
}
