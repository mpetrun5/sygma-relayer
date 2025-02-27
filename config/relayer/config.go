// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package relayer

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/ChainSafe/chainbridge-core/config/relayer"
)

type RelayerConfig struct {
	relayer.RelayerConfig
	HealthPort  uint16
	MpcConfig   MpcRelayerConfig
	BullyConfig BullyConfig
}

type MpcRelayerConfig struct {
	TopologyConfiguration TopologyConfiguration
	Port                  uint16
	KeysharePath          string
	Key                   string
}

type BullyConfig struct {
	PingWaitTime     time.Duration
	PingBackOff      time.Duration
	PingInterval     time.Duration
	ElectionWaitTime time.Duration
	BullyWaitTime    time.Duration
}

type TopologyConfiguration struct {
	EncryptionKey  string `mapstructure:"EncryptionKey" json:"encryptionKey"`
	AccessKey      string `mapstructure:"AccessKey" json:"accessKey"`
	SecKey         string `mapstructure:"SecKey" json:"secKey"`
	DocumentName   string `mapstructure:"DocumentName" default:"topology.json" json:"documentName"`
	BucketRegion   string `mapstructure:"BucketRegion" default:"us-east-1" json:"bucketRegion"`
	BucketName     string `mapstructure:"BucketName" default:"mpc-topology" json:"bucketName"`
	ServiceAddress string `mapstructure:"ServiceAddress" default:"buckets.chainsafe.io" json:"serviceAddress"`
	Path           string `mapstructure:"Path" json:"path"`
}

type RawRelayerConfig struct {
	relayer.RawRelayerConfig `mapstructure:",squash"`
	HealthPort               string              `mapstructure:"HealthPort" json:"healthPort" default:"9001"`
	MpcConfig                RawMpcRelayerConfig `mapstructure:"MpcConfig" json:"mpcConfig"`
	BullyConfig              RawBullyConfig      `mapstructure:"BullyConfig" json:"bullyConfig"`
}

type RawMpcRelayerConfig struct {
	KeysharePath          string                `mapstructure:"KeysharePath" json:"keysharePath"`
	Key                   string                `mapstructure:"Key" json:"key"`
	Port                  string                `mapstructure:"Port" json:"port" default:"9000"`
	TopologyConfiguration TopologyConfiguration `mapstructure:"TopologyConfiguration" json:"topologyConfiguration"`
}

type RawBullyConfig struct {
	PingWaitTime     string `mapstructure:"PingWaitTime" json:"pingWaitTime" default:"1s"`
	PingBackOff      string `mapstructure:"PingBackOff" json:"pingBackOff" default:"1s"`
	PingInterval     string `mapstructure:"PingInterval" json:"pingInterval" default:"1s"`
	ElectionWaitTime string `mapstructure:"ElectionWaitTime" json:"electionWaitTime" default:"2s"`
	BullyWaitTime    string `mapstructure:"BullyWaitTime" json:"bullyWaitTime" default:"25s"`
}

func (c *RawRelayerConfig) Validate() error {
	if c.MpcConfig.TopologyConfiguration.AccessKey == "" {
		return errors.New("topology configuration access key not provided")
	}

	if c.MpcConfig.TopologyConfiguration.SecKey == "" {
		return errors.New("topology configuration secret key not provided")
	}
	if c.MpcConfig.TopologyConfiguration.EncryptionKey == "" {
		return errors.New("topology configuration encryption key not provided")
	}
	return nil
}

// NewRelayerConfig parses RawRelayerConfig into RelayerConfig
func NewRelayerConfig(rawConfig RawRelayerConfig) (RelayerConfig, error) {
	config := RelayerConfig{}
	err := rawConfig.Validate()
	if err != nil {
		return config, err
	}

	logLevel, err := zerolog.ParseLevel(rawConfig.LogLevel)
	if err != nil {
		return config, fmt.Errorf("unknown log level: %s", rawConfig.LogLevel)
	}
	config.LogLevel = logLevel

	config.LogFile = rawConfig.LogFile
	config.OpenTelemetryCollectorURL = rawConfig.OpenTelemetryCollectorURL

	healthPort, err := strconv.ParseInt(rawConfig.HealthPort, 0, 16)
	if err != nil {
		return RelayerConfig{}, fmt.Errorf("unable to parse health port %v", err)
	}
	config.HealthPort = uint16(healthPort)

	mpcConfig, err := parseMpcConfig(rawConfig)
	if err != nil {
		return RelayerConfig{}, err
	}
	config.MpcConfig = mpcConfig

	bullyConfig, err := parseBullyConfig(rawConfig)
	if err != nil {
		return RelayerConfig{}, err
	}
	config.BullyConfig = bullyConfig

	return config, nil
}

func parseMpcConfig(rawConfig RawRelayerConfig) (MpcRelayerConfig, error) {
	var mpcConfig MpcRelayerConfig

	port, err := strconv.ParseInt(rawConfig.MpcConfig.Port, 0, 16)
	if err != nil {
		return MpcRelayerConfig{}, fmt.Errorf("unable to parse mpc port from config %v", err)
	}
	mpcConfig.Port = uint16(port)

	mpcConfig.TopologyConfiguration = rawConfig.MpcConfig.TopologyConfiguration
	mpcConfig.KeysharePath = rawConfig.MpcConfig.KeysharePath
	mpcConfig.Key = rawConfig.MpcConfig.Key

	return mpcConfig, nil
}

func parseBullyConfig(rawConfig RawRelayerConfig) (BullyConfig, error) {
	electionWaitTime, err := time.ParseDuration(rawConfig.BullyConfig.ElectionWaitTime)
	if err != nil {
		return BullyConfig{}, fmt.Errorf("unable to parse bully election wait time: %w", err)
	}

	pingWaitTime, err := time.ParseDuration(rawConfig.BullyConfig.PingWaitTime)
	if err != nil {
		return BullyConfig{}, fmt.Errorf("unable to parse bully ping wait time: %w", err)
	}

	pingInterval, err := time.ParseDuration(rawConfig.BullyConfig.PingInterval)
	if err != nil {
		return BullyConfig{}, fmt.Errorf("unable to parse bully ping interval time: %w", err)
	}

	pingBackOff, err := time.ParseDuration(rawConfig.BullyConfig.PingBackOff)
	if err != nil {
		return BullyConfig{}, fmt.Errorf("unable to parse bully ping back off time: %w", err)
	}

	bullyWaitTime, err := time.ParseDuration(rawConfig.BullyConfig.BullyWaitTime)
	if err != nil {
		return BullyConfig{}, fmt.Errorf("unable to parse bully wait time: %w", err)
	}

	return BullyConfig{
		PingWaitTime:     pingWaitTime,
		PingBackOff:      pingBackOff,
		PingInterval:     pingInterval,
		ElectionWaitTime: electionWaitTime,
		BullyWaitTime:    bullyWaitTime,
	}, nil
}
