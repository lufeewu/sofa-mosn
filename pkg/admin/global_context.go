package admin

import (
	"encoding/json"
	"github.com/alipay/sofa-mosn/pkg/api/v2"
	"sync"
)

var (
	mutex           sync.RWMutex
	effectiveConfig map[string]interface{}
)

func init() {
	Reset()
}

func Reset() {
	mutex.Lock()
	effectiveConfig = make(map[string]interface{})
	effectiveConfig["listener"] = make(map[string]*v2.Listener)
	effectiveConfig["cluster"] = make(map[string]*v2.Cluster)
	effectiveConfig["original_config"] = nil
	mutex.Unlock()
}

func MergeOriginalConf(config interface{}) {
	var originalConf map[string]interface{}
	data, _ := json.Marshal(config)
	json.Unmarshal(data, &originalConf)
	Set("original_config", originalConf)
}

func Set(key string, val interface{}) {
	mutex.Lock()
	effectiveConfig[key] = val
	mutex.Unlock()
}

func SetListenerConfig(listenerName string, listenerConfig *v2.Listener) {
	mutex.Lock()
	listenerConfigMap := effectiveConfig["listener"].(map[string]*v2.Listener)
	if originalConf, ok := listenerConfigMap[listenerName]; ok {
		originalConf.ListenerConfig.FilterChains = listenerConfig.ListenerConfig.FilterChains;
		originalConf.ListenerConfig.StreamFilters = listenerConfig.ListenerConfig.StreamFilters;
	} else {
		listenerConfigMap[listenerName] = listenerConfig
	}
	mutex.Unlock()
}

func SetClusterConfig(clusterName string, clusterConfig *v2.Cluster) {
	mutex.Lock()
	clusterConfigMap := effectiveConfig["cluster"].(map[string]*v2.Cluster)
	clusterConfigMap[clusterName] = clusterConfig
	mutex.Unlock()
}

func GetEffectiveConfig() map[string]interface{} {
	return effectiveConfig
}
