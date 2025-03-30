
ubuntu@nks-bastion:~$ k describe node ai-prd-mem-01-w-cloud-club
Name:               ai-prd-mem-01-w-cloud-club
Roles:              <none>
Labels:             beta.kubernetes.io/arch=amd64
                    beta.kubernetes.io/instance-type=SVR.VSVR.HIMEM.C004.M032.NET.SSD.B050.G002
                    beta.kubernetes.io/os=linux
                    failure-domain.beta.kubernetes.io/region=1
                    failure-domain.beta.kubernetes.io/zone=2
                    kubernetes.io/arch=amd64
                    kubernetes.io/hostname=ai-prd-mem-01-w-cloud-club
                    kubernetes.io/os=linux
                    ncloud.com/nks-nodepool=ai-prd-mem-01
                    node.kubernetes.io/instance-type=SVR.VSVR.HIMEM.C004.M032.NET.SSD.B050.G002
                    nodeId=9999999999
                    regionNo=1
                    topology.kubernetes.io/region=1
                    topology.kubernetes.io/zone=2
                    zoneNo=2
Annotations:        alpha.kubernetes.io/provided-node-ip: 1.2.3.4
                    csi.volume.kubernetes.io/nodeid: {"blk.csi.ncloud.com":"9999999999","nas.csi.ncloud.com":"ai-prd-mem-01-w-cloud-club"}
                    node.alpha.kubernetes.io/ttl: 0
                    volumes.kubernetes.io/controller-managed-attach-detach: true
CreationTimestamp:  Mon, 30 Sep 2024 05:03:15 +0000
Taints:             atype=mem01:NoSchedule
Unschedulable:      false
