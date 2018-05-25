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

	host, err := os.Hostname()
	if err != nil {
		log.Fatalln("failed to read hostname:", err)
	}

	r := Response{
		Hostname:    host,
		Environment: os.Environ(),
		Request:     request,
	}

	secPath := "/var/openfaas"
	if sec, exists := os.LookupEnv("secrets_path"); exists {
		secPath = sec
	}

	nss, err := kubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to read namespaces:", err)
	}

	var nsList []Namespace
	for _, n := range nss.Items {
		nsItem := Namespace{
			Name: n.GetName(),
		}

		pods, err := kubeClient.CoreV1().Pods(n.GetName()).List(metav1.ListOptions{})
		if err != nil {
			log.Fatalln("failed to read pods:", err)
		}
		nsItem.Pods = len(pods.Items)

		svc, err := kubeClient.CoreV1().Services(n.GetName()).List(metav1.ListOptions{})
		if err != nil {
			log.Fatalln("failed to read services:", err)
		}
		nsItem.Services = len(svc.Items)

		dep, err := kubeClient.AppsV1beta1().Deployments(n.GetName()).List(metav1.ListOptions{})
		if err != nil {
			log.Fatalln("failed to read deployments:", err)
		}
		nsItem.Deployments = len(dep.Items)

		nsList = append(nsList, nsItem)
	}

	r.Namespaces = nsList

	if _, err := os.Stat(secPath); err == nil {
		secrets, err := walkDir(secPath)
		if err == nil {
			r.Secrets = secrets
		}
	}

	rb, err := json.Marshal(r)
	if err != nil {
		log.Fatalln("failed to serialize respose:", err)
	}

	return string(rb)
}
