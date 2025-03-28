# 5. Application Lifecycle Management

# Rolling Updates and Rollbacks

Rollout Command

- kubectl rollout status deployment/myapp-deployment
- kubectl rollout history deployment/myapp-deployment

replicasets 확인

- kubectl get replicasets

Rollback

- kubectl rollout undo deployment/myapp-deployment

Summarize Commands

**Create**: kubectl create -f deployment-definition.yml

**Get**: kubectl get deployments

**Update**: kubectl apply -f deployment-definition.yml

               kubectl set image deployment/myapp-deployment nginx-nginx:1.9.1

**Status**: kubectl rollout status deployment/myapp-deployment

             kubectl rollout history deployment/myapp-deployment

**Rollback**: kubectl rollout undo deployment/myapp-deployment

## Practice Test - Rolling Updates and Rollbacks

```yaml
kubectl edit deployment frontend
kubectl describe deploymenr frontend
```

---

# Commands & Arguments

docker run --name ubuntu-sleeper ubuntu-sleeper 10

**pod-definition.yml**

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: ubuntu-sleeper-pod
spec:
	containers:
		- name: ubuntu-sleeper
			image: ubuntu-sleeper
			command: ["sleep2.0"]
			args: ["10"]
```

- kubectl create -f pod-definition.yml

—

FROM Ubuntu

ENTRYPOINT ["sleep"] = command: ["sleep2.0"]

CMD ["5"] - args: ["10"]

## Practice Test - Commands and Arguments

```yaml
apiVersion: v1
kind: Pod 
metadata:
  name: ubuntu-sleeper-3
spec:
  containers:
  - name: ubuntu
    image: ubuntu
    command:
      - "sleep"
      - "1200"
      
k create -f ubuntu-sleeper-3.yaml
k describe pod ubuntu-sleeper-3

apiVersion: v1
kind: Pod 
metadata:
  name: webapp-green
  labels:
      name: webapp-green
spec:
  containers:
  - name: simple-webapp
    image: kodekloud/webapp-color
    command: ["python", "app.py"]
    args: ["--color", "green"]
    
k create -f webapp-color-pod-2.yaml
k describe pod webapp-green
```

---

# Configure environment Variables in Applications

ENV Variables in Kubernetes(ENV Value Types)

**pod-definition.yaml**

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: simple-webapp-color
spec:
	containers:
		- name: simple-webapp-color
			image: simple-webapp-color
			ports:
				- containerPort: 8080

env: #Plain Key Value
	- name: APP COLOR
		value: pink

env: #ConfigMap
	- name: APP_COLOR
		valueFrom:
			configMapKeyRef:

env: #Secrets
	- name: APP_COLOR
		valueFrom:
			secretKeyRef:
```

# Configuring configmaps in Application(쿠버네티스의 구성 데이터 활용법)

Create ConfigMaps

**ConfigMap**

APP_COLOR: blue

APP_MODE: prod

Imperative 방식(파일 정의 X)

- kubectl create configmap
- kubecti create configmap app-config --from-literal-APP_COLOR=blue
- kubecti create configmap app-config --from-literal-APP_COLOR=blue --from-literal-APP_MOD=prod

Declarative 방식(파일 정의)

- kubectl create -f

config-map.yaml

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
	name: app-config
data:
	APP_COLOR: blue
	APP_MODE: prod
```

kubectl create -f config-map.yaml

## View ConfigMaps

- kubectl get configmaps
- kubectl describe configmaps

## ConfigMap in Pods

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
		**envForm:
			- configMapRef:
					name: app-config**
```

config-map.yaml

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
	**name: app-config**
data:
	APP_COLOR: blue
	APP_MODE: prod
```

= 파란 배경의 웹 응용 프로그램을 만든다.

ENV 방식

```yaml
**envForm:
	- configMapRef:
		name: app-config**
```

SINGLE ENV 방식

```yaml
env:
	- name: APP_COLOR
		valueFrom:
			configMapKeyRef:
				name: app-config
				key: APP_COLOR
```

VOLUME 방식

```yaml
volumes:
- name: app-config-volume
	configMap:
		name: app-config
```

## Practice Test - Environment Variables

```yaml
pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: webapp-color
  labels:
    name: webapp-color
spec:
  containers:
  - name: webapp-color
    image: kodekloud/webapp-color
    env:
    - name: APP_COLOR
      value: green

k create -f sample.yaml

vi config-map.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: webapp-config-map
data:
  APP_COLOR: darkblue
  APP_OTHER: disregard
  
  k create -f config-map.yaml
  
  
pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: webapp-color
  labels:
    name: webapp-color
spec:
  containers:
  - name: webapp-color
    image: kodekloud/webapp-color
    **env:
      - name: APP_COLOR
        valueFrom:
          configMapKeyRef:
            name: webapp-config-map
            key: APP_COLOR**
```

---

# Configuring ConfigMaps in Applications

쿠버네티스의 구성 데이터 활용법

Create ConfigMaps

ConfigMap

```yaml
APP_COLOR: blue
APP_MODE: prod
```

- kubectl create configmap
- kubecti create configmap app-config --from-literal-APP_COLOR=blue

---

# Configure Secrets in Application

Secret

1. Create Secret - 시크릿 생성
2. Inject into Pod- Pod에 주입

```yaml
DB_HOST: mysql
DB_User: root
DB_Password: paswrd
```

## Create Secret

Imperative(명령적)

: kubectl create secret generic app-secret --from-literal-DB_Host=mysql

Declarative(선언적)

secret-data.yaml

```yaml
apiVersion: V1
kind: Secret
metadata:
	name: app-secret
data:
	DB_Host: mysql
	DB_User: root
	DB_Password: paswrd
```

kubectl create -f secret-data.yaml

## Encode Secrets

secret data는 인코딩된 형심으로 데이터를 지정해야한다.

echo -n 'mysql' | base64

echo -n 'root' | base64

echo -n 'paswrd' | base64

```yaml
data:
	DB_Host: bXIzcWw=
	DB_User: cm9vdA==
	DB_Password: cGFzd3Jk
```

## View Secrets

kubectl get secrets
kubectl describe secrets
kubectl get secret app-secret -o yaml

## Decode Secrets

echo -n 'bXIzcWw=' | base64 —d
echo -n 'cm9vdA==' | base64 —d
echo -n 'cGFzd3Jk' | base64 —d

## Secrets in Pods

secret-data.yaml

```yaml
apiversion: v1
kind: Secret
metadata:
	name: app-secret
data:
DB_Host: bXIzcWw=
DB User: cm9vdA-=
DB_Password: cGFzd3Jk
```

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
			envFrom:
				- secretRef:
					name: app-secret
```

kubectl create -f pod-definition.yaml

```yaml
ENV
envFrom:
	- secretRef:
		name: app-secret

SINGLE ENV
env:
	- name: DB_Password
	valueFrom:
	secretKeyRef:
		name: app-secret
		key: Db_Password
		
VOLUME
volumes:
- name: app-secret-volume
	secret:
		secretName: app-secret
```

## Secrets in Pods as Volumes

```yaml
volumes:
- name: app-secret-volume
	secret:
		secretName: app-secret
```

ls /opt/app-secret-volumes
: DB_Host DB_Password DB_User

cat /opt/app-secret-volumes/DB_Password
: paswrd

## Note on Secrets

- Secrets are not Encrypted. Only encoded
: Do not Check-in Secret objects to SCM along with code.
- Secrets are not encrypted in ETCD
: Enable encryption at rest
- Anyone able to create pods/deployments in the same namespace can access the secrets
: Configure least-privilege access to Secrets - RBAC
- Consider third-party secrets store providers
: AWS Provider, Azure Provider, GCP Provider, Vault Provider

## Practice Test - Secrets

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: webapp-pod
spec:
  containers:
  - name: webapp-pod
    image: kodekloud/simple-webapp-mysql
    envFrom:
    - secretRef:
        name: db-secret
```

---

# Multi Container Pods

pod-definition.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: simple-webapp
	labels:
		name: simple-webapp
spec:
	containers:
		- name: simple-webapp
			image: simple-webapp
			ports:
				- containerPort: 8080
				
		- name: log-agent
			image: log-agent
```

## Practice Test - **Multi Container PODs**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: yellow
spec:
  containers:
  - name: lemon
    image: busybox
    command:
      - sleep
      - "1000"

  - name: gold
    image: redis
   
---
    
apiVersion: v1
kind: Pod
metadata:
  name: app
  namespace: elastic-stack
  labels:
    name: app
spec:
  containers:
  - name: app
    image: kodekloud/event-simulator
    volumeMounts:
    - mountPath: /var/log/event-simulator/
      name: log-volume

  - name: sidecar
    image: kodekloud/filebeat-configured
    volumeMounts:
    - mountPath: /var/log/event-simulator/
      name: log-volume

  volumes:
  - name: log-volume
    hostPath:
      path: /var/log/webapp
      type: DirectoryOrCreate
```

---

# InitContainers

```yaml
apiversion: v1
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
		initcontainers:
		- name: init-myservice
			image: busybox
			command: ['sh', '-c', 'git clone «some-repository-that-will-be-used-by-application>; done;']
```

: Pod가 처음 생성되면 initcontainer가 실행되고 initContainer의 프로세스는 애플리케이션을 호스팅하는 실제 컨테이너가 시작되기 전에 완료될 때까지 실행되어야 한다.

: initcontainers를 여러 개 구성할 수 있다.이 경우 각 initcontainer는 순차적으로 한번에 하나씩 실행된다.

: initContainers가 완료되지 않으면 Kubernetes는 init container가 성공할 때까지 Pod를 반복적으로 다시 시작한다.

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
		image: busybox:1.28
		command: ['sh', '-c', 'until nslookup myservice; do echo waiting for myservice; sleep2; done;']
	- name: init-mydb
		image: busybox:1.28
		command: ['sh', '-C', 'until nslookup mydb; do echo waiting for mydb; sleep2; done;']
```

## Practice Test - **Init Containers**

```yaml
k edit pod red
k edit pod orange
k replace --force -f /tmp/kubectl-edit-4183639490.yaml
```

---

# Introduction to Autoscaling

Horizontal Pod Autoscaling(HPA) - 수평확장
Vertical pod Autoscaling(VPA) - 수직확장

## Practice Test - **Manual Scaling**

```yaml
k get deployment

apiVersion: apps/v1
kind: Deployment
metadata:
  name: flask-web-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: flask-app
  template:
    metadata:
      labels:
        app: flask-app
    spec:
      containers:
      - name: flask
        image: rakshithraka/flask-web-app
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: flask-web-app-service
spec:
  type: ClusterIP
  selector:
    app: flask-app
  ports:
   - port: 80
     targetPort: 80
```

---

# Horizontal Pod Autoscaling(HPA)

## Scaling a workload the manual way(워크로드 수동수평확장)

nginx.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
	name: myapp
spec:
	replicas: 1
	selector:
		matchLabels:
			app: my-app
	template:
		metadata:
			labels:
				app: my-app
		spec:
			containers:
			- name: my-app
				image: nginx
				resource:
					requests:
						cpu: "250m"
					limits:
						cpu: "500m"
```

: kubecti top pod my-app-pod - CPU, MEMORY 301
: kubectl scale deployment my-app --replicas=3 - CPU, MEMORY 임계값 도달 시 수동 확장

문제점 : 컴퓨터 앞에 앉아 리소스 사용량을 지속적으로 모니터링해야함.
             수동으로 확장해야함, 트래픽 급증에 대해 빠르게 대응을 못함

## 따라서 Horizontal pod Autoscaling(HPA) 사용

kubectl autoscale deployment my-app --cpu-percent=50 --min=1 --max=10

생성된 HPA 상태 확인
kubectl get hpa

HPA 삭제
kubectl delete hpa my-app

## 선언적 방식

my-app-hpa.yaml

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

## Practice Test - **HPA**

```yaml
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  creationTimestamp: null
  name: nginx-deployment
spec:
  maxReplicas: 3
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: nginx-deployment
  targetCPUUtilizationPercentage: 80
status:
  currentReplicas: 0
  desiredReplicas: 0
 
k create -f autoscale.yml
k get hpa
```

---

# In-place Resize of Pods(수동 수직 확장)

nginx.yaml

```yaml
apiversion: apps/v1
kind: Deployment
metadata:
	name: myapp
spec:
	replicas: 1
	selector:
		matchLabels:
			app: my-app
	template:
		metadata:
			labels:
				app: my-app
		spec:
			containers:
			- name: my-app
				image: nginx
				resizePolicy:
					- resourceName: cpu
						restartPolicy: NotRequired
					- resourceName: memory
						restartPolicy: RestartContainer
				resource:
					requests:
						cpu: "250m" →> "1"
						memory: "256Mi"
					limits:
						cpu: "500m"
						memory: "512Mi"
```

FEATURE_GATES=InplacepodVerticalScaling=true로 설정
컨테이너를 다시 시작하지 않고도 컨테이너에 할당된 CPU 또는 메모리 리소스의 크기를 조정하는 방법.

---

# Vertical Pod AutoScaler(VPA)

## Scaling resources for a workload the manual way(워크로드 수동수직확장)

nginx.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
	name: myapp
spec:
	replicas: 1
	selector:
		matchLabels:
			app: my-app
	template:
		metadata:
			labels:
				app: my-app
		spec:
			containers:
			- name: my-app
				image: nginx
				resource:
					requests:
						cpu: "250m"
					limits:
						cpu: "500m"
	
```

: kubectl top pod my-app-pod - CPU, MEMORY 확인
 : kubectl edit deployment my-app - CPU, MEMORY 임계값 도달 시 수동 확징

## Vertical Pod Autoscaler(HPA) 사용

Horizontal과 다르게 기본제공되지 않음 / 배포해야한다.
GitHu레포지토리에 있는 Vertical Pod Autoscaler 정의 파일 적용
: kubectl apply -f https://github.com/kubernetes/autoscaler/releases/latest/download/vertical-pod-autoscalen.yam
: kubecti get pods -n kube-system | grep vpa

```yaml
VPA Recommender : 정보를 수집
VPA Updater: Recomender로부터 정보를 모니터링하거나 가져와서 실제 Pod와 비교한 다음 Pod가 임계값을 초과하면 Pod를 죽인다.
						 파드가 죽으면 자동으로 배포가 파드를 다시 생성함
VPA Admission Controller: 이후 개입하여 리소스를 업데이트하여 파드가 새로운 크기를 갖출 수 있도록 한다.
```

## 선언적 방식

my-app-vpa.yaml

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
	name: my-app-hpa
spec:
	targetRef:
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

VPA는 4가지 모드로 수행

```yaml
Off: Only recommends. Does not change anything
Initial: Only changes on Pod creation. Not later.
Recreate: Evicts pods if usage goes beyound range
Auto: Updates existing pods to recommended numbers. For now this behaves similar to Recreate.
But when support for "In-place Update of Pod Resources" is available that made will be preferred.
```

권장사항 확인 명령어
: kubectl describe vpa my-app-vpa

## Key Differences

- 사진 추가 필요

## Practice Test - Install VPA

```yaml
kubectl get crds | grep verticalpodautoscaler
k get deployment -n kube-system

```

## Practice Test - **Modifying CPU resources in VPA**
```yaml
kubectl top pod
k get vpa flask-app -o yaml
```