# Manaual Scheduling

`spec.nodeName`이라는 필드 → Kubernetes가 자동으로 추가

Scheduler 모니터링해서 저 속성 없는것 확인하고, 스케줄링 알고리즘 실행하여 Pod의 Node를 할당 

if. Node를 모니터하고 스케줄링할 Scheduler가 없다면 Pod는 pending 상태

⇒ `spec.nodeName`에 설정하여 직접 Node에 Pod 할당 가능 (생성시에만 가능)

if. 이미 존재하는 Pod를 Node에 할당하는 경우, 

`Binding` 오브젝트 생성하고, Pod의 Binding API에 Post 요청

(실제 스케쥴러가 하는 일 모방)

POST 요청시, JSON 형식의 데이터로 

```bash
curl --header "Content-Type:application/json" --request POST --data
http://$SERVER/api/v1/namespaces/default/pods/$PODNAME/binding/
```

# Labels & Selectors

그룹으로 묶는 방법

- **Labels**
    - 각 물품에 부착된 속성
- **Selectors**
    - 위의 항목들을 필터링하는 것을 도와줌

in Kubernetes, 다른 오브젝트를 필터링하고 볼 방법이 필요함

`metadata.labels` → key: value 포맷으로 지정

```bash
kubectl get pods --selector app=App1
```

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
	name: myapp-replicaset
	labels:
		app: myapp # Replicaset의 Label
		type: front-end
spec:
	replicas: 2
	selector:
		matchLabels:
			app: myapp # Pod를 발견하기 위해 사용. 하나만 있어도 ok
	template: 
		metadata: 
			name: myapp-pod
			labels:
				app: myapp # Pod에서 구성된 Label
				type: front-end
			spec:
				containers:
				- name: nginx-container
					image: nginx
```

```yaml
apiVersion: v1
kind: Service
metadata: 
	name: my-service
spec:
	selector:
		app: myapp
	ports:
		- protocol: TCP
			port: 80
			targetPort: 9376
```

### Annotation

정보 수집 목적으로 다른 세부 사항을 기록하는데 사용

# Taint And Tolerations

Node에 Pod를 어떻게 배치할지, 어떻게 제한할지

Taint와 Toleration은 보안과 침입과 상관X

⇒ 한 노드에 어떤 포드로 스케쥴링할지 제한을 설정하기 위해 사용

- **Taint**: Node에 선택되어, 모든 Pod가 Node에 할당될 수 없도록 설정
- **Toleration**: Pod에 선택되어, Taint를 견딜 수 있도록 함.

⇒ Node에 할당할 Pod 제한

MaterNode는 어떤 Pod도 스케줄링하지 않음

⇒ why? MasterNode는 자동으로 Taint 설정이 되어있음. 

```bash
kubectl describe node kubemaster | grep Taint
```

### Taint

```bash
kubectl taint nodes <node-name> <key>=<value>:<taint-effect>
```

- `taint-effect`: taint에 tolerate 하지 않는 Pod에 대한 액션
    - `NoSchedule` : Pod가 Node에 스케줄되지 않음. 기존 실행되던 Pod는 유지됨 (퇴출 ❌)
    - `PreferNoSchedule` : Pod가 Node에 스케줄되지 않도록 최대한 노력하지만 보장X
    - `NoExecute` : 새로운 Pod는 해당 Node에 위치되지 않으며, 만약 해당 Node에 존재했던 Pod가 있다면 아래 조건에 따라 처리
        - Tolerate 되지 않은 Pod, 즉시 퇴출
        - Tolerate 된 Pod 이면서, toleration spec에 tolerationSeconds 가 지정되지 않으면, 영구적 바인딩
        - Tolerate 된 Pod 이면서, toleration spec에 tolerationSeconds 에 특정 시간이 지정되어 있으면, 해당 시간 이후 Node Lifecycle Controller에 의해 퇴출
            - tolerationSeconds: Pod가 Node에 바인딩된 시간을 지정

### Tolerations

```bash
kubectl taint nodes node1 app=blue:NoSchedule
```

```yaml
apiVersion:
kind: Pod
metadata:
	name: myapp-pod
spec:
	containers:
		- name: nginx-container
			image: nginx
	tolerations:
		- key: "app"
			operator: "Equal"
			value: "blue"
			effect: "NoSchedule"
```

⇒ 쌍따옴표(””)로 인코딩 되어있어야함

Tolerataion 설정된 Pod는 Taint된 Node말고, 다른 Node도 갈 수 있음

특정 Node에 한 Pod를 제한하려면

⇒ Node Affinity 를 통해 달성!

# Node Selector

기본 설정상 모든 Pod는 어떤 노드로도 갈 수 있음

⇒ 특정 노드에서만 작동하도록 Pod의 한계를 설정

- **Node Selector**
    - Pod가 특정 Node에서 실행하도록 제한하기 위해 `spec.nodeSelector` 속성 추가
    
    ```yaml
    ...
    spec:
    	...
    	nodeSelector:
    		size: Large # 노드에 할당된 Label
    ```
    

- Node에 Label 달기

```bash
kubectl label nodes <node-name> <lable-key>=<label-value>
```

단, 복잡한 요구사항인 경우, Node Selector로 하기 어려움

⇒ Node Affinity

# Node Affinity

Pod가 특정 Node에 호스트될 수 있도록!

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: myapp-pod
spec:
	containers:
		- name: data-processor
		image: data-processor
	affinity:
		nodeAffinity:
			requiredDuringSchedulingIgnoredDuringExecution:
				nodeSelectorTerms:
				- matchExpressions:
					- key: size
						operator: In
						values:
							- Large
							- Medium-
```

- Operators
    - https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#operators

if. Node Affinity 가 주어진 식의 Node와 일치하지 않는 경우?

if. 미래의 누가 Node의 Label을 바꾼 경우? 

⇒ Node Affinity Types로 결정

- DuringScheduling: Scheduling에는 Pod가 존재하지 않다가 처음으로 만들어지는 상태
- DuringExecution: 실행 중 환경이 변하면서 Affinity에 영향 미친 경우 ~
- **Available**
    - `requiredDuringSchedulingIgnoredDuringExecution`
        - 지정된 Affinity 찾지 못하면 Scheduling 못함 (Required)
        - 실행 중 Node Affinity의 변화가 생겨도 영향을 미치지 않음 (Ignored)
    - `preferredDuringSchedulingIgnoredDuringExecution`
        - 지정된 Affinity 찾지 못해도 해당 Pod를 모든 가능한 Node에 배치 (Preferred)
        - 실행 중 Node Affinity의 변화가 생겨도 영향을 미치지 않음 (Ignored)
- **Planned**
    - `requiredDuringSchedulingRequiredDuringExecution`
        - 지정된 Affinity 찾지 못하면 Scheduling 못함 (Required)
        - 실행 중 Node Affinity에 부합하지 않은 노드에서 실행중인 Pod 쫓아냄 (Required)

![image.png](3%20Scheduling%201b8766b7e9a4800094decc0b733dfe0c/image.png)

### Node Affinity vs Taints and Tolerations

여러개의 Node, 여러개의 Pod인 상황에서

특정 Node에 특정 Pod를 할당시키려면 Taint/Tolerations 와 Node Affinity를 함께 사용!

# Resource Requirements

https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/

Scheduler는 Pod에서 요구하는 리소스 양과 Node의 리소스 양 고려

`spec.resources` 속성 추가하여 메모리와 CPU 사용에 대한 값을 지정

```yaml
...
spec:
	containers:
	...
	resources:
		requests:
			memory: "4Gi"
			cpu: 2
		limits: # 제한
			memory: "2Gi"
			cpu: 2
```

Scheduler가 이 Pod 배치해달라는 요청 받으면, 

저정도의 리소스 사용 가능한 Node를 찾음

→ Pod가 Node에 놓이면 Pod는 사용 가능한 양의 자원을 보장받음

- CPU
    - 1 CPU = 100m = 1 AWS vCPU = 1 GCP Core …

- Memory
    - 256Mi = 268435456
    - G → Gigabyte …

if, Pod가 지정된 한도의 자원을 초과한다면?

- CPU ⇒ 시스템이 CPU를 조절해 지정된 한도를 넘지 않도록 함
- Memory ⇒ 더 많은 메모리 리소스를 쓸 수 있음. → 종료됨. OOM(Out Of Memory)

### CPU 요청과 제한 작업

- No Requests/No Limits → 모든 CPU 리소스를 한 Pod가 독점할 수 있음
- No Requests/Limits
- Requests/Limits → 불필요한 CPU 제한
- Requests/No Limits → 이상적임. 사용 가능한 리소스 소모하고 정의되지 않은 포드는 제한 가능. 단, 모든 노드에 대한 요청 설정!

### Memory 요청과 제한 작업

- No Requests/No Limits → 모든 Memory 리소스를 한 Pod가 독점할 수 있음
- No Requests/Limits
- Requests/Limits → 불필요한 Memory 제한
- Requests/No Limits → CPU와 달리 Memory는 조절할 수 없으므로 Memory가 부족하면 Pod를 죽이고 Memory 할당 받아야함

### Pod 기본 설정

```yaml
apiVersion: v1
kind: LimitRange
metadata:
	name: cpu-resource-constraint
spec:
	limits:
		- default: # limit
				cpu: 500m
			defaultRequest: # request
				cpu: 500m
			max: # limit
				cpu: "1"
			min:
				cpu: 100m
			type: Container
```

```yaml
apiVersion: v1
kind: LimitRange
metadata:
	name: memory-resource-constraint
spec:
	limits:
		- default: # limit
				memory: 1Gi
			defaultRequest: # request
				memory: 1Gi
			max: # limit
				memory: 1Gi
			min:
				memory: 500Mi
			type: Container
```

제한 범위가 생기거나 업데이트된 후에 생긴 Pod에만 영향 미침!!

### Resource Quotas

Namespace 레벨에서 Resource Quotas 생성 가능 → 리소스 제한

```yaml
apiVersion: v1
kind: ResourceQuota
metadata: 
	name: my-resource-quota
spec:
	hard:
		requests.cpu: 4
		requests.memory: 4Gi
		limits.cpu: 10
		limits.memory: 10Gi
```

https://nice-engineer.tistory.com/entry/Kubernetes-Resource-Request-Limit

# Daemon Sets

여러개의 인스턴스 Pod를 배포하도록 도와줌

But, Node 마다 Pod를 하나씩 실행

클러스터에 새 Node 추가될때마다 Pod가 자동으로 Node에 추가

Node가 제거되면, Pod는 자동으로 제거

- 사용 사례
    - Monitoring Agent, Log Collector,
    - kube-proxy
    - …

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
	name: monitoring-daemon
spec:
	selector:
		matchLabels:
			app: monitoring-agent
	template: 
		metadata: 
			labels:
				app: monitoring-agent
		spec:
			containers:
			- name: monitoring-agent
				image: monitoring-agent
```

⇒ ReplicaSet과 형상이 비슷함

```bash
kubectl create -f daemon-set-definition.yaml

k get damonsets
k describe daemonsets monitoring-daemon
```

- 작동 방식
    - version 1.12까지는 `spec.nodeName` 속성을 사용
    - 이후에는 default Scheduler와 Node Affinity 사용

# Static Pods

(( MasterNode가 없는 경우 ))

kubelet을 설정하면, Pod에 관한 정보를 저장하는 서버 디렉토리를 통해 Pod 정의 파일을 읽을 수 있다.

해당 디렉토리에 Pod 정의 파일을 넣으면 kubelet은 주기적으로 해당 디렉토리에 파일을 확인하고 호스트에 Pod를 만듦

Pod가 죽지 않도록 보장.

⇒ **Static Pods**: API 서버 간섭이나 K8s 구성 요소 간섭 없이 kubelet이 스스로 만든 Pod

ReplicaSet, Deployment는 생성 X

 Static Pod를 생성하면 자동으로 Pod이름 뒤에 `-{nodename}` 이 붙는다

- kubelet.service

```bash
systemctl cat kubelet
```

```bash
ExecStart=/usr/local/bin/kubelet \\
--container-runtime=remote \\
--container-runtime-endpoint=unix:///var/run/containerd/containerd.sock \\
--pod-manifest-path=/etc/kubernetes/manifests \\
--kubeconfig=/var/lib/kubelet/kubeconfig \\
--network-plugin=cni \\
--register-node=true \\
--v=2
```

```bash
ExecStart=/usr/local/bin/kubelet \\
--container-runtime=remote \\
--container-runtime-endpoint=unix:///var/run/containerd/containerd.sock \\
--config=kubeconfig.yaml \\
--kubeconfig=/var/lib/kubelet/kubeconfig \\
--network-plugin=cni \\
--register-node=true \\
--v=2
```

- kubeconfig.yaml

```bash
staticPodPath: /etc/kubernetes/manifests
```

⇒ Kubelet 파일에서 확인

- Static Pod 생성
    - 해당 staticPodPath에 manifest 작성
    - `systemctl restart kubelet`

```bash
docker ps
```

kubernetes 무리 확보하지 못함.

kubectl 은 kube-api server와 작동

⇒ docker 명령어 사용.

### kubelet의 Pod 생성

1. Static Pod 폴더에서 Pod 정의 파일을 통해 Pod 생성 가능
2. kube-apiserver의 HTTP API endpoint를 통해 kubelet에 입력하여 Pod 생성

if, kube-apiserver가 kubelet이 만든 Static Pod를 인지?

⇒ O. 읽기 전용이라 수정이나 삭제는 불가능.

왜 쓰는가?

Static Pod를 통해 구성 요소를 Node의 Pod처럼 배포 가능

ex. MasterNode의 apiserver.yaml을 static pod로 해서 관리 편하게 ~

```bash
kubectl get pods -n kube-system 
# 컨트롤 플레인 구성 요소를 Pod로 볼 수 있음
```

# Multiple Schedulers

Schedulers 마다 이름 달라야함

- default-scheduler
    
    ```yaml
    # default-scheduler-config.yaml
    apiVersion: kubescheduler.config.k8s.io/v1
    kind: KubeSchedulerConfiguration
    profiles:
    	- schedulerName: default-scheduler
    ```
    
- my-scheduler
    
    ```yaml
    # my-scheduler-config.yaml
    apiVersion: kubescheduler.config.k8s.io/v1
    kind: KubeSchedulerConfiguration
    profiles:
    	- schedulerName: my-scheduler
    ```
    

### Deploy Additional Scheduler

kube-scheduler 바이너리를 다운로드해 여러 옵션과 함께 실행

```bash
wget https://storage.googleapis.com/kubernetes-release/release/v1.12.0/bin/linux/amd64/kube-scheduler
```

```bash
ExecStart=/usr/local/bin/kube-scheduler \\
--config=/etc/kubernetes/config/my-scheduler-config.yaml \\
```

 

⇒ 대부분의 경우 이렇게 배포하진 않음.

kubeadm 배포에서는 모든 ControlPlane 구성요소가 쿠버네티스 클러스터 안에서 Pod 혹은 Deploy로 실행됨

### Deploy Additional Schedule as a Pod

```yaml
# my-scheduler-config.yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
profiles:
	- schedulerName: my-scheduler
leaderElection:
	leaderElect: true
	resourceNamespace: kube-system
	resourceName: lo
```

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: kube-scheduler
	namespace: kube-system
spec:
	containers:
		- command:
			- kube-scheduler
			- --address=127.0.0.1
			- --kubeconfig=/etc/kubernetes/my-scheduler-config.conf
		image: k8s.gcr.io/kube-
```

- `leaderElection` 옵션
    - 동일한 스케줄러 복사본 여러개가 다른 노드에서 실행될 경우 한번에 하나만 활성화.
    - 리더 선출 옵션은 스케줄 활동을 이끌 리더를 선택

https://kubernetes.io/docs/tasks/extend-kubernetes/configure-multiple-schedulers/

```yaml
kubectl get pods -n kube-system
```

- 생성한 scheduler 어떻게 사용?
    - `spec.schedulerName` 속성 사용

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: nginx
spec:
	continaers:
		- image: nginx
			name: nginx
	schedulerName: my-scheduler
```

```bash
k create -f pod-definition.yaml
```

if, Scheduler가 제대로 설정되지 않으면, Pod는 계속 Pending 상태로 남음

```bash
kubectl get events -o wide
```

⇒ 어떤 Scheduler가 스케줄링 했는지 확인 가능

```bash
k logs my-scheduler -n kube-system
```

# Scheduler Profile

1. **Scheduling Queue**
    
    Pod 생성 → Pod Scheduling
    
    우선순위 높은 Pod → 먼저 대기자 명단
    
    - extension point: queueSort
    - Plugin: PrioritySort
2. **Filtering** 
    
    resource 등을 근거로 Filtering
    
    - extension point: prefilter, filter, postFilter
    - Plugin: NodeResourcesFit, NodeName, NodeUnschedulable
3. **Scoring** 
    
    남은 공간 근거로 Scheduler는 각 노드에 점수 매김 
    
    - extension point: preScore, score, reserve
    - Plugin: NodeResourcesFit, ImageLocality
4. **Binding**
    
    높은 점수 받은 Node 선택
    
    - extension point: premit, preBind, bind, postBind
    - Plugin: DefaultBinder

⇒ 모든 작업은 특정 플러그인으로 이루어짐

- default-scheduler
    
    ```yaml
    # default-scheduler-config.yaml
    apiVersion: kubescheduler.config.k8s.io/v1
    kind: KubeSchedulerConfiguration
    profiles:
    	- schedulerName: default-scheduler
    ```
    
- my-scheduler
    
    ```yaml
    # my-scheduler-config.yaml
    apiVersion: kubescheduler.config.k8s.io/v1
    kind: KubeSchedulerConfiguration
    profiles:
    	- schedulerName: my-scheduler
    ```
    

각각 다른 프로세스이기에 유지하기에 노력 필요.

⇒ 구성 파일에서 단일 Scheduler  내 다수 profile 구성 가능

```bash
# my-scheduler-config.yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
profiles:
	- schedulerName: my-scheduler
		plugins:
			score:
				disabled:
					- name: TaintToleration
					enabled:
					- name: MyCustomPluginA
	- schedulerName: my-scheduler-2
		plugins:
			preScore:
				disabled:
					- name: '*'
			score:
				disabled:
					- name: '*'
```

https://kubernetes.io/blog/2017/03/advanced-scheduling-in-kubernetes/

https://stackoverflow.com/questions/28857993/how-does-kubernetes-scheduler-work

[https://github.com/kubernetes/community/blob/master/contributors/devel/sig-scheduling/scheduling_code_hierarchy_overview.md](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-scheduling/scheduling_code_hierarchy_overview.md)

# Admission Controller

**kubectl → kube-apiserver → Etcd → Create Pod**

( Authentication → Authorization → Admission Controllers)

- kube-apiserver에 도달하면 Authentication 프로세스 → 일반적으로 인증서를 통해 이루어짐
    - kubectl의 경우 ~/.kube/config 에 인증서 구성
- Authentication → Authorization 프로세스
    - Role Based 제어를 통해 달성
- Authorization → Admission Controllers

### RBAC

Role Based Access Control

특정 역할이 있는 사용자에게 다양한 종류의 제한을 설정할 수 있음

API 수준에서 어떤 사용자에게 어떤 종류의 API 작업에 대한 액세스가 허용되는지 결정

but, 사용자가 객체에 대해 어떤 종류의 액세스 권한을 갖는지 정의하는 것 이상을 수행하려면 ?

⇒ Admission Controllers

### Admission Controllers

클러스터 사용 방식을 강화하기 위해 더 나은 보안 조치를 구현

구성의 유효성을 검사하는것 이상으로, 

Pod 생성 전에 요청을 변경하거나, 추가 작업을 수행하는것 제안하는 등

많은 작업 수행 가능

ex. AlwaysPullImages, DefaultStroageClass, EventRateLimit, NamespaceExists, NamespaceAutoProvision, NamespaceLifecycle, …

```bash
kube-apiserver -h | grep enable-admission-plugins
```

```bash
k exec kube-apiserver-controlplane -n kube-system -- kube-apiserver -h | grep enable-admission-plugins 
```

```bash
systemctl cat kube-apiserver.service
```

```bash
cat /etc/kubernetes/manifests/kube-apiserver.yaml
```

# Validating & Mutating Admission Controllers

- DefaultStorageClass
    - 기본적으로 활성화되어 있는 플러그인
    - PVC에서 Stroage Class 설정하지 않아도 기본값으로 default 들어감

⇒ Mutating Admission Controllers

객체가 생성되기 전에 객체 자체를 변경하거나 변형 가능

### Mutating Admission Controllers

요청을 변경할 수 있는 권한 컨트롤러

### Validating Admission Controllers

요청의 유효성을 검사하여 허용/거부 할 수 있는 권한 컨트롤러

일반적으로 Mutating → Validating

ex. NamespaceAutoProvision → NamespaceExists

### Webhook

MutatingAdmissionWebhook 

ValidatingAdmissionWebhook

⇒ 이러한 웹후크가 K8s 클러스터 내부 또는 외부에서 호스팅되는 서버를 가리키도록 구성할 수 있고, 서버는 자체 코드와 로직으로 Own Admission Controller 실행

요청이 모든 기본 Admission Controller를 통과한 후, 구성된 Webhook에 연결

Webhook에 도달하면 JSON 형식의 AdmissionReview 객체를 전달하여 Admission Webhook server 호출

요청 받으면 Admission Webhook server는 요청이 허용되는지 여부에 대한 결과 포함된 AdmissionReview 객체로 응답

`allowed` 필드가 true → 허용

`allowed` 필드가 false → 거부

어떻게 설정?

1. 자체 로직이 있는 Admission Webhook server을 배포
    1. 변경 및 유효성 검사 API를 수락하고 웹 서버가 예상하는 JSON 객체로 응답
    2. 시험에서는 서버를 작성할 필요 X
2. Webhook 구성 오브젝트를 배포(ex. deployment), Service 구성
3. Service 연결하고 요청의 유효성 검사하거나 변경하도록 클러스터 구성
    
    → 개체 생성
    
    ```yaml
    apiVersion: admissionregistration.k8s.io/v1
    kind: ValidatingWebhookConfiguration
    metadata:
    	name: "pod-policy.example.com"
    webhooks:
    	- name: "pod-policy.example.com"
    		clientConfig:
    			# url: "https://external-server.example.com"
    			service:
    				namespace: "webhook-namespace"
    				name: "webhook-service"
    			caBundle: "Ci0tL...." # 인증서 번들
    		rules:
    			- apiGroups: [""]
    				apiVersions: ["v1"]
    				operations: ["CREATE"]
    				resources: ["pods"]
    				scope: "Namespaced"
    ```
