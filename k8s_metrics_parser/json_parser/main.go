package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var Link = JSON_METRICS_LINK // https://k8s-vip.int:6443/apis/metrics.k8s.io
var Token = JSON_METRICS_TOKEN

type Items struct {
	Items []Item `json:"items"`
}
type Item struct {
	Metadata   Metadata    `json:"metadata"`
	Containers []Container `json:"containers"`
	Usage      Usage       `json:"usage"`
}
type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
type Container struct {
	Name  string `json:"name"`
	Usage Usage  `json:"usage"`
}
type Usage struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

func main() {
	http.HandleFunc("/pods/metrics", GetPodsMetrics)
	http.HandleFunc("/nodes/metrics", GetNodesMetrics)
	http.HandleFunc("/health", health)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"status\": \"UP\"}")
}

func GetPodsMetrics(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", Link+"/pods", nil)
	if err != nil {
	}
	req.Header.Set("Authorization", "Bearer "+Token)
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var items Items
	json.Unmarshal(body, &items)
	re, _ := regexp.Compile(`\d+`)
	// for _, resource := range []string{"Cpu", "Memory"} {
	fmt.Fprintf(w, "# HELP k8s pods Cpu usage"+"\n")
	fmt.Fprintf(w, "# TYPE k8s_pod_Cpuusage untyped"+"\n")
	for i := 0; i < len(items.Items); i++ {
		if len(items.Items[i].Containers) > 0 {
			for j := 0; j < len(items.Items[i].Containers); j++ {
				res := re.FindAllString(items.Items[i].Containers[j].Usage.Cpu, -1)
				fmt.Fprintf(w, "k8s_pod_Cpuusage{app=\""+items.Items[i].Containers[j].Name+"\",id=\""+items.Items[i].Metadata.Name+"\",namespace=\""+items.Items[i].Metadata.Namespace+"\"} "+res[0]+"\n")
			}
		} else {
			res := re.FindAllString(items.Items[i].Containers[0].Usage.Cpu, -1)
			fmt.Fprintf(w, "k8s_pod_Cpuusage{app=\""+items.Items[i].Containers[0].Name+"\",id=\""+items.Items[i].Metadata.Name+"\",namespace=\""+items.Items[i].Metadata.Namespace+"\"} "+res[0]+"\n")
		}
	}
	fmt.Fprintf(w, "# HELP k8s pods Memory usage"+"\n")
	fmt.Fprintf(w, "# TYPE k8s_pod_Memoryusage untyped"+"\n")
	for i := 0; i < len(items.Items); i++ {
		if len(items.Items[i].Containers) > 0 {
			for j := 0; j < len(items.Items[i].Containers); j++ {
				res := re.FindAllString(items.Items[i].Containers[j].Usage.Memory, -1)
				fmt.Fprintf(w, "k8s_pod_Memoryusage{app=\""+items.Items[i].Containers[j].Name+"\",id=\""+items.Items[i].Metadata.Name+"\",namespace=\""+items.Items[i].Metadata.Namespace+"\"} "+res[0]+"\n")
			}
		} else {
			res := re.FindAllString(items.Items[i].Containers[0].Usage.Memory, -1)
			fmt.Fprintf(w, "k8s_pod_Memoryusage{app=\""+items.Items[i].Containers[0].Name+"\",id=\""+items.Items[i].Metadata.Name+"\",namespace=\""+items.Items[i].Metadata.Namespace+"\"} "+res[0]+"\n")
		}
	}
	// }
}
func GetNodesMetrics(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", Link+"/nodes", nil)
	if err != nil {
	}
	req.Header.Set("Authorization", "Bearer "+Token)
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var items Items
	json.Unmarshal(body, &items)
	re, _ := regexp.Compile(`\d+`)
	fmt.Fprintf(w, "# HELP k8s nodes Cpu usage"+"\n")
	fmt.Fprintf(w, "# TYPE k8s_node_Cpuusage untyped"+"\n")
	for i := 0; i < len(items.Items); i++ {
		res := re.FindAllString(items.Items[i].Usage.Cpu, -1)
		fmt.Fprintf(w, "k8s_node_Cpuusage{node=\""+items.Items[i].Metadata.Name+"\"} "+res[0]+"\n")
	}
	fmt.Fprintf(w, "# HELP k8s nodes Memory usage"+"\n")
	fmt.Fprintf(w, "# TYPE k8s_node_Memoryusage untyped"+"\n")
	for i := 0; i < len(items.Items); i++ {
		res := re.FindAllString(items.Items[i].Usage.Memory, -1)
		fmt.Fprintf(w, "k8s_node_Memoryusage{node=\""+items.Items[i].Metadata.Name+"\"} "+res[0]+"\n")
	}
}
