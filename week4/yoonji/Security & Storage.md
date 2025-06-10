# Security
## Authentication 
### 클러스터에 액세스할 수 있는 사용자 유형
#### 1. Application End Users: applications 자체에 의해 내부적으로 관리됨
#### 2. Users (Admins, Developers)
쿠버네티스가 직접 관리 x. 인증서가 있는 파일이나 LDAP같은 타사 ID 서비스가 이런 사용자들을 관리
- 모든 Users의 요청은 API 서버로 가며, 관리된다.
- API 서버는 요청 처리 전 인증을 함: 고정 암호 파일 이용, 인증서 이용, LDAP 또는 Kerberos에 연결
  - 고정 암호와 토큰 파일: CSV에 사용자 목록과 암호를 생성 후, api 서버를 다시 시작하고 수정한다. 특히 kubeadm은 파일을 업데이트하면, 자동으로 kube api 서버를 재시작한다.  
#### 3. Service Accounts (Bots)
쿠버네티스가 관리. 쿠버네티스 api를 통해 service accounts를 생성하고 관리.


## TLS in Kubernetes
클러스터를 위한 인증서 생성 방법: EASYRSA, OPENSSL, CFSSL 등
- 강의에서는 OpenSSL을 이용해 인증서를 생성한다.

### 1. CA(Certificate Autority) 인증서
1. `openssl genrsa -out ca.key 2048` 으로 프라이빗 키 생성 -> `ca.key`
2. `openssl req -new -key ca.key -subj`으로 인증 요청 -> `ca.csr`
3. `openssl x509 -req in ca.csr -signkey ca.key -out ca.crt`로 인증서 서명 -> `ca.crt`

### 2. 클라이언트 인증서 생성 
위 과정과 유사하며, 위 증명서로 클러스터를 관리한다. 
- ex) 관리자 인증서의 경우, kube api 서버에 만드는 rest api 호출에서 사용자 이름과 암호 대신, 이 인증서를 사용할 수 있다.
  - 이때, 키, 인증서, CA인증서를 옵션으로 지정하는데, 이 매개변수들을 kubeconfig로 옮길수도 있다. 

### 3. 서버 인증서 생성
쿠버네티스에서도 구성 요소끼리 서로를 확인하려면, CA의 루트 인증서 복사본이 필요하다. 

#### ETCD Server
etcd는 대개 다중 서버에 걸쳐서 배포되므로, 추가 peer 인증서가 필요하며, 과정은 위와 유사하다. 

#### KUBE-API Server
kube-api 서버를 다양한 이름(`.default` 등)으로 참조하므로, 유효한 연결이 필요하다. 
1. CA 파일 전달
2. api 서버 인증서를 TLS 옵션 아래에 제공
3. api 서버에서 사용되는 클라이언트 인증서를 지정해 기타 서버에 ca 파일로 다시 연결

#### Kubelet Server
인증서의 이름은 노드 이름을 따서 지으며, 인증서가 생성되면 `kubelet-config.yaml`에 사용한다. 

---
# Storage
## Storage in Docker
> Docker가 데이터를 어디에, 어떻게 저장하고 컨테이너 파일 시스템을 어떻게 관리하는가?
1. Docker를 시스템에 설치하면, `/var/lib/docker`에 저장되며, docker가 기본값으로 모든 데이터를 저장하는 곳이다. 이때 데이터란, docker 호스트에서 실행되는 이미지 및 컨테이너와 관련된 파일을 의미한다.
2. docker 컨테이너에 의해 생성된 볼륨은 volume 폴더 아래 생성된다.


> docker의 이미지와 컨테이너 파일은 어떻게 저장하는가?
1. 도커는 이미지를 레이어드 아키텍처로 구축한다. 도커파일의 각 줄은 도커 이미지에 새 레이어를 만드며, 이전 레이어에서 변경된 것만 수정된다. 도커는 이를 통해 이미지를 빠르게 생성하고, 디스크 공간을 효율적으로 저장한다. 
2. application도 동일하게 적용 가능. docker는 이전 계층을 캐시에서 재사용해 application 이미지를 구축한다.


> 동일한 이미지 레이어는 해당 이미지로부터 생성된 여러 컨테이너 사이에 공유될 수 있다.
- **copy-on-write**: 도커파일을 저장하기 전에, 읽기 쓰기 계층의 파일 복사본을 자동으로 생성하고, 읽기 쓰기 계층의 다른 파일 버전을 수정한다.


## Volume
> pod가 사용하는 임시 저장 공간으로, 수명주기도 파드와 동일하다. 컨테이너의 데이터를 영구적으로 유지하고 싶다면, 컨테이너에 **Persistent Volume**을 추가한다.
1. `docker volume create data_volume` -> `/var/lib/docker/volumes`에 `data_volume` 디렉토리 생성
2. `docker run -V data_volume:/var/`: 컨테이너를 실행할 때 이 볼륨을 컨테이너 내부에 마운트. `:`으로 컨테이너 내의 위치 지정.
3. 모든 데이터는 볼륨에 저장되며, 데이터는 유지된다.

## Container Runtime Interface (CRI)
- 과거: 쿠버네티스는 도커를 컨터이너 런타임 엔진으로 사용. 도커와 작동하는 모든 코드가 쿠버네티스 안에 내장되어 있었음.
- 현재: rkt, cri-o와 같은 컨테이너 런타임도 지원을 열고 확장해야됨 -> CRI 필요
- **CRI**: 쿠버네티스가 컨테이너 런타임과 어떻게 통신할지 정의하는 표준으로, 따라서 어떤 컨테이너 런타임 인터페이스가 개발되더라도 cri 표준을 따르면 됨.


## Container Storage Interface (CSI)
- **CSI**: 다중 저장소를 지원하기 위해 개발됨.
- AWS EBS, portworx, disk, dell emc 등


## Container Network Interface (CNI)
- **CNI**: 다양한 네트워킹 솔루션을 지원하기 위해 개발됨.
- weaveworks, flannel, cilium 등

## Persistent Volumes (PV)
> 클러스터 관리자가 미리 만들어 둔 스토리지 리소스. 클러스터 내에 독립적으로 존재하며, 파드가 없어져도 삭제되지 않는다. 
- NFS, EBS, GCE Persistent Disk 등 외부 스토리지와 연결 가능
- 사용자가 많은 환경에서 많은 pod를 배포할 경우, 사용자는 각각의 pod에 따라 매번 스토리지를 구성해야 한다. 변경사항이 있을 때는 모든 파드에서 수정 사항을 작성하므로, 저장소를 중앙에서 관리할 필요가 있다.
