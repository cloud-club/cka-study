# Kubernetes Security Primitives - 쿠버네티스 보안 개요

클러스터 보안을 위해 어떤 위험과 조치를 취해야 하나?

### API 서버 자체에 대한 액세스 제어 : 1차 방어선

두가지 유형의 결정

1. 누가 클러스터에 접근할 수 있는가
2. 접근해서 뭘 할 수 있는가

Authentication

1. Who can access?

: 인증 메커니즘에 의해 정의 됌.

- Files - Username and Passwords
- Files - Username and Tokens
- Certificates
- External Authentication providers - LDAP
- Service Accounts

Authorization

2. What can they do?

- RBAC Authorization
- ABAC Authorization
- Node Authorization
- Webhook Mode

### TLS Certificates

: 다양한 구성요소 사이에 인증서를 설정하는 방법

### Network Policies

: 클러스터 내 응용프로그램간의 통신
: 기본 설정 상 모든 Pod는 무리 내의 다른 Pod에 접근 가능
: 네트워크 정책을 이용해 그들 사이의 액세스 권한 제한

---

# Authentication

: 인증 메커니즘을 통해 쿠버네티스 클러스터에 액세스 권한을 확보

### Accounts

- User (Admins, Developers)
- Service Accounts

쿠버네티스는 User 계정을 직접 관리하지 않음.

: 외부 소스에 의존한다

- 사용자 세부정보나 인증서가 있는 파일
- LDAP같은 타사 ID서비스

서비스 계정의 경우 쿠버네티스가 관리 가능
: 쿠버네티스 API를 이용해 서비스 계정을 생성하고 관리
: kubectl create serviceaccount sa1
: kubectl get serviceaccount

### Auth Mechanisms

kube-apiserver는 요청을 처리하기 전에 인증을 진행

: Static Password File

: Static Token File

: Certificates

: Identity Services

### Auth Mechanisms - Basic

: csv파일에 사용자 목록과 암호를 만들 수 있고 그걸 사용자 정보의 소스로 사용

- user-details.csv 파일

password123,user1,u0001
password123,user2,u0002
password123,user3,u0003
password123,user4,u0004
password123,user5,u0005

- kube-apiserver.service 파일 내용추가
    
    --basic-auth-file=user-details.csv 추가
    

### Kube-api Server Configuration

: Kube-api Server 수정

: 파일 업데이트 하면 자동으로 Kube API서버를 재시작

/etc/kubernetes/manifests/kube-apiserver.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
	creationTimestamp: null
	name: kube-apiserver
	namespace: kube-system
spec:
	containsers:
	- command:
		- kube-apiserver
		- --authorization-mode-Node, RBAC
		- --advertise-address=172.17.0.107
		- --allow-privileged-true
		- --enable-admission-plugins=NodeRestriction
		- --basic-auth-file=user-details.csv
		image: k8s-gcr.io/kube-apiserver-amd64:v1.11.3
		name: kube-apiserver
```

### Authenticate User

: 기본 자격 증명을 이용해 API 서버에 액세스

- 사용자 및 암호 지정

: curl -v -k https://master-node-ip:6443/api/v1/pods -u “user1:password123”

### Auth Mechanisms - Basic2

- user-details.csv 파일 - Static Password File
: 선택적으로 네번째 열에 그룹정보 추가 가능

password123,user1,u0001,group1
password123,user2,u0002,group1
password123,user3,u0003,group2
password123,user4,u0004,group2
password123,user5,u0005,group2

- user-token-details.csv - Static Token File
    
    : 암호대신 토큰지정
    

qwudhqwiudhgowudhqoiwduhgoiwudh1uhi1,user10,u0010,group1

aiwudh12iuhdwaiudhaiwdhudh2uhiauhduu,user11,u0011,group1

wa121asdjivgioubasiud2uhdaiuhsd2123u,user12,u0012,group2

qwdiuv980weubqiwjdbqwiudbqiwu2aysgdi,user13,u0013,group2

-- token-auth-file=user-token-details.csv

- 인증을 할 때 그 토큰을 승인차단 토큰으로 지정

: curl -v -k https://master-node-ip:6443/api/v1/pods --header "Authorization: Bearer qwudhqwiudhqowudhqoiwduhqoiwudh1uhi1"

### NOTE

: This is not a recommended authentication mechanism

: 가장 쉬운방법이어서 소개

---

# TLS CERTIFICATES - 일반 인증서에 대한 소개

### TLS - Basics

인증서 : 거래도중 상호 신뢰를 보장하기 위해 사용 됌

- 사용자가 웹 서버에 액세스하려 할 때 TLS 인증서는 사용자와 서버 사이의 통신이 암호화 되도록 한다.
- 안전한 연결성이 없으면 사용자가 온라인 뱅킹 앱에 접속할 경우 입력한 자격증명이 일반텍스트 형식으로 전송
- 해커가 네트워크 트래픽을 탐지하면 자격증명을 쉽게 추출해 사용자의 은행계좌를 해킹가능.
- 안전하지 않으니 전송되는 데이터를 암호화해야한다.

- 키로 데이터를 암호화 : 키는 무작위 숫자와 알파벳의 집합

: 데이터에 임의의 숫자를 추가하고 인식할 수 없는 포맷으로 암호화

: 그 데이터는 서버로 보내짐

: 네트워크 상에서 해커가 데이터를 얻지만 아무것도 할 수 없다.

: 하지만 데이터를 받는 서버도 마찬가지

: 키 없이는 데이터를 해독할수 없으므로 키의 복사본도 서버로 보내서 서버가 메세지를 해독하고 읽을 수 있게 함 : 키도 같은 네트워크로 전송되므로 해커들도 데이터를 해독할 여지가 있음.

= 대칭 암호화: 데이터를 암호화 하고 해독하는데 같은 키를 사용하고

그 키는 수신기와 교환해야하기 때문에 해커가 그키에 접근에 데이터를 해독할 위험이 있음

그래서 비대칭 암호화가 등장

: 데이터를 암호화하고 해독하는데 단일 키를 쓰는 대신 비대칭암호화는 한쌍의 키와 개인키, 공용 키를 사용

: 자물쇠와 열쇠쌍

: 자물쇠로 데이터를 암호화 시 관련 키로만 열수 있음

SSH 액세스를 보안하는 더 간단한 사용사례 (비대칭키 활용)

- 개인키와 공용키 생성
- 공용키로 서버 잠금
- ssh 접근 시 개인키 사용으로 서버 공용키 잠금해제

- 다른 서버 또한 공용 잠금장치 복사본을 만들어 잠금
- 가지고 있는 개인키로 ssh접근

- 다른 사용자가 서버에 접속해야 한다면
- 똑같이 개인키와 공용키 쌍을 따로 생성
- 서버에 추가 문을 만들어 공용장금 장치를 사용.

# TLS in Kubernetes

: TLS인증서를 이용한 쿠버네티스 클러스터 보안

두가지 주요 요구사항

- Server Certificates for Servers
    - 클러스터 내의 다양한 서비스를 이용해 서버 인증서를 사용
- Client Certificates for Clients
    - 클라이언트 인증서를 이용해 정체를 확인

![image.png](attachment:3318224b-f697-4d5c-83dd-80a8babd4c9b:image.png)

### Server Certificates for Servers

KUBE-API server
: HTTPS 서비스지원

- 인증서와 키페어 생성
    - apiserver.crt / apiserver.key

ETCD server
: 클러스터에 관한 모든 정보를 저장

- 인증서와 키페어 생성
    - etcdserver.crt / etcdserver.key

KUBELET server(worker node)
: HTTPS API EndPoint
: kube-apiserver가 worker node와 상호작용하기 위해 통신

- 인증서와 키페어 생성
    - kubelet.crt / kubelet.key

### Client Certificates for Clients

admin(관리자 사용자)
: 서버에 인증하려면 인증서와 키페어 필요

- 인증서와 키페어 생성
    - admin.crt / admin.key

KUBE-SCHEDULER
: kube-apiserver에 보고해서 스케줄이 필요한 Pod를 찾음
: kube-apiserver가 올바른 worker node의 Pod를 찾도록 함
: Scheduler는 admin과 같은 클라이언트일 뿐

- 인증서와 키페어 생성
    - scheduler.crt / scheduler.key

KUBE-CONTROLLER-MANAGER
: kube-apiserver 인증을 위한 인증서 필요

- 인증서와 키페어 생성
    - controller-manager.crt / controller-manager.key

KUBE-PROXY
: kube-apiserver 인증을 위한 인증서 필요

- 인증서와 키페어 생성
    - kube-proxy.crt / kube-proxy.key

KUBE-API server
: 기타 서버들과 통신하는 유일한 서버
: 해당 서버들에게 kube-apiserver는 클라이언트 이기 때문에 인증필요

- 인증서와 키페어 생성
    - apiserver-kubelet-client.crt / apiserver-kubelet-client.key
    - apiserver-etcd-client.crt / apiserver-etcd-client.key
    - kubelet-client.crt / kubelet-client.key

![image.png](attachment:9fcab221-7b13-4c59-a66e-b743da901e0f:image.png)

# TLS in Kubernetes - Certificate Creation

인증서 생성 도구
- EASYRSA
**- OPENSSL**
- CFSSL

### OPENSSL 도구를 사용해 인증서 생성

**CA 인증서**

1. OpenSSL명령으로 Private Key생성
- openssl genrsa -out ca.key 2048
- **ca.key**
2. OpenSSL Request(인증서 서명요청)
- openssl req -new -key ca.key -subj "/CN=KUBERNETES-CA" -out ca.csr
- **ca.csr**
3. 인증서 서명
- openssl x509 -req -in ca.csr -signkey ca.key -out ca.crt
- **ca.crt**

**Client 인증서(ex. ADMIN USER)**

1. OpenSSL명령으로 Private Key생성
- openssl genrsa -out admin.key 2048
- **admin.key**
2. OpenSSL Request(인증서 서명요청)
- openssl req -new -key admin.key -subj \
"/CN=kube-admin**/O=system:masters**" -out admin.csr
- **admin.csr**
- 인증서 서명 요청에서 반드시 언급해야한다.
: **/O=system:masters**
3. 인증서 서명
- CA 인증서와 CA 키 지정
- openssl x509 -req -in admin.csr –CA ca.crt -CAkey ca.key -out admin.crt
- **admin.crt**

**KUBE SCHEDULER**
인증서 이름 = SYSTEM:KUBE-SCHEDULER

**KUBE CONTROLLER MANAGER**
인증서 이름 = SYSTEM:KUBE-CONTROLLER-MANAGER

KUBE PROXY
인증서 이름 = KUBE-PROXY

= ControlPlane의 구성요소 시스템은 이름 앞에 시스템을 붙여야 한다.

### 증명서 활용방안

 - REST API 호출에서 사용자 이름과 암호대신 인증서 사용가능
: curl https://kube-apiserver:6443/api/v1/pods \
--key admin.key --cert admin.crt
--cacert ca.crt

 - 모든 매개 변수들을 kubeconfig 구성파일로 옮김

```yaml
apiVersion: v1
clusters:
- cluster:
certificate-authority: ca.crt
server: https://kube-apiserver:6443
name: kubernetes
kind: Config
users:
- name: kubernetes-admin
user:
client-certificate: admin.crt
client-key: admin.key
```

### 쿠버네티스 구성요소끼리 서로 확인을 위해 CA의 루트 인증서 복사본이 필요

: 서버나 클라이언트가 인증서를 갖고 구성할 때마다 CA 루트 인증서도 지정해야 함

**server 인증서 - ETCD SERVERS**

- 생성절차는 Client 인증서 생성절차와 동일
- ETCD 서버는 고가용성 환경에서 다중 서버에 걸쳐 클러스터로 배포될 수 있다.
    - 클러스터 내 다른 멤버간의 통신을 안전하게 하기 위해서는 추가 PEER 인증서를 생성해야 함
    - 인증서가 생성되면 ETCD서버를 시작하는 동안 지정한다.
- cat etcd.yaml
    - ETCD server 키 명시
    - CA 루트 인증서  필요
        - ETCD server에 연결된 클라이언트가 유효한지 확인

![image.png](attachment:a51cac38-57ef-4b80-890d-817985885fc4:image.png)

**server 인증서 - KUBE-API SERVER**
: 클러스터 내에서 움직이는건 뭐든 API 서버가 안다.
: 그러므로 클러스터 안에 많은 이름과 별칭이 있다.

- KUBE-API SERVER = kubernetes
- kubernetes.default
- kubernetes.default.svc
- kubernetes.default.svc.cluster.local
- IP = Kube API 서버를 실행하는 호스트의 IP주소나 그걸 실행하는 Pod

: 이 모든 이름은 반드시 KUBE API서버를 위해 생성된 인증서에 존재해야한다.

- 그래야만 KUBE-API Server를 참조하는 이들이 유효한 연결을 설정가능

1. OpenSSL명령으로 Private Key생성
- openssl genrsa -out apiserver.key 2048
- **apiserver.key**
2. OpenSSL Request(인증서 서명요청)
- openssl req -new -key apiserrver.key -subj \
"/CN=kube-apiserver" -out apiserver.csr -config openssl.cnf
****- openssl.cnf로 다른 이름 명시
- **apiserver.csr**

```bash
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation,
subjectAltName = @alt_names
[alt_names]
DNS.1 = kubernetes
DNS.2 = kubernetes.default
DNS.3 = kubernetes.default.svc
DNS.4 = kubernetes.default.svc.cluster.local
IP.1 = 10.96.0.1
IP.2 = 172.17.0.87
```

1. 인증서 서명
- openssl x509 -req -in apiserver.csr \
-CA ca.crt -CAkey ca.key -out apiserver.crt
- **apiserver.crt**

![image.png](attachment:48756555-ce59-4630-9256-7a6f8757cee7:image.png)

**server 인증서 - KUBECTL NODES**

: 클러스터의 각 노드에 대해 키 인증서 쌍이 필요
: 인증서의 이름은 노드의 이름을 따서 지음

- CA인증서 및 kubelet node인증서 지정

```yaml
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
	x509:
		clientCAFile: **"/var/lib/kubernetes/ca.pem"**
authorization:
	mode: Webhook
clusterDomain: "cluster.local"
clusterDNS:
	- "10.32.0.10"
podCIDR: "${POD_CIDR}"
resolvConf: "/run/systemd/resolve/resolv.conf"
runtimeRequestTimeout: "15m"
tlsCertFile: **"/var/lib/kubelet/kubelet-node01.crt"**
tlsPrivateKeyFile: **"/var/lib/kubelet/kubelet-node01.key"**
```

![image.png](attachment:6bcf2763-8503-4ff2-b252-66afbe4cd225:image.png)

**client 인증서 - KUBECTL NODES**
- kubectl-client.crt / kubectl-client.key
- 노드이름

- system:node:node01
- system:node:node02
- system:node:node03

- 시스템노드라는 그룹에 노드 추가

Group:SYSTEM:NODES

# View Certificate Details

### kubeadm에 의한 클러스터 프로비저닝

1. API서버 정의 파일 확인
- cat /etc/kubernetes/manifests/kube-apiserver.yaml
2. 각각의 인증서 내부 확인
- --tls-cert-file=**/etc/kubernetes/pki/apiserver.crt**

![image.png](attachment:5b42912c-4443-4fa9-999b-566fea8dd556:image.png)

**/etc/kubernetes/pki/apiserver.crt**

- openssl x509 -in /etc/kubernetes/pki/apiserver.crt -text -noout

![image.png](attachment:5503a2be-9b1f-4ff1-b9c5-0aebcac60d31:image.png)

- kube-API서버의 별칭 모두 확인
    - Subject
- 유효기간 만료일 확인
    - Not After
- 발급자 확인
    - Issuer
    - 증명서를 발행한 CA = kubeadm = kubernetes

### Inspect Service Logs

kubeadm 사용 없이 혼자 클러스터 설정 시 서비스 로그 확인

- journalctl -u etcd.service -l

![image.png](attachment:653e2fe8-64b6-4ddb-b86b-b2d98c58dca3:image.png)

### View Logs

kubeadm으로 클러스터 설정

- kubectl logs etcd-master

![image.png](attachment:027e885a-34ae-4250-ab40-c198e6d4a7da:image.png)

API server, NCD서버 다운으로 kubectl 명령어 작동 안할시 로그 확인 방법

- docker ps -a

![image.png](attachment:fd846387-a088-439a-b0e8-1aa3d01d9f75:image.png)

- docker logs 87fc

![image.png](attachment:c5265de1-4882-4179-80ed-5801db659d3b:image.png)

## Practice Test - View Certificates

```yaml
kube-api server 컨테이너 확인
: crictl ps -a

kube-api server 컨테이너 로그확인
: crictl logs container-id

root@controlplane:~# crictl logs --tail=2 1fb242055cff8  
W0916 14:19:44.771920       1 clientconn.go:1331] [core] grpc: addrConn.createTransport failed to connect to {127.0.0.1:2379 127.0.0.1 <nil> 0 <nil>}. Err: connection error: desc = "transport: authentication handshake failed: x509: certificate signed by unknown authority". Reconnecting...
E0916 14:19:48.689303       1 run.go:74] "command failed" err="context deadline exceeded"

vi kube-apiserver.yaml
- --etcd-cafile=**/etc/kubernetes/pki/etcd/ca.crt**
```

---

# Certificates API

: 인증서 관리법 / API인증서

CA Server
: 생성한 인증서, 키 파일을 저장
: 쿠버네티스 마스터 노드 = CA Server

### 인증서 서명요청 자동화

: 쿠버네티스 내장 인증서 API 사용

1. Create CertificateSigningRequest Object
사용자는 먼저 키를 만들고 인증서 서명요청 생성으로 관리자에게 요청 전송

```yaml
- openssl genrsa -out jane.key 2048
- openssl req -new -key jane.key -subj "/CN=jane" -out jane.csr
```

![image.png](attachment:15fb84b3-b6d2-49e9-b9e4-b07bd64e086f:image.png)

1. Review Requests
csr Base64 암호화 및 manifest파일 생성

```yaml
- csr Base64 암호화
: cat jane.csr | base64

- jane-csr.yaml manifest파일 생성
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
	name: jane
spec:
	groups:
	- system:authenticated
	usages:
	- digital signature
	- key encipherment
	- server auth
	request: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVLS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVFVRVNULS0FVRVNULS0tLS0KTUlJQ1dEQ0NBVUFDQVFBtLS0KTUlJQ1dEQ0NBVUFDQVFBd0V6RVJNQThHQTFVRUd0V6RVJNQThHQTFVRUF3d0libVYzTFhWelF3d0libVYzTFhWelpYSXdnZ0VpTUEwR0NTcUdTSWIzRpYSXdnZ0VpTUEwR0NTcUdTSWIzRFFFQgpBFFFQgpBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRRE8wV0pXUVVBQTRJQkR3QXdnZ0VLQW9JQkFRRE8wV0K0RYc0FKU0lyanBObzV2UklCcGxuemcrNnhjOStVVndpXK0RYc0FKU0lyanBObzV2UklCcGxuemcrrS2kwCkxmQzI3dCsxZUVuT041TXVxOTlOZXZtTUVPbnNnhjOStVVndrS2kwCkxmQzI3dCsxZUVuT0J41TXVxOTlOZXZtTUVPbnJ

- 적용
: kubectl apply -f akshay-csr.yaml
```

![image.png](attachment:a725a483-1fa3-45f8-8357-31d14141305b:image.png)

1. Approve Requests

```yaml
- 새 요청 확인
: kubectl get csr

- 요청 승인
: kubectl certificate approve jane
```

![image.png](attachment:afa0b976-feea-4f70-87df-b454049db194:image.png)

1. Share Certs to Users

```yaml
- 생성된 인증서 확인 가능
: kubectl get csr jane -o yaml

- 암호화된 인증서 decode
: Íecho “LS0…Qo=” | base64 --decode
```

![image.png](attachment:eb0864da-24c6-4e0b-8f25-a2087557c15e:image.png)

### 인증서 작업주체

모든 인증서 관련 작업은 Kube-Api Server 내 Controller Manager에 의해 실행

![image.png](attachment:fda1d880-8031-4faa-bd28-3070f44417e6:image.png)

![image.png](attachment:9ad1bf80-25ab-4336-a0e0-4296cbc0e38b:image.png)

인증서에 서명해야 한다면 CA서버의 인증서와 개인키가 필요

![image.png](attachment:938cce4f-fe15-4709-bc18-eeeb8fa4e6ce:image.png)

## Practice Test - Certificates API

```yaml
cat akshay.csr | base64 -w 0

apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: akshay
spec:
  groups:
  - system:authenticated
  request: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVFVRVNULS0tLS0KTUlJQ1ZqQ0NBVDRDQVFBd0VURVBNQTBHQTFVRUF3d0dZV3R6YUdGNU1JSUJJakFOQmdrcWhraUc5dzBCQVFFRgpBQU9DQVE4QU1JSUJDZ0tDQVFFQXBPcTJGMjFFYzVBOTlGYnltRGZ5QWtBOTlVNjVranlFNmJCK3JCeU5VOEI2CkVLQU5SMjJBS0V6NDRTbmg0SXk0NElTNEdDK01xSVR6OUk4Ylc5QlJ5L0o0K1lEdk1hVzA0VVJ5Z0dQbHdudmwKL3FOUFl2ZGlNUU5ZSUxPdzRpKzZBRklPZ3haeTBwc3dSbGdLSnJ1QjdudlMveExtbDM2ODdqd2NsYlRsWHhOQwpsanZPTHFEc05KRDVPbnQvTU42YkhONnFndEprTFY4a2I1K0FvMjNyd1U2RENRb21uMmYxSHQrdWFvVVJFRXB3ClJmYmhySFpWVngrUFB0NEI4Y1pzdnJ2MnBmYU5MMUE1UzFLNHkwR0hkeWNWYi8zUVpuOHIzcUM5Q3Z5c2tMWmIKcURBRlRKMGpiWGFuTHV2ZGVCdDRJenVIK21hN25sL1kycVh4M0F3ZU93SURBUUFCb0FBd0RRWUpLb1pJaHZjTgpBUUVMQlFBRGdnRUJBSk9mYkFVQjZJTXJtVTdxK3NiQTB1QWRXSFlGZ1ZNalhHWUE1QjBJRjF2VFZ6U0s4MTVvCnRDTjFrbWU5SDgzSG52akFWcGgrczZYRitrNDVtdmNhM0tGTzR5bFBrU0VBQitqZTNKTE1NS2JzTytlZ2tKdGcKU0FwVG1nVVdyNDRwNk1ST1JGTG1FRWFkaSt6SGxLTi9WRTBDeG8wb094NzlOTlVKQU9LYTRyQ0VoVjBITDVlWApiaGNRRnliOEN0NTQwZzY2dFlGVk1lUENVMnowMXdNeG10enE1SGFKL3F1ZXBVWmRrc1FKd0ZZRlpiV3dWOEgwCkRyVi9SdEpCUE95ancwczR6RVdkSVc3eS9BYW0rMkhFUUlBU2M1M1ZjWkNsbUltM2VibnZhYkEyS3FOT0c3RE8KSXltblhkOVp0TUZMSEgySmQycTlRVmtleUUwMGJDVDNvbkU9Ci0tLS0tRU5EIENFUlRJRklDQVRFIFJFUVVFU1QtLS0tLQo=
  signerName: kubernetes.io/kube-apiserver-client
  usages:
  - client auth
  
k get csr
kubectl certificate approve akshay
k get csr akshay -o yaml
k get csr agent-smith -o yaml
k certificate deny agent-smith
k delete csr agent-smith
```

---

# KubeConfig

옵션을 매번 입력하는 것이 아닌 kubeconfig 파일로 정의하여 옵션 명시

![image.png](attachment:322872ae-9d18-4f61-beea-a5566cac59bf:image.png)

![image.png](attachment:71a893f8-4117-45d9-826f-3847be70db61:image.png)

### KubeConfig File

![image.png](attachment:1984610a-8961-4a0f-93cc-6e4870b85534:image.png)

```yaml
apiVersion: v1
kind: Config

clusters:
- name: my-kube-playground
	cluster:
		certificate-authority: ca.crt
		server: https://my-kube-playground:6443
		
contexts:
- name: my-kube-admin@my-kube-playground
	context:
		cluster: my-kube-playground
		user: my-kube-admin
		
users:
- name: my-kube-admin
	user:
		client-certificate: admin.crt
		client-key: admin.key
```

```yaml
자주 접근하는 클러스터 자격증명, 컨텍스트 추가
apiVersion: v1
kind: Config
current-context: dev-user@google
clusters:
- name: my-kube-playground (values hidden…)
- name: development
- name: production
- name: google

contexts:
- name: my-kube-admin@my-kube-playground
- name: dev-user@google
- name: prod-user@production

users:
- name: my-kube-admin
- name: admin
- name: dev-user
- name: prod-user

current-context: dev-user@google **= 기본 컨텍스트 지정**
kubectl은 늘 Google의 컨텍스트 개발자 자격증명을 이용해 Google 클러스터에 액세스
```

사용중인 현재 파일 확인

![image.png](attachment:0d41d3e1-bb95-4eb3-8911-7695b98c5a3e:image.png)

사용자 지정 config파일 사용(기본 config파일을 변경)

![image.png](attachment:179d2b6b-e844-4c3e-a362-b2a788d6aa0d:image.png)

현재 컨텍스트 업데이트

![image.png](attachment:bc2ee2c3-f4ff-4aa2-9c1c-1a2b6641b0cb:image.png)

config namespace 설정

![image.png](attachment:11a4496f-6fa1-4a9a-ad58-58248954433b:image.png)

인증서 경로는 전체 경로를 작성

![image.png](attachment:9048ed91-5ccd-4edf-9af0-9b147006f934:image.png)

인증서 경로 지정 대신 인증서내용 디코딩 하여 직접 작성 가능

```yaml
apiVersion: v1
kind: Config

clusters:
- name: production
	cluster:
		~~certificate-authority: /etc/kubernetes/pki/ca.crt~~
		certificate-authority-data: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVFVRVN..."
		
cat ca.crt | base64
echo "LS0..." | base64
```

## Practice Test - KubeConfig

```yaml
kubectl config --kubeconfig=/root/my-kube-config use-context research

my-kube-config 파일을 기본 kubeconfig 파일로 설정하고
기존 ~/.kube/config를 덮어쓰지 않고 모든 세션에서 영구적으로 만듭니다

- vi ~/.bashrc
- export KUBECONFIG=/root/my-kube-config 추가
- source ~/.bashrc

```

---

# API Groups

목적에 따라 여러 그룹으로 그룹화

/metrics, /healthz = 클러스터의 상태를 모니터

/version, /api , /apis= 클러스터의 버전 확인

/logs = 타사 로그와 통합하는데 사용

### /api, /apis

![image.png](attachment:870d8150-9f70-47e9-9c4d-c0088b95ea78:image.png)

![image.png](attachment:7e89966a-9d21-4442-9a0b-1caa8e7a78c9:image.png)

- 쿠버네티스의 모든 리소스는 다른 API그룹으로 그룹화

![image.png](attachment:878e5c16-b862-48f8-8e46-6b8789ec4182:image.png)

사용 가능한 API 그룹 확인가능

![image.png](attachment:21a4831e-9b38-48a0-b0e2-1c2464f514b8:image.png)

모든 리소스 그룹 확인가능

![image.png](attachment:bdcb7a31-7575-4585-9cc1-6756e2a731a3:image.png)

### 클러스터 API 액세스

- 인증매커니즘을 포함하여 접근해야한다.

![image.png](attachment:00a5cc07-2de7-45bf-8e4d-b7971da6d9ab:image.png)

- kubectl proxy 사용
    - 인증 매커니즘을 수동으로 지정 안해도 됌.
    - kubectl proxy실행 시 클러스터 액세스를 위해 kube config파일의 자격증명과 인증서를 사용함

![image.png](attachment:1581ef93-6e02-49f9-b8bc-559092c511ef:image.png)

- Kube proxy ≠ Kubectl proxy
    - Kube proxy = 클러스터 내 다양한 노드에 걸쳐 Pod와 service간의 연결을 가능하게 함
    - Kubectl proxy = Kubectl 유틸리티가 Kube API서버에 액세스 하기 위해 만든 것

---

# Authorization

![image.png](attachment:bec2a189-aaa7-46f7-ace1-5a9d26c109d2:image.png)

권한부여 매커니즘

![image.png](attachment:415147f0-5a41-424a-98fc-b3c6785908dd:image.png)

- Node Authorizer
    - kubelet은 시스템 노드그룹의 일부, 이름앞에 system:node를 붙여야함
    - 사용자가 요청하면 노드 승인자가 승인함으로써 권한을 부여받음

![image.png](attachment:60b725b6-2522-4b1b-9789-2fcf40fad020:image.png)

- ABAC(특성 기반 액세스 컨트롤)
    - 사용자나 사용자 그룹을 허가모음으로 연결
    - 보안을 추가 및 변경할 때마다 정책 파일을 수동으로 수정하고 Kube API서버를 다시 시작해야함
    - 관리하기 어려움

![image.png](attachment:91eb9cde-2302-4794-a2b9-dd4615594dcf:image.png)

- RBAC(역할 기반 액세스 컨트롤)
    - 권한집합을 역할로 생성
    - 해당 역할에 사용자를 연결
    - 파일 수정이 필요없이 역할 수정만 하면됌.

![image.png](attachment:74cf52d1-651a-45f0-9a40-ddc2dc9059ca:image.png)

- Webhook
    - 승인 매커니즘을 외부에 위탁(외부에서 권한관리)
    - Open Policy Agent : 타사도구로 접근통제와 승인관리

![image.png](attachment:a00b74b9-dda6-4409-b400-c1ed590a0445:image.png)

- Mode 설정
    - Kube API서버의 인증모드 옵션을 통해 설정
    - 옵션을 설정하지 않으면 AlwaysAllow로 기본설정
    - 여러 모드 설정가능

![image.png](attachment:6016ae61-b43c-4a6a-9fb4-3f0044effa63:image.png)

- 지정된 순서대로 각각의 요청을 사용할 권한이 부여
    - 사용자가 요청을 보내면 노드 승인자가 먼저 처리
    - 노드 승인자는 노드 요청만 처리하기 때문에 요청 거부
    - RBAC에서 사용자를 확인하고 승인
    - 승인 이후 사용자는 요청된 개체에 접근 가능

![image.png](attachment:c06d5032-cd2f-4ece-ac99-bed35579a78b:image.png)

---

# RBAC

1. 역할 생성
- role.yaml

![image.png](attachment:851a849c-08e4-4b0e-83d2-cfa53eecd3ec:image.png)

1. 사용자와 역할 연결
- binding.yaml

![image.png](attachment:2c14ed83-0b6f-49d8-b241-e4673d334d55:image.png)

- View RBAC

![image.png](attachment:9b657ba8-1f5f-41ab-ae3c-86516f42ca17:image.png)

![image.png](attachment:2539e4d2-559e-4728-a350-574c422f46c3:image.png)

- Check Access(권한부여 확인 명령어)
    - 클러스터의 특정 리소스에 접근

![image.png](attachment:af06dd2e-e8ac-4508-a08a-7db3644ab816:image.png)

- Resource Names
    - namespace로 접근 제한 가능

![image.png](attachment:18eae7ba-db3c-4780-ad55-c79081962419:image.png)

## Practice Test - RBAC

```yaml
k describe pod kube-apiserver-controlplane -n kube-system
k get roles -A
k get roles -n kube-system
k describe role kube-proxy -n kube-system
k get rolebindings -n kube-system
k describe rolebinding kube-proxy -n kube-system

apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: developer
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["list", "create", "delete"]
    
    
    
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: dev-user-binding
subjects:
- kind: User
  name: dev-user
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: developer
  apiGroup: rbac.authorization.k8s.io
  
  
  apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  creationTimestamp: "2025-03-30T12:02:17Z"
  name: developer
  namespace: blue
  resourceVersion: "4006"
  uid: e6377ef3-00ed-4c80-bd9f-d99ae7ad7779
rules:
- apiGroups:
  - ""
  resourceNames:
  - dark-blue-app
  resources:
  - pods
  verbs:
  - get
  - watch
  - create
  - delete
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - create

```

---

# Cluster Roles and Role Bindings

노드는 어떤 네임스페이스와도 연결될 수 없다.

![image.png](attachment:94e7b787-33e3-4dfd-aebc-98984ae30c7e:image.png)

- 노드나 영구 볼륨 같은 클러스터 간 리소스에 어떻게 사용자 권한을 부여하나?
    - clusterroles, clusterrolebindings 사용

![image.png](attachment:0c761757-b7c6-453c-bd44-5389382cfec5:image.png)

![image.png](attachment:0aada7b7-8980-475e-be28-996a52e6d6c7:image.png)

## Practice Test - RBAC

```yaml
k get clusterrole
k get clusterrolebindings
k describe clusterrole cluster-admin
k describe clusterrolebinding cluster-admin

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: node-admin
rules:
- apiGroups: [""]
  resources: ["nodes"]
  verbs: ["create"]

  
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: michelle-binding
subjects:
- kind: User
  name: michelle
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: node-admin
  apiGroup: rbac.authorization.k8s.io
  
---
  apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: node-admin
rules:
- apiGroups: [""]
  resources: ["nodes"]
  verbs: ["get", "watch", "list", "create", "delete"]
  
  
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: michelle-binding
subjects:
- kind: User
  name: michelle
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: node-admin
  apiGroup: rbac.authorization.k8s.io
  
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: storage-admin
rules:
- apiGroups: [""]
  resources: ["persistentvolumes"]
  verbs: ["get", "watch", "list", "create", "delete"]
- apiGroups: ["storage.k8s.io"]
  resources: ["storageclasses"]
  verbs: ["get", "watch", "list", "create", "delete"]
  
  kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: michelle-storage-admin
subjects:
- kind: User
  name: michelle
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: storage-admin
  apiGroup: rbac.authorization.k8s.io
```

---

# Service Accounts

쿠버네티스 계정은 2종류

- 사용자 계정
    - Admin, Developer
- 서비스 계정
    - Prometheus, Jenkins

서비스 계정 생성

- kubectl create serviceaccount dashboard-sa
- kubectl get serviceaccount
- kubectl describe serviceaccount dashboard-sa

계정을 위한 토큰 생성

- secret에 token 정보저장
- kubectl describe secret dashboard-sa-token-kbbdm

쿠버네티스 클러스터 자체 배포

- 서비스 토큰 secret을 타사 응용프로그램의 호스팅 Pod 내 볼륨으로 자동으로 마운트 함으로 써 간단해진다.
- 쿠버네티스 API에 접속할 수 있는 토큰이 이미 Pod안에 있어서 맵이 토근을 쉽게 읽을 수 있음
- kubectl get serviceaccount

파드 생성시 secret 볼륨이 자동으로 생성

- kubectl describe pod my-kubernetes-dashboard

다른 서비스 계정 사용

- Pod 정의 파일 수정
- 기존 Pod의 서비스 계정은 수정할 수 없다 / 삭제 후 다시 만들어야함
- deployment는 수정가능

서비스 계정을 자동으로 마운트하지 않게 설정 가능

- default는 자동으로 마운트

### 1.22/1.24 Update

서비스 토큰 decode= JWT 변환

기존의 보안문제점

- 만료 날짜가 설정되어있지 않음
- 서비스 계정이 존재하는 한 JWT는 유효함
- 각 JWT는 서비스 계정 당 별개의 비밀개체를 필요로 함으로 확장성 문제가 발생

### v1.22

- KEP 1205 - Bound Service Account Tokens
    - 더 안전한 Request API
- 서비스 계정에 종속되지 않음

### V1.24

- KEP-2799: Reduction of Secret-based Service Account Tokens
- 서비스 계정 생성시 토큰 생성하지 않음
- 토큰을 생성하려면 명령어 실행
: kubectl create token dashboard-sa
- 만료 날짜가 지정됌
- secret 생성으로 만료되지 않는 토큰 또한 생성 가능

## Practice Test - Service Accounts

```yaml
k get serviceaccount
k describe serviceaccount default
kubectl create serviceaccount dashboard-sa

apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  creationTimestamp: "2025-03-31T12:01:14Z"
  generation: 1
  name: web-dashboard
  namespace: default
  resourceVersion: "836"
  uid: efc8276f-801f-4042-a938-2aa9968d846e
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      name: web-dashboard
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        name: web-dashboard
    **spec:
      serviceAccountName: dashboard-sa**
      containers:
```

---

# Image Security

![image.png](attachment:43cc45c0-9eb9-4319-8f21-0b8fd32598f4:image.png)

### Private Repository

- 개인 레지스트리에 로그인
    - docker login private-registry.io
- 개인 레지스트리의 이미지를 이용해 애플리케이션을 실행
    - docker run private registry.10/app/internal-app
    

nginx-pod.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: nginx-pods
spec:
	containers:
	- name: nginx
		image: private-registry.io/app/internal-app
	imagePullSecrets:
	- name: regcred
```

worker node에서 자격증명을 docker 실행시간에 어떻게 전달하나

- 자격증명이 포함된 비밀개체 생성

```bash
kubectl create secret docker-registy regcred \
--docker-server= private-registry.io  \
--docker-username= registry-user   \
--docker-password= registry-password   \
--docker-email= registry-user@org.com   \
```

도커 레지스트리의 secret은 regcred

## Practice Test - Image Security

```bash
spec:
      containers:
      - image: myprivateregistry.com:5000/nginx:alpine
        imagePullPolicy: IfNotPresent
        name: nginx
        
        
        
kubectl create secret docker-registry private-reg-cred --docker-username=dock_user --docker-password=dock_password --docker-server=myprivateregistry.com:5000 --docker-email=dock_user@myprivateregistry.com
imagePullSecrets:
      - name: private-reg-cred
```

---

# Security in Docker

- 컨테이너에 의해 실행되는 모든 프로세스는 사실 고유의 네임스페이스에서 실행된다.
- 다른 네임스페이스에서는 확인 불가능
- 프로세스는 다른 네임스페이스에 다른 프로세스 ID를 가질 수 있다
    - Docker가 시스템 내에서 컨테이너를 격리하는 방식
- 기본값으로 Docker는 root사용자로서 컨터이너 내의 프로세스 실행
- root가 아닌 다른 사용자로프로세스 실행
    1. 사용자 ID 명시
    - docker run —user-1000 ubuntu sleep 3600
    - ps aux
    1. Docker 이미지 자체에서 정의

- Docker는 컨테이너 내 root 사용자의 기능을 제한 하는 보안 기능 집합을 구현
- 컨테이너 안의 root사용자는 Host의 root 사용자와 다름
- 컨테이너 안의 root사용자 권한 추가
    - docker run --cap-add MAC_ADMIN ubuntu
    - docker run --cap-drop KILL ubuntu
    - docker run --privileged ubuntu

---

# SECURITY CONTEXTS

Container 레벨 or Pod 레벨 에서 보안 설정 가능

- Pod레벨에서 설정하면 Pod안의 모든 Container에 설정 적용
- Pod와 Container 레벨 둘 다 구성하면 Container설정이 Pod의 설정을 무효화

Pod 설정

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: web-pod
spec:
	securityContext:
		runAsUser: 1000
		
containers:
	- name: ubuntu
		image: ubuntu
		command: ["sleep", "3600"]
```

Container설정

-securityContext설정이 Container 설정 밑으로

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: web-pod
spec:
	containers:
		- name: ubuntu
			image: ubuntu
			command: ["sleep", "3600"]
			securityContext:
				runAsUser: 1000
				capabilities:
					add: ["MAC_ADMIN"]
```

## Practice Test - SECURITY CONTEXTS

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: ubuntu-sleeper
  namespace: default
spec:
  securityContext:
    runAsUser: 1010
  containers:
  - command:
    - sleep
    - "4800"
    image: ubuntu
    name: ubuntu-sleeper
    
    
    ---
apiVersion: v1
kind: Pod
metadata:
  name: ubuntu-sleeper
  namespace: default
spec:
  containers:
  - command:
    - sleep
    - "4800"
    image: ubuntu
    name: ubuntu-sleeper
    securityContext:
      capabilities:
        add: ["SYS_TIME"]
```

---

# Network Policies

![image.png](attachment:65cbe507-2a5d-4937-80ce-60b0f6589348:image.png)

![image.png](attachment:d692f9dd-c37b-4e49-bd1f-31876ecc23f7:image.png)

### 쿠버네티스 네트워킹의 전제조건

- 어떤 솔루션을 구현하든 Pod가 서로 통신할 수 있어야 한다.
- 기본값으로 모든 트래픽이 허용됌

![image.png](attachment:2fafc3e9-c108-43b6-a656-927bad3faf4f:image.png)

### 네트워크 정책

![image.png](attachment:54b0691d-d14f-4f15-8117-be32c668f6f9:image.png)

![image.png](attachment:bd039c42-1ac6-48e2-b02c-0102e6fb454e:image.png)

![image.png](attachment:606bc545-0455-4c90-84b9-06368140a4f6:image.png)

Allow Ingress Traffic From API Pod on Port 3306

```yaml
policy-definition.yaml

apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
	name: db-policy
spec:
	podSelector:
		matchLabels:
			role: db
	policyTypes:
	- Ingress
	ingress:
	- from:
		- podSelector:
				matchLabels:
					name: api-pod
		ports:
		- protocol: TCP
			port: 3306
```

- kubectl create -f policy-definition.yaml

---

# Developing Network Policies

- DB-Pod 3306 Allow From API Pod

![image.png](attachment:43f61a9a-6f9c-4a4c-96f8-b89fc7008b14:image.png)

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
	name: db-policy
spec:
	podSelector:
		matchLabels:
			role: db
	policyTypes:
	- Ingress
	ingress:
	- from:
		- podSelector:
				matchLabels:
					name: api-pod
		ports:
			- protocol: TCP
				port: 3306
```

- 클러스터에 API Pod가 여러개 있다면?

![image.png](attachment:7a76dbee-ce4e-49e1-800b-058675b06b2f:image.png)

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
	name: db-policy
spec:
	podSelector:
		matchLabels:
			role: db
	policyTypes:
	- Ingress
	ingress:
	- from:
		- **podSelector**:
				matchLabels:
					name: api-pod
			**namespaceSelector**:
					matchLabels:
						name: prod
		ports:
			- protocol: TCP
				port: 3306
				
	podSelector 없이 namespaceSelector만 있다면
	= 지정된 namespace내의 모든 Pod는 DB Pod에 갈 수 있다.
```

- DB Pod → Backup Server ipBlock(192.168.5.10)

![image.png](attachment:7796974e-71ed-44f6-92a2-261c090d91f4:image.png)

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
	name: db-policy
spec:
	podSelector:
		matchLabels:
			role: db
	policyTypes:
	- Ingress
	ingress:
	- from:
		- podSelector:
				matchLabels:
					name: api-pod
			namespaceSelector:
					matchLabels:
						name: prod
		- **ipBlock**:
				cidr: 192.168.5.10/32
		ports:
			- protocol: TCP
				port: 3306
	
```

- podSelector와 namespaceSelector를 같이 둔다.
    - 두 정책 모두 만족해야함(AND)
- podSelector와 namespaceSelector를 따로 둔다.
    - 각 정책 별로 만족하면 통과(OR)

![image.png](attachment:17e4207d-332a-4fd6-ad2c-83f7740ab7fe:image.png)

- DB-Pod → Backup Server Egress(80)

![image.png](attachment:8dc89a3c-0632-4db4-9b42-336e20dd0c93:image.png)

```yaml
spec:
	podSelector:
		matchLabels:
			role: db
	policyTypes:
	- Ingress
	**- Egress**
	ingress:
	- from:
		- podSelector:
				matchLabels:
					name: api-pod
		ports:
		- protocol: TCP
			port: 3306
	**egress:
	- to:
		- ipBlock:
				cidr: 192.168.5.10/32
		ports:
		- protocol: TCP
			port: 80**
```

## Practice Test - Network Policies

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: internal-policy
  namespace: default
spec:
  podSelector:
    matchLabels:
      name: internal
  policyTypes:
  - Egress
  - Ingress
  ingress:
    - {}
  egress:
  - to:
    - podSelector:
        matchLabels:
          name: mysql
    ports:
    - protocol: TCP
      port: 3306

  - to:
    - podSelector:
        matchLabels:
          name: payroll
    ports:
    - protocol: TCP
      port: 8080

  - ports:
    - port: 53
      protocol: UDP
    - port: 53
      protocol: TCP
```