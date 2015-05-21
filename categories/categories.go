package categories

const (
	// DefaultCategoryIndex ...
	DefaultCategoryIndex = 0 // _vehicules_
)

// Get ...
func Get() []string {
	return []string{
		"_vehicules_",
		"voitures",
		"motos",
		"_immobilier_",
		"_multimedia_",
		"_maison_",
		"_loisirs_",
		"_materiel_professionnel_",
		"_emploi_services_",
		"_",
		"autres",
	}
}

// IndexOf ...
func IndexOf(v string) int {
	for i, category := range Get() {
		if v == category {
			return i
		}
	}
	return -1
}
