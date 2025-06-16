// SPDX-FileCopyrightText: 2023 Angelo 'Flecart' Huang <xuanqiang.huang@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Stefano Volpe <foxy@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package utils

import (
	"fmt"
	"testing"
)

func ExampleToSnakeCase() {
	fmt.Println(ToSnakeCase("Hello World"))
	// Output: hello_world
}

func TestToSnakeCase(t *testing.T) {

	tests := []struct {
		name string
		args string
		want string
	}{
		{
			"Standard",
			"Logica per l'informatica",
			"logica_per_l_informatica",
		},
		{
			"Standard",
			"Informatica e Societa",
			"informatica_e_societa",
		},
		{
			"Standard",
			"Sistemi Operativi",
			"sistemi_operativi",
		},
		{
			"Accents",
			"Informatica e Società",
			"informatica_e_societa",
		},
		{
			"Accents",
			"à è ì ò ù",
			"a_e_i_o_u",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.args); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
