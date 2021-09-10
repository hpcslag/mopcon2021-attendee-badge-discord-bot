package main

import (
	"fmt"

	"github.com/tealeg/xlsx/v3"
)

var KeyPairMap map[string]string = make(map[string]string)

func ReadXLSXToMap(filename string, keypairmap *map[string]string) {
	wb, err := xlsx.OpenFile(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("Reading XLSX File: ", filename)
	sheet := wb.Sheets[0]

	sheet.ForEachRow(func(r *xlsx.Row) error {
		(*keypairmap)[r.GetCell(0).String()] = r.GetCell(1).String()
		return nil
	})
}
