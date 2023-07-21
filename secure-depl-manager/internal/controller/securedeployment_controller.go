/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	primev1 "github.com/facundotorraca/secure-depl/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// SecureDeploymentReconciler reconciles a SecureDeployment object
type SecureDeploymentReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	recorder record.EventRecorder
}

const TimeToRefreshInSec = 1

//+kubebuilder:rbac:groups=prime.github.com,resources=securedeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=prime.github.com,resources=securedeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=prime.github.com,resources=securedeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SecureDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *SecureDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log.Log.Info("*** Reconciling secure deployment ***", "req", req.String())

	secureDeployment, err := r.getSecureDeploymentDef(ctx, req)
	if err != nil {
		log.Log.Error(err, "failed to get secure deployment definition")
		return ctrl.Result{}, err
	}

	if !r.deploymentIsAuthorized(secureDeployment) {
		return ctrl.Result{Requeue: true, RequeueAfter: TimeToRefreshInSec * time.Second}, nil
	}

	var deployment appsv1.Deployment
	if err := r.reconcileDeployment(ctx, secureDeployment, &deployment); err != nil {
		r.recorder.Event(secureDeployment, corev1.EventTypeWarning, "Error creating deployment", err.Error())
		log.Log.Error(err, "unable to construct Deployment", "secureDeployment", secureDeployment)
		return ctrl.Result{}, err
	}
	r.recorder.Eventf(
		secureDeployment,
		corev1.EventTypeNormal,
		"Created deployment", "deployment=%s", deployment.Name)

	// Update the Status block of the CustomObject resource
	err = r.Status().Update(ctx, secureDeployment)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SecureDeploymentReconciler) reconcileDeployment(ctx context.Context, secureDeployment *primev1.SecureDeployment, deployment *appsv1.Deployment) error {
	log.Log.Info("Creating Deployment for Project", "project", secureDeployment.Name)

	*deployment = appsv1.Deployment{
		ObjectMeta: secureDeployment.ObjectMeta,
		Spec:       secureDeployment.Spec.DeploymentSpec,
	}

	findDeploymentError := r.Get(ctx, client.ObjectKeyFromObject(deployment), &appsv1.Deployment{})
	deploymentIsNew := findDeploymentError != nil && errors.IsNotFound(findDeploymentError)

	_ = ctrl.SetControllerReference(secureDeployment, deployment, r.Scheme)

	if deploymentIsNew {
		if err := r.Create(ctx, deployment); err != nil {
			log.Log.Error(err, "unable to create ConfigMap for Project", "deployment", deployment)
			return err
		}
	} else {
		if err := r.Update(ctx, deployment); err != nil {
			log.Log.Error(err, "unable to updated ConfigMap for Project", "deployment", deployment)
			return err
		}
	}

	secureDeployment.Status.DeploymentStatus = deployment.Status
	return nil
}

func (r *SecureDeploymentReconciler) deploymentIsAuthorized(secureDeployment *primev1.SecureDeployment) bool {
	res, err := http.Get(secureDeployment.Spec.AuthUrl)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return false
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	if res.StatusCode != 200 {
		fmt.Printf("deployment is not authorized - status code: %d", res.StatusCode)
		return false
	}

	fmt.Printf("deployment is authorized - status code: %d", res.StatusCode)
	return true
}

func (r *SecureDeploymentReconciler) getSecureDeploymentDef(ctx context.Context, req ctrl.Request) (*primev1.SecureDeployment, error) {
	var project primev1.SecureDeployment

	// Log req
	if err := r.Get(ctx, req.NamespacedName, &project); err != nil {
		log.Log.Error(err, "unable to fetch Project")
		// TODO delete all resources associated with this project
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return nil, err
	}

	return &project, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecureDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&primev1.SecureDeployment{}).
		Complete(r)
}
