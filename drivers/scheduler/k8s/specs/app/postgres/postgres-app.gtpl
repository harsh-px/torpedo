apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: "{{.Name}}"
spec:
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  replicas: "{{.Replicas}}"
  template:
    metadata:
      labels:
        app: "{{.Name}}"
    spec:
      containers:
      - name: postgres
        image: postgres:9.5
        imagePullPolicy: "IfNotPresent"
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_USER
          value: pgbench
        - name: POSTGRES_PASSWORD
          value: superpostgres
        - name: PGBENCH_PASSWORD
          value: superpostgres
        - name: PGDATA
          value: /var/lib/postgresql/data/pgdata
        volumeMounts:
        - mountPath: /var/lib/postgresql/data
          name: postgredb
      volumes:
      - name: postgredb
        persistentVolumeClaim:
          claimName: "{{.Name}}"
