package gen

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	_ "github.com/SAP/go-hdb/driver"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// DBConfig type holds db config info
type DBConfig struct {
	DBDialect string `json:"db_dialect"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Usr       string `json:"Usr"`
	Password  string `json:"Password"`
	Name      string `json:"name"`
}

// ServiceActivation struct
type ServiceActivation struct {
	ServiceName   string `json:"service_name"`
	ServiceActive bool   `json:"service_active"`
}

// Config type holds the generated application's configuration info
type Config struct {
	ExternalAddress     string   `json:"external_address"`
	InternalAddress     string   `json:"internal_address"`
	Env                 string   `json:"env"`
	PingCycle           uint     `json:"ping_cycle"`
	FailureThreshold    uint64   `json:"failure_threshold"`
	Pepper              string   `json:"pepper"`
	HMACKey             string   `json:"hmac_key"`
	Database            DBConfig `json:"database"`
	CertFile            string   `json:"cert_file"`
	KeyFile             string   `json:"key_file"`
	RSA256PrivKeyFile   string   `json:"rsa256_priv_key_file"`
	RSA256PubKeyFile    string   `json:"rsa256_pub_key_file"`
	RSA384PrivKeyFile   string   `json:"rsa384_priv_key_file"`
	RSA384PubKeyFile    string   `json:"rsa384_pub_key_file"`
	RSA512PrivKeyFile   string   `json:"rsa512_priv_key_file"`
	RSA512PubKeyFile    string   `json:"rsa512_pub_key_file"`
	ECDSA256PrivKeyFile string   `json:"ecdsa256_priv_key_file"`
	ECDSA256PubKeyFile  string   `json:"ecdsa256_pub_key_file"`
	ECDSA384PrivKeyFile string   `json:"ecdsa384_priv_key_file"`
	ECDSA384PubKeyFile  string   `json:"ecdsa384_pub_key_file"`
	ECDSA521PrivKeyFile string   `json:"ecdsa521_priv_key_file"`
	ECDSA521PubKeyFile  string   `json:"ecdsa521_pub_key_file"`
	JWTSignMethod       string   `json:"jwt_sign_method"` // {EC256|EC384|EC521|RS256|RS384|RS312}
	JWTLifetime         uint     `json:"jwt_lifetime"`    // minutes
	// JWTPrivKeyFile     string              `json:"jwt_priv_key_file"`
	// JWTPubKeyFile      string              `json:"jwt_pub_key_file"`
	ServiceActivations []ServiceActivation `json:"service_activations"`
}

// ConnectionInfo returns a DBConfig string
func (c DBConfig) ConnectionInfo() string {

	switch c.DBDialect {
	case "postgres":
		if c.Password == "" {
			return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Usr, c.Name)
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Usr, c.Password, c.Name)

	case "mssql":
		// "sqlserver://SA:Bunny123!!@localhost:1401?database=sqlx")
		return fmt.Sprintf("sqlserver://%s:%s@l%s:%d?database=%s", c.Usr, c.Password, c.Host, c.Port, c.Name)

	case "hdb":
		// "hdb://SMACLEOD:Blockhead1@clkhana01.lab.clockwork.ca:30047")
		return fmt.Sprintf("hdb://%s:%s@%s:%d", c.Usr, c.Password, c.Host, c.Port)

	case "sqlite":
		// "sqlite3", "testdb.sqlite"
		return fmt.Sprintf("%s", c.Name)

	case "mysql":
		// "stevem:gogogo123@tcp(192.168.1.50:3306)/sqlx?charset=utf8&parseTime=True&loc=Local")
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", c.Usr, c.Password, c.Host, c.Port, c.Name)

	default:
		panic(fmt.Errorf("dialect %s is not recognized", c.DBDialect))

	}
}

// Validate the default postgres configuration
func (c *DBConfig) Validate() error {

	connString := c.ConnectionInfo()

	// if dbc.Password == "" {
	// 	connString = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", dbc.Host, dbc.Port, dbc.User, dbc.Name)
	// } else {
	// 	connString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbc.Host, dbc.Port, dbc.User, dbc.Password, dbc.Name)
	// }

	log.Printf("default %s validation: opening connection\n", c.DBDialect)
	dbHandle, err := sql.Open(c.DBDialect, connString)
	if err != nil {
		return err
	}
	log.Printf("default %s validation: connected to %s\n", c.DBDialect, c.DBDialect)

	defer dbHandle.Close()

	log.Printf("default %s validation: getting db transaction handle\n", c.DBDialect)
	tx, err := dbHandle.Begin()
	if err != nil {
		return err
	}
	log.Printf("default %s validation: got db transaction handle\n", c.DBDialect)

	var pid1 int
	log.Printf("default %s validation: reading db PID\n", c.DBDialect)
	err = tx.QueryRow("SELECT pg_backend_pid()").Scan(&pid1)
	if err != nil {
		return err
	}
	log.Printf("default %s validation: got db PID %v\n", c.DBDialect, pid1)
	return nil
}

// GetIsProd returns a bool value indicating to the calling template that
// the cfg struct is holding production configuration.
func (cfg Config) GetIsProd() bool {
	if strings.ToLower(cfg.Env) == "prod" {
		return true
	}
	return false
}

// IsLastServiceActivationRec is used to determine whether the config.json.gotmpl
// has processed the last ServiceActivation while building the .dec / .prd config
// files.
func (cfg Config) IsLastServiceActivationRec(name string) bool {

	l := len(cfg.ServiceActivations) - 1
	sa := cfg.ServiceActivations[l]
	if sa.ServiceName == name {
		return true
	}
	return false
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

	// execute the template and create the appconf.go
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
			cfg.Env = "dev"
		case 1:
			tfDir = dstDir + "/.prd.config.json"
			cfg.Env = "prod"
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
		// fmt.Println("cfg:", cfg)
		err = at.Execute(f, cfg)
		if err != nil {
			log.Fatal("GenerateSampleConfig: ", err)
			return err
		}
		log.Println("generated:", tfDir)
	}
	return nil
}
