package regions

import (
	"fmt"
)

type Region struct {
	Name  string
	Areas []string
}

const (
	DEFAULT_REGION = "rhone"
)

// Returns regions and there areas
func Get() []Region {
	//http://www.insee.fr/fr/methodes/nomenclatures/cog/telechargement.asp
	//http://www.pillot.fr/cartographe/fic_villes.php
	var (
		AlsaceAreas = []string{
			"bas_rhin",
			"haut_rhin",
		}

		AquitaineAreas = []string{
			"dordogne",
			"gironde",
			"landes",
			"lot_et_garonne",
			"pyrenees_atlantiques",
		}

		AuvergneAreas = []string{
			"allier",
			"cantal",
			"haute_loire",
			"puy_de_dome",
		}

		BasseNormandieAreas = []string{
			"calvados",
			"manche",
			"orne",
		}

		BourgogneAreas = []string{
			"cote_d_or",
			"nievre",
			"saone_et_loire",
			"yonne",
		}

		BretagneAreas = []string{
			"cotes_d_armor",
			"finistere",
			"ille_et_vilaine",
			"morbihan",
		}

		CentreAreas = []string{
			"cher",
			"eure_et_loir",
			"indre",
			"indre_et_loire",
			"loir_et_cher",
			"loiret",
		}

		ChampagneArdenneAreas = []string{
			"ardennes",
			"aube",
			"marne",
			"haute_marne",
		}

		CorseAreas = []string{}

		FrancheComteAreas = []string{
			"doubs",
			"jura",
			"haute_saone",
			"territoire_de_belfort",
		}

		HauteNormandieAreas = []string{
			"eure",
			"seine_maritime",
		}

		IleDeFranceAreas = []string{
			"paris",
			"seine_et_marne",
			"yvelines",
			"essonne",
			"hauts_de_seine",
			"seine_saint_denis",
			"val_de_marne",
			"val_d_oise",
		}

		LanguedocRoussillonAreas = []string{
			"aude",
			"gard",
			"herault",
			"lozere",
			"pyrenees_orientales",
		}

		LimousinAreas = []string{
			"correze",
			"creuse",
			"haute_vienne",
		}

		LorraineAreas = []string{
			"meurthe_et_moselle",
			"meuse",
			"moselle",
			"vosges",
		}

		MidiPyreneesAreas = []string{
			"ariege",
			"aveyron",
			"haute_garonne",
			"gers",
			"lot",
			"hautes_pyrenees",
			"tarn",
			"tarn_et_garonne",
		}

		NordPasDeCalaisAreas = []string{
			"nord",
			"pas_de_calais",
		}

		PaysDeLaLoireAreas = []string{
			"loire_atlantique",
			"maine_et_loire",
			"mayenne",
			"sarthe",
			"vendee",
		}

		PicardieAreas = []string{
			"aisne",
			"oise",
			"somme",
		}

		PoitouCharentesAreas = []string{
			"charente",
			"charente_maritime",
			"deux_sevres",
			"vienne",
		}

		ProvenceAlpesCoteDAzurAreas = []string{
			"alpes_de_haute_provence",
			"hautes_alpes",
			"alpes_maritimes",
			"bouches_du_rhone",
			"var",
			"vaucluse",
		}

		RhoneAlpesAreas = []string{
			"ain",
			"ardeche",
			"drome",
			"isere",
			"loire",
			"rhone",
			"savoie",
			"haute_savoie",
		}

		GuadeloupeAreas = []string{}

		MartiniqueAreas = []string{}

		GuyaneAreas = []string{}

		ReunionAreas = []string{}
	)

	return []Region{
		{"alsace", AlsaceAreas},
		{"aquitaine", AquitaineAreas},
		{"auvergne", AuvergneAreas},
		{"basse_normandie", BasseNormandieAreas},
		{"bourgogne", BourgogneAreas},
		{"bretagne", BretagneAreas},
		{"centre", CentreAreas},
		{"corse", CorseAreas},
		{"franche_comte", FrancheComteAreas},
		{"champagne_ardenne", ChampagneArdenneAreas},
		{"haute_normandie", HauteNormandieAreas},
		{"ile_de_france", IleDeFranceAreas},
		{"languedoc_roussillon", LanguedocRoussillonAreas},
		{"limousin", LimousinAreas},
		{"lorraine", LorraineAreas},
		{"midi_pyrenees", MidiPyreneesAreas},
		{"nord_pas_de_calais", NordPasDeCalaisAreas},
		{"pays_de_la_loire", PaysDeLaLoireAreas},
		{"picardie", PicardieAreas},
		{"poitou_charentes", PoitouCharentesAreas},
		{"provence_alpes_cote_d_azur", ProvenceAlpesCoteDAzurAreas},
		{"rhone_alpes", RhoneAlpesAreas},
		{"guadeloupe", GuadeloupeAreas},
		{"martinique", MartiniqueAreas},
		{"guyane", GuyaneAreas},
		{"reunion", ReunionAreas},
	}
}

// Return a formated string that suite with flag package help output
func ToHelpString() string {
	s := ""

	for _, r := range Get() {
		s += "\r\n\t"
		s = s + r.Name + ": "
		a := ""
		for _, area := range r.Areas {
			if a != "" {
				a += ", "
			}
			a += area
		}
		s += a
	}
	return s
}

// Used to confirm a region or and area.
// If regionOrArea is a region, it returns only the region name.
// If regionOrArea is an area, it returns region and area.
func GetRegionAndArea(regionOrArea string) (string, string, error) {
	for _, r := range Get() {
		if regionOrArea == r.Name {
			return r.Name, "", nil
		}

		for _, area := range r.Areas {
			if regionOrArea == area {
				return r.Name, area, nil
			}
		}
	}
	return "", "", fmt.Errorf("Invalid region: '%v'", regionOrArea)
}
