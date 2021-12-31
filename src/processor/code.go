package processor

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"

	"gopkg-analyzer/src/utils"
)

func getFirstLine(file, trim string) string {
	// Try opening file
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	// Create scanner
	scanner := bufio.NewScanner(f)
	// Read first line
	_ = scanner.Scan()
	if err := scanner.Err(); err != nil {
		log.Println("Error scanning: ", err)
		return ""
	}

	// Get just the relevant part
	line := scanner.Text()
	nameIdx := strings.Index(line, trim) + len(trim)
	return strings.TrimSpace(line[nameIdx:])
}

func getModName(root string) string {
	const moduleKwd = "module"
	return getFirstLine(root+string(os.PathSeparator)+"go.mod", moduleKwd)
}

func getPackageName(file string) string {
	const packageKwd = "package"
	return getFirstLine(file, packageKwd)
}

func getPackages(goList []string) (res []string) {
	res = []string{}
	for _, fname := range goList {
		r := getPackageName(fname)
		if r == "" {
			log.Printf("File %s doesn't have package name\n", fname)
			continue
		}
		if !utils.Contains(res, r) {
			res = append(res, r)
		}
	}
	return res
}

func AnalyzeCode(rootFolder string) error {
	moduleName := getModName(rootFolder)
	if moduleName == "" {
		return errors.New("go.mod not found")
	}

	// Get mod name
	log.Println(moduleName)
	// Get go files
	fileList := getGoFiles(rootFolder)
	log.Println(fileList)

	// Analyze connections
	dependencyGraph := analyzeImports(fileList)
	log.Printf("\n%s\n", dependencyGraph)

	return nil
}
