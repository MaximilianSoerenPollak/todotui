package db

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/maximiliansoerenpollak/todotui/constants"
	"github.com/maximiliansoerenpollak/todotui/types"
)

var MemData = types.DbDataT{TaskGroups: []types.TaskGroup{}}

func Read() types.DbDataT {
	file, err := os.ReadFile(constants.OutputFile)

	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}

	data := types.DbDataT{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func Write(DbData types.DbDataT) {

	file, err := json.MarshalIndent(DbData, "", " ")
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
	err = os.WriteFile(constants.OutputFile, file, 0644)
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
}
