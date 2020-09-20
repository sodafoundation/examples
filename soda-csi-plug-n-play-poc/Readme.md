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
kubectl create -f deploy/kubernetes/ibm/
```

#### Step2:
Deploy LVM CSI operator and driver along with soda-csi-provisioner
```
kubectl create -f deploy/kubernetes/lvm/
```

#### Step3:
Deploy StorageClass and PVC with profile as 'block.csi.ibm.com'
```go
kubectl create -f deploy/kubernetes/demo
```

This will create the StorageClass but the PVC will fail as the IBM backend is not available and the soda-csi-provisioner will send an error.

`IBM soda-csi-provisioner logs`

```go


```

Whereas in LVM soda-csi-provisioner the PVC creation will skip as the profile does not matches to LVM.
`LVM soda-csi-provisioner logs`


```go


```


#### Step4:
Deploy StorageClass and PVC with profile as 'csi-lvmplugin'
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


```


Whereas in LVM soda-csi-provisioner the PVC creation will proceed as the profile matches to LVM.
`LVM soda-csi-provisioner logs`

```go


```
