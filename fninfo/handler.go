package function

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func Handle(req []byte) string {
	request := strings.TrimSpace(string(req))
	if "delay" == request {
		// Processing will take 1-5 seconds
		rand.Seed(time.Now().Unix())
		processTime := time.Duration(rand.Intn(4)+1) * time.Second
		time.Sleep(processTime)
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ns := "openfaas-fn"
	if namespace, exists := os.LookupEnv("namespace"); exists {
		ns = namespace
	}

	pods, err := kubeClient.CoreV1().Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to read pods:", err)
	}

	var podList []string
	for _, pod := range pods.Items {
		podList = append(podList, pod.GetName())
	}

	svc, err := kubeClient.CoreV1().Services(ns).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to read services:", err)
	}

	var svcList []string
	for _, s := range svc.Items {
		svcList = append(svcList, s.GetName())
	}

	sec, err := kubeClient.AppsV1beta1().Deployments(ns).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to read deployments:", err)
	}

	var secList []string
	for _, s := range sec.Items {
		secList = append(secList, s.GetName())
	}

	host, err := os.Hostname()
	if err != nil {
		log.Fatalln("failed to read hostname:", err)
	}

	r := Response{
		host,
		podList,
		svcList,
		secList,
		os.Environ(),
		request,
	}

	rb, err := json.Marshal(r)
	if err != nil {
		log.Fatalln("failed to serialize respose:", err)
	}

	return string(rb)
}

type Response struct {
	Hostname    string
	Pods        []string
	Services    []string
	Deployments []string
	Environment []string
	Request     string
}
