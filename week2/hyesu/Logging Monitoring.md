# Monitor Cluster Components

오픈소스 솔루션 많음

### Heapster

- 초기
- 사라지고, Metric Server라는 간소화된 버전 만들어짐

### Metric Server

- Cluster 당 메트릭 서버 1개
- Metric Server는 각 Node와 Pod에서 Metric 회수
- In-Memory 모니터링 솔루션

⇒ 이후, 고급 모니터링 솔루션에 의존

Node와 Pod에 대한 지표 어떻게 생성??

⇒ kubelet

- cAdvisor라는 요소 포함
    - Pod에서 성능 메트릭 회수하고 kubelet API를 통해 메트릭을 공개해 Metric Server에서 메트릭 사용 가능하도록 함

```bash
# minikube
minikube addons enable metrics-server

# others
git clone ...
kubectl create -f deploy/1.8+/
```

```bash
kubectl top node # 노드 성능 지표
Kubectl top pod # 파드 성능 지표
```

# Application Logs

### Docker

```bash
# 난수 이벤트 발생, 웹 서버 시뮬레이션
docker run kodekloud/event-simulator

docker run -d kodekloud/event-simulator

docker logs -f ecf
# -f : 실시간 옵션
```

### Kubernetes

```bash
k logs -f event-simulator-pod
```

Pod에 컨테이너 2개인 상태에서 

Pod이름으로 logs 하는 경우 → 실패. 컨테이너 명시해야함.
