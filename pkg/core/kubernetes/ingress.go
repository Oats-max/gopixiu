package kubernetes

import (
	"context"
	
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/caoyingjunz/gopixiu/api/types"
	"github.com/caoyingjunz/gopixiu/pkg/log"	
)

type IngressGetter interface {
	Ingress(cloud string) IngressInterface
}

type IngressInterface interface {
	List(ctx context.Context, listOptions types.ListOptions) (*networkingv1.IngressList, error)
	Create(ctx context.Context, listOptions *networkingv1.Ingress) error
	Get(ctx context.Context, getOptions types.GetOrDeleteOptions) (*networkingv1.Ingress, error)
	Delete(ctx context.Context, deleteOptions types.GetOrDeleteOptions) error
}

type ingress struct {
	client *kubernetes.Clientset
	cloud  string
}

func NewIngress(c *kubernetes.Clientset, cloud string) *ingress {
	return &ingress{
		client: c,
		cloud:  cloud,
	}
}

func (c *ingress) List(ctx context.Context, listOptions types.ListOptions) (*networkingv1.IngressList, error) {
	if c.client == nil {
		return nil, clientError
	}

	ing, err := c.client.NetworkingV1().Ingresses(listOptions.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Errorf("failed to list %s %s Ingress: %v", listOptions.CloudName, listOptions.Namespace, err)
		return nil, err
	}

	return ing, err
}

func (c *ingress) Create(ctx context.Context, ingress *networkingv1.Ingress) error {
	if c.client == nil {
		return clientError
	}
	ingress, err := c.client.NetworkingV1().Ingresses(ingress.Namespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Errorf("failed to create %s %s ingress: %v", ingress.Namespace, ingress.Name, err)
		return err
	}

	return nil
}

func (c *ingress) Get(ctx context.Context, getOptions types.GetOrDeleteOptions) (*networkingv1.Ingress, error) {
	if c.client == nil {
		return nil, clientError
	}
	ingress, err := c.client.NetworkingV1().Ingresses(getOptions.Namespace).Get(ctx, getOptions.ObjectName, metav1.GetOptions{})
	if err != nil {
		log.Logger.Errorf("failed to get %s ingress: %v", ingress.Namespace, err)
		return nil, err
	}
	return ingress, nil
}

func (c *ingress) Delete(ctx context.Context, deleteOptions types.GetOrDeleteOptions) error {
	if c.client == nil {
		return clientError
	}
	if err := c.client.NetworkingV1().
		Ingresses(deleteOptions.Namespace).
		Delete(ctx, deleteOptions.ObjectName, metav1.DeleteOptions{}); err != nil {
		log.Logger.Errorf("failed to delete %s ingress: %v", deleteOptions.Namespace, err)
		return err
	}

	return nil
}
