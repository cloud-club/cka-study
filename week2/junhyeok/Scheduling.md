# Manual Scheduling

노드 이름은 생성 시에만 지정가능

- pod-definiton.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: nginx
	labels:
		name: nginx
spec:
	containers:
	- name: nginx
		image: nginx
		ports:
			- containerPort: 8080
	nodeName: node02
```

이미 있는 포드에 노드를 할당하는 방법

- 바인딩 개체를 생성하고 포드의 바인딩 API에 게시요청
- Pod-bind-definition.yaml

```yaml
apiVersion: v1
kind: Binding
metadata:
	name: nginx
target:
	apiVersion: v1
	kind: Node
	name: node02
```

- curl --header "Content-Type:application/json" --request POST --data

---

# Labels & Selectors

Labels : 각 물품에 부착된 속성

Selectors : 이 항목들을 필터링 하는 것을 도와줌

Select
: kubectl get pods --selector app=App1

: kubectl get pods --selector env=dev

: kubectl get all --selector env=prod

Replicaset

Annotations

---

# Taints And Tolerations(오염과 관용)

**포드와 노드의 관계**
: 어떤 노드에 어떤 포드를 배치할지, 어떻게 제한할 수 있는지

**Taints - Node**

: kubectl taint nodes node-name key=value:taint-effect

: kubectl taint nodes node1 app=blue:NoSchedule

: kubectl taint nodes node01 key=spray:NoSchedule:mortein

**Taints 설정**

: kubectl taint nodes controlplane node-role.kubernetes.io/control-plane:NoSchedule

**Taints 해제**

: kubectl taint nodes controlplane node-role.kubernetes.io/control-plane:NoSchedule-

**Tolerations - PODs**

: kubectl taint nodes node1

```yaml
apiVersion: v1
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

**Taints - NoExecute**

Master Node에 Taint 설정 확인

: kubectl describe node kubemaster | grep Taint

---

# Node Selectors

: 쿠버네티스의 노드 선택기
: 특정노드에서만 동작하도록 Pod의 한계를 설정

1. Node Selectors

pod-definition.yml 파일

```yaml
nodeSelector:
	size: Large
```

1. Node Selector를 사용하여 포드를 만들기 전에 노드에 먼저 라벨을 붙여야함

: kubectl label nodes <node-name> <label-key>=<label-value>
: kubectl label nodes node-1 size=Large

이후 파드 생성
: kubectl create -f pod-definition.yml

Node Selectors는 유용하지만 한계가 있다.
: 더 복잡한 요구사항에는 적용불가 / Pod를 큰 노드나 중간 노드에 놓는다던가, Pod를 작지않은 곳에 놓는다던가.

---

# Node Affinity(노드 선호도)

: 쿠버네티스의 노드 화합 기능
: 특정 노드에 Pod 배치를 제한하는 고급기능제공

pod-definition.yml 파일

```yaml
affinity:
	nodeAffinity:
		requiredDuringSchedulingIgnoredDuringExecution:
			nodeSelectorTerms:
			- matchExpressions:
				- key: size
					operator: In
					Values:
					- Large
					- Medium
```

Node Affinity Types
: Available, Planned

---

# Node Affinity vs Taints and Tolerations(노드선호도 vs 오염과 관용)

- 노드 선호도를 이용해 먼저 각각의 노드에 라벨을 붙힌다.
- Pod의 노드 선택기를 설정해 Pod와 노드 연결

= 다른 Pod가 그 노드에 없다는 보장은 없음.(다른 Pod가 우리노드에 올 가능성이 있음)

= 그러므로 노드 선호도와 오염, 관용 조합을 함께 사용하여 특정 Pod에 노드를 완전히 고정

- 오염과 관용을 사용해 다른 Pod가 우리 노드에 놓이는 것을 막음
- 노드 선호도를 이용해 우리 Pod가 그 사람의 노드에 놓이는 것을 막음

---

# Resource Requirements(리소스 요구사항)

pod-definition.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: simple-webapp-color
	labels:
		name: simple-webapp-color
spec:
	containers:
	- name: simple-webapp-color
		image: simple-webapp-color
		ports:
			- containerPort: 8080
		resources:
			requests:
				memory: "1Gi"
				cpu: 1
			limits:
				memory: "2Gi"
				cpu: 2
				# cpu: 1 = 100m = vCPU: 1
```

limits 보다 많은 메모리를 소모하게 되면 그 무리는 소멸된다.

: OOM(Out Of Memory)

## LimitRange

limit-range-cpu.yaml

```yaml
apiVersion: v1
kind: LimitRange
metadata:
	name: cpu-resource-constraint
spec:
	limits:
	- default:
			cpu: 500m
		defaultRequest:
			cpu: 500m
		max:
			cpu: "1"
		min:
			cpu: 100m
		type: Container
```

limit-range-memory.yaml

```yaml
apiVersion: v1
kind: LimitRange
metadata:
	name: memory-resource-constraint
spec:
	limits:
	- default:
			memory: 1Gi
		defaultRequest:
			memory: 1Gi
		max:
			memory: 1Gi
		min:
			memory: 500Mi
		type: Container
```

## Resource Quotas

: 네임스페이스 레벨에서 쿼터 생성가능

: 네임스페이스 레벨 객체로 요청과 한계에 대한 엄격한 제한 설정을 위해 생성

resource-quota.yaml

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

---

# Daemon Sets

: Pod의 복사본을 클러스터 내 모든 노드에 항상 존재하게 함
: 노드가 삭제되면 Pod도 자동삭제
: 사용사례
- Monitoring Solutuon
- Logs Viewer
- Networking

: 클러스터에 변화가 있을 때 이 노드에서 모니터링 에이전트를 추가하거나, 제거할 필요가 없다.
: 데몬셋이 알아서 해주기 때문

## DaemonSet Definition

daemon-set-definition.yaml (= replicaset-definition.yaml과 유사)

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
	name: monitoring-agent
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

- kubectl create -f daemon-set-definition.yaml
- kubectl get daemonsets
- kubectl describe daemonsets monitoring-daemon

---

# Static PODs

kubelet
: /etc/Kubernetes/manifests 경로 밑에 pod1.yaml, pod2.yaml 파일을 넣어두면 kubelet이 해당 파일을 읽고 호스트에 Pod를 만든다.
: 디렉터리에서 파일을 삭제하면 Pod가 자동으로 삭제된다.
: kubelet이 스스로 만든 이 Pod는 api서버의 간섭이나 나머지 쿠버네티스 클러스터 구성요소의 간섭없이 Static Pod라고 한다.(설정 영향을 받지 않음)
: 복제된 집합이나 배포는 생성할 수 없다.

kubelet.service
: --pod-manifest-path=/etc/Kubernetes/manifests
: --config=kubeconfig.yaml → kubeconfig.yaml(staticPodPath: /etc/Kubernetes/manifests)

유틸리티가 없으므로 docker ps로 확인
: 노드가 클러스터의 일부일 때 = Static 방식으로도 Pod 생성 가능하고, api서버 요청으로도 Pod생성 가능하다.

kubectl get pods -n kube-system

## Static PODs vs DaemonSets

**Static PODs**: Created by the Kubelet / UseCase = Deploy Control Plane components as Static Pods

**Daemon Sets**: Created by Kube-API server(DaemonSet Controller) / UseCase = Deploy Monitoring Agents, Logging Agents on nodes

**Both**: Ignored by the Kube-Scheduler

---

# MULTIPLE SCHEDULERS

my-scheduler-2
my-scheduler-2-config.yaml

```yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: kubeSchedulerConfiguration
profiles:
- schedulerName: my-scheduler-2
```

my-scheduler
my-scheduler-config.yaml

```yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: kubeSchedulerConfiguration
profiles:
- schedulerName: my-scheduler
```

default-scheduler
scheduler-config.yaml

```yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: kubeSchedulerConfiguration
profiles:
- schedulerName: default-scheduler
```

## View Schedulers

: kubectl get pods —namespace=kube-system
pod-definition.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: nginx
spec:
	containers:
	- image: nginx
		name: nginx
		
	schedulerName: my-custom-scheduler
```

kubectl create -f pod-definition.yaml
: 스케줄러가 정상적으로 설정되지 않으면 Pod는 계속 보류상태로 남게됌.

kubectl get pods

관리자가 스케줄러 적용 확인하는 방법
: kubectl get events -o wide
: 네임스페이스에 있는 모든 이벤트를 열거하고 예정된 이벤트를 찾음

View Scheduler Logs
: kubectl logs my-custom-scheduler --name-space=kube-system

---

# SCHEDULER PROFILE

pod-definition.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: simple-webapp-color
spec:
	**priorityClassName: high-priority**
	containers:
	- name: simple-webapp-color
		image: simple-webapp-color
		resources:
			requests:
				memory: "1Gi"
				cpu: 10
```

## Scheduling Plugins

Scheduling Queue

Filtering

Scoring

Binding

my-scheduler-2
my-scheduler-2-config.yaml

```yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: kubeSchedulerConfiguration
profiles:
- schedulerName: my-scheduler-2
	plugins:
		score:
			disabled:
				- name: TaintToleration
				enabled:
				- name: MyCustomPluginA
				- name: MyCustomPluginB

- schedulerName: my-scheduler-3
	plugins:
		prescore:
			disabled:
				- name: '*'
		score:
			disabled:
				- name: '*'

- schedulerName: my-scheduler-4
```