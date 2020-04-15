{
    "apiVersion": "compute.cnrm.cloud.google.com/v1beta1",
    "kind": "ComputeBackendService",
    "metadata": {
        "name": "web-app-backend-service"
    },
    "spec": {
        "backend": [
            (.network_endpoint_groups as $negs |
             {
                group: { networkEndpointGroupRef: { external:  "https://www.googleapis.com/compute/v1/projects/$[PROJECT_ID]/zones/\(.zones[])/networkEndpointGroups/\($negs[])"}},
                balancingMode: "RATE",
                capacityScaler: 1.0,
                maxRatePerEndpoint: 100
             })
        ],
        "healthChecks": [
            {
                "healthCheckRef": {
                    "name": "web-app-health-check"
                }
            }
        ],
        "location": "global",
        "customRequestHeaders": [
            "Host:$[BACKEND_SERVICE_HOST_NAME]"
        ],
        "loadBalancingScheme": "EXTERNAL",
        "protocol": "HTTP",
        "portName": "http",
        "sessionAffinity": "NONE",
        "timeoutSec": 30,
        "connectionDrainingTimeoutSec": 300
    }
}