# MGR api migration plan

[Dashboard OpenAPI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/ceph/ceph/main/src/pybind/mgr/dashboard/openapi.yaml)

| API group   | status   | notes   | source |
|---|---|---|
| Auth  | ✅  | Call mon command   | |
| Cluster | ✅ | Call mon command  |   | |
| User  | ✅  | Call mon command  | |
| ClusterConfiguration  | ☐ | TODO  | |
