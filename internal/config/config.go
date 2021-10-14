package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ksauth/pkg/crdadaptor"
	"log"
	"sync"
	"time"
)

type RuleConfig struct {
	Available bool   `json:"available"`
	Model     string `json:"model"`
	Policy    string `json:"policy"`
}

type Config struct {
	Rules              map[string]RuleConfig `json:"rules"`
	ExternalClient     bool                  `json:"externalClient"`
	CasbinWorkspace    string                `json:"casbinWorkspace"`
	ExcludedNamespaces []string              `json:"excludedNamespaces"`
	DebugMode          bool                  `json:"debugMode"`

	CertificateFile   string `json:"certificateFile"`
	PrivateKeyFile    string `json:"privateKeyFile"`
	AuditLogFolder    string `json:"auditLogFolder"`
	ItemsNumberPerLog int    `json:"itemsNumberPerLog"`
}

var currentConfig Config
var mutex sync.Mutex

func GetAuditParam() (string, int) {
	mutex.Lock()
	defer mutex.Unlock()
	return currentConfig.AuditLogFolder, currentConfig.ItemsNumberPerLog
}

func GetClientMode() crdadaptor.ClientType {
	mutex.Lock()
	defer mutex.Unlock()
	if currentConfig.ExternalClient {
		return crdadaptor.EXTERNAL_CLIENT
	}
	return crdadaptor.INTERNAL_CLIENT
}

func GetRules() map[string]RuleConfig {
	mutex.Lock()
	defer mutex.Unlock()
	var res map[string]RuleConfig = map[string]RuleConfig{}
	for k, v := range currentConfig.Rules {
		res[k] = v
	}
	return res
}

func GetCasbinWorkSpace() string {
	mutex.Lock()
	defer mutex.Unlock()
	return currentConfig.CasbinWorkspace
}

func GetExcludedNamespaces() []string {
	mutex.Lock()
	defer mutex.Unlock()
	var res []string = make([]string, len(currentConfig.ExcludedNamespaces))
	copy(res, currentConfig.ExcludedNamespaces)
	return res
}

func GetDebugMode() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return currentConfig.DebugMode
}

func GetCrtAndKey() (string, string) {
	mutex.Lock()
	defer mutex.Unlock()
	return currentConfig.CertificateFile, currentConfig.PrivateKeyFile
}

func loadConfig(configPath string) error {
	mutex.Lock()
	fileContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("Failed to load " + configPath + " due to" + err.Error())
		return fmt.Errorf("Failed to load %s due to %s", configPath, err.Error())
	}
	err = json.Unmarshal(fileContent, &currentConfig)
	if err != nil {
		log.Fatal("Failed to load " + configPath + " due to" + err.Error())
		return fmt.Errorf("Failed to load %s due to %s", configPath, err.Error())
	}
	mutex.Unlock()
	return nil
}

func InitConfig(configPath string) error {
	err := loadConfig(configPath)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", currentConfig)

	//start hotupdate goroutine
	go func() {
		time.Sleep(1 * time.Second)
		loadConfig(configPath)
	}()

	return nil
}
