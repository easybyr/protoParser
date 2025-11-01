package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	fmt.Println("start to parse proto...")
	var curPath string
	curPath, _ = os.Getwd()
	fmt.Println("current path: " + curPath)

	var workingDir string
	var sourceDir string
	var targetDir string
	flag.StringVar(&workingDir, "d", ".", "待解析proto的目录") // ./paladin-proto
	flag.StringVar(&sourceDir, "s", "proto", "源proto目录")
	flag.StringVar(&targetDir, "t", "java-proto", "proto解析后的存放目录")
	flag.Parse()
	fmt.Printf("working dir=%s; source dir=%s; target dir=%s\n", workingDir, sourceDir, targetDir)

	var f os.FileInfo
	var err error
	f, err = os.Stat(workingDir)
	if err != nil {
		fmt.Println("待解析的proto目录不存在")
		os.Exit(-1)
	}
	if !f.IsDir() {
		fmt.Println("非工作目录")
		os.Exit(-1)
	}

	os.Chdir(workingDir)
	fmt.Printf("Now current dir is: %s\n", workingDir)
	// ~/workcode/libra/paladin-proto

	// 检查 sourceDir
	_, err = os.Stat(sourceDir)
	if err != nil {
		fmt.Printf("proto源目录不存在: %s\n", sourceDir)
		os.Exit(-1)
	}

	// 检查 targetDir
	_, err = os.Stat(targetDir)
	if err != nil {
		fmt.Printf("java proto目标路径不存在: %s，创建新的目标路径\n", targetDir)
		err = os.Mkdir(targetDir, 0777)
		if err != nil {
			fmt.Println("创建目标路径失败，程序退出!")
			os.Exit(-1)
		}
	}

	// 进入到source目录下
	os.Chdir(sourceDir)
	curPath, _ = os.Getwd()
	fmt.Printf("Now I'm in source path: %s\n", curPath)
	fileList, err := ioutil.ReadDir(curPath)
	if err != nil {
		fmt.Println("读取source目录文件异常")
		os.Exit(-1)
	}

	var fileNames []string = make([]string, 0)

	// package: com.gf.libra.pricer.proto
	const javaPackageName = "com.gf.libra.pricer.proto"

	for _, file := range fileList {
		var fileName string = file.Name()
		fileNames = append(fileNames, fileName)

		fmt.Printf("=== start to handle: %s ===\n", fileName)
		fd, fileErr := os.Open(fileName)
		if fileErr != nil {
			continue
		}
		defer fd.Close()

		var writeFilePath string = fmt.Sprintf("../%s/%s", targetDir, fileName)
		fmt.Printf("读取文件: %s，写入文件路径: %s\n", fileName, writeFilePath)
		if !checkFileExists(writeFilePath) {
			os.Create(writeFilePath)
		}

		buff := bufio.NewReader(fd)
		var isPackageField bool = false

		fw, fileWErr := os.OpenFile(writeFilePath, os.O_WRONLY, 0666)
		if fileWErr != nil {
			continue
		}
		defer fw.Close()
		buffWriter := bufio.NewWriter(fw)

		for {
			line, _, eof := buff.ReadLine()
			if eof == io.EOF {
				break
			}
			var lineStr = string(line)
			var lineTrim = strings.TrimSpace(lineStr)
			var finalLine string

			if strings.HasPrefix(lineTrim, "package") {
				isPackageField = true
			}

			if isPackageField {
				// fmt.Printf("package field, origin: %s; now: %s\n", lineStr, javaPackageName)
				finalLine = fmt.Sprintf("package %s;", javaPackageName)
				isPackageField = false
			} else if strings.HasPrefix(lineTrim, "interface") {
				var startIndex int = strings.Index(lineStr, "interface")
				var endIndex int = strings.LastIndex(lineStr, ".")
				finalLine = strings.Replace(lineStr, lineStr[startIndex:endIndex+1], "", 1)
			} else {
				finalLine = lineStr
			}

			// 开始按行写入文件
			_, fileWErr = buffWriter.WriteString(finalLine + "\n")
			if fileWErr != nil {
				fmt.Println("file write err: " + fileWErr.Error())
				panic(fileWErr)
			}
		}

		// 写文件刷新
		buffWriter.Flush()
	}
	fmt.Printf("read source files: %v\n", fileNames)

	fmt.Println("finished!")
}

func checkFileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			panic(err)
		}
	}
	return true
}
