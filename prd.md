# AIRegistry PRD (Agent & Tool Registry Service)

## 1. Overview
AIRegistry is a central control plane service for registering, managing, and resolving AI tools, agents, and MCP services.

It acts as the source of truth for:
- What tools exist  
- Where they live  
- How they are accessed  
- What policies/security apply  

It is **not in the request path**. It is used for discovery.

---

## 2. Goals
- Centralize tool/agent/service discovery  
- Enable clients to dynamically resolve services  
- Attach security + policy context at registration time  
- Keep it lightweight (no heavy registry frameworks)  
- Provide UI-driven + API-driven management  

---

## 5. Core Concepts

### 5.1 Entity Types
- Tool (e.g., Neo4j MCP tool)  
- MCP Service (external endpoint)  

---

### 5.2 Registry Entry (Example)
```json
{
  "id": "c1f4b6a2-9d2e-4a7f-bc91-3a9f2d8e1e77",
  "name": "neo4j-mcp",
  "type": "mcp_service",
  "endpoint": "http://neo4j-mcp:8080",
  "protocol": "http",
  "auth_required": true,
  "policy_ids": ["policy-time-window-1"],
  "metadata": {
    "env": "prod",
    "owner": "data-team",
    "auth": {
        "type": "basic",
        "username": "user1",
        "password": "*****"
    },
  },
  "status": "active",
  "created_at": "2026-04-08T22:30:00Z"
}
```

## 6. Stack flow
Client → AIRegistry (discovery)
Client → MCPGW → (AISEC → AIPolicy) → External MCP Service

Request Flow
- Client queries AIRegistry to discover:
    - service endpoint
    - metadata
    - policies (optional awareness)
- Client sends request to MCPGW
- MCPGW:
    - Auth validation (AISEC)
    - Policy enforcement (AIPolicy)
- Gateway forwards request to external MCP service