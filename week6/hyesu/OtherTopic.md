# JSON Path

[https://kodekloud.com/p/json-path-quiz](https://kodekloud.com/p/json-path-quiz)

[https://mmumshad.github.io/json-path-quiz/index.html#!/?questions=questionskub1](https://mmumshad.github.io/json-path-quiz/index.html#!/?questions=questionskub1)

[https://mmumshad.github.io/json-path-quiz/index.html#!/?questions=questionskub2](https://mmumshad.github.io/json-path-quiz/index.html#!/?questions=questionskub2)

https://kubernetes.io/ko/docs/reference/kubectl/jsonpath/

데이터 필터링 작업

### Step

1. Identify the kubectl command
2. Familiarize with **JSON** output
    
    ```bash
    kubectl get nodes -o json
    kubectl get pods -o json
    ```
    
3. Form the **JSON Path** query
    1. ex. `.items[0].spec.containers[0].image`
4. Use the **JSON Path** query with **kubectl** command
    
    ```bash
    kubectl get pods -o=jsonpath='{.items[0].spec.containers[0].image}'
    ```
    

### Example

```bash
kubectl get nodes -o=jsonpath='{.items[*].metadata.name}'

kubectl get nodes -o=jsonpath='{.items[*].metadata.name}{"\n"}{.items[*].status.capacity.cpu}'
# master node01
# 4      4
```

### Loops - Range

```bash
kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.capacity.cpu}{"\n"}{end}'
# master 4
# node01 4
```

### Custom Columns

```bash
kubectl get nodes -o=custom-columes=<COLUMN NAME>:<JSON Path>
```

ex. 

```bash
kubectl get nodes -o=custom-columes=NODE:.metadata.name, CPU: .status.capacity.cpu
# NODE   CPU
# master  4 
# node01  4
```

### Sort

```bash
kubectl get nodes --sort-by=.metadata.name
```
