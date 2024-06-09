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

## Architectural Styles
In developing the Dispute Resolution Engine (DRE), we have selected several architectural patterns to address key quality requirements.

The Event-Driven Architecture (EDA) enables DRE to handle large volumes of events efficiently by allowing components to react to events asynchronously. This pattern significantly enhances the system's scalability by facilitating horizontal scaling through the addition of more event processors. Additionally, EDA boosts performance by decoupling the production and consumption of events, allowing for real-time processing and efficient handling of high-throughput scenarios. Reliability is also improved, as events are typically stored in a persistent event log that can be replayed in case of system failures, ensuring data integrity and consistent operation.

Service-Oriented Architecture (SOA) is employed to enhance DRE's scalability and maintainability. By encapsulating business logic within discrete, independently deployable services, SOA allows each service to be scaled independently, effectively managing growth of the system. This modular approach also simplifies maintenance as services can be updated, replaced, or extended without impacting other parts of the system. This ensures continuous improvement and adaptability to changing requirements and technologies.

The Gatekeeper pattern serves as a security layer within DRE, addressing essential quality requirements of security, reliability, and compliance. Acting as a single entry point for all incoming requests, the Gatekeeper enforces strict access control policies to ensure only authorized requests reach internal services. This pattern enhances system reliability by incorporating load balancing, caching, and failover mechanisms, thus maintaining smooth operation even under high demand or in the event of partial system failures. Moreover, the Gatekeeper helps ensure compliance with data privacy and security regulations by validating requests against predefined policies, thereby safeguarding sensitive information and ensuring regulatory adherence.

We make use of the Flux pattern to enhance the usability of the DRE by simplifying data flow and state management. Its unidirectional data flow ensures predictable and traceable state changes. By structuring the frontend into modular components, Flux makes the system easier to maintain and debug. This approach improves the user experience, offering developers an intuitive framework and providing end-users with a responsive and reliable interface.

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

### Flux Pattern

1. Usability

## Architectural Diagram

![image](https://github.com/COS301-SE-2024/Dispute-Resolution-Engine/assets/64808970/c87ddbb1-2ab8-4f7c-bce6-53a0ae029d49)


## Class Diagram

![classDiagram](https://github.com/COS301-SE-2024/Dispute-Resolution-Engine/assets/64808970/bf3169c4-df33-4c2a-aa61-f74257ea3d42)