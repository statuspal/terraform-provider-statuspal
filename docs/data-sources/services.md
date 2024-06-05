---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "statuspal_services Data Source - statuspal"
subcategory: ""
description: |-
  Fetches the list of services in the status page.
---

# statuspal_services (Data Source)

Fetches the list of services in the status page.

## Example Usage

```terraform
# List all services of the status page with subdomain "example-com".
data "statuspal_services" "all" {
  status_page_subdomain = "example-com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `status_page_subdomain` (String) The status page subdomain of the services.

### Read-Only

- `id` (String) Placeholder identifier attribute. Ignore it, only used in testing.
- `services` (Attributes List) List of services. (see [below for nested schema](#nestedatt--services))

<a id="nestedatt--services"></a>
### Nested Schema for `services`

Read-Only:

- `auto_incident` (Boolean) Create an incident automatically when this service is down and close it if/when it comes back up.
- `auto_notify` (Boolean) Automatically notify all your subscribers about automatically created and closed incidents.
- `children_ids` (List of Number) IDs of the service's children.
- `current_incident_type` (String) Enum: `"major"` `"minor"` `"scheduled"`
The type of the (current) incident:
  - `major` - A minor incident is currently taking place.
  - `minor` - A major incident is currently taking place.
  - `scheduled` - A scheduled maintenance is currently taking place.
- `description` (String) The description of the service.
- `display_response_time_chart` (Boolean) Display response time chart?
- `display_uptime_graph` (Boolean) Display uptime graph?
- `id` (Number) The ID of the service.
- `inbound_email_id` (String) The inbound email ID.
- `incident_type` (String) Enum: `"major"` `"minor"`
The type of the (current) incident:
  - `major` - A minor incident is currently taking place.
  - `minor` - A major incident is currently taking place.
- `inserted_at` (String) Datetime at which the service was inserted.
- `is_up` (Boolean) Is the monitored service up?
- `monitoring` (String) Enum: `null` `"internal"` `"3rd_party"`
Monitoring types:
  - `major` - No monitoring.
  - `internal` - StatusPal monitoring.
  - `3rd_party` - 3rd Party monitoring.
- `name` (String) The name of the service.
- `order` (Number) Service's position in the service list.
- `parent_incident_type` (String) Enum: `"major"` `"minor"`
The type of the (current) incident:
  - `major` - A minor incident is currently taking place.
  - `minor` - A major incident is currently taking place.
- `pause_monitoring_during_maintenances` (Boolean) Pause the the service monitoring during maintenances?
- `ping_url` (String) We will send HTTP requests to this URL for monitoring every minute.
- `private` (Boolean) Private service?
- `private_description` (String) The private description of the service.
- `translations` (Attributes Map) A translations object. For example:
  ```terraform
	{
		en = {
			name = "Your service"
			description = "This is your service's description..."
		}
		fr = {
			name = "Votre service"
			description = "Voici la description de votre service..."
		}
	}
  ```
→ (see [below for nested schema](#nestedatt--services--translations))
- `updated_at` (String) Datetime at which the service was last updated.

<a id="nestedatt--services--translations"></a>
### Nested Schema for `services.translations`

Read-Only:

- `description` (String) The description of the service.
- `name` (String) The name of the service.