package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type FileType struct {
	Name string `json:"filename"`
	Extension string `json:"fileext"`
	Path string `json:"path"`
	Size int64 `json:"size"`
}

func main(){
	// inDir, err := zenity.SelectFile(
	// 	zenity.Filename(""),
	// 	zenity.Directory(),
	// 	zenity.DisallowEmpty(),
	// 	zenity.Title("Select input directory."),
	// )
	// if err != nil {
	// 	zenity.Error(
	// 		err.Error(),
	// 		zenity.Title("Error"),
	// 		zenity.ErrorIcon,
	// 	)
	// 	log.Fatal(err)
	// }

	// dev only
	inDir := "/Volumes/ns01/users/timotewb/"
	dbFile := "db.sqlite"
	fmt.Println(inDir)

	searchFiles(inDir, dbFile)


}

func searchFiles(inDir string, dbFile string){
	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir(){
			// run function
			fmt.Println(f.Name())
			searchFiles(filepath.Join(inDir, f.Name()), dbFile)

		} else {
			writeToDB(inDir, dbFile, f)
		}
	}
}

func writeToDB(inDir string, dbFile string, f fs.DirEntry){

			// check if database exists
			db, err := sql.Open("sqlite3", dbFile)
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()

			// Create a table if it doesn't exist.
			stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS files (
				id INTEGER PRIMARY KEY,
				name TEXT NOT NULL,
				extension TEXT NOT NULL,
				path TEXT NOT NULL,
				size INTEGER NOT NULL
			)`)
			if err != nil {
				log.Fatal(err)
			}
			_, err = stmt.Exec()
			if err != nil {
				log.Fatal(err)
			}

			fileInfo, err := os.Stat(filepath.Join(inDir, f.Name()))
			if err != nil {
				log.Fatal(err)
			}

			file := FileType{
				Name: f.Name(),
				Extension: filepath.Ext(f.Name()),
				Path: inDir,
				Size: fileInfo.Size(),
			}

			// Prepare the SQL statement.
			stmt, err = db.Prepare("INSERT INTO files (name, extension, path, size) VALUES (?, ?, ?, ?)")
			if err != nil {
				log.Fatal(err)
			}

			// Execute the statement.
			_, err = stmt.Exec(file.Name, file.Extension, file.Path, file.Size)
			if err != nil {
				log.Fatal(err)
			}
}