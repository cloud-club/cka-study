# Kustomize

환경별로 변경이 필요한 부분만 수정하고 Kubernetes 구성을 재사용할 수 있는 방법

⇒ 확장 가능성!

- base
    - 모든 환경에서 동일하게 적용될 구성
- overlay
    - 환경별로 동작 지정
    - base를 덮어쓸 속성 지정

Base + Overlay ⇒ Final Manifests

Base와 Overlay 합치는게 Kustomize 역할

- **Folder Structure**
    
    ![image.png](13%20Kustomize%20Basics%201d0766b7e9a48001b05be5bdfcde3f92/image.png)
    

- **장점**
    - kubectl에 기본으로 제공하므로 다른 패키지 설치할 필요 X
    - helm처럼 템플릿 시스템을 사용하지 않고, 템플릿 언어 배울 필요 없어서 편리 
    → 표준 YAML 사용

### Kustomize vs Helm

- **Helm**
    
    ```
    .
    ├── templates/
    │   ├── ingress.yaml
    │   ├── deploy.yaml
    │   └── service.yaml
    └── enviroments/
        ├── values.dev.yaml
        └── values.prod.yaml
    ```
    
    - Go 템플릿 구문을 활용하여 다양한 속성에 변수 할당
    - 더 많은 기능 제공 ex. 조건부, 루프, 함수, hook,  …

- **Kustomize**
    - 매우 간단

# Install & Setup

[https://kubectl.docs.kubernetes.io/installation/kustomize/binaries/](https://kubectl.docs.kubernetes.io/installation/kustomize/binaries/)

```bash
kustomize version --short
```

# Kustomization YAML File

- `kustomization.yaml`

```yaml
# apiVersion과 kind는 선택사항
# but, 향후 변경 사항 대비해서 하드코딩 추천
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

# kustomize에 의해 관리될 리소스
resources:
	- nginx-deploy.yaml
	- nginx-service.yaml

# Customization
commonLabels:
	company: kodeKloud
```

```bash
kustomize build .
# 출력 결과물은 최종 구성이 어떨지 보여주기만 함
# 적용 -> 이 명령의 출력을 가져와서 kubectl apply
```

### Apply Kustomize Configs

```bash
kustomize build . | kubectl apply -f -

kustomize build . | kubectl delete -f -
```

```bash
kubectl apply -k .

kubectl delete -k .
```

# Managing Directories

```
.
├── kustomization.yaml
├── api/
│   ├── api-deploy.yaml
│   ├── api-service.yaml
│   └── kustomization.yaml
└── db/
    ├── db-deploy.yaml
    ├── db-service.yaml
    └── kustomization.yaml
```

- `api/kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
	- api-deploy.yaml
	- api-service.yaml
```

- `kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
	- api/
	- db/
```

```bash
kustomize build . | k apply -f -
# or
kubectl apply -k .
```

# Transformers

구성을 수정하거나 변환

- **Common Transformation**
    - commonLabel
    - namePrefix/Suffix
    - Namespace
    - commonAnnotations

ex. `kustomization.yaml`

```yaml
commonLabels:
org: KodeKloud

namespace: lab

namePrefix: KodeKloud-

nameSuffix: -dev

commonAnnotations:
	brank: master
```

### Image Transformer

- `web-deploy.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
	name: web-deployment
spec:
	replicas: 1
	selector:
		matchLabels:
			component: web
	template:
		metadata:
			labels:
				component: web
		spec:
			containers:
				- name: web
					image: nginx
```

- `kustomization.yaml`

```yaml
images:
	- name: nginx # target image
		newName: haproxy
		
	- name: nginx
		newTag: "2.4"
		
	- name: nginx
		newName: haproxy
		newTag: "2.4"
```

# Patches

구성 수정하는 또다른 방법

- 제공해야하는 parameter
    - Operation: add/remove/replace, …
    - Target
    - Value

metadata.name이 api-deployment로 되어있는데,

web-deployment로 변경하고 싶은 경우

- `api-deploy.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
	name: api-deployment
spec:
	replicas: 1
	selector:
		matchLabels:
			component: api
	template:
		metadata:
			labels:
				component: api
		spec:
			containers:
				- name: nginx
					image: nginx
```

- `kustomization.yaml`

```yaml
patches:
	- target:
			kind: Deployment
			name: api-deployment
			
		patch: |- # 인라인 패치
			- op: replace
				path: /metadata/name # 업데이트하려는 속성이 무엇인지
				value: web-deployment
```

복제본 값 바꾸고 싶은 경우

- `kustomization.yaml`

```yaml
patches:
	- target:
			kind: Deployment
			name: api-deployment
			
		patch: |- # 인라인 패치
			- op: replace
				path: /spec/replicas # 업데이트하려는 속성이 무엇인지
				value: 5
```

### 정의 방법

- **JSON 6902**
    - 위의 예시에서 보여준 경우
    
    ```yaml
    patches:
    	- target:
    			kind: Deployment
    			name: api-deployment
    			
    		patch: |- # 인라인 패치
    			- op: replace
    				path: /spec/replicas # 업데이트하려는 속성이 무엇인지
    				value: 5
    ```
    
    - target, patch 정보 제공해야 함

- **Strategic Merge Patch**
    
    ```yaml
    patches:
    	- patch: |-
    			apiVersion: apps/v1
    			kind: Deployment
    			metadata: # 변경할 target 정보 알려줌
    				name: api-deployment
    			spec:
    				replicas: 5
    ```
    
    - 원본 배포 파일 복붙 후, 변경하고 싶지 않은 부분 모두 삭제

### Patch Type

- **Inline**
    - 지금까지 사용해온 방식

- **Seperate File**
    - yaml 파일 만들고 대신 제공
    - JSON 6902
        - `kustomization.yaml`
        
        ```yaml
        patches:
        	- path: replica-patch.yaml
        		target:
        			kind: Deployment
        			name: nginx-deployment
        ```
        
        - `replica-patch.yaml`
        
        ```yaml
        - op: replace
        	path: /spec/replicas
        	value: 5
        ```
        
    - Strategic Merge Patch
        - `kustomization.yaml`
        
        ```yaml
        patches:
        	- replica-patch.yaml
        ```
        
        - `replica-patch.yaml`
        
        ```yaml
        apiVersion: apps/v1
        kind: Deployment
        metadata:
        	name: api-deployment
        spec:
        	replicas: 5
        ```
        
    

### Patches Dictionar

- **Replace**
- **Add**
    - JSON6902
    
    ```yaml
    patches:
    	- target:
    			kind: Deployment
    			name: api-deployment
    			
    		patch: |- 
    			- op: add
    				path: /spec/template/metadata/labels/org
    				value: KodeKloud
    ```
    
    - Strategic Merge Patch
        - `kustomization.yaml`
            
            ```yaml
            patches:
            	- label-patch.yaml
            ```
            
        - `label-patch.yaml`
            
            ```yaml
            apiVersion: apps/v1
            kind: Deployment
            metadata:
            	name: api-deployment
            spec:
            	template:
            		metadata:
            			labels:
            				org: KodeKloud # 추가하고 싶은 레이블 추가
            ```
            
- **Remove**
    - JSON6902
    
    ```yaml
    patches:
    	- target:
    			kind: Deployment
    			name: api-deployment
    			
    		patch: |- 
    			- op: remove
    				path: /spec/template/metadata/labels/org
    ```
    
    - Strategic Merge Patch
        - `kustomization.yaml`
            
            ```yaml
            patches:
            	- label-patch.yaml
            ```
            
        - `label-patch.yaml`
            
            ```yaml
            apiVersion: apps/v1
            kind: Deployment
            metadata:
            	name: api-deployment
            spec:
            	template:
            		metadata:
            			labels:
            				org: null
            ```
            

### Patches List

- `api-deploy.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
	name: api-deployment
spec:
	replicas: 1
	selector:
		matchLabels:
			component: api
	template:
		metadata:
			labels:
				component: api
		spec:
			containers:
				- name: nginx
					image: nginx
```

- **Replace**
    - JSON6902
    
    ```yaml
    patches:
    	- target:
    			kind: Deployment
    			name: api-deployment
    			
    		patch: |- 
    			- op: replace
    				path: /spec/template/spec/containers/0 # 인덱스
    				value:
    					name: haproxy
    					image: haproxy
    ```
    
    - Strategic Merge Patch
        - `kustomization.yaml`
            
            ```yaml
            patches:
            	- label-patch.yaml
            ```
            
        - `label-patch.yaml`
            
            ```yaml
            apiVersion: apps/v1
            kind: Deployment
            metadata:
            	name: api-deployment
            spec:
            	template:
            		spec:
            			containers:
            				- name: nginx
            					image: haproxy
            ```
            

- **Add**
    - JSON6902
    
    ```yaml
    patches:
    	- target:
    			kind: Deployment
    			name: api-deployment
    			
    		patch: |- 
    			- op: add
    				path: /spec/template/spec/containers/- # 경로 끝에 위치 지정
    				value:
    					name: haproxy
    					image: haproxy
    ```
    

- **Remove**
    - Strategic Merge Patch
        - `kustomization.yaml`
            
            ```yaml
            patches:
            	- label-patch.yaml
            ```
            
        - `label-patch.yaml`
            
            ```yaml
            apiVersion: apps/v1
            kind: Deployment
            metadata:
            	name: api-deployment
            spec:
            	template:
            		spec:
            			containers:
            				- $patch: delete # 삭제 지시문
            					name: database
            ```
            

# Overlays

```
~/someApp
├── base 
│   ├── nginx-deploy.yaml
│   ├── kustomization.yaml
│   └── service.yaml
└── overlays
    ├── dev
    │   ├── config-map.yaml
    │   └── kustomization.yaml
    └── prod
        ├── config-map.yaml
        └── kustomization.yaml
```

- base
    - Share or default configs across all environments
- overlays
    - Environment specific configurations that add or modify base configs

- `base/kustomization.yaml`

```yaml
resources:
	- nginx-deploy.yaml
	- service.yaml
```

- `base/nginx-deploy.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
	name: nginx-deployment
spec:
	replicas: 1
```

- `overlays/dev/kustomization.yaml`

```yaml
bases:
	- ../../base

patch: |-
	- op: replace
		path: /spec/replicas
		value: 2
```

⇒ base 구성 가져오고 patch 진행!

overlays에는 base에 없는 새로운 구성 있을 수 O

⇒ 원하는 만큼 새로운 리소스 추가 가능

- `overlays/prod/kustomization.yaml`

```yaml
bases:
	- ../../base
	
resources:
	- grafana-deploy.yaml
	
patch: |-
	- op: replace
		path: /spec/replicas
		value: 2
```

# Components

여러 Overlay에 포함될 수 있는 재사용 가능한 구성 로직 정의할 수 있는 제공

```
.
├── base 
│   ├── api-deploy.yaml
│   └── kustomization.yaml
├── components/
│   ├── caching/
│   │   ├── deploy-patch.yaml
│   │   ├── redis-deploy.yaml
│   │   └── kustomization.yaml
│   └── db/
│       ├── deploy-patch.yaml
│       ├── postgres-deploy.yaml
│       └── kustomization.yaml
└── overlays
    ├── dev/
    │   └── kustomization.yaml
    ├── premium/
    │   └── kustomization.yaml
    └── standalone/
        └── kustomization.yaml
```

- `components/db/kustomization.yaml`

```yaml
apiVersion: kustomization.config.k8s.io/v1alpha1
kind: Component

resources:
	- postgres-deploy.yaml
	
secretGenerator:
	- name: postgres-cred
		literals:
			- password=postgres123

patches:
	- deploy-patch.yaml
```

- `components/db/deploy-patch.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
	name: api-deployment
spec:
	template:
		spec:
			containers:
				- name: api
					env: 
						- name: DB_PASSWORD
							valueFrom:
								secretKeyRef:
									name: postgres-cred
									key: password
```

- `overlays/dev/kustomization.yaml`

```yaml
bases:
	- ../../base

components:
	- ../../components/db
```
