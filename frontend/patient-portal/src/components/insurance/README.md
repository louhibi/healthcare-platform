# Insurance Components

This folder groups all insurance-related form components:

- `InsuranceFieldGroup.vue` – orchestrates insurance type/provider and related fields, emits aggregated payload.
- `InsuranceTypeField.vue` – select + optional custom input for insurance type.
- `InsuranceProviderField.vue` – select + optional custom input for provider, dependent on selected type.

Imports have been updated to use relative paths from this directory. Update any dynamic import loaders accordingly if paths were previously different.
