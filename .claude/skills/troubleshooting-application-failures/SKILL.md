---
name: troubleshooting-application-failures
description: Troubleshoots failing applications by discovering and analyzing CloudWatch log groups to identify error patterns, root causes, and actionable solutions. Use when an application is experiencing failures and log-based diagnosis is needed.
version: 1
---

# Application Failure Troubleshooting

## Overview

Domain expertise for diagnosing application failures through CloudWatch log analysis.
Discovers relevant log groups, searches for error patterns and stack traces, performs
root cause analysis, and generates prioritized remediation recommendations.

## Project context - read before acting

Run `aws` through the Bash tool with `--profile probability --region us-east-1`.
There is no `call_aws` tool here. This procedure is READ-ONLY.

Only two log groups exist and they hold **error lines only** (the shipper filters
`ERROR|WARN|FATAL|PANIC|panic:|Exception|error|warn|fatal` before upload,
retention 7 days):

- `/probability/back-central`
- `/probability/front-central`

Skip the discovery step, use those directly. Consequences of the filter:

- A failure that logs no matching keyword is INVISIBLE in CloudWatch. Absence of
  errors is not evidence the service is healthy.
- Request context around an error (the preceding INFO lines) is NOT there. For
  full logs use `docker logs` on the EC2 via SSM - see `.claude/rules/infra-ops.md`.
- A spike in ingest is itself the signal: the `logs-ingesta-alta` alarm exists
  because a RabbitMQ requeue loop once wrote 4.2 MB/min of pure ERROR lines.
  Before deep-diving a flood, check for a hot requeue loop
  (`.claude/rules/colas-errores-permanentes.md`).

Also check `.claude/bitacora/` for a matching past incident before diagnosing.

## Troubleshoot a failing application

To diagnose and resolve application failures using CloudWatch logs, follow the
procedure exactly. See [Application failure troubleshooting procedure](references/application-failure-troubleshooting.md).

## Troubleshooting

### No log groups found

Ask the user for specific log group names. Common patterns: `/aws/lambda/function-name`,
`/aws/apigateway/api-name`, or custom application log groups.

### Access denied errors

Verify AWS credentials have `logs:DescribeLogGroups`, `logs:DescribeLogStreams`,
`logs:StartQuery`, and `logs:GetQueryResults` permissions.

### Query timeouts

Reduce the time window or limit results. Large log groups may require multiple smaller queries.
