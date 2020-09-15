package main

import (
	"context"
	"flag"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var clientset *kubernetes.Clientset

var namespace string

func deleteJob(jobName string) error {
	err := clientset.BatchV1().Jobs(namespace).Delete(context.TODO(), jobName, metav1.DeleteOptions{})
	if err != nil {
		failOnError(err, err.Error())
	} else {
		err = deletePodsBySelector(jobName)
	}
	return err
}

func createJob(jobName string) {

	backoffLimit := int32(0)
	ttlSecondsAfterFinished := int32(20)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"job_name": jobName}},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "container-no-1",
							Image: "alpine",
							ImagePullPolicy: "Always",
							Env: []apiv1.EnvVar{
								{Name: "MyEnvVar1", Value: "MyValue1"},
								{Name: "MyEnvVar2", Value: "MyValue2"},
							},
							Command:         []string{"sh", "-c"},
							Args:            append([]string{"sleep 30"}),
						},
					},
					RestartPolicy: "Never",
				},
			},
			BackoffLimit: &backoffLimit,
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
		},
	}

	jobsClient := clientset.BatchV1().Jobs(namespace)
	log.Println("Creating job... ")
	_, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Println(err.Error())
		} else {
			failOnError(err, "Failed on job creation")
		}
	}
	//log.Printf("Created job %q.\n", result1)
}

func deletePodsBySelector(labelSelector string) error {
	err := clientset.CoreV1().Pods(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: "job_name="+labelSelector,
	})
	if err != nil {
		failOnError(err, err.Error())
	}
	return err
}

func initK8sClientset() error {
	var err error
	var config *rest.Config

	// use the current context in kubeconfig
	config, err = rest.InClusterConfig()

	if err != nil {
		if  strings.Contains(err.Error(), "unable to load in-cluster configuration") {
			var kubeconfig *string
			if home := homeDir(); home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()
			config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		} else {
			failOnError(err, err.Error())
		}
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	return err
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}

func main() {

	namespace = "default"

	err := initK8sClientset()
	if err != nil {
		failOnError(err, err.Error())
	} else {
		jobName := "benny"
		for {
			createJob(jobName)
			time.Sleep(40 * time.Second)
			log.Printf("Deleting job: '%s'\n", jobName)
			_ = deleteJob(jobName)
			time.Sleep(10 * time.Second)
		}
	}
}