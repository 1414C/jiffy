package gen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/1414C/sqac/common"
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
	var objMap map[string]json.RawMessage
	var entMap map[string]json.RawMessage
	var idMap map[string]json.RawMessage

	// read the models.json file as []byte
	raw, err := ioutil.ReadFile(mf)
	if err != nil {
		return nil, err
	}

	// unmarshal the raw model json into a slice of map[string]json.RawMessage
	// err = json.Unmarshal(raw, &objmapSlice)
	err = json.Unmarshal(raw, &objMap)
	if err != nil {
		return nil, err
	}

	// fmt.Println("objMap:", objMap) // entities[bytes],relations[bytes]

	// deal with entities
	err = json.Unmarshal(objMap["entities"], &objmapSlice) // typeName:string, properties: {}
	if err != nil {
		return nil, err
	}

	for _, entMap = range objmapSlice {
		var e Entity
		e.Header.Name = strings.Title(cleanString(string(entMap["typeName"])))
		e.Header.Value = cleanString(strings.ToLower(e.Header.Name))

		// was a start value provided for the entity-id?
		idString := string(entMap["id_properties"])
		if idString != "" {
			err = json.Unmarshal([]byte(idString), &idMap)
			if err != nil {
				return nil, err
			}
			start, err := strconv.ParseUint(string(idMap["start"]), 10, 64)
			if err != nil {
				return nil, err
			}
			e.Header.Start = start
		}

		// parse the entity properties and use them to build out the
		// content of the Entity's Info slice.
		fieldString := string(entMap["properties"])
		var fieldRecs []Info
		if fieldString != "" {
			fieldRecs, err = buildEntityColumns(fieldString, DataColumn)
			if err != nil {
				fmt.Println("fieldRec error:", err)
				return nil, err
			}
			for _, fr := range fieldRecs {
				fr.GetSqacTagLine(true)
				fr.GetJSONTagLine()
				e.Fields = append(e.Fields, fr)
			}
		}

		// get the composite index definitions, and then augment
		// the SqacTagLine values where required.
		cmpIndexString := string(entMap["compositeIndexes"])
		if cmpIndexString != "" {
			e.Fields, err = buildCompositeIndexes(cmpIndexString, e.Header.Value, e.Fields)
			if err != nil {
				return nil, err
			}
		}

		// get the relationship definitions and append them to
		// the current entity.  this is step one of two in the
		// relations construction.  all entitites need to be
		// created in order for the relations key/field info
		// to be applied.
		relString := string(entMap["relations"])
		if relString != "" {
			e.Relations, err = buildRelationsBase(e.Header.Name, relString, e.Fields)
			if err != nil {
				return nil, err
			}
		}

		// set the extension-point generation flags on the entity header
		extString := string(entMap["ext_points"])
		if extString != "" {
			err = buildExtFlags(extString, &e.Header)
			if err != nil {
				fmt.Println("ext_points error:", err)
				return nil, err
			}
		}
		entities = append(entities, e)
	}

	// complete Relations construction now that all entities have been populated
	err = completeRelations(entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func buildExtFlags(extString string, eHeader *Info) error {

	extString = cleanString(extString)

	var extMap map[string]bool
	err := json.Unmarshal([]byte(extString), &extMap)
	if err != nil {
		return err
	}

	eHeader.GenControllerExt = extMap["gen_controller"]
	eHeader.GenValidatorExt = extMap["gen_validator"]
	eHeader.GenModelExt = extMap["gen_model"]
	return nil
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
			info.RefEntity, ok = extractString(attrObjMap["ref_entity"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"ref_entity\" field")
			}
			info.RefField, ok = extractString(attrObjMap["ref_field"])
			if !ok {
				return nil, fmt.Errorf("incorrect element-type for entity \"ref_field\" field")
			}
			if len(info.RefEntity) > 0 && len(info.RefField) == 0 {
				return nil, fmt.Errorf("ref_entity and ref_field must both be provided in the model in order to define a foreign-key")
			}
			if len(info.RefEntity) == 0 && len(info.RefField) > 0 {
				return nil, fmt.Errorf("ref_entity and ref_field must both be provided in the model in order to define a foreign-key")
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
func buildCompositeIndexes(cIdxString, tn string, info []Info) ([]Info, error) {

	fmt.Println("cIdxString:", cIdxString)
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

		rawColumns := v["index"]
		if !strings.ContainsAny(rawColumns, ",") {
			return nil, fmt.Errorf("composite index error: missing ','; have %s", rawColumns)
		}

		// split out the column names for the composite index, clean them and build
		// the composite index directive string.
		cIdxDirective := fmt.Sprintf("index:idx_%s", tn)
		indexColumnNames := strings.SplitN(rawColumns, ",", 30)
		for i := range indexColumnNames {
			indexColumnNames[i] = superCleanString(indexColumnNames[i])
			indexColumnNames[i] = common.CamelToSnake(indexColumnNames[i])
			cIdxDirective = cIdxDirective + "_" + indexColumnNames[i]
		}

		fmt.Println("indexColumnNames:", indexColumnNames)
		fmt.Println("cIdxDirective:", cIdxDirective)

		// now finally update the SqacTagLines.  for each column name in the
		// index, read the list of fields in the entity ([]info).  when an
		// index-column-name matches a field-name, add the composite index
		// to that field's .SqacTagLine.
		for _, cn := range indexColumnNames {
			for i, fr := range info {
				if fr.SnakeCaseName == cn {
					tl := fr.SqacTagLine
					switch len(tl) {
					case 0:
						info[i].SqacTagLine = "sqac:\"" + cIdxDirective + "\""

					default:
						fr.SqacTagLine = strings.TrimSuffix(fr.SqacTagLine, "\"")
						info[i].SqacTagLine = fr.SqacTagLine + ";" + cIdxDirective + "\""
					}
				}
			}
		}
	}
	return info, nil
}

// buildRelations reads relations information from the input string,
// converts the string to a slice of maps, reads each map and populates
// the relation struct.
func buildRelationsBase(fromEntName, relString string, info []Info) ([]Relation, error) {

	var relations []Relation
	var relation Relation

	relString = cleanString(relString)

	var relMapSlice = make([]map[string]json.RawMessage, 0)
	err := json.Unmarshal([]byte(relString), &relMapSlice)
	if err != nil {
		fmt.Println("relations unmarshalling error:", err)
		return nil, err
	}

	for _, relMap := range relMapSlice {

		relation.RelName = cleanString(string(relMap["relName"]))
		relation.RelNameLC = strings.ToLower(relation.RelName)
		relPropString := string(relMap["properties"])
		if relString != "" {
			var relPropMap map[string]string
			err := json.Unmarshal([]byte(relPropString), &relPropMap)
			if err != nil {
				return nil, err
			}
			relation.RefKey = relPropMap["refKey"]
			relation.RelType = relPropMap["relType"]
			if relation.RelType != "hasOne" && relation.RelType != "hasMany" && relation.RelType != "belongsTo" {
				return nil, fmt.Errorf("relations relationship-type error - %s is not a valid relationship type", relation.RelType)
			}
			relation.FromEntity = fromEntName
			relation.FromEntityLC = strings.ToLower(relation.FromEntity)
			relation.ToEntity = relPropMap["toEntity"]
			relation.ToEntityLC = strings.ToLower(relation.ToEntity)
			relation.ForeignPK = relPropMap["foreignPK"]
			relations = append(relations, relation)
		}
	}
	return relations, nil
}

// completeRelations uses the completed entity definitions to
// validate and populate relations keys.
func completeRelations(entities []Entity) error {

	// for each entity-relationship , embed the ToEntity's []Info.
	// This is used for validation of user-based key selection when
	// generating the controllers for the relations.
	for i := range entities {
		for j := range entities[i].Relations {
			for _, v := range entities {
				if v.Header.Name == entities[i].Relations[j].ToEntity {
					// prevent non-persisted fields from inclusion in []ToEntInfo
					for _, f := range v.Fields {
						if !f.NoDB {
							entities[i].Relations[j].ToEntInfo = append(entities[i].Relations[j].ToEntInfo, f)
						}
					}
				}
			}
		}
	}

	// examine each relation and determine the validity of the named keys,
	// and/or set the default key names in the Relation to remove the need
	// to perform the operation inside the text/template.
	for i := range entities {
		for j := range entities[i].Relations {
			v := &entities[i].Relations[j]

			// from side
			if v.RefKey == "" {
				switch v.RelType {
				case "hasOne":
					v.RefKey = "ID"
					v.RefKeyOptional = false

				case "hasMany":
					v.RefKey = "ID"
					v.RefKeyOptional = false

				case "belongsTo":
					v.RefKey = v.ToEntity + "ID"

				default:
					// RelType is previously checked - so hard stop here
					panic(fmt.Errorf("unknown relationship type %s detected", v.RelType))
				}
			}
			// validate from RefKey
			bValid := false
			if v.RefKey != "ID" {
				for _, f := range entities[i].Fields {
					if f.Name == v.RefKey {
						bValid = true
						entities[i].Relations[j].RefKeyOptional = !f.Required
					}
				}
				if bValid == false {
					return fmt.Errorf("RefKey %s is not a valid field-name in %s relationship %s", v.RefKey, entities[i].Header.Name, v.RelName)
				}
			}

			// to side
			if v.ForeignPK == "" {
				switch v.RelType {
				case "hasOne":
					v.ForeignPK = v.FromEntity + "ID"

				case "hasMany":
					v.ForeignPK = v.FromEntity + "ID"

				case "belongsTo":
					v.ForeignPK = "ID"
					v.ForeignPKOptional = false

				default:
					// RelType is previously checked - so hard stop here
					panic(fmt.Errorf("unknown relationship type %s detected", v.RelType))
				}
			}
			// validate to ForeignPK
			if v.ForeignPK != "ID" {
				for _, f := range v.ToEntInfo {
					if f.Name == v.ForeignPK {
						bValid = true
						entities[i].Relations[j].ForeignPKOptional = !f.Required
					}
				}
				if bValid == false {
					return fmt.Errorf("ForeignPK %s is not a valid field-name in %s relationship %s", v.ForeignPK, entities[i].Header.Name, v.RelName)
				}
			}
		}
	}
	return nil
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
