package gen

import (
	"encoding/json"
	"fmt"
	"github.com/1414C/sqac/common"
	"io/ioutil"
	"reflect"
	"strings"
)

// ColType specifies which type of db item is being interpreted
type ColType int

const (
	KeyColumn ColType = iota // 0
	DataColumn
)

// ReadModelFile reads a model file
func ReadModelFile(mf string) ([]Entity, error) {

	var entities []Entity
	var objmapSlice []map[string]json.RawMessage

	// read the models.json file as []byte
	raw, err := ioutil.ReadFile(mf)
	if err != nil {
		return nil, err
	}

	// unmarshal the raw model json into a slice of map[string]json.RawMessage
	err = json.Unmarshal(raw, &objmapSlice)
	if err != nil {
		return nil, err
	}

	var objmap map[string]json.RawMessage

	// iterate over the raw Entities - not good - does not order the struct fields reliably!
	for _, objmap = range objmapSlice {

		// get the Entity header information using the known json names
		var e Entity
		e.Header.Name = strings.Title(cleanString(string(objmap["typeName"])))
		e.Header.Value = cleanString(strings.ToLower(e.Header.Name))

		// parse the entity properties and use them to build out the
		// content of the Entity's Info slice.
		fieldString := string(objmap["properties"])
		var fieldRecs []Info
		if fieldString != "" {
			fieldRecs, err = buildEntityColumns(fieldString, DataColumn)
			if err != nil {
				fmt.Println("fieldRec error:", err)
				return nil, err
			}
			for _, fr := range fieldRecs {
				fr.GetRgenTagLine(true)
				fr.GetJSONTagLine()
				// fmt.Println("fr.GetGormTagLine:", fr.GormTagLine)
				e.Fields = append(e.Fields, fr)
			}
		}

		// get the composite index definitions, and then augment
		// the GormTagLine values where required.
		cmpIndexString := string(objmap["compositeIndexes"])
		if cmpIndexString != "" {
			e.Fields, err = buildCompositeIndexes(cmpIndexString, e.Fields)
			if err != nil {
				return nil, err
			}
		}

		entities = append(entities, e)
	}
	return entities, nil
}

// buildEntityColumns sets the entity field attributes in the Entity.[]Info
func buildEntityColumns(colString string, colType ColType) ([]Info, error) {

	var colInfo []Info
	var info Info
	var ok bool

	colString = cleanString(colString)

	// build the top part
	colObjects := strings.Split(colString, "},") // ["key1_name": {"type":"string"},"key2_name": {"type":"uint"}]

	// create a map of key_name: {key_attr1: value, key_attr2: value}
	var colObjMap map[string]json.RawMessage
	err := json.Unmarshal([]byte(colString), &colObjMap)
	if err != nil {
		return nil, err
	}

	for _, co := range colObjects {

		// isolate the current column name
		keyPair := strings.SplitAfter(co, ":")
		keyPair[0] = strings.TrimPrefix(keyPair[0], "{")
		keyPair[0] = strings.TrimSuffix(keyPair[0], ":")
		info.Name = superCleanString(keyPair[0])
		if strings.ContainsAny(info.Name, "_") {
			return nil, fmt.Errorf("error: snake_case (%s) is not permitted in model field-names;  please use camelCase", info.Name)
		}

		// get the actual column information
		// create a map of [key_attr1: value, key_attr2: value] for each key
		// var attrObjMap map[string]json.RawMessage
		var attrObjMap map[string]interface{}
		v := colObjMap[info.Name]
		err = json.Unmarshal(v, &attrObjMap)
		if err != nil {
			return nil, err
		}

		// decide what to populate
		switch colType {
		// KeyColumn is deprecated for now
		case KeyColumn:
			info.IsKey = true
			info.Required = true
			info.Name = strings.Title(info.Name)
			info.SnakeCaseName = common.CamelToSnake(info.Name)

		case DataColumn:
			info.IsKey = false
			info.Name = strings.Title(info.Name)
			// info.SnakeCaseName = util.ToDBName(info.Name)
			info.SnakeCaseName = common.CamelToSnake(info.Name)
			info.Value, ok = extractString(attrObjMap["type"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"type\" field")
			}
			info.Format, ok = extractString(attrObjMap["format"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"format\" field")
			}
			info.DBType, ok = extractString(attrObjMap["db_type"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"dbType\" field")
			}
			info.Required, ok = extractBool(attrObjMap["required"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"required\" field")
			}
			info.NoDB, ok = extractBool(attrObjMap["no_db"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"no_db\" field")
			}
			info.Unique, ok = extractBool(attrObjMap["unique"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"unique\" field")
			}
			// info.Selectable, ok = extractBool(attrObjMap["selectable"])
			info.Selectable, ok = extractString(attrObjMap["selectable"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"selectable\" field")
			}
			info.Selectable = strings.ToUpper(info.Selectable)

			info.Index, ok = extractString(attrObjMap["index"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"index\" field")
			}
			info.Relation, ok = extractString(attrObjMap["relation"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"relation\" field")
			}
			info.RelationField, ok = extractString(attrObjMap["relationFld"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"relationFld\" field")
			}
			info.RelationCardinality, ok = extractString(attrObjMap["relationCrd"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"relationCrd\" field")
			}

		default:
			// do nothing
		}
		colInfo = append(colInfo, info)
	}
	return colInfo, nil
}

// buildCompositeIndexes adds the composite index sqac directives to
// the Info.SqacTagLine.
func buildCompositeIndexes(cIdxString string, info []Info) ([]Info, error) {

	var indexMap = make([]map[string]string, 0)
	err := json.Unmarshal([]byte(cIdxString), &indexMap)
	if err != nil {
		fmt.Println("composite index error:", err)
		return nil, err
	}

	// indexMap == [map[index:f1, f2] map[index:f3, f4, f5]]
	// iterate over the composite index definitions in no
	// particular order.
	for _, v := range indexMap {
		// fmt.Println("v:", v["index"])
		rawColumns := v["index"]
		if !strings.ContainsAny(rawColumns, ",") {
			return nil, fmt.Errorf("composite index error: missing ','; have %s", rawColumns)
		}

		// split out the column names for the composite index, clean them and build
		// the composite index directive string.
		cIdxDirective := "index:idx"
		indexColumnNames := strings.SplitN(rawColumns, ",", 30)
		for i := range indexColumnNames {
			indexColumnNames[i] = superCleanString(indexColumnNames[i])
			indexColumnNames[i] = common.CamelToSnake(indexColumnNames[i])
			cIdxDirective = cIdxDirective + "_" + indexColumnNames[i]
		}
		// fmt.Println("indexColumnNames:", indexColumnNames)
		// fmt.Println("cIdxDirective:", cIdxDirective)

		// now finally update the RgenTagLines.  for each column name in the
		// index, read the list of fields in the entity ([]info).  when an
		// index-column-name matches a field-name, add the composite index
		// to that field's .GormTagLine.
		for _, cn := range indexColumnNames {
			for i, fr := range info {
				if fr.SnakeCaseName == cn {
					tl := fr.RgenTagLine
					switch len(tl) {
					case 0:
						info[i].RgenTagLine = "rgen:\"" + cIdxDirective + "\""

					default:
						fr.RgenTagLine = strings.TrimSuffix(fr.RgenTagLine, "\"")
						info[i].RgenTagLine = fr.RgenTagLine + ";" + cIdxDirective + "\""
					}
				}
			}
		}
	}
	return info, nil
}

// extractString attempts to read the interface parameter
// as a string-type. if the type-assertion fails and parameter
// i has a non-nil type, throw an error.  if the type-asserstion
// fails and parameter i has a nil-value, set the default value
// for a string and indicate success.
func extractString(i interface{}) (string, bool) {

	var ok bool
	var result string

	result, ok = i.(string)
	if ok {
		result = cleanString(result)
		return result, true
	}

	// check to see what was acually passed; if nil, set the
	// initial value for the type ("") and indicate success
	ti := reflect.TypeOf(i)
	if ti == nil {
		return "", true
	}
	return "", false
}

// extractBool attempts to read the interface parameter
// as a bool-type. if the type-assertion fails and parameter
// i has a non-nil type, throw an error.  if the type-asserstion
// fails and parameter i has a nil-value, set the default value
// for a bool and indicate success.
func extractBool(i interface{}) (bool, bool) {

	var ok bool
	var result bool

	result, ok = i.(bool)
	if ok {
		return result, true
	}

	// check to see what was acually passed; if nil, set the
	// initial value for the type (false) and indicate success
	ti := reflect.TypeOf(i)
	if ti == nil {
		return false, true
	}
	return false, false
}
