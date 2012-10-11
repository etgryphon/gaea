package main 

import (
    "flag"
    "fmt"
    utl "io/ioutil"
	"os"
//	"strings"
	"os/exec"
	"log"
//	"bytes"\
	"errors"
	"path/filepath"
)

const APP_VERSION = "0.1"

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
  _, err := os.Stat(path)
  if err == nil || os.IsExist(err) { return true, path, nil }
  if os.IsNotExist(err) { return false, path, nil }
  return false, path, err
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
  fmt.Printf("Visited: %s => Base: %s\n", path, base)
  return nil
} 

/*
	@private convert the local package to be used in GAE format
*/
func convertToLocalPackage(root string, name string) (error) {
  err := filepath.Walk(root, convertFileToLocalUse)
  if err != nil { return err }
  return nil;
}

