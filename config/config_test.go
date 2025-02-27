// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package config_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	coreRelayer "github.com/ChainSafe/chainbridge-core/config/relayer"

	"github.com/ChainSafe/sygma-relayer/config"
	"github.com/ChainSafe/sygma-relayer/config/relayer"
	"github.com/stretchr/testify/suite"
)

type GetConfigTestSuite struct {
	suite.Suite
}

func TestRunGetConfigTestSuite(t *testing.T) {
	suite.Run(t, new(GetConfigTestSuite))
}

func (s *GetConfigTestSuite) SetupSuite()    {}
func (s *GetConfigTestSuite) TearDownSuite() {}
func (s *GetConfigTestSuite) SetupTest() {
	os.Clearenv()
}
func (s *GetConfigTestSuite) TearDownTest() {}

func (s *GetConfigTestSuite) Test_GetConfigFromFile_InvalidPath() {
	_, err := config.GetConfigFromFile("invalid")

	s.NotNil(err)
}

func (s *GetConfigTestSuite) Test_GetConfigFromENV() {
	_ = os.Setenv("CBH_DOM_1", "{\n      \"id\": 1,\n      \"from\": \"0xff93B45308FD417dF303D6515aB04D9e89a750Ca\",\n      \"name\": \"evm1\",\n      \"type\": \"evm\",\n      \"endpoint\": \"ws://evm1-1:8546\",\n      \"bridge\": \"0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66\",\n      \"erc20Handler\": \"0x3cA3808176Ad060Ad80c4e08F30d85973Ef1d99e\",\n      \"erc721Handler\": \"0x75dF75bcdCa8eA2360c562b4aaDBAF3dfAf5b19b\",\n      \"genericHandler\": \"0xe1588E2c6a002AE93AeD325A910Ed30961874109\",\n      \"gasLimit\": 9000000,\n      \"maxGasPrice\": 20000000000,\n      \"blockConfirmations\": 2\n    }")
	_ = os.Setenv("CBH_DOM_2", "{\n      \"id\": 2,\n      \"from\": \"0xff93B45308FD417dF303D6515aB04D9e89a750Ca\",\n      \"name\": \"evm2\",\n      \"type\": \"evm\",\n      \"endpoint\": \"ws://evm2-1:8546\",\n      \"bridge\": \"0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66\",\n      \"erc20Handler\": \"0x3cA3808176Ad060Ad80c4e08F30d85973Ef1d99e\",\n      \"erc721Handler\": \"0x75dF75bcdCa8eA2360c562b4aaDBAF3dfAf5b19b\",\n      \"genericHandler\": \"0xe1588E2c6a002AE93AeD325A910Ed30961874109\",\n      \"gasLimit\": 9000000,\n      \"maxGasPrice\": 20000000000,\n      \"blockConfirmations\": 2\n    }")

	_ = os.Setenv("CBH_RELAYER_MPCCONFIG_KEY", "test-pk")
	_ = os.Setenv("CBH_RELAYER_MPCCONFIG_KEYSHAREPATH", "/cfg/keyshares/0.keyshare")
	_ = os.Setenv("CBH_RELAYER_MPCCONFIG_PORT", "9000")

	_ = os.Setenv("CBH_RELAYER_MPCCONFIG_TOPOLOGYCONFIGURATION_ACCESSKEY", "test-access-key")
	_ = os.Setenv("CBH_RELAYER_MPCCONFIG_TOPOLOGYCONFIGURATION_SECKEY", "test-sec-key")
	_ = os.Setenv("CBH_RELAYER_MPCCONFIG_TOPOLOGYCONFIGURATION_ENCRYPTIONKEY", "test-enc-key")

	// load from ENV
	cnf, err := config.GetConfigFromENV()
	if err != nil {
		return
	}

	s.Equal(config.Config{
		RelayerConfig: relayer.RelayerConfig{
			RelayerConfig: coreRelayer.RelayerConfig{
				LogLevel: 1,
				LogFile:  "out.log",
			},
			HealthPort: 9001,
			MpcConfig: relayer.MpcRelayerConfig{
				TopologyConfiguration: relayer.TopologyConfiguration{
					AccessKey:      "test-access-key",
					SecKey:         "test-sec-key",
					DocumentName:   "topology.json",
					BucketRegion:   "us-east-1",
					BucketName:     "mpc-topology",
					ServiceAddress: "buckets.chainsafe.io",
					EncryptionKey:  "test-enc-key",
				},
				Port:         9000,
				KeysharePath: "/cfg/keyshares/0.keyshare",
				Key:          "test-pk",
			},
			BullyConfig: relayer.BullyConfig{
				PingWaitTime:     1 * time.Second,
				PingBackOff:      1 * time.Second,
				PingInterval:     1 * time.Second,
				ElectionWaitTime: 2 * time.Second,
				BullyWaitTime:    25 * time.Second,
			},
		},
		ChainConfigs: []map[string]interface{}{
			{
				"id":                 float64(1),
				"type":               "evm",
				"bridge":             "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66",
				"erc721Handler":      "0x75dF75bcdCa8eA2360c562b4aaDBAF3dfAf5b19b",
				"gasLimit":           9e+06,
				"maxGasPrice":        2e+10,
				"from":               "0xff93B45308FD417dF303D6515aB04D9e89a750Ca",
				"name":               "evm1",
				"endpoint":           "ws://evm1-1:8546",
				"erc20Handler":       "0x3cA3808176Ad060Ad80c4e08F30d85973Ef1d99e",
				"genericHandler":     "0xe1588E2c6a002AE93AeD325A910Ed30961874109",
				"blockConfirmations": float64(2),
			},
			{
				"id":                 float64(2),
				"type":               "evm",
				"bridge":             "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66",
				"erc721Handler":      "0x75dF75bcdCa8eA2360c562b4aaDBAF3dfAf5b19b",
				"gasLimit":           9e+06,
				"maxGasPrice":        2e+10,
				"from":               "0xff93B45308FD417dF303D6515aB04D9e89a750Ca",
				"name":               "evm2",
				"endpoint":           "ws://evm2-1:8546",
				"erc20Handler":       "0x3cA3808176Ad060Ad80c4e08F30d85973Ef1d99e",
				"genericHandler":     "0xe1588E2c6a002AE93AeD325A910Ed30961874109",
				"blockConfirmations": float64(2),
			},
		},
	}, cnf)
}

type ConfigTestCase struct {
	name       string
	inConfig   config.RawConfig
	shouldFail bool
	errorMsg   string
	outConfig  config.Config
}

func (s *GetConfigTestSuite) Test_GetConfigFromFile() {
	testCases := []ConfigTestCase{
		{
			name: "missing chain type",
			inConfig: config.RawConfig{
				ChainConfigs: []map[string]interface{}{{
					"name": "chain1",
				}},
				RelayerConfig: relayer.RawRelayerConfig{
					RawRelayerConfig: coreRelayer.RawRelayerConfig{
						OpenTelemetryCollectorURL: "",
						LogLevel:                  "",
						LogFile:                   "",
					},
					MpcConfig: relayer.RawMpcRelayerConfig{
						Port: "2020",
						TopologyConfiguration: relayer.TopologyConfiguration{
							AccessKey:     "access-key",
							SecKey:        "sec-key",
							EncryptionKey: "enc-key",
						},
					},
					BullyConfig: relayer.RawBullyConfig{
						PingWaitTime:     "1s",
						PingBackOff:      "1s",
						PingInterval:     "1s",
						ElectionWaitTime: "1s",
					},
				},
			},
			shouldFail: true,
			errorMsg:   "chain 'type' must be provided for every configured chain",
			outConfig:  config.Config{},
		},
		{
			name: "invalid relayer type",
			inConfig: config.RawConfig{
				RelayerConfig: relayer.RawRelayerConfig{
					RawRelayerConfig: coreRelayer.RawRelayerConfig{
						LogLevel: "invalid",
					},
					MpcConfig: relayer.RawMpcRelayerConfig{
						TopologyConfiguration: relayer.TopologyConfiguration{
							AccessKey:     "access-key",
							SecKey:        "sec-key",
							EncryptionKey: "enc-key",
						},
					},
				},
				ChainConfigs: []map[string]interface{}{{
					"type": "evm",
					"name": "chain1",
				}},
			},
			shouldFail: true,
			errorMsg:   "unknown log level: invalid",
			outConfig:  config.Config{},
		},
		{
			name: "invalid bully config",
			inConfig: config.RawConfig{
				RelayerConfig: relayer.RawRelayerConfig{
					RawRelayerConfig: coreRelayer.RawRelayerConfig{
						LogLevel: "info",
					},
					MpcConfig: relayer.RawMpcRelayerConfig{
						TopologyConfiguration: relayer.TopologyConfiguration{
							AccessKey:     "access-key",
							SecKey:        "sec-key",
							EncryptionKey: "enc-key",
						},
						Port: "2020",
					},
					BullyConfig: relayer.RawBullyConfig{
						PingWaitTime:     "2z",
						PingBackOff:      "",
						PingInterval:     "",
						ElectionWaitTime: "",
					},
				},
				ChainConfigs: []map[string]interface{}{{
					"type": "evm",
					"name": "chain1",
				}},
			},
			shouldFail: true,
			errorMsg:   "unable to parse bully ping wait time: time: unknown unit \"z\" in duration \"2z\"",
			outConfig:  config.Config{},
		},
		{
			name: "invalid topology config",
			inConfig: config.RawConfig{
				RelayerConfig: relayer.RawRelayerConfig{
					RawRelayerConfig: coreRelayer.RawRelayerConfig{
						LogLevel: "info",
					},
					MpcConfig: relayer.RawMpcRelayerConfig{
						TopologyConfiguration: relayer.TopologyConfiguration{
							AccessKey: "access-key",
							SecKey:    "",
						},
						Port: "2020",
					},
				},
				ChainConfigs: []map[string]interface{}{{
					"type": "evm",
					"name": "chain1",
				}},
			},
			shouldFail: true,
			errorMsg:   "topology configuration secret key not provided",
			outConfig:  config.Config{},
		},
		{
			name: "missing encryption key",
			inConfig: config.RawConfig{
				RelayerConfig: relayer.RawRelayerConfig{
					RawRelayerConfig: coreRelayer.RawRelayerConfig{
						LogLevel: "info",
					},
					MpcConfig: relayer.RawMpcRelayerConfig{
						TopologyConfiguration: relayer.TopologyConfiguration{
							AccessKey: "access-key",
							SecKey:    "secret-key",
						},
						Port: "2020",
					},
				},
				ChainConfigs: []map[string]interface{}{{
					"type": "evm",
					"name": "chain1",
				}},
			},
			shouldFail: true,
			errorMsg:   "topology configuration encryption key not provided",
			outConfig:  config.Config{},
		},
		{
			name: "set default values in config",
			inConfig: config.RawConfig{
				RelayerConfig: relayer.RawRelayerConfig{
					// LogLevel: use default value,
					// LogFile: use default value
					MpcConfig: relayer.RawMpcRelayerConfig{
						TopologyConfiguration: relayer.TopologyConfiguration{
							AccessKey:     "access-key",
							SecKey:        "sec-key",
							EncryptionKey: "enc-key",
						},
						// Port: use default value,
					},
				},
				ChainConfigs: []map[string]interface{}{{
					"type": "evm",
					"name": "evm1",
				}},
			},
			shouldFail: false,
			errorMsg:   "unable to parse bully ping wait time: time: unknown unit \"z\" in duration \"2z\"",
			outConfig: config.Config{
				RelayerConfig: relayer.RelayerConfig{
					RelayerConfig: coreRelayer.RelayerConfig{
						LogLevel:                  1,
						LogFile:                   "out.log",
						OpenTelemetryCollectorURL: "",
					},
					HealthPort: 9001,
					MpcConfig: relayer.MpcRelayerConfig{
						Port: 9000,
						TopologyConfiguration: relayer.TopologyConfiguration{
							AccessKey:      "access-key",
							EncryptionKey:  "enc-key",
							SecKey:         "sec-key",
							DocumentName:   "topology.json",
							BucketRegion:   "us-east-1",
							BucketName:     "mpc-topology",
							ServiceAddress: "buckets.chainsafe.io",
						},
					},
					BullyConfig: relayer.BullyConfig{
						PingWaitTime:     1 * time.Second,
						PingBackOff:      1 * time.Second,
						PingInterval:     1 * time.Second,
						ElectionWaitTime: 2 * time.Second,
						BullyWaitTime:    25 * time.Second,
					},
				},
				ChainConfigs: []map[string]interface{}{{
					"type": "evm",
					"name": "evm1",
				}},
			},
		},
		{
			name: "valid config",
			inConfig: config.RawConfig{
				RelayerConfig: relayer.RawRelayerConfig{
					RawRelayerConfig: coreRelayer.RawRelayerConfig{
						LogLevel: "debug",
						LogFile:  "custom.log",
					},
					HealthPort: "9002",
					MpcConfig: relayer.RawMpcRelayerConfig{
						TopologyConfiguration: relayer.TopologyConfiguration{
							AccessKey:     "access-key",
							SecKey:        "sec-key",
							EncryptionKey: "enc-key",
							BucketName:    "test-mpc-bucket",
						},
						Port:         "2020",
						KeysharePath: "./share.key",
						Key:          "./key.pk",
					},
					BullyConfig: relayer.RawBullyConfig{
						PingWaitTime:     "1s",
						PingBackOff:      "1s",
						PingInterval:     "1s",
						ElectionWaitTime: "1s",
						BullyWaitTime:    "1s",
					},
				},
				ChainConfigs: []map[string]interface{}{{
					"type": "evm",
					"name": "evm1",
				}},
			},
			shouldFail: false,
			errorMsg:   "unable to parse bully ping wait time: time: unknown unit \"z\" in duration \"2z\"",
			outConfig: config.Config{
				RelayerConfig: relayer.RelayerConfig{
					RelayerConfig: coreRelayer.RelayerConfig{
						LogLevel:                  0,
						LogFile:                   "custom.log",
						OpenTelemetryCollectorURL: "",
					},
					HealthPort: 9002,
					MpcConfig: relayer.MpcRelayerConfig{
						Port:         2020,
						KeysharePath: "./share.key",
						Key:          "./key.pk",
						TopologyConfiguration: relayer.TopologyConfiguration{
							AccessKey:      "access-key",
							SecKey:         "sec-key",
							EncryptionKey:  "enc-key",
							DocumentName:   "topology.json",
							BucketRegion:   "us-east-1",
							BucketName:     "test-mpc-bucket",
							ServiceAddress: "buckets.chainsafe.io",
						},
					},
					BullyConfig: relayer.BullyConfig{
						PingWaitTime:     time.Second,
						PingBackOff:      time.Second,
						PingInterval:     time.Second,
						ElectionWaitTime: time.Second,
						BullyWaitTime:    time.Second,
					},
				},
				ChainConfigs: []map[string]interface{}{{
					"type": "evm",
					"name": "evm1",
				}},
			},
		},
	}

	for _, t := range testCases {
		s.Run(t.name, func() {
			file, _ := json.Marshal(t.inConfig)
			_ = ioutil.WriteFile("test.json", file, 0644)

			conf, err := config.GetConfigFromFile("test.json")

			_ = os.Remove("test.json")

			if t.shouldFail {
				s.NotNil(err)
				s.Equal(t.errorMsg, err.Error())
			} else {
				s.Nil(err)
				s.Equal(t.outConfig, conf)
			}
		})
	}
}
