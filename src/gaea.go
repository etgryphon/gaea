package main 

import (
    flag "github.com/ogier/pflag"
    "fmt"
    utl "io/ioutil"
	"os"
	"strings"
	"os/exec"
	"log"
//	"bytes"\
	"errors"
	"path/filepath"
)

const APP_VERSION = "0.1"
var fileCount int = 0
var dirCount int = 0
var fileBytesRead int = 0
var fileBytesWritten int = 0

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
    flag.Parse() // Scan the arguments list 

    if *versionFlag {
        fmt.Println("Version:", APP_VERSION)
        return
    }
    cmd := flag.Arg(0)
    name := flag.Arg(1)
    switch cmd {
    	case "help": PrintHelpCommand(nil)
    	case "init": CreateNewProject(name)
    	case "get": GetNewImport(name)
    	default: PrintHelpCommand("Do Not Recognize command: ["+cmd+"]")
    }  
}

func PrintHelpCommand(preamble interface{}) {
  t, _ := utl.ReadFile("templates/help.tmpl")
  if preamble != nil {
  	fmt.Fprintf(os.Stdout, "%s\n%s\n", preamble, t)
  } else {
  	fmt.Fprintf(os.Stdout, "%s\n", t)
  }
  
}

func CreateNewProject(name interface{}){
  if name == nil {
  	name = "blue"
  }
  fmt.Fprintf(os.Stdout, "Creating a New Project...[%s]\n", name)
}

func GetNewImport(name string){
  var err error
  if len(name) < 1 {
    PrintHelpCommand("In order to get a package, you must have a name")
    return
  }
  // check to see if it is present in the GOPATH
  pkgIsThere,fullPath,err := checkIfPackageIsPresent(name)
  if err != nil {
    log.Fatalf("Local Package Verify Error:\n\t%s", err)
  }
  
  if !pkgIsThere {
    err := fetchExternalPackage(name)
    if err != nil {
      log.Fatalf("External Package Fetch Error:\n\t%s", err)
    }
  }
  
  // Convert package to local package
  err = convertToLocalPackage(fullPath, name)
  if err != nil {
    log.Fatalf("Package Conversion Error:\n\t%s", err)
  }
}

/*
	@private function to check for the presence of a local package
*/
func checkIfPackageIsPresent(name string) (bool, string, error) {
  gopath := os.Getenv("GOPATH")
  path := gopath+"/src/"+name+"/"
  formattedPath := filepath.FromSlash(path)
  doesExist, err := checkIfPathExists(formattedPath) 
  return doesExist,formattedPath, err
}
/*
	@private function to check the presents of a directory
*/
func checkIfPathExists(path string)(bool, error){
	_, err := os.Stat(path)
  if err == nil || os.IsExist(err) { return true, nil }
  if os.IsNotExist(err) { return false, nil }
  return false, err
}

/*
	@private function for actually fetching the package
*/
func fetchExternalPackage(name string) (error) {
  log.Println("Fetching External Package: ", name)
  cmd := exec.Command("go", "get", name)
  output,err := cmd.CombinedOutput()
  if err != nil {
    return errors.New(string(output))
  }
  return nil
}

func convertFileToLocalUse(path string, f os.FileInfo, err error) error {
  base := filepath.Base(path)
  if base[0] == '.' { return filepath.SkipDir }
  gopath := os.Getenv("GOPATH")+"/src"
  curDir,_ := os.Getwd()
  newPath := strings.Replace(path, gopath, curDir+"/pkgs", -1)
  
  // Now, check what the item is:
  info, err := os.Stat(path)
  if err != nil { return filepath.SkipDir }
  if info.IsDir() {
    dirCount += 1
  	err = os.MkdirAll(newPath, 0755)
  	if err != nil { return filepath.SkipDir }
  } else {
    fileCount += 1
  	err = translateFile(path, newPath)
  }
  return nil
}

/*
	@private function that will copy and translate files to local use
*/
func translateFile(source string, dest string)(error){
  s, err := os.Open(source) 
  if err != nil {
    return nil
  }
  _, e := os.Stat(dest)
  if e == nil { return nil }
  d, err := os.Create(dest)
  if err != nil {
    return nil
  }
  sInfo,_ := s.Stat()
  buffer := make([]byte, sInfo.Size())
 
  readsize, err := s.Read(buffer)
  fileBytesRead += readsize
  if err != nil {
  	return errors.New(fmt.Sprintf("Could not write read from %s", source))
  }
  
  written, err := d.Write(buffer)
  fileBytesWritten += written
  if err != nil {
    return errors.New(fmt.Sprintf("Could not write to %s", dest))
  }

  // makes sure that # of bytes was written/read correctly.
  if written < readsize {
    return errors.New(fmt.Sprintf("Not enough bytes written to %s", dest))
  }

  if readsize < written {
    return errors.New(fmt.Sprintf("Wrote more bytes than read to %s", dest))
  }
  return nil
}

/*
	@private convert the local package to be used in GAE format
*/
func convertToLocalPackage(root string, name string) (error) {
  fmt.Println("\n\nTransfering GOPATH package ["+name+"] to local use...")
  err := filepath.Walk(root, convertFileToLocalUse)
  if err != nil { return err }
  fmt.Println("\nTotals")
  fmt.Println("\tDirectories: ", dirCount)
  fmt.Println("\tFiles: ", fileCount)
  fmt.Println("\tBytes Read: ", fileBytesRead)
  fmt.Println("\tBytes Written: ", fileBytesWritten)
  fmt.Printf("\nTo Use it in your Google App Engine program:\n\n\timport \"./pkgs/%s\"\n\n", name)
  return nil;
}
