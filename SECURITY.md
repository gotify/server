# Security Policy

## Supported Versions

Only the latest released version.

If you found a vulnerability that only applies to older versions but has been accidentally fixed recently, please open a private advisory to let us evaluate if a backdated advisory is necessary.

If you found a vulnerability in unreleased code (Git trunk), please verify that the latest release is not affected and then use the public issue and pull request workflow to submit your research.

## Reporting a Vulnerability

Please report (suspected) security vulnerabilities to
[GitHub Advisory](https://github.com/gotify/server/security/advisories/new)
or **[gotify@protonmail.com](mailto:gotify@protonmail.com)**.
You will receive a response from us within a few days.

To reduce paperwork and align with CVE key details phrasing,
an executive summary containing the following elements is sufficient for most reports:

- The affected component (package, file, function, etc)
- The root cause (weakness in code, insecure default, misleading documentation, etc)
- The attack model (precondition, vector, impact)
- A PoC

If the issue is confirmed, we will release a
patch as soon as possible.
Additionally, we will submit findings that demonstrate the necessity for
user triage to the GitHub CNA Program.
