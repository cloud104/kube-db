package rds

import (
	"fmt"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// create an External named service object for Kubernetes
func (k *Kube) createServiceObj(s *v1.Service, namespace string, hostname string, internalname string) *v1.Service {
	s.Spec.Type = "ExternalName"
	s.Spec.ExternalName = hostname

	s.Name = internalname
	s.Annotations = map[string]string{"origin": "rds"}
	s.Namespace = namespace
	return s
}

// CreateService Creates or updates a service in Kubernetes with the new information
func (k *Kube) CreateService(namespace string, hostname string, internalname string) error {

	// create a service in kubernetes that points to the AWS RDS instance
	serviceInterface := k.Client.CoreV1().Services(namespace)

	s, sErr := serviceInterface.Get(hostname, metav1.GetOptions{})

	create := false
	if sErr != nil {
		s = &v1.Service{}
		create = true
	}
	s = k.createServiceObj(s, namespace, hostname, internalname)
	var err error
	if create {
		_, err = serviceInterface.Create(s)
	} else {
		_, err = serviceInterface.Update(s)
	}

	return err
}

func (k *Kube) DeleteService(namespace string, dbname string) error {
	serviceInterface := k.Client.CoreV1().Services(namespace)
	err := serviceInterface.Delete(dbname, &metav1.DeleteOptions{})
	if err != nil {
		// @TODO: Refactor, cuz ugly as fuck
		code := metav1.Status(err.(k8s_errors.APIStatus).Status()).Code
		if code == 404 {
			return nil
		}
		return errors.Wrap(err, fmt.Sprintf("delete of service %v failed in namespace %v", dbname, namespace))
	}
	return nil
}

func (k *Kube) GetSecret(namespace string, name string, key string) (string, error) {
	secret, err := k.Client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to fetch secret %v", name))
	}
	password := secret.Data[key]
	return string(password), nil
}
