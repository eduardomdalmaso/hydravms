---
name: vms-reports-analytics
description: High-performance report generation (PDF, Excel, CSV), optimized aggregation queries (Materialized Views, Timescale continuous aggregates), and analytical dashboards for HydraVMS.
---

# 📊 Report Generation Engine & High-Performance Analytics Skill (HydraVMS)

This skill covers the technical guidelines for report generation, asynchronous export jobs, and database query optimizations for high-throughput AI event analytics.

---

## 1. Relatórios Suportados & Formatos de Saída

```text
[ Solicitação do Operador ] (Filtros de Data, Câmeras, Tipos de Eventos)
              │
              ▼ (POST /api/v1/reports)
┌────────────────────────────────────────────────────────────────────────┐
│             HYDRAVMS ASYNC REPORT WORKER (Fila de Tarefas em Go)       │
├────────────────────────────────────────────────────────────────────────┤
│ 📄 PDF Executivo / Forense:                                            │
│    • Resumo gerencial de alarmes com gráficos de tendência e horários. │
│    • Tabela de incidentes com thumbnails de snapshots embutidos.       │
│    • Links auditáveis para download de clipes MP4 de evidência.        │
├────────────────────────────────────────────────────────────────────────┤
│ 📊 Excel (XLSX) & CSV Streaming:                                       │
│    • Exportação de milhares de eventos via streaming (`excelize` em Go)│
│    • Zero estouro de memória RAM (< 20MB de heap durante exportação).  │
├────────────────────────────────────────────────────────────────────────┤
│ 📦 Upload Direto no MinIO S3:                                          │
│    • Salvo em: `hydravms-reports/{tenant_id}/{YYYY}/{MM}/relatorio.pdf`│
│    • Disponibilizado via Presigned URL de download válida por 24h.     │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Otimização de Consultas & Agregações Analíticas

Para responder a dashboards em menos de **10ms**, mesmo com milhões de eventos:

### A. Materialized Views com Atualização Automática (Rollups de 1 Hora e 1 Dia)
```sql
CREATE MATERIALIZED VIEW mv_hourly_event_stats AS
SELECT 
    tenant_id,
    camera_id,
    event_type,
    severity,
    date_trunc('hour', triggered_at) AS event_hour,
    COUNT(*) AS total_events,
    COUNT(CASE WHEN status = 'acknowledged' THEN 1 END) AS acknowledged_count,
    COUNT(CASE WHEN status = 'resolved' THEN 1 END) AS resolved_count,
    AVG(EXTRACT(EPOCH FROM (resolved_at - triggered_at))) AS avg_resolution_time_seconds
FROM events
GROUP BY tenant_id, camera_id, event_type, severity, date_trunc('hour', triggered_at);

CREATE UNIQUE INDEX idx_mv_hourly_stats ON mv_hourly_event_stats (tenant_id, camera_id, event_type, severity, event_hour);
```

### B. Índices Compostos & Índices Parciais de Alta Seletividade
```sql
-- Busca rápida para a tela de monitoramento de eventos pendentes de atendimento
CREATE INDEX idx_events_unresolved ON events (tenant_id, triggered_at DESC) 
WHERE status = 'new';

-- Busca rápida de alvos de alta severidade (intrusões críticas)
CREATE INDEX idx_events_critical ON events (tenant_id, triggered_at DESC) 
WHERE severity = 'critical';
```

---

## 3. Métricas Analíticas Prontas para Dashboard
1. **Pico de Horários de Alarme:** Distribuição de eventos por hora do dia (00h às 23h).
2. **Heatmap de Câmeras com Mais Incidentes:** Ranking das câmeras com maior frequência de disparo.
3. **SLA Operacional (MTTA / MTTR):** Tempo médio até o operador reconhecer o alarme (*Mean Time to Acknowledge*) e tempo médio de resolução (*Mean Time to Resolve*).
4. **Taxa de Falsos Positivos:** Percentual de eventos marcados com status `false_positive` para refinamento do modelo YOLO no HydraForge.
