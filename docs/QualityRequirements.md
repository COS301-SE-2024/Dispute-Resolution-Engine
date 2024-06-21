# Quality Requirements

## 1. Reliability

The Dispute Resolution Engine must maintain high uptime due to the nature of the services we provide. Clients should have continuous access to the system to ensure they can respond promptly to time-sensitive communications. A reliable system with high uptime not only benefits our clients but also increases our throughput.

| Stimulus Source | Stimulus | Response | Response Measure | Environment | Artifact |
|-----------------|----------|----------|------------------|-------------|----------|
| System Users/Clients | Attempt to access the Dispute Resolution Engine | The system should maintain high uptime to ensure continuous access | 99.9% system uptime, fewer than 1 hour of downtime per year | Production environment | Dispute Resolution Engine |

## 2. Scalability and Performance

Scalability and performance go hand-in-hand in this system. Disputes are common occurrences and it is expected that the Dispute Resolution Engine (DRE) must be capable of handling many cases simultaneously, as well as multiple active users, per case, at any given time. Performance is essential for scalability. A system that performs well for each user can scale more effectively.

### 2.1 Quantification

Renowned for its high aptitude to function concurrently, Golang will be the reason for our highly performant backend. As for the actual server, we will be using a Terreco server to host Dispute Resolution Engine. The server boasts powerful hardware, further benefitting performance and throughput.

## 3. Usability

The Dispute Resolution Engine is designed to be user-friendly. It aims to enable easy navigation through simple and intuitive user interfaces. This ensures that users can operate the system effectively without extensive training or technical knowledge. This is crucial for a platform that is intended to accommodate users from a wide range of domains.

### 3.1 Quantification

Simple but small extra features to guide users when trying to navigate an application. To aid our clients as much as possible we plan to implement tooltips, a guided tutorial on the site and a well written user-manual of course.

## 4. Security

The Dispute Resolution Engine will handle sensitive information regularly. In addition to user credentials, participants in disputes will expect a high level of confidentiality when exchanging information. This is mainly due to documented evidence or communications that can contain exploitable or otherwise sensitive information. Therefore, a fundamental requirement of the Dispute Resolution Engine is to ensure the protection of both the personal and dispute-related information of our clients.

### 4.1 Quantification

Security is quantified by ensuring that all the OWASP checks are passed, each of which identify vunerabilities within the system. Covering these will ensure that all holes are filled regarding the security of our system. For further clarification, please review the [OWASP document](https://github.com/COS301-SE-2024/Dispute-Resolution-Engine/blob/feat/documentation/docs/OWASP.md) on our github.

## 5. Maintainability

The Dispute Resolution Engine must be easy to update and extend. It is essential that expansion of the system allows for the addition of new legal processes, dispute types, and other features. This will ensure that the system remains relevant and up-to-date with the latest legal requirements and dispute resolution practices.

### 5.1 Quantification

DRE promotes maintainability as it is important for adding functionality and new features to the system in the future. We do this by ensuring that the frontend makes use of a component library, the backend's API is fully extensible (it is trivial to add/remove endpoints). 

