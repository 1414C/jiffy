package gen

import (
	// "github.com/1414C/rgen/gen"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"
)

// Info is used to hold name-value-pairs for Entity definitions
type Info struct {
	Name  string // field name from model
	Value string // type
	// LowerCaseName       string // field name in lower-case for query path - possibly deprecated
	SnakeCaseName       string // field name in gorm database format (snake_case)
	DBType              string
	IsKey               bool
	Format              string
	Required            bool
	Unique              bool
	Index               string // unique, not-unique, ""
	Selectable          string // "eq,like,gt,lt,ge,le,ne"
	GormTagLine         string
	Relation            string
	RelationField       string
	RelationCardinality string // iota?
}

// Entity definition
type Entity struct {
	Header Info
	Fields []Info
	// CompositeIndexes []string  //
	AppPath string
}

// Static definition
type Static struct {
	SrcDir   string
	DstDir   string
	AppPath  string
	Entities []Entity
}

// constants
const (
	cStrCrt  string  = "crt_string"
	cFl32Crt float32 = 9.99
	cFl64Crt float64 = 1900.99
	cICrt    int     = 100000
	cI8Crt   int8    = 10
	cI16Crt  int16   = 100
	cI32Crt  int32   = 1000
	cI64Crt  int64   = 10000
	cUCrt    uint    = 500000
	cU8Crt   uint8   = 50
	cU16Crt  uint16  = 500
	cU32Crt  uint32  = 5000
	cU64Crt  uint64  = 50000

	cStrUpd  string  = "upd_string"
	cFl32Upd float32 = 8.88
	cFl64Upd float64 = 8888.88
	cIUpd    int     = 888888
	cI8Upd   int8    = 88
	cI16Upd  int16   = 888
	cI32Upd  int32   = 8888
	cI64Upd  int64   = 88888
	cUUpd    uint    = 999999
	cU8Upd   uint8   = 99
	cU16Upd  uint16  = 999
	cU32Upd  uint32  = 9999
	cU64Upd  uint64  = 99999
)

//=============================================================================================
// public Entity generation functions
//=============================================================================================

// CreateModelFile generates a model file for the Entity
// using the user-defined model.json file in conjunction
// with the model.gotmpl text/template.  Returns the fully-qualified
// file-name / error.
func (ent *Entity) CreateModelFile(tDir string) (fName string, err error) {

	// https://medium.com/@IndianGuru/understanding-go-s-template-package-c5307758fab0
	mt := template.New("Entity model template")
	mt, err = template.ParseFiles("templates/model.gotmpl")
	if err != nil {
		log.Fatal("Parse: ", err)
		return "", err
	}

	// create the model file path and create if required
	tDir = tDir + "/models"
	_, err = os.Stat(tDir)
	if err != nil {
		os.Mkdir(tDir, 0755)
	}

	// create the model file
	tfDir := tDir + "/" + ent.Header.Value + "m.go"
	f, err := os.Create(tfDir)
	if err != nil {
		log.Fatal("CreateModelFile: ", err)
		return "", err
	}
	defer f.Close()

	// set permissions
	err = f.Chmod(0644)
	if err != nil {
		log.Fatal("CreateModelFile: ", err)
		return "", err
	}

	// execute the template using the new model file as a target
	err = mt.Execute(f, ent)
	if err != nil {
		log.Fatal("CreateModelFile: ", err)
		return "", err
	}
	log.Println("generated:", tfDir)
	f.Close()
	return tfDir, nil
}

// CreateControllerFile generates a controller file for the Entity
// using the user-defined model.json file in conjunction
// with the controller.gotmpl text/template.  Returns the fully-qualified
// file-name / error.
func (ent *Entity) CreateControllerFile(tDir string) (fName string, err error) {
	ct := template.New("Entity controller template")
	ct, err = template.ParseFiles("templates/controller.gotmpl")
	if err != nil {
		log.Fatal("Parse: ", err)
		return "", err
	}

	// create the controller file path and create if required
	tDir = tDir + "/controllers"
	_, err = os.Stat(tDir)
	if err != nil {
		os.Mkdir(tDir, 0755)
	}

	// create the controller file
	tfDir := tDir + "/" + ent.Header.Value + "c.go"
	f, err := os.Create(tfDir)
	if err != nil {
		log.Fatal("CreateControllerFile: ", err)
		return "", err
	}
	defer f.Close()

	// set permissions
	err = f.Chmod(0644)
	if err != nil {
		log.Fatal("CreateControllerFile: ", err)
		return "", err
	}

	// execute the template using the new controller file as a target
	err = ct.Execute(f, ent)
	if err != nil {
		log.Fatal("CreateControllerFile: ", err)
		return "", err
	}
	log.Println("generated:", tfDir)
	f.Close()
	return tfDir, nil
}

//=============================================================================================
// static generation functions
//=============================================================================================

// GenerateStaticTemplates reads the ./static folder and uses Glob
// to execute each template in-turn.  Returns the fully-qualified
// file-names or an error.
func (s *Static) GenerateStaticTemplates() (fNames []string, err error) {

	// log.Println("s.DstDir:", s.DstDir)
	// log.Println("s.SrcDir:", s.SrcDir)

	tmlFiles, err := filepath.Glob(s.SrcDir + "/*" + ".gotmpl")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, f := range tmlFiles {

		st := template.New("Static template")
		st, err := template.ParseFiles(f)
		if err != nil {
			log.Fatal("Parse: ", err)
			return nil, err
		}

		// create the file-path if required
		_, err = os.Stat(s.DstDir)
		if err != nil {
			os.Mkdir(s.DstDir, 0755)
		}

		// create the static source file
		fileName := filepath.Base(f)
		fileName = strings.TrimSuffix(fileName, "tmpl")
		// log.Println(fileName)
		f, err := os.Create(s.DstDir + "/" + fileName)
		if err != nil {
			log.Fatal("generateStaticTemplates: ", err)
			return nil, err
		}
		defer f.Close()

		// set permissions
		err = f.Chmod(0755)
		if err != nil {
			log.Fatal("generateStaticTemplates: ", err)
			return nil, err
		}

		// execute the template using the new controller file as a target
		err = st.Execute(f, s)
		if err != nil {
			log.Fatal("generateStaticTemplates: ", err)
			return nil, err
		}
		fName := s.DstDir + "/" + fileName
		fNames = append(fNames, fName)
		log.Println("generated:", fName)
		f.Close()
	}
	return fNames, nil
}

//=============================================================================================
// model generation functions
//=============================================================================================

// GetGormTagLine returns a string containing a set of gorm
// directives for the column attributes.
// Called from within readmodel.go/ReadModelFile()
func (i *Info) GetGormTagLine(b bool) string {

	// set `not null`
	if i.Required {
		i.gormTagLineExtend("not null")
	}

	// set dbType if provided
	// for example: type:varchar(100)
	if i.DBType != "" {
		i.gormTagLineExtend("type:" + i.DBType)
	}

	// set unique in column directive
	if i.Unique {
		i.gormTagLineExtend("unique")
	}

	// if an index has been specified, add the relevant index directive
	switch i.Index {
	case "":
	case "unique":
		i.gormTagLineExtend("unique_index")
	case "nonUnique":
		i.gormTagLineExtend("index")
	default:
		// do nothing
	}

	if len(i.GormTagLine) > 0 {
		i.GormTagLine = "gorm:" + i.GormTagLine
		return i.GormTagLine
	}
	return ""
}

// gormTagLineExtend is used to build-out the `gorm:"..."` model directive field
// Called from within a text/template.
func (i *Info) gormTagLineExtend(s string) {
	if len(i.GormTagLine) > 0 {
		i.GormTagLine = strings.TrimSuffix(i.GormTagLine, "\"")
		i.GormTagLine = i.GormTagLine + ";" + s + "\""
		return
	}
	i.GormTagLine = "\"" + s + "\""
}

//=============================================================================================
// template functions
//=============================================================================================

// GetLowerCasePrefixLetter is a method that will be called
// from within the templates to return the first letter of
// the lower-case name of the entity.  used for the model
// -> service vars.  defaults to "e" for entity, but this
// should never occur.
// Called from within a text/template.
func (ent *Entity) GetLowerCasePrefixLetter() string {
	if len(ent.Header.Value) > 0 {
		return string(ent.Header.Value[0])
	}
	return "e"
}

// GetDateTimeStamp returns a stringified date-time in
// RFC822 format for use in template execution.
// Called from within a text/template.
func (ent *Entity) GetDateTimeStamp() string {
	return time.Now().Format(time.RFC822)
}

// GetQueryOps is used to obtain a slice of the required simple query operators
// for the entity controller template.
// Acceptable types are 'EQ','LT','LE','GT','GE','LIKE'.
// Called from within a text/template.
func (i *Info) GetQueryOps() []string {

	// checked in template already, but safety and stuff
	if len(i.Selectable) > 0 {
		os := strings.Split(i.Selectable, ",")
		for x := range os {
			os[x] = superCleanString(os[x])
		}
		return os
	}
	return nil
}

// GetHasEQOp checks to see if an Entity field has been configured
// for simple-selection via the use of the 'EQ' operator.
func (i *Info) GetHasEQOp() bool {
	if strings.Contains(i.Selectable, "EQ") {
		return true
	}
	return false
}

// GetHasNEOp checks to see if an Entity field has been configured
// for simple-selection via the use of the 'NE' operator.
func (i *Info) GetHasNEOp() bool {
	if strings.Contains(i.Selectable, "NE") {
		return true
	}
	return false
}

// GetHasLTOp checks to see if an Entity field has been configured
// for simple-selection via the use of the 'LT' operator.
func (i *Info) GetHasLTOp() bool {
	if strings.Contains(i.Selectable, "LT") {
		return true
	}
	return false
}

// GetHasLEOp checks to see if an Entity field has been configured
// for simple-selection via the use of the 'LE' operator.
func (i *Info) GetHasLEOp() bool {
	if strings.Contains(i.Selectable, "LE") {
		return true
	}
	return false
}

// GetHasGTOp checks to see if an Entity field has been configured
// for simple-selection via the use of the 'GT' operator.
func (i *Info) GetHasGTOp() bool {
	if strings.Contains(i.Selectable, "GT") {
		return true
	}
	return false
}

// GetHasGEOp checks to see if an Entity field has been configured
// for simple-selection via the use of the 'GE' operator.
func (i *Info) GetHasGEOp() bool {
	if strings.Contains(i.Selectable, "GE") {
		return true
	}
	return false
}

// GetHasLIKEOp checks to see if an Entity field has been configured
// for simple-selection via the use of the 'LIKE' operator.
func (i *Info) GetHasLIKEOp() bool {
	if strings.Contains(i.Selectable, "LIKE") {
		return true
	}
	return false
}

// GetSelStringRegex examines the Selectable field and generates a
// GET-type gorilla mux route regex based on the requested (and supported)
// operators.
func (i *Info) GetSelStringRegex() string {

	var qOps string
	ops := i.GetQueryOps()
	for _, op := range ops {
		switch op {
		case "EQ":
			qOps = qOps + "EQ|eq|"
		case "NE":
			qOps = qOps + "NE|ne|"
		case "LIKE":
			qOps = qOps + "LIKE|like|"
		default:
			// ignore all others
		}
	}
	if len(qOps) > 0 {
		qOps = superCleanString(qOps)
		qOps = strings.TrimSuffix(qOps, "|")
		// [(]+(?:EQ|eq|NE|ne|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+
		qOps = "[(]+(?:" + qOps + ")+[ ']+[a-zA-Z0-9_]+[')]+"
		return qOps
	}
	return ""
}

// GetSelBoolRegex examines the Selectable field and generates a
// GET-type gorilla mux route regex based on the requested (and supported)
// operators for the bool type.
func (i *Info) GetSelBoolRegex() string {

	var qOps string
	ops := i.GetQueryOps()
	for _, op := range ops {
		switch op {
		case "EQ":
			qOps = qOps + "EQ|eq|"
		case "NE":
			qOps = qOps + "NE|ne|"
		default:
			// ignore all others
		}
	}
	if len(qOps) > 0 {
		qOps = superCleanString(qOps)
		qOps = strings.TrimSuffix(qOps, "|")
		// [(]+(?:EQ|eq)+[ ']+(?:true|TRUE|false|FALSE)+[')]+
		qOps = "[(]+(?:" + qOps + ")+[ ']+(?:true|TRUE|false|FALSE)+[')]+"
		return qOps
	}
	return ""
}

// GetSelNumberRegex examines the Selectable field and generates a
// GET-type gorilla mux route regex based on the requested (and supported)
// operators for the numeric type (uint, int, float).
func (i *Info) GetSelNumberRegex() string {

	var qOps string
	ops := i.GetQueryOps()
	for _, op := range ops {
		switch op {
		case "EQ":
			qOps = qOps + "EQ|eq|"
		case "NE":
			qOps = qOps + "NE|ne|"
		case "LT":
			qOps = qOps + "LT|lt|"
		case "LE":
			qOps = qOps + "LE|le|"
		case "GT":
			qOps = qOps + "GT|gt|"
		case "GE":
			qOps = qOps + "GE|ge|"
		default:
			// ignore all others
		}
	}
	if len(qOps) > 0 {
		qOps = superCleanString(qOps)
		qOps = strings.TrimSuffix(qOps, "|")
		// [(]+(?:EQ|eq|LT|lt|LE|le|GT|gt|GE|ge)+[ ]+[0-9._]+[)]+

		if i.IsFloatFieldType() {
			qOps = "[(]+(?:" + qOps + ")+[ ]+[0-9._]+[)]+"
			return qOps
		}

		if i.IsUIntFieldType() {
			qOps = "[(]+(?:" + qOps + ")+[ ]+[0-9_]+[)]+"
			return qOps
		}

		if i.IsIntFieldType() {
			qOps = "[(]+(?:" + qOps + ")+[ ]+[0-9_-]+[)]+"
			return qOps
		}
	}
	return ""
}

// GetQueryComponentFuncCall is used to determine the function to call in order to
// separate simple query strings into an operator and predicate value of the
// appropriate type.
// Called from within a text/template.
func (i *Info) GetQueryComponentFuncCall() string {

	switch i.Value {
	case "string":
		return "buildStringQueryComponents(searchValue)"
	case "int":
		return "buildIntQueryComponents(searchValue)"
	case "int8":
		return "buildIntQueryComponents(searchValue)"
	case "int16":
		return "buildIntQueryComponents(searchValue)"
	case "int32":
		return "buildIntQueryComponents(searchValue)"
	case "int64":
		return "buildIntQueryComponents(searchValue)"
	case "uint":
		return "buildUIntQueryComponents(searchValue)"
	case "uint8":
		return "buildUIntQueryComponents(searchValue)"
	case "uint16":
		return "buildUIntQueryComponents(searchValue)"
	case "uint32":
		return "buildUIntQueryComponents(searchValue)"
	case "uint64":
		return "buildUIntQueryComponents(searchValue)"
	case "float32":
		return "buildFloatQueryComponents(searchValue)"
	case "float64":
		return "buildFloatQueryComponents(searchValue)"
	case "bool":
		return "buildBoolQueryComponents(searchValue)"
	default:
		return ""
	}
}

// IsStringFieldType is used in the text/templates to determine
// whether an Info.Value is of type "string", or other.
// Called from within a text/template.
func (i *Info) IsStringFieldType() bool {

	if i.Value == "string" {
		return true
	}
	return false
}

// IsBoolFieldType is used in the text/templates to determine
// whether an Info.Value is of type "bool", or other.
// Called from within a text/template.
func (i *Info) IsBoolFieldType() bool {

	if i.Value == "bool" {
		return true
	}
	return false
}

// IsNumberFieldType is used in the text/templates to determine
// whether an Info.Value is of type uint*, int* or float*.
// Called from within a text/template.
func (i *Info) IsNumberFieldType() bool {

	if strings.Contains(i.Value, "int") ||
		strings.Contains(i.Value, "float") {
		return true
	}
	return false
}

// IsFloatFieldType is used to determine whether an Info record
// has a float-type.
// Called internally via Info.IsNumberFieldType.
func (i *Info) IsFloatFieldType() bool {

	if strings.Contains(i.Value, "float") {
		return true
	}
	return false
}

// IsUIntFieldType is used to determine whether an Info record
// has a uint-type.
// Called internally via Info.IsNumberFieldType.
func (i *Info) IsUIntFieldType() bool {

	if strings.Contains(i.Value, "uint") {
		return true
	}
	return false
}

// IsIntFieldType is used to determine whether an Info record
// has an int-type.
// Called internally via Info.IsNumberFieldType.
func (i *Info) IsIntFieldType() bool {

	if i.IsUIntFieldType() {
		return false
	}
	if strings.Contains(i.Value, "int") {
		return true
	}
	return false
}

// GetDateTimeStamp returns a stringified date-time in
// RFC822 format for use in template execution.
// Called from within a text/template.
func (s *Static) GetDateTimeStamp() string {
	return time.Now().Format(time.RFC822)
}

// GetAddrConcatenatedEntities returns a string of concatenated entity
// addresses in the form of "&Entity1{}, &Entity2{}, &Entity3{}".  this
// is useful for AutoMigrate and DestructiveReset purposes.
// Called from within a text/template.
func (s *Static) GetAddrConcatenatedEntities() string {
	var result string
	for _, e := range s.Entities {
		result = result + "&" + e.Header.Name + "{}, "
	}
	result = strings.TrimSuffix(result, ", ")
	return result
}

//=============================================================================================
// main_test support methods and functions
//=============================================================================================

// BuildTestPostJSON constructs a basic JSON message
// body based on the definition of the Entity passed
// in from the template.  The intent is to create
// a body which can be edited by the developer in
// order to add more meaningful data.
// * string types will be assigned: "string_value"
// * float64 type will be assigned an incrementing
//   float value starting at 9.91
// * int types will be assigned an incrementing int
//   value starting at 10.
// * uint types will be assigned an incrementing uint
//   value starting at 10.
//
// `{"name":"test_product",
// 	"height":55.5,
// 	"cost":66.6,
// 	"supplier":"Ace Hardware",
// 	"weight":88.8,
// 	"length":44.4,
// 	"width":33.3,
// 	"name":"TEST_PRODUCT",
// 	"description":"a nice test product",
// 	"uom":"EA"}`
func (ent *Entity) BuildTestPostJSON(isUpdate bool) string {

	var result string

	// log.Println("ent.Fields:", ent.Fields)

	result = result + "`{"
	for _, f := range ent.Fields {

		switch f.Value {
		case "string":
			result = result + fmt.Sprintf("\"%s\":\"%s\",\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "float32":
			result = result + fmt.Sprintf("\"%s\":%.2f,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "float64":
			result = result + fmt.Sprintf("\"%s\":%.2f,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "int":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "int8":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "int16":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "int32":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "int64":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "uint":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "uint8":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "uint16":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "uint32":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		case "uint64":
			result = result + fmt.Sprintf("\"%s\":%d,\n", strings.ToLower(f.Name), getTestValue(isUpdate, f.Value))
		}
	}
	result = strings.TrimSuffix(result, ",\n")
	result = result + "}`"
	return result
}

// BuildTestValidationExpression is used to build a starter-validation
// statement for each entity's Create / Update tests in main_test.go.
func (ent *Entity) BuildTestValidationExpression(isUpdate bool) string {

	var result string

	for _, fi := range ent.Fields {
		switch fi.Value {
		case "string":
			result = result + fmt.Sprintf("e.%s != \"%s\" ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "float32":
			result = result + fmt.Sprintf("e.%s != %.2f ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "float64":
			result = result + fmt.Sprintf("e.%s != %.2f ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "int":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "int8":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "int16":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "int32":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "int64":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "uint":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "uint8":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "uint16":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "uint32":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		case "uint64":
			result = result + fmt.Sprintf("e.%s != %d ||\n", fi.Name, getTestValue(isUpdate, fi.Value))
		}
	}
	result = strings.TrimSuffix(result, " ||\n")
	return result
}

// get values for test and test validations - support Create and Update
func getTestValue(isUpdate bool, dataType string) interface{} {

	switch dataType {
	case "string":
		if isUpdate {
			return "string_update"
		}
		return "string_value"

	case "float32":
		if isUpdate {
			return cFl32Upd
		}
		return cFl32Crt
	case "float64":
		if isUpdate {
			return cFl64Upd
		}
		return cFl64Crt
	case "int":
		if isUpdate {
			return cIUpd
		}
		return cICrt
	case "int8":
		if isUpdate {
			return cI8Upd
		}
		return cI8Crt
	case "int16":
		if isUpdate {
			return cI16Upd
		}
		return cI16Crt
	case "int32":
		if isUpdate {
			return cI32Upd
		}
		return cI32Crt
	case "int64":
		if isUpdate {
			return cI64Upd
		}
		return cI64Crt
	case "uint":
		if isUpdate {
			return cUUpd
		}
		return cUCrt
	case "uint8":
		if isUpdate {
			return cU8Upd
		}
		return cU8Crt
	case "uint16":
		if isUpdate {
			return cU16Upd
		}
		return cU16Crt
	case "uint32":
		if isUpdate {
			return cU32Upd
		}
		return cU32Crt
	case "uint64":
		if isUpdate {
			return cU64Upd
		}
		return cU64Crt
	default:
		log.Printf("unknown data-type %s in test generation - please add support manually", dataType)
		os.Exit(-1)
	}
	return nil
}

//=============================================================================================
// local functions
//=============================================================================================
func cleanString(s string) string {
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")
	return s
}

func superCleanString(s string) string {
	s = strings.Replace(s, "\"", "", 20)
	s = strings.Replace(s, " ", "", 100)
	s = strings.Replace(s, "\n", "", 100)
	return s
}

// ExecuteGoTools runs gofmt -w and goimports on the specified file
func ExecuteGoTools(fileName string) error {

	// runtime.GOOS
	// android darwin dragonfly freebsd linux nacl netbsd openbsd plan9 solaris windows zos

	// commands to run
	tools := []string{"/bin/goimports", "/bin/gofmt"}
	toolArgs := make(map[string]string)
	toolArgs["/bin/gofmt"] = "-w"
	toolArgs["/bin/goimports"] = "-w"

	for _, t := range tools {
		goftmPath := runtime.GOROOT() + t
		_, err := os.Stat(goftmPath)
		if err != nil {
			log.Printf("could not stat %s to format %s\n", t, fileName)
			return err
		}

		_, err = os.Stat(fileName)
		if err != nil {
			log.Printf("%s attempt could not stat %s\n", t, fileName)
			return err
		}

		log.Printf("executing %s %s %s\n", goftmPath, toolArgs[t], fileName)
		// cmd1 := exec.Command(goftmPath, "-w", fileName)
		var cmd1 *exec.Cmd
		if toolArgs[t] != "" {
			cmd1 = exec.Command(goftmPath, toolArgs[t], fileName)
		} else {
			cmd1 = exec.Command(goftmPath, fileName)
		}
		err = cmd1.Run()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
