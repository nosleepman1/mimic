package faker

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	firstNames = []string{
		"Mamadou", "Fatou", "Ibrahima", "Awa", "Ousmane", "Aminata", "Cheikh", "Mariama",
		"Abdoulaye", "Khady", "Modou", "Ndeye", "Samba", "Rokhaya", "Alioune", "Coumba",
		"Babacar", "Adama", "Moussa", "Seynabou", "Boubacar", "Khadija", "El Hadji", "Binta",
	}

	lastNames = []string{
		"Diop", "Ndiaye", "Sow", "Fall", "Faye", "Diallo", "Ba", "Gueye",
		"Sy", "Cisse", "Traore", "Kane", "Sarr", "Thiam", "Mbodj", "Seck",
		"Dione", "Kone", "Camara", "Badiane", "Niang", "Toure", "Dia", "Tall",
	}

	cities = []string{
		"Dakar", "Touba", "Thies", "Saint-Louis", "Kaolack", "Ziguinchor",
		"Mbour", "Rufisque", "Tambacounda", "Kolda", "Fatick", "Louga",
	}

	companyNames = []string{
		"Teranga Tech", "Dakar Innovations", "Sunu Labs", "Baobab Systems", "Sahel Solutions",
		"SenPay Digital", "Walo Logistics", "Sine Consulting", "Nianing Cloud",
	}

	departments = []string{
		"Engineering", "Design", "Marketing", "Sales", "Customer Support", "Finance", "DevOps",
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomFullName génère un nom et prénom réalistes
func RandomFullName() string {
	fn := firstNames[rand.Intn(len(firstNames))]
	ln := lastNames[rand.Intn(len(lastNames))]
	return fmt.Sprintf("%s %s", fn, ln)
}

// RandomPhoneSN génère un numéro de téléphone au format sénégalais (+221 77/78/76/70...)
func RandomPhoneSN() string {
	prefixes := []string{"77", "78", "76", "70", "75"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	num := rand.Intn(9000000) + 1000000
	return fmt.Sprintf("+221 %s %03d %02d %02d", prefix, num/10000, (num/100)%100, num%100)
}

// RandomCurrencyXOF génère un montant en Francs CFA (XOF)
func RandomCurrencyXOF() int {
	amounts := []int{1000, 2500, 5000, 10000, 15000, 25000, 50000, 100000, 250000, 500000}
	return amounts[rand.Intn(len(amounts))]
}

// RandomCity retourne une ville
func RandomCity() string {
	return cities[rand.Intn(len(cities))]
}

// RandomCompany retourne un nom d'entreprise
func RandomCompany() string {
	return companyNames[rand.Intn(len(companyNames))]
}

// RandomDepartment retourne un département
func RandomDepartment() string {
	return departments[rand.Intn(len(departments))]
}
