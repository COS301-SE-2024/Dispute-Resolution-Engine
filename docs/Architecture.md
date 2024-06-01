# Architectural Patterns

## Quality Requirements Addressed: [Link](QualityRequirements.md)

## Architectural Constraints

### Adherence to Legal Standards and Regulations:

#### POPIA Compliance

-   The system architecture must incorporate robust data encryption mechanisms and access control measures to ensure compliance with the Protection of Personal Information Act (POPIA). This includes safeguarding personal information during data handling, storage, and processing.

#### Regular Legal Reviews

-   The system architecture should be modular and scalable to facilitate regular reviews and updates without causing disruptions to the entire system. This allows for seamless integration of changes based on legal requirements and ensures continuous compliance with relevant laws.

#### Data Protection Officer

-   Implement role-based access control within the system architecture to ensure that the Data Protection Officer (DPO) has appropriate access privileges to oversee data protection strategies and compliance activities.

### Balancing Automation and Human Oversight:

#### Human-in-the-loop System

-   The system architecture must support workflow orchestration to seamlessly integrate human oversight into critical decision-making processes. This involves designing workflows that enable efficient collaboration between automated systems and human operators.

### Limiting Bias in the Dispute Resolution Process:

#### Algorithmic Transparency

-   Architectural design should prioritize the use of explainable AI models, which provide clear explanations for their decisions, thereby reducing bias and facilitating audits. This involves structuring the system architecture to accommodate algorithms that offer transparency in decision-making processes.

## Event-Driven Architectural Pattern

### Quality Requirements Addressed:

1. Scalability

    - Event-driven architecture allows components to react to events asynchronously, enabling the system to handle large volumes of events efficiently and scale horizontally by adding more event processors.

2. Performance

    - By separating the responsibilities of the production and consumption of events, the system can process events in real-time and handle high-throughput scenarios efficiently.

3. Reliability

    - Events are typically stored in a persistent event log, which can be replayed in case of system failures, ensuring data integrity and system reliability.

## Service-Oriented Architecture (SOA) Pattern

### Quality Requirements Addressed:

1. Scalability

    - Services in an SOA can be deployed and scaled independently, allowing the system to grow and manage increased loads effectively.

2. Maintainability

    - By encapsulating business logic within discrete services, SOA makes it easier to update, replace, or extend functionalities without affecting other parts of the system.

## Gatekeeper Architectural Pattern

### Quality Requirements Addressed:

1. Security

    - The Gatekeeper pattern acts as a security layer that enforces access control policies, ensuring that only authorized requests are allowed to reach the internal services.

2. Reliability

    - By routing requests through a single entry point, the Gatekeeper pattern can provide load balancing, caching, and failover mechanisms to improve the reliability of the system.

3. Compliance

    - The Gatekeeper pattern can enforce compliance requirements by validating requests against predefined policies and ensuring that data privacy and security regulations are met.

## Flux

### Quality Requirements Addressed:

1. Usability

    - Flux architecture simplifies the data flow in the system, making it easier to understand and maintain. This enhances the usability of the system by providing a clear and predictable data flow.

## Conclusion of Choices

### Event-Driven Architecture

1. Scalability and Performance
2. Reliability

### Service-Oriented Architecture (SOA)

1. Scalability
2. Maintainability

### Gatekeeper Pattern

1. Security
2. Reliability
3. Compliance
