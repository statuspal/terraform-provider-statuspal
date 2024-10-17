# Changelog

All notable changes of the StatusPal Terraform provider will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.12] - 2024-10-17

### Fixed

- Error when trying to change an attribute from "something" to `null` or empty string.

### Removed

- `theme_configs` attribute from status page resource and data source, as it was deprecated.

### Added

- `zoom_notifications_enabled` and `allowed_email_domains` attributes to status page resource and data source.
- `monitoring_options` attribute to service resource and data source.

## [0.2.11] - 2024-10-07

### Fixed

- Issue with `terraform apply` showing the services will change when they actually don't.

## [0.2.10] - 2024-09-18

### Fixed

- Exclude omitting empty value from boolean attributes.

## [0.2.9] - 2024-09-09

### Fixed

- Omit empty values from the API request payload except the service `monitoring` attribute.

## [0.2.8] - 2024-07-19

### Fixed

- When an attribute is empty, `null` or omitted the empty value is not saved.

## [0.2.7] - 2024-07-10

### Changed

- Update dependecies.

## [0.2.6] - 2024-07-10

### Added

- `webhook_monitoring_service`, `webhook_custom_jsonpath_settings`, `inbound_email_address` and `incoming_webhook_url` attributes to the service resource and services data source.
- `webhook` type to the service `monitoring` attribute
- `bg_image`, `logo` and `favicon` readonly attributes to the status_page resource and status_pages data source.
- Validators to the service resource attributes where it was needed.

### Changed

- Polished the documentation.

## [0.2.5] - 2024-06-14

### Changed

- The `parent_id` attribute type of the `service` resource from `int64` to `string`.

## [0.2.4] - 2024-06-13

### Changed

- Better documentation.

## [0.2.3] - 2024-06-07

### Changed

- Better documentation.

## [0.2.2] - 2024-06-07

### Added

- Validators to schema attributes where it was needed.

### Changed

- Better documentation.

## [0.2.1] - 2024-06-05

### Changed

- Better documentation.

## [0.2.0] - 2024-06-05

### Added

- `services` data source and `service` resource to manage [StatusPal services](https://docs.statuspal.io/platform/services-components).

### Changed

- [statuspal Provider](https://registry.terraform.io/providers/statuspal/statuspal/latest/docs) documentation.

## [0.1.0] - 2024-06-03

### Added

- `status_pages` data source and `status_page` resource to manage [StatusPal status pages](https://www.statuspal.io/features/status-page).
