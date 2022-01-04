package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/graphql-go/graphql"
	_ "github.com/mattn/go-sqlite3"
	"pvdb.de/mastr/internal"
	"pvdb.de/mastr/internal/spec"
)

var errMissingOption = errors.New("missing mandatory argument")

func main() {
	err := serve()
	if errors.Is(err, errMissingOption) {
		flag.PrintDefaults()
		os.Exit(64)
	}
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func serve() error {
	const defaultOption = "<undefined>"
	specFileName := flag.String("spec", defaultOption, "file name of the table spec")
	sqliteFile := flag.String("database", defaultOption, "file name of the SQLite database")
	flag.Parse()

	// Ensure mandatory flags are set.
	for _, arg := range []string{
		*specFileName,
		*sqliteFile,
	} {
		if arg == defaultOption {
			return errMissingOption
		}
	}

	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return fmt.Errorf("failed to load location data: %w", err)
	}
	internal.Location = location

	export, err := spec.DecodeExport(*specFileName)
	if err != nil {
		return fmt.Errorf("failed to decode export spec: %w", err)
	}

	// Connect to the database.
	db, err := sql.Open("sqlite3", *sqliteFile)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Create GraphQL schema.
	fields := graphql.Fields{}
	for _, ed := range export {
		td := ed.Table
		args := graphql.FieldConfigArgument{}
		args[td.Primary] = &graphql.ArgumentConfig{
			Type: graphql.String,
		}

		//obj := graphql.NewObject(graphql.ObjectConfig{
		//	Name: td.Element,
		//	Fields: graphql.Fields{
		//		td.Primary: &graphql.Field{
		//			Type: graphql.String,
		//			Args: args,
		//			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//				fmt.Printf("%s = %s\n", td.Primary, p.Args[td.Primary])
		//				println(p.Info.FieldName)
		//				row := db.QueryRow(`select Ort from Marktakteur where Ort is not null limit 1`)
		//				var id string
		//				err := row.Scan(&id)
		//				return id, err
		//			},
		//		},
		//	},
		//})

		obj, err := getTable(db, &td)
		if err != nil {
			return err
		}
		fields[td.Element] = &graphql.Field{
			Type: obj,
			Args: args,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// if _, ok := p.Args[td.Primary]; ok {
				// 	return obj, nil
				// }
				// return td.Primary, nil
				//fmt.Printf("hi %s = %s\n", td.Primary, p.Args[td.Primary])
				//fmt.Printf("Resolve: %s, %v\n%v", td.Primary, p.Args[td.Primary], p)

				primary := p.Args[td.Primary].(string)
				return getObject(db, &td, primary)
			},
		}
	}

	// Schema

	// for k, _ := range fields {
	// 	println(k)
	// }
	// fields := graphql.Fields{
	// 	"hello": &graphql.Field{
	// 		Type: graphql.String,
	// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
	// 			return "world", nil
	// 		},
	// 	},
	// }
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			Marktakteur(MastrNummer: "ABR900000079818") {
				MastrNummer
				DatumLetzeAktualisierung
				Personenart
				Firmenname
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}

	return nil
}

func getTable(db *sql.DB, td *spec.Table) (*graphql.Object, error) {
	f := graphql.Fields{}
	for _, v := range td.Fields {
		f[v.Name] = &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				primary := p.Source.(*graphql.Object).PrivateName
				//fmt.Printf("getTable.Resolve: %v\n", p)
				v, err := getValue(db, td, primary, p.Info.FieldName)
				if err != nil {
					return nil, err
				}
				if v.Valid {
					return v.String, nil
				} else {
					return "", nil
				}
			},
		}
	}
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   td.Element,
		Fields: f,
	}), nil
}

func getValue(db *sql.DB, td *spec.Table, primary string, column string) (sql.NullString, error) {
	row := db.QueryRow(fmt.Sprintf(`select %s from %s where %s = ?`, column, td.Element, td.Primary), primary)
	result := sql.NullString{}
	err := row.Scan(&result)
	return result, err
}

func getObject(db *sql.DB, td *spec.Table, primary string) (*graphql.Object, error) {
	fields := internal.NewFields(td.Fields)
	headers := fields.Header()
	columns := headers[0]
	for _, header := range headers[1:] {
		columns = fmt.Sprintf("%s, %s", columns, header)
	}

	stmt := fmt.Sprintf(`select %s from %s where %s = ?`, columns, td.Element, td.Primary)
	row := db.QueryRow(stmt, primary)
	record := make([]interface{}, len(headers))
	for i, _ := range record {
		record[i] = &sql.NullString{}
	}

	if err := row.Scan(record...); err != nil {
		return nil, fmt.Errorf("getObject scan: %v", err)
	}

	m, err := fields.Map(record)
	if err != nil {
		return nil, fmt.Errorf("getObject new fields: %v", err)
	}

	f := graphql.Fields{}
	for k, v := range m {
		f[k] = &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				fmt.Printf("getObject.Resolve: %v\n", p)
				v2 := v.(sql.NullString)
				if v2.Valid {
					return v2.String, nil
				} else {
					return nil, nil
				}
			},
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   primary,
		Fields: f,
	}), nil
}
