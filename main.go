package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"errors"
	"kcexport/api"
	"gopkg.in/yaml.v3"
	flag "github.com/spf13/pflag"
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

	var context_name string
	var kc string
	flag.StringVarP(&context_name, "context", "c", "", "Provide context to export")
	flag.StringVarP(&kc, "kubeconfig", "k", fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")), "Please provide kubeconfig")
	flag.Parse()

	if context_name == "" {
		fmt.Println("You must provide the name of the context to export")
		os.Exit(1)
	}

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

	var kcObj api.Config
	yaml.Unmarshal(kcf, &kcObj)

	//context_name := os.Args[1]
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

	var new_config = api.Config{
		APIVersion: "v1",
		Kind: "Config",
		CurrentContext: context_name,
		Clusters: []api.Clusters{*cluster_data},
		Contexts: []api.Contexts{*context_data},
		Users: []api.Users{*user_data},
	}

	yC, err := yaml.Marshal(new_config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", yC)
	
}
