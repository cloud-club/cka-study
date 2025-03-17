# Cluster Architecture

### Docker vs containerd
Kubernetes는 처음에 Docker만 지원
다양한 컨테이너 런타임을 지원하기 위해 CRI(Container Runtime Interface) 도입
- CRI는 OCI 표준을 따르는 모든 런타임과 호환됨

But, Docker는 OCI를 완전하게 준수 X → dockershim 도입해서 Kubernetes에서 Docker를 임시 지원

dockershim 유지 보수 어려움 → dockershim 제거
⇒ Kubernetes에서 Docker가 공식 지원 런타임에서 제외

하지만, Docker로 만든 이미지는 여전히 containerd를 통해 사용 가능

<br>

### ETCD
ETCD is a distrivuted reliavle **key-value store** that is Simple, Secure & Fast
Node, Pods, Configs, Secrets, Accounts, Roles, Bindings 와 같은 정보들 다 저장

<br>

### Kube-API Server
요청의 인증과 유효성을 확인하고, 데이터 스토어에서 데이터를 검색하고 업데이트
기타 데이터 저장소와 **직접 상호작용**하는 유일한 구성요소

kube-scheduler, kube-controller-manager 과 kubelet과 같은 
다른 구성요소는 Kube-API Server를 이용해 각 영역의 클러스터에서 업데이트를 수행

<br>

### Kube Controller Manager
지속적으로 모니터링, 시스템 전체를 원하는 기능 상태로 만듦
ex. node-controller, replication-contoller, …


<br>

### Kube Scheduler
Pod가 어느 Node에  들어갈지만 결정 → 실제로 놓는거 X (kubelet 역할)
best node for the pod → Filter & Rank Nodes


<br>

### Kubelet
Node의 모든 활동 지휘 ⇒ Register Node, Create Pods, Monitor Node & Pods

- always manually install the kubelet on your worker node


<br>

### Kube-proxy
- pod network: 내부 가상 네트워크로, 모든 pod가 연결되는 클러스터 내 모든 노드에 걸쳐있음 → 이 네트워크로 서로 통신 가능


쿠버네티스 클러스터의 각 노드에서 실행되는 프로세스
새 service 생성될 때마다 각 노드에 적절한 규칙을 만들어 그 service로 트래픽 전달

- 방법 → iptables rule




# Workloads
### Pod
어플리케이션의 single instance
Kubernetes에서 만들 수 있는 가장 작은 오브젝트

- **YAML in Kubernetes**
    - apiVersion
        - 오브젝트 생성을 위한 Kubernetes API의 버전
        - Pod, Service → v1
        - ReplicaSet, Deployment → apps/v1
    - kind : 오브젝트 타입
    - metadata
        - name → string형태
        - labels → dictionary 형태
    - spec
        - 생성하려는 것에 대한 추가 정보 제공
        - containers → List/Array

<br>

### Replicaset
유지해야할 Pod의 갯수 보장 → **고가용성**
Pod가 하나여도 가능; Pod 고장났을때 자동으로 새 Pod 교체

Load Balancing & Scaling

- Relication Controller → 구식
- Replica Set → replica 설정하는 새로운 권장 방법


- **Labels and Selectors**
    - ReplicaSet의 역할은 pod 감시
    - 어떤 Pod를 감시할지 설정하는게 Labeling

<br>

### Deployment
Production 환경에 배포, 매끄러운 업그레이드, 롤링 업데이트, 롤백, 환경별 구성 등 가능하도록 

Deployment > Replica Set > Pod




# Services
**communication between various components within and outside of the application**

클라이언트와 파드의 연결 담당
파드들이 외부의 트래픽을 받을 수 있도록 노출 시킬 수 O
YAML파일 - spec에서 type형태로 지정

<br>

### Nodeport
각 노드에서 <NodeIP>:<NodePort>를 통해 클러스터 외부에서 접근할 수 있게 해줌

- LoadBalancing → Random Algorithm

- 만약 Node가 여러개에 각각 Pod가 배포됐다면?
    - 자동적으로, Service 하나로 모든 노드 매핑 가능!

<br>

### ClusterIP
Cluster 안에 가상의 IP 생성해서 Pod끼리 Cluster 내부 통신
⇒ 하나의 인터페이스로 단체 Pod에 접속할 수 있음 !

<br>

### LoadBalancer
Cloud Provider가 제공하는 LB 사용




# Namespace
격리된 공간

- `default`: Kubernetes 처음 설치시에 자동적으로 생성
- `kube-system`: 네트워킹이나 DNS Service 등에 요구되는것
- `kube-public`: 모든 유저가 사용할 리소스인 경우 Kubernetes가 생성

namespace 각각 누가 뭘 할건지 policy를 가질 수 있음

namespace에 리소스 할당량 할당 가능

- 다른 namespace의 service에 연결하려면
    - `<service name>.<namespace>.svc.cluster.local`
    
    ⇒ 서비스 생성될 때 DNS가 자동으로 저 포맷으로 추가되기 때문에 가능




# Imperative vs Declarative
명령적 접근 vs 선언적 접근

### 명령적 접근
Specifying what to do and **how to do**

```bash
# Create Objects
kubectl run nginx --image=nginx
kubectl create deployment nginx --image=nginx
kubectl expose deployement nginx --port 80

# Update Objects
kubectl edit deployment nginx
kubectl scale deployment nginx --replicas=5
kubectl set image deployment nginx nginx=nginx:1.18

kubectl create -f nginx.yaml
kubectl replace -f nginx.yaml
kubectl delete -f nginx.yaml

```

- 한계
    - 고급 설정의 경우, 길고 복잡한 명령 필요 ex. 다중 컨테이너
    - 한번 실행하고 끝. 어떻게 만들어졌는지 알기 어려움
    
    ⇒ Configuration file(YAML)가 도움 O

<br>

### 선언적 접근
Just specifying the final destination (**what to do**)
declare our requirements → 시스템이 알아서 ~

파일로 예상 상태를 적용함 

```bash
# Create / Updates
kubectl apply -f nginx.yaml 
```

- **과정**
    - `kubectl apply`  → Local file
    - 라이브 개체가 status 필드 가진채로 존재 
    → Live Object configuration (in the Kubernetes memory)
    - JSON format으로 converted & stored as last applied configuration
        
        → stored on the live object configuration on the Kubernetes cluster itself 
        
        - Live Object configuration의 `metadata.annotations.kubectl.kubernetes.io/last-applied-configuration`에 JSON으로 들어감 

