package main

import (
	"sort"
	"time"
	"strings"
	"fmt"
	"os"
	"path/filepath"
	"bufio"
	"log"
	"io/ioutil"
	"path"
)

type S struct {
	sql string
	shadowDura time.Duration
	prodDura time.Duration
	rate float64
}

func main() {
	// resDir := "./result"
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s resDir", os.Args[0])
	}
	resDir := os.Args[1]

	fileInfos, err := ioutil.ReadDir(resDir)
	if err != nil {
		panic(err)
	}

	const shadowFnSuffix = ".shadow"
	const prodFnSuffix = ".prod"

	comparedSQL := make(map[string]bool)
	var fnCnt int
	var reports []S
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			log.Printf("got %s, which is not result file", fileInfo)
			continue
		}

		fnBase := fileInfo.Name()
		fnExt := filepath.Ext(fnBase)
		fnNoExt := strings.TrimSuffix(fnBase, fnExt)
		// Already compared.
		if _, ok := comparedSQL[fnNoExt]; ok {
			continue
		}
		comparedSQL[fnNoExt] = true

		shadowResFilePath := path.Join(resDir, fnNoExt + shadowFnSuffix)
		shadowFile, err := os.Open(shadowResFilePath)
		if err != nil {
			panic(err)
		}
		prodResFilePath := path.Join(resDir, fnNoExt + prodFnSuffix)
		prodFile, err := os.Open(prodResFilePath)
		if err != nil {
			panic(err)
		}
		var shadowLines []string
		var prodLines []string
		scanner := bufio.NewScanner(shadowFile)
		for scanner.Scan() {
			shadowLines = append(shadowLines, scanner.Text())
		}
		scanner = bufio.NewScanner(prodFile)
		for scanner.Scan() {
			prodLines = append(prodLines, scanner.Text())
		}

		fnCnt++
		if len(shadowLines) != len(prodLines) {
			msg := fmt.Sprintf("comparing %s and %s, len is: %v vs %v\n", shadowResFilePath, prodResFilePath, len(shadowResFilePath), len(prodLines))
			panic(msg)
		}
		// ignore first line
		same := true
		// duras[fnNoExt] = []string{getDuration(shadowLines[0]), getDuration(prodLines[0])}
		report := S {
			sql: fnNoExt,
			shadowDura: getDuration(shadowLines[0]),
			prodDura: getDuration(prodLines[0]),
		}
		report.rate = report.shadowDura.Seconds()/report.prodDura.Seconds()
		reports = append(reports, report)
		for i := 1; i < len(shadowLines); i++ {
			if shadowLines[i] != prodLines[i] {
				same = false
				// fmt.Printf("result not same %s and %s, line %d, %s vs %s\n", shadowResFilePath, prodResFilePath, i, shadowLines[i], prodLines[i])
				break
			}
		}
		if !same {
			fmt.Printf("vimdiff %s %s\n", fnNoExt + prodFnSuffix, fnNoExt + shadowFnSuffix)
		}
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].rate > reports[j].rate
	})
	fmt.Println()
	fmt.Println()
	fmt.Printf("| %-70s | %-20v | %-20v | %-20v(order by) |\n", "sql fn", "shadow", "prod", "shadow/prod")
	for _, s := range reports {
		fmt.Printf("| %-70s | %-20v | %-20v | %-20v |\n", s.sql, s.shadowDura, s.prodDura, s.rate)
	}
}

func getDuration(line string) time.Duration {
	s := strings.Split(line, " ")
	duraStr := s[len(s)-1]
	dura, err := time.ParseDuration(duraStr)
	if err != nil {
		panic(err)
	}
	return dura
}
