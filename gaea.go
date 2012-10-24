package main 

import (
    flag "github.com/ogier/pflag"
    su "github.com/etgryphon/stringUp"
    "fmt"
    "io"
	"bufio"
	"os"
	"strings"
	"os/exec"
	"log"
//	"bytes"
	"errors"
	"path/filepath"
	"html/template"
)

const APP_VERSION = "0.1"
const inputDelim = '\n'
var fileCount int = 0
var dirCount int = 0
var fileBytesRead int = 0
var fileBytesWritten int = 0

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
    flag.Parse() // Scan the arguments list 
	fmt.Fprintln(os.Stdout, LOGO)
    if *versionFlag {
        fmt.Println("Version:", APP_VERSION)
        return
    }
    cmd := flag.Arg(0)
    name := flag.Arg(1)
    switch cmd {
    	case "help": printHelpCommand(nil)
    	case "init": createNewProject(name)
    	case "get": getNewImport(name)
    	case "run": runDevelopmentServer(name)
    	default: printHelpCommand("Do Not Recognize command: ["+cmd+"]")
    }  
}

func printHelpCommand(preamble interface{}) {
  if preamble != nil {
  	fmt.Fprintf(os.Stdout, "%s\n%s\n", preamble, HELP_MESSAGE)
  } else {
  	fmt.Fprintf(os.Stdout, "%s\n", HELP_MESSAGE)
  } 
}

func createNewProject(name string){
  if len(name) < 1 {
	name = readProjectName()
  }
  camelName := su.CamelCase(name)
  projPath := camelName+"Project"
  fmt.Fprintf(os.Stdout, "Creating a New Project...[%s] in directory: %s\n", camelName, projPath)
  createGAEDirectoryStructure(projPath, camelName)
}

func readProjectName()(name string){
  r := bufio.NewReader(os.Stdin)

  fmt.Fprintf(os.Stdout, "Please Enter a Name for your GAE project: ")
  name, err := r.ReadString(inputDelim)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  name = strings.Trim(name, " \n")
  return name
}

func createGAEDirectoryStructure(path string, name string){
  var err error
  
  // Create Base Project Directory
  err = createProjectDirectory("./", path, 1)
  if err != nil { panic(err) }
  // Create app.yml
  err = createProjectFile( name, "./"+path+"/", "app.yml", YML_TEMPLATE, 2)
  if err != nil { panic(err) }
  
  // Create the main folder
  err = createProjectDirectory("./"+path+"/", name, 2)
  if err != nil { panic(err) }
  err = createProjectFile( name, "./"+path+"/"+name+"/", name+".go", BASE_GAE_APP_TEMPLATE, 3)
  if err != nil { panic(err) }
  
  // Create Public Folder
  err = createProjectDirectory("./"+path+"/", "public", 2)
  if err != nil { panic(err) }
  // Create Public stylesheets folder
  err = createProjectDirectory("./"+path+"/public/", "stylesheets", 3)
  if err != nil { panic(err) }
  // Create Public Images Folder
  err = createProjectDirectory("./"+path+"/public/", "images", 3)
  if err != nil { panic(err) }
  // Create Public Javascript Folder
  err = createProjectDirectory("./"+path+"/public/", "js", 3)
  if err != nil { panic(err) }
  
}

func createProjectDirectory(path, name string, level int)(error){
  fullPath := path+name
  fmtStr := fmt.Sprintf("%%%ds\n", 1+len(name)+(level*2) )
  exists, err := checkIfPathExists(fullPath)
  if !exists {
  	fmt.Fprintf(os.Stdout, fmtStr, "+"+name)
  	err = os.Mkdir(fullPath, 0755)
  } else {
  	fmt.Fprintf(os.Stdout, fmtStr, "-"+name) 	
  }
  return err
}

func createProjectFile(pkgName,path,name,tmplInstance string, level int)(error){
  fmtStr := fmt.Sprintf("%%%ds\n", 1+len(name)+(level*2) )
  fileIO, err := os.Create(path+name)
  tmpl, err := template.New(name).Parse(tmplInstance)
  if err != nil { panic(err) }
  err = tmpl.Execute(fileIO, pkgName)
  if err != nil { panic(err) }
  fileIO.Close()
  fmt.Fprintf(os.Stdout, fmtStr, "+"+name)
  return err
}

func getNewImport(name string){
  var err error
  if len(name) < 1 {
    printHelpCommand("In order to get a package, you must have a name")
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
  printOutTransferInformation(name)
  return nil;
}

func printOutTransferInformation(name string){
  fmt.Fprintln(os.Stdout, "\nTotals")
  fmt.Fprintln(os.Stdout, "\tDirectories: ", dirCount)
  fmt.Fprintln(os.Stdout, "\tFiles: ", fileCount)
  fmt.Fprintln(os.Stdout, "\tBytes Read: ", fileBytesRead)
  fmt.Fprintln(os.Stdout, "\tBytes Written: ", fileBytesWritten)
  fmt.Fprintf(os.Stdout, "\nTo Use it in your Google App Engine program:\n\n\timport \"./pkgs/%s\"\n\n", name)
}

func runDevelopmentServer(path string){
  if len(path) < 1 {
	printHelpCommand("A GAEA project directory must be specified:")
	return
  }
  
  verified := verifyAppServerExists()
  if (!verified){ return }
  fmt.Fprintln(os.Stdout, "Running Dev App Server through GAEA...")
  cmd := exec.Command("dev_appserver.py", path)
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    log.Fatal(err)
  }
  stderr, err := cmd.StderrPipe()
  if err != nil {
    log.Fatal(err)
  }
  if err := cmd.Start(); err != nil {
    log.Fatal(err)
  }
  go io.Copy(os.Stdout, stdout)
  go io.Copy(os.Stderr, stderr)
  if err := cmd.Wait(); err != nil {
    log.Fatal(err)
  }
}

func verifyAppServerExists()(bool){
  cmd := exec.Command("which", "dev_appserver.py")
  output,err := cmd.CombinedOutput()
  if err != nil {
    fmt.Fprintf(os.Stdout, "Error: %s\n", string(output))
    return false
  }
  if len(output) == 0 {
  	fmt.Fprintln(os.Stdout, "Error: Can't find dev_appserver.py in your PATH")
    return false
  }
  return true
}

