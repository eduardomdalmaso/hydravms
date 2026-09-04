---
name: validate-project
description: Procedure for running full project validation, DDD layer compliance audits, modular frontend file size checks (<100 lines), Go unit tests, build checks, and documentation synchronization on HydraVMS.
---

# 🛡️ Project Validation & Self-Healing Skill (HydraVMS)

Use this skill to audit and validate all codebase changes in HydraVMS.

---

## 1. Automated Checklist

Execute the following commands sequentially before completing any prompt:

```bash
# 1. Frontend Line Limit (< 100 lines per file)
wc -l web/src/**/*.vue web/src/**/*.ts web/src/**/*.css web/src/**/*.js 2>/dev/null | awk '$1 > 100 { print "VIOLATION: " $2 " has " $1 " lines (>100)" }'

# 2. DDD Domain Purity (Zero net/http or database/sql in internal/domain)
grep -rnE "(net/http|database/sql|github.com/gin-gonic|github.com/labstack)" internal/domain/ && echo "VIOLAÇÃO DDD: Remova infraestrutura do domínio!"

# 3. Unit & Integration Tests
go test -v ./...

# 4. Binary Compilation
go build -o /dev/null ./cmd/hydravms
```
