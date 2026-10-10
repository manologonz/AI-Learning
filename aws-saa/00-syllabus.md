# AWS Solutions Architect – Associate (SAA-C03) — Course Syllabus

**Learner starting level:** Beginner (little/no cloud experience)
**Goal:** Pass AWS Certified Solutions Architect – Associate, then build toward Professional (SAP-C02)
**Hands-on:** Has an AWS account — every unit includes a lab
**Estimated total time:** ~10-12 weeks at ~8-10 hrs/week (≈ 95 hours)

## Exam map (SAA-C03 domain weights)

| Domain | Weight | Main units |
|---|---|---|
| 1. Design Secure Architectures | 30% | 2, 3, 6, 9, 14 |
| 2. Design Resilient Architectures | 26% | 5, 7, 8, 10, 11, 12 |
| 3. Design High-Performing Architectures | 24% | 4, 5, 6, 7, 8, 9, 10 |
| 4. Design Cost-Optimized Architectures | 20% | 1, 4, 5, 6, 13, 15 |

## Units

### Part A — Foundations
1. **Cloud & AWS Foundations, Account Safety** (~5h) — Global infrastructure, shared responsibility, pricing models, Well-Architected pillars, billing/budgets.
   - Objectives: explain Region/AZ/edge; apply the shared responsibility model to a scenario; set budgets and alarms; name the six Well-Architected pillars and map a trade-off to them
   - Prerequisites: none
2. **IAM & Access Management** (~8h) — Users, groups, roles, policies, STS, MFA, permission boundaries.
   - Objectives: read and write JSON policies; explain policy evaluation logic (explicit deny > allow > implicit deny); choose role vs user; design cross-account access
   - Prerequisites: 1
3. **Networking I: VPC** (~9h) — CIDR, subnets, route tables, IGW, NAT, security groups vs NACLs, endpoints.
   - Objectives: subnet a CIDR block; build a public/private VPC; distinguish SG vs NACL (stateful vs stateless); choose Gateway vs Interface endpoints
   - Prerequisites: 1

### Part B — Core services
4. **Compute: EC2** (~8h) — Instance families, AMIs, EBS/instance store, purchasing options, placement groups, user data.
   - Objectives: pick an instance family for a workload; choose On-Demand/Reserved/Savings Plans/Spot; select an EBS volume type; explain placement group strategies
   - Prerequisites: 2, 3
5. **Storage: S3, EBS, EFS, FSx** (~9h) — Storage classes, lifecycle, versioning, replication, encryption, file systems.
   - Objectives: choose storage type by access pattern; design lifecycle policies; secure a bucket; pick EFS vs FSx variants
   - Prerequisites: 2, 4
6. **Databases** (~9h) — RDS, Aurora, DynamoDB, ElastiCache, Redshift, and choosing among them.
   - Objectives: match workload to database; explain Multi-AZ vs read replicas; use DynamoDB capacity modes and keys; know when to cache
   - Prerequisites: 3, 4

### Part C — Scale, resilience, decoupling
7. **Load Balancing, Auto Scaling & DNS (Route 53)** (~8h) — ALB/NLB/GWLB, ASG policies, routing policies, health checks.
   - Objectives: choose the load balancer type; design scaling policies; pick a Route 53 routing policy for a scenario
   - Prerequisites: 3, 4
8. **Decoupling & Integration** (~7h) — SQS, SNS, EventBridge, Kinesis, Step Functions, API Gateway.
   - Objectives: decide queue vs topic vs stream; design fan-out; explain at-least-once vs FIFO; choose the right integration service
   - Prerequisites: 4, 7
9. **Serverless & Containers** (~7h) — Lambda, Fargate, ECS, EKS, App Runner, Cognito.
   - Objectives: know Lambda limits and concurrency; choose ECS vs EKS vs Fargate vs Lambda; design a serverless web backend
   - Prerequisites: 2, 8
10. **Content Delivery & Global Performance** (~5h) — CloudFront, Global Accelerator, edge caching, S3 Transfer Acceleration.
    - Objectives: choose CloudFront vs Global Accelerator; design cache behaviours; secure origins (OAC)
    - Prerequisites: 5, 7
11. **High Availability, Disaster Recovery & Backup** (~7h) — RTO/RPO, DR strategies, AWS Backup, multi-Region patterns.
    - Objectives: pick among backup/restore, pilot light, warm standby, multi-site; compute RTO/RPO trade-offs; design multi-AZ/multi-Region architectures
    - Prerequisites: 5, 6, 7

### Part D — Security, operations, cost
12. **Advanced Networking & Hybrid Connectivity** (~7h) — Peering, Transit Gateway, VPN, Direct Connect, PrivateLink.
    - Objectives: choose the connectivity pattern; design hub-and-spoke; avoid CIDR overlap traps
    - Prerequisites: 3
13. **Security & Encryption in Depth** (~8h) — KMS, Secrets Manager, ACM, WAF, Shield, GuardDuty, Inspector, Macie, Organizations/SCPs.
    - Objectives: design encryption at rest/in transit; pick the right security service per threat; apply SCPs
    - Prerequisites: 2, 5, 12
14. **Monitoring, Governance & Operations** (~5h) — CloudWatch, CloudTrail, Config, Systems Manager, CloudFormation.
    - Objectives: distinguish CloudWatch/CloudTrail/Config; write basic CloudFormation; use Systems Manager
    - Prerequisites: 4
15. **Data, Analytics & Migration** (~6h) — Athena, Glue, EMR, Redshift, Lake Formation, DMS, DataSync, Snow Family, Transfer Family.
    - Objectives: choose the migration tool by data size/network; design a basic data lake; pick an analytics service
    - Prerequisites: 5, 6
16. **Cost Optimization & Well-Architected Review** (~4h) — Cost Explorer, pricing traps, right-sizing, data transfer costs.
    - Objectives: reduce cost of a given architecture; identify hidden data-transfer charges; apply the pillars in trade-off questions
    - Prerequisites: all prior

### Part E — Exam readiness
17. **Exam Strategy & Scenario Drills** (~6h) — Decoding question wording, eliminating distractors, timing, 65 questions / 130 min.
    - Objectives: identify constraint keywords (least effort, most cost-effective, highly available); eliminate distractors systematically
    - Prerequisites: all prior
- **Final Exam** — 65-question cumulative exam weighted by domain, plus a capstone architecture design task.

### Part F — Professional bridge (optional, after passing SAA)
18. **Multi-Account Strategy & Governance** (Organizations, Control Tower, RAM, SCPs at scale)
19. **Complex Migration & Modernization** (7 Rs, Migration Hub, Application Discovery, hybrid)
20. **Advanced Design & Optimization** (multi-Region active-active, performance tuning, advanced cost)

## How this course works
Each unit: **lecture → hands-on lab → quiz → exam → resources**. I teach the lecture conversationally, you complete the lab in your account, then take the quiz (formative) and exam (harder, cumulative in-unit). I grade with explanations, log results in `progress.md`, and re-teach weak objectives before moving on. Before the final exam, we do spaced-repetition reviews of earlier units. Labs are designed to stay in free tier; Unit 1 sets up budgets first so you can't be surprised by a bill.
