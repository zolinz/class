package handlers

import (
	"context"
	"github.com/zolinz/class/foundation/web"
	"log"
	"net/http"
	"os"
	
)

type check struct {
	log *log.Logger
	build string
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error{

	/*if n := rand.Intn(100) ; n%2 == 0 {
		return web.NewRequestError(	errors.New("trusted error "), http.StatusBadRequest)
		//panic("oh my god")
	}*/
	status := struct {
		Status string
	}{
		Status: "OK",
	}


	return web.Respond(ctx, w, status, http.StatusOK)
	//return json.NewEncoder(w).Encode(status)
}


// liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (c check) liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.check.liveness")
	//defer span.End()

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     c.build,
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	return web.Respond(ctx, w, info, http.StatusOK)
}
