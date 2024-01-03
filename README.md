# Garten Rechnungen erstellen
## Vorbedingungen
Golang installiert [Golang Installation](https://go.dev/doc/install)

## Einrichtung
1. diese Git Projekt klonen
3. Abhängige Dateien in einen Unterordner namens _data_ hinzufügen
  a. _logo.png_ Logo welches als Absender auf den Brief gedruckt werden soll
  b. _mitgliederliste.xlsx_ Excel Datei mit den Daten siehe [Struktur in der mitgliederliste](#struktur-in-mitgliederliste)
  c. _bills_ Ordner in den die Rechnungen als Pdf erstellt werden sollen


#### Struktur im Excel **mitgliederliste.xlsx**
Damit die Daten korrekt aus dem Excel gelesen werden können müssen die folgenden Reiter mit extakt dem angegebenen Namen vorhanden sein

> Ein Beispiel Excel mit der geforderten Struktur ist im Ordner [_example_data_](./example_data/mitgliederliste.xlsx) hinterlegt.

##### _Mitgliederliste_
Die folgende Tabelle Zeigt die Mindestandorderung an das Tabellenblatt _Mitgliederliste_. Wichtig ist vor allem die Reihenfolge.
|Parz.	        | Name | Vorname | Adresse | PLZ | Ort | Tel. | Aa | | |Spr. | Vorstand |
|----------|------|---------|---------|--------|-----|------------|-----------|--- |---|-----------|----------|		
|Parzellen nummer|Name als text|Vorname als Text|Strasse und Hausnummer|Postleitzahl|Ortschaftsnamen|Telefon (wird nicht benötigt)|Anzahl Aren als Zahlenwert|Leere Spalte|Leere Spalte|Sprache 'D' oder 'F' |'J' falls Mitglied Vorstansmitglied ist ansonsten leer|	

##### _Betraege_
Diese Tabelle _muss_ explizit so ausgefüllt werden
alle Einträge sind in derselben Währung als Zahl anzugeben. Die ersten beiden Zeilen Werden für die Übersetzungen der Bezeichnung verwendet. Die erste Zeile ist Deutsch und die zweite Französisch.
In der dritten Zeile sind dann die Beträge in Franken angegeben.

|Pachtzins | Wasserbezug	| Abonnement GF | Strom | Versicherung | Mitgliederbeitrag | Reparatur Fonds | Verwaltungskosten |
|---------|--------------|--------------|------|--------|-----------|-----------|-----------|
|Loyer de la parcelle |	Consommation d'eau |	Abonnement du Journal |	elctricité | Assurance |	Cotisation | Fonds de réparation | frais de gestion|
| Zins pro Are| Wasserkosten pro Are|Kosten für das Abonement des "Gartenfreund"|Stromkosten pauschal| Versicherungskosten pauschal |Mitgliederbeitrag pauschal|Beitrag an den Reparaturfonds pauschal|Beitrag an die Verwaltungskosten pauschal|

##### _Rechnungsdetails_
Auch diese Tabelle _muss_ explizit so ausgefüllt werden
Hierbei sind die Spalten massgebend. Die erste Spalte ist die Bezeichnung, während die zweite Spalte die Werte beinhaltet.

Ab Zeile 7 werden wieder die Übersetzungen eingetragen. Hier handelt es sich um die Spaltenüberschriften und die Mengenangaben.

| | | De | Fr |
|---- | ---- | ---- | ---- |
|Name | Name in der Adresszeile|
| Adresse |	Adresse ohne Hausnummer |
| Adress Nummer | Hausnummer |
| Postleitzahl | Postleitzahl |
| Stadt | Ortsangabe ohne plz |
|Iban Nummer | gültige Iban Nummer |
| Ueberschrift |  | Überschrift im Brief auf Deutsch| Überschrift im Brief auf Französisch |
| Tabelle Anzahl | | Anzahl | Nombre |
| Tabelle Einheit | | Einheit | Unité |
| Tabelle Bezeichnung | | Bezeichnung | Dénomination |
| Tabelle Preis | | Preis | Prix |
| Tabelle Betrag | | Betrag | Montant |
| Tabelle Aaren | | Aren | Are |
| Tabelle Jahre | | Jahr | Année |


## Anwendung

### Build für eine ausführbare Datei
- build für das Gerät auf dem es Ausgeführt wird: `go build`
- Für Windows `GOOS="windows" GOARCH="amd64" go build` 

Die Ausführbare Datei wird hierdurch erstellt und heisst: "qr-invoice"
Diese kann auf den jeweiligen Systemen (und im Ordner  der auch die Unterordner _data_ und _bills_ enthält) direkt ausgeführt werden.

### Programm mit Go laufen lassen
> Vorbedingung hierfür sind die benötigten Dateien im Ordner [Vorbedingungen](#einrichtung) _data_

Programm zu starten `go run main.go`
Nach erfolgreichem Lauf befinden sich die Rechnungen unter _bills_


