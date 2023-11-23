# mastr

Werkzeug für den Umgang mit dem Marktstammdatenregister-Gesamtexport.

`mastr` kann den Gesamtdatenexport auslesen, validieren und in eine SQLite-Datenbank umwandeln.

## Website

[![Website](https://img.shields.io/website?up_message=up&down_message=down&url=https%3A%2F%2Fds.marktstammdatenregister.dev%2FMarktstammdatenregister%3Fsql%3Dselect%2B%2522up%2522&style=plastic&label=Datasette)](https://ds.marktstammdatenregister.dev)

Aus diesem Repository heraus wird https://ds.marktstammdatenregister.dev veröffentlicht.

## Download

[![Data](https://zenodo.org/badge/DOI/10.5281/zenodo.10200980.svg)](https://doi.org/10.5281/zenodo.10200980)

Der komplette Export ist in verschiedenen verarbeitungsfreundlichen Formaten [auf Zenodo](https://zenodo.org/records/10200980) zum Download verfügbar.

## Verwendung

Lade zuerst einen [Gesamtdatenexport](https://www.marktstammdatenregister.de/MaStR/Datendownload) herunter.

```
$ go build -o mastr ./cmd/main.go
$ ./mastr -export 'Gesamtdatenexport_<...>.zip' \
          -spec 'spec/Gesamtdatenexport.yaml' \
	  -database 'Marktstammdatenregister.db'
```

Mit diesen Argumenten wandelt `mastr` den heruntergeladenen Gesamtdatenexport in eine SQLite-Datenbank mit dem Namen "Marktstammdatenregister.db" um.
