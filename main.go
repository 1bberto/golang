package main

import (
	"fmt"
	"io/ioutil"
	"log"
    "os"        
    "strings"
    "strconv"
    "io"
)
type filesToMove struct {
    basePath string
    folder string
    file os.FileInfo
}
func main(){	
    var filesInFolder []filesToMove

    if(len(os.Args) != 4){
        log.Fatal("INVALID PARAMETERS, \n\t[SOURCE FOLDER] [TARGET FOLDER] [nDAYS TO KEEP]")
    }
    
    sourceFolder := os.Args[1] 
    targetFolder := os.Args[2]
    diasCorte := os.Args[3]
    
    fmt.Printf("SOURCE FOLDER -> %q\nTARGET FOLDER -> %q\nDAYS TO KEEP-> %q\n",sourceFolder, targetFolder,diasCorte)

	getFilesFromFolder(sourceFolder, "" ,&filesInFolder)

    for _, file := range filesInFolder{
        sourcePath, targetPath := createFolder(targetFolder, file)        
        moveFile(sourcePath, targetPath)
    }
}

func moveFile(sourcePath, destPath string) error {
    inputFile, err := os.Open(sourcePath)
    if err != nil {
        return fmt.Errorf("Couldn't open source file: %s", err)
    }
    outputFile, err := os.Create(destPath)
    if err != nil {
        inputFile.Close()
        return fmt.Errorf("Couldn't open dest file: %s", err)
    }
    defer outputFile.Close()
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {
        return fmt.Errorf("Writing to output file failed: %s", err)
    }
    // The copy was successful, so now delete the original file
    // err = os.Remove(sourcePath)
    // if err != nil {
    //     return fmt.Errorf("Failed removing original file: %s", err)
    // }
    return nil
}

func getFilesFromFolder(sourceFolder string, folder string, files *[]filesToMove){
    if(folder != ""){
        sourceFolder = fmt.Sprintf("%s\\%s", sourceFolder, folder)
    }

    filesInFolder, err := ioutil.ReadDir(sourceFolder)

    if err != nil {
        log.Fatal(err)
    }  

    for _, f := range filesInFolder {        
        if(!f.Mode().IsDir()){            
            *files = append(*files, filesToMove{sourceFolder,folder,f})
        }else{
            getFilesFromFolder(sourceFolder, f.Name(), files)
        }
    }
}


func createFolder(target string, fileToMove filesToMove) (sourcePath string, targetPath string) {
    parts := []string{target,strconv.Itoa(fileToMove.file.ModTime().Year()),strconv.Itoa(int(fileToMove.file.ModTime().Month()))}
    source := []string{fileToMove.basePath, fileToMove.file.Name()}
    
    if fileToMove.folder != ""{
        parts = append(parts, fileToMove.folder)
    }

    path := strings.Join(parts,"/")
    //Checks whether the directory exists or not
    if _, err := os.Stat(path); os.IsNotExist(err) {
        os.MkdirAll(path, 1)
    }

    parts = append(parts,fileToMove.file.Name())

    sourcePath = strings.Join(source, "/")
    targetPath = strings.Join(parts,"/")    
    return sourcePath,targetPath
}
