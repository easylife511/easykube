package config

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type K8s struct {
	k8sClient *kubernetes.Clientset
	config    *rest.Config
}

func (k8 *K8s) GetK8sconfig(kubeConfigFilePath string) {
	var masterUrl string
	if kubeConfigFilePath == "" {
		kubeConfigFilePath = os.Getenv("KUBECONFIG")
	}
	if kubeConfigFilePath == "" {
		kubeConfigFilePath = clientcmd.RecommendedHomeFile
	}
	_, err := os.Stat(kubeConfigFilePath)
	if err != nil {
		fmt.Printf("[Warning] Config not found: %s\n", kubeConfigFilePath)
		kubeConfigFilePath = ""
		masterUrl = "http://localhost:8080"
	}
	config, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeConfigFilePath)
	if err != nil {
		panic(err)
	}
	k8.config = config
}

func (k8 *K8s) GetK8sClient() {
	clientset, err := kubernetes.NewForConfig(k8.config)
	if err != nil {
		panic(err)
	}
	k8.k8sClient = clientset
}

// --------------------------------- node --------------------------------------
func (k8 *K8s) GetNodesList() *v1.NodeList {
	nodes, err := k8.k8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return nodes
}
func (k8 *K8s) GetNode(nodeName string) *v1.Node {
	node, err := k8.k8sClient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return node
}
func (k8 *K8s) PatchNode(nodeName string, payloadBytes []byte) *v1.Node {
	node, err := k8.k8sClient.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
	return node
}

// --------------------------------- ns --------------------------------------
func (k8 *K8s) GetNamespacesList() *v1.NamespaceList {
	namespaces, err := k8.k8sClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return namespaces
}
func (k8 *K8s) GetNamespace(nsName string) *v1.Namespace {
	ns, err := k8.k8sClient.CoreV1().Namespaces().Get(context.TODO(), nsName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return ns
}

// --------------------------------- pod --------------------------------------
func (k8 *K8s) GetPodsList(ns string) *v1.PodList {
	pods, err := k8.k8sClient.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return pods
}
func (k8 *K8s) GetPod(ns string, podName string) *v1.Pod {
	pod, err := k8.k8sClient.CoreV1().Pods(ns).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return pod
}
func (k8 *K8s) PatchPod(ns string, podName string, payloadBytes []byte) *v1.Pod {
	pod, err := k8.k8sClient.CoreV1().Pods(ns).Patch(context.TODO(), podName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
	return pod
}
func (k8 *K8s) GetPodMetricsList(ns string) *v1beta1.PodMetricsList {
	mc, err := metrics.NewForConfig(k8.config)
	if err != nil {
		panic(err)
	}
	podMetricsList, err := mc.MetricsV1beta1().PodMetricses(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return podMetricsList
}

// --------------------------------- other --------------------------------------
func (k8 *K8s) ExecCmd(cmd string, ns string, podName string, containerName string) (string, bool) {
	req := k8.k8sClient.CoreV1().RESTClient().Post().
		Namespace(ns).
		Resource("pods").
		Name(podName).
		SubResource("exec").VersionedParams(
		&v1.PodExecOptions{
			Container: containerName,
			Command:   strings.Split(cmd, " "),
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(k8.config, "POST", req.URL())
	if err != nil {
		fmt.Printf("[ERROR] create Executor failed, err: %v\n", err)
		return "", false
	}
	var stdout, stderr bytes.Buffer
	if err := exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}); err != nil || stderr.String() != "" {
		fmt.Println("stdout: ", stdout.String())
		fmt.Println("stderr: ", stderr.String())
		fmt.Printf("[ERROR] exec cmd: [ %s ] failed, err: %v\n", cmd, err)
		return "", false
	}
	return stdout.String(), true
}

func (k8 *K8s) GetPVList() *v1.PersistentVolumeList {
	pvList, err := k8.k8sClient.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return pvList
}

func (k8 *K8s) GetPVCList(ns string) *v1.PersistentVolumeClaimList {
	pvcList, err := k8.k8sClient.CoreV1().PersistentVolumeClaims(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return pvcList
}
