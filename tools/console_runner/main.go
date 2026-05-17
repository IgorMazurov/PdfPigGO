// Package main provides a console runner for benchmarking PDF document parsing.
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/content"
)

type optionalArg struct {
	shortSymbol   string
	symbol        string
	supportsValue bool
	value         *string
}

type parsedArgs struct {
	suppliedArgs      []*optionalArg
	suppliedDirectoryPath string
}

func getSupportedArgs() []*optionalArg {
	return []*optionalArg{
		{
			shortSymbol:   "nr",
			symbol:        "no-recursion",
			supportsValue: false,
		},
		{
			shortSymbol:   "o",
			symbol:        "output",
			supportsValue: true,
		},
		{
			shortSymbol:   "l",
			symbol:        "limit",
			supportsValue: true,
		},
	}
}

func tryParseArgs(args []string) (*parsedArgs, bool) {
	var path string
	var suppliedOpts []*optionalArg

	opts := getSupportedArgs()

	for i := 0; i < len(args); i++ {
		str := args[i]

		if !strings.HasPrefix(str, "-") {
			if path == "" {
				path = str
			} else {
				return nil, false
			}
		} else {
			var item *optionalArg
			for _, opt := range opts {
				flag1 := "-" + opt.shortSymbol
				flag2 := "--" + opt.symbol
				if strings.EqualFold(flag1, str) || strings.EqualFold(flag2, str) {
					item = opt
					break
				}
			}

			if item == nil {
				return nil, false
			}

			if item.supportsValue {
				if i == len(args)-1 {
					return nil, false
				}
				i++
				v := args[i]
				item.value = &v
			}

			suppliedOpts = append(suppliedOpts, item)
		}
	}

	if path == "" {
		return nil, false
	}

	return &parsedArgs{
		suppliedArgs:          suppliedOpts,
		suppliedDirectoryPath: path,
	}, true
}

type dataRecord struct {
	fileName       string
	size           int64
	words          int
	pages          int
	openCostMicros  int64
	totalCostMicros int64
	perPageMicros   float64
}

func run(args []string) int {
	if len(args) == 0 {
		fmt.Println("At least 1 argument, path to test file directory, must be provided.")
		return 7
	}

	parsed, ok := tryParseArgs(args)
	if !ok {
		strJoined := strings.Join(args, " ")
		fmt.Printf("Unrecognized arguments passed: %s\n", strJoined)
		return 7
	}

	info, err := os.Stat(parsed.suppliedDirectoryPath)
	if err != nil || !info.IsDir() {
		fmt.Printf("The provided path is not a valid directory: %s.\n", parsed.suppliedDirectoryPath)
		return 7
	}

	var maxCount *int
	for _, opt := range parsed.suppliedArgs {
		if opt.shortSymbol == "l" && opt.value != nil {
			v, err := strconv.Atoi(*opt.value)
			if err == nil {
				fmt.Printf("Limiting input files to first: %d\n", v)
				maxCount = &v
			}
		}
	}

	noRecursionMode := false
	for _, opt := range parsed.suppliedArgs {
		if opt.shortSymbol == "nr" {
			noRecursionMode = true
			break
		}
	}

	var outputOpt *optionalArg
	for _, opt := range parsed.suppliedArgs {
		if opt.shortSymbol == "o" && opt.value != nil {
			outputOpt = opt
			break
		}
	}

	fileList, err := findPdfFiles(parsed.suppliedDirectoryPath, noRecursionMode)
	if err != nil {
		fmt.Printf("Error finding PDF files: %v\n", err)
		return 7
	}

	sort.Strings(fileList)

	var runningCount int

	fmt.Printf("Found %d files.\n", len(fileList))
	fmt.Println()

	printTableColumns("File", "Size", "Words", "Pages", "Open cost (μs)", "Total cost (μs)", "Page cost (μs)")

	var dataList []*dataRecord
	var hasError bool
	var errorBuilder strings.Builder

	for _, file := range fileList {
		if maxCount != nil && runningCount >= *maxCount {
			break
		}

		numWords := 0
		numPages := 0
		var openMicros int64
		var totalPageMicros int64

		startTime := time.Now()

		pdfDoc, err := pdfpig.OpenFile(file, &content.ParsingOptions{
			UseLenientParsing: true,
			SkipMissingFonts:  true,
		})

		if err != nil {
			hasError = true
			errorBuilder.WriteString(fmt.Sprintf("Parsing document %s failed due to an error.\n", file))
			errorBuilder.WriteString(err.Error())
			errorBuilder.WriteString("\n")
			runningCount++
			continue
		}

		openMicros = time.Since(startTime).Microseconds()

		pageStartTime := time.Now()

		pages, err := pdfDoc.GetPages()
		if err != nil {
			hasError = true
			errorBuilder.WriteString(fmt.Sprintf("Parsing document %s failed due to an error.\n", file))
			errorBuilder.WriteString(err.Error())
			errorBuilder.WriteString("\n")
			pdfDoc.Close()
			runningCount++
			continue
		}

		for _, pgAny := range pages {
			numPages++
			if pg, ok := pgAny.(*content.Page); ok {
				words := pg.GetWords()
				for _, word := range words {
					if word != nil {
						numWords++
					}
				}
			}
		}

		totalPageMicros = time.Since(pageStartTime).Microseconds()

		pdfDoc.Close()

		filename := filepath.Base(file)

		fileInfo, err := os.Stat(file)
		size := int64(0)
		if err == nil {
			size = fileInfo.Size()
		}

		item := &dataRecord{
			fileName:       filename,
			openCostMicros:  openMicros,
			pages:          numPages,
			size:           size,
			words:          numWords,
			totalCostMicros: totalPageMicros + openMicros,
			perPageMicros:   math.Round(float64(totalPageMicros)/float64(max(numPages, 1))*100) / 100,
		}

		dataList = append(dataList, item)

		printTableColumns(
			item.fileName,
			item.size,
			item.words,
			item.pages,
			item.openCostMicros,
			item.totalCostMicros,
			item.perPageMicros,
		)

		runningCount++
	}

	if hasError {
		fmt.Println(errorBuilder.String())
		return 5
	}

	if outputOpt != nil && outputOpt.value != nil {
		writeOutput(*outputOpt.value, dataList)
	}

	fmt.Println("Complete! :)")

	return 0
}

func findPdfFiles(rootPath string, topLevelOnly bool) ([]string, error) {
	var files []string

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors on individual entries
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".pdf") {
			files = append(files, path)
		}
		if topLevelOnly && path != rootPath && info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	if err := filepath.Walk(rootPath, walkFn); err != nil {
		return nil, err
	}

	return files, nil
}

func writeOutput(outPath string, records []*dataRecord) {
	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open output file: %v\n", err)
		return
	}
	defer f.Close()

	w := &strings.Builder{}
	w.WriteString("File,Size,Words,Pages,Open Cost,Total Cost,Per Page\n")

	for _, record := range records {
		sizeStr := strconv.FormatInt(record.size, 10)
		wordsStr := strconv.Itoa(record.words)
		pagesStr := strconv.Itoa(record.pages)
		openCostStr := strconv.FormatInt(record.openCostMicros, 10)
		totalCostStr := strconv.FormatInt(record.totalCostMicros, 10)
		ppcStr := strconv.FormatFloat(record.perPageMicros, 'f', 2, 64)

		numericParts := []string{sizeStr, wordsStr, pagesStr, openCostStr, totalCostStr, ppcStr}
		w.WriteString(fmt.Sprintf("\"%s\",%s\n", record.fileName, strings.Join(numericParts, ",")))
	}

	f.WriteString(w.String())
}

func printTableColumns(values ...any) {
	for _, value := range values {
		valueStr := fmt.Sprint(value)

		cleaned := getCleanStr(valueStr)

		var padding string
		padChars := 16 - len(cleaned)
		if padChars > 0 {
			padding = strings.Repeat(" ", padChars)
		}

		padded := cleaned + padding

		fmt.Print("| ")
		fmt.Print(padded)
	}

	fmt.Println()
}

func getCleanStr(name string) string {
	maxLength := 16
	if len(name) <= maxLength {
		fillLength := maxLength - len(name)
		return name + strings.Repeat(" ", fillLength)
	}

	return name[:maxLength]
}

func main() {
	os.Exit(run(os.Args[1:]))
}
