# Functional Requirements

## Overview

Provide a brief overview of the system or feature being developed.

## Functional Requirements

### Core Requirements

1. **Two-Party Engagement** - Louis
   A system where disputing parties (Party A and Party B) can register, submit their dispute, and provide evidence or statements.
   - File sharing (Multimedia)
   - Text based forum
   - Group disputes as individual entities
2. **Mediator(s) Assignment** - Louis, Neil
   An algorithm for assigning a neutral mediator based on the dispute type, expertise required, and mediator availability.
   - Classification of dispute type
   - Sourcing of expertise (pulling docs from historically similar cases?)
   - Generate lists of available mediators
   - Only assigned when all evidence submitted
3. **Dispute Resolution Workflow** - Louis, Neil
   A step-by-step process that guides parties through initial submission, mediation sessions, and final resolution. This should include timelines, notification systems, and a checklist of required actions for each party.
   - "Wizard" like process
   - Templates for disputes
4. **Secure Communication Channel** - Louis
   Encrypted messaging and document exchange between the parties and the mediator, ensuring confidentiality.
   - End-to-end encrypted communication channels
5. **Natural Language Processing**: - Neil
   Process the PDF submissions such that the mediators can have a summary of the dispute statement and evidence.
   - PDF files need to be "digested" into text for NLP
   - Possible application of Sentiment Analysis
6. **Analytics profile** - Louis
   Create an analytic profile for the mediators and lawyers to see the outcomes, resolutions, summary and statistics of cases in certain areas.
   - Categorisation and classification of disputes
   - Archiving of disputes.
   - Define metrics of disputes for tracking
7. **Universal Dispute Creator** - Louis
   Allow for the dispute system to work with any field/domain.
   - Further use for templatisation/generic dispute entities
   - Determine common dispute resolution requirements

### Optional Requirements

1. **Resolution Support Tools** - Louis, Neil
   Interactive tools for mediators to help in generating settlement options, understanding dispute nuances, and guiding parties towards a resolution.
   - Knowledge base that covers all domains/cultures/entities.
2. **Resolution Archive** - Louis, Neil
   A secure, searchable archive of resolved disputes for auditing purposes and to serve as precedents where applicable.
   - Useful for Core Requirement 2.2
   - Summary generation for disputes
3. **Feedback System** - Louis, Neil
   Post-resolution feedback collection from all parties involved, including satisfaction ratings and suggestions for system improvement.
   - Feedback form generation per dispute
   - Profiling for  mediators(?) 
4. **Language Adaptability** - Louis
   Equip the platform with multilingual support to cater to a diverse user base. Implement translation services and local language processing to facilitate dispute resolutions across different linguistic backgrounds.
   - Google Translate API(?)
   - Mediators profiling could contain language metadata

### Wow Factors

1. Integration with blockchain for immutable record-keeping.
2. Cultural Sensitivity Algorithms. - Neil
    > Develop algorithms that can adapt resolution recommendations based on cultural >nuances and legal differences across regions.
3. Advanced analytics for dispute resolution insights (using the Resolution Archive and external data). - Mark
4. Customizable dispute resolution workflows.
5. AI-driven mediation suggestions.
