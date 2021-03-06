package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

func main() {
	//server := grpc.NewServer()
	//pb.RegisterAlphaRuleEngineServer(server, &AlphaRuleImpl{})
	//
	//lis, err := net.Listen("tcp", ":"+PORT)
	//if err != nil {
	//	log.Fatalf("net.Listen err: %v", err)
	//}

//	server.Serve(lis)
	InitConf()
	// handle http
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/auth/keys", apiAuthKeys)
	//http.HandleFunc("/api/dashboards/home", apiDashboardsHome)

	// serve http
	//http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("GRPC_CLIENT_PORT")), nil)
	http.ListenAndServe(fmt.Sprintf("%s:%s", getConf().Ip, getConf().Port), nil)
}

func apiAuthKeys(w http.ResponseWriter, r *http.Request)   {
	client := &http.Client{}
	req,_ := http.NewRequest("GET","http://52.41.98.206:3000/api/auth/keys?includeExpired=false",nil)
	//req.Header.Add("Authorization","Bearer eyJrIjoiR0FKRFF5S1F4aFlNSUFjNlVYQ3JUZ2N1azdNWWNaNDMiLCJuIjoiaGFuY2tzaG9uMiIsImlkIjoxfQ==\" http://52.41.98.206:3000/api/dashboards/home")

	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Add("x-grafana-org-id", "1")
	req.Header.Add("Referer", "http://52.41.98.206:3000/org/apikeys")
	req.Header.Add("Accept-Language", "zh-CN,zh-TW;q=0.9,zh;q=0.8,en-US;q=0.7,en;q=0.6")
	req.Header.Add("Cookie", "grafana_session_3000=fd65002d5e5734061c774c8c3e1825a0")
	resp,_ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(body))
	w.WriteHeader(200)
	w.Write(body)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	destUrl := fmt.Sprintf("http://%s:%s%s", getConf().RIp, getConf().RPort, r.RequestURI)
	fmt.Println(destUrl)
	r.Header.Get("")
	client := &http.Client{}
	req,_ := http.NewRequest("GET",destUrl,nil)
	//req.Header.Add("Authorization","Bearer eyJrIjoiR0FKRFF5S1F4aFlNSUFjNlVYQ3JUZ2N1azdNWWNaNDMiLCJuIjoiaGFuY2tzaG9uMiIsImlkIjoxfQ==\" http://52.41.98.206:3000/api/dashboards/home")
	for key, _ := range r.Header {
		req.Header.Add(key, r.Header.Get(key) )
	}
	if req.URL.Query().Get("Authorization") == "" {
		req.Header.Add("Authorization", getConf().Authorization)
	} else {
		req.Header.Add("Authorization", req.URL.Query().Get("Authorization"))
	}
	resp,_ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,DELETE,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "token,Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.WriteHeader(200)
	w.Write(body)
	return
}

//??????conf??????
//???????????????????????????????????????????????????
type conf struct {
	Ip   string `yaml: "ip"`
	Port   string `yaml:"port"`
	RIp    string `yaml:"rip"`
	RPort string `yaml:"rport"`
	Authorization  string `yaml:"authorization"`
}
var gConf conf

func InitConf() {
	//??????yaml????????????
	conf := getConf()
	fmt.Println(conf)

	//?????????????????????json??????
	data, err := json.Marshal(conf)
	if err != nil {
		fmt.Println("err:\t", err.Error())
		return
	}

	//?????????json???????????????
	fmt.Println("data:\t", string(data))
}

//??????Yaml????????????,
//????????????conf??????
func getConf() *conf {
	//????????? ????????????
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
