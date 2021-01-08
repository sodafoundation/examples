package ontap

import (
	"encoding/json"
	"fmt"
	drivers "github.com/netapp/trident/storage_drivers"

	log "github.com/golang/glog"
	"github.com/netapp/trident/storage_drivers/ontap"
	. "github.com/sodafoundation/dock/contrib/drivers/utils/config"
	"github.com/sodafoundation/dock/pkg/model"

	"github.com/ghodss/yaml"
	pb "github.com/sodafoundation/dock/pkg/model/proto"
	"github.com/sodafoundation/dock/pkg/utils/config"
)

const (
	defaultTgtConfDir = "/etc/tgt/conf.d"
	defaultTgtBindIp  = "127.0.0.1"
	username          = "username"
	password          = "password"
	port              = "port"
	timeoutForssh     = 60
)

type ReplicationDriver struct {
	cli              *Cli
	sanStorageDriver *ontap.SANStorageDriver
	conf             *ONTAPConfig
}

func (d *ReplicationDriver) Setup() error {
	// Read NetApp ONTAP config file
	d.conf = &ONTAPConfig{}

	p := config.CONF.OsdsDock.Backends.NetappOntapSan.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	if _, err := Parse(d.conf, p); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("unable to instantiate ontap backend.")
		}
	}()

	empty := ""
	config := &drivers.OntapStorageDriverConfig{
		CommonStorageDriverConfig: &drivers.CommonStorageDriverConfig{
			Version:           d.conf.Version,
			StorageDriverName: StorageDriverName,
			StoragePrefixRaw:  json.RawMessage("{}"),
			StoragePrefix:     &empty,
		},
		ManagementLIF: d.conf.ManagementLIF,
		DataLIF:       d.conf.DataLIF,
		IgroupName:    d.conf.IgroupName,
		SVM:           d.conf.Svm,
		Username:      d.conf.Username,
		Password:      d.conf.Password,
	}
	marshaledJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatal("unable to marshal ONTAP config:  ", err)
	}
	configJSON := string(marshaledJSON)

	// Convert config (JSON or YAML) to JSON
	configJSONBytes, err := yaml.YAMLToJSON([]byte(configJSON))
	if err != nil {
		err = fmt.Errorf("invalid config format: %v", err)
		return err
	}
	configJSON = string(configJSONBytes)

	// Parse the common config struct from JSON
	commonConfig, err := drivers.ValidateCommonSettings(configJSON)
	if err != nil {
		err = fmt.Errorf("input failed validation: %v", err)
		return err
	}

	d.sanStorageDriver = &ontap.SANStorageDriver{
		Config: *config,
	}

	// Initialize the driver.
	if err = d.sanStorageDriver.Initialize(driverContext, configJSON, commonConfig); err != nil {
		log.Errorf("could not initialize storage driver (%s). failed: %v", commonConfig.StorageDriverName, err)
		return err
	}
	log.Infof("storage driver (%s) initialized successfully.", commonConfig.StorageDriverName)

	return nil
}

// Unset
func (r *ReplicationDriver) Unset() error { return nil }

// CreateReplication
func (r *ReplicationDriver) CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error) {
	log.Infof("netapp create replication ....")

	var sourceVolume = getVolumeName(opt.PrimaryVolumeId)
	var destinationVolume = getVolumeName(opt.SecondaryVolumeId)

	var sourceServer = r.conf.Svm
	var destinationServer = r.conf.Svm

	err := r.cli.CreateSnapmirror(sourceVolume, destinationVolume, sourceServer, destinationServer)
	if err != nil {
		log.Error("replication creation failed!!")
		return nil, err
	}

	return &model.ReplicationSpec{
		Name:                           opt.Name,
		Description:                    opt.Description,
		PrimaryVolumeId:                opt.GetPrimaryVolumeId(),
		SecondaryVolumeId:              opt.GetSecondaryVolumeId(),
		AvailabilityZone:               opt.AvailabilityZone,
		PrimaryReplicationDriverData:   opt.GetPrimaryReplicationDriverData(),
		SecondaryReplicationDriverData: opt.GetSecondaryReplicationDriverData(),
		ReplicationMode:                opt.GetReplicationMode(),
		ReplicationPeriod:              opt.ReplicationPeriod,
		ReplicationBandwidth:           opt.GetReplicationBandwidth(),
		PoolId:                         opt.PoolId,
		VolumeDataList:                 opt.GetVolumeDataList(),
		Metadata:                       opt.GetMetadata(),
	}, nil
}

func (r *ReplicationDriver) DeleteReplication(opt *pb.DeleteReplicationOpts) error {
	log.Infof("netapp delete replication ....")

	var sourceVolume = getVolumeName(opt.PrimaryVolumeId)
	var destinationVolume = getVolumeName(opt.SecondaryVolumeId)

	var sourceServer = r.conf.Svm
	var destinationServer = r.conf.Svm

	err := r.cli.DeleteSnapmirror(sourceVolume, destinationVolume, sourceServer, destinationServer)
	if err != nil {
		log.Error("replication deletion failed!!")
		return err
	}

	return nil
}

func (r *ReplicationDriver) EnableReplication(opt *pb.EnableReplicationOpts) error {
	log.Infof("netapp enable replication ....")

	return nil
}

func (r *ReplicationDriver) DisableReplication(opt *pb.DisableReplicationOpts) error {
	log.Infof("netapp disable replication ....")
	return nil
}

func (r *ReplicationDriver) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	log.Infof("netapp failover replication ....")
	return nil
}
