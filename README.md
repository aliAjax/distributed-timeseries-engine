# Distributed Time Series Engine

纯 Go 的分布式时序存储与窗口查询引擎参考实现，面向设备、运行和应用测量值。服务采用标准库 HTTP，数据目录包含 `wal.log` 和可校验的 immutable JSON blocks；复制、分片迁移、保留和配额组件以单节点 simulator 形式提供。

## 快速开始

```bash
make test
make run
curl http://localhost:8080/healthz
curl -X POST http://localhost:8080/api/v1/ingest -H 'content-type: application/json' \
  -d '{"samples":[{"tenant":"demo","metric":"cpu_usage","labels":[{"Name":"host","Value":"a"}],"timestamp":"2026-08-21T00:00:00Z","value":42,"quality":1}]}'
curl --get 'http://localhost:8080/api/v1/query' --data-urlencode 'query=cpu_usage{host="a"}[1h]'
```

环境变量：`TS_HTTP_ADDR`、`TS_DATA_DIR`、`TS_MAX_SAMPLES`。关闭服务会安全停止 HTTP；WAL 在下次启动时回放并截断坏尾部。

## 架构与运行手册

写入路径：HTTP -> domain validation -> quota -> WAL fsync -> in-memory series -> block seal。查询路径：parser -> label index -> planner -> bounded repository scan -> aggregation。详细 SLO、容量估算、威胁模型与故障演练见 [docs/operations.md](docs/operations.md)。

## 质量门禁

`make fmt vet test race build smoke` 覆盖格式化、静态检查、单测、竞态检查、编译和接口烟测。项目不依赖外部数据库或消息系统，适合本地和容器演示。
