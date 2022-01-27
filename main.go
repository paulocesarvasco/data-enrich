package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"meli/models"
	"os/user"
	"path/filepath"
)

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	path := filepath.Join(dir, "input.json")
	file, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	var records models.CloudtrailData

	err = json.Unmarshal(file, &records)
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug

	}
	fmt.Printf("%+v\n", records) // output for debug

}
