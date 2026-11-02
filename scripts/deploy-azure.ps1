param(
    [Parameter(Mandatory)][ValidatePattern('^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,127}$')][string]$ImageTag,
    [switch]$Build,
    [switch]$TemporaryData
)
$ErrorActionPreference = 'Stop'
if (!$TemporaryData) { throw 'utnapp-env is Express and does not support Azure Files. Explicit -TemporaryData is required for this demo.' }
Write-Warning 'DEMO: SQLite is ephemeral. Data is lost on restart, scale-to-zero or deployment.'
$credentials = $env:AZURE_CREDENTIALS | ConvertFrom-Json
if (!$credentials.clientSecret -or !$env:ACR_PASSWORD) { throw 'AZURE_CREDENTIALS and ACR_PASSWORD are required' }
$env:ACR_PASSWORD = $env:ACR_PASSWORD.Trim()
$token = Invoke-RestMethod -Method Post -Uri "https://login.microsoftonline.com/$($credentials.tenantId)/oauth2/v2.0/token" -Body @{
    grant_type='client_credentials'; client_id=$credentials.clientId
    client_secret=$credentials.clientSecret; scope='https://management.azure.com/.default'
}
$headers = @{Authorization="Bearer $($token.access_token)"}
$root = "https://management.azure.com/subscriptions/$($credentials.subscriptionId)/resourceGroups/frsrutn"
function Request([string]$Method, [string]$Path, $Body = $null) {
    $args = @{Method=$Method; Uri="$root/$Path"; Headers=$headers}
    if ($null -ne $Body) { $args.ContentType='application/json'; $args.Body=($Body | ConvertTo-Json -Depth 30 -Compress) }
    Invoke-RestMethod @args
}
function Wait-Resource([string]$Path) {
    for ($i=0; $i -lt 90; $i++) {
        $resource = Request GET $Path
        $state = $resource.properties.provisioningState
        if ($state -eq 'Succeeded') { return $resource }
        if ($state -in @('Failed','Canceled')) { throw "Provisioning failed: $Path" }
        Start-Sleep -Seconds 5
    }
    throw "Provisioning timed out: $Path"
}
$image = "ingcalidad2026.azurecr.io/grupo10-carva:$ImageTag"
if ($Build) {
    $env:ACR_PASSWORD | docker login ingcalidad2026.azurecr.io --username ingcalidad2026 --password-stdin
    if ($LASTEXITCODE) { throw 'Registry login failed' }
    try {
        docker build --platform linux/amd64 -t $image .
        if ($LASTEXITCODE) { throw 'Image build failed' }
        docker push $image
        if ($LASTEXITCODE) { throw 'Image push failed' }
    } finally { docker logout ingcalidad2026.azurecr.io | Out-Null }
}
$environmentPath = 'providers/Microsoft.App/managedEnvironments/utnapp-env?api-version=2024-03-01'
$environment = Request GET $environmentPath
$location = $environment.location
$appPath = 'providers/Microsoft.App/containerApps/grupo10-carva?api-version=2024-03-01'
$definition = @{
    location=$location; tags=@{project='Carva';group='10'}
    properties=@{
        managedEnvironmentId=$environment.id
        configuration=@{
            ingress=@{external=$true;targetPort=8080;transport='auto';allowInsecure=$false;traffic=@(@{latestRevision=$true;weight=100})}
            secrets=@(@{name='acr-password';value=$env:ACR_PASSWORD})
            registries=@(@{server='ingcalidad2026.azurecr.io';username='ingcalidad2026';passwordSecretRef='acr-password'})
        }
        template=@{
            containers=@(@{
                name='carva';image=$image;resources=@{cpu=0.25;memory='0.5Gi'}
                env=@(@{name='CARVA_DB_PATH';value='/data/carva.db'})
                probes=@(@{type='Readiness';httpGet=@{path='/api/projects';port=8080};initialDelaySeconds=5;periodSeconds=10})
            })
            scale=@{minReplicas=0;maxReplicas=1;rules=@(@{name='http';http=@{metadata=@{concurrentRequests='10'}}})}
        }
    }
}
Request PUT $appPath $definition | Out-Null
$app = Wait-Resource $appPath
$url = "https://$($app.properties.configuration.ingress.fqdn)"
$healthy=$false
for ($i=0; $i -lt 60; $i++) {
    try { $response = Invoke-WebRequest "$url/api/projects" -TimeoutSec 30; if ($response.StatusCode -eq 200) { $healthy=$true; break } } catch { }
    Start-Sleep -Seconds 5
}
if (!$healthy) { throw "Deployment created but HTTP check failed: $url" }
Write-Output "Deployed $image to $url"
