package gen

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"strings"
// )

// // ReadModelFile2 reads a model file
// func ReadModelFile2(mf string) ([]Entity, error) {

// 	var entities []Entity
// 	var objmapSlice []map[string]json.RawMessage
// 	var objMap map[string]json.RawMessage
// 	var entMap map[string]json.RawMessage

// 	// read the models.json file as []byte
// 	raw, err := ioutil.ReadFile(mf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// unmarshal the raw model json into a slice of map[string]json.RawMessage
// 	// err = json.Unmarshal(raw, &objmapSlice)
// 	err = json.Unmarshal(raw, &objMap)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// fmt.Println("objMap:", objMap) // entities[bytes],relations[bytes]

// 	// deal with entities
// 	err = json.Unmarshal(objMap["entities"], &objmapSlice) // typeName:string, properties: {}
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, entMap = range objmapSlice {
// 		var e Entity
// 		e.Header.Name = strings.Title(cleanString(string(entMap["typeName"])))
// 		e.Header.Value = cleanString(strings.ToLower(e.Header.Name))

// 		// parse the entity properties and use them to build out the
// 		// content of the Entity's Info slice.
// 		fieldString := string(entMap["properties"])
// 		var fieldRecs []Info
// 		if fieldString != "" {
// 			fieldRecs, err = buildEntityColumns(fieldString, DataColumn)
// 			if err != nil {
// 				fmt.Println("fieldRec error:", err)
// 				return nil, err
// 			}
// 			for _, fr := range fieldRecs {
// 				fr.GetRgenTagLine(true)
// 				fr.GetJSONTagLine()
// 				e.Fields = append(e.Fields, fr)
// 			}
// 		}

// 		// get the composite index definitions, and then augment
// 		// the RgenTagLine values where required.
// 		cmpIndexString := string(entMap["compositeIndexes"])
// 		if cmpIndexString != "" {
// 			e.Fields, err = buildCompositeIndexes(cmpIndexString, e.Fields)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}

// 		entities = append(entities, e)
// 	}

// 	// deal with relationships and foreign-keys
// 	var relations []Relation
// 	var relation Relation
// 	var relmapSlice []map[string]json.RawMessage
// 	err = json.Unmarshal(objMap["relations"], &relmapSlice) // typeName:string, properties: {}
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, relMap := range relmapSlice {

// 		relation.RelName = cleanString(string(relMap["relName"]))
// 		relString := string(relMap["properties"])
// 		if relString != "" {
// 			err = buildRelation(relString, &relation)
// 			if err != nil {
// 				return nil, err
// 			}
// 			relations = append(relations, relation)
// 			fmt.Println("relation:", relation)
// 		}
// 	}

// 	// add the relations to their entity
// 	for _, v := range relations {
// 		for i := range entities {
// 			if v.FromEntity != entities[i].Header.Name {
// 				entities[i].Relations = append(entities[i].Relations, v)
// 			}
// 		}
// 	}
// 	return entities, nil
// }
