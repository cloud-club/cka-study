# Security Promitives

kube-apiserver → kube-control 유틸리티와 api에 직접 엑세스함으로써 상호작용

⇒ 1차 방어선; API 서버 자체에 대한 액세스 제어

- **who can access?**
    - Authenticate(인증) 매커니즘에 의해 정의
    - ex. password, tokens, certificates, 외부 인증 공급자, Service Account
- **what can they do?**
    - Authorization(권한) 매커니즘에 의해 정의
    - ex. RBAC Authorization, ABAC Authorization, Node Authorization, Webhook Mode

- TLS Certificates
    - K8s를 구성하는 Component → TLS 암호화로 보안
- 클러스터 내 응용 프로그램 간 통신 → 네트워크 정책을 통해 엑세스 권한 제한

# Authentication

최종 사용자의 보안은 응용 프로그램 자체에 의해 내부적으로 관리 → 제외

초점: 관리 목적으로 사용자가 클러스터에 액세스

⇒ User(Admins, Developers) & Service Accounts(Bots)

- User
    - Kubernetes는 사용자 계정 직접 관리 X → 외부 소스 의존
- Service Accounts
    - Kubernetes가 관리 O
    
    ```bash
    kubectl create sa sa1
    
    kubectl get sa
    ```
    

### User

모든 사용자 엑세스는 API 서버에 의해 관리

kube-apiserver는 요청 처리하기 전에 인증

- 인증 매커니즘
    - Static Password File
    - Static Token File
    - Certificates
    - Identify Service

- Auth Mechanisms - Basic
    - ❗️Kubernetes 1.19 이후, 기본 인증(Basic Authentication)은 폐기. Bearer Token 인증은 유지❗️
    - csv 파일에 <암호>,<사용자>,<uid>,<그룹 정보>
        - token의 경우 암호 대신 token 저장
    - kube-apiserver.service 또는 /etc/kubernetes/manifests/kube-apiserver.yaml에 추가 → 자동으로 업데이트
        - `—-basic-auth-file=user-detail.csv`
        - `—-token-auth-file=user-token-details.csv`
    
    ```bash
    curl -v -k https://master-node-ip:6443/api/v1/pods -u "user1:password123"
    
    curl -v -k https://master-node-ip:6443/api/v1/pods --header "Authorization: Bearer ~~~""
    ```
    

⇒ 불안정하기에 권장 X

https://kubernetes.io/docs/reference/access-authn-authz/authentication/#authentication-strategies

# TLS Certificates

### TLS Certifiacate란?

TLS 인증서는 사용자-서버 사이의 통신이 암호화되도록 

https://babbab2.tistory.com/4

- 대칭키 암호화
    - 암호화를 하는 키와 복호화를 하는 키가 동일한 방식
        
        ⇒ 해커가 데이터 해독할 위험 O
        
- 비대칭키 암호화
    - 한 쌍의 키 → 공개 키, 개인 키
        - ex. Private key, Public key(=Public lock)
    - 여러 서버 보호 → Public lock 복사본 만들어 원하는 만큼 서버에 두기
    - 추가 공용 Lock 생성 가능
    
- **TLS 암호화 방식**
    - 대칭키, 비대칭키 암호화 방식 함께 사용
    - 처음에 대칭키를 서로 공유하는 통신을 RSA 비대칭키 방식을 이용하고, 실제 통신을 할 때는 CPU 리소스 소모가 적은 대칭키 방식으로 데이터를 주고 받음
    
    ```bash
    openssl genrsa -out my-bank.key 1024
    # my-back.key
    
    openssl rsa -in my-bank.key -pubout > myback.pem
    # my-back.key mybank.pem
    ```
    
    ⇒ 해커가 해독할 위험 없이 안전하게 사용 가능
    
    but, 웹사이트를 위조한다면 ? ⇒ 인증서 !
    
- **Certificate**
    - 서버가 신뢰할 수 있는 진짜 서버 임을 확인하기 위해 필요한 것
    - 브라우저 → 인증서 유효성 검사 매커니즘으로 기본 제공. 서버에서 받은 인증서가 합법적인지 유효성 확인
    - CA(Certificate Authority)
        - 인증서에 서명하고 유효성을 확인해주는 알려진 기관
        - 방식
            - Certificate Signing Request (CSR)
            
            ```bash
            openssl req -new -key my-bank.key -out my-bank.csr \
            	-subj "/C=US/ST=CA/O=MyOrg, Inc./CN=my-bank.com"
            # my-bank.jey my-bank.csr
            ```
            
            - Validate Information
            - Sign and Send Certificate
    - CA가 어떻게 유효한지?
        - Public Key, Private Key
        - Public Key는 브라우저에 기본 제공
            - but, 개별적으로 호스팅된 사이트의 유효성 검사 X
                
                ⇒ 개인 CAs를 호스트 가능
                
                : 회사 내에서 내부적으로 배포할 수 있는 CA 서버
                
    - PKI(Public Key Infrastructure)
    - Naming Convention
        - Public Key → `*.crt`, `*.pem`
        - Private Key → `*.key`, `*-key.pem`

### TLS Certificates for Cluster Components

- Certificates 유형
    - Server Certificates
    - Root Certificates
    - Client Certificates

- Kubernetes는 노드들 사이 모든 통신 보안 및 암호화 필요
- Kubernetes 클러스터 내 모든 구성 요소 간의 통신도 보안 필요

⇒ Server Certificates for Servers, Client Certificates for Clients

![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image.png)

![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image%201.png)

- **Server Certificates for Servers**
    - **kube-apiserver**
        - https 서비스 → 클라이언트와의 모든 통신을 보호하기 위해 인증서 필요
        - `apiserver.crt`, `apiserver.key`
    - **etcd server**
        - 클러스터에 관한 모든 정보 저장 → 인증서 한 쌍과 자체 키 필요
        - `etcdserver.crt`, `etcdserver.key`
    - **kubelet server**
        - HTTPS API 엔드포인트 공개 → 인증서, 키페어 필요
        - `kubelet.crt`, `kubelet.key`

- **Client Certificates for Clients**
    - admin
        - kube-apiserver에 엑세스하는 클라이언트
        - `admin.crt`, `admin.key`
    - **kube-scheduler**
        - kube-apiserver에 엑세스하는 클라이언트
        - `scheduler.crt`, `scheduler.key`
    - **kube-controller-manager**
        - kube-apiserver에 엑세스하는 클라이언트
        - `controller-manager.crt`, `controller-manager.key`
    - **kube-proxy**
        - `kube-proxy.crt`, `kube-proxy.key`
    - **kube-apiserver**
        - 모든 구성 요소들과 통신 → 클라이언트
        - 같은 키 사용 가능
            - `apiserver.crt`, `apiserver.key`
        - or 새로운 인증서 한 쌍 생성해 etcd 서버에 인증
            - `apiserver-etcd-client.crt`, `apiserver-etcd-client.key`
            - `apiserver-kubelet-client.crt`, `apiserver-kubelet-client.key`

- **CA**
    - CA 자체에 인증서 및 키 존재
    - `ca.crt`, `ca.key`

### Generate Certificates

EASYRSA, OPENSSL, CFSSL 등 다양한 도구 이용 가능

강의 → OpenSSL 도구 사용

- **CA Certificates**

```bash
# Generate Keys
openssl genrsa -out ca.key 2048
#ca.key

# Certificate Signing Request
openssl req -new -key ca.key -subj \
	"/CN=KUBERNETES-CA" -out ca.csr
#ca.csr
	
# Sign Certificates
openssl x509 -req -in ca.csr -signkey ca.key -out ca.crt
#ca.crt
```

- **Client Certificates**
    - **Admin User**
        - 기본 사용자와 관리자 구분하기 위해 CSR에 그룹 정보 추가
    
    ```bash
    # Generate Keys
    openssl genrsa -out admin.key 2048
    #admin.key
    
    # Certificate Signing Request
    openssl req -new -key admin.key -subj \
    	"/CN=kube-admin/O=system:masters" -out admin.csr
    #admin.csr
    	
    # Sign Certificates
    openssl x509 -req -in admin.csr -CA ca.crt -CAkey ca.key \
    	-out admin.crt
    #admin.crt
    ```
    
    - **Kube Scheduler**
        - Control plane의 일부이므로 앞에 `SYSTEM` 키워드
    - **Kube Controller Manager**
        - Control plane의 일부이므로 앞에 `SYSTEM` 키워드
    - **Kube Proxy**
        - Control plane의 일부이므로 앞에 `SYSTEM` 키워드
    - **Kubelet**
        
        → 아래에 설명..
        

클라이언트가 서버에서 보낸 인증서를 유효성 검사하기 위해서는(또는 그 반대의 경우) 

CA의 public certificate의 복제본이 필요! 

→ 웹에서는 기본적으로 브라우저에 깔려있음

Kubernetes의 다양한 컴포넌트가 서로 검증하려면, 모두 CA의 root certificate 복제본이 필요

- **증명서 사용하는 법**
    - REST API 호출에서 인증서 사용
    
    ```bash
    curl https://kube-apiserver:6443/api/v1/pods \
    	--key admin.key --cert admin.crt
    	--cacert ca.crt
    ```
    
    - kubeconfig 사용
    
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
    

- **Server Certificates**
    - **ETCD Server**
        - ETCD 서버는 다중으로 구성될 수 있음 → 추가적인 인증서 생성(peer certificate)
        
        ```bash
        cat etcd.yaml
        ```
        
    - **Kube API Server**
        - 다양한 alias 존재
            
            ex. kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster.local, 10.96.0.1, …
            
            ⇒ 이 모든 이름은 반드시 kube-apiserver 인증서에 존재해야함
            
        
        ```bash
        openssl genrsa -out apiserver.key 2048
        #apiserver.key
        
        openssl req -new -key apiserver.key -subj \
        	"/CN=kube-apiserver" -out apiserver.csr
        #apiserver.csr
        ```
        
        ```bash
        # openssl.cnf
        ...
        [alt_names]
        DNS.1 = kubernetes
        DNS.2 = kubernetes.default
        ...
        ```
        
        ```bash
        openssl x509 -req -in apiserver.csr \
        	-CA ca.crt -CAkey ca.key -out apiserver.crt
        #apiserver.crt
        ```
        
        - **키 어디에 지정?**
            
            ![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image%202.png)
            
            - kube-apiserver가 사용하는 클라이언트 인증서 고려
            - kube-apiserver로 인증서 전달
                - CA 파일 전달되어야 함!
                - kube-apiserver 인증서를 tls 인증 옵션에 제공
                - kube-apiserver에서 사용하는 client 인증서를 지정해 etcd server에 ca 파일로 다시 연결
                - kubelets 연결하기 위해 kube-apiserver client 인증서 설정
    - **Kubelet (Server Cert)**
        - ACTPS API Server로 각 노드에서 실행
        
        ⇒ 각 노드에 대해 키-인증서 쌍 필요
        
        이때, 인증서 이름은 노드 이름 따서 
        
        ![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image%203.png)
        

- **Kubelet (Client Cert)**
    - kube-apiserver와 통신할때 사용
    - certificates name → node는 시스템 구성 요소 → `SYSTEM` 키워드 필요
    - `SYSTEM:NODES` 라는 Group 추가
    
    ![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image%204.png)
    

### View Certificate Details

- Kubernetes 클러스터 생성 방법
    - 처음부터 네이티브 배포
    - kubeadm 툴 의존 → 구성요소가 pod로 배포

⇒ 강의에서는 kubeadm에 의한 클러스터 프로비저닝

```bash
cat /etc/kubernetes/manifests/kube-apiserver.yaml
```

각 인증서 내부 살펴 해당 인증서에 관해 상세히 알아보기

ex. `apiserver.crt` 

```bash
openssl x509 -in /etc/kubernetes/pki/apiserver.crt -text -noout
```

![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image%205.png)

![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image%206.png)

- **View Logs**
    - Native로 설정한 경우 → 운영체제 로깅 기능 이용한 서비스 로그 확인
        
        ```bash
        journalctl -u etcd.service -l
        ```
        
    - kubeadm의 경우
        
        ```bash
        k logs etcd-master
        ```
        
        if. kube-apiserver나 NCD 서버 같은 핵심 구성 요소 다운된 경우 kubectl 명령어 동작 X
        
        ⇒ Docker로 내려가 로그 가져오기
        
        ```bash
        docker ps -a
        
        docker logs <container id>
        ```
        

https://github.com/mmumshad/kubernetes-the-hard-way/tree/master/tools

### Certificate Workflow & API

- CA: 생성한 키와 인증서 파일 한 쌍
- CA 서버: CA가 저장된 서버 → Kubernetes에서는 보통 master node

- **Kubernetes built-in Certificates API**
    1. Create CSR Object
        1. Kubernetes object로 생성됨!
    2. Review Requests
    3. Approve Requests
    4. Share Certs to Users

```bash
openssl genrsa -out jane.key 2048
#jane.key

openssl req -new -key jane.key -subj "/CN=jane" -out jane.csr
#jane.csr
```

```yaml
# jane-csr.yaml
apiVersion: certifiactes.k8s.io/v1
kind: CertificateSigningRequest
metadata:
	name: jane
spec:
	expirationSeconds: 600 # seconds
	usages:
		- digital signature
		- key encipherment
		- server auth
	request: # 사용자가 보낸 CSR -> base64를 통해 암호화
			
```

```bash
cat jane.csr | base64
echo "..." | base64 --decode

kubectl get csr

kubectl certificate approve jane
kubectl certificate deny jane
```

⇒ 인증서는 최종 사용자와 공유 가능

- 인증서 관련 작업
    
    ⇒ kube-controller-manager이 인증서 관련 작업
    
    - CSR-APPROVING
    - CSR-SIGNING
    
    ![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image%207.png)
    

# Security KUBECONFIG

kube-apiserver가 사용자를 인증하기 위해 curl 혹은 kubectl 명령어 사용

⇒ 매번 certificate 옵션 입력하기 어려움

⇒ kubeconfig 구성 파일 사용해서 옵션으로 명시!

`$HOME/.kube/config` 

```bash
curl https://my-kube-playground:6443/api/v1/pods \
	--key admin.key
	--cert admin.crt
	--cacert ca.crt
```

```bash
kubectl get pods \
	--server my-kube-playground:6443 \
	--client-key admin.key \
	--client-certificate admin.crt \
	--certificate-authority ca.crt
```

- 구성 요소
    - **Cluster**
        - 액세스 필요한 다양한 kubernetes 클러스터 정의
        - ex. dev, prod, google
        - 사용자 인증 → 서버 사양 부분 (server, certificate-authority)
    - **Contexts**
        - 어떤 사용자 계정이 어떤 클러스터에 액세스하기 위해 사용될지 정의
        - ex. admin@prod, dev@google
        - 사용자 인증 → 해당 클러스터와 해당 사용자 연결
    - **Users**
        - 클러스터에 액세스 권한이 있는 사용자 계정 정의
        - ex. admin, dev user, prod user, …
        - 사용자 인증 → 관리자 키와 인증서 (client-key, client-certificate)

- KubeConfig File
    - yaml 형태

```yaml
apiVersion: v1

kind: Config
current-context: dev-user@google

clusters:
	- name: my-kube-playgroud
		cluster:
			# certificate-authority: /etc/kubernetes/pki/ca.crt
			certificate-authority-data: | # ca.crt 파일을 base64로 변환해서!
				...
			server: https://my-kube-playground:6443
	- name: google
		...
contexts:
	- name: my-kube-admin@my-kube-playgroud
		context:
			cluster: my-kube-playgroud
			user: my-kube-admin
	- name: dev-user@google
		context:
			cluster: google
			user: dev-user
			namespace: test # 특정 네임스페이스 자동 설정
		...
users:
	- name: my-kube-admin
		user: 
			client-certificate: /etc/kubernetes/pki/users/admin.crt
			client-key: /etc/kubernetes/pki/users/admin.key
	- name: dev-user
		...
```

⇒ kubectl create 하지 않아도 됨! kubectl 커맨드로 자동으로 읽힘

```bash
kubectl config view

kubectl config view --kubeconfig=my-custom-config

# 현재 컨텍스트 업데이트
kubectl config use-context my-kube-admin@my-kube-playground
```

- Certificates

# API Groups

```bash
curl https://kube-master:6443/version
curl https://kube-master:6443/api/v1/pods
```

- `/version`
    - 클러스터 버전 보기 위함
- `/metrics` , `/healthz`
    - 클러스터 상태 모니터
- `/logs`
    - 서드파티 log application과 통합

- `/api`  → core
    - /v1 하에 모든 핵심 기능 존재
    - ex. namespace, pods, rc, events, endpoints, nodes, bindings, …
    
    ![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image%208.png)
    
- `/apis` → named
    - 좀 더 조직화되고 새로운 기능
    - /apps, /extensions, /networking.k8s.io, /storage.k8s.io, /authentication.k8s.io, /certificates.k8s.io
    
    ![image.png](7%20Security%201c6766b7e9a4800889acdc539002cd85/image%209.png)
    
    ⇒ Verbs
    

- **클러스터 API 액세스**
    - curl
    
    ```bash
    curl http://localhost:6443 -k \
    	--key admin.key
    	--cert admin.crt
    	--cacert ca.crt
    #사용 가능한 API 그룹 리스트 나옴
    ```
    
    - kubectl proxy
        - ACTP 프록시 서비스
        - kubectl 유틸리티가 kube-apiserver에 엑세스하기 위해 만듦
    
    ```bash
    # kubeconfig의 자격 증명과 인증서 사용
    kubectl proxy 
    #Starting to serve on 127.0.0.1:8001
    
    curl http://localhost:8001 -k
    ```
    

# Authorization

클러스터를 서로 다른 조직이나 팀으로 공유할때, 논리적으로 네임스페이스를 이용해 분할하여 사용자에 대한 접근을 제한할 수 있음

- Authorization Mode
    - Node, ABAC, RBAC, Webhook, AlwayAllow, AlwaysDeny

- **Node**
    - 클러스터 내 액세스
    - `system:node:node01`이라는 시스템 노드 이름과 `SYSTEM:NODES`라는 시스템 노드 그룹의 일부로 사용자가 요청 → 노드 승인자가 승인하고 권한 부여
- **ABAC**
    - Attribute-Based Access Control
    - 사용자나 사용자 그룹을 set of Permissions로 연결
    - JSON 형식으로 policy 집합과 파일 생성
    - 보안 추가하거나 변경할때마다  이 정책 파일 수동으로 수정하고 kube-apiserver 다시 시작
    
    ⇒ 관리하기 어려움..
    
- **RBAC**
    - Role-Based Access Controls
    - 요구하는 권한 집합으로 역할 생성 → 유저와 연결
- **Webhook**
    - 외부에서 권한 관리하고 싶은 경우
- **AlwayAllow, AlwaysDeny**

⇒ /etc/kubernetes/manifest/kube-apiserver 파일의 `—-authorizatoin-mode` 를 통해 설정 (default: AlwaysAllow)

if. 다양한 모드로 설정한 경우,

`—-authorization-mode=Node,RBAC,Webhook`

지정된 순서대로 각각 요청 사용할 권한 부여

### RBAC

Role을 만들고 RoleBinding으로 사용자를 그 Role에 링크

- `developer-role.yaml`

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
	name: developer
rules:
	- apiGroups: [""] # core group인 경우 비워도 됨
		resources: ["pods"] # 제공할 리소스
		verbs: ["list", "get", "create", "update", "delete"]
		resourceNames: ["blue", "pink"]# 특정 리소스만 액세스
		
	- apiGroups: [""] 
		resources: ["configMap"]
		verbs: ["create"]
```

```bash
k create role developer --verb=list,get,delete \
	--resource=pods
```

```bash
kubectl create -f developer-role.yaml
```

- `devuser-developer-binding.yaml`

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
	name: dev-user-binding
subjects: # 사용자 세부 정보 지정
	- kind: User
		name: dev-user
		apiGroup: rbac.authorization.k8s.io
roleRef: # 우리가 만든 롤 세부사항 제공
	kind: Role
	name: developer
	apiGroup: rbac.authorization.k8s.io
```

```bash
k create rolebinding dev-user-binding --role=developer \
	--user=dev-user
```

```bash
kubectl create -f dev-user-binding.yaml
```

```bash
kubectl get roles
kubectl get rolebindings

kubectl describe role developer
kubectl describe rolebinding dev-user-binding
```

if. 사용자가 클러스터의 특정 리소스에 접근하고 싶다면?

```bash
kubectl auth can-i create deployment
#yes

kubectl auth can-i delete nodes
#no
```

if. 관리자인 경우 다른 유저의 권한 확인 가능

```bash
kubectl auth can-i create deployments --as dev-user
#no

kubectl auth can-i create pods --as dev-user
#yes
```

# Cluster Roles

Role과 Rolebinding은 namespace scoped

but, 노드는 클러스터 범위 리소스

⇒ nodes, PV, clusterroles, clusterrolebinding, csr, namespace

```bash
# cluster scoped 설정
kubectl api-resources --namespaced=false
```

- `cluster-admin-role.yaml`

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
	name: cluster-administrator
rules:
- apiGroups: [""]
	resources: ["nodes"]
	verbs: ["list", "get"]
```

- `cluster-admin-role-binding.yaml`

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
	name: cluster-admin-role-binding
subjects:
	- kind: User
		name: cluster-admin
		apiGroup: rbac.authorization.k8s.io
roleRef:
	kind: ClusterRole
	name: cluster-administrator
	apiGroup: rbac.authorization.k8s.io
```

cluster role과 binding은 cluster scoped 라고 했지만,

namespaced resource를 위한 cluster role도 생성 가능!

# Service Account

Kubernetes account → User Account & Service Account

- User → 사용자가 사용
- Service → 기계가 사용

응용 프로그램이 kubernetes api를 쿼리하려면 인증되어야 함 → SA 사용
