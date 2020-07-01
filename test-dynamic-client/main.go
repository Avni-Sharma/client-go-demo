package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var namespace = "default"

func main() {
	// var kubeconfig *string
	home := homeDir()
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kubeconfig"))
	if err != nil {
		panic(err.Error())
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	dpl := "secret-demo"

	// construct a deployment using unstructured.Unstructured
	secretRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}

	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name": dpl,
			},
			"data": map[string]interface{}{
				"helo":  []byte("Hello"),
				"world": []byte("World"),
			},
		},
	}
	// create deployment
	// fmt.Println("Creating secret...")
	result, err := clientset.Resource(secretRes).Namespace(namespace).Create(deployment, metav1.CreateOptions{})
	fmt.Println("\nRESULTTTTTTTTTTT %v\n", result)
	if err != nil {
		panic(err)
	}

	fmt.Println("updated deployment")
	existingSecret, err := clientset.Resource(secretRes).Namespace(namespace).Get("secret-demo", metav1.GetOptions{})

	existingSecret2, err := clientset.Resource(secretRes).Namespace(namespace).Get("mysecret", metav1.GetOptions{})

	existingData, _, _ := unstructured.NestedMap(existingSecret.Object, "data")
	fmt.Println("\n DATA1 %v \n", existingData)
	mysecret, _, _ := unstructured.NestedMap(existingSecret2.Object, "data")
	fmt.Println("\n DATA2 %v \n", mysecret)

	eq := reflect.DeepEqual(existingData, mysecret)
	if eq {
		fmt.Println("DO NOT UPDATE")
	} else {
		fmt.Println("UPDATE")
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func int32Ptr(i int32) *int32 { return &i }

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
