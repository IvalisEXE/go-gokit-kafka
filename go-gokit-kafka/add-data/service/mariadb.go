package service

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" //mysql driver
)

const (
	addDataPegawai = `INSERT INTO pegawai (id, NamaDepan, NamaBelakang, Alamat) VALUES (?,?,?,?);`
)

type dbReadWriter struct {
	db *sql.DB
}

//NewDBReadWriter is ..
func NewDBReadWriter(url string, schema string, user string, password string) ReadWriter {
	schemaURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, url, schema)
	db, err := sql.Open("mysql", schemaURL)
	if err != nil {
		panic(err)
	}
	return &dbReadWriter{db: db}
}

func (rw *dbReadWriter) WriteDataPegawai(pe Pegawai) error {

	log.Println("<-><-><-><-><-><->SAVING TO DATABASE DISPATCH ")

	tx, err := rw.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, er := tx.Exec(addDataPegawai, pe.ID, pe.NamaDepan, pe.NamaBelakang, pe.Alamat)
	if er != nil {
		return er
	}

	return tx.Commit()
}
