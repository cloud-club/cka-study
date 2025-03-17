# Core Concept

### Cluster Architecture
- Master
: ETCD CLUSTER
: kube-apiserver
: Kube Controller-Manager
: kube-scheduler

- Worker Nodes
: kubelet
: Kube-proxy
: Container Runtime Engine

### Docker-vs-ContainerD
![alt text](image.png)

### ETCD
: distributed reliable key-value store that is Simple, Secure, & Fast

### Kube-API Server
kubectl 명령실행 - kube-api server 요청인증 및 확인 - ETCD 클러스터에서 요청된 정보와 함께 응답

- Authenticate User
- Validate Request
- Retrieve data
- Update ETCD
- Scheduler
- Kubelet

### Kube Controller-Manager
#### Controller
- Watch Status
- Remediate Situation
- Node Monitor Period = 5s
- Node Monitor Grace Period = 40s
- POD Eviction Timeout = 5m

### Kube Scheduler
- 알맞은 컨테이너를 알맞은 배에 실어야한다. / 배가 컨테이너를 충분히 수용할 수 있는지 확인필요
- 배마다 목적지가 다를 수 있다. / 컨테이너가 올바른 목적지로 갈 수 있도록 올바른 배에 실렸는지 확인필요

쿠버네티스에서 스케줄러는 어떤 노드를 결정한다.
스케줄러는 이 파드를 어떻게 할당하나
- 스케줄러는 각 파드를 보고, 그것에 가장 적합한 노드를 찾으려고 한다.

### kubelet
- 배의 선장

### YAML in Kubernetes
- **apiVersion**
    - 오브젝트 생성을 위한 Kubernetes API의 버전
    - Pod, Service : v1
    - ReplicaSet, Deployment : app/v1

- **kind**
    - 만들려고 하는 객체
    - Pod, Service, ReplicaSet, Deployment

- **metadata**
    - 데이터
    - 이름, 라벨(app, type) 등과 같은 객체 위에 있음

- **spec**
    - 만들려는 객체에 따라 추가 정보 제공
    - containers(name, image)

- 생성 명령어
    - kubectl create -f pod-definition.yml

- Pod 확인 명령어
    - kubectl get pods
- Pod 자세한 정보 확인
    - kubectl describe pod myapp-pod