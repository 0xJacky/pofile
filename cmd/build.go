package cmd

/*
Copyright © 2022 0xJacky <me@jackyu.cn>
*/
import (
	"encoding/json"
	"fmt"
	"github.com/0xJacky/pofile/profile"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

var dir string
var file string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build translations.json",
	Long:  `Build translations.json from po file(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		if dir != "" {
			buildFromDir()
		} else if file != "" {
			buildFromFile()
		}
	},
}

func buildFromDir() {
	dict := make(profile.Dict)

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("file to parse %q: %v\n", path, err)
			return err
		}
		if strings.Contains(info.Name(), ".po") {
			p, err := profile.Parse(path)
			if err != nil {
				return err
			}
			// fmt.Println(p.Header.Language)
			dict[p.Header.Language] = p.ToDict()
		}
		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	bytes, err := json.Marshal(dict)
	if err != nil {
		log.Fatal("json marshal error", err)
	}
	err = ioutil.WriteFile("translations.json", bytes, 0644)
	if err != nil {
		log.Fatal("write file error", err)
	}
}

func buildFromFile() {
	dict := make(profile.Dict)
	p, err := profile.Parse(file)
	if err != nil {
		log.Fatalln("[profile parse err]", file, err)
	}
	// fmt.Println(p.Header.Language)
	dict[p.Header.Language] = p.ToDict()
	bytes, err := json.Marshal(dict)
	if err != nil {
		log.Fatal("json marshal error", err)
	}
	err = ioutil.WriteFile("translations.json", bytes, 0644)
	if err != nil {
		log.Fatal("write file error", err)
	}
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	buildCmd.PersistentFlags().StringVarP(&dir, "dir", "d", "", "input a directory")
	buildCmd.PersistentFlags().StringVarP(&file, "file", "p", "", "input a po file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
