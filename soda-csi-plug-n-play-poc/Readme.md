# soda-csi-plug-n-play POC

This is a poc for soda-csi-plug-n-play to experiment the possibilities of soda-csi-provisioner being used to provision the hetrogeneous csi storage solutions.

This is an experimental feature and currently provisioning of PVC has been demonstrated in this demo.

You can follow the below steps to make a POC or watch the video of this POC over [here](https://youtu.be/ytXY_dKQCYg).

## Setup

![](./Soda-CSI-Plugin.png)

In this setup we are using two csi drivers in the same k8s env, we will use the same StorageClass with different profile ID and sod-csi-porvisioner will dynamically provision the storage's in applicable driver.

#### Step1:
Deploy IBM CSI operator and driver along with soda-csi-provisioner
```
kubectl create -f deploy/kubernetes/ibm/ibm-block-csi-operator.yaml

kubectl create -f deploy/kubernetes/ibm/csi.ibm.com_v1_ibmblockcsi_cr.yaml 

kubectl get pods
NAME                                     READY   STATUS    RESTARTS   AGE
ibm-block-csi-controller-0               5/5     Running   0          44s
ibm-block-csi-node-9fzdv                 3/3     Running   0          44s
ibm-block-csi-operator-bdfb89bdd-sh977   1/1     Running   0          79s
 
```

#### Step2:
Deploy LVM CSI operator and driver along with soda-csi-provisioner
```
kubectl create -f deploy/kubernetes/lvm/

kubectl get pods
NAME                                     READY   STATUS    RESTARTS   AGE
csi-attacher-0                           2/2     Running   0          5s
csi-lvmplugin-nkh5s                      2/2     Running   0          4s
csi-provisioner-0                        1/1     Running   0          4s
ibm-block-csi-controller-0               5/5     Running   0          74s
ibm-block-csi-node-9fzdv                 3/3     Running   0          74s
ibm-block-csi-operator-bdfb89bdd-sh977   1/1     Running   0          109s

```

#### Step3:
Deploy StorageClass and PVC with profile as 'block.csi.ibm.com'
```go
kubectl create -f deploy/kubernetes/demo
```

This will create the StorageClass but the PVC will fail as the IBM backend is not available and the soda-csi-provisioner will send an error.

`IBM soda-csi-provisioner logs`

```go
kubectl logs ibm-block-csi-controller-0 csi-provisioner 
I0920 21:04:24.859595       1 feature_gate.go:243] feature gates: &{map[]}
I0920 21:04:24.859676       1 csi-provisioner.go:107] Version: v1.6.0-0-g321fa5c1c-dirty
I0920 21:04:24.859697       1 csi-provisioner.go:121] Building kube configs for running in cluster...
I0920 21:04:24.867210       1 connection.go:153] Connecting to unix:///var/lib/csi/sockets/pluginproxy/csi.sock
I0920 21:04:26.952475       1 common.go:111] Probing CSI driver for readiness
I0920 21:04:26.952501       1 connection.go:182] GRPC call: /csi.v1.Identity/Probe
I0920 21:04:26.952507       1 connection.go:183] GRPC request: {}
I0920 21:04:27.032893       1 connection.go:185] GRPC response: {}
I0920 21:04:27.035408       1 connection.go:186] GRPC error: <nil>
I0920 21:04:27.035420       1 csi-provisioner.go:165] Detected CSI driver soda-csi
W0920 21:04:27.035436       1 metrics.go:142] metrics endpoint will not be started because `metrics-address` was not specified.
I0920 21:04:27.035449       1 connection.go:182] GRPC call: /csi.v1.Identity/GetPluginCapabilities
I0920 21:04:27.035645       1 connection.go:183] GRPC request: {}
I0920 21:04:27.037894       1 connection.go:185] GRPC response: {"capabilities":[{"Type":{"Service":{"type":1}}}]}
I0920 21:04:27.038841       1 connection.go:186] GRPC error: <nil>
I0920 21:04:27.038854       1 connection.go:182] GRPC call: /csi.v1.Controller/ControllerGetCapabilities
I0920 21:04:27.038858       1 connection.go:183] GRPC request: {}
I0920 21:04:27.044049       1 connection.go:185] GRPC response: {"capabilities":[{"Type":{"Rpc":{"type":1}}},{"Type":{"Rpc":{"type":5}}},{"Type":{"Rpc":{"type":2}}}]}
I0920 21:04:27.046370       1 connection.go:186] GRPC error: <nil>
I0920 21:04:27.046725       1 controller.go:709] Using saving PVs to API server in background
I0920 21:04:27.046862       1 reflector.go:153] Starting reflector *v1.StorageClass (1h0m0s) from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.046868       1 reflector.go:153] Starting reflector *v1.PersistentVolumeClaim (15m0s) from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.046888       1 reflector.go:188] Listing and watching *v1.PersistentVolumeClaim from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.046877       1 reflector.go:188] Listing and watching *v1.StorageClass from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.146948       1 shared_informer.go:227] caches populated
I0920 21:04:27.146973       1 shared_informer.go:227] caches populated
I0920 21:04:27.146986       1 controller.go:799] Starting provisioner controller soda-csi_ibm-block-csi-controller-0_5892083e-79a4-4bbf-a622-eba90a09a439!
I0920 21:04:27.147019       1 clone_controller.go:58] Starting CloningProtection controller
I0920 21:04:27.147045       1 clone_controller.go:74] Started CloningProtection controller
I0920 21:04:27.147098       1 volume_store.go:97] Starting save volume queue
I0920 21:04:27.147140       1 reflector.go:153] Starting reflector *v1.StorageClass (15m0s) from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.147160       1 reflector.go:188] Listing and watching *v1.StorageClass from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.147206       1 reflector.go:153] Starting reflector *v1.PersistentVolume (15m0s) from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.147219       1 reflector.go:188] Listing and watching *v1.PersistentVolume from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.147548       1 reflector.go:153] Starting reflector *v1.PersistentVolumeClaim (15m0s) from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.147563       1 reflector.go:188] Listing and watching *v1.PersistentVolumeClaim from pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108
I0920 21:04:27.249045       1 shared_informer.go:227] caches populated
I0920 21:04:27.249398       1 controller.go:848] Started provisioner controller soda-csi_ibm-block-csi-controller-0_5892083e-79a4-4bbf-a622-eba90a09a439!
I0920 21:06:35.168794       1 controller.go:1284] provision "default/demo-pvc-file-system" class "soda-high-io": started
I0920 21:06:35.171967       1 connection.go:182] GRPC call: /csi.v1.Identity/GetPluginInfo
I0920 21:06:35.171980       1 connection.go:183] GRPC request: {}
I0920 21:06:35.173139       1 event.go:281] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"demo-pvc-file-system", UID:"702cc677-f023-4587-a50b-36d27a854975", APIVersion:"v1", ResourceVersion:"17886520", FieldPath:""}): type: 'Normal' reason: 'Provisioning' External provisioner is provisioning volume for claim "default/demo-pvc-file-system"
I0920 21:06:35.174366       1 connection.go:185] GRPC response: {"name":"block.csi.ibm.com","vendor_version":"1.3.0"}
I0920 21:06:35.174959       1 connection.go:186] GRPC error: <nil>
I0920 21:06:35.174969       1 controller.go:437] The Backend Driver Name is : block.csi.ibm.com 
I0920 21:06:35.174975       1 controller.go:438] The provisioner.DriverName  is : soda-csi 
I0920 21:06:35.174983       1 controller.go:468] The parameters in the StorageClass are  : profile ===== block.csi.ibm.com
I0920 21:06:35.175008       1 controller.go:613] CreateVolumeRequest {Name:pvc-702cc677-f023-4587-a50b-36d27a854975 CapacityRange:required_bytes:1073741824  VolumeCapabilities:[mount:<fs_type:"ext4" > access_mode:<mode:SINGLE_NODE_WRITER > ] Parameters:map[profile:block.csi.ibm.com] Secrets:map[] VolumeContentSource:<nil> AccessibilityRequirements:<nil> XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}


```

Whereas in LVM soda-csi-provisioner the PVC creation will skip as the profile does not matches to LVM.
`LVM soda-csi-provisioner logs`


```go
kubectl logs csi-provisioner-0 csi-provisioner 
W0920 21:05:34.172439       1 deprecatedflags.go:53] Warning: option provisioner="csi-lvmplugin" is deprecated and has no effect
I0920 21:05:34.172498       1 feature_gate.go:243] feature gates: &{map[Topology:true]}
I0920 21:05:34.172517       1 csi-provisioner.go:107] Version: v1.6.0-0-g321fa5c1c-dirty
I0920 21:05:34.172532       1 csi-provisioner.go:121] Building kube configs for running in cluster...
I0920 21:05:34.173769       1 round_trippers.go:423] curl -k -v -XGET  -H "Accept: application/json, */*" -H "User-Agent: csi-provisioner/v0.0.0 (linux/amd64) kubernetes/$Format" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ink4YVNPdXNndzdsaGlZRlJ1YWRTXzNmNTdlTTJJakdMUUVoVm9NYUhRaGMifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImNzaS1wcm92aXNpb25lci10b2tlbi1rbXNuNSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJjc2ktcHJvdmlzaW9uZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI3ZWJiY2FiMS01OTExLTQ4OGYtYmU3My04MzIzMzU3NTg5MWYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDpjc2ktcHJvdmlzaW9uZXIifQ.Pz4x9G0sjiDzNce_CE1JJ-Zw8ulPk38RPwKpt7Qs0Z3G2JDLNYvCnfKwDKV7KWfEBz6BbX15Zpz4hiuAufznolRpQ2Eh6eFYKWykZEiaLOZh23sjbME4zzonJpTfXBcrRxLbOURpuZ7lxX3p4PNDJ3tEfrMP9COejIZALSNYiiSGqmDR2JjAcR_y1PSJHcbfZUxve1vMm-DrcJrn6cBh_etI-kX80vqIIQNb77hn_BIs1tivhDWq3RPapWHrmh7l5n_GqAlklGFCsARp0QPguQRNfyHBABmfiWYcgTkem98cSSZwmvqZkjM6ZkYjG4l3ALi7aN5BYACxOHgKuQii7w" 'https://10.96.0.1:443/version?timeout=32s'
I0920 21:05:34.179965       1 round_trippers.go:443] GET https://10.96.0.1:443/version?timeout=32s 200 OK in 6 milliseconds
I0920 21:05:34.179993       1 round_trippers.go:449] Response Headers:
I0920 21:05:34.180001       1 round_trippers.go:452]     Cache-Control: no-cache, private
I0920 21:05:34.180006       1 round_trippers.go:452]     Content-Type: application/json
I0920 21:05:34.180010       1 round_trippers.go:452]     Content-Length: 263
I0920 21:05:34.180014       1 round_trippers.go:452]     Date: Sun, 20 Sep 2020 21:05:34 GMT
I0920 21:05:34.180044       1 request.go:1017] Response Body: {
  "major": "1",
  "minor": "16",
  "gitVersion": "v1.16.9",
  "gitCommit": "a17149e1a189050796ced469dbd78d380f2ed5ef",
  "gitTreeState": "clean",
  "buildDate": "2020-04-16T11:36:15Z",
  "goVersion": "go1.13.9",
  "compiler": "gc",
  "platform": "linux/amd64"
}
I0920 21:05:34.180221       1 connection.go:153] Connecting to unix:///var/lib/kubelet/plugins/csi-lvmplugin/csi.sock
I0920 21:05:34.181061       1 common.go:111] Probing CSI driver for readiness
I0920 21:05:34.181078       1 connection.go:182] GRPC call: /csi.v1.Identity/Probe
I0920 21:05:34.181083       1 connection.go:183] GRPC request: {}
I0920 21:05:34.187264       1 connection.go:185] GRPC response: {}
I0920 21:05:34.187737       1 connection.go:186] GRPC error: <nil>
I0920 21:05:34.187748       1 csi-provisioner.go:165] Detected CSI driver soda-csi
W0920 21:05:34.187754       1 metrics.go:142] metrics endpoint will not be started because `metrics-address` was not specified.
I0920 21:05:34.187763       1 connection.go:182] GRPC call: /csi.v1.Identity/GetPluginCapabilities
I0920 21:05:34.187767       1 connection.go:183] GRPC request: {}
I0920 21:05:34.192064       1 connection.go:185] GRPC response: {"capabilities":[{"Type":{"Service":{"type":1}}},{"Type":{"Service":{"type":2}}}]}
I0920 21:05:34.194037       1 connection.go:186] GRPC error: <nil>
I0920 21:05:34.194050       1 connection.go:182] GRPC call: /csi.v1.Controller/ControllerGetCapabilities
I0920 21:05:34.194055       1 connection.go:183] GRPC request: {}
I0920 21:05:34.195653       1 connection.go:185] GRPC response: {"capabilities":[{"Type":{"Rpc":{"type":1}}}]}
.
.
.
.
.

1 connection.go:182] GRPC call: /csi.v1.Identity/GetPluginInfo
I0920 21:06:35.183464       1 connection.go:183] GRPC request: {}
I0920 21:06:35.184738       1 event.go:281] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"demo-pvc-file-system", UID:"702cc677-f023-4587-a50b-36d27a854975", APIVersion:"v1", ResourceVersion:"17886520", FieldPath:""}): type: 'Normal' reason: 'Provisioning' External provisioner is provisioning volume for claim "default/demo-pvc-file-system"
I0920 21:06:35.185107       1 request.go:1017] Request Body: {"kind":"Event","apiVersion":"v1","metadata":{"name":"demo-pvc-file-system.163699f5039cd664","namespace":"default","creationTimestamp":null},"involvedObject":{"kind":"PersistentVolumeClaim","namespace":"default","name":"demo-pvc-file-system","uid":"702cc677-f023-4587-a50b-36d27a854975","apiVersion":"v1","resourceVersion":"17886520"},"reason":"Provisioning","message":"External provisioner is provisioning volume for claim \"default/demo-pvc-file-system\"","source":{"component":"soda-csi_csi-provisioner-0_ce2a73c8-1c0e-4e9c-9ceb-56ed7547fc37"},"firstTimestamp":"2020-09-20T21:06:35Z","lastTimestamp":"2020-09-20T21:06:35Z","count":1,"type":"Normal","eventTime":null,"reportingComponent":"","reportingInstance":""}
I0920 21:06:35.185733       1 round_trippers.go:423] curl -k -v -XPOST  -H "Accept: application/json, */*" -H "Content-Type: application/json" -H "User-Agent: csi-provisioner/v0.0.0 (linux/amd64) kubernetes/$Format" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ink4YVNPdXNndzdsaGlZRlJ1YWRTXzNmNTdlTTJJakdMUUVoVm9NYUhRaGMifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImNzaS1wcm92aXNpb25lci10b2tlbi1rbXNuNSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJjc2ktcHJvdmlzaW9uZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI3ZWJiY2FiMS01OTExLTQ4OGYtYmU3My04MzIzMzU3NTg5MWYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDpjc2ktcHJvdmlzaW9uZXIifQ.Pz4x9G0sjiDzNce_CE1JJ-Zw8ulPk38RPwKpt7Qs0Z3G2JDLNYvCnfKwDKV7KWfEBz6BbX15Zpz4hiuAufznolRpQ2Eh6eFYKWykZEiaLOZh23sjbME4zzonJpTfXBcrRxLbOURpuZ7lxX3p4PNDJ3tEfrMP9COejIZALSNYiiSGqmDR2JjAcR_y1PSJHcbfZUxve1vMm-DrcJrn6cBh_etI-kX80vqIIQNb77hn_BIs1tivhDWq3RPapWHrmh7l5n_GqAlklGFCsARp0QPguQRNfyHBABmfiWYcgTkem98cSSZwmvqZkjM6ZkYjG4l3ALi7aN5BYACxOHgKuQii7w" 'https://10.96.0.1:443/api/v1/namespaces/default/events'
I0920 21:06:35.187288       1 connection.go:185] GRPC response: {"name":"csi-lvmplugin"}
I0920 21:06:35.188694       1 connection.go:186] GRPC error: <nil>
I0920 21:06:35.188705       1 controller.go:437] The Backend Driver Name is : csi-lvmplugin 
I0920 21:06:35.188709       1 controller.go:438] The provisioner.DriverName  is : soda-csi 
I0920 21:06:35.188715       1 controller.go:468] The parameters in the StorageClass are  : profile ===== block.csi.ibm.com
I0920 21:06:35.188735       1 controller.go:1358] provision "default/demo-pvc-file-system" class "soda-high-io": volume provision ignored: ignored because PVC doesnot match the current driver name : soda-csi with expected block.csi.ibm.com
I0920 21:06:35.188749       1 controller.go:1047] Stop provisioning, removing PVC 702cc677-f023-4587-a50b-36d27a854975 from claims in progress

```


#### Step4:
Deploy StorageClass and PVC with profile as 'csi-lvmplugin'
```
kubectl delete -f deploy/kubernetes/demo/
```
```
vi deploy/kubernetes/demo/storageClass.yaml


Replace   profile: block.csi.ibm.com 
                    with
          profile: csi-lvmplugin
```

```go
kubectl create -f deploy/kubernetes/demo
```

IBM soda-csi-provisioner the PVC creation will skip as the profile does not matches to IBM

`IBM soda-csi-provisioner logs`

```go
kubectl logs ibm-block-csi-controller-0 csi-provisioner
.
.
.
.
I0920 21:08:40.742155       1 controller.go:1284] provision "default/demo-pvc-file-system" class "soda-high-io": started
I0920 21:08:40.744071       1 connection.go:182] GRPC call: /csi.v1.Identity/GetPluginInfo
I0920 21:08:40.744089       1 connection.go:183] GRPC request: {}
I0920 21:08:40.744630       1 event.go:281] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"demo-pvc-file-system", UID:"702cc677-f023-4587-a50b-36d27a854975", APIVersion:"v1", ResourceVersion:"17886520", FieldPath:""}): type: 'Normal' reason: 'Provisioning' External provisioner is provisioning volume for claim "default/demo-pvc-file-system"
I0920 21:08:40.745490       1 connection.go:185] GRPC response: {"name":"block.csi.ibm.com","vendor_version":"1.3.0"}
I0920 21:08:40.745920       1 connection.go:186] GRPC error: <nil>
I0920 21:08:40.745930       1 controller.go:437] The Backend Driver Name is : block.csi.ibm.com 
I0920 21:08:40.745934       1 controller.go:438] The provisioner.DriverName  is : soda-csi 
I0920 21:08:40.745940       1 controller.go:468] The parameters in the StorageClass are  : profile ===== block.csi.ibm.com
I0920 21:08:40.745960       1 controller.go:613] CreateVolumeRequest {Name:pvc-702cc677-f023-4587-a50b-36d27a854975 CapacityRange:required_bytes:1073741824  VolumeCapabilities:[mount:<fs_type:"ext4" > access_mode:<mode:SINGLE_NODE_WRITER > ] Parameters:map[profile:block.csi.ibm.com] Secrets:map[] VolumeContentSource:<nil> AccessibilityRequirements:<nil> XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}
I0920 21:08:40.746020       1 connection.go:182] GRPC call: /csi.v1.Controller/CreateVolume
I0920 21:08:40.746027       1 connection.go:183] GRPC request: {"capacity_range":{"required_bytes":1073741824},"name":"pvc-702cc677-f023-4587-a50b-36d27a854975","parameters":{"profile":"block.csi.ibm.com"},"volume_capabilities":[{"AccessType":{"Mount":{"fs_type":"ext4"}},"access_mode":{"mode":1}}]}
I0920 21:08:40.835658       1 connection.go:185] GRPC response: {}
I0920 21:08:40.836071       1 connection.go:186] GRPC error: rpc error: code = InvalidArgument desc = Validation error has occurred : pool parameter is missing.
I0920 21:08:40.836097       1 controller.go:685] CreateVolume failed, supports topology = false, node selected false => may reschedule = false => state = Finished: rpc error: code = InvalidArgument desc = Validation error has occurred : pool parameter is missing.
I0920 21:08:40.836143       1 controller.go:1051] Final error received, removing PVC 702cc677-f023-4587-a50b-36d27a854975 from claims in progress
W0920 21:08:40.836153       1 controller.go:916] Retrying syncing claim "702cc677-f023-4587-a50b-36d27a854975", failure 7
E0920 21:08:40.836249       1 controller.go:939] error syncing claim "702cc677-f023-4587-a50b-36d27a854975": failed to provision volume with StorageClass "soda-high-io": rpc error: code = InvalidArgument desc = Validation error has occurred : pool parameter is missing.
I0920 21:08:40.836274       1 event.go:281] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"demo-pvc-file-system", UID:"702cc677-f023-4587-a50b-36d27a854975", APIVersion:"v1", ResourceVersion:"17886520", FieldPath:""}): type: 'Warning' reason: 'ProvisioningFailed' failed to provision volume with StorageClass "soda-high-io": rpc error: code = InvalidArgument desc = Validation error has occurred : pool parameter is missing.
I0920 21:10:48.836536       1 controller.go:1284] provision "default/demo-pvc-file-system" class "soda-high-io": started
I0920 21:10:48.932969       1 connection.go:182] GRPC call: /csi.v1.Identity/GetPluginInfo
I0920 21:10:48.932988       1 connection.go:183] GRPC request: {}
I0920 21:10:48.933462       1 event.go:281] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"demo-pvc-file-system", UID:"702cc677-f023-4587-a50b-36d27a854975", APIVersion:"v1", ResourceVersion:"17886520", FieldPath:""}): type: 'Normal' reason: 'Provisioning' External provisioner is provisioning volume for claim "default/demo-pvc-file-system"
I0920 21:10:48.934311       1 connection.go:185] GRPC response: {"name":"block.csi.ibm.com","vendor_version":"1.3.0"}
I0920 21:10:48.934745       1 connection.go:186] GRPC error: <nil>
I0920 21:10:48.934753       1 controller.go:437] The Backend Driver Name is : block.csi.ibm.com 
I0920 21:10:48.934758       1 controller.go:438] The provisioner.DriverName  is : soda-csi 
I0920 21:10:48.934763       1 controller.go:468] The parameters in the StorageClass are  : profile ===== block.csi.ibm.com
I0920 21:10:48.934780       1 controller.go:613] CreateVolumeRequest {Name:pvc-702cc677-f023-4587-a50b-36d27a854975 CapacityRange:required_bytes:1073741824  VolumeCapabilities:[mount:<fs_type:"ext4" > access_mode:<mode:SINGLE_NODE_WRITER > ] Parameters:map[profile:block.csi.ibm.com] Secrets:map[] VolumeContentSource:<nil> AccessibilityRequirements:<nil> XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}
I0920 21:10:48.934839       1 connection.go:182] GRPC call: /csi.v1.Controller/CreateVolume
I0920 21:10:48.934846       1 connection.go:183] GRPC request: {"capacity_range":{"required_bytes":1073741824},"name":"pvc-702cc677-f023-4587-a50b-36d27a854975","parameters":{"profile":"block.csi.ibm.com"},"volume_capabilities":[{"AccessType":{"Mount":{"fs_type":"ext4"}},"access_mode":{"mode":1}}]}
I0920 21:10:48.939231       1 connection.go:185] GRPC response: {}
I0920 21:10:48.939632       1 connection.go:186] GRPC error: rpc error: code = InvalidArgument desc = Validation error has occurred : pool parameter is missing.
I0920 21:10:48.939658       1 controller.go:685] CreateVolume failed, supports topology = false, node selected false => may reschedule = false => state = Finished: rpc error: code = InvalidArgument desc = Validation error has occurred : pool parameter is missing.
I0920 21:10:48.939688       1 controller.go:1051] Final error received, removing PVC 702cc677-f023-4587-a50b-36d27a854975 from claims in progress
W0920 21:10:48.939699       1 controller.go:916] Retrying syncing claim "702cc677-f023-4587-a50b-36d27a854975", failure 8
E0920 21:10:48.939719       1 controller.go:939] error syncing claim "702cc677-f023-4587-a50b-36d27a854975": failed to provision volume with StorageClass "soda-high-io": rpc error: code = InvalidArgument desc = Validation error has occurred : pool parameter is missing.
I0920 21:10:48.940015       1 event.go:281] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"demo-pvc-file-system", UID:"702cc677-f023-4587-a50b-36d27a854975", APIVersion:"v1", ResourceVersion:"17886520", FieldPath:""}): type: 'Warning' reason: 'ProvisioningFailed' failed to provision volume with StorageClass "soda-high-io": rpc error: code = InvalidArgument desc = Validation error has occurred : pool parameter is missing.
I0920 21:11:34.149840       1 reflector.go:432] pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108: Watch close - *v1.StorageClass total 3 items received
I0920 21:11:38.149693       1 reflector.go:432] pkg/mod/k8s.io/client-go@v0.17.0/tools/cache/reflector.go:108: Watch close - *v1.PersistentVolumeClaim total 5 items received
I0920 21:12:20.232798       1 controller.go:1284] provision "default/demo-pvc-file-system" class "soda-high-io": started
I0920 21:12:20.235428       1 controller.go:1293] provision "default/demo-pvc-file-system" class "soda-high-io": persistentvolume "pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd" already exists, skipping
I0920 21:12:20.235556       1 controller.go:1047] Stop provisioning, removing PVC f3912c2f-16ca-46b1-b44b-22d5372d7bdd from claims in progress


```


Whereas in LVM soda-csi-provisioner the PVC creation will proceed as the profile matches to LVM.
`LVM soda-csi-provisioner logs`

```go
kubectl logs csi-provisioner-0 csi-provisioner
.
.
.
.
I0920 21:12:20.177569       1 connection.go:183] GRPC request: {}
I0920 21:12:20.180253       1 round_trippers.go:443] POST https://10.96.0.1:443/api/v1/namespaces/default/events 201 Created in 2 milliseconds
I0920 21:12:20.180275       1 round_trippers.go:449] Response Headers:
I0920 21:12:20.180281       1 round_trippers.go:452]     Content-Length: 900
I0920 21:12:20.180285       1 round_trippers.go:452]     Date: Sun, 20 Sep 2020 21:12:20 GMT
I0920 21:12:20.180290       1 round_trippers.go:452]     Cache-Control: no-cache, private
I0920 21:12:20.180294       1 round_trippers.go:452]     Content-Type: application/json
I0920 21:12:20.180335       1 request.go:1017] Response Body: {"kind":"Event","apiVersion":"v1","metadata":{"name":"demo-pvc-file-system.16369a4556dcc2c0","namespace":"default","selfLink":"/api/v1/namespaces/default/events/demo-pvc-file-system.16369a4556dcc2c0","uid":"f754f898-cafd-490e-a9d2-aa0c9abc2916","resourceVersion":"17886992","creationTimestamp":"2020-09-20T21:12:20Z"},"involvedObject":{"kind":"PersistentVolumeClaim","namespace":"default","name":"demo-pvc-file-system","uid":"f3912c2f-16ca-46b1-b44b-22d5372d7bdd","apiVersion":"v1","resourceVersion":"17886990"},"reason":"Provisioning","message":"External provisioner is provisioning volume for claim \"default/demo-pvc-file-system\"","source":{"component":"soda-csi_csi-provisioner-0_ce2a73c8-1c0e-4e9c-9ceb-56ed7547fc37"},"firstTimestamp":"2020-09-20T21:12:20Z","lastTimestamp":"2020-09-20T21:12:20Z","count":1,"type":"Normal","eventTime":null,"reportingComponent":"","reportingInstance":""}
I0920 21:12:20.180088       1 connection.go:185] GRPC response: {"name":"csi-lvmplugin"}
I0920 21:12:20.180885       1 connection.go:186] GRPC error: <nil>
I0920 21:12:20.180894       1 controller.go:437] The Backend Driver Name is : csi-lvmplugin 
I0920 21:12:20.180898       1 controller.go:438] The provisioner.DriverName  is : soda-csi 
I0920 21:12:20.180911       1 controller.go:468] The parameters in the StorageClass are  : profile ===== csi-lvmplugin
W0920 21:12:20.180931       1 topology.go:343] No topology keys found on any node
I0920 21:12:20.180943       1 controller.go:613] CreateVolumeRequest {Name:pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd CapacityRange:required_bytes:1073741824  VolumeCapabilities:[mount:<fs_type:"ext4" > access_mode:<mode:SINGLE_NODE_WRITER > ] Parameters:map[profile:csi-lvmplugin] Secrets:map[] VolumeContentSource:<nil> AccessibilityRequirements:<nil> XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}
I0920 21:12:20.181065       1 connection.go:182] GRPC call: /csi.v1.Controller/CreateVolume
I0920 21:12:20.181075       1 connection.go:183] GRPC request: {"capacity_range":{"required_bytes":1073741824},"name":"pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd","parameters":{"profile":"csi-lvmplugin"},"volume_capabilities":[{"AccessType":{"Mount":{"fs_type":"ext4"}},"access_mode":{"mode":1}}]}
I0920 21:12:20.186800       1 connection.go:185] GRPC response: {"volume":{"capacity_bytes":1073741824,"volume_id":"pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd"}}
I0920 21:12:20.187638       1 connection.go:186] GRPC error: <nil>
I0920 21:12:20.187659       1 controller.go:695] create volume rep: {CapacityBytes:1073741824 VolumeId:pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd VolumeContext:map[] ContentSource:<nil> AccessibleTopology:[] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}
I0920 21:12:20.187698       1 controller.go:777] successfully created PV pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd for PVC demo-pvc-file-system and csi volume name pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd
I0920 21:12:20.187712       1 controller.go:793] successfully created PV {GCEPersistentDisk:nil AWSElasticBlockStore:nil HostPath:nil Glusterfs:nil NFS:nil RBD:nil ISCSI:nil Cinder:nil CephFS:nil FC:nil Flocker:nil FlexVolume:nil AzureFile:nil VsphereVolume:nil Quobyte:nil AzureDisk:nil PhotonPersistentDisk:nil PortworxVolume:nil ScaleIO:nil Local:nil StorageOS:nil CSI:&CSIPersistentVolumeSource{Driver:soda-csi,VolumeHandle:pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd,ReadOnly:false,FSType:ext4,VolumeAttributes:map[string]string{storage.kubernetes.io/csiProvisionerIdentity: 1600635934196-8081-soda-csi,},ControllerPublishSecretRef:nil,NodeStageSecretRef:nil,NodePublishSecretRef:nil,ControllerExpandSecretRef:nil,}}
I0920 21:12:20.187848       1 controller.go:1392] provision "default/demo-pvc-file-system" class "soda-high-io": volume "pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd" provisioned
I0920 21:12:20.187880       1 controller.go:1409] provision "default/demo-pvc-file-system" class "soda-high-io": succeeded
I0920 21:12:20.187893       1 volume_store.go:154] Saving volume pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd
I0920 21:12:20.188670       1 request.go:1017] Request Body: {"kind":"PersistentVolume","apiVersion":"v1","metadata":{"name":"pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd","creationTimestamp":null,"annotations":{"pv.kubernetes.io/provisioned-by":"soda-csi"}},"spec":{"capacity":{"storage":"1Gi"},"csi":{"driver":"soda-csi","volumeHandle":"pvc-f3912c2f-16ca-46b1-b44b-22d5372d7bdd","fsType":"ext4","volumeAttributes":{"storage.kubernetes.io/csiProvisionerIdentity":"1600635934196-8081-soda-csi"}},"accessModes":["ReadWriteOnce"],"claimRef":{"kind":"PersistentVolumeClaim","namespace":"default","name":"demo-pvc-file-system","uid":"f3912c2f-16ca-46b1-b44b-22d5372d7bdd","apiVersion":"v1","resourceVersion":"17886990"},"persistentVolumeReclaimPolicy":"Delete","storageClassName":"soda-high-io","volumeMode":"Filesystem"},"status":{}}


```

***Note***: 
The Deployment scripts for plugins and operators has been take from https://github.com/IBM/ibm-block-csi-operator and https://github.com/wavezhang/k8s-csi-lvm . 


#### Experimenting with other CSI drivers
##### Ceph RBD CSI driver
###### Step1:
Deploy Ceph RBD CSI driver along with soda-csi-provisioner
```
kubectl create -f deploy/kubernetes/cephcsi/rbd

kubectl get pods
NAME                                         READY   STATUS    RESTARTS   AGE
csi-rbdplugin-6pw9z                          3/3     Running   0          7s
csi-rbdplugin-provisioner-6b8b9d99fd-x4wn6   7/7     Running   0          7s
 
```
###### Step2:
Deploy Secret with UserID and UserKey of the CEPH user
```go
kubectl create -f deploy/kubernetes/ceph/secret.yaml
```
```go
kubectl get secret csi-rbd-secret -o yaml
apiVersion: v1
data:
userID: a3ViZQ==
userKey: QVFEalpZRmdPdy9GRnhBQVNXcFJKeitIMHlJM0xVNHg2Q3QyekE9PQ==
kind: Secret
metadata:
creationTimestamp: "2021-04-27T04:14:33Z"
managedFields:
- apiVersion: v1
fieldsType: FieldsV1
fieldsV1:
f:data:
.: {}
f:userID: {}
f:userKey: {}
f:type: {}
manager: kubectl-create
operation: Update
time: "2021-04-27T04:14:33Z"
name: csi-rbd-secret
namespace: default
resourceVersion: "229355"
uid: 6db441a0-53ed-480c-806d-2c2982301022
type: Opaque

```

###### Step3:
Deploy StorageClass, PVC and POD with profile ID having "driver": "rbd.csi.ceph.com" in CustomProperties of a Profile.
```go
kubectl create -f deploy/kubernetes/ceph
```
```go
kubectl get sc csi-rbd-ceph-sc -o yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
creationTimestamp: "2021-04-27T04:15:54Z"
managedFields:
- apiVersion: storage.k8s.io/v1
fieldsType: FieldsV1
fieldsV1:
f:mountOptions: {}
f:parameters:
.: {}
f:attachMode: {}
f:clusterID: {}
f:csi.storage.k8s.io/controller-expand-secret-name: {}
f:csi.storage.k8s.io/controller-expand-secret-namespace: {}
f:csi.storage.k8s.io/fstype: {}
f:csi.storage.k8s.io/node-stage-secret-name: {}
f:csi.storage.k8s.io/node-stage-secret-namespace: {}
f:csi.storage.k8s.io/provisioner-secret-name: {}
f:csi.storage.k8s.io/provisioner-secret-namespace: {}
f:imageFeatures: {}
f:imageFormat: {}
f:pool: {}
f:profile: {}
f:provisioner: {}
f:reclaimPolicy: {}
f:volumeBindingMode: {}
manager: kubectl-create
operation: Update
time: "2021-04-27T04:15:54Z"
name: csi-rbd-ceph-sc
resourceVersion: "229530"
uid: 818d0dfc-fe46-43f2-aa7a-255d8f5d93d6
mountOptions:
- discard
parameters:
attachMode: rw
clusterID: 4ac5251b-114f-4044-9bec-2d27fadad502
csi.storage.k8s.io/controller-expand-secret-name: csi-rbd-secret
csi.storage.k8s.io/controller-expand-secret-namespace: default
csi.storage.k8s.io/fstype: ext4
csi.storage.k8s.io/node-stage-secret-name: csi-rbd-secret
csi.storage.k8s.io/node-stage-secret-namespace: default
csi.storage.k8s.io/provisioner-secret-name: csi-rbd-secret
csi.storage.k8s.io/provisioner-secret-namespace: default
imageFeatures: layering
imageFormat: "2"
pool: osdsrbd
profile: f63a7fcc-6e83-49cc-81c1-8627fc8e0f52
provisioner: soda-csi
reclaimPolicy: Delete
volumeBindingMode: Immediate


```

StorageClass is created as shown above and the PVC creation will proceed as the profile matches to CEPH.
```go
Ceph RBD soda-csi-provisioner logs
```
```go
kubectl logs csi-rbdplugin-provisioner-68b5cdf677-rlhqq csi-provisioner
.
.
.
.
W0427 04:14:53.391436       1 feature_gate.go:235] Setting GA feature gate Topology=false. It will be removed in a future release.
I0427 04:14:53.391654       1 feature_gate.go:243] feature gates: &{map[Topology:false]}
I0427 04:14:53.391844       1 csi-provisioner.go:132] Version:
I0427 04:14:53.393313       1 csi-provisioner.go:155] Building kube configs for running in cluster...
I0427 04:14:54.573753       1 connection.go:153] Connecting to unix:///csi/csi-provisioner.sock
I0427 04:14:56.722717       1 common.go:111] Probing CSI driver for readiness
I0427 04:14:56.722774       1 connection.go:182] GRPC call: /csi.v1.Identity/Probe
I0427 04:14:56.722789       1 connection.go:183] GRPC request: {}
I0427 04:14:56.975155       1 connection.go:185] GRPC response: {}
I0427 04:14:56.975305       1 connection.go:186] GRPC error: <nil>
I0427 04:14:56.975333       1 csi-provisioner.go:199] Detected CSI driver soda-csi
I0427 04:14:56.975356       1 connection.go:182] GRPC call: /csi.v1.Identity/GetPluginCapabilities
I0427 04:14:56.975368       1 connection.go:183] GRPC request: {}
I0427 04:14:57.049960       1 connection.go:185] GRPC response: {"capabilities":[{"Type":{"Service":{"type":1}}},{"Type":{"VolumeExpansion":{"type":1}}},{"Type":{"Service":{"type":2}}}]}
I0427 04:14:57.050220       1 connection.go:186] GRPC error: <nil>
I0427 04:14:57.050247       1 connection.go:182] GRPC call: /csi.v1.Controller/ControllerGetCapabilities
I0427 04:14:57.050260       1 connection.go:183] GRPC request: {}
I0427 04:14:57.051604       1 connection.go:185] GRPC response: {"capabilities":[{"Type":{"Rpc":{"type":1}}},{"Type":{"Rpc":{"type":5}}},{"Type":{"Rpc":{"type":7}}},{"Type":{"Rpc":{"type":9}}}]}
I0427 04:14:57.051754       1 connection.go:186] GRPC error: <nil>
I0427 04:14:57.231707       1 csi-provisioner.go:241] CSI driver does not support PUBLISH_UNPUBLISH_VOLUME, not watching VolumeAttachments
I0427 04:14:57.333851       1 controller.go:756] Using saving PVs to API server in background
I0427 04:14:57.392443       1 reflector.go:219] Starting reflector *v1.PersistentVolumeClaim (15m0s) from k8s.io/client-go/informers/factory.go:134
I0427 04:14:57.392512       1 reflector.go:255] Listing and watching *v1.PersistentVolumeClaim from k8s.io/client-go/informers/factory.go:134
I0427 04:14:57.392539       1 reflector.go:219] Starting reflector *v1.StorageClass (1h0m0s) from k8s.io/client-go/informers/factory.go:134
I0427 04:14:57.393093       1 reflector.go:255] Listing and watching *v1.StorageClass from k8s.io/client-go/informers/factory.go:134
I0427 04:14:57.634197       1 shared_informer.go:270] caches populated
I0427 04:14:57.634232       1 shared_informer.go:270] caches populated
I0427 04:14:57.634266       1 controller.go:835] Starting provisioner controller soda-csi_csi-rbdplugin-provisioner-68b5cdf677-rlhqq_79ddf7a3-6411-4213-a9ba-a19b6ca12d44!
I0427 04:14:57.634379       1 clone_controller.go:66] Starting CloningProtection controller
I0427 04:14:57.634442       1 clone_controller.go:84] Started CloningProtection controller
I0427 04:14:57.634730       1 volume_store.go:97] Starting save volume queue
I0427 04:14:57.634790       1 reflector.go:219] Starting reflector *v1.PersistentVolumeClaim (15m0s) from github.com/kubernetes-csi/external-provisioner/pkg/controller/clone_controller.go:82
I0427 04:14:57.635099       1 reflector.go:255] Listing and watching *v1.PersistentVolumeClaim from github.com/kubernetes-csi/external-provisioner/pkg/controller/clone_controller.go:82
I0427 04:14:57.635049       1 reflector.go:219] Starting reflector *v1.StorageClass (15m0s) from sigs.k8s.io/sig-storage-lib-external-provisioner/v6/controller/controller.go:872
I0427 04:14:57.635287       1 reflector.go:255] Listing and watching *v1.StorageClass from sigs.k8s.io/sig-storage-lib-external-provisioner/v6/controller/controller.go:872
I0427 04:14:57.635200       1 reflector.go:219] Starting reflector *v1.PersistentVolume (15m0s) from sigs.k8s.io/sig-storage-lib-external-provisioner/v6/controller/controller.go:869
I0427 04:14:57.635728       1 reflector.go:255] Listing and watching *v1.PersistentVolume from sigs.k8s.io/sig-storage-lib-external-provisioner/v6/controller/controller.go:869
I0427 04:14:57.734504       1 shared_informer.go:270] caches populated
I0427 04:14:57.735220       1 controller.go:884] Started provisioner controller soda-csi_csi-rbdplugin-provisioner-68b5cdf677-rlhqq_79ddf7a3-6411-4213-a9ba-a19b6ca12d44!
I0427 04:15:54.469345       1 controller.go:1332] provision "default/raw-block-pvc" class "csi-rbd-ceph-sc": started
I0427 04:15:54.470475       1 event.go:282] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"raw-block-pvc", UID:"8ea03417-9275-45c8-b92e-58c467ccbf8d", APIVersion:"v1", ResourceVersion:"229533", FieldPath:""}): type: 'Normal' reason: 'Provisioning' External provisioner is provisioning volume for claim "default/raw-block-pvc"
I0427 04:15:56.978651       1 connection.go:182] GRPC call: /csi.v1.Identity/GetPluginInfo
I0427 04:15:56.978691       1 connection.go:183] GRPC request: {}
I0427 04:15:56.980624       1 connection.go:185] GRPC response: {"name":"rbd.csi.ceph.com","vendor_version":"canary"}
I0427 04:15:56.980736       1 connection.go:186] GRPC error: <nil>
I0427 04:15:56.990564       1 controller.go:801] CreateVolumeRequest name:"pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d" capacity_range:<required_bytes:1073741824 > volume_capabilities:<block:<> access_mode:<mode:SINGLE_NODE_WRITER > > parameters:<key:"attachMode" value:"rw" > parameters:<key:"clusterID" value:"4ac5251b-114f-4044-9bec-2d27fadad502" > parameters:<key:"csi.storage.k8s.io/pv/name" value:"pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d" > parameters:<key:"csi.storage.k8s.io/pvc/name" value:"raw-block-pvc" > parameters:<key:"csi.storage.k8s.io/pvc/namespace" value:"default" > parameters:<key:"imageFeatures" value:"layering" > parameters:<key:"imageFormat" value:"2" > parameters:<key:"pool" value:"osdsrbd" > parameters:<key:"profile" value:"f63a7fcc-6e83-49cc-81c1-8627fc8e0f52" > secrets:<key:"userID" value:"kube" > secrets:<key:"userKey" value:"AQDjZYFgOw/FFxAASWpRJz+H0yI3LU4x6Ct2zA==" >
I0427 04:15:56.991220       1 connection.go:182] GRPC call: /csi.v1.Controller/CreateVolume
I0427 04:15:56.991234       1 connection.go:183] GRPC request: {"capacity_range":{"required_bytes":1073741824},"name":"pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d","parameters":{"attachMode":"rw","clusterID":"4ac5251b-114f-4044-9bec-2d27fadad502","csi.storage.k8s.io/pv/name":"pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d","csi.storage.k8s.io/pvc/name":"raw-block-pvc","csi.storage.k8s.io/pvc/namespace":"default","imageFeatures":"layering","imageFormat":"2","pool":"osdsrbd","profile":"f63a7fcc-6e83-49cc-81c1-8627fc8e0f52"},"secrets":"***stripped***","volume_capabilities":[{"AccessType":{"Block":{}},"access_mode":{"mode":1}}]}
I0427 04:15:57.529205       1 connection.go:185] GRPC response: {"volume":{"capacity_bytes":1073741824,"volume_context":{"attachMode":"rw","clusterID":"4ac5251b-114f-4044-9bec-2d27fadad502","csi.storage.k8s.io/pv/name":"pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d","csi.storage.k8s.io/pvc/name":"raw-block-pvc","csi.storage.k8s.io/pvc/namespace":"default","imageFeatures":"layering","imageFormat":"2","imageName":"csi-vol-44027e9e-a70f-11eb-a408-0242ac110003","journalPool":"osdsrbd","pool":"osdsrbd","profile":"f63a7fcc-6e83-49cc-81c1-8627fc8e0f52"},"volume_id":"0001-0024-4ac5251b-114f-4044-9bec-2d27fadad502-0000000000000001-44027e9e-a70f-11eb-a408-0242ac110003"}}
I0427 04:15:57.529540       1 connection.go:186] GRPC error: <nil>
I0427 04:15:57.529567       1 controller.go:832] create volume rep: {CapacityBytes:1073741824 VolumeId:0001-0024-4ac5251b-114f-4044-9bec-2d27fadad502-0000000000000001-44027e9e-a70f-11eb-a408-0242ac110003 VolumeContext:map[attachMode:rw clusterID:4ac5251b-114f-4044-9bec-2d27fadad502 csi.storage.k8s.io/pv/name:pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d csi.storage.k8s.io/pvc/name:raw-block-pvc csi.storage.k8s.io/pvc/namespace:default imageFeatures:layering imageFormat:2 imageName:csi-vol-44027e9e-a70f-11eb-a408-0242ac110003 journalPool:osdsrbd pool:osdsrbd profile:f63a7fcc-6e83-49cc-81c1-8627fc8e0f52] ContentSource:<nil> AccessibleTopology:[] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}
I0427 04:15:57.529671       1 controller.go:908] successfully created PV pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d for PVC raw-block-pvc and csi volume name 0001-0024-4ac5251b-114f-4044-9bec-2d27fadad502-0000000000000001-44027e9e-a70f-11eb-a408-0242ac110003
I0427 04:15:57.529693       1 controller.go:924] successfully created PV {GCEPersistentDisk:nil AWSElasticBlockStore:nil HostPath:nil Glusterfs:nil NFS:nil RBD:nil ISCSI:nil Cinder:nil CephFS:nil FC:nil Flocker:nil FlexVolume:nil AzureFile:nil VsphereVolume:nil Quobyte:nil AzureDisk:nil PhotonPersistentDisk:nil PortworxVolume:nil ScaleIO:nil Local:nil StorageOS:nil CSI:&CSIPersistentVolumeSource{Driver:rbd.csi.ceph.com,VolumeHandle:0001-0024-4ac5251b-114f-4044-9bec-2d27fadad502-0000000000000001-44027e9e-a70f-11eb-a408-0242ac110003,ReadOnly:false,FSType:,VolumeAttributes:map[string]string{attachMode: rw,clusterID: 4ac5251b-114f-4044-9bec-2d27fadad502,csi.storage.k8s.io/pv/name: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d,csi.storage.k8s.io/pvc/name: raw-block-pvc,csi.storage.k8s.io/pvc/namespace: default,imageFeatures: layering,imageFormat: 2,imageName: csi-vol-44027e9e-a70f-11eb-a408-0242ac110003,journalPool: osdsrbd,pool: osdsrbd,profile: f63a7fcc-6e83-49cc-81c1-8627fc8e0f52,storage.kubernetes.io/csiProvisionerIdentity: 1619496897051-8081-soda-csi,},ControllerPublishSecretRef:nil,NodeStageSecretRef:&SecretReference{Name:csi-rbd-secret,Namespace:default,},NodePublishSecretRef:nil,ControllerExpandSecretRef:&SecretReference{Name:csi-rbd-secret,Namespace:default,},}}
I0427 04:15:57.568786       1 controller.go:1439] provision "default/raw-block-pvc" class "csi-rbd-ceph-sc": volume "pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d" provisioned
I0427 04:15:57.568847       1 controller.go:1456] provision "default/raw-block-pvc" class "csi-rbd-ceph-sc": succeeded
I0427 04:15:57.568871       1 volume_store.go:154] Saving volume pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d
I0427 04:15:57.658619       1 volume_store.go:157] Volume pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d saved
I0427 04:15:57.658794       1 controller.go:1093] Claim processing succeeded, removing PVC 8ea03417-9275-45c8-b92e-58c467ccbf8d from claims in progress
I0427 04:15:57.658896       1 controller.go:1332] provision "default/raw-block-pvc" class "csi-rbd-ceph-sc": started
I0427 04:15:57.658954       1 controller.go:1341] provision "default/raw-block-pvc" class "csi-rbd-ceph-sc": persistentvolume "pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d" already exists, skipping
I0427 04:15:57.659013       1 controller.go:1095] Stop provisioning, removing PVC 8ea03417-9275-45c8-b92e-58c467ccbf8d from claims in progress
I0427 04:15:57.659192       1 event.go:282] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"raw-block-pvc", UID:"8ea03417-9275-45c8-b92e-58c467ccbf8d", APIVersion:"v1", ResourceVersion:"229533", FieldPath:""}): type: 'Normal' reason: 'ProvisioningSucceeded' Successfully provisioned volume pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d


```

Ceph driver is triggered for volume creation.

`Ceph RBD driver logs`

```go
kubectl logs  csi-rbdplugin-provisioner-68b5cdf677-rlhqq csi-rbdplugin
.
.
.
.
I0427 04:15:56.980031       1 identityserver-default.go:36] ID: 16 Using default GetPluginInfo
I0427 04:15:56.980104       1 utils.go:173] ID: 16 GRPC response: {"name":"rbd.csi.ceph.com","vendor_version":"canary"}
I0427 04:15:56.994114       1 utils.go:162] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d GRPC call: /csi.v1.Controller/CreateVolume
I0427 04:15:56.995447       1 utils.go:166] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d GRPC request: {"capacity_range":{"required_bytes":1073741824},"name":"pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d","parameters":{"attachMode":"rw","clusterID":"4ac5251b-114f-4044-9bec-2d27fadad502","csi.storage.k8s.io/pv/name":"pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d","csi.storage.k8s.io/pvc/name":"raw-block-pvc","csi.storage.k8s.io/pvc/namespace":"default","imageFeatures":"layering","imageFormat":"2","pool":"osdsrbd","profile":"f63a7fcc-6e83-49cc-81c1-8627fc8e0f52"},"secrets":"***stripped***","volume_capabilities":[{"AccessType":{"Block":{}},"access_mode":{"mode":1}}]}
I0427 04:15:57.008680       1 rbd_util.go:976] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d setting disableInUseChecks: false image features: [layering] mounter: rbd
I0427 04:15:57.229215       1 omap.go:84] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d got omap values: (pool="osdsrbd", namespace="", name="csi.volumes.default"): map[]
I0427 04:15:57.236272       1 omap.go:148] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d set omap keys (pool="osdsrbd", namespace="", name="csi.volumes.default"): map[csi.volume.pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d:44027e9e-a70f-11eb-a408-0242ac110003])
I0427 04:15:57.238187       1 omap.go:148] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d set omap keys (pool="osdsrbd", namespace="", name="csi.volume.44027e9e-a70f-11eb-a408-0242ac110003"): map[csi.imagename:csi-vol-44027e9e-a70f-11eb-a408-0242ac110003 csi.volname:pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d csi.volume.owner:default])
I0427 04:15:57.238244       1 rbd_journal.go:472] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d generated Volume ID (0001-0024-4ac5251b-114f-4044-9bec-2d27fadad502-0000000000000001-44027e9e-a70f-11eb-a408-0242ac110003) and image name (csi-vol-44027e9e-a70f-11eb-a408-0242ac110003) for request name (pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d)
I0427 04:15:57.249912       1 rbd_util.go:228] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d rbd: create osdsrbd/csi-vol-44027e9e-a70f-11eb-a408-0242ac110003 size 1024M (features: [layering]) using mon 192.168.1.59:6789
I0427 04:15:57.308369       1 controllerserver.go:476] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d created volume pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d backed by image csi-vol-44027e9e-a70f-11eb-a408-0242ac110003
I0427 04:15:57.527605       1 omap.go:148] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d set omap keys (pool="osdsrbd", namespace="", name="csi.volume.44027e9e-a70f-11eb-a408-0242ac110003"): map[csi.imageid:2b8edabf6b50])
I0427 04:15:57.528163       1 utils.go:173] ID: 17 Req-ID: pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d GRPC response: {"volume":{"capacity_bytes":1073741824,"volume_context":{"attachMode":"rw","clusterID":"4ac5251b-114f-4044-9bec-2d27fadad502","csi.storage.k8s.io/pv/name":"pvc-8ea03417-9275-45c8-b92e-58c467ccbf8d","csi.storage.k8s.io/pvc/name":"raw-block-pvc","csi.storage.k8s.io/pvc/namespace":"default","imageFeatures":"layering","imageFormat":"2","imageName":"csi-vol-44027e9e-a70f-11eb-a408-0242ac110003","journalPool":"osdsrbd","pool":"osdsrbd","profile":"f63a7fcc-6e83-49cc-81c1-8627fc8e0f52"},"volume_id":"0001-0024-4ac5251b-114f-4044-9bec-2d27fadad502-0000000000000001-44027e9e-a70f-11eb-a408-0242ac110003"}}

```


***Note***: 
The Deployment scripts for ceph rbd driver are taken from https://github.com/ceph/ceph-csi . 


## Simple steps to integrate and experiment any csi driver

### Step 1
Retrieve the deployments scripts/mechanism for the respective csi driver.
Update the csi provisioner image used with the soda csi provisioner image(sodafoundation/soda-csi-provisioner:v1.6.0).

### Step 2
Deploy the respective csi driver. If the driver is already deployed, it can be patched by editing deployment after replacing the image (step 1). 

### Step 3
Create soda profile corresponding the driver getting added . Then involves mentioning driver name in profile. Reference provided [here](https://youtu.be/ytXY_dKQCYg). This will help to identify the driver for which volume provisioning is requested for.

### Step 4
Create storage class specifying the above profile as shown in below example. Then proceed with Volume provisioning (pvc).
```go
kubectl get sc soda-high-io -o yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  creationTimestamp: "2021-01-07T10:31:11Z"
  name: soda-high-io
  resourceVersion: "13934"
  selfLink: /apis/storage.k8s.io/v1/storageclasses/soda-high-io
  uid: 75e47bca-50d3-11eb-a5ff-080027310244
parameters:
  profile: rbd.csi.ceph.com
provisioner: soda-csi
reclaimPolicy: Delete
volumeBindingMode: Immediate

```

