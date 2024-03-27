package main

import (
	"flag"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/1414C/jiffy/gen"
)

func main() {

	projectPath := flag.String("p", "/exp", "project path starting in &GOPATH/src")
	modelFile := flag.String("m", "", "use the model file located at a fully-qualified path")
	modelDirectory := flag.String("md", "", "process all model files in the fully-qualified directory location")
	rsaBits := flag.Uint("rb", 2048, "length of generated RSA keys")

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
		log.Fatal("the -m and -md flag are mutually exclusive. exiting...")
	}

	// check that at least one of -m or -md have been included in the arg list
	if *modelFile == "" && *modelDirectory == "" {
		log.Fatal("either the -m or -md flag must be specified in order to build a project. exiting...")
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

	// walk over the template files to ensure they have been packaged
	// into the binary.

	generatedFiles := make([]string, 0)

	// iterate over the entities to create model and controller files
	for i, ent := range entities {

		ent.AppPath = appPath
		entities[i].AppPath = appPath
		fn, err := ent.CreateModelFile(*projectPath)
		if err != nil {
			log.Fatal(err)
		}
		generatedFiles = append(generatedFiles, fn)

		fn, err = ent.CreateModelExtensionPointsFile(*projectPath)
		if err != nil {
			log.Fatal(err)
		}
		generatedFiles = append(generatedFiles, fn)

		fn, err = ent.CreateControllerFile(*projectPath)
		if err != nil {
			log.Fatal(err)
		}
		generatedFiles = append(generatedFiles, fn)

		fn, err = ent.CreateControllerExtensionPointsFile(*projectPath)
		if err != nil {
			log.Fatal(err)
		}
		generatedFiles = append(generatedFiles, fn)
	}

	// generate static model source files
	s := gen.Static{
		SrcDir:   "/static/models",
		DstDir:   *projectPath + "/models",
		AppPath:  appPath,
		Entities: entities,
	}
	fs, err := s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static model extension-point source files
	s = gen.Static{
		SrcDir:  "/static/models/ext",
		DstDir:  *projectPath + "/models/ext",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static controller source files
	s = gen.Static{
		SrcDir:  "/static/controllers",
		DstDir:  *projectPath + "/controllers",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static controller extension-point source files
	s = gen.Static{
		SrcDir:  "/static/controllers/ext",
		DstDir:  *projectPath + "/controllers/ext",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// iterate over the entities to create their relations
	// via the generation of entity-specific controllers.
	for _, ent := range entities {
		fn, err := ent.CreateControllerRelationsFile(*projectPath, entities)
		if err != nil {
			log.Fatal(err)
		}
		generatedFiles = append(generatedFiles, fn)
	}

	// generate static middleware source files
	s = gen.Static{
		SrcDir:  "/static/middleware",
		DstDir:  *projectPath + "/middleware",
		AppPath: appPath,
		ECDSA:   []string{"256", "384", "521"},
		RSA:     []string{"256", "384", "512"},
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// ensure the group folder exists
	d := *projectPath + "/group"
	_, err = os.Stat(d)
	if err != nil {
		os.Mkdir(d, 0755)
	}

	// generate static group-membership client files
	s = gen.Static{
		SrcDir:  "/static/group/gmcl",
		DstDir:  *projectPath + "/group/gmcl",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static group-membership common files
	s = gen.Static{
		SrcDir:  "/static/group/gmcom",
		DstDir:  *projectPath + "/group/gmcom",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static group-membership server files
	s = gen.Static{
		SrcDir:  "/static/group/gmsrv",
		DstDir:  *projectPath + "/group/gmsrv",
		AppPath: appPath,
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate static main application source files
	s = gen.Static{
		SrcDir:   "/static/appobj",
		DstDir:   *projectPath + "/appobj",
		AppPath:  appPath,
		Entities: entities,
		ECDSA:    []string{"256", "384", "521"},
		RSA:      []string{"256", "384", "512"},
	}
	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate main application source files
	// (main.go, main_test.go,)
	s = gen.Static{
		SrcDir:   "/static",
		DstDir:   *projectPath,
		AppPath:  appPath,
		Entities: entities,
	}

	fs, err = s.GenerateStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fs...)

	// generate the go.mod file
	s = gen.Static{
		SrcDir:  "/static/modules",
		DstDir:  *projectPath,
		AppPath: appPath,
	}

	err = s.GenerateGoMod()
	if err != nil {
		log.Fatal(err)
	}

	// JWT key generation
	keyConf := gen.KeyConfig{
		RSABits: *rsaBits,
		// ECDSACurve: "P384", // deprecated
		ECDSA:     []string{"256", "384", "521"},
		RSA:       []string{"256", "384", "512"},
		TargetDir: *projectPath + "/jwtkeys",
	}

	// need to pass the generated key info back to this level in order to populate the initial config file
	conf := gen.Config{}
	err = keyConf.GenerateJWTKeys(&conf)
	if err != nil {
		log.Fatal(err)
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
		log.Println(err)
	}

	// complete the initial app configuration
	conf.ExternalAddress = "127.0.0.1:3000"
	conf.InternalAddress = "127.0.0.1:4444"
	conf.Env = "def"
	conf.PingCycle = 1        // seconds
	conf.FailureThreshold = 5 // suspect status counts
	conf.Pepper = "secret-pepper-key"
	conf.Database = dbConf
	conf.CertFile = "" // https
	conf.KeyFile = ""  // https
	conf.JWTSignMethod = "ES384"
	conf.JWTLifetime = 120 // minutes
	conf.ECDSA384PrivKeyFile = "jwtkeys/ecdsa/ec384.priv.pem"
	conf.ECDSA384PubKeyFile = "jwtkeys/ecdsa/ec384.pub.pem"

	// default the services to active
	service := gen.ServiceActivation{}
	for _, v := range entities {
		service.ServiceName = v.Header.Name
		service.ServiceActive = true
		conf.ServiceActivations = append(conf.ServiceActivations, service)
	}

	fn, err := conf.GenerateAppConf(*projectPath + "/appobj")
	if err != nil {
		log.Fatal(err)
	}
	generatedFiles = append(generatedFiles, fn)

	// generate a sample .config.json file
	err = conf.GenerateSampleConfig(*projectPath)
	if err != nil {
		log.Fatal(err)
	}

	// generate sample Docker configuration / Dockerfile and docker-entrypoint.sh
	err = conf.GenerateSampleDockerConfig(*projectPath + "/docker-sample")
	if err != nil {
		log.Fatal(err)
	}

	// run gofmt / goimports on all generated .go files
	for _, f := range generatedFiles {
		err = gen.ExecuteGoTools(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	// run go mod tidy against the generated app
	err = os.Chdir(*projectPath)
	if err != nil {
		log.Fatal(err)
	}
	cmd2 := exec.Command("go", "mod", "tidy")
	err = cmd2.Run()
	if err != nil {
		log.Fatal(err)
	}
}
