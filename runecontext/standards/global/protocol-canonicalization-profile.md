---
schema_version: 1
id: global/protocol-canonicalization-profile
title: Protocol Canonicalization Profile
status: active
suggested_context_bundles:
    - protocol-foundation
---

# Protocol Canonicalization Profile

Use RFC 8785 JCS for hashed or signed protocol objects.

- Sign the detached payload's RFC 8785 JCS bytes, not wrapper fields or language-native serialization
- Accept RFC 8785-compatible JSON numbers, including decimal and exponent forms, and normalize them through JCS/ECMAScript number formatting rules
- Normalize `-0` to `0`
- Sort object keys by UTF-16 code units as required by RFC 8785 JCS, including non-ASCII keys
- Trusted persisted and signed RuneCode protocol surfaces are currently object-rooted; object and array roots are the supported canonicalized top-level forms in current wrappers
- Prove parity with checked-in golden canonical JSON and hash fixtures
- Keep Go and JS fixture checks aligned on the same canonical byte contract
