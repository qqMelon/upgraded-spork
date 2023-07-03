package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Tag struct {
	Name string `json:"name"`
}

func main() {
	// URL de l'API GitHub pour récupérer les tags
	apiURL := "https://api.github.com/repos/tukui-org/ElvUI/tags"

	// Récupération des tags
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Erreur lors de la récupération des tags: %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}
	defer resp.Body.Close()

	// Vérification du code de statut HTTP
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erreur lors de la récupération des tags. Code de statut: %d\n", resp.StatusCode)
		time.Sleep(3 * time.Second)
		return
	}

	// Lecture de la réponse JSON
	var tags []Tag
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture de la réponse JSON: %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}

	// Vérification si des tags sont disponibles
	if len(tags) == 0 {
		fmt.Println("Aucun tag trouvé pour le référentiel.")
		time.Sleep(3 * time.Second)
		return
	}

	// Récupération du dernier tag
	lastTag := tags[0].Name

	// Construction de l'URL du fichier zip
	zipURL := fmt.Sprintf("https://github.com/tukui-org/ElvUI/archive/%s.zip", lastTag)

	// Téléchargement du fichier zip
	zipResp, err := http.Get(zipURL)
	if err != nil {
		fmt.Printf("Erreur lors du téléchargement du fichier zip: %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}
	defer zipResp.Body.Close()

	// Création du fichier local pour enregistrer le zip
	file, err := os.Create(fmt.Sprintf("%s.zip", lastTag))
	if err != nil {
		fmt.Printf("Erreur lors de la création du fichier local: %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}
	defer file.Close()

	// Copie du contenu du fichier zip téléchargé vers le fichier local
	_, err = io.Copy(file, zipResp.Body)
	if err != nil {
		fmt.Printf("Erreur lors de la copie du contenu du fichier zip: %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}

	fmt.Printf("Le fichier zip du dernier tag (%s) a été téléchargé avec succès.\n", lastTag)
	fmt.Println("Debut de la decompression du fichier ...")

	err = unzip(fmt.Sprintf("%s.zip", lastTag), "./AddOns/")
	if err != nil {
		fmt.Printf("Erreur lors de la decompression du fichier zip: %s\n", err)
		time.Sleep(3 * time.Second)
		return
	}

	fmt.Println("L'installation de la derniere version de ElVui est un succes !")
	time.Sleep(3 * time.Second)
}

func unzip(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}

	return nil
}
