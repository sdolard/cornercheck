package regions

import (
	"fmt"
	"strings"
	"testing"
)

func TestToHelpString(t *testing.T) {
	want := []string{
		"\r\n\talsace: bas_rhin, haut_rhin",
		"\r\n\taquitaine: dordogne, gironde, landes, lot_et_garonne, pyrenees_atlantiques",
		"\r\n\tauvergne: allier, cantal, haute_loire, puy_de_dome",
		"\r\n\tbasse_normandie: calvados, manche, orne",
		"\r\n\tbourgogne: cote_d_or, nievre, saone_et_loire, yonne",
		"\r\n\tbretagne: cotes_d_armor, finistere, ille_et_vilaine, morbihan",
		"\r\n\tcentre: cher, eure_et_loir, indre, indre_et_loire, loir_et_cher, loiret",
		"\r\n\tcorse: ",
		"\r\n\tfranche_comte: doubs, jura, haute_saone, territoire_de_belfort",
		"\r\n\tchampagne_ardenne: ardennes, aube, marne, haute_marne",
		"\r\n\thaute_normandie: eure, seine_maritime",
		"\r\n\tile_de_france: paris, seine_et_marne, yvelines, essonne, hauts_de_seine, seine_saint_denis, val_de_marne, val_d_oise",
		"\r\n\tlanguedoc_roussillon: aude, gard, herault, lozere, pyrenees_orientales",
		"\r\n\tlimousin: correze, creuse, haute_vienne",
		"\r\n\tlorraine: meurthe_et_moselle, meuse, moselle, vosges",
		"\r\n\tmidi_pyrenees: ariege, aveyron, haute_garonne, gers, lot, hautes_pyrenees, tarn, tarn_et_garonne",
		"\r\n\tnord_pas_de_calais: nord, pas_de_calais",
		"\r\n\tpays_de_la_loire: loire_atlantique, maine_et_loire, mayenne, sarthe, vendee",
		"\r\n\tpicardie: aisne, oise, somme",
		"\r\n\tpoitou_charentes: charente, charente_maritime, deux_sevres, vienne",
		"\r\n\tprovence_alpes_cote_d_azur: alpes_de_haute_provence, hautes_alpes, alpes_maritimes, bouches_du_rhone, var, vaucluse",
		"\r\n\trhone_alpes: ain, ardeche, drome, isere, loire, rhone, savoie, haute_savoie",
		"\r\n\tguadeloupe: ",
		"\r\n\tmartinique: ",
		"\r\n\tguyane: ",
		"\r\n\treunion: ",
	}
	if get := ToHelpString(); strings.Join(want, "") != get {
		t.Errorf("ToHelpString() = '%v', want '%v'", get, want)
	}
}

func TestGetRegionAndArea(t *testing.T) {
	wants := []struct {
		region       string
		regionOrArea string
	}{
		{"alsace", "bas_rhin"},
		{"rhone_alpes", "haute_savoie"},
		{"reunion", "reunion"},
	}

	for _, wanted := range wants {
		if region, _, _ := GetRegionAndArea(wanted.regionOrArea); region != wanted.region {
			t.Errorf("TestGetRegionAndArea() = '%v', want '%v'", region, wanted.region)
		}
	}

	invalidRegion := "an invalid region"
	errMessage := fmt.Sprintf("Invalid region: '%v'", invalidRegion)
	if _, _, err := GetRegionAndArea(invalidRegion); err.Error() != errMessage {
		t.Errorf("TestGetRegionAndArea() = %v, want an error: %v", err.Error(), errMessage)
	}
}
