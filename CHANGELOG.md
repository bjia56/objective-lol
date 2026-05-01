# Changelog

## [0.0.3] - 2026-04-30

### Changed
- Class member variable type annotations are now optional, enabling dynamic typing consistent with global and local variable declarations

### Fixed
- Member variable getter errors are now propagated as interpreter exceptions instead of being silently swallowed
- Setter errors from native (foreign-language bridge) implementations are now propagated as runtime exceptions rather than silently ignored
- Error messages for failed function and method calls now include the function/method name

## [0.0.2] - 2025-09-16

### Changed
- Modified octal prefix from `0` to `0o`

## [0.0.1] - 2025-09-14

### Added
- Initial release of Objective-LOL interpreter and LSP server