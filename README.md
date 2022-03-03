# mastr

[![ds.marktstammdatenregister.dev](https://img.shields.io/website?down_color=lightgrey&down_message=down&label=datasette&style=flat-square&up_color=green&up_message=up&url=https%3A%2F%2Fds.marktstammdatenregister.dev%2FMarktstammdatenregister%3Fsql%3Dselect%2B%2527gh%2527)](https://ds.marktstammdatenregister.dev)

Werkzeug für den Umgang mit dem Marktstammdatenregister-Gesamtexport.

`mastr` kann den Gesamtdatenexport auslesen, validieren und in eine SQLite-Datenbank umwandeln.

Aus diesem Repository heraus wird https://ds.marktstammdatenregister.dev veröffentlicht.

## Verwendung

Lade zuerst einen [Gesamtdatenexport](https://www.marktstammdatenregister.de/MaStR/Datendownload) herunter.

```
$ go build -o mastr ./cmd/main.go
$ ./mastr -export 'Gesamtdatenexport_<...>.zip' \
          -spec 'spec/Gesamtdatenexport.yaml' \
	  -database 'Marktstammdatenregister.db'
```

Mit diesen Argumenten wandelt `mastr` den heruntergeladenen Gesamtdatenexport in eine SQLite-Datenbank mit dem Namen "Marktstammdatenregister.db" um.
