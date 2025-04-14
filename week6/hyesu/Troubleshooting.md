# Application Failure

애플리케이션이 어떻게 구성되어 있는지 미리 적어놓는게 좋음!

위에서 부터 ~ 

Application → Dependent

- **Check Accessibility**
    
    ```bash
    curl http://web-service-ip:node-port
    ```
    
- **Check Service Status**
    - selector과 pod 일치하는지 확인
- **Check Pod**
    - pod 상태와 재시작 횟수 확인
    
    ```bash
    kubectl get pod
    
    kubectl describe pod web
    kubectl logs web -f --previous
    ```
    
- **Check Dependent Service**
- **Check Dependent Applications**

# Control Plane Failure

- **Check Node Status**
    
    ```bash
    kubectl get nodes
    
    kubectl get pods
    ```
    
- **Check Controlplane Pods**
    
    ```bash
    kubectl get pods -n kube-system
    ```
    
- **Check Controlplane Services**
    
    ```bash
    service kube-apiserver status
    
    service kube-controller-manager status
    
    service kube-scheduler status
    
    service kubelet status
    
    service kube-proxy status
    ```
    
- **Check Service Logs**
    
    ```bash
    kubectl logs kube-apiserver-master -n kube-system
    
    sudo journalctl -u kube-apiserver
    ```
    

# Worker Node Failure

- **Check Node Status**
    
    ```bash
    kubectl get nodes
    
    kubectl describe node worker-1
    # Condition -> Type과 Status 확인!
    ```
    
- **Check Node**
    - if. 노드가 master와 끊긴 경우, CPU와 메모리, 디스크 공간 확인
    
    ```bash
    top
    
    df -h
    ```
    
    - kubelet 상태 확인
    
    ```bash
    service kubelet status
    
    sudo journalctl -u kubelet
    ```
    
- **Check Certificates**
    
    ```bash
    openssl x509 -in /var/lib/kubelet/worker-1.crt -text
    ```
    

# Network

### Network Plugin

1. **Weave Net**
    - 설치:
        
        ```bash
        bash
        복사편집
        kubectl apply -f https://github.com/weaveworks/weave/releases/download/v2.8.1/weave-daemonset-k8s.yaml
        
        ```
        
    - 공식 문서: [https://kubernetes.io/docs/concepts/cluster-administration/addons/#networking-and-network-policy](https://kubernetes.io/docs/concepts/cluster-administration/addons/#networking-and-network-policy)
2. **Flannel**
    - 설치:
        
        ```bash
        bash
        복사편집
        kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml
        
        ```
        
    - 참고: Flannel은 NetworkPolicy를 지원하지 않음.
3. **Calico**
    - 설치:
        
        ```bash
        bash
        복사편집
        curl -O https://raw.githubusercontent.com/projectcalico/calico/v3.25.0/manifests/calico.yaml
        kubectl apply -f calico.yaml
        
        ```
        
    - 가장 강력한 CNI 플러그인 중 하나로 평가됨.

> 시험 팁 (CKA/CKAD): CNI 설치를 직접 요구하지 않음. 요구 시 정확한 URL 제공됨.
> 

---

### CoreDNS

- **기본 DNS 서버**: CoreDNS
- **영향 요소**: Pod/Service 수, 캐시 크기, 쿼리량
- **리소스**:
    - `Deployment`: coredns
    - `Service`: kube-dns
    - `ConfigMap`: coredns (Corefile 구성)
- **포트**: DNS는 기본적으로 **53번 포트** 사용

- **주요 설정 예시 (Corefile)**

```
kubernetes cluster.local in-addr.arpa ip6.arpa {
   pods insecure
   fallthrough in-addr.arpa ip6.arpa
   ttl 30
}
forward . /etc/resolv.conf
```

- **DNS 문제 해결 팁**
    1. CoreDNS Pod가 Pending 상태 → **CNI 설치 확인**
    2. CrashLoopBackOff 상태 → SELinux/Docker 이슈 해결:
        - Docker 업그레이드
        - SELinux 비활성화
        - `allowPrivilegeEscalation: true` 설정 변경
        - `resolv.conf` 루프 해결 (예: `/run/systemd/resolve/resolv.conf` 또는 `8.8.8.8` 지정)
    3. `kube-dns` 서비스에 endpoint 없는 경우:
        
        ```bash
        kubectl -n kube-system get ep kube-dns
        ```
        

### kube-proxy

- **역할**: 노드에서 네트워크 프록시 역할 수행, Virtual IP → 실제 Pod로 트래픽 전달
- **위치**: DaemonSet으로 실행 (`kube-system` 네임스페이스)
- **명령어 예시**:
    
    ```bash
    /usr/local/bin/kube-proxy --config=/var/lib/kube-proxy/config.conf --hostname-override=$(NODE_NAME)
    ```
    

- **주요 설정 항목**
    - `clusterCIDR`
    - `mode` (iptables or ipvs)
    - `bindAddress`
    - `kube-config`

- **문제 해결**
1. kube-proxy Pod 상태 확인
2. 로그 확인
3. ConfigMap 및 config.conf 확인
4. 포트 확인
    
    ```bash
    netstat -plan | grep kube-proxy
    ```
    

### 📚 참고 문서

- 서비스 디버깅: [https://kubernetes.io/docs/tasks/debug-application-cluster/debug-service/](https://kubernetes.io/docs/tasks/debug-application-cluster/debug-service/)
- DNS 문제 해결: [https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/](https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/)
