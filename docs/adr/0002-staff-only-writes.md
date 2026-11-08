# Restrict shipment and hub writes to trusted staff

Shipment and hub records have no customer ownership relationship. Public signup must not assign staff privileges, and only trusted staff may create, change, or delete shipments and hubs. Public reads remain available; registered customers and the visible demo account are read-only. Allowing any registered account to write would expose every record, while customer-owned writes would require a separate ownership model.

The backend seeds an admin using deployment-provided credentials and a separate read-only demo account whose login is visible in the frontend. Preserve existing accounts and shipments, but demote all pre-existing admin accounts during the one-time production cutover; only the newly seeded admin retains write access.
