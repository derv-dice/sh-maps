package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sh-maps/utils"
	"time"

	"github.com/gobuffalo/packr/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	CfgName         = "config.json"
	staticPath      = "../static"
	templatesPath   = "../templates"
	MapsStaticPath  = `map-static`
	BuildingCfgPath = `b-cfg`
	paramName       = "name"
	mapJsonKeyTmpl  = `"map":"%s.svg"`
)

var (
	Cfg                   = &Config{}
	Static                *packr.Box
	Templates             *packr.Box
	BuildingConfigs       = BuildingConfigMap{}
	UniversityDescription = &UniDescription{}
	rxKeyMap              = regexp.MustCompile(`"map"\s*:\s*"(?P<name>[A-Za-z0-9_-]+)\.svg"`)
)

type Config struct {
	Server struct {
		RemoteAddr string `yaml:"remote_addr"`
		Addr       string `yaml:"addr"`
		Https      struct {
			Enabled bool   `yaml:"enabled"`
			Cert    string `yaml:"cert"`
			Key     string `yaml:"key"`
		} `yaml:"https"`
	} `yaml:"server"`
	Log struct {
		NoLog      bool   `yaml:"no_log"`
		TraceLevel string `yaml:"trace_level"`
		LogDir     string `yaml:"log_dir"`
	} `yaml:"log"`
	MapsDir string `yaml:"maps_dir"`
}

type BuildingConfigMap map[string]*BuildingConfig

type UniDescription struct {
	Abbreviation string `json:"abbreviation"` // Аббревиатура вуза
	FullName     string `json:"full_name"`    // Полное название вуза
	Icon         string `json:"icon"`         // Ссылка на герб вуза
}

type MapsCfg struct {
	Description *UniDescription   `json:"description"`
	Buildings   []*BuildingConfig `json:"buildings"`
}

type BuildingConfig struct {
	Name           string `json:"name"`
	Dir            string `json:"dir"`
	Address        string `json:"address"`
	BuildingConfig []byte `json:"-"`
}

func (c *Config) Load() (err error) {
	log.Println("Загрузка статики")
	Static = packr.New("static", staticPath)

	log.Println("Загрузка шаблонов")
	Templates = packr.New("templates", templatesPath)

	log.Println("Загрузка конфига")
	var data []byte
	if data, err = ioutil.ReadFile("config.yml"); err != nil {
		return
	}
	return yaml.Unmarshal(data, &c)
}

// LogToFile - Включить логирование в файл
func LogToFile(dir string) (err error) {
	// Создаем директорию логов, если ее нет
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			return
		}
	}

	// В качестве имени файла ставим текущее время
	name := time.Now().Format("02-01-2006_15-04-05") + ".log"
	var file io.Writer
	if file, err = os.Create(path.Join("logs", name)); err != nil {
		return
	}

	// Выключаем вывод стандартного логера в созданный файл и терминал
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	var lvl log.Level
	if lvl, err = log.ParseLevel(Cfg.Log.TraceLevel); err != nil {
		log.SetLevel(lvl)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return
}

// PrepareBuildingsConfig - подготовка конфигов карт
func PrepareBuildingsConfig() (err error) {
	if Cfg.MapsDir == "" {
		return fmt.Errorf("не указана директория с файлами карт (svg, config.json)")
	}

	if _, err = os.Stat(Cfg.MapsDir); os.IsNotExist(err) {
		return err
	}

	// Парсинг главного конфига карт в корне директории
	var data []byte
	if data, err = os.ReadFile(filepath.Join(Cfg.MapsDir, CfgName)); err != nil {
		return
	}
	mc := MapsCfg{Buildings: []*BuildingConfig{}}
	if err = json.Unmarshal(data, &mc); err != nil {
		return
	}

	UniversityDescription = mc.Description

	// Подготовка конфигов, описывающих конкретное здание/корпус универа
	for i := range mc.Buildings {
		bDir := filepath.Join(Cfg.MapsDir, mc.Buildings[i].Dir)
		var bCfgStr string
		if bCfgStr, err = prepareSingleBuildingCfg(
			mc.Buildings[i].Dir,
			filepath.Join(bDir, CfgName),
		); err != nil {
			return err
		}

		mc.Buildings[i].BuildingConfig = []byte(bCfgStr)
		BuildingConfigs[mc.Buildings[i].Dir] = mc.Buildings[i]
	}
	return
}

// prepareSingleBuildingCfg - Добавление в конфиги url сервера перед названиями файлов .svg
func prepareSingleBuildingCfg(dir, filename string) (prepared string, err error) {
	// Чтение конфига построчно и замена имен файлов на ссылки 4a.svg --> http://.../4a.svg
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()

		if rxKeyMap.MatchString(line) {
			pm := utils.GetRxParams(rxKeyMap, line)
			line = rxKeyMap.ReplaceAllString(line, fmt.Sprintf(mapJsonKeyTmpl, fmt.Sprintf("%s/%s",Cfg.Server.RemoteAddr, path.Join(MapsStaticPath, dir, pm[paramName]))))
		}
		prepared += line
	}

	err = s.Err()
	return
}
