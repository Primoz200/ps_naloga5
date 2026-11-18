// Package vreme omogoča shranjevanje vremenskih postaj, dodajanje meritev in izpis meritev ter poročil.
package vreme

import (
	"fmt"
	"io"
	"math"
	"sort"
)

// Postaja predstavlja vremensko postajo. Njena polja so zasebna, zato jo
// ustvarite s funkcijo DodajPostajo in z meritvami upravljate prek tega paketa.
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// DodajPostajo v slovar postaj doda postajo z dano oznako, imenom in krajem.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) {
	postaje[oznaka] = Postaja{ime: ime, kraj: kraj}
}

// DodajMeritev doda temperaturo postaji
func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura, minTemp, maxTemp float64) error {
	if math.IsNaN(temperatura) || math.IsInf(temperatura, 0) || temperatura < minTemp || temperatura > maxTemp {
		return fmt.Errorf("neveljavna temperatura %.2f °C (dovoljen razpon %.2f–%.2f °C)", temperatura, minTemp, maxTemp)
	}
	p, ok := postaje[oznaka]
	if !ok {
		return fmt.Errorf("postaja %q ne obstaja", oznaka)
	}
	p.meritve = append(p.meritve, temperatura)
	postaje[oznaka] = p
	return nil
}

// IzpisMeritev izpiše oznake postaj in vse njihove meritve v pisalnik.
func IzpisMeritev(w io.Writer, postaje map[string]Postaja) {
	fmt.Fprintln(w, "MERITVE:")
	for _, oznaka := range oznake(postaje) {
		p := postaje[oznaka]
		fmt.Fprintf(w, "%s - %s %s: %v\n", oznaka, p.ime, p.kraj, p.meritve)
	}
}

// IzpisPorocila izpiše statistično poročilo za postaje, ki imajo vsaj minMeritev meritev. Postaje z manj meritvami dobijo obvestilo.
func IzpisPorocila(w io.Writer, postaje map[string]Postaja, minMeritev int) {
	for _, oznaka := range oznake(postaje) {
		p := postaje[oznaka]
		mi, mx, avg, ok := statistika(p, minMeritev)
		if !ok {
			fmt.Fprintf(w, "%s: premalo meritev\n", p.ime)
			continue
		}
		razred := "Mrzlo"
		switch {
		case avg >= 25:
			razred = "Vroče"
		case avg >= 15:
			razred = "Toplo"
		case avg >= 5:
			razred = "Hladno"
		}
		fmt.Fprintf(w, "%s: min %.2f C, max %.2f C, povprečje %.2f C -> %s\n", p.ime, mi, mx, avg, razred)
	}
}

func statistika(p Postaja, minMeritev int) (najnizja, najvisja, povprecna float64, ok bool) {
	if minMeritev < 1 || len(p.meritve) < minMeritev {
		return 0, 0, 0, false
	}
	najnizja, najvisja = p.meritve[0], p.meritve[0]
	var vsota float64
	for _, v := range p.meritve {
		vsota += v
		if v < najnizja {
			najnizja = v
		}
		if v > najvisja {
			najvisja = v
		}
	}
	return najnizja, najvisja, vsota / float64(len(p.meritve)), true
}

func oznake(postaje map[string]Postaja) []string {
	seznam := make([]string, 0, len(postaje))
	for oznaka := range postaje {
		seznam = append(seznam, oznaka)
	}
	sort.Strings(seznam)
	return seznam
}
