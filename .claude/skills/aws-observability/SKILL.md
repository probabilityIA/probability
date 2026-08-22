---
name: aws-observability
description: >-
  CloudWatch observability on AWS: Log Insights queries, alarms (metric, composite,
  anomaly detection), dashboards, custom metrics and EMF, metric filters, log
  retention, and CloudTrail operational auditing. Applies to CloudWatch, alarms,
  dashboards, log groups, log queries, EMF, metric math, and CloudTrail lookups.
  Not for application logging design or security threat detection.
metadata:
  version: "2-probability"
---

# AWS Observability

## Overview

Domain expertise for CloudWatch metrics, logs, alarms, dashboards and CloudTrail
auditing. Trimmed for this project: the upstream Application Signals, ADOT,
X-Ray, Synthetics and Dynamic Instrumentation material was removed because none
of it is deployed here (no EKS, no ADOT collector, no instrumented services).

**Note:** reference files contain quota values and feature matrices that change.
When precision matters, confirm against current AWS documentation.

## Project context - read before acting

This repo has its own operating rules that OVERRIDE any generic advice below.
See `.claude/rules/infra-ops.md`. The ones that collide most often:

1. **Every command needs `--profile probability --region us-east-1`.**
2. **Do NOT send CloudTrail to CloudWatch Logs.** The trail `probability-trail`
   is deliberately S3-only; Logs ingestion is billed. Same reason data events
   stay off.
3. **Log ingestion is the cost driver, not storage.** The 5 GB free tier is
   ingest. Container logs are grep-filtered to error lines before shipping
   (~0.24% of volume). Never propose shipping unfiltered logs.
4. **Any new alarm must set `AlarmActions`** to
   `arn:aws:sns:us-east-1:476702565908:probability-alertas`. An alarm without it
   notifies nobody - that exact bug was live until 2026-08-21.
5. **EC2 detailed monitoring stays disabled** (1-minute metrics are billed).
6. Free-tier ceilings already in use: 10 alarms (7 used), 3 dashboards (0 used).

Existing setup: 7 alarms, SNS topic `probability-alertas` -> secamc93@gmail.com,
log groups `/probability/back-central` and `/probability/front-central`
(7-day retention), RDS Enhanced Monitoring in `RDSOSMetrics`.

## Routing

| User need | Action |
|-----------|--------|
| Writing Log Insights queries | Read [log-insights.md](references/log-insights.md) |
| Configuring alarms (metric, composite, anomaly) | Read [alarms.md](references/alarms.md) |
| Publishing custom metrics or using EMF | Read [metrics.md](references/metrics.md) |
| Building dashboards | Read [dashboards.md](references/dashboards.md) |
| Debugging observability issues | Read [troubleshooting.md](references/troubleshooting.md) - starts with the 5 most common fixes |
| CloudTrail operational auditing | Read [cloudtrail.md](references/cloudtrail.md) |
| Spans multiple areas | Read the most specific reference first, then consult others |

## Files

| File | Content |
|------|---------|
| [alarms.md](references/alarms.md) | Metric, composite, anomaly detection alarms - configuration, constraints, recommended defaults |
| [log-insights.md](references/log-insights.md) | Query syntax, commands, functions, known issues, reusable query library |
| [metrics.md](references/metrics.md) | Custom metrics, EMF spec, metric filters, high-resolution, retention |
| [dashboards.md](references/dashboards.md) | Widget types, cross-account/region, dynamic labels, sharing |
| [troubleshooting.md](references/troubleshooting.md) | Error -> cause -> fix for observability services |
| [cloudtrail.md](references/cloudtrail.md) | Operational auditing, event types, S3+Athena queries |
| [alarm-template.ts](assets/alarm-template.ts) | CDK alarm+dashboard reference (this project uses Terraform - read for the alarm thresholds, not the CDK) |
