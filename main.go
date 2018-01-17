package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/1414C/rgen/gen"
)

func main() {

	projectPath := flag.String("p", "/exp/usrgen92", "project path starting in &GOPATH/src")
	modelFile := flag.String("m", "", "model file relative to application base directory")
	modelDirectory := flag.String("md", "", "process all model files in the specified directory")

	flag.Parse()
	if *projectPath == "" {
		log.Fatal("project path must be provided via the -p flag.  exiting...")
	}

	// verify that the project path exists under $GOPATH/src/
	appPath := strings.TrimPrefix(*projectPath, "/")
	*projectPath = build.Default.GOPATH + "/src" + *projectPath
	_, err := os.Stat(*projectPath)
	if err != nil {
		os.Mkdir(*projectPath, 0755)
	}

	// check that -m and -md have not both been included in the arg list
	if *modelFile != "" && *modelDirectory != "" {
		log.Fatal("the m and md flag are mutually exclusive. exiting...")
	}

	// read the JSON models file to get the Entity definitions if a single
	// model file has been specified via the -m flag
	var entities []gen.Entity
	if *modelFile != "" {
		entities, err = gen.ReadModelFile(*modelFile) // gen.GetEntities()
		if err != nil {
			log.Fatal(err, "exiting...")
		}
	}

	// stat the modelDirectory, then read all of the JSON files if multiple
	// model files have been specified via the -md flag
	if *modelDirectory != "" {
		_, err = os.Stat(*modelDirectory)
		if err != nil {
			log.Fatal(err, "exiting...")
		}

		files, err := ioutil.ReadDir(*modelDirectory)
		if err != nil {
			log.Fatal(err, "exiting...")
		}

		for _, mf := range files {
			if strings.HasSuffix(strings.ToUpper(mf.Name()), "JSON") {
				fEntities, err := gen.ReadModelFile(*modelDirectory + "/" + mf.Name())
				if err != nil {
					log.Fatal(err, "exiting...")
				}

				for _, v := range fEntities {
					entities = append(entities, v)
				}
				fEntities = nil
			}
		}
	}

	// perform a cursory check for duplicate entity names
	mapEntities := make(map[string]bool)
	for _, v := range entities {
		mapEntities[v.Header.Name] = false
	}

	for _, v := range entities {
		b := mapEntities[v.Header.Name]
		if b == false {
			mapEntities[v.Header.Name] = true
		} else {
			log.Fatalf("duplicate entity %s found in model files.  please check the model sources and try again.  exiting...\n", v.Header.Name)
		}
	}

	// for _, e := range entities {
	// 	for _, r := range e.Relations {
	// 		fmt.Println("relation:", r.RelName)
	// 		fmt.Println("r.RelType:", r.RelType)
	// 		fmt.Println("r.RefKey:", r.RefKey)
	// 		fmt.Println("r.ForeignPK:", r.ForeignPK)
	// 		fmt.Println("")
	// 		fmt.Println("")
	// 	}
	// }

	generatedFiles := make([]string, 0)

	// iterate over the entities to create model and controller files
	for i, ent := range entities {

		ent.AppPath = appPath
		entities[i].AppPath = appPath
		fn, err := ent.CreateModelFile(*projectPath)
		if err != nil {
			fmt.Println(err)
		}
		generatedFiles = append(generatedFiles, fn)

		fn, err = ent.CreateModelExtensionPointsFile(*projectPath)
		if err != nil {
			fmt.Println(err)
		}
		generatedFiles = append(generatedFiles, fn)

		fn, err = ent.CreateControllerFile(*projectPath)
		if err != nil {
			fmt.Println(err)
		}
		generatedFiles = append(generatedFiles, fn)

		fn, err = ent.CreateControllerExtensionPointsFile(*projectPath)
		if err != nil {
			fmt.Println(err)
		}
		generatedFiles = append(generatedFiles, fn)
	}

	// generate static model source files
	s := gen.Static{
		SrcDir:   "static/models",
		DstDir:   *projectPath + "/models",
		AppPath:  appPath,
		Entities: entities,
	}
	fs, err := s.GenerateStaticTemplates()
	if err != nil {
		fmt.Println(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static model extension-point source files
	s = gen.Static{
		SrcDir:  "static/models/ext",
		DstDir:  *projectPath + "/models/ext",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		fmt.Println(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static controller source files
	s = gen.Static{
		SrcDir:  "static/controllers",
		DstDir:  *projectPath + "/controllers",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		fmt.Println(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static controller extension-point source files
	s = gen.Static{
		SrcDir:  "static/controllers/ext",
		DstDir:  *projectPath + "/controllers/ext",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		fmt.Println(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// iterate over the entities to create their relations
	// via the generation of entity-specific controllers.
	for _, ent := range entities {
		fn, err := ent.CreateControllerRelationsFile(*projectPath, entities)
		if err != nil {
			fmt.Println(err)
		}
		generatedFiles = append(generatedFiles, fn)
	}

	// generate static middleware source files
	s = gen.Static{
		SrcDir:  "static/middleware",
		DstDir:  *projectPath + "/middleware",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		fmt.Println(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static main application source files
	s = gen.Static{
		SrcDir:   "static/appobj",
		DstDir:   *projectPath + "/appobj",
		AppPath:  appPath,
		Entities: entities,
	}
	fs, err = s.GenerateStaticTemplates()
	// err = s.GenerateAppObjFile()
	if err != nil {
		fmt.Println(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate main application source files
	// (main.go, main_test.go,)
	s = gen.Static{
		SrcDir:   "static",
		DstDir:   *projectPath,
		AppPath:  appPath,
		Entities: entities,
	}
	// err = s.GenerateMain()
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		fmt.Println(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static util package
	s = gen.Static{
		SrcDir:  "templates/util",
		DstDir:  *projectPath + "/util",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		fmt.Println(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// JWT key generation
	keyConf := gen.KeyConfig{
		RSABits:    0,
		ECDSACurve: "P384",
		TargetDir:  *projectPath + "/jwtkeys",
	}

	err = keyConf.GenerateJWTKeys()
	if err != nil {
		fmt.Println(err)
	}

	// test default DB connection
	dbConf := gen.DBConfig{
		DBDialect: "postgres",
		Host:      "localhost",
		Port:      5432,
		Usr:       "godev",
		Password:  "gogogo123",
		Name:      "glrestgen",
	}

	err = dbConf.Validate()
	if err != nil {
		fmt.Println(err)
	}

	// generate the initial app configuration
	conf := gen.Config{
		Port:           3000,
		Env:            "def",
		Pepper:         "secret-pepper-key",
		Database:       dbConf,
		CertFile:       "", // https
		KeyFile:        "", // https
		JWTPrivKeyFile: "jwtkeys/private.pem",
		JWTPubKeyFile:  "jwtkeys/public.pem",
	}

	// default the services to active
	service := gen.ServiceActivation{}
	for _, v := range entities {
		service.ServiceName = v.Header.Name
		service.ServiceActive = true
		conf.ServiceActivations = append(conf.ServiceActivations, service)
	}

	fn, err := conf.GenerateAppConf(*projectPath + "/appobj")
	if err != nil {
		fmt.Println(err)
	}
	generatedFiles = append(generatedFiles, fn)

	// generate a sample .config.json file
	err = conf.GenerateSampleConfig(*projectPath)
	if err != nil {
		fmt.Println(err)
	}

	// run gofmt / goimports on all generated .go files
	for _, f := range generatedFiles {
		err = gen.ExecuteGoTools(f)
		if err != nil {
			fmt.Println(err)
		}
	}
}
