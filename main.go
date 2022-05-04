package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	// "go_replace_extended/helper" // TODO: see if we can build with the original code with something link helper.Main()
)

var (
	replaceMe   = os.Getenv("REPLACE_ME")   // "MY_AWS_REGION,ap-southeast-2|INTERNAL_ZONE_ID,YYYYYYYY"
	installApps = os.Getenv("INSTALL_APPS") // "datadog,kyverno,gp3"
)

func RunShellCmd(command string) (stdout string, stderr string) {
	cmd := exec.Command("sh", "-c", command)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()
	return outb.String(), errb.String()
}

func replaceStrings() {
	// Replace Strings in the file
	replaceMeList := strings.Split(replaceMe, "|")
	for _, wordSet := range replaceMeList {
		findReplace := strings.Split(wordSet, ",")
		cmd := fmt.Sprintf("gor %s -r %s", findReplace[0], findReplace[1])
		stdout, stderr := RunShellCmd(cmd)
		fmt.Println(stdout)
		if stderr != "" {
			log.Fatalf("ERROR - Command %v failed with error %v \n", cmd, stderr)
		}
	}
	// log.Print("All strings replaced")
}

// https://stackoverflow.com/a/67629473
func findFiles(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func removeFiles() {
	// Remove the files not in appsList
	appsList := strings.Split(installApps, ",")

	for _, appFile := range findFiles(".", ".yaml") {
		appName := strings.Split(appFile, ".")[0]
		if contains(appsList, appName) == false {
			e := os.Remove(appFile)
			if e != nil {
				log.Fatal(e)
			} else {
				// log.Printf("File %v removed", appFile)
			}
		}
	}
}

func main() {
	if replaceMe == "" {
		log.Fatalln("ERROR - REPLACE_ME env variable cannot be empty")
	}

	if installApps != "" {
		removeFiles()
	}

	replaceStrings()
}
