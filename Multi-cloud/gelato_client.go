package main

import (
    "fmt"
	"os"

	"github.com/opensds/multi-cloud/client"
)

const (
    fname = "out.txt"
)

var (
	c *client.Client
)

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}

func GelatoUpload(bucket, filename string) {
	// <bucket name> <object key> <object>
	resp, err := c.UploadObject(bucket, filename, filename)
	if err != nil {
		fmt.Println("Error Upload", err)
	}

	fmt.Println("Response:", resp)
}

func GelatoDownload(bucket, filename string) {
	// <bucket name> <object key>
	err := c.DownloadObject(bucket, filename)
	if err != nil {
		fmt.Println("Error Download", err)
	}
}

func init () {

	fmt.Println ("Initialize multicloud")

	// Environment variables to be exported

	// export MULTI_CLOUD_IP=192.168.20.158
	// export MICRO_SERVER_ADDRESS=:8089
	// export OS_AUTH_AUTHSTRATEGY=keystone
	// export OS_ACCESS_KEY=ZNRJARg7wkfm9wxzuIeD
	// export OPENSDS_ENDPOINT=http://192.168.20.158:50040
	// export OPENSDS_AUTH_STRATEGY=keystone
	// export OS_AUTH_URL=http://192.168.20.158/identity
	// export OS_USERNAME=admin
	// export OS_PASSWORD=opensds@123
	// export OS_TENANT_NAME=admin
	// export OS_PROJECT_NAME=admin
	// export OS_USER_DOMIN_ID=default

	ip, ok := os.LookupEnv("MULTI_CLOUD_IP")
	if !ok {
		fmt.Errorf("ERROR: You must provide the ip by setting " +
			"the environment variable MULTI_CLOUD_IP")
	}

	cfg := &client.Config{
		Endpoint: "http://" + ip + os.Getenv(client.MicroServerAddress),
	}
	authStrategy := os.Getenv(client.OsAuthAuthstrategy)

	switch authStrategy {
	case client.Keystone:
		cfg.AuthOptions = client.LoadKeystoneAuthOptions()
	case client.Noauth:
		cfg.AuthOptions = client.NewNoauthOptions("adminTenantId")
	default:
		cfg.AuthOptions = client.NewNoauthOptions("tenantId")
	}

	c = client.NewClient(cfg)

	// c.ListBackends()
}
