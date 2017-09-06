kind: PersistentVolumeClaim
apiVersion: v1
metadata:
   name: "{{.Name}}"
   annotations:
     volume.beta.kubernetes.io/storage-class: "{{.StorageClass}}"
spec:
   accessModes:
     - ReadWriteOnce
   resources:
     requests:
       storage: "{{.Size}}Gi"