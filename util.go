package main

import (
	"fmt"
	"context"
	"time"
	"k8s.io/kubernetes/pkg/kubelet/util"
	"k8s.io/kubernetes/pkg/kubelet/apis/podresources"
	kubeletpodresourcesv1alpha1 "k8s.io/kubernetes/pkg/kubelet/apis/podresources/v1alpha1"
)

const (
	// Kubelet internal cgroup name for node allocatable cgroup.
	defaultNodeAllocatableCgroup = "kubepods"
	// defaultPodResourcesPath is the path to the local endpoint serving the podresources GRPC service.
	defaultPodResourcesPath    = "/var/lib/kubelet/pod-resources"
	defaultPodResourcesTimeout = 10 * time.Second
	defaultPodResourcesMaxSize = 1024 * 1024 * 16 // 16 Mb
)

func getNodeDevices() (*kubeletpodresourcesv1alpha1.ListPodResourcesResponse, error) {

	endpoint, err := util.LocalEndpoint(defaultPodResourcesPath, podresources.Socket)
	if err != nil {
		return nil, fmt.Errorf("Error getting local endpoint: %v", err)
	}

	client, conn, err := podresources.GetClient(endpoint, defaultPodResourcesTimeout, defaultPodResourcesMaxSize)
	if err != nil {
		return nil, fmt.Errorf("Error getting grpc client: %v", err)
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.List(ctx, &kubeletpodresourcesv1alpha1.ListPodResourcesRequest{})
	if err != nil {
		return nil, fmt.Errorf("%v.Get(_) = _, %v", client, err)
	}

	return resp, nil

}


func main() {
	podResources, err := getNodeDevices()
	if err != nil{
		fmt.Errorf("Error in getting the node devices %v", err)
	}
	fmt.Println("pod resources %v", podResources)
}
