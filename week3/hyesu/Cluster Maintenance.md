# Operating System Upgrade

Node를 제거하는 시나리오

if. Node가 5분 이상 다운되면 해당 Node에서 Pod가 종료

Pod가 복구되길 기다리는 시간은 Pod Eviction Timeout으로 Controller 관리자에게 5분이라는 기본 값 설정

⇒ 안전한 방법: 모든 Workload를 의도적으로 Drain

작업이 클러스터 내 다른 노드로 이동하도록~

```bash
kubectl drain node-1
```

drain 하면 Pod가 정상적으로 종료되고 다른 Node에서 재현

```bash
kubectl uncordon node-1
```

⇒ 다시 Pod 스케줄 할 수 있음

단, 이미 이동된 pod가 자동으로 돌아가진 않음

```bash
kubectl cordon node-2
```

⇒ 단순히 해당 노드에 새 Pod가 스케줄링 되지 않도록 함

# Kubernetes Releases

```bash
kubectl get nodes # 버전 확인 가능
```

- 버전
    - Major
    - Minor : 새로운 기능
    - Patch : 버그 패치

- 순서
    - alpha : 버그 고치고 개선한 뒤 ~ → 버그 존재할 수 있음
    - beta
    - main stable release

- 참고
    - 패키지 다운받으면 kube-apiserver, Contoller manager. kube-scheduler, kubelet, kube-proxy, kubectl 버전 똑같음
    - 단, ETCD Cluster, CoreDNS는 다른 프로젝트로, 고유한 버전 존재함

**References**

[https://kubernetes.io/docs/concepts/overview/kubernetes-api/](https://kubernetes.io/docs/concepts/overview/kubernetes-api/)

Here is a link to kubernetes documentation if you want to learn more about this topic (You don't need it for the exam though):

[https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)

[https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md)

# Cluster Upgrade Process

- kube-apiserver는 주요 구성요소이고 다른 구성 요소들과 통신하는 구성요소
    
     ⇒ 다른 구성요소는 kube-apiserver보다 높으면 안됨
    
- controller-manager, kube-scheduler는 1단계 낮아도 됨
- kubelet, kube-proxy는 2단계 낮아도 됨
    
    ex. kube-apiserver 1.10
    
    controller-manager, kube-scheduler 1.10 1.9
    
    kubelet, kube-proxy 1.10 1.9 1.8
    
- kubectl 은 kube-apiserver보다 +-1 가능

Kubernetes는 공식적으로 3개 release만 지원

한 번에 마이너 버전 하나씩 업그레이드하는 것 권장

# Cluster Upgrade Process

- kube-apiserver는 주요 구성요소이고 다른 구성 요소들과 통신하는 구성요소
    
     ⇒ 다른 구성요소는 kube-apiserver보다 높으면 안됨
    
- controller-manager, kube-scheduler는 1단계 낮아도 됨
- kubelet, kube-proxy는 2단계 낮아도 됨
    
    ex. kube-apiserver 1.10
    
    controller-manager, kube-scheduler 1.10 1.9
    
    kubelet, kube-proxy 1.10 1.9 1.8
    
- kubectl 은 kube-apiserver보다 +-1 가능

Kubernetes는 공식적으로 3개 release만 지원

한 번에 마이너 버전 하나씩 업그레이드하는 것 권장

### Cluster Upgrade

1. Master Node Upgrade
    - Master이 업그레이드 되는 동안 구성 요소 잠시 다운됨
2. Worker Node Upgrade
    - Master 업그레이드 동안 영향 X
    - 전략
        - 한꺼번에 업그레이드 → 앱 접속 불가능
        - 한번에 노드 하나씩 업그레이드
        - 클러스터에 새 버전의 노드 추가

```bash
kubeadm upgrade plan # 클러스터 업그레이드할 수 있는 버전 확인
```

클러스터 업그레이드 전에 kubeadm 업그레이드 필요!

- **Master node 업그레이드**

현재 1.11 → 1.13

```bash
apt-get upgrade -y kubeadm=1.12.0-00 # kubeadm 1.12로 업그레이드

kubeadm upgrade apply v1.12.0

kubectl get nodes
# 결과물 보면 kube-apiserver로 등록된 각각 노드에서 kubelet 버전 보여줌
# v1.11.3

# 셋업에 따라 마스터 노드에서 kubelet이 실행될수도, 아닐수도 있음
# 이 경우, kubeadm과 함께 배포된 클러스터는 마스터 노드에 kubelet이 있음

apt-get upgrade -y kubelet=1.12.0-00
systemctl restart kubelet

kubectl get nodes
# master의 version은 1.12.0
# workernode의 version은 여전히 1.11.3
```

- **Worker node 업그레이드**

하나씩 업데이트

```bash
# Pod를 옮겨야 함!

kubectl drain node-1 # master node 명령

# 각 node마다 명령
apt-get upgrade -y kubeadm=1.12.0=00
apt-get upgrade -y kubelet=1.12.0=00
kubeadm upgrade node config --kubelet-version v1.12.0
systemctl restart kubelet

kubectl uncordon node-1 # master node 명령
```

https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-upgrade/

kubelet은 수동으로 업그레이드 해야함!

1. kubeadm 업그레이드
2. kubectl drain
3. kubelet과 kubectl → apt upgrade
4. sudo systemctl restart kubelet
5. kubectl uncordon

# Backup and Restore

Resource Configuration, ETCD Cluster, Persistent Volumes, …

- **Resource Configuration**
    - 선언적 접근법 선호
    - 선언적 접근을 위한 파일  → repo에 저장 권장
    
    if, 필요한 리소스를 기록하지 않았다면?
    
    ⇒ kube-apiserver 쿼리!
    
    kubectl 을 통해 kube-apiserver를 쿼리하거나 API 서버에 직접 엑세스 함으로써 클러스터에 생성된 모든 개체에 대한 리소스 구성을 복사해 저장
    
    ```bash
    kubectl get all --all-namespace -o yaml > all-deploy-services.yaml
    ```
    

- **ETCD**
    - 클러스터 자체, 노드 및 클러스터 내부에서 생성된 모든 리소스 저장
    - ETCD는 모든 데이터가 저장된 장소 명시함
        
        ```bash
        # etcd.service
        ExecStart=/user/local/bin/etcd \\
        	...
        	--data-dir=/var/lib/etcd
        
        ```
        
    - build-in snapshot solution 존재
        
        ```bash
        ETCDCTL_API=3 etcdctl \
        	snapshot save snapshot.db
        	
        ETCDCTL_API=3 etcdctl \
        	snapshot status snapshot.db
        ```
        
        ```bash
        # 복원
        service kube-apiserver stop
        
        ETCDCTL_API=3 etcdctl \
        	snapshot restore snapshot.db \
        	--data-dir /var/lib/etcd-from-backup
        ```
        
        ```bash
        # etcd.service
        ExecStart=/user/local/bin/etcd \\
        	...
        	--data-dir=/var/lib/etcd-from-backup
        
        ```
        
        ```bash
        systemctl daemon-reload
        service etcd restart
        service kube-apiserver start
        ```
        
        - etcd 명령어는 인증서 파일 지정 !!
        

### Working with ETCDCTL

etcdctl은 etcd의 CLI(Client) 도구이며, Kubernetes에서는 etcd v3를 사용

✅ 사용 전 `ETCDCTL_API=3` 환경 변수를 설정해야 함

```bash
bash
복사편집
export ETCDCTL_API=3

```

✅ etcd는 Master 노드에서 Static Pod로 실행되며, 기본 포트는 `127.0.0.1:2379`

✅ TLS가 활성화되어 있으므로 백업 및 복구 시 다음 옵션이 필수

- `-cacert` : CA 인증서 파일
- `-cert` : 클라이언트 인증서 파일
- `-key` : 클라이언트 키 파일
- `-endpoints=[127.0.0.1:2379]` : etcd 엔드포인트 설정

✅ 백업 명령어 (`snapshot save`)

```bash
bash
복사편집
etcdctl snapshot save <backup-file> --cacert=<ca.pem> --cert=<cert.pem> --key=<key.pem> --endpoints=[127.0.0.1:2379]

```

✅ 복구 명령어 (`snapshot restore`)

```bash
bash
복사편집
etcdctl snapshot restore <backup-file> -h

```

➡️ `-h` 옵션을 사용하면 추가적인 설정 옵션 확인 가능.

📌 추가 정보:

백업 및 복구 관련 자세한 설명은 Kubernetes Backup and Restore Lab 솔루션 영상을 참고. 🚀
