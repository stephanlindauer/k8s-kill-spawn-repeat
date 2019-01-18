package cmd

import (
	"github.com/gorilla/mux"
	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"
	"time"

	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type App struct {
	Kubeconfig string
	Namespace  string
	Client     *v1beta1.ExtensionsV1beta1Client
}

func (app *App) ScalingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)

	deployments := app.Client.Deployments(vars["namespace"])
	deployment, err := deployments.Get(vars["deployment"], v1.GetOptions{})
	must(err)

	scaleInt64, err := strconv.ParseInt(vars["scale"], 0, 32)
	must(err)
	var scaleInt32 = int32(Max(0,Min(1,scaleInt64)) )

	deployment.Spec.Replicas = &scaleInt32

	deployments.Update(deployment)
}

func (app *App) Run(cmd *cobra.Command, args []string) {
	app.Client = app.Kubernetes()

	r := mux.NewRouter()

	r.HandleFunc("/{namespace}/{deployment}/{scale}", app.ScalingHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Fatal(srv.ListenAndServe())
}

func (app *App) Kubernetes() *v1beta1.ExtensionsV1beta1Client {
	var config *rest.Config
	var err error

	if app.Kubeconfig == "" {
		config, err = rest.InClusterConfig()
		must(err)
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", app.Kubeconfig)
		must(err)
	}

	kube, err := v1beta1.NewForConfig(config)
	must(err)

	return kube
}

func (app *App) Bind(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&app.Kubeconfig, "kubeconfig", "",
		"Path to the kubeconfig file. Keep empty for in-cluster configuration.")
}

func NewRootCommand() *cobra.Command {
	cmd := cmdutil.NewRootCommand(new(App))
	cmd.Use = "k8s-kill-spawn-repeat --kubeconfig <Path> --namespace <Namespace>"
	cmd.Short = "This app scale deployments according to incoming GET requests."
	return cmd
}
