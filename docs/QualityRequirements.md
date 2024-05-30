# Quality Requirements

## 1. Reliability

The Dispute Resolution Engine must maintain high uptime due to the nature of the services we provide. Clients should have continuous access to the system to ensure they can respond promptly to time-sensitive communications. A reliable system with high uptime not only benefits our clients but also increases our throughput.

### 1.1 Quanitifications

We will be using a VPS server to host Dispute Resolution Engine via domains.co.za. Their servers boast an outstanding 99.9% uptime as well as free monthly backups, ensuring that our server is always up and can recover from data corruption.

## 2. Scalability and Performance

Scalability and performance go hand-in-hand in this system. Disputes are common occurrences and it is expected that the Dispute Resolution Engine (DRE) must be capable of handling many cases simultaneously, as well as multiple active users, per case, at any given time. Performance is essential for scalability. A system that performs well for each user can scale more effectively.

### 2.1 Quantification

Renowned for its high aptitude to function concurrently, Golang will be the reason for our highly perfromante backend. As for the actual server, we will be using [stubbed] to host Dispute Resolution Engine. The server boasts powerful hardware, further benefitting performance and throughput.

## 3. Usability

The Dispute Resolution Engine is designed to be user-friendly. It aims to enable easy navigation through simple and intuitive user interfaces. This ensures that users can operate the system effectively without extensive training or technical knowledge. This is crucial for a platform that is intended to accommodate users from a wide range of domains.

### 3.1 Quantification

Simple but small extra features to guide users when trying to navigate an application. To aid our clients as much as possible we plan to implement tooltips, a guided tutorial on the site and a well written user-manual of course.

## 4. Security

The Dispute Resolution Engine will handle sensitive information regularly. In addition to user credentials, participants in disputes will expect a high level of confidentiality when exchanging information. This is mainly due to documented evidence or communications that can contain exploitable or otherwise sensitive information. Therefore, a fundamental requirement of the Dispute Resolution Engine is to ensure the protection of both the personal and dispute-related information of our clients.

## 5. Compliance

<!--!need to clarify with Neil what should be covered here -->

Disputes involve numerous legal processes, so it is vital to ensure that DRE meets legal, regulatory, and contractual obligations relevant to the handling of disputes.

## 6. Maintainability

The Dispute Resolution Engine must be easy to update and extend. It is essential that expansion of the system allows for the addition of new legal processes, dispute types, and other features. This will ensure that the system remains relevant and up-to-date with the latest legal requirements and dispute resolution practices.
