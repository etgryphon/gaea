package main

import (
	"bufio"
	"fmt"
	su "github.com/etgryphon/stringUp"
	flag "github.com/ogier/pflag"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"errors"
	"html/template"
	"path/filepath"
	"go/parser"
	"go/token"
	"bytes"
)

const APP_VERSION = "0.1"
const inputDelim = '\n'

var fileCount int = 0
var dirCount int = 0
var fileBytesRead int = 0
var fileBytesWritten int = 0

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

// Files
var FileSep = string(os.PathSeparator)

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
	case "help":
		printHelpCommand(nil)
	case "init":
		createNewProject(name)
	case "get":
		getNewImport(name)
	case "run":
		runDevelopmentServer(name)
	default:
		printHelpCommand("Do Not Recognize command: [" + cmd + "]")
	}
}

func printHelpCommand(preamble interface{}) {
	if preamble != nil {
		fmt.Fprintf(os.Stdout, "%s\n%s\n", preamble, HELP_MESSAGE)
	} else {
		fmt.Fprintf(os.Stdout, "%s\n", HELP_MESSAGE)
	}
}

func createNewProject(name string) {
	if len(name) < 1 {
		name = readProjectName()
	}
	camelName := su.CamelCase(name)
	projPath := camelName + "Project"
	fmt.Fprintf(os.Stdout, "Creating a New Project...[%s] in directory: %s\n", camelName, projPath)
	createGAEDirectoryStructure(projPath, camelName)
}

func readProjectName() (name string) {
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

func createGAEDirectoryStructure(path string, name string) {
	var err error

	// Create Base Project Directory
	err = createProjectDirectory("."+FileSep, path, 1)
	if err != nil {
		panic(err)
	}
	// Create app.yml
	projRoot := "."+FileSep+path+FileSep
	err = createProjectFile(name, projRoot, "app.yml", YML_TEMPLATE, 2)
	if err != nil {
		panic(err)
	}

	// Create the main folder
	err = createProjectDirectory(projRoot, name, 2)
	if err != nil {
		panic(err)
	}
	err = createProjectFile(name, projRoot+name+FileSep, name+".go", BASE_GAE_APP_TEMPLATE, 3)
	if err != nil {
		panic(err)
	}

	// Create Public Folder
	err = createProjectDirectory(projRoot, "public", 2)
	if err != nil {
		panic(err)
	}
	// Create Public stylesheets folder
	err = createProjectDirectory(projRoot+"public"+FileSep, "stylesheets", 3)
	if err != nil {
		panic(err)
	}
	// Create Public Images Folder
	err = createProjectDirectory(projRoot+"public"+FileSep, "images", 3)
	if err != nil {
		panic(err)
	}
	// Create Public Javascript Folder
	err = createProjectDirectory(projRoot+"public"+FileSep, "js", 3)
	if err != nil {
		panic(err)
	}
}

func createProjectDirectory(path, name string, level int) error {
	fullPath := path + name
	fmtStr := fmt.Sprintf("%%%ds\n", 1+len(name)+(level*2))
	exists, err := checkIfPathExists(fullPath)
	if !exists {
		fmt.Fprintf(os.Stdout, fmtStr, "+"+name)
		err = os.Mkdir(fullPath, 0755)
	} else {
		fmt.Fprintf(os.Stdout, fmtStr, "-"+name)
	}
	return err
}

func createProjectFile(pkgName, path, name, tmplInstance string, level int) error {
	fmtStr := fmt.Sprintf("%%%ds\n", 1+len(name)+(level*2))
	fileIO, err := os.Create(path + name)
	tmpl, err := template.New(name).Parse(tmplInstance)
	if err != nil {
	  panic(err)
	}
	err = tmpl.Execute(fileIO, pkgName)
	if err != nil {
	  panic(err)
	}
	fileIO.Close()
	fmt.Fprintf(os.Stdout, fmtStr, "+"+name)
	return err
}

func getNewImport(name string) {
	var err error
	if len(name) < 1 {
		printHelpCommand("In order to get a package, you must have a name")
		return
	}
	// check to see if it is present in the GOPATH
	pkgIsThere, srcPath, _ := checkIfPackageIsPresent(name)
	if !pkgIsThere {
		err := fetchExternalPackage(name)
		if err != nil {
			log.Fatalf("External Package Fetch Error:\n\t%s", err)
		}
	}

	// Convert package to local package
	err = convertToLocalPackage(srcPath, name)
	if err != nil {
		log.Fatalf("Package Conversion Error:\n\t%s", err)
	}
}

/*
	@private function to check for the presence of a local package
*/
func checkIfPackageIsPresent(name string) (doesExist bool, srcPath, formattedPath string) {
  doesExist = false
	gopaths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	fmtString := FileSep+"src"+FileSep
	for _,x := range gopaths {
	  srcPath = x+FileSep+"src"
	  path := x + fmtString + name + FileSep
	  formattedPath = filepath.FromSlash(path)
	  fmt.Fprintln(os.Stdout, "dir > ", formattedPath)
	  there, _ := checkIfPathExists(formattedPath)
	  if (there) {
	  	doesExist = true
	  	break
	  }
	}
	return 
}

/*
	@private function to check the presents of a directory
*/
func checkIfPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/*
	@private function for actually fetching the package
*/
func fetchExternalPackage(name string) error {
	log.Println("Fetching External Package: ", name)
	cmd := exec.Command("go", "get", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

func convertFileToLocalUse(gopath, path string, f os.FileInfo, err error) error {
	base := filepath.Base(path)
	if base[0] == '.' {
		return filepath.SkipDir
	}
	curDir, _ := os.Getwd()
	fmt.Fprintln(os.Stdout, "Path >", path)
	newPath := strings.Replace(path, gopath, curDir+FileSep+"pkgs"+FileSep, -1)

	// Now, check what the item is:
	info, err := os.Stat(path)
	if err != nil {
		return filepath.SkipDir
	}
	if info.IsDir() {
		dirCount += 1
		err = os.MkdirAll(newPath, 0755)
		if err != nil {
			return filepath.SkipDir
		}
	} else {
		fileCount += 1
		err = translateFile(gopath, path, newPath)
	}
	return nil
}

/*
	@private function that will copy and translate files to local use
*/
func translateFile(gopath, source, dest string) error {
	s, err := os.Open(source)
	if err != nil {
		return nil
	}
	_, e := os.Stat(dest)
	if e == nil {
		return nil
	}
	d, err := os.Create(dest)
	if err != nil {
		return nil
	}
	sInfo, _ := s.Stat()
	buffer := make([]byte, sInfo.Size())

	_, err = s.Read(buffer)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write read from %s", source))
	}
	
  // If it is a go file then we have to re-write the imports to 
  // have the correct '/pkg' prefix
	if (filepath.Ext(source) == ".go"){
		fset := token.NewFileSet()	// positions are relative to fset
	
		// Parse the file containing this very example
		// but stop after processing the imports.
		f, err := parser.ParseFile(fset, "", buffer, parser.ImportsOnly)
		if err != nil {
			fmt.Println(err)
		}
	
		// Print the imports from the file's AST.
		fmt.Println("File: ", source)
		
		for _,s := range f.Imports {
		  end := len(s.Path.Value)-2
			pkg := s.Path.Value[1:end]
			found, srcPath, _ := checkIfPackageIsPresent(pkg)
			if found {
				err = convertToLocalPackage(srcPath, pkg)
				if err != nil {
					log.Fatalf("Internal Package Conversion Error:\n\t%s", err)
				}
				idx := bytes.Index(buffer, []byte(pkg))
				newBuffer := make([]byte, len(buffer)+4)
				newBuffer = buffer[0:idx]
				newBuffer = append(newBuffer,'p','k','g','s',os.PathSeparator)
				newBuffer = append(newBuffer, buffer[idx:]...)
				buffer = newBuffer
			}
		}
	}

	_, err = d.Write(buffer)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write to %s", dest))
	}
	return nil
}

/*
	@private convert the local package to be used in GAE format
*/
func convertToLocalPackage(root string, name string) error {
	fmt.Println("\n\nTransfering GOPATH package [" + name + "] to local use...")
	fmt.Fprintln(os.Stdout, "GOPATH > ", root)
	curriedConvertFileToLocalUse := func(path string, f os.FileInfo, err error) error {
		return convertFileToLocalUse(root, path, f, err)
	}
	err := filepath.Walk(root+FileSep+name, curriedConvertFileToLocalUse)
	if err != nil {
		return err
	}
	printOutTransferInformation(name)
	return nil
}

func printOutTransferInformation(name string) {
	fmt.Fprintln(os.Stdout, "\nTotals")
	fmt.Fprintln(os.Stdout, "\tDirectories: ", dirCount)
	fmt.Fprintln(os.Stdout, "\tFiles: ", fileCount)
	fmt.Fprintf(os.Stdout, "\nTo Use it in your Google App Engine program:\n\n\timport \".%spkgs%s%s\"\n\n", FileSep, FileSep, name)
}

func runDevelopmentServer(path string) {
	if len(path) < 1 {
		printHelpCommand("A GAEA project directory must be specified:")
		return
	}

	verified := verifyAppServerExists()
	if !verified {
		return
	}
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

func verifyAppServerExists() bool {
	cmd := exec.Command("which", "dev_appserver.py")
	output, err := cmd.CombinedOutput()
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
