# Application Lifecycle 
## Rollout 
![image](https://github.com/user-attachments/assets/748dfa8f-f057-44a6-ae03-2424396b157d)

1. 처음 배포를 생성하면 롤아웃이 일어난다.
2. 새로운 롤아웃은 새로운 배포 revision을 생성한다. revision은 배포에 일어난 변화를 추적할 수 있게 한다. 
3. 앱이 업그레이드되면, 새 롤아웃이 트리거된다. (revision1 -> revision2)


## Deployment Strategy
![image](https://github.com/user-attachments/assets/11c8dffa-6e50-4497-a682-42d89ecd9715)

### 1. Recreate
인스턴스가 5개일때, 새로운 버전을 동시에 띄워 기존 것을 파괴한 후 새 것을 배포한다. 
- 단점: 새 버전이 업데이트되기 까지, 다운된다. (사용자 접근 X)


### 2. Rolling Update(기본)
하나씩 새 버전으로 올린다. 업그레이드하는 동안 다운되지 않는다. 

## ConfigMap
1. ConfigMap 생성
- `kubectl create configmap <name> --form-literal=<key>=<value>`
- config-map.yaml 생성 후 `kubectl create -f config-map.yaml`

2. Inject into Pod
- `envFrom` - `configMapRef` - `name` 으로 configmap 지정


## Configure Secrets in Applications
1. create secret
-  `kubectl create secret generic <name> --from-literal=<key>=<value>`
-  secret-data.yaml 생성 후 `kubectl create -f secret-data.yaml`
2. inject into pod
- `envFrom` - `secretRef` - `name` 으로 secret 지정


## AutoScaling
- 수직 확장: 애플리케이션을 중단하고 더 많은 cpu, memory 리소스 추가 후 전원을 킨다.
- 수평 확장: 서버를 추가해 애플리케이션의 인스턴스를 더 많이 실행한다. 
- 오케스트레이션의 목적: 컨테이너 형태로 애플리케이션을 호스팅하고 확장 및 축소

### 쿠버네티스 스케일링의 유형 
#### 1. 워크로드 확장(Scaling Workloads)
- 수직 확장: 기존 파드에 할당된 리소스 늘리기
  - 수동: `kubectl edit` 으로 파드와 관련된 리소스 요청을 변경
  - 자동: VPA (Vertical Pod AutoScaler)
- 수평 확장: 더 많은 파드 생성
  - 수동: `kubectl scale`로 파드 수 조절
  - 자동: HPA (Horizontal Pod AutoScaler)

#### 2. 클러스터 자체의 확장(Scaling Cluster Infra)
- 수직 확장: 기존 노드에서 리소스 늘리기
  - 수동: 일반적인 접근 방식 X
- 수평 확장: 클러스터에 노드 추가
  - 수동: 새 노드를 수동으로 프로비저닝한 다음, `kubeadm join`으로 클러스터에 노드 추가
  - 자동: Cluster AutoScaler

## HPA (Horizontal Pod AutoScaler)
1. imperative 방식
`kubectl autoscale deployment my-app --cpu-percent=50 -min=1 --max=10`
- 메트릭 서버를 지속적으로 가져와 사용량을 모니터링하고, 사용량에 따라 레플리카셋 수를 조절한다. 
- HPA의 상태는 `kubectl get hpa`로 실행한다.
- HPA가 더이상 필요하지 않다면, `kubectl delete hpa`로 실행한다. 

1. declarative 방식
- 파일 생성: `my-app-hpa.yaml`


---
# Cluster Maintenance
## Software Versions
쿠버네티스 버전 `v1.11.3`의 구조
- 1: major
- 11: minor (Features, Functionalities)
- 3: patch

etcd, coreDNS는 버전이 다를 수 있으나, kube-apiserver, controller-manager, kube-scheduler, kubelet, kube-proxy, kubectl은 버전이 비슷하다. 
- 단, kube-apiserver > controller-manager = kube-scheduler = kubectl > kubelet = kube-proxy
- apiserver가 핵심이기 때문.
