package main

import (
	"context"
	"database/sql"
	//"encoding/binary"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"github.com/twpayne/go-geom/encoding/wkt"
	"io"
	"log"
	"os"
	//"strings"
)

func convert(r *csv.Reader, stmt *sql.Stmt) error {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		g, err := ewkbhex.Decode(record[0])
		if err != nil {
			return err
		}
		switch t := g.(type) {
		case *geom.LineString:
			g2 := geom.LineString(*t)
			coords := [][][]geom.Coord{[][]geom.Coord{g2.Coords()}}
			g = geom.NewMultiPolygon(geom.XY).MustSetCoords(coords)
		case *geom.Polygon:
			g2 := geom.Polygon(*t)
			coords := [][][]geom.Coord{g2.Coords()}
			g = geom.NewMultiPolygon(geom.XY).MustSetCoords(coords)
		case *geom.MultiPolygon: // ignore
		default:
			return fmt.Errorf("cannot handle geometry %v", t)
		}
		//h, err := ewkbhex.Encode(g, binary.LittleEndian)
		h, err := wkt.Marshal(g)
		if err != nil {
			return err
		}
		//tags := strings.ReplaceAll(record[1], `\\`, `\`)
		//_, err = stmt.Exec(h, tags)
		_, err = stmt.Exec(h)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

type entrypoint struct {
	lib  string
	proc string
}

var LibNames = []entrypoint{
	//{"mod_spatialite", "sqlite3_modspatialite_init"},
	//{"mod_spatialite.dylib", "sqlite3_modspatialite_init"},
	//{"libspatialite.so", "sqlite3_modspatialite_init"},
	//{"libspatialite.so.7", "spatialite_init_ex"},
	//{"libspatialite.so", "spatialite_init_ex"},
	//{"/nix/store/0497qxf4msd68yxpfzfkgnficxr84vns-libspatialite-5.0.1/lib/mod_spatialite", "sqlite3_modspatialite_init"},
	{"/nix/store/0497qxf4msd68yxpfzfkgnficxr84vns-libspatialite-5.0.1/lib/libspatialite", "spatialite_init_ex"},
	//{"./libspatialite", "spatialite_init_ex"},
}

var ErrSpatialiteNotFound = errors.New("shaxbee/go-spatialite: spatialite extension not found.")

func register() error {
	sql.Register("spatialite", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			for _, v := range LibNames {
				err := conn.LoadExtension(v.lib, v.proc)
				if err == nil {
					return nil
				}
				fmt.Printf("%v", err)
			}
			return ErrSpatialiteNotFound
		},
	})
	return nil
}

func main() {
	outputDbFileName := flag.String("output", "<undefined>", "file name of the SQLite database to write to")
	flag.Parse()

	//checkErr(register())

	r := csv.NewReader(os.Stdin)
	r.Comma = '\t'
	r.LazyQuotes = true
	r.ReuseRecord = true

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_synchronous=OFF", *outputDbFileName))
	checkErr(err)
	defer db.Close()

	ctx := context.Background()
	opts := sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
		ReadOnly:  false,
	}
	tx, err := db.BeginTx(ctx, &opts)
	checkErr(err)
	defer tx.Commit()

	for _, s := range []string{
		`create table buildings (geometry_wkt text not null)`,
		//`select AddGeometryColumn('buildings', 'geometry', 4326, 'GEOMETRY', 'XY', 1)`,
		//`alter table buildings drop column dummy`,
	} {
		stmt, err := tx.Prepare(s)
		checkErr(err)
		_, err = stmt.Exec()
		checkErr(err)
	}
	//stmt, err := tx.Prepare(`insert into buildings (geometry, tags) values (MultiPolygonFromText(?), ?)`)
	stmt, err := tx.Prepare(`insert into buildings (geometry_wkt) values (?)`)
	checkErr(err)
	err = convert(r, stmt)
	checkErr(err)
	//stmt, err = tx.Prepare(`select CreateSpatialIndex('buildings', 'geometry')`)
	//checkErr(err)
	//_, err = stmt.Exec()
	//checkErr(err)
}
