# 6. Cluster Maintenance

## Practice Test - **OS Upgrades**

```yaml
k get pods -o wide

kubectl drain node01 --ignore-daemonsets

k uncordon node01

kubectl cordon node01

```

## Practice Test - **Cluster Upgrade Process**