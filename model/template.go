package model

import (
	"bytes"
	"fmt"
	"text/template"

	cparser "github.com/cartabinaria/config-parser-go"
)

type globaVar struct {
	ChatGroups map[string]map[string]string
	Domains    map[string]string
}

func FillActionsTemplate(actions []byte) ([]byte, error) {

	if len(Degrees) == 0 {
		panic("Degrees empty, could not continue. This is probably caused by not parsing degrees before actions\n")
	}

	SUB_VAR := getVariables()

	tmpl, err := template.New("tmp").Parse(string(actions))
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, SUB_VAR)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func getVariables() globaVar {

	degreeInformatica, found := Degrees["informatica"]
	if !found {
		panic("Informatica not found in degrees. Could not continue")
	}

	degreeIngegneria, found := Degrees["ingegneria"]
	if !found {
		panic("ingegneria not found in degrees. Could not continue")
	}

	degreeManagement, found := Degrees["informatica_per_il_management"]
	if !found {
		fmt.Println(Degrees)
		panic("informatica_per_il_management not found in degrees. Could not continue")
	}

	degreeIngegneriaMag, found := Degrees["ingegneria_informatica_magistral"] // Without the final 'e' ?
	if !found {
		panic("ingegneria_informatica_magistrale not found in degrees. Could not continue")
	}

	degreeInformaticaMag, found := Degrees["informatica_magistrale"]
	if !found {
		panic("informatica_magistrale not found in degrees. Could not continue")
	}

	degreeInformaticaScienzeInfoMag, found := Degrees["ingegneria_e_scienze_informatich"]
	if !found {
		panic("ingegneria-e-scienze-informatiche-magistrale not found in degrees. Could not continue")
	}

	degreeAI, found := Degrees["artificial_intelligence"]
	if !found {
		panic("artificial_intelligence not found in degrees. Could not continue")
	}

	lab, found := Degrees["lab"]
	if !found {
		panic("lab not found in degrees. Could not continue")
	}

	v := globaVar{
		ChatGroups: map[string]map[string]string{
			"Informatica": {
				"Global": "",
				"First":  cparser.MustGetYear(degreeInformatica, 1).Chat,
				"Second": cparser.MustGetYear(degreeInformatica, 2).Chat,
				"Third":  cparser.MustGetYear(degreeInformatica, 3).Chat,
			},
			"Informatica_per_il_management": {
				"Global": degreeManagement.Chat, // Need to be modified
			},
			"Ingegneria": {
				"Global": degreeIngegneria.Chat,
			},
			"Ingegneria_informatica_magistrale": {
				"Global": degreeIngegneriaMag.Chat,
			},
			"Informatica_magistrale": {
				"Global": degreeInformaticaMag.Chat,
			},
			"Ingegneria_e_scienze_informatiche_magistrale": {
				"Global": degreeInformaticaScienzeInfoMag.Chat,
			},
			"Artificial_intelligence": {
				"First":  cparser.MustGetYear(degreeAI, 1).Chat,
				"Second": cparser.MustGetYear(degreeAI, 2).Chat,
			},
			"Lab": {
				"Global": lab.Chat,
			},
		},
		// The following domains are hard-coded, we could write them in config...
		Domains: map[string]string{
			"ADMStaffBase": "students.cs.unibo.it",
			"GithubBase":   "github.com/cartabinaria",
			"GithubPages":  "cartabinaria.github.io",
		},
	}

	// Add https:// to links for ChatGroups
	for i := range v.ChatGroups {
		for j := range v.ChatGroups[i] {
			v.ChatGroups[i][j] = "https://" + v.ChatGroups[i][j]
		}
	}

	return v
}
