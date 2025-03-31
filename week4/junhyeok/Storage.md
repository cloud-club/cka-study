# Storage in Docker

### File system

- /var/lib/docker
    - aufs
    - containers
    - image
    - volumes

### Layered architecture

- Dockerfile

```yaml
FROM Ubuntu

RUN apt-get update && apt-get -y install python

RUN pip install flask flask-mysql

COPY . /opt/source-code

ENTRYPOINT FLASK_APP=/opt/source-code/app.py flask run
```

docker build Docerfile -t mmumshad/my-custom-app

- Layer1. Base Ubuntu Layer - 120MB
- Layer2. Changes in apt packages - 306MB
- Layer3. Changes in pip packages - 6.3MB
- Layer4. Source code - 229B
- Layer5. Update Entrypoint - 0B

- Dockerfile2

```yaml
FROM Ubuntu

RUN apt-get update && apt-get -y install python

RUN pip install flask flask-mysql

COPY . /opt/source-code

ENTRYPOINT FLASK_APP=/opt/source-code/app.py flask run
```

docker build Docerfile -t mmumshad/my-custom-app

- Layer1. Base Ubuntu Layer - 0MB
- Layer2. Changes in apt packages - 0MB
- Layer3. Changes in pip packages - 0MB
- Layer4. Source code - 229B
- Layer5. Update Entrypoint - 0B

= Docker는 캐시에서 재사용해 최신 소스 코드를 업데이트 해 빨리 응용 프로그램 이미지를 재구축 한다.

= 재건축과 업데이트 작업 시 시간이 많이 절약됌

= Image Layers - Read Only

docker run mmumshad/my-custom-app

- Layer6. Container Layer - Read, Write

### COPY-ON-WRITE

Image Layer의 파일을 Copy 해서 Container Layer에서 수정 가능

- 컨테이너 삭제 시 Container Layer에 저장된 모든 데이터 삭제됌

### volumes

데이터를 유지하기 위해 Docker Host에 volume 생성(볼륨 마운트)

- docker volume create data_volume
- /var/lib/docker
    - volumes
        - data_volume

컨테이너 내부에 마운트

- docker run -v data_volume:/var/lib/mysql mysql
- sql 저장 기본 위치

= 컨테이너가 파괴되어도 데이터는 살아 있다.

Docker Host 외부 저장소에 데이터 저장(바인트 마운트)

위치 : /data/mysql

- docker run -v /data/mysql:/var/lib/mysql mysql

-v 옵션은 옛날 방식

- docker run —mount type=bind, source=/data/mysql, target=/var/lib/mysql mysql

### Storage drivers

- AUFS - Ubuntu의 기본
- ZFS
- BTRFS
- Device Mapper
- Overlay
- Overlay2

Docker는 운영체제에 근거해 자동으로 사용 가능한 최고의 저장소 드라이버를 선택함

### Volume Driver Plugins in Docker

- Local / Azure File Storage / Convoy / DigitalOcean Block Storage / RexRay …

AWS EBS

```yaml
docker run -it \
	--name mysql
	--volume-driver rexray/ebs
	--mount src=ebs-vo1, target=/var/lib/mysql mysql
```

---

# Container Storage Interface(CSI)

= 다중 저장소 솔루션을 지원하도록 개발됌

= 쿠버네티스 특성이 아닌 유니버셜 표준이어야함.

- portworx / Amazon EBS / DELL EMC / GlusterFS

Pod가 생성되고 Volume이 필요할 때 쿠버네티스가 CreateVolume RPC를 호출해 Volume 이름 같은 세부사항을 전달해야 한다.

- CreateVolume
- DeleteVolume
- ControllerPublishVolume

---

# Volumes

- Pod에서 생성된 데이터가 Volume에 저장됌
- 삭제된 후에도 데이터는 남아있음

### Volumes & Mounts

컨테이너 내부의 디렉터리에 volume 마운트

![image.png](attachment:885497c2-6e89-4f8f-b8dc-f3c4aaa2804a:image.png)

```yaml
apiVersion: v1
kind: Pod
metadata:
	name: random-number-generator
spec:
	containers:
	- image: alpine
		name: alpine
		command: ["/bin/sh","-c"]
		args: ["shuf -i 0-100 -n 1 >> /opt/number.out;"]
		volumeMounts:
		- mountPath: /opt
			name: data-volume

	volumes:
	- name: data-volume
		hostPath:
			path: /data
			type: Directory
```

### Volume Types

= 쿠버네티스는 다양한 유형의 저장소 솔루션을 지원

![image.png](attachment:80be4706-0e96-4e5c-b972-47939721e006:image.png)

Amazon EBS volume구성

```yaml
	volumes:
	- name: data-volume
		awsElasticBlockStore:
			volumeID: <volume-id>
			fsType: ext4
```

![image.png](attachment:af32aee9-4270-440a-a93a-4fd8802b6368:image.png)

---

# Persistent Volumes

관리자가 거대한 저장소 풀을 생성, 사용자가 요구에 따라 나눠 사용

![image.png](attachment:73ccd2ef-4761-4ec5-8568-212d5cea0955:image.png)

```yaml
pv-definition.yaml

apiVersion: v1
kind: PersistentVolume
metadata:
	name: pv-vol1
	
spec: 
	accessModes:
		- ReadWriteOnce
	capacity:
		storage: 1Gi
		
	hostPath:
		path: /tmp/data
		
		
kubectl create -f pv-definition.yaml
kubectl get persistentvolume

- 저장소 솔루션 치환(hostPath)
	awsElasticBlockStore:
		volumeID: <volume-id>
		fsType: ext4
```

- accessModes : 호스트에 볼륨이 어떻게 마운트 되어야 하는지 정의
- capacity storage : 예약될 저장소 양 명시

---

# Persistent Volumes Claims

관리자는 Persistent Volumes 세트를 만들고, 사용자는 저장소를 사용하기 위해 Persistent Volumes Claims를 만든다.

- Persistent Volumes Claims이 생성되면 쿠버네티스가 Claims에 Persistent Volumes을 묶는다.

![image.png](attachment:04d2b16d-6c79-4e61-8207-7085de14d2a3:image.png)

![image.png](attachment:a7051e37-1a4e-4d01-83b1-c2973bc68910:image.png)

```yaml
pvc-definition.yaml

apiVersion: v1
kind: PersistentVolumeClaim
metadata:
	name: myclaim
spec:
	accessModes:
		- ReadWriteOnce
	resources:
		requests:
			storage: 500Mi
			
kubectl create -f pvc-definition.yaml
kybectl get persistentvolumeclaim
```

- accessModes 일치
- 요청한 용량 500Mi는 스토리지 볼륨으로 설정 가능

![image.png](attachment:afbe942c-6ee1-4c2d-bd3d-32680f26e904:image.png)

- View PVCs
    - kubectl get persistentvolumeclaim
- Delete PVCs
    - kubectl delete persistentvolumeclaim myclaim
- 옵션
    - persistentVolumeReclaimPolicy: Retain
        - 기본값: 유지
        - 관리자가 수동으로 삭제할 때 까지 지속적으로 볼륨이 남음
    - persistentVolumeReclaimPolicy: Delete
        - Claim이 삭제되는 순간 Volume도 삭제
    - persistentVolumeReclaimPolicy: Recycle
        - 데이터 Volume이 다른 Claim에 사용되기 전에 삭제

### Using PVCs in Pods

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
    - name: myfrontend
      image: nginx
      volumeMounts:
      - mountPath: "/var/www/html"
        name: mypd
  volumes:
    - name: mypd
      persistentVolumeClaim:
        claimName: myclaim
```

## Practice Test - Persistent Volumes and Persistent Volume Claims

```yaml

```

---

# Storage Class

Static Provisioning

- Google 클라우드에서 수동으로 디스크 프로비전
- 수동으로 PV 정의 파일 생성

Dynamic Provisioning

- Google 클라우드에서 저장소를 자동으로 프로비저닝

```yaml
PV
pv-definition.yaml

apiVersion: v1
kind: PersistentVolume
metadata:
	name: pv-vol1
spec: 
	accessModes:
		- ReadWriteOnce
	capacity:
		storage: 500Mi
		
gcePersistentDisk:
	pdName: pd-disk
	fsType: ext4
```

```yaml
SC
sc-definition.yaml

apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
	**name: google-storage**
provisioner: kubernetes.io/gce-pd
```

```yaml
PVC
pvc-definition.yaml

apiVersion: v1
kind: PersistentVolumeClaim
metadata:
	name: myclaim
spec:
	accessModes:
		- ReadWriteOnce
	**storageClassName: google-storage**
	resources:
		requests:
			storage: 500Mi
```

```yaml
Pod
pod-definition.yaml

apiVersion: v1
kind: Pod
metadata:
	name: random-number-generator
spec:
	containers:
	- image: alpine
		name: alpine
		command: ["/bin/sh","-c"]
		args: ["shuf -i 0-100 -n 1 >> /opt/number.out;"]
		volumeMounts:
		- mountPath: /opt
			name: data-volume
	volumes:
	- name: data-volume
		persistentVolumeClaim:
			claimName: myclaim
```

- PVC가 생성되면 연관된 저장소 클래스는 정의 프로비저너를 이용해 GCP에 요구되는 사이즈의 새 디스크를 프로비저닝 한다. 이후 영구볼륨을 생성해 볼륨에 PVC를 묶는다.
- 수동으로 PV를 생성할 필요가 없음