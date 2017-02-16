package annotator

import (
	"time"

	"github.com/op/go-logging"

	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/runtime"
	"k8s.io/client-go/1.5/pkg/watch"
	"k8s.io/client-go/1.5/rest"
	"k8s.io/client-go/1.5/tools/cache"

	mbapi "github.com/microscaling/microbadger/api"
)

var (
	log = logging.MustGetLogger("annotator")
)

const (
	resyncPeriod = 5 * time.Minute
)

type Operator struct {
	kclient *kubernetes.Clientset
	podInf  cache.SharedIndexInformer
}

// New creates a new operator.
func New(cfg *rest.Config) (*Operator, error) {
	kclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	o := &Operator{
		kclient: kclient,
	}

	o.podInf = cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (runtime.Object, error) {
				return kclient.Pods(api.NamespaceAll).List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return kclient.Pods(api.NamespaceAll).Watch(options)
			},
		},
		&v1.Pod{}, resyncPeriod, cache.Indexers{},
	)

	o.podInf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: o.handleAddPod,
	})

	return o, nil
}

// Run the operator.
func (o *Operator) Run(stopc <-chan struct{}) error {
	go o.podInf.Run(stopc)

	<-stopc
	return nil
}

// Annotate new pods as they are detected.
func (o *Operator) handleAddPod(obj interface{}) {
	pod := obj.(*v1.Pod)

	err := o.annotatePod(pod)
	if err != nil {
		log.Errorf("Failed to annotate pod %s: %v", pod.Name, err)
	}
}

func (o *Operator) annotatePod(pod *v1.Pod) error {
	// FIXME Hack to only annotate the first container
	c := pod.Spec.Containers[0]
	image := c.Image

	l, err := mbapi.GetLabels(image)
	if err == nil {
		// Get a fresh copy of the ingress before updating.
		np, err := o.kclient.Pods(pod.Namespace).Get(pod.Name)
		if err != nil {
			log.Errorf("Failed to get pod: %v", err)
			return err
		}

		a := np.ObjectMeta.Annotations

		for k, v := range l {
			a[k] = v
		}

		np.ObjectMeta.SetAnnotations(a)

		_, err = o.kclient.Pods(np.Namespace).Update(np)
		if err != nil {
			log.Errorf("Error updating pod: %v", err)
			return err
		}
	}

	return nil
}
