# OS Upgrades

작업이 클러스터내 다른 노드로 이동

- kubectl drain node-1

해당 노드에 새 포드가 스케쥴링 되지 않도록 함

- kubectl cordon node-2

스케줄링 재개

- kubectl uncordon node-1

### Kubernetes Software Version

버전확인

- kubectl get nodes

## Practice Test - **OS Upgrades**

```yaml
k get pods -o wide

kubectl drain node01 --ignore-daemonsets

k uncordon node01

kubectl cordon node01
```

---

# Cluster Upgrade Process

### kubeadm - upgrade

kubeadm 업그레이드 계획 확인 가능

- kubeadm upgrade plan

마스터 노드 업그레이드

1. kubeadm 도구 업그레이드 : apt-get upgrade -y kubeadm=1.12.0-00
2. kubeadm 업그레이드 적용

: kubeadm upgrade apply v1.12.0

3. 노드버전 확인

: kubectl get nodes

4. kubelet 업그레이드 및 kubelet 재시작

: apt-get upgrade -y kubelet=1.12.0-00

: systemctl restart kubelet

5. 노드버전 재확인

: kubectl get nodes

워커 노드 업그레이드

1. node-1 작업 다른 노드로 이동
: kubectl drain node-1
2. kubeadm 도구 업그레이드
: apt-get upgrade -y kubeadm-1.12.0-00
3. kubelet 업그레이드
: apt-get upgrade -y kubelet-1.12.0-00
4. 새 kubelet 버전을 위해 노드 구성 업데이트
: kubeadm upgrade node config --kubelet-version v1.12.0
5. kubelet 재시작
: systemctl restart kubelet
6. node-1노드에 스케줄링 재개
: kubectl uncordon node-1
node-2
kubectl drain node-2
kubectl uncordon node-2
node-3
kubectl drain node-3
kubectl uncordon node-3

# Backup and Restore Methods

### Resource Configuration

: 리소스 파일은 Github에 저장하는게 좋음

: 모든 namespace yaml 파일 저장

: kubectl get all --all-namespaces -o yaml › all-deploy-services.yam

### ETCD Cluster

: 클러스터 자체에 관한 정보와 노드 및 클러스터 내부에서 생성된 모든 리소스 저장

- 백업

: ETCDCTL_APT-3 etcdctl snapshot save snapshot. db

: ETCDCTL_API=3 etcdctl snapshot status snapshot.db

- 복원

: service kube-apiserver stop

: ETCDCTL_API=3 etcdctl snapshot restore snapshot.db --data-dir /var/lib/etcd-from-backup

: etcd.service 파일 수정

: --data-dir=/var/lib/etcd-from-backup

: systemctl daemon-reload : service etcd restart : service kube-apiserver start

- ETCDCTL

: ETCDCTL_API 버전이 3으로 설정되어 있는지 확인필요

: export ETCDCTL_API=3

: etcdctl version

: ETCD 데이터베이스는 TLS가 활성화 되어 있기 때문에 다음 옵션은 필수

: --cacert

: --cert

: --endpoints=[127.0.0.1:2379]

: --key

스냅샷 복원에 대한 도움말

: etcdctl snapshot restore -h

: ETCDCTL_API=3 etcactl

snapshot restore snapshot.db \

- -endpoints=https://127.0.0.1:2379 \
- -cacert=/etc/etcd/ca.crt
- -cert=/etc/etcd/etcd-server.crt
- - key-/etc/etcd/etcd-server.key

### 시험 팁

: 시험에서는 연습시험과 같이 자신이 한것이 올바른지 아닌지 알 수 없다. 직접확인해야한다.

: 예를들어, 특정 이미지로 포드를 생성하는 것이라는 질문이라면

kubectl describe pod 명령을 실행하여 파드가 올바른 이름과 올바른 이미지로 생성되었는지 확인해야한다.

## Practice Test - **Cluster Upgrade Process**

```yaml
kubeadm upgrade plan

k cordon controlplane

sudo kubeadm upgrade apply v1.32.0

kubectl uncordon controlplane

kubectl drain node01 --ignore-daemonsets
```