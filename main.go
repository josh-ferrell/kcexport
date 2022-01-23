package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"github.com/ghodss/yaml"
	"kcexport/api"
	"encoding/json"
	"errors"
)

func getContext(context_name string, contexts []api.Contexts) (*api.Contexts, error) {
	for _, context := range contexts {
		if context.Name == context_name {
			return &context, nil
		}
	}
	return nil, errors.New("Context not found")
}

func getCluster(cluster_name string, clusters []api.Clusters) (*api.Clusters, error) {
	for _, cluster := range clusters {
		if cluster.Name == cluster_name {
			return &cluster, nil
		}
	}
	return nil, errors.New("Cluster not found")
}

func getUser(user_name string, users []api.Users) (*api.Users, error) {
	for _, user := range users {
		if user.Name == user_name {
			return &user, nil
		}
	}
	return nil, errors.New("User not found")
}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("You must provide the name of the context to export")
		os.Exit(1)
	}

	kc := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))

	_, err := os.Stat(kc)

	if err != nil {
		fmt.Println("%s does not exist")
		os.Exit(1)
	}

	kcf, err := ioutil.ReadFile(kc)
	if err != nil {
		fmt.Println("Cannot open kubeconfig: %s", kc)
		os.Exit(1)
	}

	// Works to produce JSON
	kcYAML, err := yaml.YAMLToJSON(kcf)
	if err != nil {
		fmt.Println("Error converting YAML")
		os.Exit(1)
	}

	var kcObj api.Config
	json.Unmarshal(kcYAML, &kcObj)

	context_name := os.Args[1]
	context_data, err := getContext(context_name, kcObj.Contexts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	user_name := context_data.Context.User
	user_data, err := getUser(user_name, kcObj.Users)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cluster_name := context_data.Context.Cluster
	cluster_data, err := getCluster(cluster_name, kcObj.Clusters)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var new_config api.Config
	new_config.APIVersion = "v1"
	new_config.Kind = "Config"
	new_config.CurrentContext = context_name
	new_config.Clusters = append(new_config.Clusters, *cluster_data)
	new_config.Contexts = append(new_config.Contexts, *context_data)
	new_config.Users = append(new_config.Users, *user_data)

	yC, err := yaml.Marshal(new_config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", yC)
	
}
