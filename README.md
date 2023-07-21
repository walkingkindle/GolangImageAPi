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



