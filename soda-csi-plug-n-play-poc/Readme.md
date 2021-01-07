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
Deploy StorageClass and PVC with profile as 'rbd.csi.ceph.com'
```go
kubectl create -f deploy/kubernetes/demo
```
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

StorageClass is created as shown above but the PVC will fail as the Ceph backend is not available and the soda-csi-provisioner will send an error.

`Ceph RBD soda-csi-provisioner logs`

```go
kubectl logs csi-rbdplugin-provisioner-6b8b9d99fd-x4wn6 csi-provisioner 
I0107 10:31:19.485001       1 controller.go:1284] provision "default/demo-pvc-file-system" class "soda-high-io": started
I0107 10:31:19.488100       1 connection.go:182] GRPC call: /csi.v1.Identity/GetPluginInfo
I0107 10:31:19.488114       1 connection.go:183] GRPC request: {}
I0107 10:31:19.489971       1 event.go:281] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"demo-pvc-file-system", UID:"75e213de-50d3-11eb-a5ff-080027310244", APIVersion:"v1", ResourceVersion:"13950", FieldPath:""}): type: 'Normal' reason: 'Provisioning' External provisioner is provisioning volume for claim "default/demo-pvc-file-system"
I0107 10:31:19.490741       1 connection.go:185] GRPC response: {"name":"rbd.csi.ceph.com","vendor_version":"canary"}
I0107 10:31:19.491201       1 connection.go:186] GRPC error: <nil>
I0107 10:31:19.491975       1 controller.go:462] The Backend Driver Name and provisioner name is : rbd.csi.ceph.com , soda-csi
I0107 10:31:19.492053       1 controller.go:466] The parameters in the StorageClass are  : profile ===== rbd.csi.ceph.com

I0107 10:31:19.515052       1 controller.go:620] CreateVolumeRequest {Name:pvc-75e213de-50d3-11eb-a5ff-080027310244 CapacityRange:required_bytes:2147483648  VolumeCapabilities:[mount:<fs_type:"ext4" > access_mode:<mode:SINGLE_NODE_WRITER > ] Parameters:map[profile:rbd.csi.ceph.com] Secrets:map[] VolumeContentSource:<nil> AccessibilityRequirements:<nil> XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}
I0107 10:31:19.515216       1 connection.go:182] GRPC call: /csi.v1.Controller/CreateVolume
I0107 10:31:19.515269       1 connection.go:183] GRPC request: {"capacity_range":{"required_bytes":2147483648},"name":"pvc-75e213de-50d3-11eb-a5ff-080027310244","parameters":{"csi.storage.k8s.io/pv/name":"pvc-75e213de-50d3-11eb-a5ff-080027310244","csi.storage.k8s.io/pvc/name":"demo-pvc-file-system","csi.storage.k8s.io/pvc/namespace":"default","profile":"rbd.csi.ceph.com"},"volume_capabilities":[{"AccessType":{"Mount":{"fs_type":"ext4"}},"access_mode":{"mode":1}}]}

```

Ceph driver is triggered for volume creation.

`Ceph RBD driver logs`

```go
kubectl logs csi-rbdplugin-provisioner-6b8b9d99fd-x4wn6 csi-rbdplugin
I0107 10:31:06.447514       1 utils.go:132] ID: 19 GRPC call: /csi.v1.Identity/Probe
I0107 10:31:06.447889       1 utils.go:133] ID: 19 GRPC request: {}
I0107 10:31:06.449711       1 utils.go:138] ID: 19 GRPC response: {}
I0107 10:31:19.490255       1 utils.go:132] ID: 20 GRPC call: /csi.v1.Identity/GetPluginInfo
I0107 10:31:19.490396       1 utils.go:133] ID: 20 GRPC request: {}
I0107 10:31:19.490462       1 identityserver-default.go:36] ID: 20 Using default GetPluginInfo
I0107 10:31:19.490541       1 utils.go:138] ID: 20 GRPC response: {"name":"rbd.csi.ceph.com","vendor_version":"canary"}
I0107 10:31:19.520057       1 utils.go:132] ID: 21 Req-ID: pvc-75e213de-50d3-11eb-a5ff-080027310244 GRPC call: /csi.v1.Controller/CreateVolume
I0107 10:31:19.535219       1 utils.go:133] ID: 21 Req-ID: pvc-75e213de-50d3-11eb-a5ff-080027310244 GRPC request: {"capacity_range":{"required_bytes":2147483648},"name":"pvc-75e213de-50d3-11eb-a5ff-080027310244","parameters":{"csi.storage.k8s.io/pv/name":"pvc-75e213de-50d3-11eb-a5ff-080027310244","csi.storage.k8s.io/pvc/name":"demo-pvc-file-system","csi.storage.k8s.io/pvc/namespace":"default","profile":"rbd.csi.ceph.com"},"volume_capabilities":[{"AccessType":{"Mount":{"fs_type":"ext4"}},"access_mode":{"mode":1}}]}

```


***Note***: 
The Deployment scripts for ceph rbd driver are taken from https://github.com/ceph/ceph-csi . 


