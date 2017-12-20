package kube

import (
	"fmt"

	"k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/mbarison/brigade/pkg/brigade"
)

// GetWorker returns the worker description.
//
// This will return an error if no worker is found for the build, which can
// happen when a build is scheduled, but not yet started.
func (s *store) GetWorker(buildID string) (*brigade.Worker, error) {
	labels := labels.Set{"heritage": "brigade", "build": buildID}
	listOption := meta.ListOptions{LabelSelector: labels.AsSelector().String()}
	pods, err := s.client.CoreV1().Pods(s.namespace).List(listOption)
	if err != nil {
		return nil, err
	}
	if len(pods.Items) < 1 {
		return nil, fmt.Errorf("could not find worker for build %s: no pod exists with label %s", buildID, labels.AsSelector().String())
	}
	return NewWorkerFromPod(pods.Items[0]), nil
}

// NewWorkerFromPod creates a new *Worker from a pod definition.
func NewWorkerFromPod(pod v1.Pod) *brigade.Worker {
	l := pod.Labels
	worker := &brigade.Worker{
		ID:        pod.Name,
		BuildID:   l["build"],
		ProjectID: l["project"],
		Status:    brigade.JobStatus(pod.Status.Phase),
	}

	if (worker.Status != brigade.JobPending) && (worker.Status != brigade.JobUnknown) {
		worker.StartTime = pod.Status.StartTime.Time
	}

	if len(pod.Status.ContainerStatuses) > 0 {
		cs := pod.Status.ContainerStatuses[0]
		if cs.State.Terminated != nil {
			worker.EndTime = cs.State.Terminated.FinishedAt.Time
			worker.ExitCode = cs.State.Terminated.ExitCode
		}
	}

	return worker
}
