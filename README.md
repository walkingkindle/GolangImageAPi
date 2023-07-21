# GolangImageAPi

Zadatak 1:

Napisati API upotrebom standardne biblioteke ili u frameworku po izboru koji ima sledece mogucnosti:
    * upload slike u folder
        * ako slika postoji odbijte upload razumnom porukom
        * naziv slike je SHA256(slike) tako da cemo onemoguciti duplikate
    * listanje uploadovanih slika gde je izlaz JSON output sa nazivima slika.
    * sortiraj nazive A-Z rastuce
    * brisanje odabrane slike, brise se po nazivu slike
    * download odabrane slike, preuzima se po nazivu slike

Pitanja na koja treba dati odgovore:
    * koji su potencijalni izazovi ovakvog mikroservisa u cilju skaliranja?
    * kako bi obradio velike fajlove od po 1-2GB?

Dodatno (nije obavezno):
    * Go unit testovi
    * Dockerizovati mikroservis



------------------------------------------------------------------------------------------------------------------------
1. Pokrenuti Goimgapi.exe
2. U browseru ući na localhost:8080


    * kako bi obradio velike fajlove od po 1-2GB?
Odgovor: Pretvorio bi fail-ove u manje chunk-ove ili koristio multipart/form-data kao html upload ulaz.
Ako bi pretvorio u manje chunk-ove ovako bi otprilike izgledao moj kod:
	buf := make([]byte, 4*1024) // 4KB buffer size
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading file:", err)
			}
			break
		}
		_, err = out.Write(buf[:n])
		if err != nil {
			fmt.Println("Error writing file:", err)
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

    * koji su potencijalni izazovi ovakvog mikroservisa u cilju skaliranja?

Odgovor:
Za ovakvu vrstu izazova koristio bi nesto kao horizontalno skaliranje, tj veci broj servera koji bi delili workload izmedju fajlova, regularno pratio perfomanse programa optimizovao infrastrukturu da ne bi dolazilo do bagova, ili koristio neku vrstu cloud resenja i prostora na cloud sistemima koji rastu bez manuelne intervencije.
