package processor

import (
	"bufio"
	"log"
	"os"
	"strings"

	"gopkg-analyzer/src/structs"
	"gopkg-analyzer/src/utils"
)

func parseImports(f *bufio.Scanner) (res []string) {
	// Skip first line
	_ = f.Scan()

	res = []string{}
	const importKwd = "import"
	insideImports := false
	for f.Scan() {
		if err := f.Err(); err != nil {
			log.Printf("Error parsing imports: %s\n", err)
			return
		}

		line := f.Text()

		if line == "" {
			continue
		}

		if !insideImports {
			if !strings.Contains(line, importKwd) {
				return
			}

			insideImports = strings.Contains(line, importKwd)

			// if entered, check if contains ( or "
			if insideImports {
				if strings.Contains(line, "(") {
					// Entered imports, just go to next line
					continue
				} else if strings.Contains(line, "\"") {
					// One line of import
					s := strings.Index(line, "\"")
					e := strings.LastIndex(line, "\"")
					res = append(res, line[s+1:e])
					insideImports = false
					continue
				} else {
					// Something went wrong
					log.Printf("What in tarnation??? %s\n", line)
					return
				}
			}
		} else {
			if strings.Contains(line, ")") {
				insideImports = false
				continue
			}

			if !strings.Contains(line, "\"") {
				continue
			}

			s := strings.Index(line, "\"")
			e := strings.LastIndex(line, "\"")
			res = append(res, line[s+1:e])
		}
	}
	return res
}

func getImportList(fname string) (res []string) {
	res = []string{}

	// Try opening file
	f, err := os.Open(fname)
	if err != nil {
		return res
	}
	defer f.Close()

	// Create scanner
	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		log.Println("Error scanning: ", err)
		return res
	}

	// Parse imports
	res = parseImports(scanner)
	log.Println(res)

	return res
}

func analyzeImports(fileList []string) *structs.SquareGraph {
	// Get unique package names
	pkgs := getPackages(fileList)
	log.Println(pkgs)

	// Create graph
	graph := structs.NewGraph(len(pkgs)).HideDiagonal().AddKeys(pkgs, true)

	// Analyze the imports of each file
	for _, fname := range fileList {
		pkg := getPackageName(fname)
		importList := getImportList(fname)
		for _, imp := range importList {
			// Remove path
			imp = imp[strings.LastIndex(imp, "/")+1:]

			// Skip if the import is from a third party
			// TODO INCLUDE IMPORTS OF THIRD PARTY
			if !utils.Contains(pkgs, imp) {
				continue
			}

			// Conect nodes on the graph
			graph.Connect(pkg, imp)
		}
	}

	return graph
}
