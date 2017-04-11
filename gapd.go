package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type GAPD_JSON struct {
	packagePaths []string `json:"packagePaths"`
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func printUsage() {
	fmt.Println("Usage: gapd <command> (options)")
	fmt.Println("")
	fmt.Println("gapd init")
	fmt.Println("gapd install <package-path>")
}

func main() {
	dir, _ := os.Getwd()
	workingDirectory := strings.Replace(dir, " ", "\\ ", -1) + "/"
	gapdFilePath := workingDirectory + "gapd.json"

	args := os.Args
	len_args := len(args)

	switch len_args {
	case 0:
		fallthrough
	case 1:
		fmt.Println("Welcome to gapd!")
		fmt.Println("")
		printUsage()
	default:
		command := args[1]

		switch command {
		case "init":
			_, err := os.Open(gapdFilePath)
			if err != nil {
				ioutil.WriteFile(gapdFilePath, []byte("{\n\t\"packagePaths: []\"\n}"), 0644)
				fmt.Println("\"gapd.json\" file is created.")
			} else {
				fmt.Println("\"gapd.json\" file is already exist.")
			}

			fmt.Println("")
			fmt.Println("Maybe you will need following command.")
			fmt.Println(" : gapd install <package-path>")
		case "install":
			gapd_json_byte, err := ioutil.ReadFile(gapdFilePath)
			if err != nil {
				fmt.Println("\"gapd.json\" file is not exist.")
				fmt.Println("")
				fmt.Println("You have to use following command first.")
				fmt.Println(" : gapd init")
			} else if len_args == 2 {
				gapd_json := GAPD_JSON{}
				err := json.Unmarshal(gapd_json_byte, &gapd_json)
				checkError(err)

				for _, path := range gapd_json.packagePaths {
					err = exec.Command("go", "get", path).Run()
					checkError(err)
					fmt.Println(path, "is installed!")
				}

				gapd_json_byte, err = json.Marshal(map[string]interface{}{"packagePaths": gapd_json.packagePaths})
				checkError(err)

				ioutil.WriteFile(gapdFilePath, gapd_json_byte, 0644)
			} else if len_args == 3 {
				path := args[2]

				gapd_json := GAPD_JSON{}
				err := json.Unmarshal(gapd_json_byte, &gapd_json)
				checkError(err)

				for _, exist_path := range gapd_json.packagePaths {
					if exist_path == path {
						err = exec.Command("go", "get", path).Run()
						checkError(err)
						fmt.Println(path, "is installed!")
						return
					}
				}

				gapd_json.packagePaths = append(gapd_json.packagePaths, path)
				err = exec.Command("go", "get", path).Run()
				checkError(err)
				fmt.Println(path, "is installed!")

				gapd_json_byte, err = json.Marshal(map[string]interface{}{"packagePaths": gapd_json.packagePaths})
				checkError(err)

				ioutil.WriteFile(gapdFilePath, gapd_json_byte, 0644)
			} else {
				printUsage()
			}
		default:
			printUsage()
		}
	}
}
