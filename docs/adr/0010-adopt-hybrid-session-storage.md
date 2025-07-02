# Adopt Hybrid Session Storage

## Status

Accepted

## Context

Our Redis-based session storage solution (ADR-0009) has served us well, but as ShopFlow has grown globally, we've encountered new requirements that necessitate a more sophisticated approach to session management.

**Current Challenges with Redis-Only Approach:**
* **Geographic Distribution**: Users across multiple regions experience latency accessing centralized Redis
* **Compliance Requirements**: New regulatory requirements mandate durable storage of session data for audit trails
* **Cost Optimization**: Redis memory costs are high for storing inactive sessions that are rarely accessed
* **Disaster Recovery**: Need for cross-region session replication for business continuity
* **Session Analytics**: Limited ability to perform analytics on historical session data

**New Business Requirements:**
* Support for users in North America, Europe, and Asia-Pacific regions
* Compliance with GDPR, CCPA, and SOX auditing requirements
* Cost optimization for sessions with different access patterns
* 99.99% availability target including disaster recovery scenarios
* Session analytics for business intelligence and fraud detection

**Technical Requirements:**
* Sub-10ms latency for active sessions
* Cross-region replication with eventual consistency
* Durable storage for compliance and analytics
* Cost-effective storage for inactive sessions
* Seamless failover between regions

## Decision

We will implement a hybrid session storage architecture that combines Redis for hot session data with DynamoDB for durable and cold session storage.

### Hybrid Architecture Design

```mermaid
C4Component
    title Component Diagram - Hybrid Session Storage Architecture
    
    Container_Boundary(client_layer, "Client Applications") {
        Component(webApp, "Web App", "React", "Global customer interface")
        Component(mobileApp, "Mobile App", "React Native", "Mobile applications")
        Component(adminApp, "Admin Portal", "React", "Administrative tools")
    }
    
    Container_Boundary(api_layer, "API Layer") {
        Component(apiGateway, "API Gateway", "Kong", "Global load balancing")
        Component(sessionManager, "Session Manager", "Node.js", "Session orchestration")
        Component(authService, "Auth Service", "Node.js", "Authentication logic")
    }
    
    Container_Boundary(hot_storage, "Hot Storage (Redis)") {
        Component(redisUS, "Redis US-East", "Redis Cluster", "Active US sessions")
        Component(redisEU, "Redis EU-West", "Redis Cluster", "Active EU sessions")
        Component(redisAP, "Redis AP-Southeast", "Redis Cluster", "Active APAC sessions")
    }
    
    Container_Boundary(cold_storage, "Cold Storage (DynamoDB)") {
        Component(dynamoGlobal, "DynamoDB Global", "DynamoDB", "Durable session storage")
        Component(dynamoAnalytics, "Analytics Table", "DynamoDB", "Session analytics data")
    }
    
    Container_Boundary(sync_layer, "Synchronization") {
        Component(sessionSync, "Session Sync", "Lambda", "Hot/Cold sync service")
        Component(replication, "Cross-Region Sync", "DynamoDB Streams", "Region replication")
    }
    
    Rel(webApp, apiGateway, "HTTPS requests")
    Rel(mobileApp, apiGateway, "HTTPS requests")
    Rel(adminApp, apiGateway, "HTTPS requests")
    
    Rel(apiGateway, sessionManager, "Session operations")
    Rel(sessionManager, authService, "Authentication")
    
    Rel(sessionManager, redisUS, "Hot session R/W (US)")
    Rel(sessionManager, redisEU, "Hot session R/W (EU)")
    Rel(sessionManager, redisAP, "Hot session R/W (APAC)")
    
    Rel(sessionManager, dynamoGlobal, "Cold session R/W")
    Rel(sessionSync, dynamoGlobal, "Sync operations")
    Rel(sessionSync, redisUS, "Promote/demote sessions")
    Rel(sessionSync, redisEU, "Promote/demote sessions")
    Rel(sessionSync, redisAP, "Promote/demote sessions")
    
    Rel(dynamoGlobal, replication, "Cross-region replication")
    Rel(dynamoGlobal, dynamoAnalytics, "Analytics pipeline")
```

### Session Lifecycle Management

**Hot Session (Redis):**
* Active sessions accessed within last 30 minutes
* Stored in regional Redis clusters for optimal latency
* Automatic expiration after 2 hours of inactivity
* Cross-region backup for disaster recovery

**Cold Session (DynamoDB):**
* Sessions older than 30 minutes of inactivity
* Durable storage with cross-region replication
* Used for session resurrection and compliance
* Analytics and reporting capabilities

**Session Promotion/Demotion:**
* **Promotion**: Cold → Hot when session becomes active
* **Demotion**: Hot → Cold after inactivity threshold
* **Automatic**: Background service manages transitions
* **Seamless**: Invisible to application logic

### Session State Transition Flow

```mermaid
stateDiagram-v2
    [*] --> Creating
    Creating --> Hot : Session Created
    Hot --> Active : User Activity
    Active --> Hot : Continue Activity
    Active --> Inactive : No Activity (5 min)
    Inactive --> Active : User Returns
    Inactive --> Cold : Inactivity Timeout (30 min)
    Cold --> Hot : Session Access
    Cold --> Expired : TTL Expired
    Hot --> Expired : Max TTL Reached
    Expired --> [*] : Cleanup

    note right of Hot
        Redis Storage
        Sub-10ms latency
        2 hour max TTL
    end note

    note right of Cold
        DynamoDB Storage
        Durable persistence
        30 day retention
    end note
```

### Dynamic Session Access Flow

```mermaid
sequenceDiagram
    participant User
    participant API as API Gateway
    participant SM as Session Manager
    participant Redis
    participant DynamoDB
    participant Sync as Sync Service

    User->>API: Request with session token
    API->>SM: Validate session
    
    alt Session in Redis (Hot)
        SM->>Redis: Get session data
        Redis-->>SM: Session found
        SM-->>API: Valid session
    else Session not in Redis
        SM->>DynamoDB: Check cold storage
        alt Session in DynamoDB
            DynamoDB-->>SM: Session found
            SM->>Redis: Promote to hot storage
            SM-->>API: Valid session
            Note over SM: Session promoted to hot
        else Session not found
            SM-->>API: Invalid session
        end
    end
    
    API-->>User: Response
    
    Note over Sync: Background process
    Sync->>Redis: Check inactive sessions
    Sync->>DynamoDB: Demote to cold storage
    Sync->>Redis: Remove from hot storage
```

### Implementation Strategy

**Phase 1: DynamoDB Foundation (4 weeks)**
* Set up DynamoDB tables with global secondary indexes
* Implement session sync service using Lambda
* Create cross-region replication with DynamoDB Global Tables
* Build monitoring and alerting for hybrid system

**Phase 2: Hot/Cold Logic (3 weeks)**
* Implement session manager with hot/cold awareness
* Build promotion/demotion algorithms
* Add regional Redis cluster support
* Implement failover mechanisms

**Phase 3: Migration (2 weeks)**
* Gradual migration from Redis-only to hybrid model
* A/B testing with percentage of traffic
* Performance validation and optimization
* Full production rollout

**Phase 4: Optimization (3 weeks)**
* Fine-tune promotion/demotion thresholds
* Implement advanced caching strategies
* Add session analytics and reporting
* Performance optimization and cost analysis

### Configuration Parameters

```yaml
session_config:
  hot_storage:
    ttl: 7200  # 2 hours
    promotion_threshold: 1800  # 30 minutes
    regions:
      - us-east-1
      - eu-west-1
      - ap-southeast-1
  
  cold_storage:
    ttl: 2592000  # 30 days
    compliance_retention: 7776000  # 90 days
    analytics_retention: 31536000  # 1 year
  
  sync:
    batch_size: 100
    sync_interval: 300  # 5 minutes
    cross_region_delay: 30  # seconds
```

## Consequences

**Positive:**
* **Performance**: Sub-5ms latency for active sessions via regional Redis
* **Cost Optimization**: 60% reduction in Redis costs by offloading inactive sessions
* **Compliance**: Durable DynamoDB storage meets audit and regulatory requirements
* **Global Scale**: Regional Redis clusters provide optimal performance worldwide
* **Disaster Recovery**: Cross-region replication ensures business continuity
* **Analytics**: Rich session data available for business intelligence
* **Flexibility**: Can adjust hot/cold thresholds based on usage patterns

**Negative:**
* **Complexity**: More sophisticated architecture requiring careful orchestration
* **Consistency**: Eventual consistency between hot and cold storage layers
* **Operational Overhead**: Multiple storage systems require specialized monitoring
* **Development Complexity**: Session manager logic more complex than single-store approach
* **Debugging**: Distributed session state can be challenging to troubleshoot

**Neutral:**
* **Migration Effort**: Significant one-time effort to implement hybrid system
* **Team Training**: Engineering team needs expertise in both Redis and DynamoDB
* **Monitoring**: Enhanced observability required for multi-tier storage
* **Cost Model**: Different cost structure with usage-based DynamoDB pricing

### Success Metrics

**Performance Targets:**
* < 5ms latency for hot session access (95th percentile)
* < 50ms latency for cold session promotion (95th percentile)
* 99.99% availability across all regions
* < 1% session loss during region failover

**Cost Targets:**
* 50%+ reduction in total session storage costs
* Optimal cost per active user across regions
* Predictable scaling costs with usage growth

**Compliance Targets:**
* 100% session audit trail retention
* < 1 second session data retrieval for compliance queries
* Cross-region data sovereignty compliance

### Monitoring and Alerting

* **Hot/Cold Ratio**: Track percentage of sessions in each tier
* **Promotion/Demotion Rates**: Monitor session lifecycle transitions
* **Cross-Region Latency**: Track replication performance
* **Cost Optimization**: Monitor cost per session across storage tiers
* **Availability Metrics**: Track uptime across all regions and storage layers

---

*This ADR represents our current approach to session storage, combining the performance benefits of Redis with the durability and cost advantages of DynamoDB in a globally distributed architecture.*