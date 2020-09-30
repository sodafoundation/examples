# KubeEdge and SODA Integration PoC

  
# Introduction

This provides a basic integration of [KubeEdge](https://kubeedge.io/en/) project and [SODA](https://sodafoundation.io/). It gives SODA at Edge feasibility for heterogeneous data/storage management by showcasing a workload from Edge Computing platform (KubeEdge) provisioning the volume and mouting using SODA CSI plugin interface.
###  Github:
SODA : [https://github.com/sodafoundation](https://github.com/sodafoundation)
KubeEdge : [https://github.com/kubeedge](https://github.com/kubeedge)

### SODA @ Edge Design
Please refer the the SODA@Edge Design Draft at https://github.com/sodafoundation/design-specs/blob/master/specs/soda-edge/SODAEdgeAnalysisAndDesign.md

  
# Basic Environment

-   Ubutnu 18.04LTS fresh installed Virtual Machine running under VirtualBox (6.1 version with guest additions added) [Similar Ubuntu18.04LTS machine should be ok]
    
-   Root login (#) [This is where it is tested. If super user, should also be ok. In that case please use sudo with some of the commands as needed]
    
-   login to machine as root
    
-   pwd will be /root
    
-   KubeEdge (both cloud core and edge core) and SODA installations are done on a single VM in this poc
    

# KubeEdge Setup

Version Considered : latest (26th Sep 2020) or 1.4.0 release

## Prerequisite Installation

##### go lang 1.14+ needed
-  	 Ref : [https://golang.org/doc/install](https://golang.org/doc/install)
-   Create $HOME/go folder
      
	- cd /root (assuming this as home folder) 
	- mkdir go (this is your workspace)
	
-   go download
	- Downloaded [https://golang.org/dl/go1.15.2.linux-amd64.tar.gz](https://golang.org/dl/go1.15.2.linux-amd64.tar.gz)
	Or 
	- curl -O  [https://golang.org/dl/go1.15.2.linux-amd64.tar.gz](https://golang.org/dl/go1.15.2.linux-amd64.tar.gz)
-   extract
	-   tar -xvf go1.15.2.linux-amd64.tar.gz
-   Adjust the access rights and move to /usr/local
	-   chown -R root:root ./go
	-   mv go /usr/local
-   export the paths
	-   export GOPATH=$HOME/go
	-   export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
	OR
	-   Update ~/.profile with above exports and do
	-   source ~/.profile
`  		[Better to add in .profile file, so can use each exporter]
-   go version
	-   go version go1.15.2 linux/amd64

##### git
-   apt-get update
-   apt install git
-   git version
	-   git version 2.17.1
  
#####    Iptables
-   apt-get install iptables
	-   iptables is already the newest version (1.6.1-2ubuntu2).
#####    openssl
   -   apt get install openssl
	   -   openssl is already the newest version (1.1.1-1ubuntu2.1~18.04.6).
    
  #####   sudo
   - apt-get install sudo
   - sudo --version
	  - sudo version 1.8.21p2
       
   ##### make
   -   sudo apt-get update
   -   sudo apt-get install make
   -   make --version
      GNU Make 4.1
     
   ##### gcc
-   sudo apt update
-   sudo apt install build-essential
-   sudo apt-get install manpages-dev
-   gcc –version
    gcc (Ubuntu 7.5.0-3ubuntu1~18.04) 7.5.0

##### ifconfig
-   apt install net-tools
    
#####   Kubectl
    
-   Ref: https://kubernetes.io/docs/tasks/tools/install-kubectl/
 -   sudo apt-get update && sudo apt-get install -y apt-transport-https gnupg2
 -   curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
 -   echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
 -   sudo apt-get update
 -   sudo apt-get install -y kubectl
 -   Check version
	 -   kubectl version --client
    		-   Client Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.2", GitCommit:"f5743093fd1c663cb0cbc89748f730662345d44d", GitTreeState:"clean", BuildDate:"2020-09-16T13:41:02Z", GoVersion:"go1.15", Compiler:"gc", Platform:"linux/amd64"}
    
##### docker 18.09+ needed
-   Ref : [https://docs.docker.com/engine/install/ubuntu/](https://docs.docker.com/engine/install/ubuntu/)
   -   sudo apt-get update
   -   sudo apt-get install \\
        apt-transport-https \\
        ca-certificates \\
         curl \\
        gnupg-agent \\
       Software-properties-common
-   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
-   sudo apt-key fingerprint 0EBFCD88
	-   pub rsa4096 2017-02-22 [SCEA]
    9DC8 5822 9FC7 DD38 854A E2D8 8D81 803C 0EBF CD88
uid [ unknown] Docker Release (CE deb) <docker@docker.com>
sub rsa4096 2017-02-22 [S]

-   Add repo
	-   sudo add-apt-repository \\
    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \\
    $(lsb_release -cs) \\
    stable"

-   Install docker
	-   sudo apt-get update
    -   sudo apt-get install docker-ce docker-ce-cli containerd.io
    -   docker --version
	       Docker version – 19.03.13
    
##### kind v0.7.0 needed
   -   go get sigs.k8s.io/kind@v0.7.0
   -   go get will put kind in $(go env GOPATH)/bin.
    
## KubeEdge Installation
-   Get the latest version of kubeedge
	-   git clone https://github.com/kubeedge/kubeedge.git 		$GOPATH/src/github.com/kubeedge/kubeedge
    -   cd kubeedge/kubeedge/hack
   -   We will use the local setup to install cloud core and edge core to the same VM
    -   run the local setup script
    -   ./local-up-kubeedge.sh
    -   This will do the complete setup and basic configurations to start cloud core and edge core
    -   You may face an error : sigs.k8s.io/kind@v0.7.0: is explicitly required in go.mod, but not marked as explicit in vendor/modules.txt. Follow the below steps to handle this:
	    -   go mod vendor
	    -   Again run ./local-up-kubeedge.sh
      -   Installation log for your reference
    

    kubectl cluster-info --context kind-test Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂 wait the control-plane ready... node/test-control-plane condition met daemonset.apps "kindnet" deleted namespace/kubeedge created creating the device crd... customresourcedefinition.apiextensions.k8s.io/devices.devices.kubeedge.io created customresourcedefinition.apiextensions.k8s.io/devicemodels.devices.kubeedge.io created creating the objectsync crd... customresourcedefinition.apiextensions.k8s.io/clusterobjectsyncs.reliablesyncs.kubeedge.io created customresourcedefinition.apiextensions.k8s.io/objectsyncs.reliablesyncs.kubeedge.io created Generating RSA private key, 2048 bit long modulus (2 primes) ...............................+++++ .......................................+++++
    
    e is 65537 (0x010001)
    
    Signature ok
    
    subject=C = CN, ST = Zhejiang, L = Hangzhou, O = KubeEdge
    
    Getting CA Private Key
    
    start cloudcore...
    
    2020-09-28 11:12:00.091586 I | INFO: Install client plugin, protocol: rest
    
    2020-09-28 11:12:00.091730 I | INFO: Installed service discovery plugin: edge
    
    I0928 11:12:00.092154 2169 util.go:397] Looking for default routes with IPv4 addresses
    
    I0928 11:12:00.092194 2169 util.go:402] Default route transits interface "enp0s3"
    
    I0928 11:12:00.092333 2169 util.go:212] Interface enp0s3 is up
    
    I0928 11:12:00.092399 2169 util.go:259] Interface "enp0s3" has 2 addresses :[192.168.1.100/24 fe80::c479:3d92:5d85:bef2/64].
    
    I0928 11:12:00.092447 2169 util.go:228] Checking addr
    192.168.1.100/24.
    
    I0928 11:12:00.092473 2169 util.go:235] IP found 192.168.1.100
    
    I0928 11:12:00.092498 2169 util.go:265] Found valid IPv4 address
    192.168.1.100 for interface "enp0s3".
    
    I0928 11:12:00.092523 2169 util.go:408] Found active IP 192.168.1.100
    
    start edgecore...
    
    Local KubeEdge cluster is running. Press Ctrl-C to shut it down.
    
    Logs:
    
    /tmp/cloudcore.log
    
    /tmp/edgecore.log
    
      
    
    To start using your kubeedge, you can run:
    
      
    
    export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/usr/local/go/bin:/usr/local/go/bin:/root/go/bin:/usr/local/go/bin:/root/go/bin:/root/go/bin
    
    export KUBECONFIG=/root/.kube/config
    
    kubectl get nodes

 -   Export the path
	    -   export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/usr/local/go/bin:/usr/local/go/bin:/root/go/bin:/usr/local/go/bin:/root/go/bin:/root/go/bin
	    -   export KUBECONFIG=/root/.kube/config
    
		-   [You can add this to ~/.profile] to get effect in all the terminals
    
-   Verify KubeEdge Installation
    -   kubectl get nodes
NAME STATUS ROLES AGE VERSION
edge-node Ready agent,edge 23m v1.18.6-kubeedge-v1.4.0 beta.0.157+3757d046121a13-dirty
test-control-plane Ready master 25m v1.17.0
	-   There are two nodes master and agent/edge
    

# SODA Setup

Version Considered : Faroe v1.0.0 [It can work with latest or Greenland release also]

Followed quick installation based on : [https://docs.sodafoundation.io/soda-gettingstarted/quickstart/](https://docs.sodafoundation.io/soda-gettingstarted/quickstart/)

## Prerequisite installation

-   apt-get update && apt-get install -y git make curl wget libltdl7 libseccomp2 libffi-dev gawk
    -   Some versions may be already there.
    -   Hence, please keep the latest version (unless specified)
   -   Install Docker
	    -   Skipped as we have already installed.
	    -   Version mentioned for SODA and the current installation are difference
	    -   But leave like Install docker-compose:that.
   -   Install dockercompose
	    -   SKIPPED
  -   Golang
    	-   Recommended 1.13.x
    	-   We have already installed 1.15.1
    	- Keep this later version.
      

## SODA Faroe Version Installation
(we use the Faroe release. Any later version should work)
-   Get SODA Release binaries
	-   cd root/
    -   mkdir soda
    -   cd soda
    -   wget https://github.com/sodafoundation/installer/releases/download/v1.0.0/installer-v1.0.0.tar.gz
    -   OR Download the installer binaries from https://github.com/sodafoundation/installer/releases/tag/v1.0.0
    -   tar xvzf installer-v1.0.0.tar.gz
    -   cd installer-v1.0.0/ansible/
    -   chmod +x ./install_ansible.sh && ./install_ansible.sh
    -   ansible –version # Ansible version 2.4.x is required
    -   ansible 2.5.1
    
-   Take the host ip address and set
    -   You can use internal IP of the machine (192.x.x.x usually) OR loop back address (127.0.0.1)
    -   In the current testing the machine IP address taken-192.168.1.100
    -   export host ip
    -   export HOST_IP=192.168.1.100 (yours may be different)

-   Modify Configuration Files
	-   Path for the below config files: ./group_vars/
    (In the current setup, full path : /root/soda/installer-v1.0.0/ansible/group_vars)
	-   common.yml
	    -   host_ip: 192.168.1.100
	    -   deploy_project: hotpot
	    (We need only hotpot where api, controller, dock will be installed)
	    -   repo_branch: master
	    -   release_version: v1.0.0
	-   sushi.yml
		-   sushi_plugin_type: csi
			(Need to support csi nbp)
		-   enabled_backends: lvm,nfs
		    (Ensure that lvm is there in the list)
	-   osdsdb.yml
	    -   etcd_host: 192.168.1.100
		  -   etcd_port: 3379
	    -   etcd_peer_port: 3380
	 -   auth.yml
		    -   opensds_auth_strategy: noauth
  
-   Start another etcd server for SODA
    -   Already there will be another etcd running with default ports of 2379 and peer port of 2380 from KubeEdge setup. This has auth enabled and hence, soda client may not be failing to connect. Hence we need a separate etcd server running on another port.
    -   For this poc, we use etcd ports as 3379 and peer port 3380
    -   mkdir /root/etcd
    -   Download etcd of the correct version
    -   curl -L https://storage.googleapis.com/etcd/v3.2.25/etcd-v3.2.25-linux-amd64.tar.gz -o /root/etcd/etcd-v3.2.25-linux-amd64.tar.gz
    -   cd /root/etcd
    -   tar xzvf etcd-v3.2.25-linux-amd64.tar.gz
    -   Run etcd
	    -   cd etcd-v3.2.25-linux-amd64
	    -   ./etcd --advertise-client-urls http://192.168.1.100:3379 --listen-client-urls http://192.168.1.100:3379 –listen-peer-urls 192.168.1.100:3380 &
    (Running etcd in the background)
    -   press Enter to get the command prompt back
    -   Check whether etcd is running correctly
    -   ps -aux|grep etcd
    -   You should see our etcd process with the same arguments running on the same IP which we gave. The IP and ports are important.
    -   There may be another etcd running (which got created from KubeEdge installation)
    -   So we will have two etcd running in this setup.
    
###   Install SODA

-   Export SODA Env variables
	-   export OPENSDS_ENDPOINT=http://192.168.1.100:50040
    -   export HOST_IP=192.168.1.100
   -   cd /root/soda/installer-v1.0.0/ansible/
   -   ansible-playbook site.yml -i local.hosts -vvv
	      This will install all needed projects for SODA
	   -   Tips: Any error or you want to clean up
	   -   ansible-playbook clean.yml -i local.hosts -vvv
   -   Basic verification of SODA Installation
	   -   Verify the process are running
	    -   ps -aux | grep osds
    
    root 1453 0.0 0.0 14424 1112 pts/2 S+ 10:28 0:00 grep --color=auto osds
    
    root 19244 0.0 0.2 649996 18784 pts/1 Sl Sep29 0:05 bin/osdslet
    
    root 19334 0.0 0.3 726220 26460 pts/1 Sl Sep29 0:08 bin/osdsapiserver
    
    root 20616 0.0 0.4 933544 37248 pts/1 Sl Sep29 0:11 bin/osdsdock

-   3 osds processes would be running
 -   vi /etc/opensds/opensds.conf
 -   This file should have all the configurations given during the installation updated.
    

# Integration and Testing

This section provides the specific configurations and changes needed to connect SODA and KubeEdge. It also discusses open issues.

-   There are two nodes created by KubeEdge
    -   kubectl get nodes
       NAME STATUS ROLES AGE VERSION
edge-node Ready agent,edge 23m v1.18.6-kubeedge-v1.4.0-beta.0.157+3757d046121a13-dirty
test-control-plane Ready master 25m v1.17.0

-   Important information / issue with edge-node connectivity issue
	-   Ideally the following pods should be running on edge-node
		-   csi-attacher-opensdsplugin-block-0, csi-nodeplugin-opensdsplugin-block-xx, csi-provisioner-opensdsplugin-block-0
	   -   But, it is not running. It gives error on environment variables as “unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined”
	   -   Based on the discussion with KubeEdge team, there is some issue in injecting the env variables to edge-node. The fix is in progress and available at [https://github.com/kubeedge/kubeedge/pull/1834](https://github.com/kubeedge/kubeedge/pull/1834)
	   -   Currently, we did *not* test with this patch. We can do this and verify as well //FIXME
   -   Due to the above issue, we tested with the master node whether the pods can run and provision the volume through SODA CSI plugin and on CSI backends (lvm under this test) through SODA.
		-   For this we need to do node selector and assign the node for our pods.
		-   Setting the node selector label for the KubeEdge Nodes
	    -   See the kubeedge nodes
		    -   kubectl get nodes
	    -   set the label ntype=edge for edge-node
		    -   kubectl label node edge-node ntype=edge
	    -   set the label ntype=master for test-control-plane
			   -   kubectl label node test-control-plane ntype=master
    -   Remove the SODA pods
    -   Clean up file and block pods
	    -   kubectl delete -f /opt/opensds-sushi-linux-amd64/csi/deploy/kubernetes/block
	    -   kubectl delete -f /opt/opensds-sushi-linux-amd64/csi/deploy/kubernetes/file
    -   Update the node assignment to the soda pod yamls
	    -   We are assigning the test-control-plane node as the node on which the pod to be deployed
	    -   cd /opt/opensds-sushi-linux-amd64/csi/deploy/kubernetes/block/
	    -   vi csi-attacher-opensdsplugin.yaml
		    -   Add:
			   
    - nodeSelector:
    				ntype: master
			
      -   You can add this at the end of the file. Same indentation as containers key in the file.
	    -   vi csi-provisioner-opensdsplugin.yaml
		    -   Add:
			    -   nodeSelector:
				    ntype: master
	
			-   You can add this at the end of the file. Same indentation as containers key in the file.
	    -   No need to modify csi-nodeplugin-opensdsplugin.yaml. It will be deployed to both the nodes
    
-   Compatible version information update in the pod files
    -   vi csi-attacher-opensdsplugin.yaml
	    -   image: quay.io/k8scsi/csi-attacher:v1.2.1
    -   vi csi-provisioner-opensdsplugin.yaml
	    -   image: quay.io/k8scsi/csi-provisioner:v1.6.0
    -   vi csi-nodeplugin-opensdsplugin.yaml
	    -   image: quay.io/k8scsi/csi-node-driver-registrar:v1.2.0
   -   Update rbac file for API changes due to version updates
	    -   vi csi-attacher-rbac.yaml
		    -   Add the following lines, after csinodeinfos api group
			    - apiGroups: ["storage.k8s.io"]
					resources: ["csinodes"]
					verbs: ["get", "list", "watch"]
-   Create SODA Profile file for the backend
	-   osdsctl profile create '{"name": "default_block", "description": "default policy", "storageType": "block"}'
    -   It will create a profile with the following information
    
	

    +--------------------------+--------------------------------------+
    
    	| Property | Value |
    
    	+--------------------------+--------------------------------------+
       	| Id | ffd2ff1a-ff1a-48bd-8683-0bbd519bdb48 |
       	| CreatedAt | 2020-09-28T19:05:27 |
       	| Name | default_block |
       	| Description | default policy |
       	| StorageType | block |
       	| ProvisioningProperties | { |
       	| | "dataStorage": { |
       	| | "compression": false, |
       	| | "deduplication": false |
       	| | }, |
       	| | "ioConnectivity": {} |
       	| | } |
       	| | |
       	| ReplicationProperties | { |
       	| | "dataProtection": { |
       	| | "isIsolated": false |
       	| | }, |
       	| | "replicaInfos": {} |
       	| | } |
       	| | |
       	| SnapshotProperties | { |
       	| | "schedule": {}, |
       	| | "retention": {}, |
       	| | "topology": {} |
       	| | } |
       	| | |
       	| DataProtectionProperties | { |
       	| | "dataProtection": { |
       	| | "isIsolated": false |
       	| | } |
       	| | } |
       	| | |
       	| CustomProperties | null |
       	| | |
       	+--------------------------+--------------------------------------+

-   osdsctl profile list
    
	    +--------------------------------------+---------------+----------------+-------------+
	    	| Id | Name | Description | StorageType |
	    	+--------------------------------------+---------------+----------------+-------------+
	    	| ffd2ff1a-ff1a-48bd-8683-0bbd519bdb48 | default_block | default policy | block |
	    	+--------------------------------------+---------------+----------------+-------------+

-   Please note the SODA profile id
    -   Ffd2ff1a-ff1a-48bd-8683-0bbd519bdb48
    
-   Update the profile ID in sample application pod yaml - here it is for an nginx application
	-   vi /opt/opensds-sushi-linux-amd64/csi/examples/kubernetes/block/nginx.yaml
	    -   Update the profile id as follows:
	    -   profile: ffd2ff1a-ff1a-48bd-8683-0bbd519bdb48
    
-   Ensure there is no volume or pods available
	-   kubectl get pods -o wide
	    No resources found in default namespace.
	-   osdsctl volume list
	

		    +----+------+-------------+------+--------+-----------+------------------+
		        | Id | Name | Description | Size | Status | ProfileId | AvailabilityZone |
		        +----+------+-------------+------+--------+-----------+------------------+
		        +----+------+-------------+------+--------+-----------+------------------+
-   osdsctl volume attachment list

	    +----+--------+----------+--------+----------------+
	    | Id | HostId | VolumeId | Status | AccessProtocol |
	    +----+--------+----------+--------+----------------+
	    +----+--------+----------+--------+----------------+

-   Deploy the SODA CSI plugin for block storage
    -   kubectl create -f /opt/opensds-sushi-linux-amd64/csi/deploy/kubernetes/block
Log:

		    service/csi-attacher-opensdsplugin-block created
		    
		    statefulset.apps/csi-attacher-opensdsplugin-block created
		    
		    serviceaccount/csi-attacher-block created
		    
		    clusterrole.rbac.authorization.k8s.io/external-attacher-runner-block created
		    
		    clusterrolebinding.rbac.authorization.k8s.io/csi-attacher-role-block created
		    
		    configmap/csi-configmap-opensdsplugin-block created
		    
		    daemonset.apps/csi-nodeplugin-opensdsplugin-block created
		    
		    serviceaccount/csi-nodeplugin-block created
		    
		    clusterrole.rbac.authorization.k8s.io/csi-nodeplugin-block created
		    
		    clusterrolebinding.rbac.authorization.k8s.io/csi-nodeplugin-block created
		    
		    service/csi-provisioner-opensdsplugin-block created
		    
		    statefulset.apps/csi-provisioner-opensdsplugin-block created
		    
		    serviceaccount/csi-provisioner-block created
		    
		    clusterrole.rbac.authorization.k8s.io/external-provisioner-runner-block created
		    
		    clusterrolebinding.rbac.authorization.k8s.io/csi-provisioner-role-block created
		    
		    service/csi-snapshotter-opensdsplugin-block created
		    
		    statefulset.apps/csi-snapshotter-opensdsplugin-block created
		    
		    serviceaccount/csi-snapshotter-block created
		    
		    clusterrole.rbac.authorization.k8s.io/external-snapshotter-runner-block created
		    
		    clusterrolebinding.rbac.authorization.k8s.io/csi-snapshotter-role-block created

  
  

-   Deploy the nginx test application which creates and uses block storage using SODA CSI plugin
    -   kubectl create -f /opt/opensds-sushi-linux-amd64/csi/examples/kubernetes/block/nginx.yaml
    Log:

		    storageclass.storage.k8s.io/csi-sc-opensdsplugin-block created
		    
		    persistentvolumeclaim/csi-pvc-opensdsplugin-block created
		    
		    pod/nginx-block created

  -   osdsctl volume list
    
    +--------------------------------------+------------------------------------------+-------------+------+------------------+--------+--------------------------------------+
    
    | Id | Name | Description | Size | AvailabilityZone | Status | ProfileId |
    
    +--------------------------------------+------------------------------------------+-------------+------+------------------+--------+--------------------------------------+
    
    | ad344046-b96d-441b-a51f-b11e29b1fda9 | pvc-3e9b2312-4804-4271-81e6-1cdc6fca8384 | | 1 | default | inUse | ffd2ff1a-ff1a-48bd-8683-0bbd519bdb48 |
    
    +--------------------------------------+------------------------------------------+-------------+------+------------------+--------+--------------------------------------+

-   The block volume will be created from ngnix application using SODA CSI
 
 -   osdsctl volume attachment list

    +--------------------------------------+--------------------------------------+--------------------------------------+-----------+----------------+
    
    | Id | HostId | VolumeId | Status | AccessProtocol |
    
    +--------------------------------------+--------------------------------------+--------------------------------------+-----------+----------------+
    
    | 445680ad-6a0d-4cce-aef2-898c5cb97a03 | 4c9d483b-07ae-4a1f-95f3-16c6bf3c7b23 | ad344046-b96d-441b-a51f-b11e29b1fda9 | available | iscsi |
    
    +--------------------------------------+--------------------------------------+--------------------------------------+-----------+----------------+

-   The volume attachment is successful
- Verify the volume from kubectl - it is created  
	- kubectl get pv
	    NAME CAPACITY ACCESS MODES RECLAIM POLICY STATUS CLAIM STORAGECLASS REASON AGE
	    
	    pvc-3e9b2312-4804-4271-81e6-1cdc6fca8384 1Gi RWX Delete Bound default/csi-pvc-opensdsplugin-block csi-sc-opensdsplugin-block 4m23s
	- You can see the pv   
- Verify the volume is attached and mounted - check from kubectl
	-   kubectl get pvc
  NAME STATUS VOLUME CAPACITY ACCESS MODES STORAGECLASS AGE
csi-pvc-opensdsplugin-block Bound pvc-3e9b2312-4804-4271-81e6-1cdc6fca8384 1Gi RWX csi-sc-opensdsplugin-block 4m50s
	-   The volume is attached and mounted successfully
    
- See all the pods
	-   kubectl get pods -o wide
NAME READY STATUS RESTARTS AGE IP NODE NOMINATED NODE READINESS GATES
csi-attacher-opensdsplugin-block-0 3/3 Running 0 5m45s 10.244.0.21 test-control-plane <none> <none>
csi-nodeplugin-opensdsplugin-block-97gs4 0/2 Pending 0 5m45s <none> edge-node <none> <none>
csi-nodeplugin-opensdsplugin-block-blcdk 2/2 Running 0 5m45s 172.17.0.2 test-control-plane <none> <none>
csi-provisioner-opensdsplugin-block-0 2/2 Running 0 5m45s 10.244.0.22 test-control-plane <none> <none>
csi-snapshotter-opensdsplugin-block-0 1/2 Error 5 5m45s <none> edge-node <none> <none>
nginx-block 0/1 ContainerCreating 0 5m33s <none> test-control-plane <none> <none>

-   Ignore csi-snapshotter pod error - we are not using this or not configured
-   Other SODA (opensds) CSI plugins are running ok
-   Nginx-block pod status is ContainerCreating
	-   Issue!
	   -   We observed that the node publishing has some issues from the node side.
	    -   However the poc shows clearly that, SODA and KubeEdge can work through the SODA CSI plugin and connect to heterogeneous storage backends.
    -   Please check the summary and open issues comments below for next steps
    

# Summary

## Conclusion

-   PoC successfully shows the SODA and KubeEdge integration through CSI interface
    
-   It provides the design PoC confirmation that KubeEdge can connect to heterogeneous storage backends through SODA CSI plugin (currently shown block storage. However it proves file as well)
    
-   Currently the PoC is done on the same VM with KubeEdge (Cloud/Edge) and SODA components
    
-   This gives the first step of integration of SODA at Edge.
    
-   Further analysis and integration with other edge computing platforms could be done
    
-   This can help to further refining the SODA Edge Architecture and deployments
    

## Open issues

-   SODA CSI plugin deployment to edge-node of KubeEdge
    

-   Already patch in KubeEdge in progress. That can help to fix this.
    

-   [https://github.com/kubeedge/kubeedge/pull/1834](https://github.com/kubeedge/kubeedge/pull/1834)
    

-   Analyse and fix this issue.
    

-   Node publishing issue with master node
    

-   Nginx pod under creation and not running
    
-   If the edge-node deployment issue can be fixed, then we do not need to debug this further, unless we want to deploy/do storage directly from cloud node/master in Kubeedge (need to check the use case)
    

## Next Steps

-   Issue fix as mentioned above
    
-   E2E SODA Edge with KubeEdge Cloud and Edge on different nodes and SODA installed with node where the edge node is running
    
-   Heterogeneous storage working at Edge
    
-   Architecture refinement and further analysis of SODA@Edge

