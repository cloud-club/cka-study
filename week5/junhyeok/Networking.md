# CNI

![image.png](attachment:61212d93-122e-44b4-8ed0-92358a35ba89:image.png)

1. 쿠버네티스가 docker 컨테이너를 만들 때 비네트워크에 컨테이너를 만든다.
2. 구성된 CNI 플러그인 호출

## Cluster Networking

- 각 노드는 네트워크에 연결된 인터페이스가 최소한 하나는 있어야 한다.

![image.png](attachment:80df0c6f-c335-4744-9f91-c630da24c248:image.png)

![image.png](attachment:89bd90c3-813f-4128-ac43-6df6699e7abf:image.png)

## Practice Test - Certificates API

```yaml
ssh node01
ip addr
netstat -plnt
```

---

# Pod Networking

- pod는 IP주소를 가지고 있어야 한다.
- 모든 Pod는 같은 노드에 있는 다른 pod들과 소통이 가능해야 한다.
- 모든 Pod는 다른 노드에 있는 Pod들과 NAT없이 소통이 가능해야 한다.

## CNI in kubernetes

View CNI configuration

- ls /opt/cni/bin
- /etc/cni/net.d

## CNI WEAVE

: CNI 플러그인

Deploy Weave

- Service or Pod로 배포
- 

daemonset으로 배포

- kubectl apply -f "https://cloud. weave.works/k8s/net?ks-version=$(kubectl version | base64 | tr -d '\n')”

: 주어진 종류의 Pod하나가 클러스터 내 모든 노드에 배포되는지 확인

각 노드에 배포된 weave peers 확인

- kubectl get pods -n kube-system weave-net-5gcmb

weave-net-~

- kubectl logs weave-net-5gcmb weave -n kube-system

## IP Address Management (IPAM) - Weave

: IP 주소관리와 쿠버네티스는 어떻게 작동할까

노드에서 가상 브리지 네트워크가 어떻게 IP 서브넷을 할당하느냐

pod도 어떻게 IP가 할당되는가 어디에 저장되며 중복IP는 누가 배정하는가

CNI의 임무

- 컨테이너에 IP를 할당

cat /etc/cni/net.d/net-script.conf

---

# Service Networking

service는 클러스터전체 개념. 클러스터의 모든 노드에 걸쳐 존재

- kubectl get pods -o wide
- kubectl get service
- kubectl -L -t nat | grep (serviceName)
- cat /var/log/kube-proxy.log

---

# DNS in kubernetes / CoreDNS in kubernetes

쿠버네티스는 클러스터를 설정할 때 기본 탑재된 DNS서버를 배포함

- 기본 DNS 동작

dns server(10.96.0.10)

```yaml
web 10.244.2.5
test 10.244.1.4
db 10.244.2.15
```

각 서버 resolv.conf 파일에 nameserver 추가

vi /etc/resolv.conf

nameserver 10.96.0.10

DNS 서버는 CoreDNS

- 클러스터에서 CoreDNS 설정

: 쿠버네티스 클러스터의 Kube-system 네임스페이스에 Pod로 배포

- cat /etc/coredns/Corefile
- kubectl get configmap -n kube-system

CoreDNS 솔루션을 배포할 때 서비스도 생성

- 서비스는 자동으로 kubeedns로 명명
- kubectl get service -n kube-system
- pod의 DNS 설정은 쿠버네티스가 자동으로 함

- cat /var/lib/kubelet/config-yaml
: DNS서버와 도메인의 Ip 확인가능

- host web-service
: DNS 확인 명령어

---

# INGRESS

- 사용자가 외부에서 액세스할 수 있는 단일 URL을 사용하여 애플리케이션에 액세스할 수 있도록 지원
- URL 경로에 따라 클러스터 내의 여러 서비스로 트래픽을 라우팅하도록 구성할 수 있다.
- 동시에 SSL 보안 구현도 할 수 있다.
- ingress가 있더라도 클러스터 외부에서 액세스할 수 있도록 노출해야함
: Nodeport, 클라우드 네이티브 LoadBalancer
- ingress controller에서 SsL 및 URL기반 라우팅 구성의 모든 부하분산 작업을 수행

Ingress Controller

: 배포하는 솔루션

Ingress Resource

: 구성하는 규칙 집합

쿠버네티스에는 ingress controller가 없으므로 controller를 배포해야 한다.

**Ingress Controller**

![image.png](attachment:bc4a6475-1551-4133-bc0d-9578a3a1239f:image.png)

**Ingress Resource**

![image.png](attachment:9b7ae24f-9c71-4fe4-9cc1-4ac5f489bafb:image.png)

![image.png](attachment:feaf804e-2dd0-4c8d-b5de-fb3f4e8b78ff:image.png)

- kubectl create -f ingress-wear-watch.yaml
- kubectl describe ingress ingress-wear-watch

![image.png](attachment:bf62025c-726e-4973-8b0b-a4212136a636:image.png)

**rule 하나일때와 둘일 때 차이**

![image.png](attachment:6a98d1f7-3fcd-4301-9f0b-47b82955b96f:image.png)
