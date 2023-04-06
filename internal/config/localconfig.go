package config

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
	"sigs.k8s.io/yaml"
)

// LocalConfig is a local Argo CD config file
type LocalConfig struct {
	Server Server `json:"server"`
	User   User   `json:"user"`
}

type Server struct {
	Server string `json:"server"`
}

// User contains user authentication information
type User struct {
	OrganizationId string `json:"organizationId"`
	Name           string `json:"name"`
	AuthToken      string `json:"auth-token,omitempty"`
	RefreshToken   string `json:"refresh-token,omitempty"`
}

// Claims returns the standard claims from the JWT claims
func (u *User) Claims() (*jwt.RegisteredClaims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	claims := jwt.RegisteredClaims{}
	_, _, err := parser.ParseUnverified(u.AuthToken, &claims)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}

// ReadLocalConfig loads up the local configuration file. Returns nil if config does not exist
func ReadLocalConfig(path string) (*LocalConfig, error) {
	var err error
	var config LocalConfig

	err = UnmarshalLocalFile(path, &config)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return &config, nil
}

// WriteLocalConfig writes a new local configuration file.
func WriteLocalConfig(config LocalConfig, configPath string) error {
	err := os.MkdirAll(path.Dir(configPath), os.ModePerm)
	if err != nil {
		return err
	}
	return MarshalLocalYAMLFile(configPath, config)
}

func DeleteLocalConfig(configPath string) error {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return err
	}
	return os.Remove(configPath)
}

func (l *LocalConfig) GetServer() Server {
	return l.Server
}

func (l *LocalConfig) UpsertServer(server Server) {
	l.Server = server
}

func (l *LocalConfig) GetUser() User {
	return l.User
}

func (l *LocalConfig) UpsertUser(user User) {
	l.User = user
}

// DefaultConfigDir returns the local configuration path for settings such as cached authentication tokens.
func DefaultConfigDir() (string, error) {
	// Manually defined config directory
	configDir := os.Getenv("TKS_CONFIG_DIR")
	if configDir != "" {
		return configDir, nil
	}

	homeDir, err := getHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(homeDir, ".config", "tks"), nil
}

func getHomeDir() (string, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}

		homeDir = usr.HomeDir
	}

	return homeDir, nil
}

// DefaultLocalConfigPath returns the local configuration path for settings such as cached authentication tokens.
func DefaultLocalConfigPath() (string, error) {
	dir, err := DefaultConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, "config"), nil
}

// Get username from subject in a claim
func GetUsername(subject string) string {
	parts := strings.Split(subject, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return subject
}

// UnmarshalReader is used to read manifests from stdin
func UnmarshalReader(reader io.Reader, obj interface{}) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return unmarshalObject(data, obj)
}

// unmarshalObject tries to convert a YAML or JSON byte array into the provided type.
func unmarshalObject(data []byte, obj interface{}) error {
	// first, try unmarshaling as JSON
	// Based on technique from Kubectl, which supports both YAML and JSON:
	//   https://mlafeldt.github.io/blog/teaching-go-programs-to-love-json-and-yaml/
	//   http://ghodss.com/2014/the-right-way-to-handle-yaml-in-golang/
	// Short version: JSON unmarshaling won't zero out null fields; YAML unmarshaling will.
	// This may have unintended effects or hard-to-catch issues when populating our application object.
	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &obj)
	if err != nil {
		return err
	}

	return err
}

// MarshalLocalYAMLFile writes JSON or YAML to a file on disk.
// The caller is responsible for checking error return values.
func MarshalLocalYAMLFile(path string, obj interface{}) error {
	yamlData, err := yaml.Marshal(obj)
	if err == nil {
		err = os.WriteFile(path, yamlData, 0600)
	}
	return err
}

// UnmarshalLocalFile retrieves JSON or YAML from a file on disk.
// The caller is responsible for checking error return values.
func UnmarshalLocalFile(path string, obj interface{}) error {
	data, err := os.ReadFile(path)
	if err == nil {
		err = unmarshalObject(data, obj)
	}
	return err
}

func Unmarshal(data []byte, obj interface{}) error {
	return unmarshalObject(data, obj)
}

// UnmarshalRemoteFile retrieves JSON or YAML through a GET request.
// The caller is responsible for checking error return values.
func UnmarshalRemoteFile(url string, obj interface{}) error {
	data, err := ReadRemoteFile(url)
	if err == nil {
		err = unmarshalObject(data, obj)
	}
	return err
}

// ReadRemoteFile issues a GET request to retrieve the contents of the specified URL as a byte array.
// The caller is responsible for checking error return values.
func ReadRemoteFile(url string) ([]byte, error) {
	var data []byte
	resp, err := http.Get(url)
	if err == nil {
		defer func() {
			_ = resp.Body.Close()
		}()
		data, err = io.ReadAll(resp.Body)
	}
	return data, err
}
