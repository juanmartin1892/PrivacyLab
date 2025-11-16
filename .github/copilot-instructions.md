# GitHub Copilot instructions for this repository

## General style

- Do not add any emojis anywhere (not in code, comments, commit messages, or documentation).
- Write all comments, explanations, and documentation in clear, concise English.
- Prefer explicit, readable code over “clever” one-liners.
- Use meaningful, descriptive names for variables, functions, types, and modules.
- Keep functions small and focused on a single responsibility.
- When proposing changes, preserve the existing architecture and conventions of this repository.

## Project context: PrivacyLab

- This repository belongs to the PrivacyLab project, which focuses on privacy-preserving systems and secure data processing.
- Assume the codebase may include:
  - Cryptographic primitives and protocols.
  - Homomorphic encryption or related privacy-preserving techniques.
  - Backend services, APIs, and infrastructure components.
- When suggesting designs, always consider:
  - Data minimization.
  - Defense in depth.
  - Clear separation between trusted and untrusted components.

## Security and correctness

- Never introduce example secrets, hard-coded credentials, or test keys that look realistic.
- Do not suggest disabling security features “just for testing” (for example, bypassing TLS verification, removing authentication, or weakening cipher suites).
- Prefer explicit input validation and error handling.
- When touching cryptographic code:
  - Do not invent new cryptographic schemes.
  - Prefer well-established libraries and primitives already used in the project.
  - Avoid patterns that can lead to side-channel leaks.
- If you are not sure about a cryptographic detail, say so explicitly and suggest consulting official documentation or domain experts.

## Code and documentation

- Match the project’s existing language, tooling, and frameworks.
- When generating code:
  - Add brief, meaningful comments only where they improve understanding.
  - Avoid redundant comments that merely restate the code.
- When generating documentation (for example, README sections, design documents, or guides):
  - Start with a short summary of the purpose and scope.
  - Use clear headings and bullet points to structure the content.
  - Explain assumptions, limitations, and security considerations.
  - If there is a readme.md create a resume and add it to the readme.md file.

## Tests and examples

- Whenever you generate or modify non-trivial code, propose appropriate tests:
  - Unit tests for core logic.
  - Integration tests where interactions between components are important.
- Prefer deterministic tests that are easy to run in CI.
- When providing example code:
  - Make it minimal but complete enough to compile and run.
  - Highlight any configuration or environment requirements.

## Explanations and reviews

- When the user asks for an explanation:
  - Explain step by step, using precise but accessible technical language.
  - Prefer short paragraphs over long, dense blocks of text.
- When reviewing or refactoring code:
  - Focus on correctness, security, readability, and maintainability.
  - Clearly separate “must fix” issues from optional improvements.

## Things to avoid

- Do not use emojis or informal chat-style language in any generated content.
- Do not generate code that relies on undocumented behavior or hidden side effects.
- Do not suggest manual steps that contradict existing automation or CI/CD workflows in this repository.