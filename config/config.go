package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

const (
	ProPath = "./"
)

var (
	Config DLZConfig //全局配置文件信息
)

type DLZConfig struct {
	configPath string                 //配置文件路径
	configMap  map[string]interface{} //配置文件中的所有配置信息
	AppPath    string                 //项目app运行路径
	Addr       string                 //运行端口
	Debug      bool                   //是否开启模式运行
	//EnableDoc       bool                   `yaml:"enable_doc"`     //是否api开启文档
	DefaultPassword string `yaml:default_password` //默认密码
	MaxMediaSize    string `yaml:"max_media_size"`
	CookieExpires   int64  `yaml:"cookie_expires"` //cookie过期时间
	//SwaggerPath     string                 //swagger-ui路径
	StaticPath string               //资源文件的路径
	Domain     string `yaml:domain` //域名
	WebApps    string               //静态项目文件放置位置
	Redis      Redis                //redis配置信息
	Swagger    SwaggerConf          //swagger配置信息
	Website    string               //网址
}

func LoadConfig(confPath string) (DLZConfig, error) {
	var conf DLZConfig
	var err error
	conf, err = BuildConfig(confPath)
	if err == nil {
		log.Println("conffile load success")
	}
	log.Println("appPath:", conf.AppPath)
	log.Println("staticPath:", conf.StaticPath)
	log.Println("webapps:", conf.WebApps)
	log.Println("Website:", conf.Website)
	log.Println("Domain:", conf.Domain)
	Config = conf
	return conf, err
}

//构建配置文件信息
func BuildConfig(confPath string) (conf DLZConfig, err error) {
	if confPath == "" {
		err = errors.New("no configPath load")
	} else {
		confPath, _ := filepath.Abs(confPath)
		confFile, _ := ioutil.ReadFile(confPath)
		err := yaml.Unmarshal(confFile, &conf)
		if err != nil {
			return conf, err
		} else {
			confMap := make(map[string]interface{})
			err := yaml.Unmarshal(confFile, &confMap)
			if err != nil {
				return conf, err
			} else {
				conf.configMap = confMap
				conf.configPath = confPath
				if conf.AppPath == "" {
					conf.AppPath = filepath.Dir(confPath)
				}
				if conf.StaticPath == "" {
					appPath, err := filepath.Abs(ProPath)
					if err != nil {
						panic(err)
					} else {
						conf.StaticPath = filepath.Join(appPath, "static")
					}
				} else {
					conf.StaticPath = getFullPath(conf.AppPath, conf.StaticPath)
				}
				if conf.WebApps == "" {
					//conf.WebApps = filepath.Join(conf.AppPath, "../webapps")
					appPath, err := filepath.Abs(ProPath)
					if err != nil {
						panic(err)
					} else {
						conf.WebApps = filepath.Join(appPath, "webapps")
					}
				} else {
					conf.WebApps = getFullPath(conf.AppPath, conf.WebApps)
				}
				if conf.Swagger.Enabledoc {
					if conf.Swagger.Swaggerpath == "" {
						//conf.Swagger.Swaggerpath = getFullPath(conf.AppPath, "../swagger-ui")
						appPath, err := filepath.Abs(ProPath)
						if err != nil {
							panic(err)
						} else {
							conf.StaticPath = filepath.Join(appPath, "swagger-ui")
						}
					} else {
						conf.Swagger.Swaggerpath = getFullPath(conf.AppPath, conf.Swagger.Swaggerpath)
					}
				}
				if conf.Website == "" {
					conf.Website = conf.configMap["domain"].(string)
				}
			}
		}
	}
	return
}

/**
 * 获取配置文件的地址
 */
func (dc *DLZConfig) GetConfigPath() string {
	return dc.configPath
}

/**
 * 获取配置文件中key的值
 */
func (dc *DLZConfig) Get(key string) interface{} {
	if val, ok := dc.configMap[key]; ok {
		return val
	}
	return nil
}

/**
 * 获取配置文件中某个key的值，如果没有，返回默认值
 */
func (dc *DLZConfig) GetString(key string, defaultStr string) string {
	val := dc.Get(key)
	if val == nil || val.(string) == "" {
		return defaultStr
	}
	return val.(string)
}

/**
 * 获取配置文件中的int配置
 */
func (dc *DLZConfig) GetInt(key string, defaultInt int) int {
	val := dc.Get(key)
	if val == nil {
		return defaultInt
	}
	return val.(int)
}

/**
 *获取相对于配置文件的绝对路径
 */
func (dc *DLZConfig) GetAbsPath(appPath string) string {
	return getFullPath(filepath.Dir(dc.configPath), appPath)
}

/**
 * 获取全路径
 */
func getFullPath(appPath string, path string) string {
	if path == "" {
		return appPath
	}
	if filepath.IsAbs(path) {
		return path
	} else {
		return filepath.Join(appPath, path)
	}
}

/**
 * 获取相对工程所在目录的绝对路径
 */
func (dc *DLZConfig) GetAbsProPath(path string) string {
	basePath := dc.StaticPath
	files := strings.Split(path, "/static/")
	absProPath := filepath.Join(basePath, files[1])
	log.Println("absProPath = ", absProPath)
	return absProPath
}

// redis配置
type Redis struct {
	Addr        string
	Password    string
	Database    string
	UniqueIdKey string
	Prefix      string
}

type SwaggerConf struct {
	Enabledoc   bool   //是否开启api文档
	Apipackage  string //api包所在位置,使用逗号分隔
	Mainpackage string //项目主入口文件
	Swaggerpath string //swagger-ui文件所在位置
	Controller  string //controller类所在子包名
}

/**
 * 获取配置文件中某个key的值，如果没有，返回默认值
 */
func (dc *DLZConfig) GetBool(key string) bool {
	val := dc.Get(key)
	if val == nil {
		return false
	}
	return val.(bool)
}
