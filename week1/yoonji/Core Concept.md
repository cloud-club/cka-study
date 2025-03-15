# Cluster Architecture
![image](https://github.com/user-attachments/assets/53f0b9ae-4db3-4ede-bdab-ddc580b88132)

1. 마스터노드
- ETCD Cluster, kube-apiserver, Kube Controller Manager, kube-scheduler
2. 워커 노드
- kubelet, kube-proxy, container runtime engine



# Master
### ETCD 
> 분산된, 신뢰할 수 있는 key-value 저장소
- key-value 저장소: 새 정보가 추가될 때마다 테이블 전체가 영향을 받는 데이터베이스와 달리, 한 파일만 변경하면 다른 파일들은 그대로이다.
- ETCD를 실행하면 2379 포트의 신호를 받는 서비스가 실행된다.

#### 쿠버네티스의 ETCD
> Node, Pod, Configs, Secrets 등 클러스터에 관한 정보를 저장한다.

kubectl을 실행할 때 모든 정보는 etcd에서 받아온다. (ex. 노드, 레플리카셋 등을 추가할 때마다 업데이트)
https://${INTERNAL_IP}:2379가 etcd의 주소가 된다. 

### Kube-API Server
#### 명령 수행 과정
1. kubectl 명령을 실행하면, kube-apiserver에 도달
2. 요청과 유효성 확인
3. etcd에서 데이터를 가져와 응답

#### 파드 생성 과정
1. 요청의 인증, 유효성 확인
2. kube-apiserver가 노드에 할당하지 않고 파드 생성
3. etcd 서버에 정보 업데이트
4. scheduler manager가 apiserver를 지속적으로 모니터링하다가, 할당되지 않은 파드를 확인
5. scheduler manager가 적절한 노드를 식별해 파드 할당 후 kube-apiserver와 통신
6. etcd에 정보 업데이트
7. api server는 해당 정보를 node -> kubelet에 전달
8. kubelet은 노드에 파드를 생성하고, container runtime engine을 통해 이미지 배포
9. kubelet은 상태를 api server로 업데이트+ etcd도 정보 업데이트
    

### Kube Controller Manager
> 다양한 컨트롤러들을 관리한다. 컨트롤러는 시스템 내 구성 요소들의 상태를 지속적으로 모니터링하고, 시스템 전체를 desired state로 만든다. 

- k8s에서 deployment, service, namespace 등은 컨트롤러를 통해 구현된다.
- 컨트롤러 매니저는 프로세스이다. 
- 컨트롤러들은 컨트롤러 매니저로 패키지화되어 있다. 컨트롤러 매니저를 설치하면 다른 컨트롤러도 설치된다.  
- 

### Kube Scheduler
> 노드, 파드의 스케줄링을 담당한다.
> 스케줄러 매니저는 어떤 파드가 어떤 노드에 들어갈지만 결정한다. 단, 파드를 해당 노드에 배치하는 건 kubelet이 한다.

cpu와 메모리에 따라 가장 적합한 노드를 탐색한다. 


# Worker Node
### Kubelet
> 파드의 상태와 컨테이너를 지속적으로 모니터링하고 kube-apiserver에 보고한다.

kubeadm을 사용하면 자동으로 kubelet을 배포하지 않는다. 워커노드에 수동으로 설치해야 한다. 

### Kube Proxy
> 각 노드에서 실행되는 프로세스로, 새 서비스가 생성될 때마다 각 노드에 적절한 규칙을 만들어 그 서비스로 트래픽을 전달한다. 

각 노드에 생성된 iptables 규칙을 이용해 서비스의 IP로 향하는 트래픽을 향하게 한다. 


# Pods
> 쿠버네티스에서 만들수 있는 가장 작은 단위

- 트래픽이 증가해 추가 인스턴스가 필요하다면, 노드를 증가하거나, 노드 내 추가 파드를 배포한다. 
- 파드는 보통 컨테이너와 일대일로 연결된다.
- (드물지만) 그러나 파드 하나에 여러 컨테이너도 가능하다. 보통 helper container로 존재한다. (ex. 업로드한 파일 처리의 기능 수행)
- 두 컨테이너는 같은 네트워크 공간을 공유하며 서로 통신할 수 있다. -> localhost



### Pods with yaml
항상 4개의 필드를 포함한다. (루트 레벨)
1. apiVersion: 쿠버네티스 api 버전.
   - pod(v1), Service(v1), ReplicaSet(apps/v1), Deployment(apps/v1)
2. kind
3. metadata: name, labels
4. spec: containers(name, image), 
- container는 dictionary 형태. -(dash)로 구분

# ReplicaSets
> 하나가 실패해도 한 개 이상의 인스턴스가 동시에 실행되도록 한다.
> 여러개의 파드를 만들어 로드밸런싱이 되도록 한다.

# Deployments
> replicaSet의 상위 개념인 deployment를 통해 인스턴스를 매끄럽게 업그레이드하도록 한다.
> 

- 롤링 업데이트: 하나씩 버전 업그레이드
- deployment와 replicaset의 파일 구성은 매우 유사하다.
- replicaset은 파드를 생성하며, deployment는 자동으로 replicaset을 생성해 deployment라는 k8s의 새 오브젝트를 생성한다.
