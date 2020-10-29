package database

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/lobsterk/otus-abf/internal/config"
	"strconv"
)

func AerospikeOpenClusterConnection(configs []config.AerospikeConfig, policy *aerospike.ClientPolicy) (*aerospike.Client, error) {
	var hosts []*aerospike.Host
	for _, conf := range configs {
		port, err := strconv.Atoi(conf.Port)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, aerospike.NewHost(conf.Host, port))
	}
	return aerospike.NewClientWithPolicyAndHost(policy, hosts...)
}
