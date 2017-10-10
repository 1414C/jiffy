package gen

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"text/template"
)

// PostgresConfig type holds pg config info
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Name     string `json:"name"`
}

// Config type holds the generated application's configuration info
type Config struct {
	Port           int            `json:"port"`
	Env            string         `json:"env"`
	Pepper         string         `json:"pepper"`
	HMACKey        string         `json:"hmac_key"`
	Database       PostgresConfig `json:"database"`
	CertFile       string         `json:"cert_file"`
	KeyFile        string         `json:"key_file"`
	JWTPrivKeyFile string         `json:"jwt_priv_key_file"`
	JWTPubKeyFile  string         `json:"jwt_pub_key_file"`
}

// Validate the default postgres configuration
func (pgc *PostgresConfig) Validate() error {

	var connString string

	if pgc.Password == "" {
		connString = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", pgc.Host, pgc.Port, pgc.User, pgc.Name)
	} else {
		connString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pgc.Host, pgc.Port, pgc.User, pgc.Password, pgc.Name)
	}

	log.Println("default postgres validation: opening connection")
	dbHandle, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	log.Println("default postgres validation: connected to postgres")

	defer dbHandle.Close()

	log.Println("default postgres validation: getting db transaction handle")
	tx, err := dbHandle.Begin()
	if err != nil {
		return err
	}
	log.Println("default postgres validation: got db transaction handle")

	var pid1 int
	log.Println("default postgres validation: reading db PID")
	err = tx.QueryRow("SELECT pg_backend_pid()").Scan(&pid1)
	if err != nil {
		return err
	}
	log.Println("default postgres validation: got db PID:", pid1)
	return nil
}

// GenerateAppConf generates the default application configuration
// source file appconf.go.
func (cfg *Config) GenerateAppConf(dstDir string) (fName string, err error) {

	at := template.New("Application configuration template")
	at, err = template.ParseFiles("templates/appconf.gotmpl")
	if err != nil {
		log.Fatal("Parse: ", err)
		return "", err
	}

	// check the destination file-path and create if required
	_, err = os.Stat(dstDir)
	if err != nil {
		os.Mkdir(dstDir, 0755)
	}

	// create the appconf.go file
	tfDir := dstDir + "/appconf.go"
	f, err := os.Create(tfDir)
	if err != nil {
		log.Fatal("GenerateAppConf: ", err)
		return "", err
	}
	defer f.Close()

	// set permissions
	err = f.Chmod(0755)
	if err != nil {
		log.Fatal("GenerateAppConf: ", err)
		return "", err
	}

	// execute the appconf.gotmpl template using new file appconf.go as a target
	err = at.Execute(f, cfg)
	if err != nil {
		log.Fatal("GenerateAppObjFile: ", err)
		return "", err
	}
	f.Close()
	log.Println("generated:", tfDir)
	return tfDir, nil
}

// GenerateSampleConfig creates a sample .config.json file
// to hold the production application config.
func (cfg *Config) GenerateSampleConfig(dstDir string) error {

	at := template.New("sample json configuration template")
	at, err := template.ParseFiles("templates/config.json.gotmpl")
	if err != nil {
		log.Fatal("Parse: ", err)
		return err
	}

	// check the destination file-path and create if required
	_, err = os.Stat(dstDir)
	if err != nil {
		os.Mkdir(dstDir, 0755)
	}

	var tfDir string

	for i := 0; i < 2; i++ {

		switch i {
		case 0:
			tfDir = dstDir + "/.dev.config.json"
		case 1:
			tfDir = dstDir + "/.prd.config.json"
		default:

		}

		// create the .xxx.config.json file
		f, err := os.Create(tfDir)
		if err != nil {
			log.Fatal("GenerateSampleConfig: ", err)
			return err
		}
		defer f.Close()

		// set permissions
		err = f.Chmod(0755)
		if err != nil {
			log.Fatal("GenerateSampleConfig: ", err)
			return err
		}

		// execute the config.json.gotmpl template using new file .xxx.config.json as a target
		err = at.Execute(f, nil)
		if err != nil {
			log.Fatal("GenerateSampleConfig: ", err)
			return err
		}
		log.Println("generated:", tfDir)
	}
	return nil
}
