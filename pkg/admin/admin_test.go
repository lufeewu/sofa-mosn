package admin

import (
	"encoding/json"
	"fmt"
	"github.com/alipay/sofa-mosn/pkg/api/v2"
	"github.com/juju/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func loadConfig(path string) map[string]interface{} {
	log.Println("load config from : ", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("load config failed, ", err)
		os.Exit(1)
	}
	config := make(map[string]interface{})
	// translate to lower case
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalln("json unmarshal config failed, ", err)
		os.Exit(1)
	}
	return config
}

func startServer(configFilePath string) *http.Server {
	var initConfig interface{}
	if configFilePath != "" {
		curDir, _ := os.Getwd()
		configPath := path.Join(curDir, configFilePath);
		initConfig = loadConfig(configPath)
	}
	return Start(initConfig)
}

func getEffectiveConfig() (string, error) {
	resp, err := http.Get("http://localhost:8888/api/v1/config_dump")
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("call admin api failed response status: %d, %s", resp.StatusCode, string(b)))
	}

	if err != nil {
		return "", err
	}
	return string(b), nil
}

func TestSetConfig(t *testing.T) {
	type args struct {
		listenerConfigMap map[string]*v2.Listener
		clusterConfigMap  map[string]*v2.Cluster
	}

	listenerconf1 := &v2.Listener{
		ListenerConfig: v2.ListenerConfig{
			Name:                                  "test_listener",
			AddrConfig:                            "127.0.0.1:2045",
			BindToPort:                            true,
			HandOffRestoredDestinationConnections: false,
			LogPath:                               "stdout",
			FilterChains: []v2.FilterChain{
				{
					Filters: []v2.Filter{
						{
							Name: "proxy",
							Config: map[string]interface{}{
								"downstream_protocol": "Http1",
								"upstream_protocol":   "Http2",
								"virtual_hosts": []map[string]interface{}{
									{
										"name":    "clientHost",
										"domains": []interface{}{"*"},
										"routers": []map[string]interface{}{
											{
												"match": map[string]interface{}{
													"prefix": "/",
												},
												"route": map[string]interface{}{
													"cluster_name": "clientCluster",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	listenerconf2 := &v2.Listener{
		ListenerConfig: v2.ListenerConfig{
			Name:                                  "test_listener",
			AddrConfig:                            "127.0.0.1:2045",
			BindToPort:                            false,
			HandOffRestoredDestinationConnections: true,
			LogPath:                               "stdout",
			FilterChains: []v2.FilterChain{
				{
					Filters: []v2.Filter{
						{
							Name: "proxy",
							Config: map[string]interface{}{
								"downstream_protocol": "Http1",
								"upstream_protocol":   "Http2",
								"virtual_hosts": []map[string]interface{}{
									{
										"name":    "clientHost",
										"domains": []interface{}{"*"},
										"routers": []map[string]interface{}{
											{
												"match": map[string]interface{}{
													"prefix": "/xxx",
												},
												"route": map[string]interface{}{
													"cluster_name": "clientCluster",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	clusterConf1 := &v2.Cluster{
		Name:                 "clientCluster",
		ClusterType:          "SIMPLE",
		SubType:              "",
		LbType:               "LB_RANDOM",
		MaxRequestPerConn:    1024,
		ConnBufferLimitBytes: 32768,
		OutlierDetection:     v2.OutlierDetection{},
		HealthCheck: v2.HealthCheck{
			HealthCheckConfig: v2.HealthCheckConfig{
				Protocol: "",
				TimeoutConfig: v2.DurationConfig{
					Duration: 0,
				},
				IntervalConfig: v2.DurationConfig{
					Duration: 0,
				},
				IntervalJitterConfig: v2.DurationConfig{
					Duration: 0,
				},
			},
		},
		Spec: v2.ClusterSpecInfo{},
		LBSubSetConfig: v2.LBSubsetConfig{
			FallBackPolicy:  0,
			DefaultSubset:   nil,
			SubsetSelectors: nil,
		},
		TLS: v2.TLSConfig{},
		Hosts: []v2.Host{
			{
				HostConfig: v2.HostConfig{
					Address: "127.0.0.1:2046",
					Weight:  1,
					MetaDataConfig: v2.MetadataConfig{
						MetaKey: v2.LbMeta{
							LbMetaKey: nil,
						},
					},
				},
			},
		},
	}
	clusterConf2 := &v2.Cluster{
		Name:                 "clientCluster",
		ClusterType:          "SIMPLE",
		SubType:              "",
		LbType:               "LB_RANDOM",
		MaxRequestPerConn:    1024,
		ConnBufferLimitBytes: 32768,
		OutlierDetection:     v2.OutlierDetection{},
		HealthCheck: v2.HealthCheck{
			HealthCheckConfig: v2.HealthCheckConfig{
				Protocol: "",
				TimeoutConfig: v2.DurationConfig{
					Duration: 0,
				},
				IntervalConfig: v2.DurationConfig{
					Duration: 0,
				},
				IntervalJitterConfig: v2.DurationConfig{
					Duration: 0,
				},
			},
		},
		Spec: v2.ClusterSpecInfo{},
		LBSubSetConfig: v2.LBSubsetConfig{
			FallBackPolicy:  0,
			DefaultSubset:   nil,
			SubsetSelectors: nil,
		},
		TLS: v2.TLSConfig{},
		Hosts: []v2.Host{
			{
				HostConfig: v2.HostConfig{
					Address: "127.0.0.1:2045",
					Weight:  1,
					MetaDataConfig: v2.MetadataConfig{
						MetaKey: v2.LbMeta{
							LbMetaKey: nil,
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "add listener config",
			args: args{
				listenerConfigMap: map[string]*v2.Listener{
					"test_listener_1": listenerconf1,
				},
			},
			want: `{"cluster":{},"listener":{"test_listener":{"name":"test_listener","address":"127.0.0.1:2045","bind_port":true,"handoff_restoreddestination":false,"log_path":"stdout","filter_chains":[{"tls_context":{"status":false,"type":"","extend_verify":null},"filters":[{"type":"proxy","config":{"downstream_protocol":"Http1","upstream_protocol":"Http2","virtual_hosts":[{"domains":["*"],"name":"clientHost","routers":[{"match":{"prefix":"/"},"route":{"cluster_name":"clientCluster"}}]}]}}]}]}},"original_config":null}`,
		},
		{
			name: "update listener config",
			args: args{
				listenerConfigMap: map[string]*v2.Listener{
					"test_listener_1": listenerconf1,
					"test_listener_2": listenerconf2,
				},
			},
			want: `{"cluster":{},"listener":{"test_listener":{"name":"test_listener","address":"127.0.0.1:2045","bind_port":true,"handoff_restoreddestination":false,"log_path":"stdout","filter_chains":[{"tls_context":{"status":false,"type":"","extend_verify":null},"filters":[{"type":"proxy","config":{"downstream_protocol":"Http1","upstream_protocol":"Http2","virtual_hosts":[{"domains":["*"],"name":"clientHost","routers":[{"match":{"prefix":"/xxx"},"route":{"cluster_name":"clientCluster"}}]}]}}]}]}},"original_config":null}`,
		},
		{
			name: "add cluster config",
			args: args{
				clusterConfigMap: map[string]*v2.Cluster{
					"clientCluster_1": clusterConf1,
				},
			},
			want: `{"cluster":{"clientCluster":{"name":"clientCluster","type":"SIMPLE","sub_type":"","lb_type":"LB_RANDOM","max_request_per_conn":1024,"conn_buffer_limit_bytes":32768,"circuit_breakers":null,"outlier_detection":{"Consecutive5xx":0,"Interval":0,"BaseEjectionTime":0,"MaxEjectionPercent":0,"ConsecutiveGatewayFailure":0,"EnforcingConsecutive5xx":0,"EnforcingConsecutiveGatewayFailure":0,"EnforcingSuccessRate":0,"SuccessRateMinimumHosts":0,"SuccessRateRequestVolume":0,"SuccessRateStdevFactor":0},"health_check":{"protocol":"","timeout":"0s","interval":"0s","interval_jitter":"0s","healthy_threshold":0,"unhealthy_threshold":0},"spec":{},"lb_subset_config":{"fall_back_policy":0,"default_subset":null,"subset_selectors":null},"tls_context":{"status":false,"type":"","extend_verify":null},"hosts":[{"address":"127.0.0.1:2046","weight":1,"metadata":{"filter_metadata":{"mosn.lb":null}}}]}},"listener":{},"original_config":null}`,
		},
		{
			name: "update cluster config",
			args: args{
				clusterConfigMap: map[string]*v2.Cluster{
					"clientCluster_1": clusterConf1,
					"clientCluster_2": clusterConf2,
				},
			},
			want: `{"cluster":{"clientCluster":{"name":"clientCluster","type":"SIMPLE","sub_type":"","lb_type":"LB_RANDOM","max_request_per_conn":1024,"conn_buffer_limit_bytes":32768,"circuit_breakers":null,"outlier_detection":{"Consecutive5xx":0,"Interval":0,"BaseEjectionTime":0,"MaxEjectionPercent":0,"ConsecutiveGatewayFailure":0,"EnforcingConsecutive5xx":0,"EnforcingConsecutiveGatewayFailure":0,"EnforcingSuccessRate":0,"SuccessRateMinimumHosts":0,"SuccessRateRequestVolume":0,"SuccessRateStdevFactor":0},"health_check":{"protocol":"","timeout":"0s","interval":"0s","interval_jitter":"0s","healthy_threshold":0,"unhealthy_threshold":0},"spec":{},"lb_subset_config":{"fall_back_policy":0,"default_subset":null,"subset_selectors":null},"tls_context":{"status":false,"type":"","extend_verify":null},"hosts":[{"address":"127.0.0.1:2045","weight":1,"metadata":{"filter_metadata":{"mosn.lb":null}}}]}},"listener":{},"original_config":null}`,
		},
	}

	srv := startServer("")
	defer func() {
		if err := srv.Close(); err != nil {
			fmt.Errorf("server close error: %s", err)
		}
	}()
	time.Sleep(time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.listenerConfigMap {
				SetListenerConfig(v.Name, v)
			}
			for _, v := range tt.args.clusterConfigMap {
				SetClusterConfig(v.Name, v)
			}
			if got, err := getEffectiveConfig(); err != nil || strings.Compare(got, tt.want) != 0 {
				if err != nil {
					t.Errorf("getEffectiveConfig failed with error: %s", err)
				} else {
					t.Errorf("getEffectiveConfig() = %v, want %v", got, tt.want)
				}
			}
			Reset()
		})
	}
}
