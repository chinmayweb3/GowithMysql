package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	fmt.Println("ths is starting")
	cfg := mysql.Config{
		User:   "root",
		Passwd: "root",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN()) 
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected")

	albumByArtist("John ")
	albumByID(1)

	addAlb := Album{
		ID:     5,
		Title:  "The Dark Side of the Moon",
		Artist: "Pink Floyd",
		Price:  9.99,
	}

	addAlbum(addAlb)
}

// albumByArtist returns all albums by Artist name passed in as a parameter
func albumByArtist(name string) ([]Album, error) {

	var albums []Album
	rows, err := db.Query("SELECT * from album WHERE artist = ?", name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			log.Fatal(err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumByArtist has error: %v", err)
	}
	return albums, nil

}

// albumByID returns an album with ID
func albumByID(id int64) (Album, error) {

	var alb Album
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		log.Fatalln("albumByID error : ", err)
	}
	if err := row.Err(); err != nil {
		return alb, fmt.Errorf("albumByID has error : %v", err)
	}
	return alb, nil

}

// addAlbum adds a new album to the database
func addAlbum(alb Album) (int64, error) {
	fmt.Print("startgin \n")
	result, err := db.Exec("INSERT INTO album ( title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum has error 1 : %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum has error 2 : %v", err)
	}
	fmt.Printf("row has been added %d", id)
	return id, nil

}
