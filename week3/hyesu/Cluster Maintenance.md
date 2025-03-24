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
kubeadm upgrade plan
```

클러스터 업그레이드 전에 kubeadm 업그레이드 필요!

- 업그레이드 예제

현재 1.11 → 1.13

```bash
apt-get upgrade -y kubeadm=1.12.0-00 # kubeadm 1.12로 업그레이드

```
