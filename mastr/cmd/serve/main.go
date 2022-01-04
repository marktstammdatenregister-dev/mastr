package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
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
		obj, err := getTable(db, &td)
		if err != nil {
			return err
		}
		f := internal.NewFields(td.Fields)
		t, err := f.GraphqlType(td.Primary)
		if err != nil {
			return err
		}
		fields[td.Element] = &graphql.Field{
			Type: obj,
			Args: graphql.FieldConfigArgument{
				td.Primary: &graphql.ArgumentConfig{
					Type: t,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				primary := p.Args[td.Primary]
				fmt.Printf("top level Resolve: %s\n", primary)
				p.Context = context.WithValue(context.Background(), "toplevel", primary)
				return getObject(db, &td, primary)
			},
		}
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: fields}),
	})
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
			EinheitWind(EinheitMastrNummer: "SEE900002935310") {
				EinheitMastrNummer
				DatumLetzteAktualisierung
				Ort
				Hausnummer
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

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/graphql", h)
	return http.ListenAndServe(":8080", nil)
}

func getTable(db *sql.DB, td *spec.Table) (*graphql.Object, error) {
	f := graphql.Fields{}
	for _, field := range td.Fields {
		fields := internal.NewFields(td.Fields)
		t, err := fields.GraphqlType(field.Name)
		if err != nil {
			return nil, err
		}
		f[field.Name] = &graphql.Field{
			Type: t,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				//primary := p.Source.(*graphql.Object).PrivateName
				//primary := p.Args[td.Primary].(string)
				primary := p.Context.Value("toplevel").(string)
				fmt.Printf("getTable.Resolve: %s\n", primary)
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
	result := sql.NullString{}
	stmt, err := db.Prepare(fmt.Sprintf(`select %s from %s where %s = ?`, column, td.Element, td.Primary))
	if err != nil {
		return result, err
	}
	defer stmt.Close()
	fmt.Println(stmt)
	row := stmt.QueryRow(primary)
	err = row.Scan(&result)
	return result, err
}

func getObject(db *sql.DB, td *spec.Table, primary interface{}) (*graphql.Object, error) {
	// fields := internal.NewFields(td.Fields)
	// headers := fields.Header()
	// columns := headers[0]
	// for _, header := range headers[1:] {
	// 	columns = fmt.Sprintf("%s, %s", columns, header)
	// }

	// // TODO: Prepare one statement
	// stmt := fmt.Sprintf(`select %s from %s where %s = ?`, columns, td.Element, td.Primary)
	// row := db.QueryRow(stmt, primary)
	// dest := fields.ScanDest()
	// if err := row.Scan(dest...); err != nil {
	// 	return nil, fmt.Errorf("getObject scan: %v", err)
	// }

	// m, err := fields.Map(dest)
	// if err != nil {
	// 	return nil, fmt.Errorf("getObject new fields: %v", err)
	// }

	// f := graphql.Fields{}
	// for k, v := range m {
	// 	t, err := fields.GraphqlType(k)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	f[k] = &graphql.Field{
	// 		Type: t,
	// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
	// 			fmt.Printf("getObject.Resolve: %v\n", p)
	// 			v2 := v.(sql.NullString)
	// 			if v2.Valid {
	// 				return v2.String, nil
	// 			} else {
	// 				return nil, nil
	// 			}
	// 		},
	// 	}
	// }

	return graphql.NewObject(graphql.ObjectConfig{
		Name: fmt.Sprintf("%v", primary),
		// TODO: Try adding Resolve here -- perhaps that's how we can achieve caching?
		//Fields: f,
	}), nil
}
