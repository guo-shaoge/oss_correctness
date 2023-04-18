package main

import (
	"os"
	"time"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"log"
)

func runSQL(fromDir string, sqlFiles []string, resDir string, resFnSuffix string,
	user, host, port, passwd string) error {
	for i, sqlFn := range sqlFiles {
		sql := fmt.Sprintf("source %s;", path.Join(fromDir, sqlFn))
		logfile := path.Join(resDir, sqlFn + ".result" + resFnSuffix)

		start := time.Now()
		log.Printf("\n\nstart time: %v, iter: %v, sqlFn: %v\n", start, i, sqlFn)

		out, err := exec.Command("mycli", "-u", user, "-h", host, "-P", port, "-D", "gharchive_dev",
	"--ssl-ca", "/etc/ssl/certs/ca-certificates.crt", "--ssl-verify-server-cert", "-p", passwd, "--execute", sql, "--csv").CombinedOutput()
		if err != nil {
			log.Fatalf("failed: %v", string(out))
		}

		msg := fmt.Sprintf("succeed: end time: %v, duration: %v\n\n", time.Now(), time.Since(start))
		out = append([]byte(msg), out...)
		err = os.WriteFile(logfile, out, 0666)
		if err != nil {
			log.Fatalf("got error when write output: %v", err)
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s sql_dir res_dir", os.Args[0])
	}
	fromDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	targetDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		panic(err)
	}

	const timeLayout = "2006_01_02_15_04_05"
	resDir := path.Join(targetDir, time.Now().Format(timeLayout))
	if err := os.Mkdir(resDir, 0755); err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(fromDir)
    if err != nil {
        log.Fatal(err)
    }

	var sqlFiles []string
    for _, file := range files {
		if !file.IsDir() {
			sqlFiles = append(sqlFiles, file.Name())
		}
    }


	// prodCmdTemplate := "mycli -u '3EDFHZJX5iSzvfr.gh_debug' -h gateway01.us-west-2.prod.aws.tidbcloud.com -P 4000 -D test --ssl-ca=/etc/ssl/certs/ca-certificates.crt --ssl-verify-server-cert -pvsPK2GFU4HRAgWVBhoYu --execute %s --csv &> %s"
	// shadowCmdTemplate := "mycli -u '3EDFHZJX5iSzvfr.shadow-ro.c7' -h gateway01.us-west-2.prod.aws.tidbcloud.com -P 4000 -D test --ssl-ca=/etc/ssl/certs/ca-certificates.crt --ssl-verify-server-cert -p1bed8f53a6d716e6a5b5fb1ee28afbd7 --execute %s --csv &> %s "

	runSQL(fromDir, sqlFiles, resDir, ".prod", "3EDFHZJX5iSzvfr.gh_debug", "gateway01.us-west-2.prod.aws.tidbcloud.com", "4000", "vsPK2GFU4HRAgWVBhoYu")
	runSQL(fromDir, sqlFiles, resDir, ".shadow", "3EDFHZJX5iSzvfr.shadow-ro.c7", "gateway01.us-west-2.prod.aws.tidbcloud.com", "4000", "1bed8f53a6d716e6a5b5fb1ee28afbd7")
}
