package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Tag struct {
	Name string `json:"name"`
}

func main() {
	apiURL := "https://api.github.com/repos/tukui-org/ElvUI/tags"

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Error while receive tags : %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error while receive tags. Code statut: %d\n", resp.StatusCode)
		time.Sleep(3 * time.Second)
		return
	}

	var tags []Tag
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		fmt.Printf("Error on reading JSON response : %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}

	if len(tags) == 0 {
		fmt.Println("No tag found on repo.")
		time.Sleep(3 * time.Second)
		return
	}

	lastTag := tags[0].Name

	zipURL := fmt.Sprintf("https://github.com/tukui-org/ElvUI/archive/%s.zip", lastTag)

	zipResp, err := http.Get(zipURL)
	if err != nil {
		fmt.Printf("Error while downloading zip file : %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}
	defer zipResp.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s.zip", lastTag))
	if err != nil {
		fmt.Printf("Error while creating local file : %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, zipResp.Body)
	if err != nil {
		fmt.Printf("Error while copying zip files : %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}

	fmt.Printf("Zip file on latest tag (%s) downloaded with success.\n", lastTag)
	fmt.Println("Starting decompressing file ...")

	err = unzip(fmt.Sprintf("%s.zip", lastTag), "./AddOns/")
	if err != nil {
		fmt.Printf("Erreur lors de la decompression du fichier zip: %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}

	fmt.Println("Latest version of ElvUI installation is succed !")
	time.Sleep(3 * time.Second)
}

func unzip(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("Error while opening zip file : %s", err)
	}
	defer reader.Close()

	commonPrefix := findCommonPrefix(reader.File)
	if commonPrefix == "" {
		return fmt.Errorf("None common prefix found")
	}

	for _, file := range reader.File {
		if !strings.HasPrefix(file.Name, commonPrefix) {
			continue
		}

		// Build destination path with common prefix removed
		relPath, err := filepath.Rel(commonPrefix, file.Name)
		if err != nil {
			return fmt.Errorf("Error while wrint relatif path %s : %s", file.Name, err)
		}

		path := filepath.Join(dest, relPath)

		if file.FileInfo().IsDir() {
			// Create recursivly dir if not exist
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return fmt.Errorf("Error while creating dir %s : %s", path, err)
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("Error while opening files %s in zip : %s", file.Name, err)
		}
		defer fileReader.Close()

		// Create recursivly dir if not exist
		dir := filepath.Dir(path)
		os.MkdirAll(dir, os.ModePerm)

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("Error while create file %s : %s", path, err)
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return fmt.Errorf("%s : %s", file.Name, err)
		}
	}

	// Remove all files and dirs except ElvUI, ElvUI_Libraries and ElvUI_Options
	if err := cleanUpExcept(dest, []string{"ElvUI", "ElvUI_Libraries", "ElvUI_Options"}); err != nil {
		return fmt.Errorf("Error when trying to delete useless file : %s", err)
	}

	// Remove zip file
	err = os.Remove(zipFile)
	if err != nil {
		return fmt.Errorf("Error while deleting zip file : %s", err)
	}

	return nil
}

func cleanUpExcept(dir string, keep []string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullPath := filepath.Join(dir, file.Name())

		if !contains(keep, file.Name()) {
			if file.IsDir() {
				if err := os.RemoveAll(fullPath); err != nil {
					return fmt.Errorf("Error when trying to delete directroy %s : %s", fullPath, err)
				}
			} else {
				if err := os.Remove(fullPath); err != nil {
					return fmt.Errorf("Error when trying to delete file %s : %s", fullPath, err)
				}
			}
		}
	}

	return nil
}

func contains(list []string, item string) bool {
	for _, val := range list {
		if val == item {
			return true
		}
	}
	return false
}

func findCommonPrefix(files []*zip.File) string {
	if len(files) == 0 {
		return ""
	}
	prefix := files[0].Name
	for _, file := range files[1:] {
		for i := 0; i < len(prefix) && i < len(file.Name); i++ {
			if prefix[i] != file.Name[i] {
				prefix = prefix[:i]
				break
			}
		}
	}
	return prefix
}
