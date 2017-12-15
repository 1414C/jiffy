package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"strings"

	"github.com/1414C/rgen/gen"
)

func main() {

	projectPath := flag.String("p", "/exp/start_test", "project path starting in &GOPATH/src")
	modelFile := flag.String("m", "./model.scratch2.json", "model file relative to application base directory")

	flag.Parse()
	if *projectPath == "" {
		os.Exit(-1)
	}

	// verify that the project path exists under $GOPATH/src/
	appPath := strings.TrimPrefix(*projectPath, "/")
	*projectPath = build.Default.GOPATH + "/src" + *projectPath
	_, err := os.Stat(*projectPath)
	if err != nil {
		os.Mkdir(*projectPath, 0755)
	}

	// read the JSON models file to get the Entity definitions
	entities, err := gen.ReadModelFile(*modelFile) // gen.GetEntities()
	if err != nil {
		fmt.Println(err)
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

		fn, err = ent.CreateControllerFile(*projectPath)
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
