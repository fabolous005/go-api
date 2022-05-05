package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type ParsedJson struct {
	Test1 struct {
		Bla string `json:"bla"`
		Foo string `json:"foo"`
	} `json:"test1"`
	Test2 []string `json:"test2"`
}

type Myyaml struct {
	Http struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
	Path struct {
		JSON   string `yaml:"Json"`
		Ok     string `yaml:"OK"`
		Params string `yaml:"Params"`
	} `yaml:"path"`
	Log struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"log"`
}

func logging(myyaml string) {
	level, err := logger.ParseLevel(myyaml)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(level)
	logger.SetFormatter(&logger.JSONFormatter{})
	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(file)
}

func parseYaml() Myyaml {
	var myyaml Myyaml
	yamlFile, err := ioutil.ReadFile("/workspaces/rest/myyaml.yml")
	test := yaml.Unmarshal([]byte(yamlFile), &myyaml)
	if test != nil {
		panic(err)
	}
	return myyaml
}

func returnRes(w http.ResponseWriter, r *http.Request) {
	logger.Info("called returnRes")
	w.WriteHeader(200)
}

func param(w http.ResponseWriter, r *http.Request) {
	logger.Info("called param")
	query := r.URL.Query().Get("param1")
	if len(query) <= 0 {
		w.WriteHeader(400)
		w.Write([]byte("param error"))
		logger.Error("param error")
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(query))
}

func postJson(w http.ResponseWriter, r *http.Request) {
	logger.Info("called postJson")
	w.Header().Set("Content-Type", "application/json")
	var parsedJson ParsedJson
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parsedJson)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("json error"))
		logger.Error("json error")
		return
	}
	w.WriteHeader(200)
	responseString := parsedJson.Test1.Foo + ": " + strings.Join(parsedJson.Test2, ",")
	fmt.Println(responseString)
}

func main() {
	myyaml := parseYaml()
	logging(myyaml.Log.Level)
	r := mux.NewRouter()
	r.HandleFunc(fmt.Sprint(myyaml.Path.Params), param)
	r.HandleFunc(fmt.Sprint(myyaml.Path.Ok), returnRes)
	r.HandleFunc(fmt.Sprint(myyaml.Path.JSON), postJson).Methods("POST")
	http.ListenAndServe(":"+fmt.Sprint(myyaml.Http.Port), r)
}
