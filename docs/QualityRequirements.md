# Quality Requirements

## 1. Reliability

The Dispute Resolution Engine must maintain high uptime due to the nature of the services we provide. Clients should have continuous access to the system to ensure they can respond promptly to time-sensitive communications. A reliable system with high uptime not only benefits our clients but also increases our throughput.

| Stimulus Source | Stimulus | Response | Response Measure | Environment | Artifact |
|-----------------|----------|----------|------------------|-------------|----------|
| System Users/Clients | Attempt to access the Dispute Resolution Engine | The system should maintain high uptime to ensure continuous access | 99.9% system uptime, fewer than 1 hour of downtime per year | Production environment | Dispute Resolution Engine |

## 2. Scalability and Performance

Scalability and performance go hand-in-hand in this system. Disputes are common occurrences and it is expected that the Dispute Resolution Engine (DRE) must be capable of handling many cases simultaneously, as well as multiple active users, per case, at any given time. Performance is essential for scalability. A system that performs well for each user can scale more effectively.

| Stimulus Source     | Stimulus                                             | Response                                                                                          | Response Measure                        | Environment         | Artifact                  |
|---------------------|------------------------------------------------------|---------------------------------------------------------------------------------------------------|-----------------------------------------|---------------------|---------------------------|
| System Users/Clients | High volume of disputes and multiple active users per case | The system should handle multiple cases simultaneously and perform efficiently for each user       | System supports X concurrent users and Y disputes without performance degradation | Production environment | Dispute Resolution Engine |


## 3. Usability

The Dispute Resolution Engine is designed to be user-friendly. It aims to enable easy navigation through simple and intuitive user interfaces. This ensures that users can operate the system effectively without extensive training or technical knowledge. This is crucial for a platform that is intended to accommodate users from a wide range of domains.

| Stimulus Source     | Stimulus                                        | Response                                                                                         | Response Measure                                           | Environment         | Artifact                  |
|---------------------|-------------------------------------------------|--------------------------------------------------------------------------------------------------|------------------------------------------------------------|---------------------|---------------------------|
| System Users/Clients | Attempt to navigate and use the system           | The system should provide a simple and intuitive user interface requiring minimal training        | User satisfaction scores above X, average task completion time within Y minutes | Production environment | Dispute Resolution Engine |


## 4. Security

The Dispute Resolution Engine will handle sensitive information regularly. In addition to user credentials, participants in disputes will expect a high level of confidentiality when exchanging information. This is mainly due to documented evidence or communications that can contain exploitable or otherwise sensitive information. Therefore, a fundamental requirement of the Dispute Resolution Engine is to ensure the protection of both the personal and dispute-related information of our clients.

| Stimulus Source             | Stimulus                                      | Response                                                                                         | Response Measure                                                                                                 | Environment         | Artifact                  |
|-----------------------------|-----------------------------------------------|--------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------|---------------------|---------------------------|
| Malicious actors/Unauthorised users | Attempt to access sensitive information       | The system should ensure the protection of personal and dispute-related information               | Unauthorized access attempts identified and flagged, all OWASP checks passed                                      | Production environment | Dispute Resolution Engine |


## 5. Maintainability

The Dispute Resolution Engine must be easy to update and extend. It is essential that expansion of the system allows for the addition of new legal processes, dispute types, and other features. This will ensure that the system remains relevant and up-to-date with the latest legal requirements and dispute resolution practices.

| Stimulus Source     | Stimulus                                      | Response                                                                                         | Response Measure                                                                                              | Environment         | Artifact                  |
|---------------------|-----------------------------------------------|--------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------|---------------------|---------------------------|
| Development Team    | Requirement to update or extend the system    | The system should allow for easy addition of new legal processes, dispute types, and other features | Frontend uses a component library, backend API is fully extensible with trivial addition/removal of endpoints | Development environment | Dispute Resolution Engine |


