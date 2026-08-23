# Operations

## SLO and capacity

单节点目标：健康检查 p99 < 50ms，单批 1k 样本写入 p99 < 500ms，查询在 100k 点以内 p99 < 2s。默认每批最多 100,000 个样本，查询结果限制 100,000 点。实际容量由磁盘和标签基数决定，`/api/v1/metrics` 暴露当前用量。

## Failure drills

1. 向 `data/wal.log` 追加截断尾部，重启后 replay 会截断到最后完整帧。
2. 修改 block JSON 校验字段，启动时将其放入 bad block quarantine，不阻断其他数据。
3. 在副本 simulator 中调用 `SetHealthy` 模拟 quorum 丢失，写入返回超时。
4. 保留 worker 先写 tombstone，经过 delete delay 后再删除，审计系统可记录 tombstone 事件。

## Security and recovery

HTTP 读取体限制为 8MiB，所有请求有 15 秒 context deadline；生产部署应在网关启用 TLS、租户认证和网络策略。WAL 与 blocks 应放在持久卷并定期快照，恢复顺序为 WAL replay -> index rebuild -> block health scan。
