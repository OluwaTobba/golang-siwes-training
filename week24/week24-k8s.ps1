if ((kubectl config current-context) -ne "docker-desktop") { Write-Host "Switch to the docker-desktop context first: kubectl config use-context docker-desktop"; return }

# Free the ports used by the compose stack
docker compose down

# Make sure the image exists (same Dockerfile as compose)
$gv = (Select-String -Path go.mod -Pattern '^go (\d+\.\d+)').Matches[0].Groups[1].Value
docker build --build-arg GO_VERSION=$gv -t inventory-api:v1.0.0 .

New-Item -ItemType Directory -Force k8s | Out-Null

@'
apiVersion: v1
kind: Secret
metadata:
  name: inventory-secret
type: Opaque
stringData:
  db-password: secret
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  labels:
    app: postgres
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:16-alpine
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_USER
              value: postgres
            - name: POSTGRES_DB
              value: siwes_db
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: inventory-secret
                  key: db-password
          readinessProbe:
            exec:
              command: ["pg_isready", "-U", "postgres"]
            initialDelaySeconds: 5
            periodSeconds: 5
          volumeMounts:
            - name: data
              mountPath: /var/lib/postgresql/data
              subPath: pgdata
      volumes:
        - name: data
          emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: postgres-svc
spec:
  selector:
    app: postgres
  ports:
    - port: 5432
      targetPort: 5432
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: inventory-config
data:
  DB_HOST: "postgres-svc"
  DB_PORT: "5432"
  DB_USER: "postgres"
  DB_NAME: "siwes_db"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inventory-api
  labels:
    app: inventory-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: inventory-api
  template:
    metadata:
      labels:
        app: inventory-api
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
      containers:
        - name: inventory-api
          image: inventory-api:v1.0.0
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8080
            - containerPort: 50051
          envFrom:
            - configMapRef:
                name: inventory-config
          env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: inventory-secret
                  key: db-password
          resources:
            requests:
              cpu: "100m"
              memory: "64Mi"
            limits:
              cpu: "500m"
              memory: "128Mi"
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8080
            initialDelaySeconds: 3
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: inventory-svc
spec:
  selector:
    app: inventory-api
  ports:
    - name: http
      port: 8080
      targetPort: 8080
    - name: grpc
      port: 50051
      targetPort: 50051
'@ | Set-Content -Encoding ascii k8s/inventory.yaml

kubectl apply -f k8s/inventory.yaml
kubectl rollout status deployment/postgres --timeout=180s
kubectl rollout status deployment/inventory-api --timeout=180s
kubectl get pods,svc

# Port-forward REST and gRPC, then test
$pf1 = Start-Process kubectl -ArgumentList "port-forward","svc/inventory-svc","8080:8080" -PassThru -WindowStyle Hidden
$pf2 = Start-Process kubectl -ArgumentList "port-forward","svc/inventory-svc","50051:50051" -PassThru -WindowStyle Hidden
Start-Sleep -Seconds 5

$base = "http://localhost:8080"
$sku = "K8S-" + (Get-Date -Format HHmmss)
$new = Invoke-RestMethod -Method Post "$base/api/v1/products" -ContentType "application/json" -Body (@{name="Router"; sku=$sku; price=250; stock=15} | ConvertTo-Json)
"Created in Kubernetes: id=$($new.id) sku=$($new.sku) stock=$($new.stock)"
Invoke-RestMethod "$base/api/v1/products" | Format-Table id, name, sku, price, stock
go run ./cmd/client

# Self-healing: kill an API pod
$victim = kubectl get pods -l app=inventory-api -o jsonpath="{.items[0].metadata.name}"
kubectl delete pod $victim
Start-Sleep -Seconds 10
kubectl get pods -l app=inventory-api

Stop-Process -Id $pf1.Id, $pf2.Id -Force -ErrorAction SilentlyContinue