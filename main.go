package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ncruces/zenity"
)

type FileType struct {
	Name      string `json:"filename"`
	Extension string `json:"fileext"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
}

func main() {
	inDir, err := zenity.SelectFile(
		zenity.Filename(""),
		zenity.Directory(),
		zenity.DisallowEmpty(),
		zenity.Title("Select input directory."),
	)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	dbDir, err := zenity.SelectFile(
		zenity.Filename(""),
		zenity.Directory(),
		zenity.DisallowEmpty(),
		zenity.Title("Select Database directory."),
	)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	// executablePath, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }
	// dbDir := filepath.Dir(executablePath)

	dbFile := filepath.Join(dbDir, "go-tools-archive-to-db.sqlite")
	tableName := makeValidTableName(filepath.Base(inDir))
	fmt.Println("--------------------------------------------------")
	fmt.Println("Start processing for:")
	fmt.Printf(" - %s\n", inDir)
	fmt.Printf(" - %s\n", dbFile)
	fmt.Printf(" - %s\n", tableName)
	fmt.Println("--------------------------------------------------")

	zenity.Info(`- `+inDir+`
	 - `+dbFile+`
	 - `+tableName, zenity.Title("Details"))

	// setup database
	// check if database exists
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}
	defer db.Close()

	// drop table if exists
	_, err = db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS " + tableName + "_err")
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	// Create a table if it doesn't exist.
	stmt, err := db.Prepare(`CREATE TABLE ` + tableName + ` (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		extension TEXT NOT NULL,
		path TEXT NOT NULL,
		size INTEGER NOT NULL
	)`)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	// Create a table if it doesn't exist.
	stmt, err = db.Prepare(`CREATE TABLE ` + tableName + `_err (
		id INTEGER PRIMARY KEY,
		file TEXT NOT NULL,
		message TEXT NOT NULL
	)`)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	// Search files
	searchFiles(inDir, dbFile, tableName, db)

	// Complete messages
	zenity.Info("Folder search complete!",
		zenity.Title("Complete"),
		zenity.InfoIcon,
	)

}

func searchFiles(inDir string, dbFile string, tableName string, db *sql.DB) {

	fmt.Println(inDir)

	files, err := os.ReadDir(inDir)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			// run function
			searchFiles(filepath.Join(inDir, f.Name()), dbFile, tableName, db)

		} else {
			writeToDB(inDir, dbFile, tableName, f, db)
		}
	}
}

func writeToDB(inDir string, dbFile string, tableName string, f fs.DirEntry, db *sql.DB) {

	fileInfo, err := os.Stat(filepath.Join(inDir, f.Name()))
	if err != nil {
		errorDB(tableName+"_err", filepath.Join(inDir, f.Name()), err.Error(), db)
		return
	}

	// define and set data
	file := FileType{
		Name:      f.Name(),
		Extension: filepath.Ext(f.Name()),
		Path:      inDir,
		Size:      fileInfo.Size(),
	}

	// insert into db
	stmt, err := db.Prepare("INSERT INTO " + tableName + " (name, extension, path, size) VALUES (?, ?, ?, ?)")
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}
	// Execute the statement.
	_, err = stmt.Exec(file.Name, file.Extension, file.Path, file.Size)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}
}

func makeValidTableName(str string) string {
	// Replace invalid characters with underscore
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	str = reg.ReplaceAllString(str, "_")

	// Ensure the first character is a letter
	if !unicode.IsLetter(rune(str[0])) {
		str = "t_" + str
	}

	return str
}

func errorDB(tableName string, fileNamePath string, errMsg string, db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO " + tableName + " (file, message) VALUES (?, ?)")
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}
	// Execute the statement.
	_, err = stmt.Exec(fileNamePath, errMsg)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}
}
