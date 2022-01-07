package conf

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//定义conf类型
//类型里的属性，全是配置文件里的属性
type conf struct {
	Host   string `yaml: "host"`
	User   string `yaml:"user"`
	Pwd    string `yaml:"pwd"`
	Dbname string `yaml:"dbname"`
}
var gConf conf

func InitConf() {
	//读取yaml配置文件
	conf := getConf()
	fmt.Println(conf)

	//将对象，转换成json格式
	data, err := json.Marshal(conf)
	if err != nil {
		fmt.Println("err:\t", err.Error())
		return
	}

	//最终以json格式，输出
	fmt.Println("data:\t", string(data))
}

//读取Yaml配置文件,
//并转换成conf对象
func getConf() *conf {
	//应该是 绝对地址
	yamlFile, err := ioutil.ReadFile("./conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
		return &gConf
	}
	err = yaml.Unmarshal(yamlFile, &gConf)

	if err != nil {
		fmt.Println(err.Error())
		return &gConf
	}

	return &gConf
}
