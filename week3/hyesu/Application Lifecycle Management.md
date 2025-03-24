# Updates and Rollback Deployment

- 롤아웃과 버전관리
    - 처음 배포 생성 → 롤아웃 트리거 → 배포 revision 생성
    - 앱 업그레이드 → 롤아웃 트리거 → 새 배포 revision 생성
    
    ⇒ 배포에 일어난 변화 추적 & 이전 버전으로 돌릴 수 O
    

```bash
kubectl rollout status deployment/myapp-deployment

# see rivision and history rivision
kubectl rollout history deployment/myapp-deployment
```

### Deploy Strategy

- **Recreate Strategy**
    - 한꺼번에 모두 파괴하고 새 인스턴스 배포
    - 문제
        - 구 버전 다운되고 새 버전 업되기 전 기간에 사용자 접근 불가능
- **Rolling Update(default)**
    - 하나씩 구버전 내리고, 새버전 올림
    - 앱 다운 X, 업그레이드 원활

⇒ `kubectl describe`를 통해 StrategyType 확인 가능!

### Upgrades

Deployment → new Replica Set

예전 버전의 Replica Set Pod 제거하면서 하나씩 new Replica Set의 Pod 띄움

⇒ Rolling Update Stragegy

### Rollback

Deployment는 이전 revision으로 롤백 가능

```bash
kubectl rollout undo deployment/my-deployment
```

# Application Commands

CKA 자격증에 필수 주제는 아님

### Docker

```bash
docker run ubuntu
docker ps 
# 컨테이너가 안보임. 왜?
```

⇒ VM과 달리 컨테이너는 운영체제 자체를 호스팅 X

특정 작업이나 프로세스를 실행.

작업이 끝나면 컨테이너 빠져나옴

⇒ 컨테이너는 그 안의 과정이 살아있어야만 살 수 있음

if. 컨테이너 안의 웹 서비스 멈추거나 충돌하면 컨테이너 나가게 됨

컨테이너 내 실행되는 프로세스는 누가 정의?

⇒ Dockerfile의 CMD 커맨드

- ubuntu 이미지의 경우 CMD [”bash”]
    - bash는 terminal의 input 듣는 shell
    - terminal 못 찾으면 탈출
- 기본값으로 Docker는 실행 중일때 컨테이너에 터미널 연결 X

컨테이너 시작 명령을 어떻게 지정?

⇒ 이미지에 지정된 기본 명령 재정의

```bash
docker run ubuntu [COMMAND]
```

if. 변화를 지속하고 싶을 경우?

⇒ Ubuntu 이미지 가져와서 새 명령 명시

```docker
FROM ubuntu

# CMD command param1 # 쉘 양식
CMD ["command", "param1"] # JSON 포맷
```

```docker
docker build -t ubuntu-sleeper .

docker run ubuntu-sleeper
```

if. 잠자는 시간 변경하고 싶다면?

```docker
docker run ubuntu-sleeper 10
```

⇒ ENTRYPOINT 명령 필요

```docker
FROM ubuntu

ENTRYPOINT ["sleep"]
```

if. params 주어지지 않은 경우?

⇒ 오류 발생. 기본값 설정 필요

```docker
FROM ubuntu

ENTRYPOINT ["sleep"]

CMD ["5"]
```

if. 런타임동안 진입점 수정하고 싶은 경우?

```bash
docker run --name ubuntu-sleeper \
		--entrypoint sleep2.0 
		ubuntu-sleeper 10
```

### Kubernetes

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: ubuntu-sleeper-pod
spec:
	containers:
		- name: ubuntu-sleeper
			image: ubuntu-sleeper
			command: ["sleep2.0"] # docker의 ENTRYPOINT와 일치
			args: ["10"] # docker의 CMD와 일치
```

# Configure Applications

### ENV Variables in Kubernetes

- Plain Key Value

```yaml
...
spec:
	containers:
			...
			env:
				- name: APP_COLOR
					value: pink
```

- ConfigMaps

```yaml
...
spec:
	containers:
			...
			env:
				- name: APP_COLOR
					valueFrom:
						configMapKeyRef:
```

- Secrets

```yaml
...
spec:
	containers:
			...
			env:
				- name: APP_COLOR
					valueFrom:
						secretKeyRef: 
```

### ConfigMaps

Kubernetes의 key-value 쌍의 구성 데이터 전달하는데 사용

Pod 생성 → Pod에 Configmap 삽입해 key-value 쌍이 환경 변수로 사용될 수 있도록 함

- 단계
    - configmap 생성
    - pod에 주입

- **Create ConfigMap**
    - 명령적 접근
    
    ```bash
    kubectl create configmap \
    	<config-name> --from-literal=<key1>=<value1> \
    								--from-literal=<key2>=<value2>
    # --from-literal: key-value 쌍 지정하는데 사용
    ```
    
    ```bash
    kubectl create configmap \
    	<config-name> --from-file=<path-to-file>
    ```
    
    - 선언적 접근
    
    ```bash
    kubectl create configmap -f configmap.yaml
    ```
    
    ```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
    	name: app-config
    data:
    	APP_COLOR: blue
    	APP_MODE: prod
    ```
    

- View ConfigMaps

```bash
kubectl get configmaps

kubectl describe configmaps
```

- **ConfigMap in Pods**
    - 전체 ENV 넣는 방법
        
        ```yaml
        ...
        spec:
        	containers:
        			...
        			envFrom:
        				- configMapRef:
        						name: app-config
        ```
        
    - Single ENV 넣는 방법
        
        ```yaml
        ...
        spec:
        	containers:
        			...
        			env:
        				- name: APP_COLOR
        					valueFrom:
        						configMapKeyRef:
        							name: app-config
        							key: APP_COLOR
        ```
        
    - Volume
        
        ```yaml
        volumes:
        - name: app-config-volume
        	configMap:
        		name: app-config
        ```
        

### Secrets

민감한 정보 저장하는데 사용

base64로 인코딩된 형식으로 저장

```bash
# 인코딩 방법
echo -n 'mysql' | base64

# 디코딩 방법
echo -n 'bXlzcWw==' | base64 --decode
```

- 단계
    - Secret 생성
    - Pod에 주입

- **Create Secret**
    - 명령적 방식
        
        ```bash
        kubectl create secret generic \
        	<secret-name> --from-literal=<key>=<value> \
        								--from-literal=<key>=<value>
        ```
        
        ```bash
        kubectl create secret generic \
        	<secret-name> --from-file=<path-to-file>
        ```
        
    - 선언적 방식
        
        ```bash
        kubectl create secret -f secret.yaml
        ```
        
        ```yaml
        apiVersion: v1
        kind: Secret
        metadata:
        name: app-secret
        data: 
        	DB_HOST: bXlzcWw== # base64 값으로 인코딩된 값을 사용해야 함
        ```
        
        ```yaml
        apiVersion: v1
        kind: Secret
        metadata:
        name: app-secret
        stringData: 
        	DB_HOST: mysql
        ```
        

- View Secrets

```bash
kubectl get secrets

kubectl describe secrets
```

- **Secrets in Pods**
    - 전체 ENV 넣는 방법
        
        ```yaml
        ...
        spec:
        	containers:
        			...
        			envFrom:
        				- secretRef:
        						name: app-secret
        ```
        
    - Single ENV 넣는 방법
        
        ```yaml
        ...
        spec:
        	containers:
        			...
        			env:
        				- name: DB_HOST
        					valueFrom:
        						secretRef:
        							name: app-secret
        							key: DB_HOST
        ```
        
    - Volume
        
        ```yaml
        volumes:
        - name: app-secret-volume
        	secret:
        		name: app-secret
        ```
        
        ```bash
        ls /opt/app-secret-volues
        ```
        
    
- **Note on Secrets**
    - Secret은 암호화 X. only encoded
    - Secret은 ETCD에서 암호화 X
        - https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/
        - EncryptionConfiguration 사용해서 Secret 암호화
        - `resourves.providers`의 순서가 암호화에 사용됨을 의미!
        - 암호화 이후 새로 만드는것에만 적용! → `k replace -f -`
    - 같은 namespace에서 create pods/deployments 할 수 있는 사람은 secret에 접근 가능 → RBAC 구성
    - third-party Secret store provider 고려
        - ex. AWS Provider, Azure Provider, GCP Provider, Vault Provider

- 참고
    - Secrets를 안전하게 관리하는 방법
        - Secret 객체 정의 파일을 소스 코드 저장소에 체크인하지 않기
        - ETCD의 Encryption at Rest(저장 시 암호화)를 활성화하여, Secrets을 암호화하여 저장
    - 쿠버네티스의 Secret 처리 방식
        - Secret은 해당 노드의 Pod에서 필요할 때만 노드로 전송됨
        - Kubelet은 Secret을 tmpfs(메모리 기반 파일 시스템)에 저장하여 디스크에 기록되지 않도록 함
        - Pod가 삭제되면, Kubelet은 해당 Secret의 로컬 복사본을 삭제함

https://www.youtube.com/watch?v=MTnQW9MxnRI

# Multi-Container Pods

같은 LifeCycle, Network, Storage 공유 → 서로를 localhost로 부를 수 있음

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: simple-webapp
	labels:
		name: simple-webapp
spec:
	container: # 배열
	- name: simple-webapp
		image: simple-webapp
		ports:
			- containerPort: 8080
	- name: log-agent
		image: log-agent
```

- Design Patterns
    - SideCar
    - Adapter
    - Ambassador
    
    ⇒ CKAD 범위
    

# Init Containers

한 번만 실행되면 되는 작업을 실행해야 하는 경우 사용

ex. 애플리케이션이 사용할 코드나 바이너리를 저장소에서 가져오는 작업, 실제 애플리케이션이 시작되기 전에 외부 서비스나 데이터베이스가 준비될 때까지 대기하는 작업

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
spec:
  containers:
  - name: myapp-container
    image: busybox:1.28
    command: ['sh', '-c', 'echo The app is running! && sleep 3600']
  initContainers:
  - name: init-myservice
    image: busybox
    command: ['sh', '-c', 'git clone <some-repository-that-will-be-used-by-application> ; done;']
```

POD가 생성될 때 먼저 InitContainer가 실행

InitContainer가 완료될 때까지 Kubernetes는 애플리케이션 컨테이너를 실행하지 X

- **여러 개의 InitContainer 설정**

InitContainer는 여러 개 설정할 수도 있으며, 이 경우 순차적으로 하나씩 실행

모든 InitContainer가 실행을 완료해야 본 컨테이너가 실행

- **InitContainer의 실패**

만약 InitContainer 실행이 실패하면, 

Kubernetes는 POD를 계속 재시작하며 InitContainer가 정상적으로 완료될 때까지 반복

https://kubernetes.io/docs/concepts/workloads/pods/init-containers/

# Self Healing Applications

Kubernetes는 ReplicaSet과 Replication Controller를 통해 Self-Healing 지원

Replication Controller는 POD 내의 애플리케이션이 충돌(Crash)하면 자동으로 POD를 재생성하여, 항상 일정한 개수의 복제본(Replicas)이 실행되도록 보장

Kubernetes는 Liveness Probe 및 Readiness Probe를 제공 → POD 내부 애플리케이션의 상태를 확인하고, 필요할 경우 조치를 취할 수 있도록 지원

⇒ CKAD 범위

# AutoScaling in Kubernetes

- Scaling
    - HPA
        - 더 많은 인스턴스 또는 서버를 추가
    - VPA
        - 더 많은 리소스 추가

### Scaling Cluster Infra

- **Horizontal** → 클러스터에 노드를 더 추가
    - Manual → 새 노드를 수동으로 프로비저닝하고 `kubeadm join`으로 클러스터에 새 노드 추가
    - Automated → Cluster Autoscaler
- **Vertical** → 기존  노드의 리소스 늘림
    - 실행하는 애플리케이션 중단하고 다시 올려야해서 Kubernetes에서 일반적으로 사용하는 접근 방식은 아님

### Scaling Workloads

- **Horizontal** → 더 많은 Pod 생성
    - Manual → `kubectl scale`
    - Automated → Horizontal Pod Autoscaler(HPA)
- **Vertical** → 재시작해서 기존 Pod에 할당된 리소스 늘림
    - Manual → `kubectl edit` 으로 리소스 제한이나 요청 변경
    - Automated → Vertical Pod Autoscaler(VPA)

|  | **Scaling Cluster Infra** |  | **Scaling Workload** |  |
| --- | --- | --- | --- | --- |
|  | Horizontal Scaling | Vertical Scaling | Horizontal Scaling | Vertical Scaling |
| **Manual 
수동** | `kubeadm join` | X | `kubectl scale` | `kubectl edit` |
| **Automated
자동** | Cluster Autoscaler | X | HPA | VPA |

### HPA

Horizontal Pod Autoscaler

- 수동으로 수평 확장하는 법
    
    ```bash
    # Cluster에서 메트릭 서버가 실행 중인 경우
    kubectl top pod my-app-pod 
    
    kubectl scale deployment my-app --replicas=3
    ```
    
    ⇒ 지속적으로 모니터링 해야하기에 빠르게 대응 못할수도 있음
    
    ⇒ HPA 사용
    
    if. 클러스터 리소스 부족한 경우 → 일부 Pod만 실행되고, **나머지는 Pending 상태**
    
- **HPA**
    - top 명령어를 통해 메트릭을 지속적으로 모니터링
    - CPU, Memory, Custom metric에 따라 Deploy, Replicaset, Statefulset Pod수 자동으로 늘리거나 줄임
    - CPU 메모리 또는 Memory 사용량이 많아지면 HPA는 더 많은 Pod 생성
    - 리소스 절약하기 위해 여분의 Pod 제거 → Balance
    - 여러 metric 추적 가능

```bash
kubectl autoscale deployment my-app \
	--cpu-percent=50 --min=1 --max=10
```

```bash
kubectl get hpa

kubectl delete hpa my-app
```

- 선언적 방법

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
	name: my-app-hpa
spec:
	scaleTargetRef:
		apiVersion: apps/v1
		kind: Deployment
		name: my-app
	minReplicas: 1
	maxReplicas: 10
	metrics:
		- type: Resource
			resource:
				name: cpu
				target:
					type: Utilization
					averageUtilization: 50
```

⇒ Kubernetes 1.23부터 built-in으로 사용 가능

내부 metric-server도 가능하지만, 

external adapter에서 metric 참조 가능

### In-place Resize of Pods

제자리 크기 조정

Kubernetes 1.32 버전부터 

Deployment에서 리소스 요구 사항을 변경하는 경우, 

- 기본 동작: 기존 Pod 삭제하고 새 변경 사항 포함된 새 Pod 스핀업
    
    ⇒ statefuleset에서 문제 발생할 여지 O
    
    ⇒ In-place update pod resource라는 개선 작업 진행중
    

https://kubernetes.io/docs/concepts/workloads/autoscaling/#requirements-for-in-place-resizing

```bash
FEATURE_GATES=InPlacePodVerticalScaling=true
```

```yaml
...
kind: Deployment
spec:
	...
	template:
		...
		spec:
			containers:
			- name: my-app
				...
				resizePolicy:
					- resourceName: cou
						restartPolicy: NotRequired
						resourceName: memory
						restartPolicy: RestartContainer
```

- **Limitations**
    - CPU와 Memory 리소스에만 적용됨
    - Pod QoS 클래스 변경 X
    - In-place Resize 방식으로는 컨테이너와 임시 컨테이너의 크기 조정 불가능
    - resource 요청과 제한은 한번 설정하면 변경 X
    - 컨테이너 메모리 제한은 사용량 이하로 줄일 수 X
        - 만약, 해당 상태에 놓이면 원하는 메모리 제한이 가능해질때까지 크기 조정 상태가 계속 진행 중으로 유지
    - Windows Pod는 resized 불가능

### VPA

- 수동으로 수직 확장하는 법
    
    ```bash
    # Cluster에서 메트릭 서버가 실행 중인 경우
    kubectl top pod my-app-pod 
    
    kubectl edit deploy my-app
    ```
    
- **VPA**
    - Metric 지속적으로 모니터링
    - Deployment에서 Pod에 할당된 리소스를 자동으로 늘리거나 줄여서 균형 맞춤

자동으로 제공되지 않음

⇒ GitHub repo에 있는 VPA 정의 파일 apply

```bash
k get pods -n kube-system | grep vpa
# vpa-admission-controller-xxxxx
# vpa-recommender-xxxxx
# vpa-updater-xxxxx
```

- **구성 요소**
    - `VPA-Recommender`
        - kubernetes metric API에서 리소스 사용량을 지속적으로 모니터링
        - Pod의 과거 및 실시간 사용량 데이터를 수집
        - 최적의 CPU 및 메모리 값에 대한 권장 사항 제공
        
        ⇒ Pod를 직접 수정하지 않고 변경 사항만 제안.
        
    - `VPA-Updater`
        - 최적이 아닌 리소스로 실행중인 파드 감지
        - 업데이트가 필요할때 해당 파드 종료
        - VPA-Recommender로부터 정보 가져와서 Pod 모니터링
    - `VPA Admission Controller`
        - Pod 생성 프로세스에 개입
        - VPA-Recommender의 권장 사항 다시 사용하여 시작시 권장 CPU 및 메모리 값을 적용하도록 파드 사양 변경
            
            ⇒ 새로 생성된 파드가 올바른 리소스 요청으로 시작
            

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
	name: my-app-vpa
spec:
	targetRef: # 모니터링 대상
		apiVersion: apps/v1
		kind: Deployment
		name: my-app
	updatePolicy:
		updateMode: "Auto"
	resourcePolicy:
		containerPolicies:
		- containerName: "my-app"
			minAllowed:
				cpu: "250m"
			maxAllowed:
				cpu: "2"
			controlledResources: ["cpu"]
```

- **VPA Mode**

| Mode | 설명 |
| --- | --- |
| Off | 변경 사항만 recommends. 아무 작업도 수행 X → `VPA-Recommender`만 작동 |
| Initial | 생성 시에만 파드 변경. `VPA-Recommender`이 변경 권장.  |
| Recreate | 리소스 소비가 지정된 범위 초과하면 `VPA-Updater` 가 개입해서 기존 Pod 종료하고 재생성 |
| Auto | 기존 Pod를 권장 수로 업데이트 
(Recreate mode와 유사 but, In-place update of Pod Resources 지원할 수 있는 경우에는 Auto 모드 선호) |

```bash
kubectl describe vpa my-app-vpa
```

### Key Differences VPA, HPA

| Feature | HPA | VPA |
| --- | --- | --- |
| **Scaling Method** | 부하에 기반해서 Pod 생성/삭제 | 개별 Pod의 CPU & Memory 증가 |
| **Pod Behavior** | Pod running 유지 | 새로운 리소스로 Pod Restart |
| **트래픽 폭증 처리** | ✅ 즉시 Pod 추가 | ❌ VPA는 Pod 재시작 필요 |
| **비용 최적화** | ✅ 필요없는 낭비 Pod 피함 | ✅ over-provisioning된 CPU/Memory 방지 |
| **Best For** | Web apps, microservices, stateless services | Statefule workloads, CPU/Memory-heavy apps(DBs, ML workloads) |
| Example Use Cases | web servers (Nginx, API services), message queues, microservices | Databases(MySQL, PostgreSQL), JVM-based apps, AI/MK workloads |
