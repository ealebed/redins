package k8s

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/ealebed/redins/helpers"
	"github.com/ealebed/redins/redis"
)

func initK8sClient() (*kubernetes.Clientset, error) {
	inCluster := helpers.GetEnvironmentVariableAsBool("IN_CLUSTER", true)

	var config *rest.Config

	if inCluster {
		config, _ = initInClusterClient()
	} else {
		config, _ = initOutOfClusterClient()
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func initInClusterClient() (*rest.Config, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func initOutOfClusterClient() (*rest.Config, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// NewInformerFactory returns new SharedInformerFactory
func NewInformerFactory() informers.SharedInformerFactory {
	cs, _ := initK8sClient()
	kubeInformerFactory := informers.NewSharedInformerFactoryWithOptions(cs, time.Second*30,
		informers.WithNamespace("default"),
		informers.WithTweakListOptions(func(options *v1.ListOptions) {
			options.LabelSelector = "app=ads-redis-statistic"
		}))
	podInformer := kubeInformerFactory.Core().V1().Pods().Informer()

	_, err := podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: onAdd,
	})
	if err != nil {
		fmt.Printf("Failed to add event handler: %v\n", err)
	}

	return kubeInformerFactory
}

func onAdd(obj interface{}) {
	rc := redis.InitRedisClient(
		helpers.GetEnvironmentVariableAsString("REDIS_HOST", "127.0.0.1:6379"),
		helpers.GetEnvironmentVariableAsString("REDIS_PASSWORD", ""),
		helpers.GetEnvironmentVariableAsInteger("REDIS_DB", 4))

	if _, err := rc.Connect(); err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}

	flowRulesValue := "[{\"resource\":\"loopme.grpc.ssp.v0.AdsTxtRecordService/GetAdsTxtRelationships\"," +
		"\"count\":100.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"loopme.grpc.ssp.v0.PublisherAccountService/GetPublisherById\"," +
		"\"count\":5.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"loopme.grpc.ssp.v1.PublisherAccountService/GetPublisherById\"," +
		"\"count\":5.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"loopme.grpc.ssp.v0.BundleLegacyService/GetBundleByKey\"," +
		"\"count\":20.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"loopme.lsm.ssp.v0.BundleService/GetBundleById\"," +
		"\"count\":20.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"loopme.lsm.ssp.v0.BundleService/QueryBundle\"," +
		"\"count\":20.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"loopme.grpc.ssp.v0.AppLegacyService/GetAppById\"," +
		"\"count\":10.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"loopme.grpc.ssp.v0.AppLegacyService/GetAppIdByKey\"," +
		"\"count\":10.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"loopme.grpc.ssp.v0.AppLegacyService/GetAppIdByContainerKey\"," +
		"\"count\":16.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"loopme.grpc.ssp.v0.AppLegacyService/GetAppByContainerKey\"," +
		"\"count\":10.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"ExchangeThrottleRateService/GetThrottleRatesByKeys\"," +
		"\"count\":20.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"dsp-fetcher\",\"count\":25.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"exchange-fetcher\",\"count\":300.0,\"grade\":\"THREAD\",\"limit-app\":\"default\"}," +
		"{\"resource\":\"kafka_dmp_ads_requests_info\",\"count\":500.0,\"grade\":\"QPS\",\"limit-app\":\"default\"}]"

	if err := rc.SetValue("flow-rules-key", flowRulesValue); err != nil {
		fmt.Printf("Failed to set value in Redis: %v\n", err)
		_ = rc.Disconnect()
		return
	}

	fmt.Println(rc.QueryValue("flow-rules-key"))
	if err := rc.Disconnect(); err != nil {
		fmt.Printf("Failed to disconnect from Redis: %v\n", err)
	}
}
