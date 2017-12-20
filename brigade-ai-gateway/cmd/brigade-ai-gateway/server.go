package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"gopkg.in/gin-gonic/gin.v1"
	"k8s.io/api/core/v1"

	"github.com/mbarison/brigade/pkg/storage/kube"
)

var (
	kubeconfig              string
	master                  string
	namespace               string
	buildForkedPullRequests bool
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.StringVar(&namespace, "namespace", defaultNamespace(), "kubernetes namespace")
}

func main() {
	flag.Parse()

	//clientset
	_, err := kube.GetClient(master, kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	//store := kube.New(clientset, namespace)

	router := gin.New()
	router.Use(gin.Recovery())

	events := router.Group("/events")
	{
		events.Use(gin.Logger())
		events.POST("/training", Handle)
	}

	router.GET("/healthz", healthz)

	router.Run(":7766")
}

// Handle routes a webhook to its appropriate handler.
//
// It does this by sniffing the event from the header, and routing accordingly.
func Handle(c *gin.Context) {
	event := c.Request.Header.Get("X-Training-Event")
	switch event {
	case "ping":
		log.Print("Received ping from Scheduler")
		c.JSON(200, gin.H{"message": "OK"})
		return
	//case "push", "pull_request", "create", "release", "status", "commit_comment", "pull_request_review":
	//	s.handleEvent(c, event)
	//	return
	default:
		// Issue #127: Don't return an error for unimplemented events.
		log.Printf("Unsupported event %q", event)
		c.JSON(200, gin.H{"message": "Ignored"})
		return
	}
}

func defaultNamespace() string {
	if ns, ok := os.LookupEnv("BRIGADE_NAMESPACE"); ok {
		return ns
	}
	return v1.NamespaceDefault
}

func healthz(c *gin.Context) {
	c.String(http.StatusOK, http.StatusText(http.StatusOK))
}
